// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/collector"
)

// CollectorCreate Create a new telemetry collector from the provided definition
func (a *API) CollectorCreate(ctx *Context, params collector.CollectorCreateParams) middleware.Responder {
	return sparks.ErrNotImplemented("collectorCreate")
}

// CollectorDelete Returns the request collected object
func (a *API) CollectorDelete(ctx *Context, params collector.CollectorDeleteParams) middleware.Responder {
	return sparks.ErrNotImplemented("collectorDelete")
}

// CollectorGet Returns the request collected object
func (a *API) CollectorGet(ctx *Context, params collector.CollectorGetParams) middleware.Responder {
	return sparks.ErrNotImplemented("collectorGet")
}

// CollectorList Returns an array of telemetry collectors defined in the system
func (a *API) CollectorList(ctx *Context, params collector.CollectorListParams) middleware.Responder {
	return sparks.ErrNotImplemented("collectorList")
}

// CollectorUpdate Update the properties of the specified collector
func (a *API) CollectorUpdate(ctx *Context, params collector.CollectorUpdateParams) middleware.Responder {
	return sparks.ErrNotImplemented("collectorUpdate")
}
