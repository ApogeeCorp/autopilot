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
	"bufio"
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"gitlab.com/ModelRocket/sparks/types/ptr"
)

var (
	// ErrNotAuthorized defines a not authorized error
	ErrNotAuthorized = NewError().Status(http.StatusUnauthorized).Format("not authorized")

	// ErrInvalidParameter defines an invalid parameter error
	ErrInvalidParameter = NewError().Status(http.StatusBadRequest).Format("invalid parameter")

	// ErrNotImplemented defines a not implemented error
	ErrNotImplemented = NewError().Status(http.StatusNotImplemented).FormatFunc("method '%s' not implemented")

	// ErrResourceConflict is returned when an object conflicts with one that already exists
	ErrResourceConflict = NewError().Status(http.StatusConflict).FormatFunc("'%s' conflicts with one that already exists")

	// ErrNotFound is is returned when an object was not found
	ErrNotFound = NewError().Status(http.StatusNotFound).FormatFunc("'%s' not found")

	// ErrServerError is reserved for server errors
	ErrServerError = NewError().Status(http.StatusInternalServerError).FormatFunc("server error: %s")
)

type (
	// Error is the common error struct compatible with swagger response types
	Error struct {
		Code    *string           `json:"code,omitempty"`
		Message *string           `json:"message,omitempty"`
		Detail  map[string]string `json:"detail,omitempty"`
		status  int
		err     error
		length  int
	}
)

// NewError returns a new models.Error from a go error
func NewError(err ...error) *Error {
	var msg *string
	var code *string
	var e error

	detail := make(map[string]string)
	status := http.StatusInternalServerError

	if len(err) > 0 {
		e = err[0]

		switch typedErr := (e).(type) {
		case *pq.Error:
			msg = ptr.String(typedErr.Message)
			code = ptr.String(typedErr.Code)
			detail = map[string]string{
				"reason": typedErr.Detail,
				"column": typedErr.Column,
			}

			switch typedErr.Code {
			case "23505":
				status = http.StatusConflict
			}

		case *Error:
			return typedErr

		default:
			if e == sql.ErrNoRows || typedErr == gorm.ErrRecordNotFound {
				status = http.StatusNotFound
			}
			msg = ptr.String(e.Error())
		}
	}
	return &Error{
		Message: msg,
		Code:    code,
		Detail:  detail,
		status:  status,
		err:     e,
	}
}

// Format formats the error message
func (e *Error) Format(msg string, args ...interface{}) *Error {
	e.Message = ptr.String(fmt.Sprintf(msg, args...))
	return e
}

// FormatFunc func formats a error message functor
func (e *Error) FormatFunc(msg string) func(args ...interface{}) *Error {
	return func(args ...interface{}) *Error {
		if len(args) > 0 {
			val := reflect.ValueOf(args[0])
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			if val.Kind() == reflect.Struct {
				args[0] = val.Type().String()
			}
		}
		return e.Format(msg, args...)
	}
}

// Reason sets the detail for the error
func (e *Error) Reason(args ...interface{}) *Error {
	if len(args) == 0 {
		return e
	}
	if len(args) == 1 {
		switch a := args[0].(type) {
		case error:
			e.Detail["reason"] = a.Error()
			return e
		case string:
			e.Detail["reason"] = a
			return e
		case []interface{}:
			for i, v := range a {
				e.Detail[fmt.Sprintf("error_%d", i)] = fmt.Sprint(v)
			}
		}
		val := reflect.ValueOf(args[0])
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.IsValid() && val.Kind() == reflect.Struct {
			if e.err == gorm.ErrRecordNotFound {
				e.Message = ptr.String(fmt.Sprintf("%s not found", val.Type().String()))
			} else {
				e.Detail["model"] = val.Type().String()
			}
		}
	} else {
		if len(args) > 1 {
			f, ok := args[0].(string)
			if ok {
				e.Detail["reason"] = fmt.Sprintf(f, args[1:]...)
			} else {
				e.Detail["reason"] = fmt.Sprint(args...)
			}
		} else {
			e.Detail["reason"] = fmt.Sprint(args[0])
		}
	}

	return e
}

// Status sets the http status
func (e *Error) Status(status int) *Error {
	e.status = status
	return e
}

// ErrorCode sets the error code
func (e *Error) ErrorCode(code string) *Error {
	e.Code = ptr.String(code)
	return e
}

// Header returns the http header
func (e *Error) Header() http.Header {
	return make(http.Header)
}

// StatusCode returns the http statuscode
func (e *Error) StatusCode() int {
	return e.status
}

// BodyLength returns the http body length in bytes
func (e *Error) BodyLength() int {
	return e.length
}

// Err returns the go error
func (e *Error) Error() string {
	msg := ""
	if e.Message != nil {
		msg = *e.Message
	}

	if reason, ok := e.Detail["reason"]; ok {
		msg = msg + ": " + reason
	}

	return msg
}

func (e *Error) String() string {
	return e.Error()
}

// WriteResponse implements the middleware.Responder interface
func (e *Error) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
	var b bytes.Buffer
	pr = runtime.JSONProducer()
	out := bufio.NewWriter(&b)
	rw.Header().Set("content-type", "application/json")
	if err := pr.Produce(out, e); err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
	} else {
		out.Flush()
		rw.WriteHeader(e.status)
		rw.Write(b.Bytes())
		e.length = b.Len()
	}
}

// Err returns the underlining error
func (e *Error) Err() error {
	return e.err
}

// Validate handles the strfmt validation for the StringArray object
func (*Error) Validate(formats strfmt.Registry) error {
	return nil
}

// Scan implements the sql.Scanner interface
func (e *Error) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Invalid data type for Error")
	}
	return json.Unmarshal(bytes, e)
}

// Value implements the driver.Valuer interface.
func (e Error) Value() (driver.Value, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}
