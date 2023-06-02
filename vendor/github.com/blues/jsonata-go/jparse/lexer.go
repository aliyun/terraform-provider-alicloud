// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jparse

import (
	"fmt"
	"unicode/utf8"
)

const eof = -1

type tokenType uint8

const (
	typeEOF tokenType = iota
	typeError

	typeString   // string literal, e.g. "hello"
	typeNumber   // number literal, e.g. 3.14159
	typeBoolean  // true or false
	typeNull     // null
	typeName     // field name, e.g. Price
	typeNameEsc  // escaped field name, e.g. `Product Name`
	typeVariable // variable, e.g. $x
	typeRegex    // regular expression, e.g. /ab+/

	// Symbol operators
	typeBracketOpen
	typeBracketClose
	typeBraceOpen
	typeBraceClose
	typeParenOpen
	typeParenClose
	typeDot
	typeComma
	typeColon
	typeSemicolon
	typeCondition
	typePlus
	typeMinus
	typeMult
	typeDiv
	typeMod
	typePipe
	typeEqual
	typeNotEqual
	typeLess
	typeLessEqual
	typeGreater
	typeGreaterEqual
	typeApply
	typeSort
	typeConcat
	typeRange
	typeAssign
	typeDescendent

	// Keyword operators
	typeAnd
	typeOr
	typeIn
)

func (tt tokenType) String() string {
	switch tt {
	case typeEOF:
		return "(eof)"
	case typeError:
		return "(error)"
	case typeString:
		return "(string)"
	case typeNumber:
		return "(number)"
	case typeBoolean:
		return "(boolean)"
	case typeName, typeNameEsc:
		return "(name)"
	case typeVariable:
		return "(variable)"
	case typeRegex:
		return "(regex)"
	default:
		if s := symbolsAndKeywords[tt]; s != "" {
			return s
		}
		return "(unknown)"
	}
}

// symbols1 maps 1-character symbols to the corresponding
// token types.
var symbols1 = [...]tokenType{
	'[': typeBracketOpen,
	']': typeBracketClose,
	'{': typeBraceOpen,
	'}': typeBraceClose,
	'(': typeParenOpen,
	')': typeParenClose,
	'.': typeDot,
	',': typeComma,
	';': typeSemicolon,
	':': typeColon,
	'?': typeCondition,
	'+': typePlus,
	'-': typeMinus,
	'*': typeMult,
	'/': typeDiv,
	'%': typeMod,
	'|': typePipe,
	'=': typeEqual,
	'<': typeLess,
	'>': typeGreater,
	'^': typeSort,
	'&': typeConcat,
}

type runeTokenType struct {
	r  rune
	tt tokenType
}

// symbols2 maps 2-character symbols to the corresponding
// token types.
var symbols2 = [...][]runeTokenType{
	'!': {{'=', typeNotEqual}},
	'<': {{'=', typeLessEqual}},
	'>': {{'=', typeGreaterEqual}},
	'.': {{'.', typeRange}},
	'~': {{'>', typeApply}},
	':': {{'=', typeAssign}},
	'*': {{'*', typeDescendent}},
}

const (
	symbol1Count = rune(len(symbols1))
	symbol2Count = rune(len(symbols2))
)

func lookupSymbol1(r rune) tokenType {
	if r < 0 || r >= symbol1Count {
		return 0
	}
	return symbols1[r]
}

func lookupSymbol2(r rune) []runeTokenType {
	if r < 0 || r >= symbol2Count {
		return nil
	}
	return symbols2[r]
}

func lookupKeyword(s string) tokenType {
	switch s {
	case "and":
		return typeAnd
	case "or":
		return typeOr
	case "in":
		return typeIn
	case "true", "false":
		return typeBoolean
	case "null":
		return typeNull
	default:
		return 0
	}
}

// A token represents a discrete part of a JSONata expression
// such as a string, a number, a field name, or an operator.
type token struct {
	Type     tokenType
	Value    string
	Position int
}

