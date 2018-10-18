// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package telemetry provides the interfaces for telemetry providers
package telemetry

// Provider is the telemetry provider interface
type Provider interface {
	Name() string
	Init(params map[string]string) error
}
