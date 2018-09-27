// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package Name is the internal implementation of the AutopilotAPI
package autopilot

import (
	"github.com/sirupsen/logrus"
)

type API struct {
  Log *logrus.Logger
}

func (a *API) Initialize() error {
	return nil
}

// InitializeContext is called after authorization and before the API method.
// This method is used to setup the context for the next call.
func (a *API) InitializeContext(principal provider.AuthToken, r *http.Request) (*Context, middleware.Responder) {
	return &Context{principal}, nil
}
