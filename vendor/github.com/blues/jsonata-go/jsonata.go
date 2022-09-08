// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jsonata

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"
	"unicode"

	"github.com/blues/jsonata-go/jlib"
	"github.com/blues/jsonata-go/jparse"
	"github.com/blues/jsonata-go/jtypes"
)

var (
	globalRegistryMutex sync.RWMutex
	globalRegistry      map[string]reflect.Value
)

// An Extension describes custom functionality added to a
// JSONata expression.
type Extension struct {

	// Func is a Go function that implements the custom
	// functionality and returns either one or two values.
	// The second return value, if provided, must be an
	// error.
	Func interface{}

	// UndefinedHandler is a function that determines how
	// this extension handles undefined arguments. If
	// UndefinedHandler is non-nil, it is called before
	// Func with the same arguments. If the handler returns
	// true, Func is not called and undefined is returned
	// instead.
	UndefinedHandler jtypes.ArgHandler

	// EvalContextHandler is a function that determines how
	// this extension handles missing arguments. If
	// EvalContextHandler is non-nil, it is called before
	// Func with the same arguments. If the handler returns
	// true, the evaluation context is inserted as the first
	// argument when Func is called.
	EvalContextHandler jtypes.ArgHandler
}

// RegisterExts registers custom functions for use in JSONata
// expressions. It is designed to be called once on program
// startup (e.g. from an init function).
//
// Custom functions registered at the package level will be
// available to all Expr objects. To register custom functions
// with specific Expr objects, use the RegisterExts method.
func RegisterExts(exts map[string]Extension) error {

	values, err := processExts(exts)
	if err != nil {
		return err
	}

	updateGlobalRegistry(values)
	return nil
}

// RegisterVars registers custom variables for use in JSONata
// expressions. It is designed to be called once on program
// startup (e.g. from an init function).
//
// Custom variables registered at the package level will be
// available to all Expr objects. To register custom variables
// with specific Expr objects, use the RegisterVars method.
func RegisterVars(vars map[string]interface{}) error {

	values, err := processVars(vars)
	if err != nil {
		return err
	}

	updateGlobalRegistry(values)
	return nil
}

// An Expr represents a JSONata expression.
type Expr struct {
	node     jparse.Node
	registry map[string]reflect.Value
}

// Compile parses a JSONata expression and returns an Expr
// that can be evaluated against JSON data. If the input is
// not a valid JSONata expression, Compile returns an error
// of type jparse.Error.
func Compile(expr string) (*Expr, error) {

	node, err := jparse.Parse(expr)
	if err != nil {
		return nil, err
	}

	e := &Expr{
		node: node,
	}

	globalRegistryMutex.RLock()
	e.updateRegistry(globalRegistry)
	globalRegistryMutex.RUnlock()

	return e, nil
}

// MustCompile is like Compile except it panics if given an
// invalid expression.
func MustCompile(expr string) *Expr {

	e, err := Compile(expr)
	if err != nil {
		panicf("could not compile %s: %s", expr, err)
	}

	return e
}

// Eval executes a JSONata expression against the given data
// source. The input is typically the result of unmarshaling
// a JSON string. The output is an object suitable for
// marshaling into a JSON string. Use EvalBytes to skip the
// unmarshal/marshal steps and work solely with JSON strings.
//
// Eval can be called multiple times, with different input
// data if required.
func (e *Expr) Eval(data interface{}) (interface{}, error) {
	input, ok := data.(reflect.Value)
	if !ok {
		input = reflect.ValueOf(data)
	}

	result, err := eval(e.node, input, e.newEnv(input))
	if err != nil {
		return nil, err
	}

	if !result.IsValid() {
		return nil, ErrUndefined
	}

	if !result.CanInterface() {
		return nil, fmt.Errorf("Eval returned a non-interface value")
	}

	if result.Kind() == reflect.Ptr && result.IsNil() {
		return nil, nil
	}

	return result.Interface(), nil
}

