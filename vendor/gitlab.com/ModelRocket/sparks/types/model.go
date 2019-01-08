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
	strfmt "github.com/go-openapi/strfmt"
	"github.com/lib/pq/hstore"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Model is the base datebase model type used by Model Rocket
type Model struct {
	// Universally unique identifier
	// Read Only: true
	ID strfmt.UUID4 `json:"uuid,omitempty"`

	// Date and time of object creation
	// Read Only: true
	CreatedAt strfmt.DateTime `json:"_created_at,omitempty" sql:"column:_created_at"`

	// User uuid that created the object
	// Read Only: true
	CreatedBy string `json:"_created_by,omitempty" sql:"column:_created_by"`

	// Last date of object modification
	// Read Only: true
	UpdatedAt strfmt.DateTime `json:"_updated_at,omitempty" sql:"column:_updated_at"`

	// User that updated the object
	// Read Only: true
	UpdatedBy string `json:"_updated_by,omitempty" sql:"column:_updated_by"`

	// meta
	Meta *hstore.Hstore `json:"meta,omitempty"`

	// tags
	Tags StringArray `json:"tags,omitempty"`
}

// Validate validates this base model
func (m *Model) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Model) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Model) UnmarshalBinary(b []byte) error {
	var res Model
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
