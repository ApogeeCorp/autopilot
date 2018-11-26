// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	"github.com/go-openapi/runtime/middleware"
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/collector"
	"github.com/libopenstorage/autopilot/api/autopilot/types"
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

	switch col.Type {
	case types.CollectorTypePrometheus:
	default:
		return sparks.ErrInvalidParameter
	}

	return nil
}