// EvalBytes is like Eval but it accepts and returns byte slices
// instead of objects.
func (e *Expr) EvalBytes(data []byte) ([]byte, error) {

	var v interface{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	v, err = e.Eval(v)
	if err != nil {
		return nil, err
	}

	return json.Marshal(v)
}

// RegisterExts registers custom functions for use during
// evaluation. Custom functions registered with this method
// are only available to this Expr object. To make custom
// functions available to all Expr objects, use the package
// level RegisterExts function.
func (e *Expr) RegisterExts(exts map[string]Extension) error {

	values, err := processExts(exts)
	if err != nil {
		return err
	}

	e.updateRegistry(values)
	return nil
}

// RegisterVars registers custom variables for use during
// evaluation. Custom variables registered with this method
// are only available to this Expr object. To make custom
// variables available to all Expr objects, use the package
// level RegisterVars function.
func (e *Expr) RegisterVars(vars map[string]interface{}) error {

	values, err := processVars(vars)
	if err != nil {
		return err
	}

	e.updateRegistry(values)
	return nil
}

// String returns a string representation of an Expr.
func (e *Expr) String() string {
	if e.node == nil {
		return ""
	}
	return e.node.String()
}

func (e *Expr) updateRegistry(values map[string]reflect.Value) {

	for name, v := range values {
		if e.registry == nil {
			e.registry = make(map[string]reflect.Value, len(values))
		}
		e.registry[name] = v
	}
}

func (e *Expr) newEnv(input reflect.Value) *environment {

	tc := timeCallables(time.Now())

	env := newEnvironment(baseEnv, len(tc)+len(e.registry)+1)

	env.bind("$", input)
	env.bindAll(tc)
	env.bindAll(e.registry)

	return env
}

var (
	milisT = mustGoCallable("millis", Extension{
		Func: func(millis int64) int64 {
			return millis
		},
	})

	nowT = mustGoCallable("now", Extension{
		Func: func(millis int64, picture jtypes.OptionalString, tz jtypes.OptionalString) (string, error) {
			return jlib.FromMillis(millis, picture, tz)
		},
	})
)

func timeCallables(t time.Time) map[string]reflect.Value {

	ms := t.UnixNano() / int64(time.Millisecond)

	millis := &partialCallable{
		callableName: callableName{
			name: "millis",
		},
		fn: milisT,
		args: []jparse.Node{
			&jparse.NumberNode{
				Value: float64(ms),
			},
		},
	}

	now := &partialCallable{
		callableName: callableName{
			name: "now",
		},
		fn: nowT,
		args: []jparse.Node{
			&jparse.NumberNode{
				Value: float64(ms),
			},
			&jparse.PlaceholderNode{},
			&jparse.PlaceholderNode{},
		},
	}

	return map[string]reflect.Value{
		"millis": reflect.ValueOf(millis),
		"now":    reflect.ValueOf(now),
	}
}

func processExts(exts map[string]Extension) (map[string]reflect.Value, error) {

	var m map[string]reflect.Value

	for name, ext := range exts {

		if !validName(name) {
			return nil, fmt.Errorf("%s is not a valid name", name)
		}

		callable, err := newGoCallable(name, ext)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid function: %s", name, err)
		}

		if m == nil {
			m = make(map[string]reflect.Value, len(exts))
		}
		m[name] = reflect.ValueOf(callable)
	}

	return m, nil
}

func processVars(vars map[string]interface{}) (map[string]reflect.Value, error) {

	var m map[string]reflect.Value

	for name, value := range vars {

		if !validName(name) {
			return nil, fmt.Errorf("%s is not a valid name", name)
		}

		if !validVar(value) {
			return nil, fmt.Errorf("%s is not a valid variable", name)
		}

		if m == nil {
			m = make(map[string]reflect.Value, len(vars))
		}
		m[name] = reflect.ValueOf(value)
	}

	return m, nil
}

func updateGlobalRegistry(values map[string]reflect.Value) {

	globalRegistryMutex.Lock()

	for name, v := range values {
		if globalRegistry == nil {
			globalRegistry = make(map[string]reflect.Value, len(values))
		}
		globalRegistry[name] = v
	}

	globalRegistryMutex.Unlock()
}

func validName(s string) bool {

	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !isLetter(r) && !isDigit(r) && r != '_' {
			return false
		}
	}

	return true
}

func validVar(v interface{}) bool {
	// TODO: Variable validation.
	return true
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return (r >= '0' && r <= '9') || unicode.IsDigit(r)
}
