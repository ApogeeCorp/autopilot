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

// New returns a new prometheus instance
func New(params telemetry.Params) (telemetry.Provider, error) {
	return &prometheus{
		params: params,
		url:    params.String("url"),
	}, nil
}

// Query implements the telemetry.Provider.Query interface method
func (p *prometheus) Query(params telemetry.Params) ([]telemetry.Vector, error) {
	client := &http.Client{}

	base, err := url.Parse(p.url)
	if err != nil {
		return nil, err
	}

	base.Path = path.Join(base.Path, params.String("path", "/query"))

	req, _ := http.NewRequest("GET", base.String(), nil)

	q := req.URL.Query()

	if query, ok := params.IsSetV("query"); ok {
		q.Add("query", query.String())
	}

	if v, ok := params.IsSetV("start"); ok {
		start, err := time.Parse(time.RFC3339, v.String())
		if err != nil {
			return nil, err
		}
		q.Add("start", fmt.Sprint(start.UTC().Unix()))
	}

	if v, ok := params.IsSetV("end"); ok {
		end, err := time.Parse(time.RFC3339, v.String())
		if err != nil {
			return nil, err
		}
		q.Add("end", fmt.Sprint(end.UTC().Unix()))
	}

	if step, ok := params.IsSetV("step"); ok {
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

	return p.Parse(data)
}

// Parse implements the telemetry.Provider.Parse interface method
func (p *prometheus) Parse(data []byte) ([]telemetry.Vector, error) {
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

func (p *prometheus) Exec(policy *telemetry.StoragePolicy) (bool, error) {
	log.Infof("checking policy %q", policy.Name)
	return false, nil
}

func init() {
	telemetry.Register("prometheus", New)
}
