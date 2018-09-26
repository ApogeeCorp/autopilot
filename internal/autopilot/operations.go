//
//  MODEL ROCKET LLC CONFIDENTIAL
//  _________________
//   Copyright (c) 2018 - 2019 MODEL ROCKET LLC
//   All Rights Reserved.
//
//   NOTICE:  All information contained herein is, and remains
//   the property of MODEL ROCKET LLC and its suppliers,
//   if any.  The intellectual and technical concepts contained
//   herein are proprietary to MODEL ROCKET LLC
//   and its suppliers and may be covered by U.S. and Foreign Patents,
//   patents in process, and are protected by trade secret or copyright law.
//   Dissemination of this information or reproduction of this material
//   is strictly forbidden unless prior written permission is obtained
//   from MODEL ROCKET LLC.
//

package autopilot

import (
	sparks "github.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/rest/operations/operations"
)

// TelemetrySourceList Returns an array of telemetry sources defined in the system
func (a *API) TelemetrySourceList(ctx *Context, params operations.TelemetrySourceListParams) middleware.Responder {
	return sparks.ErrNotImplemented("telemetrySourceList")
}
