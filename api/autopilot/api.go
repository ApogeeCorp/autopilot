// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

// Package autopilot is the internal implementation of the AutopilotAPI
package autopilot

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/libopenstorage/autopilot/engine"
	"github.com/sirupsen/logrus"
)

// API is the acme API interface implementation
type API struct {
	Log    *logrus.Logger
	DB     *gorm.DB
	engine *engine.Engine
}

// Initialize initializes the API before the server starts handling request
func (a *API) Initialize() error {
	a.engine = &engine.Engine{
		Log: a.Log,
	}
	return nil
}

// InitializeContext is called after authorization and before the API method.
// This method is used to setup the context for the next call.
func (a *API) InitializeContext(principal interface{}, r *http.Request) (*Context, error) {
	return &Context{Username: principal.(string)}, nil
}
