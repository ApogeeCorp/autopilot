// Code generated by go-swagger (hiro); DO NOT EDIT.

// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/validate"
)

// Collector A collector pulls data from a telemetry source, parses,
// and reformats the data to be consumed by the autopilot engine.
//
// swagger:discriminator Collector type
type Collector interface {
	runtime.Validatable

	// The emitters to use after processing the samples
	Emitters() []string
	SetEmitters([]string)

	// The collector name
	Name() string
	SetName(string)

	// The interval the collector will
	ScheduleInterval() *string
	SetScheduleInterval(*string)

	// type
	Type() CollectorType
	SetType(CollectorType)

	// The collector url
	URL() string
	SetURL(string)
}

type collector struct {
	emittersField []string

	nameField string

	scheduleIntervalField *string

	typeField CollectorType

	urlField string
}

// Emitters gets the emitters of this polymorphic type
func (m *collector) Emitters() []string {
	return m.emittersField
}

// SetEmitters sets the emitters of this polymorphic type
func (m *collector) SetEmitters(val []string) {
	m.emittersField = val
}

// Name gets the name of this polymorphic type
func (m *collector) Name() string {
	return m.nameField
}

// SetName sets the name of this polymorphic type
func (m *collector) SetName(val string) {
	m.nameField = val
}

// ScheduleInterval gets the schedule interval of this polymorphic type
func (m *collector) ScheduleInterval() *string {
	return m.scheduleIntervalField
}

// SetScheduleInterval sets the schedule interval of this polymorphic type
func (m *collector) SetScheduleInterval(val *string) {
	m.scheduleIntervalField = val
}

// Type gets the type of this polymorphic type
func (m *collector) Type() CollectorType {
	return "Collector"
}

// SetType sets the type of this polymorphic type
func (m *collector) SetType(val CollectorType) {

}

// URL gets the url of this polymorphic type
func (m *collector) URL() string {
	return m.urlField
}

// SetURL sets the url of this polymorphic type
func (m *collector) SetURL(val string) {
	m.urlField = val
}

// UnmarshalCollectorSlice unmarshals polymorphic slices of Collector
func UnmarshalCollectorSlice(reader io.Reader, consumer runtime.Consumer) ([]Collector, error) {
	var elements []json.RawMessage
	if err := consumer.Consume(reader, &elements); err != nil {
		return nil, err
	}

	var result []Collector
	for _, element := range elements {
		obj, err := unmarshalCollector(element, consumer)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
	}
	return result, nil
}

// UnmarshalCollector unmarshals polymorphic Collector
func UnmarshalCollector(reader io.Reader, consumer runtime.Consumer) (Collector, error) {
	// we need to read this twice, so first into a buffer
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return unmarshalCollector(data, consumer)
}

func unmarshalCollector(data []byte, consumer runtime.Consumer) (Collector, error) {
	buf := bytes.NewBuffer(data)
	buf2 := bytes.NewBuffer(data)

	// the first time this is read is to fetch the value of the type property.
	var getType struct {
		Type string `json:"type"`
	}
	if err := consumer.Consume(buf, &getType); err != nil {
		return nil, err
	}

	if err := validate.RequiredString("type", "body", getType.Type); err != nil {
		return nil, err
	}

	// The value of type is used to determine which type to create and unmarshal the data into
	switch getType.Type {
	case "Collector":
		var result collector
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil

	case "Prometheus":
		var result Prometheus
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil

	}
	return nil, errors.New(422, "invalid type value: %q", getType.Type)

}

// Validate validates this collector
func (m *collector) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
