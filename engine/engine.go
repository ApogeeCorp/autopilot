// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/template"
	"github.com/go-openapi/strfmt"
	"github.com/libopenstorage/autopilot/api/autopilot/types"
	"github.com/libopenstorage/autopilot/telemetry/providers/prometheus"
	"github.com/mithrandie/csvq/lib/action"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/query"
	"github.com/sirupsen/logrus"
)

// Engine is the autopilot recommendation engine
type Engine struct {
	Log        *logrus.Logger
	Prometheus *prometheus.Prometheus
}

func getFields(rule *types.Rule, samplePath string) (fieldStr string, fileStr string) {
	var files []string
	fields := regexp.MustCompile(`\bpx_(\w)*`).FindAllString(rule.Expr+" "+rule.Proposal+" "+rule.Issue, -1)
	fieldStr = strings.Join(fields, ",")
	if strings.Contains(fieldStr, "px_volume") {
		files = append(files, "`"+samplePath+"/volume.csv`")
		fields = append(fields, "volume")
	}
	if strings.Contains(fieldStr, "px_node") || strings.Contains(fieldStr, "px_network") || strings.Contains(fieldStr, "px_cluster") || strings.Contains(fieldStr, "px_proc") {
		files = append(files, "`"+samplePath+"/node.csv`")
	}
	if strings.Contains(fieldStr, "px_pool") {
		files = append(files, "`"+samplePath+"/pool.csv`")
		fields = append(fields, "pool")
	}
	if strings.Contains(fieldStr, "px_disk") {
		files = append(files, "`"+samplePath+"/disk.csv`")
		fields = append(fields, "disk")
	}
	fileStr = strings.Join(files, ",")
	fieldStr = strings.Join(fields, ",")
	return fieldStr, fileStr
}

func formatAsDate(timestamp string) string {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	unixTimeUTC := time.Unix(i, 0)          //gives unix time stamp in utc
	return unixTimeUTC.Format(time.RFC3339) // converts utc time to RFC3339 format
}

func formatFloat(value string) string {
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%.1f", floatVal)
}

func (e *Engine) executePromQLRule(rule *types.Rule) (*types.Recommendation, error) {
	e.Log.Debugf("Executing PromQL Rule %v", rule)
	return e.Prometheus.Query("http://70.0.69.141:9090/api/v1/query", rule)
}

// Recommend returns a recommendation from the engine based on the rules and sample
func (e *Engine) Recommend(rules []*types.Rule, samplePath string) (recommendations []*types.Recommendation, error error) {
	flags := cmd.GetFlags()
	flags.SetFormat("", "out.json")
	fmap := template.FuncMap{
		"formatAsDate": formatAsDate,
		"formatFloat":  formatFloat,
	}

	for _, rule := range rules {

		if rule.Type == "promql" {
			recommendation, err := e.executePromQLRule(rule)
			if err == nil {
				recommendations = append(recommendations, recommendation)
			} else {
			}
		} else {
			e.Log.Debugf("Processing Rule %s, %v", samplePath, rule.Expr)
			proc := query.NewProcedure()
			fieldStr, fileStr := getFields(rule, samplePath)
			outFileStr := samplePath + "/" + rule.Name + ".json"
			queryStr := "select timestamp, cluster, instance, node, (" + rule.Expr[0:len(rule.Expr)-4] + ") as result, " + fieldStr + " FROM " + fileStr + " WHERE " + rule.Expr
			//e.Log.Debugf("Query String %s", queryStr)

			err := action.Run(proc, queryStr, samplePath, outFileStr)
			if err != nil {
				e.Log.Debugf("Could not run actions %v", err)
			}
			b, err := ioutil.ReadFile(outFileStr)
			if err != nil {
				e.Log.Debugf("Could not read Recommendations %v", err)
			}
			if len(b) > 5 {
				var recommendation types.Recommendation
				recommendation.Timestamp = strfmt.DateTime(time.Now())
				var results []map[string]interface{}
				json.Unmarshal([]byte(b), &results)
				//e.Log.Debugf("Executed Procedure %v", results)
				for _, result := range results {
					var proposal types.Proposal
					proposal.Rule = rule.Name
					proposal.ClusterID = result["cluster"].(string)
					proposal.NodeID = result["node"].(string)
					if result["volume"] != nil {
						proposal.VolumeID = result["volume"].(string)
					}
					t := template.Must(template.New("Issue").Funcs(fmap).Parse("{{.timestamp | formatAsDate}}: " + rule.Issue + ` ` + rule.Proposal))
					//prop := template.Must(template.New("Proposal").Parse(rule.Proposal)
					var proposalValue bytes.Buffer
					err := t.Execute(&proposalValue, result)
					if err != nil {
						e.Log.Debugf("Could not parse issue with result %v", result)
					}
					proposal.Action = proposalValue.String()
					recommendation.Proposals = append(recommendation.Proposals, &proposal)
				}
				recommendations = append(recommendations, &recommendation)
			}
		}

	}
	e.Log.Debugf("Finished With Recommendation \n\n")

	return recommendations, nil
}
