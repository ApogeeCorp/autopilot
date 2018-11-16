// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/rule"
)

// RuleCreate Create a new telemetry rule from the provided definition
func (a *API) RuleCreate(ctx *Context, params rule.RuleCreateParams) middleware.Responder {
	return sparks.ErrNotImplemented("ruleCreate")
}

// RuleDelete Returns the request collected object
func (a *API) RuleDelete(ctx *Context, params rule.RuleDeleteParams) middleware.Responder {
	return sparks.ErrNotImplemented("ruleDelete")
}

// RuleGet Returns the request collected object
func (a *API) RuleGet(ctx *Context, params rule.RuleGetParams) middleware.Responder {
	return sparks.ErrNotImplemented("ruleGet")
}

// RuleList Returns an array of telemetry rules defined in the system
func (a *API) RuleList(ctx *Context, params rule.RuleListParams) middleware.Responder {
	return sparks.ErrNotImplemented("ruleList")
}

// RuleUpdate Update the properties of the specified rule
func (a *API) RuleUpdate(ctx *Context, params rule.RuleUpdateParams) middleware.Responder {
	return sparks.ErrNotImplemented("ruleUpdate")
}
