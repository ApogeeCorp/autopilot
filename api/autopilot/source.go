// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/source"
)

// SourceCreate Create a new telemetry source from the provided definition
func (a *API) SourceCreate(ctx *Context, params source.SourceCreateParams) middleware.Responder {
	return sparks.ErrNotImplemented("sourceCreate")
}

// SourceDelete Returns the request collected object
func (a *API) SourceDelete(ctx *Context, params source.SourceDeleteParams) middleware.Responder {
	return sparks.ErrNotImplemented("sourceDelete")
}

// SourceGet Returns the request collected object
func (a *API) SourceGet(ctx *Context, params source.SourceGetParams) middleware.Responder {
	return sparks.ErrNotImplemented("sourceGet")
}

// SourceList Returns an array of telemetry sources defined in the system
func (a *API) SourceList(ctx *Context, params source.SourceListParams) middleware.Responder {
	return sparks.ErrNotImplemented("sourceList")
}

// SourceUpdate Update the properties of the specified source
func (a *API) SourceUpdate(ctx *Context, params source.SourceUpdateParams) middleware.Responder {
	return sparks.ErrNotImplemented("sourceUpdate")
}

// SourcePoll Poll a source and collect a sample manually
func (a *API) SourcePoll(ctx *Context, params source.SourcePollParams) middleware.Responder {
	return sparks.ErrNotImplemented("sourcePoll")
}
