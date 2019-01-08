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

import "github.com/spf13/cast"

// Value defines a value wrapper
type Value struct {
	v interface{}
}

// NewValue returns a new value with cast helpers
func NewValue(v interface{}) Value {
	return Value{v: v}
}

// String casts the value to a string
func (v Value) String() string {
	return cast.ToString(v.v)
}

// StringPtr casts the value to a string pointer
func (v Value) StringPtr() *string {
	tmp := cast.ToString(v.v)
	return &tmp
}

// StringSlice casts the value to a string slice
func (v Value) StringSlice() []string {
	return cast.ToStringSlice(v.v)
}

// Bool casts the value to a bool
func (v Value) Bool() bool {
	return cast.ToBool(v.v)
}

// Int64 casts the value to an int64
func (v Value) Int64() int64 {
	return cast.ToInt64(v.v)
}

// Float64 casts the value to a float64
func (v Value) Float64() float64 {
	return cast.ToFloat64(v.v)
}
