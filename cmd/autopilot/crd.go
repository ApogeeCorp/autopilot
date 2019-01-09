package main

import (
	"reflect"
	"time"

	autopilot "github.com/libopenstorage/autopilot/pkg/apis/autopilot"
	autopilotv1 "github.com/libopenstorage/autopilot/pkg/apis/autopilot/v1alpha1"
	"github.com/portworx/sched-ops/k8s"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
)

const (
	validateCRDInterval time.Duration = 5 * time.Second
	validateCRDTimeout  time.Duration = 1 * time.Minute
)

func crdInstallAction(c *cli.Context) error {

	resource := k8s.CustomResource{
		Name:    autopilotv1.StoragePolicyResourceName,
		Plural:  autopilotv1.StoragePolicyResourcePlural,
		Group:   autopilot.GroupName,
		Version: autopilot.Version,
		Scope:   apiextensionsv1beta1.NamespaceScoped,
		Kind:    reflect.TypeOf(autopilotv1.StoragePolicy{}).Name(),
	}

	err := k8s.Instance().CreateCRD(resource)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	log.Debugf("CRD for %s created successfully", resource.Name)

	return k8s.Instance().ValidateCRD(resource, validateCRDTimeout, validateCRDInterval)
}
