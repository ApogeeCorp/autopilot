// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/task"
)

// TaskList Returns an array of tasks
func (a *API) TaskList(ctx *Context, params task.TaskListParams) middleware.Responder {
	return sparks.ErrNotImplemented("taskList")
}
