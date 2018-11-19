// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewRecommendationsGetParams creates a new RecommendationsGetParams object
// with the default values initialized.
func NewRecommendationsGetParams() RecommendationsGetParams {

	var (
		// initialize parameters with default values

		typeVarDefault = string("prometheus")
	)

	return RecommendationsGetParams{
		Type: &typeVarDefault,
	}
}

// RecommendationsGetParams contains all the bound params for the recommendations get operation
// typically these are obtained from a http.Request
//
// swagger:parameters recommendationsGet
type RecommendationsGetParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The sample to create
	  Required: true
	  In: formData
	*/
	Sample io.ReadCloser
	/*The provider type to process the sample with
	  In: query
	  Default: "prometheus"
	*/
	Type *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewRecommendationsGetParams() beforehand.
func (o *RecommendationsGetParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		if err != http.ErrNotMultipart {
			return errors.New(400, "%v", err)
		} else if err := r.ParseForm(); err != nil {
			return errors.New(400, "%v", err)
		}
	}

	sample, sampleHeader, err := r.FormFile("sample")
	if err != nil {
		res = append(res, errors.New(400, "reading file %q failed: %v", "sample", err))
	} else if err := o.bindSample(sample, sampleHeader); err != nil {
		// Required: true
		res = append(res, err)
	} else {
		o.Sample = &runtime.File{Data: sample, Header: sampleHeader}
	}

	qType, qhkType, _ := qs.GetOK("type")
	if err := o.bindType(qType, qhkType, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindSample binds file parameter Sample.
//
// The only supported validations on files are MinLength and MaxLength
func (o *RecommendationsGetParams) bindSample(file multipart.File, header *multipart.FileHeader) error {
	return nil
}

// bindType binds and validates parameter Type from query.
func (o *RecommendationsGetParams) bindType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewRecommendationsGetParams()
		return nil
	}

	o.Type = &raw

	if err := o.validateType(formats); err != nil {
		return err
	}

	return nil
}

// validateType carries on validations for parameter Type
func (o *RecommendationsGetParams) validateType(formats strfmt.Registry) error {

	if err := validate.Enum("type", "query", *o.Type, []interface{}{"prometheus"}); err != nil {
		return err
	}

	return nil
}
