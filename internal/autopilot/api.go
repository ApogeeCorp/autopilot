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
