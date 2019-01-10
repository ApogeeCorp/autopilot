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

	"github.com/go-openapi/strfmt"
	"github.com/lib/pq"
)

// UUIDArray is a proper string array type
type UUIDArray []strfmt.UUID

// Strings returns a string slice from the uuid array
func (a UUIDArray) Strings() []string {
	rval := make([]string, 0)
	for _, id := range a {
		rval = append(rval, id.String())
	}
	return rval
}

// Scan implements the sql.Scanner interface
func (a *UUIDArray) Scan(src interface{}) error {
	return pq.GenericArray{A: a}.Scan(src)
}

// Value implements the driver.Valuer interface.
func (a UUIDArray) Value() (driver.Value, error) {
	return pq.GenericArray{A: a}.Value()
}

// Validate handles the strfmt validation for the StringArray object
func (*UUIDArray) Validate(formats strfmt.Registry) error {
	return nil
}
