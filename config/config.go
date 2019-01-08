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
	Providers  []types.Provider   `json:"providers"`
	Collectors []*types.Collector `json:"collectors"`
	Monitors   []*types.Monitor   `json:"monitors"`
	Emitters   []*types.Emitter   `json:"emitters"`
	Rules      []*types.Rule      `json:"rules"`
	DataDir    string             `json:"data_dir"`
	Listen     string             `json:"listen"`
}

// UnmarshalJSON handles the unmarshaling of the config
func (m *Config) UnmarshalJSON(raw []byte) error {
	var data struct {
		Providers  json.RawMessage    `json:"providers"`
		Collectors []*types.Collector `json:"collectors"`
		Monitors   []*types.Monitor   `json:"monitors"`
		Emitters   []*types.Emitter   `json:"emitters"`
		Rules      []*types.Rule      `json:"rules"`
		Listen     string             `json:"listen"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)

	if err := dec.Decode(&data); err != nil {
		return err
	}

	m.Collectors = data.Collectors
	m.Monitors = data.Monitors
	m.Emitters = data.Emitters
	m.Rules = data.Rules

	rdr := bytes.NewReader([]byte(data.Providers))
	provs, err := types.UnmarshalProviderSlice(rdr, runtime.JSONConsumer())
	if err != nil {
		return err
	}

	m.Providers = provs

	return nil
}

// GetRule returns a config rule by name
func (m *Config) GetRule(name string) (*types.Rule, bool) {
	for _, r := range m.Rules {
		if r.Name == name {
			return r, true
		}
	}

	return nil, false
}
