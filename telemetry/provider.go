// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package telemetry

import (
	"github.com/libopenstorage/autopilot/api/autopilot/types"
	sparks "gitlab.com/ModelRocket/sparks/types"
)

type (
	// Provider defines a simple interface for telemetry providers to collect and extract data
	Provider interface {
		// Collect executes a query on the provider and collects the data to the staging path
		Collect(host string, params Params, outPath string) error
		// Query executes a query directly on the provider and returns the data in its native format
		Query(host string, rule *types.Rule) (*types.Recommendation, error)
	}

	// Params is an alias for a map helper
	Params = sparks.Params
)
