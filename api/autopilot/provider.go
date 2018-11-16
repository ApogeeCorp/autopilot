// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/provider"
)

// ProviderCreate Create a new telemetry provider from the provided definition
func (a *API) ProviderCreate(ctx *Context, params provider.ProviderCreateParams) middleware.Responder {
	return sparks.ErrNotImplemented("providerCreate")
}

// ProviderDelete Returns the request collected object
func (a *API) ProviderDelete(ctx *Context, params provider.ProviderDeleteParams) middleware.Responder {
	return sparks.ErrNotImplemented("providerDelete")
}

// ProviderGet Returns the request collected object
func (a *API) ProviderGet(ctx *Context, params provider.ProviderGetParams) middleware.Responder {
	return sparks.ErrNotImplemented("providerGet")
}

// ProviderList Returns an array of telemetry providers defined in the system
func (a *API) ProviderList(ctx *Context, params provider.ProviderListParams) middleware.Responder {
	return sparks.ErrNotImplemented("providerList")
}

// ProviderUpdate Update the properties of the specified provider
func (a *API) ProviderUpdate(ctx *Context, params provider.ProviderUpdateParams) middleware.Responder {
	return sparks.ErrNotImplemented("providerUpdate")
}
