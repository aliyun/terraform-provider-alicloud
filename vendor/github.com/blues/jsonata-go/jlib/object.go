// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jlib

import (
	"fmt"
	"reflect"

	"github.com/blues/jsonata-go/jtypes"
)

// typeInterfaceMap is the reflect.Type for map[string]interface{}.
// Go uses this type (by default) when decoding JSON objects.
// JSONata objects also use this type. As such it is a common
// map type in JSONata expressions.
var typeInterfaceMap = reflect.MapOf(typeString, jtypes.TypeInterface)

// toInterfaceMap attempts to cast a reflect.Value to a
// map[string]interface{}. This is useful for performance
// because direct map access is significantly faster than
// the reflect functions MapKeys and MapIndex.
func toInterfaceMap(v reflect.Value) (map[string]interface{}, bool) {
	if v.Type() == typeInterfaceMap && v.CanInterface() {
		return v.Interface().(map[string]interface{}), true
	}
	return nil, false
}

// Each applies the function fn to each name/value pair in
// the object obj and returns the results in an array. The
// order of the items in the array is undefined.
//
// obj must be a map or a struct. If it is a struct, any
// unexported fields are ignored.
//
// fn must be a Callable that takes one, two or three
// arguments. The first argument is the value of a name/value
// pair. The second and third arguments, if applicable, are
// the value and the source object respectively.
func Each(obj reflect.Value, fn jtypes.Callable) (interface{}, error) {

	var each func(reflect.Value, jtypes.Callable) ([]interface{}, error)

	obj = jtypes.Resolve(obj)

	switch {
	case jtypes.IsMap(obj):
		each = eachMap
	case jtypes.IsStruct(obj) && !jtypes.IsCallable(obj):
		each = eachStruct
	default:
		return nil, fmt.Errorf("argument must be an object")
	}

	if argc := fn.ParamCount(); argc < 1 || argc > 3 {
		return nil, fmt.Errorf("function must take 1, 2 or 3 arguments")
	}

	results, err := each(obj, fn)
	if err != nil {
		return nil, err
	}

	switch len(results) {
	case 0:
		return nil, jtypes.ErrUndefined
	case 1:
		return results[0], nil
	default:
		return results, nil
	}
}

func eachMap(v reflect.Value, fn jtypes.Callable) ([]interface{}, error) {

	size := v.Len()
	if size == 0 {
		return nil, nil
	}

	var results []interface{}

	argv := make([]reflect.Value, fn.ParamCount())

	for _, k := range v.MapKeys() {

		for i := range argv {
			switch i {
			case 0:
				argv[i] = v.MapIndex(k)
			case 1:
				argv[i] = k
			case 2:
				argv[i] = v
			}
		}

		res, err := fn.Call(argv)
		if err != nil {
			return nil, err
		}

		if res.IsValid() && res.CanInterface() {
			if results == nil {
				results = make([]interface{}, 0, size)
			}
			results = append(results, res.Interface())
		}
	}

	return results, nil
}

func eachStruct(v reflect.Value, fn jtypes.Callable) ([]interface{}, error) {

	size := v.NumField()
	if size == 0 {
		return nil, nil
	}

	var results []interface{}

	t := v.Type()
	argv := make([]reflect.Value, fn.ParamCount())

	for i := 0; i < size; i++ {

		field := t.Field(i)
		if field.PkgPath != "" {
			// Skip unexported fields.
			continue
		}

		for j := range argv {
			switch j {
			case 0:
				argv[j] = v.Field(i)
			case 1:
				argv[j] = reflect.ValueOf(field.Name)
			case 2:
				argv[j] = v
			}
		}

		res, err := fn.Call(argv)
		if err != nil {
			return nil, err
		}

		if res.IsValid() && res.CanInterface() {
			if results == nil {
				results = make([]interface{}, 0, size)
			}
			results = append(results, res.Interface())
		}
	}

	return results, nil
}

// Sift returns a map containing name/value pairs from the
// object obj that satisfy the predicate function fn.
//
// obj must be a map or a struct. If it is a map, the keys
// must be of type string. If it is a struct, any unexported
// fields are ignored.
//
// fn must be a Callable that takes one, two or three
// arguments. The first argument is the value of a name/value
// pair. The second and third arguments, if applicable, are
// the value and the source object respectively.
func Sift(obj reflect.Value, fn jtypes.Callable) (interface{}, error) {

	var sift func(reflect.Value, jtypes.Callable) (map[string]interface{}, error)

	obj = jtypes.Resolve(obj)

	switch {
	case jtypes.IsMap(obj):
		sift = siftMap
	case jtypes.IsStruct(obj) && !jtypes.IsCallable(obj):
		sift = siftStruct
	default:
		return nil, fmt.Errorf("argument must be an object")
	}

	if argc := fn.ParamCount(); argc < 1 || argc > 3 {
		return nil, fmt.Errorf("function must take 1, 2 or 3 arguments")
	}

	results, err := sift(obj, fn)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, jtypes.ErrUndefined
	}

	return results, nil
}

