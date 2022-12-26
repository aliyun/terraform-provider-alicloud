// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Package jparse converts JSONata expressions to abstract
// syntax trees. Most clients will not need to work with
// this package directly.
//
// Usage
//
// Call the Parse function, passing a JSONata expression as
// a string. If an error occurs, it will be of type Error.
// Otherwise, Parse returns the root Node of the AST.
package jparse
