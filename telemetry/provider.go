// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package telemetry

type (
	// Provider defines a simple interface for telemetry providers to collect and extract data
	Provider interface {
		// Collect executes a query on the provider and collects the data to the staging path
		Collect(params map[string]string, outPath string) error
	}
)
