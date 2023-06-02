// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jlib

import (
	"fmt"
	"reflect"

	"github.com/blues/jsonata-go/jtypes"
)

// Map (golint)
func Map(v reflect.Value, f jtypes.Callable) (interface{}, error) {

	v = forceArray(jtypes.Resolve(v))

	var results []interface{}

	argc := clamp(f.ParamCount(), 1, 3)

	for i := 0; i < arrayLen(v); i++ {

		argv := []reflect.Value{v.Index(i), reflect.ValueOf(i), v}

		res, err := f.Call(argv[:argc])
		if err != nil {
			return nil, err
		}
		if res.IsValid() && res.CanInterface() {
			results = append(results, res.Interface())
		}
	}

	return results, nil
}

// Filter (golint)
func Filter(v reflect.Value, f jtypes.Callable) (interface{}, error) {

	v = forceArray(jtypes.Resolve(v))

	var results []interface{}

	argc := clamp(f.ParamCount(), 1, 3)

	for i := 0; i < arrayLen(v); i++ {

		item := v.Index(i)
		argv := []reflect.Value{item, reflect.ValueOf(i), v}

		res, err := f.Call(argv[:argc])
		if err != nil {
			return nil, err
		}
		if Boolean(res) && item.IsValid() && item.CanInterface() {
			results = append(results, item.Interface())
		}
	}

	return results, nil
}

// Reduce (golint)
func Reduce(v reflect.Value, f jtypes.Callable, init jtypes.OptionalValue) (interface{}, error) {

	v = forceArray(jtypes.Resolve(v))

	var res reflect.Value

	if f.ParamCount() != 2 {
		return nil, fmt.Errorf("second argument of function \"reduce\" must be a function that takes two arguments")
	}

	i := 0
	switch {
	case init.IsSet():
		res = jtypes.Resolve(init.Value)
	case arrayLen(v) > 0:
		res = v.Index(0)
		i = 1
	}

	var err error
	for ; i < arrayLen(v); i++ {
		res, err = f.Call([]reflect.Value{res, v.Index(i)})
		if err != nil {
			return nil, err
		}
	}

	if !res.IsValid() || !res.CanInterface() {
		return nil, jtypes.ErrUndefined
	}

	return res.Interface(), nil
}

// Single returns the one and only one value in the array parameter that satisfy
// the function predicate (i.e. function returns Boolean true when passed the
// value). Returns an error if the number of matching values is not exactly
// one.
// https://docs.jsonata.org/higher-order-functions#single
func Single(v reflect.Value, f jtypes.Callable) (interface{}, error) {
	filteredValue, err := Filter(v, f)
	if err != nil {
		return nil, err
	}

	switch reflect.TypeOf(filteredValue).Kind() {
	case reflect.Slice:
		// Since Filter() returned a slice, if there is either zero or
		// more than one item in the slice, return a error, otherwise
		// return the item
		s := reflect.ValueOf(filteredValue)
		if s.Len() != 1 {
			return nil, fmt.Errorf("number of matching values returned by single() must be 1, got: %d", s.Len())
		}
		return s.Index(0).Interface(), nil

	default:
		// Filter returned a single value, so use that
		return reflect.ValueOf(filteredValue).Interface(), nil
	}
}

func clamp(n, min, max int) int {
	switch {
	case n < min:
		return min
	case n > max:
		return max
	default:
		return n
	}
}
