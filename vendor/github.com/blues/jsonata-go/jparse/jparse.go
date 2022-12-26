// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jparse

// The JSONata parser is based on Pratt's Top Down Operator
// Precededence algorithm (see https://tdop.github.io/). Given
// a series of tokens representing a JSONata expression and the
// following metadata, it converts the tokens into an abstract
// syntax tree:
//
// 1. Functions that convert tokens to nodes based on their
//    type and position (see 'nud' and 'led' in Pratt).
//
// 2. Binding powers (i.e. operator precedence values) for
//    infix operators (see 'lbp' in Pratt).
//
// This metadata is defined below.

// A nud (short for null denotation) is a function that takes
// a token and returns a node representing that token's value.
// The parsing algorithm only calls the nud function for tokens
// in the prefix position. This includes simple values like
// strings and numbers, complex values like arrays and objects,
// and prefix operators like the negation operator.
type nud func(*parser, token) (Node, error)

// An led (short for left denotation) is a function that takes
// a token and a node representing the left hand side of an
// infix operation, and returns a node representing that infix
// operation. The parsing algorithm only calls the led function
// for tokens in the infix position, e.g. the mathematical
// operators.
type led func(*parser, token, Node) (Node, error)

// nuds defines nud functions for token types that are valid
// in the prefix position.
var nuds = [...]nud{
	typeString:      parseString,
	typeNumber:      parseNumber,
	typeBoolean:     parseBoolean,
	typeNull:        parseNull,
	typeRegex:       parseRegex,
	typeVariable:    parseVariable,
	typeName:        parseName,
	typeNameEsc:     parseEscapedName,
	typeBracketOpen: parseArray,
	typeBraceOpen:   parseObject,
	typeParenOpen:   parseBlock,
	typeMult:        parseWildcard,
	typeMinus:       parseNegation,
	typeDescendent:  parseDescendent,
	typePipe:        parseObjectTransformation,
	typeIn:          parseName,
	typeAnd:         parseName,
	typeOr:          parseName,
}

// leds defines led functions for token types that are valid
// in the infix position.
var leds = [...]led{
	typeParenOpen:    parseFunctionCall,
	typeBracketOpen:  parsePredicate,
	typeBraceOpen:    parseGroup,
	typeCondition:    parseConditional,
	typeAssign:       parseAssignment,
	typeApply:        parseFunctionApplication,
	typeConcat:       parseStringConcatenation,
	typeSort:         parseSort,
	typeDot:          parseDot,
	typePlus:         parseNumericOperator,
	typeMinus:        parseNumericOperator,
	typeMult:         parseNumericOperator,
	typeDiv:          parseNumericOperator,
	typeMod:          parseNumericOperator,
	typeEqual:        parseComparisonOperator,
	typeNotEqual:     parseComparisonOperator,
	typeLess:         parseComparisonOperator,
	typeLessEqual:    parseComparisonOperator,
	typeGreater:      parseComparisonOperator,
	typeGreaterEqual: parseComparisonOperator,
	typeIn:           parseComparisonOperator,
	typeAnd:          parseBooleanOperator,
	typeOr:           parseBooleanOperator,
}

// bps defines binding powers for token types that are valid
// in the infix position. The parsing algorithm requires that
// all infix operators (as defined by the leds variable above)
// have a non-zero binding power.
//
// Binding powers are calculated from a 2D slice of token types
// in which the outer slice is ordered by operator precedence
// (highest to lowest) and each inner slice contains token
// types of equal operator precedence.
var bps = initBindingPowers([][]tokenType{
	{
		typeParenOpen,
		typeBracketOpen,
	},
	{
		typeDot,
	},
	{
		typeBraceOpen,
	},
	{
		typeMult,
		typeDiv,
		typeMod,
	},
	{
		typePlus,
		typeMinus,
		typeConcat,
	},
	{
		typeEqual,
		typeNotEqual,
		typeLess,
		typeLessEqual,
		typeGreater,
		typeGreaterEqual,
		typeIn,
		typeSort,
		typeApply,
	},
	{
		typeAnd,
	},
	{
		typeOr,
	},
	{
		typeCondition,
	},
	{
		typeAssign,
	},
})

const (
	nudCount = tokenType(len(nuds))
	ledCount = tokenType(len(leds))
)

func lookupNud(tt tokenType) nud {
	if tt >= nudCount {
		return nil
	}
	return nuds[tt]
}

func lookupLed(tt tokenType) led {
	if tt >= ledCount {
		return nil
	}
	return leds[tt]
}

func lookupBp(tt tokenType) int {
	if tt >= ledCount {
		return 0
	}
	return bps[tt]
}

// Parse builds the abstract syntax tree for a JSONata expression
// and returns the root node. If the provided expression is not
// valid, Parse returns an error of type Error.
func Parse(expr string) (root Node, err error) {

	// Handle panics from parseExpression.
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*Error); ok {
				root, err = nil, e
				return
			}
			panic(r)
		}
	}()

	p := newParser(expr)
	node := p.parseExpression(0)

	if p.token.Type != typeEOF {
		return nil, newError(ErrSyntaxError, p.token)
	}

	return node.optimize()
}

