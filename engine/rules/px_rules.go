package rules

import (
	"fmt"

	"github.com/spf13/viper"
)

/*
rules:
- alert: PortworxVolumeUsageCritical
  expr: 100 * (px_volume_usage_bytes / px_volume_capacity_bytes) > 80
  for: 5m
  labels:
	issue: Portworx volume {{$labels.volumeid}} usage on {{$labels.host}} is high.
	severity: critical
  annotations:
	description: Portworx volume {{$labels.volumeid}} on {{$labels.host}} is over
	  80% used for more than 10 minutes.
	summary: Portworx volume capacity is at {{$value}}% used.
type (
	Rules struct {
		Rules []Rule `yaml:"rules"`
	}
	type Rule struct {
		Record      string            `yaml:"record,omitempty"`
		Alert       string            `yaml:"alert,omitempty"`
		Expr        string            `yaml:"expr"`
		// For         time.Duration    `yaml:"for,omitempty"`
		Labels      map[string]string `yaml:"labels,omitempty"`
		Annotations map[string]string `yaml:"annotations,omitempty"`
	}

	Rule struct {
		Alert string `yaml:"alert"`
		Expr string `yaml:"expr"`
		For string `yaml:"for"`
		Expr string `yaml:"alert"`
		Expr string `yaml:"alert"`
		Expr string `yaml:"alert"`
		For string
	}
)
*/

type (
	// Rules to execute against CSV Files Parsed
	Rules interface{}
)

// ReadRules from the config file
func ReadRules() {
	viper.SetConfigName("autopilot-default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	fmt.Println("Got Rules")
	fmt.Println(viper.Get("rules"))
}
