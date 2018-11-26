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

// SampleDelete Delete a sample from the disk
func (a *API) SampleDelete(ctx *Context, params sample.SampleDeleteParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleDelete")
}

// SampleList Returns an array of samples
func (a *API) SampleList(ctx *Context, params sample.SampleListParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleList")
}
