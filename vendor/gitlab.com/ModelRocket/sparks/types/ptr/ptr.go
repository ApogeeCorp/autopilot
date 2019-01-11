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

// Package ptr provides some helpers for converting types to pointers
package ptr

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Int64 is a safe int64 pointer wrapper
func Int64(val interface{}) *int64 {
	rval := int64(0)
	if val == nil {
		return &rval
	}
	switch v := val.(type) {
	case *int64:
		return v
	case int64:
		rval = v
	case int:
		rval = int64(v)
	case int32:
		rval = int64(v)
	}
	return &rval
}

// Float64 is a safe float64 pointer wrapper
func Float64(val interface{}) *float64 {
	rval := float64(0)
	if val == nil {
		return &rval
	}
	switch v := val.(type) {
	case float64:
		rval = v
	case *float64:
		return v
	case *int64:
		rval = float64(*v)
	case int64:
		rval = float64(v)
	case int:
		rval = float64(v)
	case int32:
		rval = float64(v)
	}
	return &rval
}

// Duration is a safe time.Duration pointer wrapper
func Duration(val interface{}) *time.Duration {
	rval := time.Duration(0)
	if val == nil {
		return &rval
	}
	switch v := val.(type) {
	case *int64:
		rval = time.Duration(*v)
	case int64:
		rval = time.Duration(v)
	case time.Duration:
		return &v
	case *time.Duration:
		return v
	}
	return &rval
}

// Bool returns a bool pointer
func Bool(val interface{}, def ...bool) *bool {
	defVal := false
	if len(def) > 0 {
		defVal = def[0]
	}
	if val == nil {
		return &defVal
	}

	switch b := val.(type) {
	case bool:
		return &b
	case *bool:
		if b == nil {
			return &defVal
		}
		return b
	case string:
		r := false
		tmp, err := strconv.ParseBool(b)
		if err != nil {
			return &r
		}
		return &tmp
	case *string:
		if b == nil {
			return &defVal
		}
		r := false
		tmp, err := strconv.ParseBool(*b)
		if err != nil {
			return &r
		}
		return &tmp
	}

	return &defVal
}

// String is a safe string pointer wrapper
func String(val interface{}) *string {
	if val == nil {
		rval := ""
		return &rval
	}
	switch v := val.(type) {
	case *string:
		return v
	case string:
		return &v
	case []byte:
		rval := string(v)
		return &rval
	}

	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	tmp := v.Interface()
	rval := fmt.Sprint(tmp)
	return &rval
}
