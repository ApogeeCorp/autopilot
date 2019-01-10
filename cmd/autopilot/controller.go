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
type crdController struct{}

// Handle updates for StoragePolicy objects
func (c *crdController) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *autopilotv1.StoragePolicy:
		spLock.Lock()
		defer spLock.Unlock()

		if event.Deleted {
			log.Debugf("storage policicy %q delete", o.Name)
			delete(storagePolicies, o.Name)
		} else {
			log.Debugf("storage policy %q update", o.Name)
			storagePolicies[o.Name] = o
		}
	}
	return nil
}

func startController() error {
	ctl := &crdController{}

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
		ctl); err != nil {
		return err
	}

	return controller.Run()
}
