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
	"strconv"
	"strings"
	"unicode"
)

// Params is a simple type alias for a Map
type Params = StringMap

// ParseStringParams parses a parameter string and returns a params object
// Example string: foo=bar aparam="value"
func ParseStringParams(s string) Params {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	// splitting string by space but considering quoted section
	items := strings.FieldsFunc(s, f)

	// create and fill the map
	m := make(Params)
	for _, item := range items {
		x := strings.SplitN(item, "=", 2)
		update := func(val interface{}) {
			if m[x[0]] != nil {
				if arr, ok := m[x[0]].([]interface{}); ok {
					m[x[0]] = append(arr, val)
				} else {
					tmp := m[x[0]]
					m[x[0]] = []interface{}{tmp, val}
				}
			} else {
				m[x[0]] = val
			}
		}

		if len(x) > 1 {
			if v, err := strconv.Unquote(x[1]); err == nil {
				update(v)
			} else {
				update(x[1])
			}
		} else {
			update(true)
		}
	}
	return m
}
