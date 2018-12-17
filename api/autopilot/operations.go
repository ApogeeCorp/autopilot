// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations"
	"github.com/libopenstorage/autopilot/api/autopilot/types"
)

// CollectorList Returns an array of telemetry collectors defined in the system
func (a *API) CollectorList(ctx *Context, params operations.CollectorListParams) middleware.Responder {
	return operations.NewCollectorListOK().WithPayload(a.Config.Collectors)
}

// CollectorPoll Poll a collector for the given period directly
func (a *API) CollectorPoll(ctx *Context, params operations.CollectorPollParams) middleware.Responder {
	return sparks.ErrNotImplemented("collectorPoll")
}

// EmitterList Returns an array of telemetry emitters defined in the system
func (a *API) EmitterList(ctx *Context, params operations.EmitterListParams) middleware.Responder {
	return operations.NewEmitterListOK().WithPayload(a.Config.Emitters)
}

// ProviderList Returns an array of telemetry providers defined in the system
func (a *API) ProviderList(ctx *Context, params operations.ProviderListParams) middleware.Responder {
	return operations.NewProviderListOK().WithPayload(a.Config.Providers)
}

// ProviderQuery Query a provider directly
func (a *API) ProviderQuery(ctx *Context, params operations.ProviderQueryParams) middleware.Responder {
	return sparks.ErrNotImplemented("providerQuery")
}

// RecommendationsGet Create a new telemetry sample from the provided definition and get recommendations
func (a *API) RecommendationsGet(ctx *Context, params operations.RecommendationsGetParams) middleware.Responder {
	rules := make([]*types.Rule, 0)

	for _, rule := range params.Rules {
		rule, ok := a.Config.GetRule(rule)
		if !ok {
			return sparks.ErrInvalidParameter.Reason("rule")
		}
		rules = append(rules, rule)
	}

	recs, err := a.Engine.Recommend(params.Provider, rules)
	if err != nil {
		return sparks.NewError(err)
	}

	return operations.NewRecommendationsGetOK().WithPayload(recs)
}

// RuleList Returns an array of telemetry rules defined in the system
func (a *API) RuleList(ctx *Context, params operations.RuleListParams) middleware.Responder {
	return operations.NewRuleListOK().WithPayload(a.Config.Rules)
}
