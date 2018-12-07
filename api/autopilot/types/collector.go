// Code generated by go-swagger (hiro); DO NOT EDIT.

// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Collector A collector pulls data from a telemetry source, parses,
// and reformats the data to be consumed by the autopilot engine.
//
// swagger:model Collector
type Collector struct {

	// The emitters to use after processing the samples
	Emitters []string `json:"emitters"`

	// The collector name
	Name string `json:"name,omitempty"`

	// json data object
	Params map[string]interface{} `json:"params,omitempty"`

	// The interval the collector will run at
	ScheduleInterval *string `json:"schedule_interval,omitempty"`

	// type
	Type CollectorType `json:"type,omitempty"`

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

func (m *Collector) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	if err := m.Type.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("type")
		}
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