// lexer converts a JSONata expression into a sequence of tokens.
// The implmentation is based on the technique described in Rob
// Pike's 'Lexical Scanning in Go' talk.
type lexer struct {
	input   string
	length  int
	start   int
	current int
	width   int
	err     error
}

// newLexer creates a new lexer from the provided input. The
// input is tokenized by successive calls to the next method.
func newLexer(input string) lexer {
	return lexer{
		input:  input,
		length: len(input),
	}
}

// next returns the next token from the provided input. When
// the end of the input is reached, next returns EOF for all
// subsequent calls.
//
// The allowRegex argument determines how the lexer interprets
// a forward slash character. Forward slashes in JSONata can
// either be the start of a regular expression or the division
// operator depending on their position. If allowRegex is true,
// the lexer will treat a forward slash like a regular
// expression.
func (l *lexer) next(allowRegex bool) token {

	l.skipWhitespace()

	ch := l.nextRune()
	if ch == eof {
		return l.eof()
	}

	if allowRegex && ch == '/' {
		l.ignore()
		return l.scanRegex(ch)
	}

	if rts := lookupSymbol2(ch); rts != nil {
		for _, rt := range rts {
			if l.acceptRune(rt.r) {
				return l.newToken(rt.tt)
			}
		}
	}

	if tt := lookupSymbol1(ch); tt > 0 {
		return l.newToken(tt)
	}

	if ch == '"' || ch == '\'' {
		l.ignore()
		return l.scanString(ch)
	}

	if ch >= '0' && ch <= '9' {
		l.backup()
		return l.scanNumber()
	}

	if ch == '`' {
		l.ignore()
		return l.scanEscapedName(ch)
	}

	l.backup()
	return l.scanName()
}

// scanRegex reads a regular expression from the current position
// and returns a regex token. The opening delimiter has already
// been consumed.
func (l *lexer) scanRegex(delim rune) token {

	var depth int

Loop:
	for {
		switch l.nextRune() {
		case delim:
			if depth == 0 {
				break Loop
			}
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		case '\\':
			if r := l.nextRune(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.error(ErrUnterminatedRegex, string(delim))
		}
	}

	l.backup()
	t := l.newToken(typeRegex)
	l.acceptRune(delim)
	l.ignore()

	// Convert JavaScript-style regex flags to Go format,
	// e.g. /ab+/i becomes /(?i)ab+/.
	if l.acceptAll(isRegexFlag) {
		flags := l.newToken(0)
		t.Value = fmt.Sprintf("(?%s)%s", flags.Value, t.Value)
	}

	return t
}

// scanString reads a string literal from the current position
// and returns a string token. The opening quote has already been
// consumed.
func (l *lexer) scanString(quote rune) token {
Loop:
	for {
		switch l.nextRune() {
		case quote:
			break Loop
		case '\\':
			if r := l.nextRune(); r != eof {
				break
			}
			fallthrough
		case eof:
			return l.error(ErrUnterminatedString, string(quote))
		}
	}

	l.backup()
	t := l.newToken(typeString)
	l.acceptRune(quote)
	l.ignore()
	return t
}

// scanNumber reads a number literal from the current position
// and returns a number token.
func (l *lexer) scanNumber() token {

	// JSON does not support leading zeroes. The integer part of
	// a number will either be a single zero, or a non-zero digit
	// followed by zero or more digits.
	if !l.acceptRune('0') {
		l.accept(isNonZeroDigit)
		l.acceptAll(isDigit)
	}
	if l.acceptRune('.') {
		if !l.acceptAll(isDigit) {
			// If there are no digits after the decimal point,
			// don't treat the dot as part of the number. It
			// could be part of the range operator, e.g. "1..5".
			l.backup()
			return l.newToken(typeNumber)
		}
	}
	if l.acceptRunes2('e', 'E') {
		l.acceptRunes2('+', '-')
		l.acceptAll(isDigit)
	}
	return l.newToken(typeNumber)
}

// scanEscapedName reads a field name from the current position
// and returns a name token. The opening quote has already been
// consumed.
func (l *lexer) scanEscapedName(quote rune) token {
Loop:
	for {
		switch l.nextRune() {
		case quote:
			break Loop
		case eof, '\n':
			return l.error(ErrUnterminatedName, string(quote))
		}
	}

	l.backup()
	t := l.newToken(typeNameEsc)
	l.acceptRune(quote)
	l.ignore()
	return t
}

// scanName reads from the current position and returns a name,
// variable, or keyword token.
func (l *lexer) scanName() token {

	isVar := l.acceptRune('$')
	if isVar {
		l.ignore()
	}

	for {
		ch := l.nextRune()
		if ch == eof {
			break
		}

		// Stop reading if we hit whitespace...
		if isWhitespace(ch) {
			l.backup()
			break
		}

		// ...or anything that looks like an operator.
		if lookupSymbol1(ch) > 0 || lookupSymbol2(ch) != nil {
			l.backup()
			break
		}
	}

	t := l.newToken(typeName)

	if isVar {
		t.Type = typeVariable
	} else if tt := lookupKeyword(t.Value); tt > 0 {
		t.Type = tt
	}

	return t
}

func (l *lexer) eof() token {
	return token{
		Type:     typeEOF,
		Position: l.current,
	}
}

func (l *lexer) error(typ ErrType, hint string) token {
	t := l.newToken(typeError)
	l.err = newErrorHint(typ, t, hint)
	return t
}

func (l *lexer) newToken(tt tokenType) token {
	t := token{
		Type:     tt,
		Value:    l.input[l.start:l.current],
		Position: l.start,
	}
	l.width = 0
	l.start = l.current
	return t
}

func (l *lexer) nextRune() rune {

	if l.err != nil || l.current >= l.length {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.current:])
	l.width = w
	l.current += w
	/*
		if r == '\n' {
			l.line++
		}
	*/
	return r
}

