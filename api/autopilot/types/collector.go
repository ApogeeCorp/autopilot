// Code generated by go-swagger (hiro); DO NOT EDIT.

// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Collector A collector pulls data from a telemetry source, parses,
// and reformats the data to be consumed by the autopilot engine.
//
// swagger:model Collector
type Collector struct {

	// The emitters to use after processing the samples
	Emitters []string `json:"emitters"`

	// The interval the collector will run at
	Interval string `json:"interval,omitempty"`

	// The collector name
	Name string `json:"name,omitempty"`

	// json data object
	Params map[string]interface{} `json:"params,omitempty"`

	// The collector client to use
	// Enum: [prometheus]
	Type string `json:"type,omitempty"`

	// The collector url
	URL string `json:"url,omitempty"`
}

// Validate validates this collector
func (m *Collector) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var collectorTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["prometheus"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		collectorTypeTypePropEnum = append(collectorTypeTypePropEnum, v)
	}
}

const (

	// CollectorTypePrometheus captures enum value "prometheus"
	CollectorTypePrometheus string = "prometheus"
)

// prop value enum
func (m *Collector) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, collectorTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Collector) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Collector) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Collector) UnmarshalBinary(b []byte) error {
	var res Collector
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
