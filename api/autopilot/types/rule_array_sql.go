// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Scan implements the sql.Scanner interface
func (r *RuleArray) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Invalid data type for RuleArray")
	}
	return json.Unmarshal(bytes, r)
}

// Value implements the driver.Valuer interface.
func (r RuleArray) Value() (driver.Value, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}
