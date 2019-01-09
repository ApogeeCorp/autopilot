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
	clientset "github.com/libopenstorage/autopilot/pkg/client/clientset/versioned"
	listers "github.com/libopenstorage/autopilot/pkg/client/listers/autopilot/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

// Controller is the k8s controller interface for autopilot resources
type Controller struct {
	kubeclientset kubernetes.Interface
	apclientset   clientset.Interface

	storagePolicyLister listers.StoragePolicyLister
	storagePolicySynced cache.InformerSynced
	workqueue           workqueue.RateLimitingInterface
	recorder            record.EventRecorder
}
