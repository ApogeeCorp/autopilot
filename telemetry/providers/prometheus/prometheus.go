/*
Copyright 2019 Openstorage.org

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package prometheus

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	meta "github.com/libopenstorage/autopilot/pkg/apis/autopilot/v1alpha1"
	"github.com/libopenstorage/autopilot/telemetry"
	log "github.com/sirupsen/logrus"
)

type (
	// results is the complete set of metrics
	results struct {
		// Status is the results status
		Status string `json:"status"`

		// Data is the data for the results
		Data struct {
			ResultType string             `json:"resultType"`
			Results    []telemetry.Vector `json:"result"`
		} `json:"data"`

		// ErrorType is the prometheus error type
		ErrorType string `json:"errorType"`

		// Error is the error message
		Error string `json:"error"`
	}

	prometheus struct {
		params telemetry.Params
		url    string
	}
)

var promQLMetricLookup = map[string]string{
	"openstorage.io/condition.volume.usage_percentage": "100 * (px_volume_usage_bytes / px_volume_capacity_bytes)",
	"openstorage.io/condition.volume.capacity_gb":      "ps_volume_fs_capacity_bytes / 1000000000",
}

var promQLOperatorLookup = map[meta.LabelSelectorOperator]string{
	"gt": ">",
	"lt": "<",
	"eq": "=",
}

// New returns a new prometheus instance
func New(params telemetry.Params) (telemetry.Provider, error) {
	return &prometheus{
		params: params,
		url:    params.String("url"),
	}, nil
}

// query implements the telemetry.Provider.Query interface method
func (p *prometheus) query(params telemetry.Params) ([]telemetry.Vector, error) {
	client := &http.Client{}

	base, err := url.Parse(p.url)
	if err != nil {
		return nil, err
	}

	base.Path = path.Join(base.Path, params.String("path", "/query"))

	req, _ := http.NewRequest("GET", base.String(), nil)

	q := req.URL.Query()

	if query, ok := params.GetValue("query"); ok {
		q.Add("query", query.String())
	}

	if v, ok := params.GetValue("start"); ok {
		start, err := time.Parse(time.RFC3339, v.String())
		if err != nil {
			return nil, err
		}
		q.Add("start", fmt.Sprint(start.UTC().Unix()))
	}

	if v, ok := params.GetValue("end"); ok {
		end, err := time.Parse(time.RFC3339, v.String())
		if err != nil {
			return nil, err
		}
		q.Add("end", fmt.Sprint(end.UTC().Unix()))
	}

	if step, ok := params.GetValue("step"); ok {
		q.Add("step", step.String())
	}

	req.URL.RawQuery = q.Encode()

	log.Debugf("prometheus: executing query %s", req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get data: %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return p.parse(data)
}

// parse implements the telemetry.Provider.Parse interface method
func (p *prometheus) parse(data []byte) ([]telemetry.Vector, error) {
	results := &results{}
	err := json.Unmarshal(data, results)
	if err != nil {
		return nil, err
	}

	switch results.Status {
	case "success":
		return results.Data.Results, nil
	case "error":
		return nil, fmt.Errorf("%s: %s", results.ErrorType, results.Error)
	default:
		return nil, errors.New("invalid return data")
	}
}

func (p *prometheus) LookupMetric(metric string) string {
	return promQLMetricLookup[metric]
}

func (p *prometheus) LookupOperator(operator meta.LabelSelectorOperator) string {
	return promQLOperatorLookup[operator]
}

func (p *prometheus) ConditionToQuery(condition *meta.LabelSelectorRequirement) string {
	return p.LookupMetric(condition.Key) + " " + p.LookupOperator(condition.Operator) + " " + condition.Values[0]
}

func (p *prometheus) Query(policy *telemetry.StoragePolicy) ([]telemetry.Vector, error) {
	rval := make([]telemetry.Vector, 0)

	for _, c := range policy.Spec.Conditions {
		log.Infof("Condition % #v", c)
		m := make(telemetry.Params)
		m["query"] = p.ConditionToQuery(c)
		log.Infof("   Prometheus Query %#v", m["query"])
		vectors, err := p.query(m)
		if err != nil {
			log.Infof("Error Executing Policy %q, %s, % #v", policy.Name, c.Key, err)
			return nil, err
		}
		if len(vectors) > 0 {
			rval = append(rval, vectors...)
		}
	}

	return rval, nil
}

func init() {
	telemetry.Register("prometheus", New)
}
