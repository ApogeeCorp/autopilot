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

import "strings"

// StrSlicer is the sprocket StrSlicer interface
type StrSlicer interface {
	ContainsAll(values ...string) bool
	ContainsAny(values ...string) bool
	ContainsStringIgnoreCase(value string) bool
	Remove(r string) StrSlicer
	Slice() []string
	Length() int
}

type strslice []string

// StrSlice returns a sprocket StrSlice
func StrSlice(s []string) StrSlicer {
	return strslice(s)
}

// ContainsAll returns true if the slice contains all of the values
func (s strslice) ContainsAll(values ...string) bool {
	sm := s.Map()

	for _, v := range values {
		if _, ok := sm[v]; !ok {
			return false
		}
	}

	return true
}

func (s strslice) ContainsStringIgnoreCase(value string) bool {
	for _, v := range s {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}

func (s strslice) Slice() []string {
	return []string(s)
}

func (s strslice) Length() int {
	return len([]string(s))
}

// ContainsAny returns true if the slice contains any of the values
func (s strslice) ContainsAny(values ...string) bool {
	sm := s.Map()

	for _, v := range values {
		if _, ok := sm[v]; ok {
			return true
		}
	}

	return false
}

// Map converts a string slice to a map
func (s strslice) Map(filter ...MapFilter) map[string]interface{} {
	return SliceToMap(s.Slice(), filter...)
}

// Remove removes a string from the slice
func (s strslice) Remove(r string) StrSlicer {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// StringFilter filters a string
func StringFilter(ss []string, p func(string) string) []string {
	for i := range ss {
		ss[i] = p(ss[i])
	}
	return ss
}
