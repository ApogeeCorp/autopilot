// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package store

import (
	"encoding/csv"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/libopenstorage/autopilot/telemetry"
	"github.com/spf13/cast"
)

// Store is the internal autopilot storage writer
type Store struct {
	path string
	lock sync.Mutex
}

// NewStore returns the storage interface
func NewStore(path string) *Store {
	return &Store{
		path: path,
	}
}

// Write implements the store.Writer.Write
func (s *Store) Write(vectors []telemetry.Vector) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	/*timeseries, alerts := transformToRows(vectors)

	base := filepath.Join(stagingPath, startDate.Format("2006-01-02"), startDate.Format("1504"))
	if err := os.MkdirAll(base, 0770); err != nil {
		return err
	}

	writeCSV(timeseries, base, Volume)
	writeCSV(timeseries, base, Disk)
	writeCSV(timeseries, base, Pool)
	writeCSV(timeseries, base, Node)
	writeAlertCSV(base, alerts)
	*/
	return nil
}

// keys Return keys of the given map - used to pull the values from a particular VolumeID, Disk, Pool, Proc of a BaseMetric
func keys(m map[string]string) []string {
	rval := make([]string, 0)
	for k := range m {
		rval = append(rval, k)
	}
	sort.Strings(rval)
	return rval
}

// createCSVRows creates the rows for a particular MAP in the CSVMetrics
func createCSVRows(br Row, m map[string]map[string]string, label string) ([][]string, []string) {
	rows := make([][]string, 0)
	rowHeaders := make([]string, 0)

	for rowKey, rowValue := range m {
		var row = []string{strconv.FormatUint(uint64(br.Timestamp), 10), br.Cluster, br.Instance, br.Node}
		var rowHeader = []string{Timestamp, Cluster, Instance, Node}
		if label != Cluster && label != Node {
			row = append(row, rowKey)
			rowHeader = append(rowHeader, label)
		}
		for _, k := range keys(rowValue) {
			v := rowValue[k]
			row = append(row, v)
			rowHeader = append(rowHeader, k)
		}
		rows = append(rows, row)
		rowHeaders = rowHeader
	}

	return rows, rowHeaders
}

// writeCSV for a particular metric category, create the CSV file
func writeCSV(timeSeries map[Row]*Metrics, base, name string) error {
	base = path.Join(base, name)
	f, err := os.Create(base + ".csv")
	if err != nil {
		return err
	}

	defer f.Close()
	w := csv.NewWriter(f)
	var wroteHeader = false

	// We want to make sure the CSV files are in order of timestamp, so we need to sort the keys here
	keys := make([]Row, len(timeSeries))
	i := 0
	for k := range timeSeries {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Timestamp < keys[j].Timestamp
	})

	for _, br := range keys {
		bm := timeSeries[br]
		var m map[string]map[string]string
		if name == Volume {
			m = bm.Volume
		} else if name == Disk {
			m = bm.Disk
		} else if name == Pool {
			m = bm.Pool
		} else if name == Node {
			m = bm.Node
		}

		csvRows, csvHeader := createCSVRows(br, m, name)
		if wroteHeader == false {
			if err := w.Write(csvHeader); err != nil {
				return err
			}
			wroteHeader = true
		}
		for _, csvRow := range csvRows {
			if err := w.Write(csvRow); err != nil {
				return err
			}
		}
	}
	w.Flush()
	return nil
}

