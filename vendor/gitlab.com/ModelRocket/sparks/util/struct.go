// Copyright © 2016 Charles Phillips <charles@doublerebel.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"reflect"
	"strings"
)

func Expand(value map[string]interface{}) map[string]interface{} {
	return ExpandPrefixed(value, "")
}

func ExpandPrefixed(value map[string]interface{}, prefix string) map[string]interface{} {
	m := make(map[string]interface{})
	ExpandPrefixedToResult(value, prefix, m)
	return m
}

func ExpandPrefixedToResult(value map[string]interface{}, prefix string, result map[string]interface{}) {
	if prefix != "" {
		prefix += "."
	}
	for k, val := range value {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		key := k[len(prefix):]
		idx := strings.Index(key, ".")
		if idx != -1 {
			key = key[:idx]
		}
		if _, ok := result[key]; ok {
			continue
		}
		if idx == -1 {
			result[key] = val
			continue
		}

		// It contains a period, so it is a more complex structure
		result[key] = ExpandPrefixed(value, k[:len(prefix)+len(key)])
	}
}

func Flatten(value interface{}, tags ...string) map[string]interface{} {
	return FlattenPrefixed(value, "", tags...)
}

func FlattenPrefixed(value interface{}, prefix string, tags ...string) map[string]interface{} {
	if len(tags) == 0 {
		tags = []string{"json"}
	}

	m := make(map[string]interface{})
	FlattenPrefixedToResult(value, prefix, m, tags)
	return m
}

func FlattenPrefixedToResult(value interface{}, prefix string, m map[string]interface{}, tags []string) {
	base := ""
	if prefix != "" {
		base = prefix + "."
	}

	original := reflect.ValueOf(value)
	kind := original.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		original = reflect.Indirect(original)
		kind = original.Kind()
	}

	if !original.IsValid() {
		if prefix != "" {
			m[prefix] = nil
		}
		return
	}

	t := original.Type()

	switch kind {
	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			break
		}
		for _, childKey := range original.MapKeys() {
			childValue := original.MapIndex(childKey)
			key := strings.ToLower(childKey.String())
			FlattenPrefixedToResult(childValue.Interface(), base+key, m, tags)
		}
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			childValue := original.Field(i)

			childKey := t.Field(i).Name

			field, _ := t.FieldByName(childKey)

			omit := false
			override := ""

			for _, tag := range tags {
				tagName, tagOpts := parseTag(field.Tag.Get(tag))
				if tagName == "-" {
					omit = true
					break
				}
				if tagName != "" {
					childKey = tagName
				} else {
					continue
				}

				// if the value is a zero value and the field is marked as omitempty do
				// not include
				if _, ok := tagOpts.Get("omitempty"); ok {
					zero := reflect.Zero(childValue.Type()).Interface()
					current := childValue.Interface()

					if reflect.DeepEqual(current, zero) {
						omit = true
					}
				}
				if b, ok := tagOpts.Get("base"); ok {
					override = b + "."
				}
				break
			}

			if omit {
				continue
			}

			if childValue.Kind() == reflect.Ptr && !childValue.IsNil() {
				childValue = childValue.Elem()
			}

			next := base + childKey
			if override != "" {
				next = override + childKey
			}
			FlattenPrefixedToResult(childValue.Interface(), next, m, tags)
		}
	default:
		if prefix != "" {
			prefix = strings.ToLower(prefix)
			m[prefix] = value
		}
	}
}

// tagOptions contains a slice of tag options
type tagOptions map[string]string

// Has returns true if the given option is available in tagOptions
func (t tagOptions) Get(opt string) (string, bool) {
	for tagOpt, v := range t {
		if tagOpt == opt {
			return v, true
		}
	}

	return "", false
}

// parseTag splits a struct field's tag into its name and a list of options
// which comes after a name. A tag is in the form of: "name,option1,option2".
// The name can be neglectected.
func parseTag(tag string) (string, tagOptions) {
	// tag is one of followings:
	// ""
	// "name"
	// "name,opt"
	// "name,opt,opt2"
	// ",opt"

	res := strings.Split(tag, ",")

	opts := make(tagOptions)
	for _, opt := range res[1:] {
		parts := strings.Split(opt, "=")
		if len(parts) == 1 {
			opts[parts[0]] = "true"
		} else if len(parts) == 2 {
			opts[parts[0]] = parts[1]
		}
	}

	return res[0], opts
}
