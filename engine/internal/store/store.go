// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package store is the internal storage formats for autopilot
// TODO: this is a crappy name
package store

import (
	"github.com/libopenstorage/autopilot/telemetry"
)

type (
	// Writer defines an autopilot storage writer interface
	Writer interface {
		Write(vectors []telemetry.Vector) error
	}
)
