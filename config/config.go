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

package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// MetricsProvider provides metrics data to autopilot
type MetricsProvider struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Params   string `yaml:"params"`
	PollRate string `yaml:"poll_rate"`
}

// Config defines the autopilot configuration structure
type Config struct {
	Providers []MetricsProvider `yaml:"providers"`
}

// ReadFile reads a configuration file
func ReadFile(f string) (*Config, error) {
	config := &Config{}

	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
