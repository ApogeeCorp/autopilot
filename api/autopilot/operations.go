// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/operations"
)

// CollectorCreate Create a new telemetry collector from the provided definition
func (a *API) CollectorCreate(ctx *Context, params operations.CollectorCreateParams) middleware.Responder {
	return sparks.ErrNotImplemented("collectorCreate")
}

// CollectorList Returns an array of telemetry collectors defined in the system
func (a *API) CollectorList(ctx *Context, params operations.CollectorListParams) middleware.Responder {
	return sparks.ErrNotImplemented("collectorList")
}
