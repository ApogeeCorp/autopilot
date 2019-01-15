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
	"errors"
	"io/ioutil"

	"github.com/libopenstorage/autopilot/metrics"
	log "github.com/sirupsen/logrus"

	"github.com/libopenstorage/autopilot/config"
	autopilot "github.com/libopenstorage/autopilot/pkg/apis/autopilot/v1alpha1"
	"github.com/urfave/cli"
	sparks "gitlab.com/ModelRocket/sparks/types"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

func policyTestAction(c *cli.Context) error {
	cfg, err := config.ReadFile(c.GlobalString("config"))
	if err != nil {
		return err
	}

	if c.NArg() < 1 {
		return errors.New("missing policy document path")
	}

	path := c.Args().Get(0)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	scheme := runtime.NewScheme()
	if err = autopilot.AddToScheme(scheme); err != nil {
		return err
	}

	deserializer := serializer.NewCodecFactory(scheme).UniversalDeserializer()
	obj, _, err := deserializer.Decode(data, nil, nil)
	if err != nil {
		return err
	}

	policy, ok := obj.(*autopilot.StoragePolicy)
	if !ok {
		return errors.New("invalid storage policy object")
	}

	for _, p := range cfg.Providers {
		prov, err := metrics.NewProvider(p.Type, p.Params)
		if err != nil {
			return err
		}

		vecs, err := prov.Query(policy)
		if err != nil {
			return err
		}

		if len(vecs) != 0 {
			return executePolicy(policy, vecs)
		}
		log.Infof("no policies matched")
	}

	return nil
}

func executePolicy(policy *autopilot.StoragePolicy, vecs []metrics.Vector) error {
	for _, exp := range policy.Spec.Object.MatchExpressions {
		values := sparks.Slice(&exp.Values)
		for _, vec := range vecs {
			switch policy.Spec.Object.Type {
			case "openstorage.io/object.volume":
				if values.Contains(*vec.Metric.VolumeName) {
					log.Infof("should execute action %s on volume %s", policy.Spec.Action.Name, *vec.Metric.VolumeName)

					// TODO: blacklist actionable policies
				}
			}
		}
	}

	return nil
}
