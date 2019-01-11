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

	"github.com/libopenstorage/autopilot/config"
	"github.com/libopenstorage/autopilot/metrics"
	_ "github.com/libopenstorage/autopilot/metrics/providers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	sparks "gitlab.com/ModelRocket/sparks/types"
)

func main() {
	app := cli.NewApp()

	app.Name = "autopilot"
	app.Version = "0.0.1"
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
		var shutdown = make(chan os.Signal)

		cfg, err := config.ReadFile(c.GlobalString("config"))
		if err != nil {
			return err
		}

		signal.Notify(shutdown, syscall.SIGTERM)
		signal.Notify(shutdown, syscall.SIGINT)

		// install our CRD
		if err := crdInstallAction(c); err != nil {
			return err
		}

		controller := newController()

		// start the controller
		if err := controller.start(); err != nil {
			return err
		}

		pollRate, err := time.ParseDuration(cfg.PollRate)
		if err != nil {
			return err
		}

		log.Infof("starting the metrics poller (%s)", cfg.PollRate)

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
				log.Debug("beginning metrics polling")

				controller.lock()

				for name, prov := range provs {
					for _, pol := range controller.storagePolicies {
						log.Debugf("querying provider %s for policy %s", name, pol.Name)
						vecs, err := prov.Query(pol)
						if err != nil {
							log.Errorln(err)
							continue
						}

						log.Debugf("policy %s has %d match(es) on provider %s", pol.Name, len(vecs), name)

						if err := executePolicy(pol, vecs); err != nil {
							log.Errorln(err)
						}
					}
				}
				controller.unlock()

				log.Debug("metrics polling done")

				ticker.Reset()

			case <-shutdown:
				log.Infof("shutting down")
				return nil
			}
		}
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "crd",
			Usage: "Manage auto-pilot CRDs",
			Subcommands: []cli.Command{
				cli.Command{
					Name:   "install",
					Action: crdInstallAction,
					Usage:  "publish the autopilot crds to the k8s cluster",
				},
			},
		},
		cli.Command{
			Name:  "policy",
			Usage: "Manage auto-pilot policy objects",
			Subcommands: []cli.Command{
				cli.Command{
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
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	case "json":
		fallthrough
	default:
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}

	// setup the log level
	if level, err := log.ParseLevel(c.String("log-level")); err == nil {
		log.SetLevel(level)
	}

	return nil
}
