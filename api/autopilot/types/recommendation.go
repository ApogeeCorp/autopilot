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
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Recommendation A recommendation is a list of recommended arbitrations to be emitted by the system
//
// swagger:model Recommendation
type Recommendation struct {

	// The recommendation values mapping rule.name -> formatted proposal
	Proposals []*Proposal `json:"proposals"`

	// The recommendation timestamp
	// Format: date-time
	Timestamp strfmt.DateTime `json:"timestamp,omitempty"`
}

// Validate validates this recommendation
func (m *Recommendation) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateProposals(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Recommendation) validateProposals(formats strfmt.Registry) error {

	if swag.IsZero(m.Proposals) { // not required
		return nil
	}

	for i := 0; i < len(m.Proposals); i++ {
		if swag.IsZero(m.Proposals[i]) { // not required
			continue
		}

		if m.Proposals[i] != nil {
			if err := m.Proposals[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("proposals" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Recommendation) validateTimestamp(formats strfmt.Registry) error {

	if swag.IsZero(m.Timestamp) { // not required
		return nil
	}

	if err := validate.FormatOf("timestamp", "body", "date-time", m.Timestamp.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Recommendation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Recommendation) UnmarshalBinary(b []byte) error {
	var res Recommendation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
