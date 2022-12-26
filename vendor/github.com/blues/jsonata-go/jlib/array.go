// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jlib

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"

	"github.com/blues/jsonata-go/jtypes"
)

// Count (golint)
func Count(v reflect.Value) int {
	v = jtypes.Resolve(v)

	if !jtypes.IsArray(v) {
		if v.IsValid() {
			return 1
		}
		return 0
	}

	return v.Len()
}

// Distinct returns the values passed in with any duplicates removed.
func Distinct(v reflect.Value) interface{} {
	v = jtypes.Resolve(v)

	// To match the behavior of jsonata-js, if this is a string we should
	// return the entire string and not dedupe the individual characters
	if jtypes.IsString(v) {
		return v.String()
	}

	if jtypes.IsArray(v) {
		items := arrayify(v)
		visited := make(map[interface{}]struct{})
		distinctValues := reflect.MakeSlice(reflect.SliceOf(typeInterface), 0, 0)

		for i := 0; i < items.Len(); i++ {
			item := jtypes.Resolve(items.Index(i))

			if jtypes.IsMap(item) {
				// We can't hash a map, so convert it to a
				// string that is hashable
				mapItem := fmt.Sprint(item.Interface())
				if _, ok := visited[mapItem]; ok {
					continue
				}
				visited[mapItem] = struct{}{}
				distinctValues = reflect.Append(distinctValues, item)

				continue
			}

			if _, ok := visited[item.Interface()]; ok {
				continue
			}

			visited[item.Interface()] = struct{}{}
			distinctValues = reflect.Append(distinctValues, item)
		}
		return distinctValues.Interface()
	}

	return nil
}

// Append (golint)
func Append(v1, v2 reflect.Value) (interface{}, error) {
	if !v2.IsValid() && v1.IsValid() && v1.CanInterface() {
		return v1.Interface(), nil
	}

	if !v1.IsValid() && v2.IsValid() && v2.CanInterface() {
		return v2.Interface(), nil
	}

	v1 = arrayify(v1)
	v2 = arrayify(v2)

	len1 := v1.Len()
	len2 := v2.Len()

	results := reflect.MakeSlice(reflect.SliceOf(typeInterface), 0, len1+len2)

	appendSlice := func(vs reflect.Value, length int) {
		for i := 0; i < length; i++ {
			if item := vs.Index(i); item.IsValid() {
				results = reflect.Append(results, item)
			}
		}
	}

	appendSlice(v1, len1)
	appendSlice(v2, len2)

	return results.Interface(), nil
}

// Reverse (golint)
func Reverse(v reflect.Value) (interface{}, error) {
	v = arrayify(v)
	length := v.Len()

	results := reflect.MakeSlice(v.Type(), 0, length)

	for i := length - 1; i >= 0; i-- {
		if item := v.Index(i); item.IsValid() {
			results = reflect.Append(results, item)
		}
	}

	return results.Interface(), nil
}

// Sort (golint)
func Sort(v reflect.Value, swap jtypes.OptionalCallable) (interface{}, error) {
	v = jtypes.Resolve(v)

	switch {
	case !v.IsValid():
		return nil, jtypes.ErrUndefined
	case !jtypes.IsArray(v):
		if v.CanInterface() {
			return []interface{}{v.Interface()}, nil
		}
	case swap.Callable != nil:
		return sortArrayFunc(v, swap.Callable)
	case jtypes.IsArrayOf(v, jtypes.IsNumber):
		return sortNumberArray(v), nil
	case jtypes.IsArrayOf(v, jtypes.IsString):
		return sortStringArray(v), nil
	}

	return nil, fmt.Errorf("argument 1 of function sort must be an array of strings or numbers")
}

func sortNumberArray(v reflect.Value) []interface{} {
	size := v.Len()
	results := make([]interface{}, 0, size)

	for i := 0; i < size; i++ {
		if n, ok := jtypes.AsNumber(v.Index(i)); ok {
			results = append(results, n)
		}
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].(float64) < results[j].(float64)
	})

	return results
}

func sortStringArray(v reflect.Value) []interface{} {
	size := v.Len()
	results := make([]interface{}, 0, size)

	for i := 0; i < size; i++ {
		if s, ok := jtypes.AsString(v.Index(i)); ok {
			results = append(results, s)
		}
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].(string) < results[j].(string)
	})

	return results
}

