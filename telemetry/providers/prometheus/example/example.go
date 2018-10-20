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

	//p := &prometheus.Prometheus{}

	//timeseries, alerts := p.TransformToRows(&results)
	//fmt.Printf("TimeSeries Length %# v", pretty.Formatter(timeseries))
	/*
		p.WriteCSV(timeseries, prometheus.Volume)
		p.WriteCSV(timeseries, prometheus.Disk)
		p.WriteCSV(timeseries, prometheus.Pool)
		p.WriteCSV(timeseries, prometheus.Proc)
		p.WriteCSV(timeseries, prometheus.Cluster)
		p.WriteCSV(timeseries, prometheus.Node)

		p.WriteAlertCSV(alerts)
	*/
}
