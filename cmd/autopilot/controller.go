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
	"context"
	"reflect"
	"sync"
	"time"

	autopilot "github.com/libopenstorage/autopilot/pkg/apis/autopilot"
	autopilotv1 "github.com/libopenstorage/autopilot/pkg/apis/autopilot/v1alpha1"
	"github.com/libopenstorage/stork/pkg/controller"
	log "github.com/sirupsen/logrus"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	resyncPeriod = 30 * time.Second
)

// crdController is the k8s controller interface for autopilot resources
type crdController struct {
	storagePolicies map[string]*autopilotv1.StoragePolicy
	spLock          sync.RWMutex
}

// Handle updates for StoragePolicy objects
func (c *crdController) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *autopilotv1.StoragePolicy:
		c.spLock.Lock()
		defer c.spLock.Unlock()

		if event.Deleted {
			delete(c.storagePolicies, o.Name)
			log.Infof("policy %s/%s/%s deleted", o.APIVersion, o.Kind, o.Name)
		} else {
			if tmp, ok := c.storagePolicies[o.Name]; !ok {
				c.storagePolicies[o.Name] = o
				log.Infof("policy %s/%s/%s added", o.APIVersion, o.Kind, o.Name)
			} else if tmp.GetResourceVersion() != o.GetResourceVersion() {
				c.storagePolicies[o.Name] = o
				log.Infof("policy %s/%s/%s updated", o.APIVersion, o.Kind, o.Name)
			}
		}
	}
	return nil
}

func newController() *crdController {
	return &crdController{
		storagePolicies: make(map[string]*autopilotv1.StoragePolicy),
	}
}

func (c *crdController) start() error {
	if err := controller.Init(); err != nil {
		return err
	}

	if err := controller.Register(
		&schema.GroupVersionKind{
			Group:   autopilot.GroupName,
			Version: autopilot.Version,
			Kind:    reflect.TypeOf(autopilotv1.StoragePolicy{}).Name(),
		},
		"",
		resyncPeriod,
		c); err != nil {
		return err
	}

	return controller.Run()
}
