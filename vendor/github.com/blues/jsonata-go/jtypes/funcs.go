// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Package jtypes (golint)
package jtypes

import (
	"reflect"
)

// Resolve (golint)
func Resolve(v reflect.Value) reflect.Value {
	for {
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr:
			if !v.IsNil() {
				v = v.Elem()
				break
			}
			fallthrough
		default:
			return v
		}
	}
}

// IsBool (golint)
func IsBool(v reflect.Value) bool {
	return v.Kind() == reflect.Bool || resolvedKind(v) == reflect.Bool
}

// IsString (golint)
func IsString(v reflect.Value) bool {
	return v.Kind() == reflect.String || resolvedKind(v) == reflect.String
}

// IsNumber (golint)
func IsNumber(v reflect.Value) bool {
	return isFloat(v) || isInt(v) || isUint(v)
}

// IsCallable (golint)
func IsCallable(v reflect.Value) bool {
	v = Resolve(v)
	return v.IsValid() &&
		(v.Type().Implements(TypeCallable) || reflect.PtrTo(v.Type()).Implements(TypeCallable))
}

// IsArray (golint)
func IsArray(v reflect.Value) bool {
	return isArrayKind(v.Kind()) || isArrayKind(resolvedKind(v))
}

func isArrayKind(k reflect.Kind) bool {
	return k == reflect.Slice || k == reflect.Array
}

// IsArrayOf (golint)
func IsArrayOf(v reflect.Value, hasType func(reflect.Value) bool) bool {
	if !IsArray(v) {
		return false
	}

	v = Resolve(v)
	for i := 0; i < v.Len(); i++ {
		if !hasType(v.Index(i)) {
			return false
		}
	}

	return true
}

// IsMap (golint)
func IsMap(v reflect.Value) bool {
	return resolvedKind(v) == reflect.Map
}

// IsStruct (golint)
func IsStruct(v reflect.Value) bool {
	return resolvedKind(v) == reflect.Struct
}

// AsBool (golint)
func AsBool(v reflect.Value) (bool, bool) {
	v = Resolve(v)

	switch {
	case IsBool(v):
		return v.Bool(), true
	default:
		return false, false
	}
}

// AsString (golint)
func AsString(v reflect.Value) (string, bool) {
	v = Resolve(v)

	switch {
	case IsString(v):
		return v.String(), true
	default:
		return "", false
	}
}

// AsNumber (golint)
func AsNumber(v reflect.Value) (float64, bool) {
	v = Resolve(v)

	switch {
	case isFloat(v):
		return v.Float(), true
	case isInt(v), isUint(v):
		return v.Convert(typeFloat64).Float(), true
	default:
		return 0, false
	}
}

// AsCallable (golint)
func AsCallable(v reflect.Value) (Callable, bool) {
	v = Resolve(v)

	if v.IsValid() && v.Type().Implements(TypeCallable) && v.CanInterface() {
		return v.Interface().(Callable), true
	}

	if v.IsValid() && reflect.PtrTo(v.Type()).Implements(TypeCallable) && v.CanAddr() && v.Addr().CanInterface() {
		return v.Addr().Interface().(Callable), true
	}

	return nil, false
}

func isInt(v reflect.Value) bool {
	return isIntKind(v.Kind()) || isIntKind(resolvedKind(v))
}

func isIntKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func isUint(v reflect.Value) bool {
	return isUintKind(v.Kind()) || isUintKind(resolvedKind(v))
}

func isUintKind(k reflect.Kind) bool {
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func isFloat(v reflect.Value) bool {
	return isFloatKind(v.Kind()) || isFloatKind(resolvedKind(v))
}

func isFloatKind(k reflect.Kind) bool {
	switch k {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func resolvedKind(v reflect.Value) reflect.Kind {
	return Resolve(v).Kind()
}
