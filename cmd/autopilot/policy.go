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
	"regexp"

	"github.com/libopenstorage/autopilot/metrics"
	"github.com/portworx/sched-ops/k8s"
	"github.com/sirupsen/logrus"

	autopilot "github.com/libopenstorage/autopilot/pkg/apis/autopilot/v1alpha1"
	"github.com/libopenstorage/autopilot/pkg/log"
	"github.com/urfave/cli"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var policyActionNameRegex = regexp.MustCompile(`^(.+)/(.+)`)

func policyTestAction(c *cli.Context) error {
	/*cfg, err := config.ReadFile(c.GlobalString("config"))
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

		objects, err := getObjectsForPolicy(policy)
		if err != nil {
			return err
		}

		for _, object := range objects {
			if isConditionMetOnObject(object, vecs) {
				// TODO offload this to a scheduled task
				if err := executePolicyAction(policy, object); err != nil {
					return err
				}
			}

			logrus.Debugf("[debug] condition not met for policy: %s for object: %v", policy.Name, object)
		}
	}*/

	return nil
}

func (c *crdController) isObjectInCoolDown(object string) bool {
	c.probationLock.Lock()
	defer c.probationLock.Unlock()

	_, present := c.objectsInProbation[object]
	return present
}

func (c *crdController) markObjectForCoolDown(object string) error {
	c.probationLock.Lock()
	defer c.probationLock.Unlock()

	if err := c.probation.Add(object, nil, true); err != nil {
		return err
	}

	c.objectsInProbation[object] = nil

	return nil
}

func (c *crdController) executePolicyAction(policy *autopilot.StoragePolicy, object string) error {
	logrus.Infof("should execute action %s on object %s", policy.Spec.Action.Name, object)
	actionObjectType, actionType := parseObjectTypeFromActionName(policy.Spec.Action.Name)

	if len(actionObjectType) == 0 {
		return fmt.Errorf("failed to get action object type for policy: %s", policy.Name)
	}

	if len(actionType) == 0 {
		return fmt.Errorf("failed to get action type for policy: %s", policy.Name)
	}

	log.StoragePolicyLog(policy).Infof("action type: %s, action object type: %s", actionType, actionObjectType)

	switch actionObjectType {
	case autopilot.PolicyActionVolume:
		log.StoragePolicyLog(policy).Debugf("running volume policy action")
		return c.executeVolumeAction(policy, actionType, object)
	default:
		err := fmt.Errorf("unsupported policy action: %s", policy.Spec.Action.Name)
		log.StoragePolicyLog(policy).Errorf(err.Error())
		return err
	}
}

func (c *crdController) executeVolumeAction(policy *autopilot.StoragePolicy, actionType string, volumeID string) error {
	switch actionType {
	case autopilot.PolicyActionVolumeResize:
		log.StoragePolicyLog(policy).Infof("Performing resize on vol: %s", volumeID)
		if err := c.resizeVolume(policy, volumeID); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported action: %s on volume: %s", actionType, volumeID)
	}

	c.recorder.Event(policy,
		v1.EventTypeNormal,
		string(autopilot.StoragePolicyActionTriggered),
		fmt.Sprintf("action: %s triggered successfully on volume: %s",
			actionType, volumeID))
	return nil
}

func (c *crdController) resizeVolume(policy *autopilot.StoragePolicy, volumeID string) error {
	pv, err := k8s.Instance().GetPersistentVolume(volumeID)
	if err != nil {
		return err
	}

	claimRef := pv.Spec.ClaimRef
	if claimRef == nil {
		return fmt.Errorf("failed to get PVC from PV as claim reference is nil")
	}

	pvcName := claimRef.Name
	pvcNamespace := claimRef.Namespace
	pvc, err := k8s.Instance().GetPersistentVolumeClaim(pvcName, pvcNamespace)
	if err != nil {
		return err
	}

	storageSize := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	// TODO resize by user given factor
	extraAmount, _ := resource.ParseQuantity("2Gi")
	storageSize.Add(extraAmount)
	pvc.Spec.Resources.Requests[v1.ResourceStorage] = storageSize

	_, err = k8s.Instance().UpdatePersistentVolumeClaim(pvc)
	if err != nil {
		return err
	}

	log.StoragePolicyLog(policy).Infof("successfully resized PVC: [%s] %s by %v for PV: %s",
		pvcNamespace, pvcName, extraAmount, volumeID)

	return nil

}

func parseObjectTypeFromActionName(actionName string) (string, string) {
	matches := policyActionNameRegex.FindStringSubmatch(actionName)
	if len(matches) == 3 {
		return matches[1], matches[2]
	}

	return "", ""
}

func getObjectsForPolicy(policy *autopilot.StoragePolicy) ([]string, error) {
	objects := make([]string, 0)

	switch policy.Spec.Object.Type {

	case autopilot.PolicyObjectTypeVolume:
		pvcs, err := k8s.Instance().GetPersistentVolumeClaims(policy.GetNamespace(), policy.Spec.Object.MatchLabels)
		if err != nil {
			return nil, err
		}

		for _, pvc := range pvcs.Items {
			pvName, err := k8s.Instance().GetVolumeForPersistentVolumeClaim(&pvc)
			if err != nil {
				return nil, err
			}

			objects = append(objects, pvName)
		}

	default:
		return nil, fmt.Errorf("unsupported object type: %s for policy", policy.Spec.Object.Type)
	}

	return objects, nil
}

func isConditionMetOnObject(object string, vecs []metrics.Vector) bool {
	if len(vecs) == 0 {
		return false
	}

	found := 0
	for _, vec := range vecs {
		// TODO can't assume volume type here
		if object == *vec.Metric.VolumeName {
			found = found + 1
			continue
		}
	}

	// all vecs need to have the object
	return found == len(vecs)
}
