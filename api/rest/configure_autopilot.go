// Code generated by hiro; DO NOT EDIT.

// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dre1080/recover"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"
	"gitlab.com/ModelRocket/sparks/cloud/provider"

	interpose "github.com/carbocation/interpose/middleware"
	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/libopenstorage/autopilot/api/rest/operations"
)

type contextKey string

const AuthKey contextKey = "Auth"

// OperationsAPI
type OperationsAPI interface {
	// SourceList is Returns an array of telemetry sources defined in the system
	SourceList(ctx *autopilot.Context, params operations.SourceListParams) middleware.Responder
}

type AutopilotAPI interface {
	OperationsAPI
	// Initialize is called during handler creation to perform and changes during startup
	Initialize() error

	// InitializeContext is call before the api methods are executed
	InitializeContext(principal provider.AuthToken, r *http.Request) (*autopilot.Context, error)
}

// Config is configuration for Handler
type Config struct {
	AutopilotAPI
	Logger *logrus.Logger
	// InnerMiddleware is for the handler executors. These do not apply to the swagger.json document.
	// The middleware executes after routing but before authentication, binding and validation
	InnerMiddleware func(http.Handler) http.Handler
}

// Handler returns an http.Handler given the handler configuration
// It mounts all the business logic implementers in the right routing.
func Handler(c Config) (http.Handler, error) {
	spec, err := loads.Analyzed(swaggerCopy(SwaggerJSON), "")
	if err != nil {
		return nil, fmt.Errorf("analyze swagger: %v", err)
	}
	api := operations.NewAutopilotAPI(spec)
	api.ServeError = errors.ServeError
	api.Logger = c.Logger.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.SourceListHandler = operations.SourceListHandlerFunc(func(params operations.SourceListParams) middleware.Responder {
		ctx, err := c.InitializeContext(principal, params.HTTPRequest)
		if err != nil {
			return sparks.NewError(err)
		}
		return c.AutopilotAPI.SourceList(ctx, params)
	})
	api.ServerShutdown = func() {}

	logMiddleware := func(handler http.Handler) http.Handler {
		handlePanic := recover.New(&recover.Options{
			Log: c.Logger.Error,
		})

		logViaLogrus := interpose.NegroniLogrus()

		if c.InnerMiddleware != nil {
			handler = c.InnerMiddleware(handler)
		}

		return handlePanic(logViaLogrus(handler))
	}

	if err := c.AutopilotAPI.Initialize(); err != nil {
		return nil, err
	}

	return api.Serve(logMiddleware), nil
}

// swaggerCopy copies the swagger json to prevent data races in runtime
func swaggerCopy(orig json.RawMessage) json.RawMessage {
	c := make(json.RawMessage, len(orig))
	copy(c, orig)
	return c
}
