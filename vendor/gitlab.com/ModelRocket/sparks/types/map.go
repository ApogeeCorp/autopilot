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

import "reflect"

// Mapper provides a map wrapper interface for working easily with unknown map types
type Mapper interface {
	IsSet(key interface{}) bool
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Sub(key interface{}) Mapper
	Delete(key interface{})
	Copy() Mapper
	Without(keys ...interface{}) Mapper
}

type genericMap struct {
	v reflect.Value
}

// Map returns a Mapper from the passed map
func Map(m interface{}) Mapper {
	switch t := m.(type) {
	case map[string]interface{}:
		return StringMap(t)
	}
	return &genericMap{v: reflect.ValueOf(m)}
}

func (m genericMap) IsSet(key interface{}) bool {
	return m.v.MapIndex(reflect.ValueOf(key)).IsValid()
}

func (m genericMap) Set(key, value interface{}) {
	m.v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
}

func (m genericMap) Get(key interface{}) interface{} {
	v := m.v.MapIndex(reflect.ValueOf(key))
	if !v.CanInterface() {
		return nil
	}
	return v.Interface()
}

func (m genericMap) Sub(key interface{}) Mapper {
	v := m.v.MapIndex(reflect.ValueOf(key))
	if !v.CanInterface() {
		return nil
	}
	if v.Kind() != reflect.Map {
		return nil
	}
	return &genericMap{v}
}

func (m genericMap) Delete(key interface{}) {
	m.v.SetMapIndex(reflect.ValueOf(key), reflect.Value{})
}

func (m genericMap) Copy() Mapper {
	dst := reflect.New(reflect.TypeOf(m.v))

	for _, k := range m.v.MapKeys() {
		ov := m.v.MapIndex(k)
		dst.SetMapIndex(k, ov)
	}
	return &genericMap{dst}
}

func (m genericMap) Without(keys ...interface{}) Mapper {
	dst := m.Copy()

	for _, key := range keys {
		dst.Delete(key)
	}

	return dst
}
