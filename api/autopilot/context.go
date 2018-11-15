// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the 
// root of this project.

package autopilot

import (
  "gitlab.com/ModelRocket/sparks/cloud/provider"
)

// Context is the API method context it is passed to the API calls with their
// parameters, the default struct has an AuthToken, but this is not required.
type Context struct {
  AuthToken provider.AuthToken
}
