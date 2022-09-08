// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jlib

import (
	"fmt"
	"reflect"

	"github.com/blues/jsonata-go/jtypes"
)

// Sum returns the total of an array of numbers. If the array is
// empty, Sum returns 0.
func Sum(v reflect.Value) (float64, error) {

	if !jtypes.IsArray(v) {
		if n, ok := jtypes.AsNumber(v); ok {
			return n, nil
		}
		return 0, fmt.Errorf("cannot call sum on a non-array type")
	}

	v = jtypes.Resolve(v)

	var sum float64

	for i := 0; i < v.Len(); i++ {
		n, ok := jtypes.AsNumber(v.Index(i))
		if !ok {
			return 0, fmt.Errorf("cannot call sum on an array with non-number types")
		}
		sum += n
	}

	return sum, nil
}

// Max returns the largest value in an array of numbers. If the
// array is empty, Max returns 0 and an undefined error.
func Max(v reflect.Value) (float64, error) {

	if !jtypes.IsArray(v) {
		if n, ok := jtypes.AsNumber(v); ok {
			return n, nil
		}
		return 0, fmt.Errorf("cannot call max on a non-array type")
	}

	v = jtypes.Resolve(v)
	if v.Len() == 0 {
		return 0, jtypes.ErrUndefined
	}

	var max float64

	for i := 0; i < v.Len(); i++ {
		n, ok := jtypes.AsNumber(v.Index(i))
		if !ok {
			return 0, fmt.Errorf("cannot call max on an array with non-number types")
		}
		if i == 0 || n > max {
			max = n
		}
	}

	return max, nil
}

// Min returns the smallest value in an array of numbers. If the
// array is empty, Min returns 0 and an undefined error.
func Min(v reflect.Value) (float64, error) {

	if !jtypes.IsArray(v) {
		if n, ok := jtypes.AsNumber(v); ok {
			return n, nil
		}
		return 0, fmt.Errorf("cannot call min on a non-array type")
	}

	v = jtypes.Resolve(v)
	if v.Len() == 0 {
		return 0, jtypes.ErrUndefined
	}

	var min float64

	for i := 0; i < v.Len(); i++ {
		n, ok := jtypes.AsNumber(v.Index(i))
		if !ok {
			return 0, fmt.Errorf("cannot call min on an array with non-number types")
		}
		if i == 0 || n < min {
			min = n
		}
	}

	return min, nil
}

// Average returns the mean of an array of numbers. If the array
// is empty, Average returns 0 and an undefined error.
func Average(v reflect.Value) (float64, error) {

	if !jtypes.IsArray(v) {
		if n, ok := jtypes.AsNumber(v); ok {
			return n, nil
		}
		return 0, fmt.Errorf("cannot call average on a non-array type")
	}

	v = jtypes.Resolve(v)
	if v.Len() == 0 {
		return 0, jtypes.ErrUndefined
	}

	var sum float64

	for i := 0; i < v.Len(); i++ {
		n, ok := jtypes.AsNumber(v.Index(i))
		if !ok {
			return 0, fmt.Errorf("cannot call average on an array with non-number types")
		}
		sum += n
	}

	return sum / float64(v.Len()), nil
}
