// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Package jlib implements the JSONata function library.
package jlib

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/blues/jsonata-go/jtypes"
)

func init() {
	// Seed random numbers for Random() and Shuffle().
	rand.Seed(time.Now().UnixNano())
}

var typeBool = reflect.TypeOf((*bool)(nil)).Elem()
var typeCallable = reflect.TypeOf((*jtypes.Callable)(nil)).Elem()
var typeString = reflect.TypeOf((*string)(nil)).Elem()
var typeNumber = reflect.TypeOf((*float64)(nil)).Elem()

// StringNumberBool (golint)
type StringNumberBool reflect.Value

// ValidTypes (golint)
func (StringNumberBool) ValidTypes() []reflect.Type {
	return []reflect.Type{
		typeBool,
		typeString,
		typeNumber,
	}
}

// StringCallable (golint)
type StringCallable reflect.Value

// ValidTypes (golint)
func (StringCallable) ValidTypes() []reflect.Type {
	return []reflect.Type{
		typeString,
		typeCallable,
	}
}

func (s StringCallable) toInterface() interface{} {
	if v := reflect.Value(s); v.IsValid() && v.CanInterface() {
		return v.Interface()
	}
	return nil
}

// TypeOf implements the jsonata $type function that returns the data type of
// the argument
func TypeOf(x interface{}) (string, error) {
	v := reflect.ValueOf(x)
	if jtypes.IsCallable(v) {
		return "function", nil
	}
	if jtypes.IsString(v) {
		return "string", nil
	}
	if jtypes.IsNumber(v) {
		return "number", nil
	}
	if jtypes.IsArray(v) {
		return "array", nil
	}
	if jtypes.IsBool(v) {
		return "boolean", nil
	}
	if jtypes.IsMap(v) {
		return "object", nil
	}

	switch x.(type) {
	case *interface{}:
		return "null", nil
	}

	xType := reflect.TypeOf(x).String()
	return "", fmt.Errorf("unknown type %s", xType)
}
