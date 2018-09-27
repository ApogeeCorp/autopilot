// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/rest/operations/operations"
)

// SourceList Returns an array of telemetry sources defined in the system
func (a *API) SourceList(ctx *Context, params operations.SourceListParams) middleware.Responder {
	return sparks.ErrNotImplemented("sourceList")
}