func (l *lexer) backup() {
	// TODO: Support more than one backup operation.
	// TODO: Store current rune so that when nextRune
	// is called again, we don't need to repeat the call
	// to DecodeRuneInString.
	l.current -= l.width
}

func (l *lexer) ignore() {
	l.start = l.current
}

func (l *lexer) acceptRune(r rune) bool {
	return l.accept(func(c rune) bool {
		return c == r
	})
}

func (l *lexer) acceptRunes2(r1, r2 rune) bool {
	return l.accept(func(c rune) bool {
		return c == r1 || c == r2
	})
}

func (l *lexer) accept(isValid func(rune) bool) bool {
	if isValid(l.nextRune()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptAll(isValid func(rune) bool) bool {
	var b bool
	for l.accept(isValid) {
		b = true
	}
	return b
}

func (l *lexer) skipWhitespace() {
	l.acceptAll(isWhitespace)
	l.ignore()
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r', '\v':
		return true
	default:
		return false
	}
}

func isRegexFlag(r rune) bool {
	switch r {
	case 'i', 'm', 's':
		return true
	default:
		return false
	}
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isNonZeroDigit(r rune) bool {
	return r >= '1' && r <= '9'
}

// symbolsAndKeywords maps operator token types back to their
// string representations. It's only used by tokenType.String
// (and one test).
var symbolsAndKeywords = func() map[tokenType]string {

	m := map[tokenType]string{
		typeAnd:  "and",
		typeOr:   "or",
		typeIn:   "in",
		typeNull: "null",
	}

	for r, tt := range symbols1 {
		if tt > 0 {
			m[tt] = fmt.Sprintf("%c", r)
		}
	}

	for r, rts := range symbols2 {
		for _, rt := range rts {
			m[rt.tt] = fmt.Sprintf("%c", r) + fmt.Sprintf("%c", rt.r)
		}
	}

	return m
}()
