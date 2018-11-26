// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"encoding/json"

	"github.com/go-openapi/runtime"
	"github.com/libopenstorage/autopilot/api/autopilot/types"
)

// Config defines the autopilot configuration structure
type Config struct {
	Rules         []*types.Rule    `json:"rules"`
	CollectorsRaw json.RawMessage  `json:"collectors"`
	Emitters      []*types.Emitter `json:"emitters"`
	collectors    []types.Collector
}

// Collectors unmarshals the typed collectors properly from the json
func (c *Config) Collectors() []types.Collector {
	return c.collectors
}

// UnmarshalJSON handles the custom unmarshaling of the config
func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	collectors, err := types.UnmarshalCollectorSlice(bytes.NewReader(c.CollectorsRaw), runtime.JSONConsumer())
	if err != nil {
		return err

	}
	c.collectors = collectors

	return nil
}