func siftMap(v reflect.Value, fn jtypes.Callable) (map[string]interface{}, error) {

	size := v.Len()
	if size == 0 {
		return nil, nil
	}

	var results map[string]interface{}

	argv := make([]reflect.Value, fn.ParamCount())

	for _, k := range v.MapKeys() {

		key, ok := jtypes.AsString(k)
		if !ok {
			return nil, fmt.Errorf("object key must evaluate to a string, got %v (%s)", k, k.Kind())
		}

		val := v.MapIndex(k)
		if !val.IsValid() || !val.CanInterface() {
			// Skip undefined or non-interfaceable values. We
			// already know we don't want them in the results,
			// so we can bypass the function call.
			continue
		}

		for i := range argv {
			switch i {
			case 0:
				argv[i] = val
			case 1:
				argv[i] = k
			case 2:
				argv[i] = v
			}
		}

		res, err := fn.Call(argv)
		if err != nil {
			return nil, err
		}

		if Boolean(res) {
			if results == nil {
				results = make(map[string]interface{}, size)
			}
			results[key] = val.Interface()
		}
	}

	return results, nil
}

func siftStruct(v reflect.Value, fn jtypes.Callable) (map[string]interface{}, error) {

	size := v.NumField()
	if size == 0 {
		return nil, nil
	}

	var results map[string]interface{}

	t := v.Type()
	argv := make([]reflect.Value, fn.ParamCount())

	for i := 0; i < size; i++ {

		key := t.Field(i).Name
		val := v.Field(i)
		if !val.IsValid() || !val.CanInterface() {
			// Skip undefined or non-interfaceable values. We
			// already know we don't want them in the results,
			// so we can bypass the function call. This also
			// filters out unexported fields (as they are
			// non-interfaceable).
			continue
		}

		for j := range argv {
			switch j {
			case 0:
				argv[j] = val
			case 1:
				argv[j] = reflect.ValueOf(key)
			case 2:
				argv[j] = v
			}
		}

		res, err := fn.Call(argv)
		if err != nil {
			return nil, err
		}

		if Boolean(res) {
			if results == nil {
				results = make(map[string]interface{}, size)
			}
			results[key] = val.Interface()
		}
	}

	return results, nil
}

// Keys returns an array of the names in the object obj.
// The order of the returned items is undefined.
//
// obj must be a map, a struct or an array. If obj is a map,
// its keys must be of type string. If obj is a struct, any
// unexported fields are ignored. And if obj is an array,
// Keys returns the unique set of names from each object
// in the array.
func Keys(obj reflect.Value) (interface{}, error) {

	results, err := keys(obj)
	if err != nil {
		return nil, err
	}

	switch len(results) {
	case 0:
		return nil, jtypes.ErrUndefined
	case 1:
		return results[0], nil
	default:
		return results, nil
	}
}

func keys(v reflect.Value) ([]string, error) {

	v = jtypes.Resolve(v)

	switch {
	case jtypes.IsMap(v):
		return keysMap(v)
	case jtypes.IsStruct(v) && !jtypes.IsCallable(v):
		return keysStruct(v)
	case jtypes.IsArray(v):
		return keysArray(v)
	default:
		return nil, nil
	}
}

func keysMap(v reflect.Value) ([]string, error) {

	if v.Len() == 0 {
		return nil, nil
	}

	if m, ok := toInterfaceMap(v); ok {
		return keysMapFast(m), nil
	}

	results := make([]string, v.Len())

	for i, k := range v.MapKeys() {

		key, ok := jtypes.AsString(k)
		if !ok {
			return nil, fmt.Errorf("object key must evaluate to a string, got %v (%s)", k, k.Kind())
		}

		results[i] = key
	}

	return results, nil
}

func keysMapFast(m map[string]interface{}) []string {

	results := make([]string, 0, len(m))
	for key := range m {
		results = append(results, key)
	}

	return results
}

func keysStruct(v reflect.Value) ([]string, error) {

	size := v.NumField()
	if size == 0 {
		return nil, nil
	}

	var results []string

	t := v.Type()
	for i := 0; i < size; i++ {

		field := t.Field(i)
		if field.PkgPath != "" {
			// Skip unexported fields.
			continue
		}

		if results == nil {
			results = make([]string, 0, size)
		}
		results = append(results, field.Name)
	}

	return results, nil
}

