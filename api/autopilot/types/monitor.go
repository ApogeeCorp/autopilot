// Code generated by go-swagger (hiro); DO NOT EDIT.

// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.
//

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Monitor A monitor executes and analyzes rules at the specified interval, emmiting alerts
//
// swagger:model Monitor
type Monitor struct {

	// The interval to monitor
	Interval *string `json:"interval,omitempty"`

	// The monitor name
	Name string `json:"name,omitempty"`

	// the monitor provider additional params
	Params map[string]interface{} `json:"params,omitempty"`

	// The provider to monitor
	Provider string `json:"provider,omitempty"`

	// The rules to execute on the provider
	Rules []string `json:"rules"`
}

// Validate validates this monitor
func (m *Monitor) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Monitor) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Monitor) UnmarshalBinary(b []byte) error {
	var res Monitor
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
