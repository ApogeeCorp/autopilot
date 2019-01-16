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
	"reflect"
)

// Slicer is a generic slice manipulation interface helper
type Slicer interface {
	Append(vals ...interface{})
	Insert(index int, val interface{})
	IndexOf(val interface{}) int
	Remove(index int)
	Map(out interface{}, filters ...MapFilter)
	Contains(values ...interface{}) bool
	ContainsAny(values ...interface{}) bool
}

type slice struct {
	v reflect.Value
}

// MapFilter returns true and the key and value, or false if the value should be skipped
type MapFilter func(index int, val interface{}) (string, interface{}, bool)

// Slice returns a slicer object from a pointer to a slice
// 	v := make([]string, 0)
//	util.Slice(&v).Append("foo")
func Slice(s interface{}) Slicer {
	return &slice{
		v: reflect.ValueOf(s).Elem(),
	}
}

func (s *slice) Append(vals ...interface{}) {
	for _, val := range vals {
		s.v.Set(reflect.Append(s.v, reflect.ValueOf(val)))
	}
}

func (s *slice) Insert(index int, val interface{}) {
	s.v.Set(reflect.AppendSlice(s.v.Slice(0, index+1), s.v.Slice(index, s.v.Len())))
	s.v.Index(index).Set(reflect.ValueOf(val))
}

func (s *slice) Remove(index int) {
	s.v.Set(reflect.AppendSlice(s.v.Slice(0, index), s.v.Slice(index+1, s.v.Len())))
}

// Map maps the slice to the destination using the filters
// Destination can be a slice or a map
func (s *slice) Map(out interface{}, filters ...MapFilter) {
	rval := reflect.ValueOf(out).Elem()

	for _, filter := range filters {
		for i := 0; i < s.v.Len(); i++ {
			if k, v, skip := filter(i, s.v.Index(i).Interface()); !skip {
				switch rval.Kind() {
				case reflect.Slice:
					Slice(out).Append(v)
				case reflect.Map:
					Map(out).Set(k, v)
				}
			}
		}
	}
}

func (s *slice) IndexOf(val interface{}) int {
	for i := 0; i < s.v.Len(); i++ {
		if reflect.DeepEqual(s.v.Index(i), reflect.ValueOf(val)) {
			return i
		}
	}

	return -1
}

func (s *slice) Contains(values ...interface{}) bool {
	for _, val := range values {
		if s.IndexOf(val) < 0 {
			return false
		}
	}
	return true
}

func (s *slice) ContainsAny(values ...interface{}) bool {
	for _, val := range values {
		if s.IndexOf(val) >= 0 {
			return true
		}
	}
	return false
}
