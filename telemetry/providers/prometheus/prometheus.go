package prometheus

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	// Cluster stats
	Cluster = "Cluster"
	// Node stats
	Node = "Node"
	// Disk stats
	Disk = "Disk"
	// Volume stats
	Volume = "Volume"
	// Pool stats
	Pool = "Pool"
	// Proc stats - on the PX processes
	Proc = "Proc"
	// Alerts metrics are a special case for outputs
	Alerts = "ALERTS"
	// Instance field in the CSV
	Instance = "instance"
	// Timestamp field in the CSV
	Timestamp = "timestamp"
)

type (
	// Metric is the metric that comes from Prometheus apiproxy
	Metric struct {
		Name          string  `json:"__name__"`
		Cluster       string  `json:"cluster"`
		Instance      string  `json:"instance"`
		Node          string  `json:"node_id"`
		NodeName      string  `json:"node"`
		Job           *string `json:"job,omitempty"`
		Volume        *string `json:"volumeid,omitempty"`
		VolumeName    *string `json:"volumename,omitempty"`
		VolumePVC     *string `json:"volumepvc,omitempty"`
		Disk          *string `json:"disk,omitempty"`
		Pool          *string `json:"pool,omitempty"`
		AlertName     *string `json:"alertname,omitempty"`
		AlertState    *string `json:"alertstate,omitempty"`
		AlertSeverity *string `json:"severity,omitempty"`
		AlertIssue    *string `json:"issue,omitempty"`
		Proc          *string `json:"proc,omitempty"`
	}

	// Result for a single metric and value
	Result struct {
		Metric Metric `json:"metric"`
		// Value is [timestamp, value] such as "value": [1231231233.232, "0"], its mixed type
		Value []interface{} `json:"value"`
	}

	// ClusterResults is the complete set of metrics for the cluster
	ClusterResults struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string   `json:"resultType"`
			Results    []Result `json:"result"`
		} `json:"data"`
	}

	// Structures for Underlying Representation for ML

	// CSVRow is the common characteristics of every metric - this defines a row in the CSV
	CSVRow struct {
		Timestamp uint32
		Cluster   string
		Instance  string
		Node      string
	}

	// CSVMetrics contain ALL the BaseAttributes for a row in the CSV Each map is keyed by the Name of the field (VolumeID, Disk, Pool, Proc)
	// This probably could contain a map of maps of maps, but its becoming unreadable with that level of nesting.
	CSVMetrics struct {
		// Cluster are the cluster based metrics keyed on cluster field
		Cluster map[string]map[string]string
		// Node are the node based metrics keyed on cluster field
		Node map[string]map[string]string
		// Volume are the volume based metrics keyed on Volume field
		Volume map[string]map[string]string
		// Disk are the disk based metrics keyed on Disk field
		Disk map[string]map[string]string
		// Pool are the pool based metrics keyed on Pool field
		Pool map[string]map[string]string
		// Proc are the proc based metrics keyed on Proc field
		Proc map[string]map[string]string
	}

	// AlertRow is a special type of metric that we use for the Output of the ML.  This is one of the predictors
	AlertRow struct {
		csvRow        CSVRow
		AlertName     string
		AlertState    string
		AlertSeverity string
		AlertIssue    string
		AlertValue    string
	}
)

// NewCSVMetrics constructor to initialize internal maps
func NewCSVMetrics() *CSVMetrics {
	return &CSVMetrics{
		Cluster: make(map[string]map[string]string),
		Node:    make(map[string]map[string]string),
		Volume:  make(map[string]map[string]string),
		Disk:    make(map[string]map[string]string),
		Pool:    make(map[string]map[string]string),
		Proc:    make(map[string]map[string]string),
	}
}

// NewAlertRow Create a new AlertRow
func NewAlertRow(csvRow CSVRow, result *Result) *AlertRow {
	return &AlertRow{
		csvRow:        csvRow,
		AlertName:     *result.Metric.AlertName,
		AlertState:    *result.Metric.AlertState,
		AlertSeverity: *result.Metric.AlertSeverity,
		AlertIssue:    *result.Metric.AlertIssue,
		AlertValue:    result.Value[1].(string),
	}
}

// Keys Return keys of the given map - used to pull the values from a particular VolumeID, Disk, Pool, Proc of a BaseMetric
func Keys(m map[string]string) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// CreateCSVRows creates the rows for a particular MAP in the CSVMetrics
func CreateCSVRows(br CSVRow, m map[string]map[string]string, label string) (rows [][]string, rowHeaders []string) {
	for rowKey, rowValue := range m {
		var row = []string{strconv.FormatUint(uint64(br.Timestamp), 10), br.Cluster, br.Instance, br.Node}
		var rowHeader = []string{Timestamp, Cluster, Instance, Node}
		if label != "Cluster" && label != "Node" {
			row = append(row, rowKey)
			rowHeader = append(rowHeader, label)
		}
		for _, k := range Keys(rowValue) {
			v := rowValue[k]
			row = append(row, v)
			rowHeader = append(rowHeader, k)
		}
		rows = append(rows, row)
		rowHeaders = rowHeader
	}
	return rows, rowHeaders
}

