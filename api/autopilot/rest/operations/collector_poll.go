// Code generated by hiro; DO NOT EDIT.

// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.
//

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// CollectorPollHandlerFunc turns a function with the right signature into a collector poll handler
type CollectorPollHandlerFunc func(CollectorPollParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn CollectorPollHandlerFunc) Handle(params CollectorPollParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// CollectorPollHandler interface for that can handle valid collector poll params
type CollectorPollHandler interface {
	Handle(CollectorPollParams, interface{}) middleware.Responder
}

// NewCollectorPoll creates a new http.Handler for the collector poll operation
func NewCollectorPoll(ctx *middleware.Context, handler CollectorPollHandler) *CollectorPoll {
	return &CollectorPoll{Context: ctx, Handler: handler}
}

/*CollectorPoll swagger:route GET /collectors/{collector}/poll collectorPoll

Poll a collector

Poll a collector for the given period directly

*/
type CollectorPoll struct {
	Context *middleware.Context
	Handler CollectorPollHandler
}

func (o *CollectorPoll) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCollectorPollParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
