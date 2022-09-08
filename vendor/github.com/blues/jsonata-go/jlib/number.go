// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jlib

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/blues/jsonata-go/jtypes"
)

var reNumber = regexp.MustCompile(`^-?(([0-9]+))(\.[0-9]+)?([Ee][-+]?[0-9]+)?$`)

// Number converts values to numbers. Numeric values are returned
// unchanged. Strings in legal JSON number format are converted
// to the number they represent. Boooleans are converted to 0 or 1.
// All other types trigger an error.
func Number(value StringNumberBool) (float64, error) {
	v := reflect.Value(value)
	if b, ok := jtypes.AsBool(v); ok {
		if b {
			return 1, nil
		}
		return 0, nil
	}

	if n, ok := jtypes.AsNumber(v); ok {
		return n, nil
	}

	s, ok := jtypes.AsString(v)
	if ok && reNumber.MatchString(s) {
		if n, err := strconv.ParseFloat(s, 64); err == nil {
			return n, nil
		}
	}

	return 0, fmt.Errorf("unable to cast %q to a number", s)
}

// Round rounds its input to the number of decimal places given
// in the optional second parameter. By default, Round rounds to
// the nearest integer. A negative precision specifies which column
// to round to on the left hand side of the decimal place.
func Round(x float64, prec jtypes.OptionalInt) float64 {
	// Adapted from gonum's floats.RoundEven.
	// https://github.com/gonum/gonum/tree/master/floats

	if x == 0 {
		// Make sure zero is returned
		// without the negative bit set.
		return 0
	}
	// Fast path for positive precision on integers.
	if prec.Int >= 0 && x == math.Trunc(x) {
		return x
	}
	intermed := multByPow10(x, prec.Int)
	if math.IsInf(intermed, 0) {
		return x
	}
	if isHalfway(intermed) {
		correction, _ := math.Modf(math.Mod(intermed, 2))
		intermed += correction
		if intermed > 0 {
			x = math.Floor(intermed)
		} else {
			x = math.Ceil(intermed)
		}
	} else {
		if x < 0 {
			x = math.Ceil(intermed - 0.5)
		} else {
			x = math.Floor(intermed + 0.5)
		}
	}

	if x == 0 {
		return 0
	}

	return multByPow10(x, -prec.Int)
}

// Power returns x to the power of y.
func Power(x, y float64) (float64, error) {
	res := math.Pow(x, y)
	if math.IsInf(res, 0) || math.IsNaN(res) {
		return 0, fmt.Errorf("the power function has resulted in a value that cannot be represented as a JSON number")
	}
	return res, nil
}

// Sqrt returns the square root of a number. It returns an error
// if the number is less than zero.
func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, fmt.Errorf("the sqrt function cannot be applied to a negative number")
	}
	return math.Sqrt(x), nil
}

// Random returns a random floating point number between 0 and 1.
func Random() float64 {
	return rand.Float64()
}

// multByPow10 multiplies a number by 10 to the power of n.
// It does this by converting back and forth to strings to
// avoid floating point rounding errors, e.g.
//
//     4.525 * math.Pow10(2) returns 452.50000000000006
func multByPow10(x float64, n int) float64 {
	if n == 0 || math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}

	s := fmt.Sprintf("%g", x)

	chunks := strings.Split(s, "e")
	switch len(chunks) {
	case 1:
		s = chunks[0] + "e" + strconv.Itoa(n)
	case 2:
		e, _ := strconv.Atoi(chunks[1])
		s = chunks[0] + "e" + strconv.Itoa(e+n)
	default:
		return x
	}

	x, _ = strconv.ParseFloat(s, 64)
	return x
}

func isHalfway(x float64) bool {
	_, frac := math.Modf(x)
	frac = math.Abs(frac)
	return frac == 0.5 || (math.Nextafter(frac, math.Inf(-1)) < 0.5 && math.Nextafter(frac, math.Inf(1)) > 0.5)
}
