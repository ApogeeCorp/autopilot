/*
Copyright 2019 Openstorage.org

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kubernetes/kubernetes/pkg/api/legacyscheme"
	"github.com/libopenstorage/autopilot/config"
	"github.com/libopenstorage/autopilot/metrics"
	_ "github.com/libopenstorage/autopilot/metrics/providers"
	autopilot "github.com/libopenstorage/autopilot/pkg/apis/autopilot/v1alpha1"
	"github.com/libopenstorage/autopilot/pkg/log"
	"github.com/libopenstorage/autopilot/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	sparks "gitlab.com/ModelRocket/sparks/types"
	api_v1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

const (
	eventComponentName = "autopilot"
)

func main() {
	app := cli.NewApp()

	app.Name = "autopilot"
	app.Version = version.Version
	app.Usage = "Autopilot Storage Optimization Engine"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config,f",
			Usage:  "set the configuration file path",
			EnvVar: "CONFIG_FILE",
			Value:  "/etc/autopilot/config.yaml",
		},
		cli.StringFlag{
			Name:   "data-dir,d",
			Usage:  "set the data directory for the process",
			EnvVar: "DATA_DIR",
			Value:  "/var/run/autopilot",
		},
		cli.StringFlag{
			Name:   "log-level",
			Usage:  "set the log level",
			EnvVar: "LOG_LEVEL",
			Value:  "info",
		},
		cli.StringFlag{
			Name:   "log-format",
			Usage:  "set the log format",
			EnvVar: "LOG_FORMAT",
			Value:  "text",
		},
		cli.StringFlag{
			Name:   "kube-config",
			Usage:  "set the kubernetes config path",
			EnvVar: "KUBECONFIG",
		},
		cli.StringFlag{
			Name:   "kube-master-url",
			Usage:  "set the kubernetes master url",
			EnvVar: "KUBERNETES_MASTER_URL",
		},
	}

	app.Before = setupLog

	app.Action = func(c *cli.Context) error {
		logrus.Infof("Starting autopilot version: %s", app.Version)
		var shutdown = make(chan os.Signal, 1)

		cfg, err := config.ReadFile(c.GlobalString("config"))
		if err != nil {
			return err
		}

		signal.Notify(shutdown, syscall.SIGTERM)
		signal.Notify(shutdown, syscall.SIGINT)

		config, err := rest.InClusterConfig()
		if err != nil {
			logrus.Fatalf("Error getting cluster config: %v", err)
		}

		k8sClient, err := clientset.NewForConfig(config)
		if err != nil {
			logrus.Fatalf("Error getting client, %v", err)
		}

		eventBroadcaster := record.NewBroadcaster()
		eventBroadcaster.StartRecordingToSink(&v1.EventSinkImpl{Interface: v1.New(k8sClient.CoreV1().RESTClient()).Events("")})
		recorder := eventBroadcaster.NewRecorder(legacyscheme.Scheme, api_v1.EventSource{Component: eventComponentName})

		// install our CRD
		if err := crdInstallAction(c); err != nil {
			return err
		}

		controller := newController(recorder, cfg)

		// start the controller
		if err := controller.start(); err != nil {
			return err
		}

		pollRate, err := time.ParseDuration(cfg.PollRate)
		if err != nil {
			return err
		}

		logrus.Infof("starting the metrics poller (%s)", cfg.PollRate)

		ticker := sparks.NewTicker(pollRate)

		provs := make(map[string]metrics.Provider)

		for _, prov := range cfg.Providers {
			inst, err := metrics.NewProvider(prov.Type, prov.Params)
			if err != nil {
				return err
			}

			provs[prov.Name] = inst
		}

		for {
			select {
			case <-ticker.C():
				controller.lock()

				for name, prov := range provs {
					for _, pol := range controller.storagePolicies {
						vecs, err := prov.Query(pol)
						if err != nil {
							log.StoragePolicyLog(pol).Errorln(err)
							continue
						}

						if len(vecs) == 0 {
							log.StoragePolicyLog(pol).Debugf("no vectors matched")
							break
						}

						log.StoragePolicyLog(pol).Debugf("has %d match(es) on provider %s", len(vecs), name)
						objects, err := getObjectsForPolicy(pol)
						if err != nil {
							log.StoragePolicyLog(pol).Errorf(err.Error())
							return err
						}

						for _, object := range objects {
							if isConditionMetOnObject(object, vecs) {
								// TODO improve this
								conditionStr := ""
								for i, cond := range pol.Spec.Conditions {
									conditionStr = conditionStr + fmt.Sprintf("%d => %s %s %s\t", i+1, cond.Key, cond.Operator, cond.Values)
								}
								controller.recorder.Event(pol,
									api_v1.EventTypeNormal,
									string(autopilot.StoragePolicyConditonMet),
									fmt.Sprintf("conditions: %s met on object: %s",
										conditionStr, object))

								if controller.isObjectInCoolDown(object) {
									continue
								}

								if err := controller.executePolicyAction(pol, object); err != nil {
									log.StoragePolicyLog(pol).Errorln(err)
									controller.recorder.Event(pol,
										api_v1.EventTypeWarning,
										string(autopilot.StoragePolicyActionFailed),
										err.Error())
									return err
								}

								if err := controller.markObjectForCoolDown(object); err != nil {
									log.StoragePolicyLog(pol).Errorln(err)
									controller.recorder.Event(pol,
										api_v1.EventTypeWarning,
										string(autopilot.StoragePolicyActionFailed),
										err.Error())
									return err
								}
							} else {
								log.StoragePolicyLog(pol).Debugf("condition not met for object: %v", object)
							}
						}
					}
				}
				controller.unlock()
				ticker.Reset()

			case <-shutdown:
				logrus.Infof("shutting down")
				return nil
			}
		}
	}

	app.Commands = []cli.Command{
		{
			Name:  "crd",
			Usage: "Manage auto-pilot CRDs",
			Subcommands: []cli.Command{
				{
					Name:   "install",
					Action: crdInstallAction,
					Usage:  "publish the autopilot crds to the k8s cluster",
				},
			},
		},
		{
			Name:  "policy",
			Usage: "Manage auto-pilot policy objects",
			Subcommands: []cli.Command{
				{
					Name:      "test",
					Action:    policyTestAction,
					Usage:     "Test a policy document using the configuration",
					UsageText: "test <file>",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func setupLog(c *cli.Context) error {
	// setup the log format
	switch c.String("log-format") {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	case "json":
		fallthrough
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}

	// setup the log level
	if level, err := logrus.ParseLevel(c.String("log-level")); err == nil {
		logrus.SetLevel(level)
	}

	return nil
}
