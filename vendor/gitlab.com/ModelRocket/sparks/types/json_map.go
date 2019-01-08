/*************************************************************************
 * MIT License
 * Copyright (c) 2018 Model Rocket
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cast"
)

// JSONMap is a helper for the postgres jsonb datatype
type JSONMap map[string]interface{}

// IsSet returns true if the parameter is set
func (p JSONMap) IsSet(key string) bool {
	_, ok := p[key]
	return ok
}

// IsSetV returns true if the parameter is set and the value
func (p JSONMap) IsSetV(key string) (Value, bool) {
	v, ok := p[key]
	return NewValue(v), ok
}

// Set sets a value in the map
func (p JSONMap) Set(key string, value interface{}) {
	p[key] = value
}

// Get returns the value and if the key is set
func (p JSONMap) Get(key string, def ...interface{}) (interface{}, bool) {
	var dval interface{}
	if len(def) > 0 {
		dval = def[0]
	}
	val, ok := p[strings.ToLower(key)]
	if !ok {
		return dval, false
	}

	return val, true
}

// Sub returns a sub JSONMap for the key
func (p JSONMap) Sub(key string) JSONMap {
	if tmp, ok := p[key]; ok {
		if p, ok := tmp.(map[string]interface{}); ok {
			return JSONMap(p)
		}
	}
	return JSONMap{}
}

// String returns a string value for the param, or the optional default
func (p JSONMap) String(key string, def ...string) string {
	rval, ok := p.Get(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return ""
	}

	return cast.ToString(rval)
}

// StringPtr returns a string ptr or nil
func (p JSONMap) StringPtr(key string, def ...string) *string {
	rval, ok := p.Get(key)
	if !ok {
		if len(def) > 0 {
			return &def[0]
		}
		return nil
	}
	tmp := cast.ToString(rval)
	return &tmp
}

// StringSlice returns a string value for the param, or the optional default
func (p JSONMap) StringSlice(key string) []string {
	rval, ok := p.Get(key)
	if !ok {
		return []string{}
	}

	return cast.ToStringSlice(rval)
}

// Bool parses and returns the boolean value of the parameter
func (p JSONMap) Bool(key string, def ...bool) bool {
	rval, ok := p.Get(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return false
	}

	return cast.ToBool(rval)
}

// Int64 returns the int value or 0 if not set
func (p JSONMap) Int64(key string, def ...int64) int64 {
	rval, ok := p.Get(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	if val, ok := rval.(json.Number); ok {
		if rval, err := val.Int64(); err == nil {
			return rval
		}
	}
	return cast.ToInt64(rval)
}

// Float64 returns the float value or 0 if not set
func (p JSONMap) Float64(key string, def ...float64) float64 {
	rval, ok := p.Get(key)
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	if val, ok := rval.(json.Number); ok {
		if rval, err := val.Float64(); err == nil {
			return rval
		}
	}
	return cast.ToFloat64(rval)
}

// Delete removes a key
func (p JSONMap) Delete(key string) {
	delete(p, key)
}

// Copy does a shallow copy
func (p JSONMap) Copy() JSONMap {
	rval := make(JSONMap)
	for k, v := range p {
		rval[k] = v
	}
	return rval
}

// Without removes the keys, returns a shallow copy
func (p JSONMap) Without(keys ...string) JSONMap {
	rval := p.Copy()
	for _, key := range keys {
		delete(rval, key)
	}
	return rval
}

// Validate handles the strfmt validation for the StringArray object
func (*JSONMap) Validate(formats strfmt.Registry) error {
	return nil
}

// Scan implements the sql.Scanner interface
func (p *JSONMap) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Invalid data type for JSONMap")
	}
	return json.Unmarshal(bytes, p)
}

// Value implements the driver.Valuer interface.
func (p JSONMap) Value() (driver.Value, error) {
	data, err := json.Marshal(map[string]interface{}(p))
	if err != nil {
		return nil, err
	}
	return string(data), nil
}
