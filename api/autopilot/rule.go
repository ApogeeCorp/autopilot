// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/rule"
	"github.com/libopenstorage/autopilot/api/autopilot/types"
)

// RuleCreate Create a new telemetry rule from the provided definition
func (a *API) RuleCreate(ctx *Context, params rule.RuleCreateParams) middleware.Responder {
	if err := a.DB.Create(params.Rule).Error; err != nil {
		return sparks.NewError(err)
	}
	return rule.NewRuleCreateCreated().WithPayload(params.Rule)
}

// RuleDelete Returns the request collected object
func (a *API) RuleDelete(ctx *Context, params rule.RuleDeleteParams) middleware.Responder {
	if err := a.DB.Delete(&types.RuleSet{}, "id=?", params.RuleID).Error; err != nil {
		return sparks.NewError(err)
	}
	return rule.NewRuleDeleteNoContent()
}

// RuleGet Returns the request collected object
func (a *API) RuleGet(ctx *Context, params rule.RuleGetParams) middleware.Responder {
	r := &types.RuleSet{}
	if err := a.DB.First(r, "id=?", params.RuleID).Error; err != nil {
		return sparks.NewError(err)
	}
	return rule.NewRuleGetOK().WithPayload(r)
}

// RuleList Returns an array of telemetry rules defined in the system
func (a *API) RuleList(ctx *Context, params rule.RuleListParams) middleware.Responder {
	rules := make([]*types.RuleSet, 0)

	if err := a.DB.Find(rules).Error; err != nil {
		return sparks.NewError(err)
	}
	return rule.NewRuleListOK().WithPayload(rules)
}

// RuleUpdate Update the properties of the specified rule
func (a *API) RuleUpdate(ctx *Context, params rule.RuleUpdateParams) middleware.Responder {
	if err := a.DB.Model(params.Rule).Where("id=?", params.RuleID).Updates(params.Rule).Error; err != nil {
		return sparks.NewError(err)
	}
	return rule.NewRuleUpdateNoContent()
}
