package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/libopenstorage/autopilot/telemetry/providers/prometheus"
)

func main() {
	// reading data from JSON File
	data, err := ioutil.ReadFile("./prometheus_day_6h.json")
	if err != nil {
		fmt.Println(err)
	}
	// Unmarshal JSON data
	var results prometheus.ClusterResults
	err = json.Unmarshal([]byte(data), &results)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf(" RESULTS %# v\n", pretty.Formatter(results))

	timeseries, alerts := prometheus.TransformToRows(&results)
	//fmt.Printf("TimeSeries Length %# v", pretty.Formatter(timeseries))

	prometheus.WriteCSV(timeseries, prometheus.Volume)
	prometheus.WriteCSV(timeseries, prometheus.Disk)
	prometheus.WriteCSV(timeseries, prometheus.Pool)
	prometheus.WriteCSV(timeseries, prometheus.Proc)
	prometheus.WriteCSV(timeseries, prometheus.Cluster)
	prometheus.WriteCSV(timeseries, prometheus.Node)

	prometheus.WriteAlertCSV(alerts)
}
