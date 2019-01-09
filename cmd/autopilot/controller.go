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
	"github.com/libopenstorage/autopilot/pkg/controller"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	resyncPeriod = 30 * time.Second
)

// Controller is the k8s controller interface for autopilot resources
type Controller struct {
}

// Init Initialize the migration controller
func (c *Controller) Init() error {
	return controller.Register(
		&schema.GroupVersionKind{
			Group:   autopilot.GroupName,
			Version: autopilot.Version,
			Kind:    reflect.TypeOf(autopilotv1.StoragePolicy{}).Name(),
		},
		"",
		resyncPeriod,
		c)
}

// Handle updates for StoragePolicy objects
func (c *Controller) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *autopilotv1.StoragePolicy:
		log.Debugf("%s => %s (%v)", o.Kind, o.Spec.Object.Type, event.Deleted)
	}
	return nil
}
