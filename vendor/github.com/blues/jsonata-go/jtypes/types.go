// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Package jtypes provides types and utilities for third party
// extension functions.
package jtypes

import (
	"errors"
	"reflect"
)

var undefined reflect.Value

var (
	typeBool    = reflect.TypeOf((*bool)(nil)).Elem()
	typeInt     = reflect.TypeOf((*int)(nil)).Elem()
	typeFloat64 = reflect.TypeOf((*float64)(nil)).Elem()
	typeString  = reflect.TypeOf((*string)(nil)).Elem()

	// TypeOptional (golint)
	TypeOptional = reflect.TypeOf((*Optional)(nil)).Elem()
	// TypeCallable (golint)
	TypeCallable = reflect.TypeOf((*Callable)(nil)).Elem()
	// TypeConvertible (golint)
	TypeConvertible = reflect.TypeOf((*Convertible)(nil)).Elem()
	// TypeVariant (golint)
	TypeVariant = reflect.TypeOf((*Variant)(nil)).Elem()
	// TypeValue (golint)
	TypeValue = reflect.TypeOf((*reflect.Value)(nil)).Elem()
	// TypeInterface (golint)
	TypeInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

// ErrUndefined (golint)
var ErrUndefined = errors.New("undefined")

// Variant (golint)
type Variant interface {
	ValidTypes() []reflect.Type
}

// Callable (golint)
type Callable interface {
	Name() string
	ParamCount() int
	Call([]reflect.Value) (reflect.Value, error)
}

// Convertible (golint)
type Convertible interface {
	ConvertTo(reflect.Type) (reflect.Value, bool)
}

// Optional (golint)
type Optional interface {
	IsSet() bool
	Set(reflect.Value)
	Type() reflect.Type
}

type isSet bool

// IsSet (golint)
func (opt *isSet) IsSet() bool {
	return bool(*opt)
}

// OptionalBool (golint)
type OptionalBool struct {
	isSet
	Bool bool
}

// NewOptionalBool (golint)
func NewOptionalBool(value bool) OptionalBool {
	opt := OptionalBool{}
	opt.Set(reflect.ValueOf(value))
	return opt
}

// Set (golint)
func (opt *OptionalBool) Set(v reflect.Value) {
	opt.isSet = true
	opt.Bool = v.Bool()
}

// Type (golint)
func (opt *OptionalBool) Type() reflect.Type {
	return typeBool
}

// OptionalInt (golint)
type OptionalInt struct {
	isSet
	Int int
}

// NewOptionalInt (golint)
func NewOptionalInt(value int) OptionalInt {
	opt := OptionalInt{}
	opt.Set(reflect.ValueOf(value))
	return opt
}

// Set (golint)
func (opt *OptionalInt) Set(v reflect.Value) {
	opt.isSet = true
	opt.Int = int(v.Int())
}

// Type (golint)
func (opt *OptionalInt) Type() reflect.Type {
	return typeInt
}

// OptionalFloat64 (golint)
type OptionalFloat64 struct {
	isSet
	Float64 float64
}

// NewOptionalFloat64 (golint)
func NewOptionalFloat64(value float64) OptionalFloat64 {
	opt := OptionalFloat64{}
	opt.Set(reflect.ValueOf(value))
	return opt
}

// Set (golint)
func (opt *OptionalFloat64) Set(v reflect.Value) {
	opt.isSet = true
	opt.Float64 = v.Float()
}

// Type (golint)
func (opt *OptionalFloat64) Type() reflect.Type {
	return typeFloat64
}

// OptionalString (golint)
type OptionalString struct {
	isSet
	String string
}

// NewOptionalString (golint)
func NewOptionalString(value string) OptionalString {
	opt := OptionalString{}
	opt.Set(reflect.ValueOf(value))
	return opt
}

// Set (golint)
func (opt *OptionalString) Set(v reflect.Value) {
	opt.isSet = true
	opt.String = v.String()
}

// Type (golint)
func (opt *OptionalString) Type() reflect.Type {
	return typeString
}

// OptionalInterface (golint)
type OptionalInterface struct {
	isSet
	Interface interface{}
}

// NewOptionalInterface (golint)
func NewOptionalInterface(value interface{}) OptionalInterface {
	opt := OptionalInterface{}
	opt.Set(reflect.ValueOf(value))
	return opt
}

// Set (golint)
func (opt *OptionalInterface) Set(v reflect.Value) {
	opt.isSet = true
	opt.Interface = v.Interface()
}

// Type (golint)
func (opt *OptionalInterface) Type() reflect.Type {
	return TypeInterface
}

// OptionalValue (golint)
type OptionalValue struct {
	isSet
	Value reflect.Value
}

// NewOptionalValue (golint)
func NewOptionalValue(value reflect.Value) OptionalValue {
	opt := OptionalValue{}
	opt.Set(reflect.ValueOf(value))
	return opt
}

// Set (golint)
func (opt *OptionalValue) Set(v reflect.Value) {
	opt.isSet = true
	opt.Value = v.Interface().(reflect.Value)
}

// Type (golint)
func (opt *OptionalValue) Type() reflect.Type {
	return TypeValue
}

// OptionalCallable (golint)
type OptionalCallable struct {
	isSet
	Callable Callable
}

// NewOptionalCallable (golint)
func NewOptionalCallable(value Callable) OptionalCallable {
	opt := OptionalCallable{}
	opt.Set(reflect.ValueOf(value))
	return opt
}

// Set (golint)
func (opt *OptionalCallable) Set(v reflect.Value) {
	opt.isSet = true
	opt.Callable = v.Interface().(Callable)
}

// Type (golint)
func (opt *OptionalCallable) Type() reflect.Type {
	return TypeCallable
}

// ArgHandler (golint)
type ArgHandler func([]reflect.Value) bool

// ArgCountEquals (golint)
func ArgCountEquals(n int) ArgHandler {
	return func(argv []reflect.Value) bool {
		return len(argv) == n
	}
}

// ArgUndefined (golint)
func ArgUndefined(i int) ArgHandler {
	return func(argv []reflect.Value) bool {
		return len(argv) > i && argv[i] == undefined
	}
}