// writeAlertCSV for the alerts csv
func writeAlertCSV(base string, alerts []*AlertRow) error {
	base = path.Join(base, "alert.csv")
	f, err := os.Create(base)
	if err != nil {
		return err
	}

	defer f.Close()
	w := csv.NewWriter(f)
	if err := w.Write([]string{Timestamp, Cluster, Instance, Node, "px_alert_name", "px_alert_state", "px_alert_severity", "px_alert_value"}); err != nil {
		return err
	}
	for _, alert := range alerts {
		var row []string
		row = append(row, strconv.FormatUint(uint64(alert.Timestamp), 10), alert.Cluster, alert.Instance, alert.Node, alert.AlertName, alert.AlertState, alert.AlertSeverity, alert.AlertValue)
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return nil

}

// transformToRows takes the Prometheus API Calls Clustervecs and flattens it to the structure that can exported as CSV
func transformToRows(vectors []telemetry.Vector) (map[Row]*Metrics, []*AlertRow) {
	// Go through the vecs for the cluster and generate the appropriate CSV files
	timeseries := make(map[Row]*Metrics)
	alerts := make([]*AlertRow, 0)

	for _, vec := range vectors {
		values := vec.Values
		if vec.Values == nil {
			values = append(values, vec.Value)
		}
		for _, value := range values {
			csvRow := Row{
				Timestamp: cast.ToUint32(value[0]),
				Cluster:   vec.Metric.Cluster,
				Instance:  vec.Metric.Instance,
				Node:      vec.Metric.Node,
			}
			csvMetrics, ok := timeseries[csvRow]
			if !ok {
				csvMetrics = &Metrics{}
				timeseries[csvRow] = csvMetrics
			}
			if vec.Metric.Name == Alerts {
				alert := &AlertRow{
					Row:           csvRow,
					AlertName:     *vec.Metric.AlertName,
					AlertState:    *vec.Metric.AlertState,
					AlertSeverity: *vec.Metric.AlertSeverity,
					AlertIssue:    *vec.Metric.AlertIssue,
					AlertValue:    value[1].(string),
				}
				alerts = append(alerts, alert)
			} else if vec.Metric.Volume != nil {
				if csvMetrics.Volume == nil {
					csvMetrics.Volume = make(map[string]map[string]string)
				}
				if csvMetrics.Volume[*vec.Metric.Volume] == nil {
					csvMetrics.Volume[*vec.Metric.Volume] = make(map[string]string)
				}
				csvMetrics.Volume[*vec.Metric.Volume][vec.Metric.Name] = value[1].(string)
			} else if vec.Metric.Disk != nil {
				if csvMetrics.Disk == nil {
					csvMetrics.Disk = make(map[string]map[string]string)
				}
				if csvMetrics.Disk[*vec.Metric.Disk] == nil {
					csvMetrics.Disk[*vec.Metric.Disk] = make(map[string]string)
				}
				csvMetrics.Disk[*vec.Metric.Disk][vec.Metric.Name] = value[1].(string)
			} else if vec.Metric.Pool != nil {
				if csvMetrics.Pool == nil {
					csvMetrics.Pool = make(map[string]map[string]string)
				}
				if csvMetrics.Pool[*vec.Metric.Pool] == nil {
					csvMetrics.Pool[*vec.Metric.Pool] = make(map[string]string)
				}
				csvMetrics.Pool[*vec.Metric.Pool][vec.Metric.Name] = value[1].(string)
			} else if vec.Metric.Proc != nil {
				if csvMetrics.Node == nil {
					csvMetrics.Node = make(map[string]map[string]string)
				}
				if csvMetrics.Node[vec.Metric.Node] == nil {
					csvMetrics.Node[vec.Metric.Node] = make(map[string]string)
				}
				csvMetrics.Node[vec.Metric.Node][vec.Metric.Name+"_"+*vec.Metric.Proc] = value[1].(string)
			} else if strings.HasPrefix(vec.Metric.Name, "px_node_stats") == true || strings.HasPrefix(vec.Metric.Name, "px_network_") == true ||
				strings.HasPrefix(vec.Metric.Name, "px_cluster_") == true {
				if csvMetrics.Node == nil {
					csvMetrics.Node = make(map[string]map[string]string)
				}
				if csvMetrics.Node[vec.Metric.Node] == nil {
					csvMetrics.Node[vec.Metric.Node] = make(map[string]string)
				}
				csvMetrics.Node[vec.Metric.Node][vec.Metric.Name] = value[1].(string)
			}
		}
	}
	return timeseries, alerts
}