func keysArray(v reflect.Value) ([]string, error) {

	size := v.Len()
	if size == 0 {
		return nil, nil
	}

	kresults := make([][]string, 0, size)

	for i := 0; i < size; i++ {
		results, err := keys(v.Index(i))
		if err != nil {
			return nil, err
		}
		kresults = append(kresults, results)
	}

	size = 0
	for _, k := range kresults {
		size += len(k)
	}

	if size == 0 {
		return nil, nil
	}

	seen := map[string]bool{}
	results := make([]string, 0, size)

	for _, k := range kresults {
		for _, s := range k {
			if !seen[s] {
				seen[s] = true
				results = append(results, s)
			}
		}
	}

	return results, nil
}

// Merge merges an array of objects into a single object that
// contains all of the name/value pairs from the array objects.
// If a name appears multiple times, values from objects later
// in the array override those from earlier.
//
// objs must be an array of maps or structs. Maps must have
// keys of type string. Unexported struct fields are ignored.
func Merge(objs reflect.Value) (interface{}, error) {

	var size int
	var merge func(map[string]interface{}, reflect.Value) error

	objs = jtypes.Resolve(objs)

	switch {
	case jtypes.IsMap(objs):
		size = objs.Len()
		merge = mergeMap
	case jtypes.IsStruct(objs) && !jtypes.IsCallable(objs):
		size = objs.NumField()
		merge = mergeStruct
	case jtypes.IsArray(objs):
		for i := 0; i < objs.Len(); i++ {
			obj := jtypes.Resolve(objs.Index(i))
			switch {
			case jtypes.IsMap(obj):
				size += obj.Len()
			case jtypes.IsStruct(obj):
				size += obj.NumField()
			default:
				return nil, fmt.Errorf("argument must be an object or an array of objects")
			}
		}
		merge = mergeArray
	default:
		return nil, fmt.Errorf("argument must be an object or an array of objects")
	}

	results := make(map[string]interface{}, size)
	if err := merge(results, objs); err != nil {
		return nil, err
	}

	return results, nil
}

func mergeMap(dest map[string]interface{}, src reflect.Value) error {

	if m, ok := toInterfaceMap(src); ok {
		mergeMapFast(dest, m)
		return nil
	}

	for _, k := range src.MapKeys() {

		key, ok := jtypes.AsString(k)
		if !ok {
			return fmt.Errorf("object key must evaluate to a string, got %v (%s)", k, k.Kind())
		}

		if val := src.MapIndex(k); val.IsValid() && val.CanInterface() {
			dest[key] = val.Interface()
		}
	}

	return nil
}

func mergeMapFast(dest, src map[string]interface{}) {
	for k, v := range src {
		if v != nil {
			dest[k] = v
		}
	}
}

func mergeStruct(dest map[string]interface{}, src reflect.Value) error {

	t := src.Type()

	for i := 0; i < src.NumField(); i++ {

		field := t.Field(i)
		if field.PkgPath != "" {
			// Skip unexported fields.
			continue
		}

		if val := src.Field(i); val.IsValid() && val.CanInterface() {
			dest[field.Name] = val.Interface()
		}
	}

	return nil
}

func mergeArray(dest map[string]interface{}, src reflect.Value) error {

	var merge func(map[string]interface{}, reflect.Value) error

	for i := 0; i < src.Len(); i++ {

		item := jtypes.Resolve(src.Index(i))

		switch {
		case jtypes.IsMap(item):
			merge = mergeMap
		case jtypes.IsStruct(item) && !jtypes.IsCallable(item):
			merge = mergeStruct
		default:
			continue
		}

		if err := merge(dest, item); err != nil {
			return err
		}
	}

	return nil
}

// Spread (golint)
func Spread(v reflect.Value) (interface{}, error) {

	var results []interface{}

	switch {
	case jtypes.IsMap(v):
		v = jtypes.Resolve(v)
		keys := v.MapKeys()
		for _, k := range keys {
			if k.Kind() != reflect.String {
				return nil, fmt.Errorf("object key must evaluate to a string, got %v (%s)", k, k.Kind())
			}
			if v := v.MapIndex(k); v.CanInterface() {
				results = append(results, map[string]interface{}{
					k.String(): v.Interface(),
				})
			}
		}
	case jtypes.IsStruct(v) && !jtypes.IsCallable(v):
		v = jtypes.Resolve(v)
		for i := 0; i < v.NumField(); i++ {
			k := v.Type().Field(i).Name
			v := v.FieldByIndex([]int{i})
			if v.CanInterface() {
				results = append(results, map[string]interface{}{
					k: v.Interface(),
				})
			}
		}
	case jtypes.IsArray(v):
		v = jtypes.Resolve(v)
		for i := 0; i < v.Len(); i++ {
			res, err := Spread(v.Index(i))
			if err != nil {
				return nil, err
			}
			switch res := res.(type) {
			case []interface{}: // Check for []interface{} first because it will also match interface{}.
				results = append(results, res...)
			case interface{}:
				results = append(results, res)
			}
		}
	default:
		if v.IsValid() && v.CanInterface() {
			return v.Interface(), nil
		}
	}

	return results, nil
}
