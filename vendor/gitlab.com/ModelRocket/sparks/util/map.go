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

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/doublerebel/bellows"
	"github.com/fatih/structs"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
)

var (
	decoderHooks = mapstructure.ComposeDecodeHookFunc(mapStructDecoder)
)

// MapFilter returns true and the value, or false if the value should be skipped
type MapFilter func(key *string, val interface{}) (interface{}, bool)

func mapStructDecoder(src reflect.Type, dst reflect.Type, data interface{}) (interface{}, error) {
	if dst.Kind() == reflect.Ptr && src.Kind() == reflect.String {
		if dst.Elem().Kind() == reflect.Struct {
			out := reflect.New(dst.Elem()).Interface()
			if err := json.Unmarshal([]byte(data.(string)), out); err != nil {
				return nil, err
			}
			return out, nil
		}
	}

	if dst.Kind() == reflect.Struct && src.Kind() == reflect.String {
		out := reflect.New(dst).Interface()
		if err := json.Unmarshal([]byte(data.(string)), out); err != nil {
			return nil, err
		}
		val := reflect.ValueOf(out).Elem().Interface()
		return val, nil
	}

	return data, nil
}

// MapStruct maps the struct into the interface
func MapStruct(obj interface{}, input interface{}, tag ...string) error {
	config := &mapstructure.DecoderConfig{
		Result:           obj,
		WeaklyTypedInput: true,
		TagName:          "json",
		DecodeHook:       decoderHooks,
		Metadata:         &mapstructure.Metadata{},
	}

	if len(tag) > 0 {
		config.TagName = tag[0]
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}

	return nil
}

func StringMap(obj interface{}) map[string]string {
	rval := make(map[string]string)
	s := structs.New(obj)
	s.TagName = "json"

	for k, v := range s.Map() {
		switch t := v.(type) {
		case string:
			rval[k] = t
		case *string:
			if t != nil {
				rval[k] = *t
			}
		default:
			val := reflect.ValueOf(t)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			if val.CanInterface() {
				rval[k] = fmt.Sprint(val.Interface())
			}
		}
	}
	return rval
}

func MapStringToInterface(val map[string]string) map[string]interface{} {
	rval := make(map[string]interface{})

	for k, v := range val {
		rval[k] = v
	}

	return rval
}

func StructMap(obj interface{}, tag ...string) map[string]interface{} {
	if obj == nil {
		return map[string]interface{}{}
	}
	rval := make(map[string]interface{})
	s := structs.New(obj)
	s.TagName = "json"

	if len(tag) == 1 {
		s.TagName = tag[0]
	}

	for k, v := range s.Map() {
		switch t := v.(type) {
		case string:
			rval[k] = t
		case *string:
			if t != nil {
				rval[k] = *t
			}
		default:
			val := reflect.ValueOf(t)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			if val.CanInterface() {
				rval[k] = val.Interface()
			}
		}
	}
	return rval
}

// CopyStruct copies from one struct to another
func CopyStruct(src, dest interface{}, override bool) error {
	if override {
		mergo.Merge(dest, src, mergo.WithOverride)
	}
	return mergo.Merge(dest, src)
}

// SliceToMap converts a slice to a map
func SliceToMap(s interface{}, filter ...MapFilter) map[string]interface{} {
	src := reflect.ValueOf(s)
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	if src.Kind() != reflect.Slice {
		return make(map[string]interface{})
	}

	m := make(map[string]interface{})

	for i := 0; i < src.Len(); i++ {
		val := src.Index(i)
		v := val.Interface()

		var key string
		switch tv := v.(type) {
		case string:
			key = tv
		case *string:
			key = *tv
		case int64:
			key = fmt.Sprint(tv)
		case *int64:
			key = fmt.Sprint(*tv)
		case int32:
			key = fmt.Sprint(tv)
		case *int32:
			key = fmt.Sprint(*tv)
		default:
			if sv, ok := v.(fmt.Stringer); ok {
				key = sv.String()
			} else if sv, ok := v.(fmt.GoStringer); ok {
				key = sv.GoString()
			} else {
				key = val.Type().Name()
			}
		}
		if len(filter) == 0 {
			m[key] = v
		} else if mv, ok := filter[0](&key, v); ok {
			m[key] = mv
		}
	}
	return m
}

// NormalizeMap normalizes the map values to non-pointer types
func NormalizeMap(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr && !val.IsNil() {
			v = val.Elem().Interface()
		}
		m[k] = v
	}
	return m
}

// KVMap turns a separated list of values into a map, i.e. username=foo;password=bar
func KVMap(vals string, sep ...string) map[string]string {
	rval := make(map[string]string)
	_sep := ";"
	if len(sep) > 0 {
		_sep = sep[0]
	}

	parts := strings.Split(vals, _sep)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)

		if len(kv) == 2 {
			rval[kv[0]] = kv[1]
		}
	}

	return rval
}

// EnvMap creates a nested map from the os.Environ()
func EnvMap(prefix ...string) map[string]interface{} {
	rval := make(map[string]interface{})

	vals := os.Environ()
	replacer := strings.NewReplacer("_", ".")

	for _, v := range vals {
		parts := strings.SplitN(v, "=", 2)
		key := replacer.Replace(parts[0])

		if strings.HasPrefix(key, ".") {
			continue
		}

		key = strings.ToLower(key)

		rval[key] = parts[1]
	}

	if len(prefix) > 0 {
		return bellows.ExpandPrefixed(rval, prefix[0])
	}

	return bellows.Expand(rval)
}

func MergeStringMaps(dst, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func PrettySprint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}
	return "{}"
}

func MapPairs(cols []string, vals []interface{}, p ...string) map[string]interface{} {
	rval := make(map[string]interface{}, 0)
	prefix := ""
	if len(p) > 0 {
		prefix = p[0]
	}
	for i, k := range cols {
		k = strings.TrimPrefix(k, prefix+"_")
		rval[k] = vals[i]
	}

	return rval
}

func MapSeries(cols []string, values [][]interface{}, prefix ...string) []map[string]interface{} {
	rval := make([]map[string]interface{}, 0)
	hasPrefix := false

	if len(prefix) > 0 {
		hasPrefix = true
	}

	for _, vals := range values {
		tmp := make(map[string]interface{})
		for i, col := range cols {
			if hasPrefix {
				col = strings.TrimPrefix(col, prefix[0])
			}
			if vals[i] != nil {
				tmp[col] = vals[i]
			}
		}

		rval = append(rval, tmp)
	}

	return rval
}
