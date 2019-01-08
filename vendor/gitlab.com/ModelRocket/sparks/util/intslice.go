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

package util

import "fmt"

// Int64Slicer is the sprocket Int64Slicer interface
type Int64Slicer interface {
	ContainsAll(values ...int64) bool
	ContainsAny(values ...int64) bool
	Remove(r int64) Int64Slicer
	Slice() []int64
}

type int64slice []int64

// Int64Slice returns a sprocket Int64Slice
func Int64Slice(s []int64) Int64Slicer {
	return int64slice(s)
}

// ContainsAll returns true if the slice contains all of the values
func (s int64slice) ContainsAll(values ...int64) bool {
	sm := s.Map()

	for _, v := range values {
		if _, ok := sm[fmt.Sprint(v)]; !ok {
			return false
		}
	}

	return true
}

func (s int64slice) Slice() []int64 {
	return []int64(s)
}

// ContainsAny returns true if the slice contains any of the values
func (s int64slice) ContainsAny(values ...int64) bool {
	sm := s.Map()

	for _, v := range values {
		if _, ok := sm[fmt.Sprint(v)]; ok {
			return true
		}
	}

	return false
}

// Map converts a int64 slice to a map
func (s int64slice) Map(filter ...MapFilter) map[string]interface{} {
	return SliceToMap(s.Slice(), filter...)
}

// Remove removes a int64 from the slice
func (s int64slice) Remove(r int64) Int64Slicer {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
