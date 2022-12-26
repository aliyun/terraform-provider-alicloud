// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jparse

import (
	"fmt"
	"regexp"
)

// ErrType describes the type of an error.
type ErrType uint

// Error types returned by the parser.
const (
	_ ErrType = iota
	ErrSyntaxError
	ErrUnexpectedEOF
	ErrUnexpectedToken
	ErrMissingToken
	ErrPrefix
	ErrInfix
	ErrUnterminatedString
	ErrUnterminatedRegex
	ErrUnterminatedName
	ErrIllegalEscape
	ErrIllegalEscapeHex
	ErrInvalidNumber
	ErrNumberRange
	ErrEmptyRegex
	ErrInvalidRegex
	ErrGroupPredicate
	ErrGroupGroup
	ErrPathLiteral
	ErrIllegalAssignment
	ErrIllegalParam
	ErrDuplicateParam
	ErrParamCount
	ErrInvalidUnionType
	ErrUnmatchedOption
	ErrUnmatchedSubtype
	ErrInvalidSubtype
	ErrInvalidParamType
)

var errmsgs = map[ErrType]string{
	ErrSyntaxError:        "syntax error: '{{token}}'",
	ErrUnexpectedEOF:      "unexpected end of expression",
	ErrUnexpectedToken:    "expected token '{{hint}}', got '{{token}}'",
	ErrMissingToken:       "expected token '{{hint}}' before end of expression",
	ErrPrefix:             "the symbol '{{token}}' cannot be used as a prefix operator",
	ErrInfix:              "the symbol '{{token}}' cannot be used as an infix operator",
	ErrUnterminatedString: "unterminated string literal (no closing '{{hint}}')",
	ErrUnterminatedRegex:  "unterminated regular expression (no closing '{{hint}}')",
	ErrUnterminatedName:   "unterminated name (no closing '{{hint}}')",
	ErrIllegalEscape:      "illegal escape sequence \\{{hint}}",
	ErrIllegalEscapeHex:   "illegal escape sequence \\{{hint}}: \\u must be followed by a 4-digit hexadecimal code point",
	ErrInvalidNumber:      "invalid number literal {{token}}",
	ErrNumberRange:        "invalid number literal {{token}}: value out of range",
	ErrEmptyRegex:         "invalid regular expression: expression cannot be empty",
	ErrInvalidRegex:       "invalid regular expression {{token}}: {{hint}}",
	ErrGroupPredicate:     "a predicate cannot follow a grouping expression in a path step",
	ErrGroupGroup:         "a path step can only have one grouping expression",
	ErrPathLiteral:        "invalid path step {{hint}}: paths cannot contain nulls, strings, numbers or booleans",
	ErrIllegalAssignment:  "illegal assignment: {{hint}} is not a variable",
	ErrIllegalParam:       "illegal function parameter: {{token}} is not a variable",
	ErrDuplicateParam:     "duplicate function parameter: {{token}}",
	ErrParamCount:         "invalid type signature: number of types must match number of function parameters",
	ErrInvalidUnionType:   "invalid type signature: unsupported union type '{{hint}}'",
	ErrUnmatchedOption:    "invalid type signature: option '{{hint}}' must follow a parameter",
	ErrUnmatchedSubtype:   "invalid type signature: subtypes must follow a parameter",
	ErrInvalidSubtype:     "invalid type signature: parameter type {{hint}} does not support subtypes",
	ErrInvalidParamType:   "invalid type signature: unknown parameter type '{{hint}}'",
}

var reErrMsg = regexp.MustCompile("{{(token|hint)}}")

// Error describes an error during parsing.
type Error struct {
	Type     ErrType
	Token    string
	Hint     string
	Position int
}

func newError(typ ErrType, tok token) error {
	return newErrorHint(typ, tok, "")
}

func newErrorHint(typ ErrType, tok token, hint string) error {
	return &Error{
		Type:     typ,
		Token:    tok.Value,
		Position: tok.Position,
		Hint:     hint,
	}
}

func (e Error) Error() string {

	s := errmsgs[e.Type]
	if s == "" {
		return fmt.Sprintf("parser.Error: unknown error type %d", e.Type)
	}

	return reErrMsg.ReplaceAllStringFunc(s, func(match string) string {
		switch match {
		case "{{token}}":
			return e.Token
		case "{{hint}}":
			return e.Hint
		default:
			return match
		}
	})
}

func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}
