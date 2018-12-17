// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package store

import "github.com/libopenstorage/autopilot/telemetry"

const (
	// Cluster stats
	Cluster = "cluster"
	// Node stats
	Node = "node"
	// Disk stats
	Disk = "disk"
	// Volume stats
	Volume = "volume"
	// Pool stats
	Pool = "pool"
	// Proc stats - on the PX processes
	Proc = "proc"
	// Alerts metrics are a special case for outputs
	Alerts = "ALERTS"
	// Instance field in the CSV
	Instance = "instance"
	// Timestamp field in the CSV
	Timestamp = "timestamp"
)

type (
	// Structures for Underlying Representation for ML

	// Row is the common characteristics of every metric - this defines a row in the CSV
	Row struct {
		Timestamp uint32
		Cluster   string
		Instance  string
		Node      string
	}

	// AlertRow is a special type of metric that we use for the Output of the ML.  This is one of the predictors
	AlertRow struct {
		Row
		AlertName     string
		AlertState    string
		AlertSeverity string
		AlertIssue    string
		AlertValue    string
	}

	// Metrics contain ALL the BaseAttributes for a row in the CSV Each map is keyed by the Name of the field (VolumeID, Disk, Pool, Proc)
	// This probably could contain a map of maps of maps, but its becoming unreadable with that level of nesting.
	// Cluster and Proc metrics each are node specific, so they will be located in the Node map
	Metrics struct {
		// Node are the node based metrics keyed on node field
		Node map[string]map[string]string
		// Volume are the volume based metrics keyed on Volume field
		Volume map[string]map[string]string
		// Disk are the disk based metrics keyed on Disk field
		Disk map[string]map[string]string
		// Pool are the pool based metrics keyed on Pool field
		Pool map[string]map[string]string
	}
)

func newAlertRow(row Row, vector telemetry.Vector, value string) *AlertRow {
	return &AlertRow{
		Row:           row,
		AlertName:     *vector.Metric.AlertName,
		AlertState:    *vector.Metric.AlertState,
		AlertSeverity: *vector.Metric.AlertSeverity,
		AlertIssue:    *vector.Metric.AlertIssue,
		AlertValue:    value,
	}
}
