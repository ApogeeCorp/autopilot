// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	"os"
	"path"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gosimple/slug"
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/collector"
	"github.com/libopenstorage/autopilot/api/autopilot/types"
	"github.com/libopenstorage/autopilot/telemetry/providers/prometheus"
)

// CollectorList Returns an array of telemetry collectors defined in the system
func (a *API) CollectorList(ctx *Context, params collector.CollectorListParams) middleware.Responder {
	return collector.NewCollectorListOK().WithPayload(a.Config.Collectors)
}

// CollectorPoll Polls a collector for the current data period
func (a *API) CollectorPoll(ctx *Context, params collector.CollectorPollParams) middleware.Responder {
	var col *types.Collector

	for _, c := range a.Config.Collectors {
		if c.Name == params.Collector {
			col = c
			break
		}
	}

	if col == nil {
		return sparks.ErrNotFound("collector")
	}

	dataPath := path.Join(a.DataDir, "collectors", slug.Make(col.Name))
	if err := os.MkdirAll(dataPath, 0770); err != nil {
		return sparks.NewError(err)
	}

	switch col.Type {
	case types.CollectorTypePrometheus:
		prov := &prometheus.Prometheus{Log: a.Log}
		if err := prov.Collect(col.URL, col.Params, dataPath); err != nil {
			return sparks.NewError(err)
		}

	default:
		return sparks.NewError().Format("Invalid collector type")
	}

	return nil
}
