// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package telemetry

import (
	"github.com/libopenstorage/autopilot/api/autopilot/types"
	sparks "gitlab.com/ModelRocket/sparks/types"
)

// The autopilot telemetry format is based on the prometheus metrics data format

type (
	// Provider defines a simple interface for telemetry providers to collect and extract data
	Provider interface {
		// Query returns a results vector from a direct query to the provider
		Query(params Params) ([]Vector, error)

		// Parse returns a result vector from the raw data
		Parse(data []byte) ([]Vector, error)
	}

	// Params is an alias for a map helper
	Params = sparks.Params

	// Attributes is an alias for a map helper
	Attributes = sparks.Params

	// NewFunc is a function registered with the telemetry layer for creating a new
	// instance of the provider.
	NewFunc func(types.Provider) (Provider, error)

	// Metric is the metric that comes from Prometheus apiproxy
	Metric struct {
		Name          string  `json:"__name__"`
		Cluster       string  `json:"cluster"`
		Instance      string  `json:"instance"`
		Node          string  `json:"node_id"`
		NodeName      string  `json:"node"`
		Job           *string `json:"job,omitempty"`
		Volume        *string `json:"volumeid,omitempty"`
		VolumeName    *string `json:"volumename,omitempty"`
		VolumePVC     *string `json:"pvc,omitempty"`
		Namespace     *string `json:"namespace,omitempty"`
		Disk          *string `json:"disk,omitempty"`
		Pool          *string `json:"pool,omitempty"`
		AlertName     *string `json:"alertname,omitempty"`
		AlertState    *string `json:"alertstate,omitempty"`
		AlertSeverity *string `json:"severity,omitempty"`
		AlertIssue    *string `json:"issue,omitempty"`
		Proc          *string `json:"proc,omitempty"`
	}

	// Vector for a single metric and value
	Vector struct {
		Metric Metric `json:"metric"`
		// Value is [timestamp, value] such as "value": [1231231233.232, "0"], its mixed type
		Value []interface{} `json:"value,omitempty"`
		// Values is for range queries its an array of Value above
		Values [][]interface{} `json:"values,omitempty"`
	}
)
