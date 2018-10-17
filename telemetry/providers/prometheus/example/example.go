package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	prometheus "gitlab.com/ModelRocket/portworx/autopilot/telemetry/providers/prometheus"
)

func main() {
	// reading data from JSON File
	data, err := ioutil.ReadFile("./prometheus.json")
	if err != nil {
		fmt.Println(err)
	}
	// Unmarshal JSON data
	var results prometheus.ClusterResults
	err = json.Unmarshal([]byte(data), &results)
	if err != nil {
		fmt.Println(err)
	}

	timeseries, alerts := prometheus.TransformToRows(&results)

	prometheus.WriteCSV(timeseries, prometheus.Volume)
	prometheus.WriteCSV(timeseries, prometheus.Disk)
	prometheus.WriteCSV(timeseries, prometheus.Pool)
	prometheus.WriteCSV(timeseries, prometheus.Proc)
	prometheus.WriteCSV(timeseries, prometheus.Cluster)
	prometheus.WriteCSV(timeseries, prometheus.Node)

	prometheus.WriteAlertCSV(alerts)
}
