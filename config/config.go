// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package config

import "github.com/libopenstorage/autopilot/api/autopilot/types"

// Config defines the autopilot configuration structure
type Config struct {
	Rules      []*types.Rule      `yaml:"rules"`
	Collectors []*types.Collector `yaml:"collectors"`
	Emitters   []*types.Emitter   `yaml:"emitters"`
}