// WriteCSV for a particular metric category, create the CSV file
func WriteCSV(timeSeries map[CSVRow]*CSVMetrics, name string) {
	f, err := os.Create("./" + name + ".csv")
	if err != nil {
		fmt.Println("Error Creating File", name, err)
	}

	defer f.Close()
	w := csv.NewWriter(f)
	var wroteHeader = false

	for br, bm := range timeSeries {
		// This is a bit of a hack to avoid reflection. If the CSVMetrics had a map instead of fields, we could avoid this.
		var m map[string]map[string]string
		if name == Volume {
			m = bm.Volume
		} else if name == Node {
			m = bm.Node
		} else if name == Disk {
			m = bm.Disk
		} else if name == Pool {
			m = bm.Pool
		} else if name == Proc {
			m = bm.Proc
		} else if name == Cluster {
			m = bm.Cluster
		}

		csvRows, csvHeader := CreateCSVRows(br, m, name)
		if wroteHeader == false {
			w.Write(csvHeader)
			wroteHeader = true
		}
		for _, csvRow := range csvRows {
			w.Write(csvRow)
		}
	}
	w.Flush()
}

// WriteAlertCSV for the alerts csv
func WriteAlertCSV(alerts []*AlertRow) {
	f, err := os.Create("./Alerts.csv")
	if err != nil {
		fmt.Println("Error Creating Alerts File", err)
	}

	defer f.Close()
	w := csv.NewWriter(f)
	w.Write([]string{Timestamp, Cluster, Instance, Node, "alert_name", "alert_state", "alert_severity", "alert_value"})
	for _, alert := range alerts {
		var row []string
		row = append(row, strconv.FormatUint(uint64(alert.csvRow.Timestamp), 10), alert.csvRow.Cluster, alert.csvRow.Instance, alert.csvRow.Node, alert.AlertName, alert.AlertState, alert.AlertSeverity, alert.AlertValue)
		w.Write(row)
	}
	w.Flush()
}

// TransformToRows takes the Prometheus API Calls ClusterResults and flattens it to the structure that can exported as CSV
func TransformToRows(results *ClusterResults) (timeseries map[CSVRow]*CSVMetrics, alerts []*AlertRow) {
	// Go through the results for the cluster and generate the appropriate CSV files
	timeseries = make(map[CSVRow]*CSVMetrics)

	for _, result := range results.Data.Results {
		csvRow := CSVRow{
			Timestamp: uint32(result.Value[0].(float64)),
			Cluster:   result.Metric.Cluster,
			Instance:  result.Metric.Instance,
			Node:      result.Metric.Node,
		}
		csvMetrics, ok := timeseries[csvRow]
		if !ok {
			csvMetrics = new(CSVMetrics)
			timeseries[csvRow] = csvMetrics
		}
		if result.Metric.Name == Alerts {
			alert := NewAlertRow(csvRow, &result)
			alerts = append(alerts, alert)
		} else if result.Metric.Volume != nil {
			if csvMetrics.Volume == nil {
				csvMetrics.Volume = make(map[string]map[string]string)
			}
			if csvMetrics.Volume[*result.Metric.Volume] == nil {
				csvMetrics.Volume[*result.Metric.Volume] = make(map[string]string)
			}
			csvMetrics.Volume[*result.Metric.Volume][result.Metric.Name] = result.Value[1].(string)
		} else if result.Metric.Disk != nil {
			if csvMetrics.Disk == nil {
				csvMetrics.Disk = make(map[string]map[string]string)
			}
			if csvMetrics.Disk[*result.Metric.Disk] == nil {
				csvMetrics.Disk[*result.Metric.Disk] = make(map[string]string)
			}
			csvMetrics.Disk[*result.Metric.Disk][result.Metric.Name] = result.Value[1].(string)
		} else if result.Metric.Pool != nil {
			if csvMetrics.Pool == nil {
				csvMetrics.Pool = make(map[string]map[string]string)
			}
			if csvMetrics.Pool[*result.Metric.Pool] == nil {
				csvMetrics.Pool[*result.Metric.Pool] = make(map[string]string)
			}
			csvMetrics.Pool[*result.Metric.Pool][result.Metric.Name] = result.Value[1].(string)
		} else if result.Metric.Proc != nil {
			if csvMetrics.Proc == nil {
				csvMetrics.Proc = make(map[string]map[string]string)
			}
			if csvMetrics.Proc[*result.Metric.Proc] == nil {
				csvMetrics.Proc[*result.Metric.Proc] = make(map[string]string)
			}
			csvMetrics.Proc[*result.Metric.Proc][result.Metric.Name] = result.Value[1].(string)
		} else if strings.HasPrefix(result.Metric.Name, "px_cluster_") == true {
			if csvMetrics.Cluster == nil {
				csvMetrics.Cluster = make(map[string]map[string]string)
			}
			if csvMetrics.Cluster[result.Metric.Cluster] == nil {
				csvMetrics.Cluster[result.Metric.Cluster] = make(map[string]string)
			}
			csvMetrics.Cluster[result.Metric.Cluster][result.Metric.Name] = result.Value[1].(string)
			// the px_node_status_<node_id>_status is ambigious, this will be fixed later so lets just ignore it
			//		} else if strings.HasPrefix(result.Metric.Name, "px_node_") == true || strings.HasPrefix(result.Metric.Name, "px_network_") == true {
		} else if strings.HasPrefix(result.Metric.Name, "px_node_stats") == true || strings.HasPrefix(result.Metric.Name, "px_network_") == true {
			if csvMetrics.Node == nil {
				csvMetrics.Node = make(map[string]map[string]string)
			}
			if csvMetrics.Node[result.Metric.Node] == nil {
				csvMetrics.Node[result.Metric.Node] = make(map[string]string)
			}
			csvMetrics.Node[result.Metric.Node][result.Metric.Name] = result.Value[1].(string)
		}
	}
	return timeseries, alerts
}
