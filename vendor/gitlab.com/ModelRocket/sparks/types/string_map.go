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

// StringMap is a wrapper on a map[string]interface{} that implements the Mapper interface
type StringMap map[string]interface{}

// IsSet returns true if the parameter is set
func (p StringMap) IsSet(key interface{}) bool {
	_, ok := p[cast.ToString(key)]
	return ok
}

// Set sets a value in the map
func (p StringMap) Set(key, value interface{}) {
	p[cast.ToString(key)] = value
}

// Get returns the value and if the key is set
func (p StringMap) Get(key interface{}) interface{} {
	val, ok := p[strings.ToLower(cast.ToString(key))]
	if !ok {
		return nil
	}

	return val
}

// Sub returns a sub StringMap for the key
func (p StringMap) Sub(key interface{}) Mapper {
	if tmp, ok := p[cast.ToString(key)]; ok {
		switch p := tmp.(type) {
		case map[string]interface{}:
			return StringMap(p)
		default:
			return Map(p)
		}
	}
	return StringMap{}
}

// Delete removes a key
func (p StringMap) Delete(key interface{}) {
	delete(p, cast.ToString(key))
}

// Copy does a shallow copy
func (p StringMap) Copy() Mapper {
	rval := make(StringMap)
	for k, v := range p {
		rval[k] = v
	}
	return rval
}

// Without removes the keys, returns a shallow copy
func (p StringMap) Without(keys ...interface{}) Mapper {
	rval := p.Copy().(StringMap)
	for _, key := range keys {
		rval.Delete(key)
	}
	return rval
}

// GetValue returns a Value and if the key is set
func (p StringMap) GetValue(key string) (Value, bool) {
	v, ok := p[key]
	return NewValue(v), ok
}

// String returns a string value for the param, or the optional default
func (p StringMap) String(key string, def ...string) string {
	rval := p.Get(key)
	if rval == nil {
		if len(def) > 0 {
			return def[0]
		}
		return ""
	}

	return cast.ToString(rval)
}

// StringPtr returns a string ptr or nil
func (p StringMap) StringPtr(key string, def ...string) *string {
	rval := p.Get(key)
	if rval == nil {
		if len(def) > 0 {
			return &def[0]
		}
		return nil
	}
	tmp := cast.ToString(rval)
	return &tmp
}

// StringSlice returns a string value for the param, or the optional default
func (p StringMap) StringSlice(key string) []string {
	rval := p.Get(key)
	if rval == nil {
		return []string{}
	}

	return cast.ToStringSlice(rval)
}

// Bool parses and returns the boolean value of the parameter
func (p StringMap) Bool(key string, def ...bool) bool {
	rval := p.Get(key)
	if rval == nil {
		if len(def) > 0 {
			return def[0]
		}
		return false
	}

	return cast.ToBool(rval)
}

// Int64 returns the int value or 0 if not set
func (p StringMap) Int64(key string, def ...int64) int64 {
	rval := p.Get(key)
	if rval == nil {
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
func (p StringMap) Float64(key string, def ...float64) float64 {
	rval := p.Get(key)
	if rval == nil {
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

// Validate handles the strfmt validation for the StringArray object
func (*StringMap) Validate(formats strfmt.Registry) error {
	return nil
}

// Scan implements the sql.Scanner interface
func (p *StringMap) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Invalid data type for StringMap")
	}
	return json.Unmarshal(bytes, p)
}

// Value implements the driver.Valuer interface.
func (p StringMap) Value() (driver.Value, error) {
	data, err := json.Marshal(map[string]interface{}(p))
	if err != nil {
		return nil, err
	}
	return string(data), nil
}
