// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jlib

import "fmt"

// ErrType (golint)
type ErrType uint

// ErrNanInf (golint)
const (
	_ ErrType = iota
	ErrNaNInf
)

// Error (golint)
type Error struct {
	Type ErrType
	Func string
}

// Error (golint)
func (e Error) Error() string {

	var msg string

	switch e.Type {
	case ErrNaNInf:
		msg = "cannot convert NaN/Infinity to string"
	default:
		msg = "unknown error"
	}

	return fmt.Sprintf("%s: %s", e.Func, msg)
}

func newError(name string, typ ErrType) *Error {
	return &Error{
		Func: name,
		Type: typ,
	}
}
