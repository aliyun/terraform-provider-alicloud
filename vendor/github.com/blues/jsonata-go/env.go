// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jsonata

import (
	"errors"
	"math"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/blues/jsonata-go/jlib"
	"github.com/blues/jsonata-go/jparse"
	"github.com/blues/jsonata-go/jtypes"
)

type environment struct {
	parent  *environment
	symbols map[string]reflect.Value
}

func newEnvironment(parent *environment, size int) *environment {
	return &environment{
		parent:  parent,
		symbols: make(map[string]reflect.Value, size),
	}
}

func (s *environment) bind(name string, value reflect.Value) {
	if s.symbols == nil {
		s.symbols = make(map[string]reflect.Value)
	}
	s.symbols[name] = value
}

func (s *environment) bindAll(values map[string]reflect.Value) {

	if len(values) == 0 {
		return
	}

	for name, value := range values {
		s.bind(name, value)
	}
}

func (s *environment) lookup(name string) reflect.Value {

	if v, ok := s.symbols[name]; ok {
		return v
	}
	if s.parent != nil {
		return s.parent.lookup(name)
	}

	return undefined
}

var (
	defaultUndefinedHandler = jtypes.ArgUndefined(0)
	defaultContextHandler   = jtypes.ArgCountEquals(0)

	argCountEquals1 = jtypes.ArgCountEquals(1)
)