type parser struct {
	lexer lexer
	token token
	// The following function pointers are a workaround
	// for an initialisation loop compile error. See the
	// comment in newParser.
	lookupNud func(tokenType) nud
	lookupLed func(tokenType) led
	lookupBp  func(tokenType) int
}

func newParser(input string) parser {

	p := parser{
		lexer: newLexer(input),

		// Because the nuds/leds arrays refer to functions that
		// call the parser methods, the parser methods cannot
		// directly refer to the nuds/leds arrays. Specifically,
		// calling the nud/led lookup functions from the parser
		// causes an initialisation loop, e.g.
		//
		//     nuds refers to
		//     parseArray refers to
		//     parser.parseExpression refers to
		//     lookupNud refers to
		//     nuds
		//
		// To avoid this, the parser accesses the nud/led lookup
		// functions via function pointers set at runtime.
		lookupNud: lookupNud,
		lookupLed: lookupLed,
		lookupBp:  lookupBp,
	}

	// Set current token to the first token in the expression.
	p.advance(true)
	return p
}

// parseExpression is the central function of the Pratt
// algorithm. It handles dispatch to the various nud/led
// functions (which may call back into parseExpression
// and the other parser methods).
//
// Note that the parser methods, parseExpression included,
// panic instead of returning errors. Panics are caught
// by the top-level Parse function and returned to the
// caller as errors. This makes the nud/led functions
// nicer to write without sacrificing the public API.
func (p *parser) parseExpression(rbp int) Node {

	if p.token.Type == typeEOF {
		panic(newError(ErrUnexpectedEOF, p.token))
	}

	t := p.token
	p.advance(false)

	nud := p.lookupNud(t.Type)
	if nud == nil {
		panic(newError(ErrPrefix, t))
	}

	lhs, err := nud(p, t)
	if err != nil {
		panic(err)
	}

	for rbp < p.lookupBp(p.token.Type) {

		t := p.token
		p.advance(true)

		led := p.lookupLed(t.Type)
		if led == nil {
			panic(newError(ErrInfix, t))
		}

		lhs, err = led(p, t, lhs)
		if err != nil {
			panic(err)
		}
	}

	return lhs
}

// advance requests the next token from the lexer and updates
// the parser's current token pointer. It panics if the lexer
// returns an error token.
func (p *parser) advance(allowRegex bool) {
	p.token = p.lexer.next(allowRegex)
	if p.token.Type == typeError {
		panic(p.lexer.err)
	}
}

// consume is like advance except it first checks that the
// current token is of the expected type. It panics if that
// is not the case.
func (p *parser) consume(expected tokenType, allowRegex bool) {

	if p.token.Type != expected {

		typ := ErrUnexpectedToken
		if p.token.Type == typeEOF {
			typ = ErrMissingToken
		}

		panic(newErrorHint(typ, p.token, expected.String()))
	}

	p.advance(allowRegex)
}

// bp returns the binding power for the given token type.
func (p *parser) bp(t tokenType) int {
	return p.lookupBp(t)
}

// initBindingPowers calculates binding power values for the
// given token types and returns them as an array. The specific
// values are not important. All that matters for parsing is
// whether one token's binding power is higher than another's.
//
// Token types are provided as a slice of slices. The outer
// slice is ordered by operator precedence, highest to lowest.
// Token types within each inner slice have the same operator
// precedence.
func initBindingPowers(tokenTypes [][]tokenType) [ledCount]int {

	// Binding powers must:
	//
	//   1. be non-zero
	//   2. increase with operator precedence
	//   3. be separated by more than one (because we subtract
	//      1 from the binding power for right-associative
	//      operators).
	//
	// This function produces a minimum binding power of 10.
	// Values increase by 10 as operator precedence increases.

	var bps [ledCount]int

	for offset, tts := range tokenTypes {

		bp := (len(tokenTypes) - offset) * 10

		for _, tt := range tts {
			if bps[tt] != 0 {
				panicf("initBindingPowers: token type %d [%s] appears more than once", tt, tt)
			}
			bps[tt] = bp
		}
	}

	validateBindingPowers(bps)

	return bps
}

// validateBindingPowers sanity checks the values calculated
// by initBindingPowers. Every token type in the leds array
// should have a binding power. No other token type should
// have a binding power.
func validateBindingPowers(bps [ledCount]int) {

	for tt := tokenType(0); tt < ledCount; tt++ {
		if leds[tt] != nil && bps[tt] == 0 {
			panicf("validateBindingPowers: token type %d [%s] does not have a binding power", tt, tt)
		}
		if leds[tt] == nil && bps[tt] != 0 {
			panicf("validateBindingPowers: token type %d [%s] should not have a binding power", tt, tt)
		}
	}
}
