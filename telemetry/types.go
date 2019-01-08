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

package telemetry

import (
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

	// NewFunc is a function registered with the telemetry layer for creating a new
	// instance of the provider.
	NewFunc func(Params) (Provider, error)

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
