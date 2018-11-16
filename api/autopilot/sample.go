// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/sample"
)

// RecommendationsGet Returns the recommendations for a particular sample
func (a *API) RecommendationsGet(ctx *Context, params sample.RecommendationsGetParams) middleware.Responder {
	return sparks.ErrNotImplemented("recommendationsGet")
}

// SampleCreate Create a new telemetry sample from the provided definition
func (a *API) SampleCreate(ctx *Context, params sample.SampleCreateParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleCreate")
}

// SampleDelete Returns the request collected object
func (a *API) SampleDelete(ctx *Context, params sample.SampleDeleteParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleDelete")
}

// SampleGet Returns the request collected object
func (a *API) SampleGet(ctx *Context, params sample.SampleGetParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleGet")
}

// SampleList Returns an array of telemetry samples defined in the system
func (a *API) SampleList(ctx *Context, params sample.SampleListParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleList")
}

// SampleUpdate Update the properties of the specified sample
func (a *API) SampleUpdate(ctx *Context, params sample.SampleUpdateParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleUpdate")
}
