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

	"github.com/libopenstorage/autopilot/api/autopilot/types"
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
		types.Prometheus
	}
)

// New returns a new prometheus instance
func New(prov types.Provider) (telemetry.Provider, error) {
	prom, ok := prov.(*types.Prometheus)
	if !ok {
		return nil, errors.New("invalid provider type")
	}
	return &prometheus{
		Prometheus: *prom,
	}, nil
}

// Query implements the telemetry.Provider.Query interface method
func (p *prometheus) Query(params telemetry.Params) ([]telemetry.Vector, error) {
	client := &http.Client{}

	base, err := url.Parse(p.URL)
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

	log.Debugf("provider %s: executing query %s", types.ProviderTypePrometheus, req.URL.String())

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

func init() {
	telemetry.Register(types.ProviderTypePrometheus, New)
}
