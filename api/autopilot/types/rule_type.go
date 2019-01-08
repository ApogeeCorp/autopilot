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
	"github.com/go-openapi/validate"
)

// RuleType Types of rules
// swagger:model RuleType
type RuleType string

type RuleTypeScalar int

var (

	// RuleTypePromql captures enum value "promql"
	RuleTypePromql RuleType = "promql"

	// RuleTypeApql captures enum value "apql"
	RuleTypeApql RuleType = "apql"

	// RuleTypeAplearn captures enum value "aplearn"
	RuleTypeAplearn RuleType = "aplearn"

	RuleTypeScalarLookup = map[RuleType]RuleTypeScalar{

		RuleTypePromql: RuleTypePromqlScalar,

		RuleTypeApql: RuleTypeApqlScalar,

		RuleTypeAplearn: RuleTypeAplearnScalar,
	}
)

const (
	RuleTypePromqlScalar RuleTypeScalar = 0

	RuleTypeApqlScalar RuleTypeScalar = 1

	RuleTypeAplearnScalar RuleTypeScalar = 2
)

func (m RuleType) ScalarValue() RuleTypeScalar {
	return RuleTypeScalarLookup[m]
}

// for schema
var ruleTypeEnum []interface{}

func init() {
	var res []RuleType
	if err := json.Unmarshal([]byte(`["promql","apql","aplearn"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		ruleTypeEnum = append(ruleTypeEnum, v)
	}
}

func (m RuleType) validateRuleTypeEnum(path, location string, value RuleType) error {
	if err := validate.Enum(path, location, value, ruleTypeEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this rule type
func (m RuleType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateRuleTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}