var baseEnv = initBaseEnv(map[string]Extension{

	// String functions

	"string": {
		Func:               jlib.String,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"length": {
		Func:               utf8.RuneCountInString,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"substring": {
		Func:               jlib.Substring,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: contextHandlerSubstring,
	},
	"substringBefore": {
		Func:               jlib.SubstringBefore,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: contextHandlerSubstringBeforeAfter,
	},
	"substringAfter": {
		Func:               jlib.SubstringAfter,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: contextHandlerSubstringBeforeAfter,
	},
	"uppercase": {
		Func:               strings.ToUpper,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"lowercase": {
		Func:               strings.ToLower,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"pad": {
		Func:               jlib.Pad,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: contextHandlerPad,
	},
	"trim": {
		Func:               jlib.Trim,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"contains": {
		Func:               jlib.Contains,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: argCountEquals1,
	},
	"split": {
		Func:               jlib.Split,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: contextHandlerSplit,
	},
	"join": {
		Func:               jlib.Join,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"match": {
		Func:               jlib.Match,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: contextHandlerMatch,
	},
	"replace": {
		Func:               jlib.Replace,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: contextHandlerReplace,
	},
	"formatNumber": {
		Func:               jlib.FormatNumber,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: contextHandlerFormatNumber,
	},
	"formatBase": {
		Func:               jlib.FormatBase,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"base64encode": {
		Func:               jlib.Base64Encode,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"base64decode": {
		Func:               jlib.Base64Decode,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"decodeUrl": {
		Func:               jlib.DecodeURL,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"decodeUrlComponent": {
		Func:               jlib.DecodeURL,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"encodeUrl": {
		Func:               jlib.EncodeURL,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"encodeUrlComponent": {
		Func:               jlib.EncodeURLComponent,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},

	// Number functions

	"number": {
		Func:               jlib.Number,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"abs": {
		Func:               math.Abs,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"floor": {
		Func:               math.Floor,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"ceil": {
		Func:               math.Ceil,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"round": {
		Func:               jlib.Round,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"power": {
		Func:               jlib.Power,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: argCountEquals1,
	},
	"sqrt": {
		Func:               jlib.Sqrt,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"random": {
		Func:               jlib.Random,
		UndefinedHandler:   nil,
		EvalContextHandler: nil,
	},

	// Number aggregation functions

	"sum": {
		Func:               jlib.Sum,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"max": {
		Func:               jlib.Max,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"min": {
		Func:               jlib.Min,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"average": {
		Func:               jlib.Average,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},

	// Boolean functions

	"boolean": {
		Func:               jlib.Boolean,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"not": {
		Func:               jlib.Not,
		UndefinedHandler:   nil,
		EvalContextHandler: defaultContextHandler,
	},
	"exists": {
		Func:               jlib.Exists,
		UndefinedHandler:   nil,
		EvalContextHandler: nil,
	},

	// Array functions

	"distinct": {
		Func:               jlib.Distinct,
		UndefinedHandler:   nil,
		EvalContextHandler: nil,
	},
	"count": {
		Func:               jlib.Count,
		UndefinedHandler:   nil,
		EvalContextHandler: nil,
	},
	"reverse": {
		Func:               jlib.Reverse,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"sort": {
		Func:               jlib.Sort,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"shuffle": {
		Func:               jlib.Shuffle,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"zip": {
		Func:               jlib.Zip,
		UndefinedHandler:   nil,
		EvalContextHandler: nil,
	},
	"append": {
		Func:               jlib.Append,
		UndefinedHandler:   undefinedHandlerAppend,
		EvalContextHandler: nil,
	},
	"map": {
		Func:               jlib.Map,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"filter": {
		Func:               jlib.Filter,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"reduce": {
		Func:               jlib.Reduce,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},
	"single": {
		Func:               jlib.Single,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},

	// Object functions

	"each": {
		Func:               jlib.Each,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"sift": {
		Func:               jlib.Sift,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: argCountEquals1,
	},
	"keys": {
		Func:               jlib.Keys,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"lookup": {
		Func:               lookup,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"spread": {
		Func:               jlib.Spread,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"merge": {
		Func:               jlib.Merge,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: nil,
	},

	// Date functions
	// The date functions $now and $millis are not included
	// in the base environment because they use the current
	// time. They're added to the evaluation environment at
	// runtime.

	"fromMillis": {
		Func:               jlib.FromMillis,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},
	"toMillis": {
		Func:               jlib.ToMillis,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},

	"type": {
		Func:               jlib.TypeOf,
		UndefinedHandler:   defaultUndefinedHandler,
		EvalContextHandler: defaultContextHandler,
	},

	// Misc functions

	"error": {
		Func:               throw,
		UndefinedHandler:   nil,
		EvalContextHandler: nil,
	},
})

func initBaseEnv(exts map[string]Extension) *environment {

	env := newEnvironment(nil, len(exts))

	for name, ext := range exts {
		fn := mustGoCallable(name, ext)
		env.bind(name, reflect.ValueOf(fn))
	}

	return env
}

func mustGoCallable(name string, ext Extension) *goCallable {

	callable, err := newGoCallable(name, ext)
	if err != nil {
		panicf("%s is not a valid function: %s", name, err)
	}

	return callable
}

// Local functions (not from external packages)

func lookup(v reflect.Value, name string) (interface{}, error) {

	res, err := evalName(&jparse.NameNode{Value: name}, v, nil)
	if err != nil {
		return nil, err
	}

	if seq, ok := asSequence(res); ok {
		res = seq.Value()
	}

	if res.IsValid() && res.CanInterface() {
		return res.Interface(), nil
	}

	return nil, nil
}

func throw(msg string) (interface{}, error) {
	return nil, errors.New(msg)
}

// Undefined handlers

func undefinedHandlerAppend(argv []reflect.Value) bool {
	return len(argv) == 2 && argv[0] == undefined && argv[1] == undefined
}

// Context handlers

func contextHandlerSubstring(argv []reflect.Value) bool {

	// If substring() is called with one or two numeric arguments,
	// use the evaluation context as the first argument.
	switch len(argv) {
	case 1:
		return jtypes.IsNumber(argv[0])
	case 2:
		return jtypes.IsNumber(argv[0]) && jtypes.IsNumber(argv[1])
	default:
		return false
	}
}

func contextHandlerSubstringBeforeAfter(argv []reflect.Value) bool {

	// If subStringBefore() or subStringAfter() are called with
	// one string argument, use the evaluation context as the first
	// argument.
	return len(argv) == 1 && jtypes.IsString(argv[0])
}

func contextHandlerPad(argv []reflect.Value) bool {

	// If pad() is called with a single number, or a number and
	// a string, use the evaluation context as the first argument.
	switch len(argv) {
	case 1:
		return jtypes.IsNumber(argv[0])
	case 2:
		return jtypes.IsNumber(argv[0]) && jtypes.IsString(argv[1])
	default:
		return false
	}
}

func contextHandlerSplit(argv []reflect.Value) bool {

	// If split() is called with a single string/regex, or a
	// string/regex and a number, use the evaluation context as
	// the first argument.
	switch len(argv) {
	case 1:
		return isStringOrCallable(argv[0])
	case 2:
		return isStringOrCallable(argv[0]) && jtypes.IsNumber(argv[1])
	default:
		return false
	}
}

func contextHandlerMatch(argv []reflect.Value) bool {

	// If match() is called with a single regex, or a regex and
	// a number, use the evaluation context as the first argument.
	switch len(argv) {
	case 1:
		return jtypes.IsCallable(argv[0])
	case 2:
		return jtypes.IsCallable(argv[0]) && jtypes.IsNumber(argv[1])
	default:
		return false
	}
}

func contextHandlerReplace(argv []reflect.Value) bool {

	// If replace() is called with a string/regex and a string/Callable,
	// or a string/regex, a string/Callable, and a number, use the
	// evaluation context as the first argument.
	switch len(argv) {
	case 2:
		return isStringOrCallable(argv[0]) && isStringOrCallable(argv[1])
	case 3:
		return isStringOrCallable(argv[0]) && isStringOrCallable(argv[1]) && jtypes.IsNumber(argv[2])
	default:
		return false
	}
}

func contextHandlerFormatNumber(argv []reflect.Value) bool {

	// If formatNumber() is called with one or two arguments, and
	// the first argument is a string, use the evaluation context
	// as the first argument.
	switch len(argv) {
	case 1, 2:
		return jtypes.IsString(argv[0])
	default:
		return false
	}
}

func isStringOrCallable(v reflect.Value) bool {
	return jtypes.IsString(v) || jtypes.IsCallable(v)
}
