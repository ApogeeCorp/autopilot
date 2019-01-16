package log

import (
	autopilot "github.com/libopenstorage/autopilot/pkg/apis/autopilot/v1alpha1"
	"github.com/sirupsen/logrus"
)

// StoragePolicyLog Format a log message with storage pilicy information
func StoragePolicyLog(policy *autopilot.StoragePolicy) *logrus.Entry {
	if policy != nil {
		fields := logrus.Fields{
			"Name":      policy.Name,
			"Namespace": policy.Namespace,
		}
		return logrus.WithFields(fields)
	}
	return logrus.WithFields(logrus.Fields{
		"StoragePolicy": policy,
	})
}
