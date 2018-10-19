// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package config provides the configuration structures used in the autopilot
package config

// Query defines a range query
type Query struct {
	// Expression is the query expression passed to prometheus
	Expression string `yaml:"expr"`

	// Step is the query resolution
	Step string `yaml:"step"`
}

type Config struct {
	// URL is the prometheus host url
	URL string `yaml:"url"`

	// Queries is the list of configured queries to perform
	Queries []Query `yaml:"queries"`

	// ModelPath is location of classification models
	ModelPath string `yaml:"model_path"`

	// OutPath is the location to output recommendations
	OutPath string `yaml:"out_path"`
}
