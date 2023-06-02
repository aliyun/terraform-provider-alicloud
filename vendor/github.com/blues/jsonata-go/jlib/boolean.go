// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jlib

import (
	"reflect"

	"github.com/blues/jsonata-go/jtypes"
)

// Boolean (golint)
func Boolean(v reflect.Value) bool {

	v = jtypes.Resolve(v)

	if b, ok := jtypes.AsBool(v); ok {
		return b
	}

	if s, ok := jtypes.AsString(v); ok {
		return s != ""
	}

	if n, ok := jtypes.AsNumber(v); ok {
		return n != 0
	}

	if jtypes.IsArray(v) {
		for i := 0; i < v.Len(); i++ {
			if Boolean(v.Index(i)) {
				return true
			}
		}
		return false
	}

	if jtypes.IsMap(v) {
		return v.Len() > 0
	}

	return false
}

// Not (golint)
func Not(v reflect.Value) bool {
	return !Boolean(v)
}

// Exists (golint)
func Exists(v reflect.Value) bool {
	return v.IsValid()
}
