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

	p := &prometheus.Prometheus{}
	base := "Prometheus" //fmt.Sprintf("%s/prometheus", stagingPath)

	timeseries, alerts := p.TransformToRows(&results)
	//fmt.Printf("TimeSeries Length %# v", pretty.Formatter(timeseries))

	p.WriteCSV(timeseries, base, prometheus.Volume)
	p.WriteCSV(timeseries, base, prometheus.Disk)
	p.WriteCSV(timeseries, base, prometheus.Pool)
	p.WriteCSV(timeseries, base, prometheus.Node)

	p.WriteAlertCSV(base, alerts)

}
