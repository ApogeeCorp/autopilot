// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/rule"
)

// RuleList Returns an array of telemetry rules defined in the system
func (a *API) RuleList(ctx *Context, params rule.RuleListParams) middleware.Responder {
	return rule.NewRuleListOK().WithPayload(a.Config.Rules)
}
