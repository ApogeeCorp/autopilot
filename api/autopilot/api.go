// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the 
// root of this project.

// Package autopilot is the internal implementation of the AutopilotAPI
package autopilot

import (
	"github.com/sirupsen/logrus"
)

// API is the acme API interface implementation
type API struct {
  Log *logrus.Logger
}

// Initialize initializes the API before the server starts handling request
func (a *API) Initialize() error {
	return nil
}

// InitializeContext is called after authorization and before the API method.
// This method is used to setup the context for the next call.
func (a *API) InitializeContext(principal provider.AuthToken, r *http.Request) (*Context, error) {
	return &Context{principal}, nil
}