func sortArrayFunc(v reflect.Value, fn jtypes.Callable) (interface{}, error) {
	size := v.Len()
	results := make([]interface{}, 0, size)

	for i := 0; i < size; i++ {
		if item := v.Index(i); item.CanInterface() {
			results = append(results, item.Interface())
		}
	}

	swapFunc := func(lhs, rhs interface{}) (bool, error) {

		args := []reflect.Value{
			reflect.ValueOf(lhs),
			reflect.ValueOf(rhs),
		}

		v, err := fn.Call(args)
		if err != nil {
			return false, err
		}

		b, ok := jtypes.AsBool(v)
		if !ok {
			return false, fmt.Errorf("argument 2 of function sort must be a function that returns a boolean, got %v (%s)", v, v.Kind())
		}

		return b, nil
	}

	return mergeSort(results, swapFunc)
}

func mergeSort(values []interface{}, swapFunc func(interface{}, interface{}) (bool, error)) ([]interface{}, error) {
	n := len(values)
	if n < 2 {
		return values, nil
	}

	pos := n / 2
	lhs, err := mergeSort(values[:pos], swapFunc)
	if err != nil {
		return nil, err
	}
	rhs, err := mergeSort(values[pos:], swapFunc)
	if err != nil {
		return nil, err
	}

	return merge(lhs, rhs, swapFunc)
}

func merge(lhs, rhs []interface{}, swapFunc func(interface{}, interface{}) (bool, error)) ([]interface{}, error) {
	results := make([]interface{}, len(lhs)+len(rhs))

	for i := range results {

		if len(rhs) == 0 {
			results = append(results[:i], lhs...)
			break
		}

		if len(lhs) == 0 {
			results = append(results[:i], rhs...)
			break
		}

		swap, err := swapFunc(lhs[0], rhs[0])
		if err != nil {
			return nil, err
		}

		if swap {
			results[i] = rhs[0]
			rhs = rhs[1:]
		} else {
			results[i] = lhs[0]
			lhs = lhs[1:]
		}
	}

	return results, nil
}

// Shuffle (golint)
func Shuffle(v reflect.Value) interface{} {
	v = forceArray(jtypes.Resolve(v))

	length := arrayLen(v)
	results := make([]interface{}, length)

	for i := 0; i < length; i++ {

		j := rand.Intn(i + 1)

		if i != j {
			results[i] = results[j]
		}

		item := v.Index(i)
		if item.IsValid() && item.CanInterface() {
			results[j] = item.Interface()
		}
	}

	return results
}

// Zip (golint)
func Zip(vs ...reflect.Value) (interface{}, error) {
	var size int

	if len(vs) == 0 {
		return nil, fmt.Errorf("cannot call zip with no arguments")
	}

	for i := 0; i < len(vs); i++ {

		vs[i] = forceArray(jtypes.Resolve(vs[i]))
		if !vs[i].IsValid() {
			return []interface{}{}, nil
		}

		if i == 0 || arrayLen(vs[i]) < size {
			size = arrayLen(vs[i])
		}
	}

	result := make([]interface{}, size)

	for i := 0; i < size; i++ {

		inner := make([]interface{}, len(vs))

		for j := 0; j < len(vs); j++ {
			v := vs[j].Index(i)
			if v.IsValid() && v.CanInterface() {
				inner[j] = v.Interface()
			}
		}

		result[i] = inner
	}

	return result, nil
}

func forceArray(v reflect.Value) reflect.Value {
	v = jtypes.Resolve(v)
	if !v.IsValid() || jtypes.IsArray(v) {
		return v
	}
	vs := reflect.MakeSlice(reflect.SliceOf(v.Type()), 0, 1)
	vs = reflect.Append(vs, v)
	return vs
}

func arrayLen(v reflect.Value) int {
	if jtypes.IsArray(v) {
		return v.Len()
	}
	return 0
}

var typeInterface = reflect.TypeOf((*interface{})(nil)).Elem()

func arrayify(v reflect.Value) reflect.Value {
	switch {
	case jtypes.IsArray(v):
		return jtypes.Resolve(v)
	case !v.IsValid():
		return reflect.MakeSlice(reflect.SliceOf(typeInterface), 0, 0)
	default:
		return reflect.Append(reflect.MakeSlice(reflect.SliceOf(typeInterface), 0, 1), v)
	}
}
