package openstorage

import (
	"errors"

	"github.com/libopenstorage/autopilot/api/autopilot/types"
	"github.com/libopenstorage/autopilot/telemetry"
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

	openstorage struct {
		types.OpenStorage
	}
)

// New returns a new prometheus instance
func New(prov types.Provider) (telemetry.Provider, error) {
	ops, ok := prov.(*types.OpenStorage)
	if !ok {
		return nil, errors.New("invalid provider type")
	}
	return &openstorage{
		OpenStorage: *ops,
	}, nil
}

// Query implements the telemetry.Provider.Query interface method
func (p *openstorage) Query(params telemetry.Params) ([]telemetry.Vector, error) {
	return nil, nil
}

// Parse implements the telemetry.Provider.Parse interface method
func (p *openstorage) Parse(data []byte) ([]telemetry.Vector, error) {
	return nil, nil
}

func init() {
	telemetry.Register(types.ProviderTypeOpenStorage, New)
}
