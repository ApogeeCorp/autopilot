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
	"compress/flate"
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/pquerna/ffjson/ffjson"
)

type (
	// ResponseOptions defines responder options
	ResponseOptions struct {
		Request      *http.Request
		Encoding     string
		ContentType  string
		Binary       bool
		Acceleration bool
	}

	// Response defines a simple middleware.Responder struct
	Response struct {
		payload       interface{}
		status        int
		options       ResponseOptions
		body          []byte
		headerWritten bool
		headers       http.Header
	}
)

var (
	// NoContentResponse is a no content HTTP success response
	NoContentResponse = NewResponse().SetStatus(http.StatusNoContent)

	// SeeOtherResponse returns a redirect response
	SeeOtherResponse = func(loc string) *Response {
		return NewResponse().SetStatus(http.StatusSeeOther).SetHeader("Location", loc)
	}
)

// NewResponse returns a new response object
func NewResponse(options ...*ResponseOptions) *Response {
	ops := ResponseOptions{}
	if len(options) > 0 {
		ops = *options[0]
	}

	return &Response{
		status:  http.StatusOK,
		options: ops,
		headers: make(map[string][]string),
	}
}

// Copy copies a response
func (r *Response) Copy() *Response {
	body := make([]byte, len(r.body))
	copy(body, r.body)

	return &Response{
		payload: r.payload,
		status:  r.status,
		options: r.options,
		body:    body,
		headers: r.headers,
	}
}

// Length returns the response length
func (r *Response) Length() int {
	if r.body == nil {
		return 0
	}
	return len(r.body)
}

// Body returns the body
func (r *Response) Body() []byte {
	return r.body
}

// StatusCode returns the response http status code
func (r *Response) StatusCode() int {
	return r.status
}

// SetPayload sets the payload
func (r *Response) SetPayload(payload interface{}) *Response {
	r.payload = payload
	return r
}

// SetOptions sets the options for a response
func (r *Response) SetOptions(options ...*ResponseOptions) *Response {
	ops := ResponseOptions{}
	if len(options) > 0 {
		// TODO: merge the options
		ops = *options[0]
	}

	r.options = ops

	return r
}

// SetStatus sets the response http status
func (r *Response) SetStatus(s int) *Response {
	r.status = s
	return r
}

// SetAcceleration enables json acceleration
func (r *Response) SetAcceleration(acc bool) *Response {
	r.options.Acceleration = acc
	return r
}

// Header implements the http.ResponseWriter
func (r *Response) Header() http.Header {
	return r.headers
}

// SetHeader sets additional response headers
func (r *Response) SetHeader(name, value string) *Response {
	r.headers.Set(name, value)
	return r
}

// SetHeaderMap sets additional response headers
func (r *Response) SetHeaderMap(vals map[string]string) *Response {
	for k, v := range vals {
		r.headers.Set(k, v)
	}
	return r
}

// SetHeaders sets all of the headers
func (r *Response) SetHeaders(h http.Header) *Response {
	r.headers = h
	return r
}

// Write implements the http.ResponseWriter.Write() method
func (r *Response) Write(body []byte) (int, error) {
	r.body = append(r.body, body...)

	if !r.headerWritten {
		r.WriteHeader(http.StatusOK)
	}

	return len(body), nil
}

// WriteHeader implements the http.ResponseWriter.WriteHeader() method
func (r *Response) WriteHeader(statusCode int) {
	r.status = statusCode

	ct := r.Header().Get("Content-Type")
	if ct == "" {
		ct = http.DetectContentType([]byte(r.body))
		r.Header().Set("Content-Type", ct)
	}

	r.headerWritten = true
}

// WriteResponse implements the middleware.Responder.WriteResponse() method
func (r *Response) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
	buf := bytes.NewBuffer(r.body)

	// default encoding is gzip
	if r.options.Encoding == "*" {
		r.options.Encoding = "gzip"
	}

	if r.options.ContentType != "" {
		rw.Header().Add("content-type", r.options.ContentType)
	}

	for k, s := range r.headers {
		for _, v := range s {
			rw.Header().Add(k, v)
		}
	}

	if strings.Contains(r.options.Encoding, "gzip") {
		rw.Header().Add("content-encoding", "gzip")
		comp := gzip.NewWriter(buf)
		if r.options.Binary {
			if data, ok := r.payload.([]byte); ok {
				rw.Header().Add("X-Content-Length", fmt.Sprint(len(data)))
				r.body = data
				comp.Write(data)
			}
		} else {
			var b bytes.Buffer
			out := bufio.NewWriter(&b)
			if err := pr.Produce(out, r.payload); err != nil {
				if err != nil {
					panic(err)
				}
			}
			data := b.Bytes()

			rw.Header().Add("X-Content-Length", fmt.Sprint(len(data)))
			r.body = data
			comp.Write(data)
		}
		comp.Flush()
		comp.Close()
	} else if strings.Contains(r.options.Encoding, "deflate") {
		rw.Header().Add("content-encoding", "deflate")
		comp, _ := flate.NewWriter(buf, 9)

		if r.options.Binary {
			if data, ok := r.payload.([]byte); ok {
				rw.Header().Add("X-Content-Length", fmt.Sprint(len(data)))
				r.body = data
				comp.Write(data)
			}
		} else {
			var b bytes.Buffer
			out := bufio.NewWriter(&b)
			if err := pr.Produce(out, r.payload); err != nil {
				if err != nil {
					panic(err)
				}
			}
			data := b.Bytes()
			rw.Header().Add("X-Content-Length", fmt.Sprint(len(data)))
			r.body = data
			comp.Write(data)
		}
		comp.Flush()
		comp.Close()
	}

	if buf.Len() > 0 {
		rw.Write(buf.Bytes())
	} else if r.payload != nil {
		if r.options.Binary {
			if data, ok := r.payload.([]byte); ok {
				rw.Header().Add("X-Content-Length", fmt.Sprint(len(data)))
				rw.Write(data)
				r.body = data
				return
			}
		}
		if r.options.Acceleration {
			buf, err := ffjson.Marshal(r.payload)
			if err == nil {
				r.body = buf
				rw.Write(buf)
				return
			}
		}
		var b bytes.Buffer
		out := bufio.NewWriter(&b)
		if err := pr.Produce(out, r.payload); err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(err.Error()))
		} else {
			out.Flush()
			rw.WriteHeader(r.status)
			r.body = b.Bytes()
			rw.Write(r.body)
		}
	} else {
		rw.WriteHeader(r.status)
	}
}
