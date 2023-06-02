// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package jsonata

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/blues/jsonata-go/jlib"
	"github.com/blues/jsonata-go/jparse"
	"github.com/blues/jsonata-go/jtypes"
)

type callableName struct {
	name string
}

func (n callableName) Name() string {
	return n.name
}

func (n *callableName) SetName(s string) {
	n.name = s
}

type callableMarshaler struct{}

func (callableMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

type goCallableParam struct {
	t        reflect.Type
	isOpt    bool
	optType  *goCallableParam
	isVar    bool
	varTypes []goCallableParam
}

func newGoCallableParam(typ reflect.Type) goCallableParam {

	param := goCallableParam{
		t: typ,
	}

	isOpt := reflect.PtrTo(typ).Implements(jtypes.TypeOptional)
	if isOpt {
		o := reflect.New(typ).Interface().(jtypes.Optional)
		p := newGoCallableParam(o.Type())
		param.isOpt = true
		param.optType = &p
	}

	isVar := typ.Implements(jtypes.TypeVariant)
	if isVar {
		var ps []goCallableParam
		types := reflect.Zero(typ).Interface().(jtypes.Variant).ValidTypes()
		if n := len(types); n > 0 {
			ps = make([]goCallableParam, n)
			for i := range ps {
				ps[i] = newGoCallableParam(types[i])
			}
		}
		param.isVar = true
		param.varTypes = ps
	}

	return param
}

// A goCallable represents a built-in or third party Go function.
// It implements the Callable interface.
type goCallable struct {
	callableName
	callableMarshaler
	fn               reflect.Value
	params           []goCallableParam
	isVariadic       bool
	undefinedHandler jtypes.ArgHandler
	contextHandler   jtypes.ArgHandler
	context          reflect.Value
}

func newGoCallable(name string, ext Extension) (*goCallable, error) {

	if err := validateGoCallableFunc(ext.Func); err != nil {
		return nil, err
	}

	v := reflect.ValueOf(ext.Func)
	t := v.Type()

	params := makeGoCallableParams(t)
	if err := validateGoCallableParams(params, t.IsVariadic()); err != nil {
		return nil, err
	}

	return &goCallable{
		callableName: callableName{
			name: name,
		},
		fn:               v,
		params:           params,
		isVariadic:       t.IsVariadic(),
		undefinedHandler: ext.UndefinedHandler,
		contextHandler:   ext.EvalContextHandler,
	}, nil
}

var typeError = reflect.TypeOf((*error)(nil)).Elem()

func validateGoCallableFunc(fn interface{}) error {

	v := reflect.ValueOf(fn)

	if v.Kind() != reflect.Func {
		return fmt.Errorf("func must be a Go function")
	}

	t := v.Type()
	switch t.NumOut() {
	case 1:
	case 2:
		if !t.Out(1).Implements(typeError) {
			return fmt.Errorf("func must return an error as its second value")
		}
	default:
		return fmt.Errorf("func must return either 1 or 2 values")
	}

	return nil
}

func validateGoCallableParams(params []goCallableParam, isVariadic bool) error {

	var hasOptionals bool

	for i, p := range params {

		if p.isOpt && p.isVar {
			return fmt.Errorf("parameters cannot be both optional and variant")
		}

		if hasOptionals && !p.isOpt {
			return fmt.Errorf("a non-optional parameter cannot follow an optional parameter")
		}

		if p.isOpt {
			if p.optType.isOpt {
				return fmt.Errorf("optional parameters cannot have an optional underlying type")
			}
			if isVariadic && i == len(params)-1 {
				return fmt.Errorf("optional parameters cannot be variadic")
			}
			hasOptionals = true
		}

		if p.isVar {
			if !jtypes.TypeValue.ConvertibleTo(p.t) {
				return fmt.Errorf("variant parameter types must be derived from reflect.Value")
			}
			if len(p.varTypes) < 2 {
				return fmt.Errorf("variant parameters must have at least two valid types")
			}
			for _, t := range p.varTypes {
				if t.isOpt || t.isVar {
					return fmt.Errorf("a variant parameter's valid types cannot be optional or variant")
				}
			}
		}
	}

	return nil
}

func makeGoCallableParams(typ reflect.Type) []goCallableParam {

	paramCount := typ.NumIn()
	if paramCount == 0 {
		return nil
	}

	isVariadic := typ.IsVariadic()
	params := make([]goCallableParam, paramCount)

	for i := range params {

		t := typ.In(i)
		if isVariadic && i == paramCount-1 {
			// The type of the final parameter in a variadic
			// function is a slice of the declared type. Call
			// Elem to get the declared type.
			t = t.Elem()
		}

		params[i] = newGoCallableParam(t)
	}

	return params
}

func (c *goCallable) SetContext(context reflect.Value) {
	c.context = context
}

func (c *goCallable) ParamCount() int {
	return len(c.params)
}

func (c *goCallable) Call(argv []reflect.Value) (reflect.Value, error) {

	var err error

	argv, err = c.validateArgCount(argv)
	if err != nil {
		if err == jtypes.ErrUndefined {
			err = nil
		}
		return undefined, err
	}

	argv, err = c.validateArgTypes(argv)
	if err != nil {
		return undefined, err
	}

	results := c.fn.Call(argv)

	if len(results) == 2 && !results[1].IsNil() {
		err := results[1].Interface().(error)
		if err == jtypes.ErrUndefined {
			err = nil
		}
		return undefined, err
	}

	return results[0], nil
}

func (c *goCallable) validateArgCount(argv []reflect.Value) ([]reflect.Value, error) {

	argc := len(argv)

	if c.contextHandler != nil && c.contextHandler(argv) {
		// TODO: Return an error if the evaluation context
		// is not the correct type.
		newargv := make([]reflect.Value, 1, len(argv)+1)
		newargv[0] = c.context
		argv = append(newargv, argv...)
	}

	if c.undefinedHandler != nil && c.undefinedHandler(argv) {
		// TODO: Validate the other arguments before doing
		// this. Otherwise we mask errors with the other
		// arguments.
		return nil, jtypes.ErrUndefined
	}

	paramCount := len(c.params)

	for i := len(argv); i < paramCount; i++ {
		if !c.params[i].isOpt {
			break
		}
		argv = append(argv, undefined)
	}

	if c.isVariadic && len(argv) < paramCount-1 {
		return nil, newArgCountError(c, argc)
	}

	if !c.isVariadic && len(argv) != paramCount {
		return nil, newArgCountError(c, argc)
	}

	return argv, nil
}

func (c *goCallable) validateArgTypes(argv []reflect.Value) ([]reflect.Value, error) {

	var ok bool
	paramCount := len(c.params)

	for i, v := range argv {

		v = jtypes.Resolve(v)

		// The preceding call to Resolve dereferences pointers.
		// This is fine for most types but we need to restore
		// pointer type Callables.
		if v.Kind() == reflect.Struct &&
			reflect.PtrTo(v.Type()).Implements(jtypes.TypeCallable) {
			if v.CanAddr() {
				v = v.Addr()
			}
		}

		j := i
		// Variadic functions can have more arguments than
		// parameters. Use the type of the final parameter
		// to process any extra arguments.
		if j >= paramCount {
			j = paramCount - 1
		}

		v, ok = processGoCallableArg(v, c.params[j])
		if !ok {
			return nil, newArgTypeError(c, i+1)
		}

		argv[i] = v
	}

	return argv, nil
}

var (
	typeString    = reflect.TypeOf((*string)(nil)).Elem()
	typeByteSlice = reflect.TypeOf((*[]byte)(nil)).Elem()
)

func processGoCallableArg(arg reflect.Value, param goCallableParam) (reflect.Value, bool) {

	if arg == undefined {
		return processUndefinedArg(param)
	}

	if param.isOpt {
		return processOptionalArg(arg, param)
	}

	if param.isVar {
		return processVariantArg(arg, param)
	}

	argType := arg.Type()
	paramType := param.t

	switch {
	case argType == paramType:
		return arg, true
	case argType.AssignableTo(paramType):
		return arg, true
	case paramType == jtypes.TypeValue:
		return reflect.ValueOf(arg), true
	case argType.ConvertibleTo(paramType):
		// Only allow conversion to a string if the source type
		// is a byte slice. Go can convert other types (such as
		// integers) to strings but this is not supported in
		// JSONata.
		if paramType == typeString && argType != typeByteSlice {
			break
		}
		return arg.Convert(paramType), true
	case argType.Implements(jtypes.TypeConvertible):
		if arg.CanInterface() {
			return arg.Interface().(jtypes.Convertible).ConvertTo(paramType)
		}
	}

	return undefined, false
}

func processUndefinedArg(param goCallableParam) (reflect.Value, bool) {

	switch {
	case param.isOpt, param.t == jtypes.TypeInterface, param.t == jtypes.TypeValue:
		return reflect.Zero(param.t), true
	default:
		return undefined, false
	}
}

func processOptionalArg(arg reflect.Value, param goCallableParam) (reflect.Value, bool) {

	v, ok := processGoCallableArg(arg, *param.optType)
	if !ok {
		return undefined, false
	}

	opt := reflect.New(param.t).Interface().(jtypes.Optional)
	opt.Set(v)

	return reflect.ValueOf(opt).Elem(), true
}

func processVariantArg(arg reflect.Value, param goCallableParam) (reflect.Value, bool) {

	for _, t := range param.varTypes {
		if v, ok := processGoCallableArg(arg, t); ok {
			return reflect.ValueOf(v).Convert(param.t), true
		}
	}

	return undefined, false
}

// A lambdaCallable represents a user-defined JSONata function
// created with the 'function' keyword.
type lambdaCallable struct {
	callableName
	callableMarshaler
	body       jparse.Node
	paramNames []string
	typed      bool
	params     []jparse.Param
	env        *environment
	context    reflect.Value
}

func (f *lambdaCallable) ParamCount() int {
	return len(f.paramNames)
}

func (f *lambdaCallable) Call(argv []reflect.Value) (reflect.Value, error) {

	argv, err := f.validateArgs(argv)
	if err != nil {
		return undefined, err
	}

	// Create a local scope for this function's arguments.
	env := newEnvironment(f.env, len(f.paramNames))

	// Add the function arguments to the local scope.
	// If there are fewer arguments than parameter names,
	// default unset parameters to undefined. If there
	// are more arguments than parameter names, ignore
	// the extraneous arguments.
	for i, name := range f.paramNames {

		var v reflect.Value

		if i < len(argv) {
			v = argv[i]
		}

		env.bind(name, v)
	}

	// Evaluate the function body.
	return eval(f.body, f.context, env)
}

func (f *lambdaCallable) validateArgs(argv []reflect.Value) ([]reflect.Value, error) {

	// An untyped lambda can take any number of arguments
	// of any type. No further processing is required.
	if !f.typed {
		return argv, nil
	}

	var err error

	if argv, err = f.validateArgCount(argv); err != nil {
		return nil, err
	}

	if argv, err = f.validateArgTypes(argv); err != nil {
		return nil, err
	}

	return f.wrapVariadicArgs(argv), nil
}

func (f *lambdaCallable) validateArgCount(argv []reflect.Value) ([]reflect.Value, error) {

	// argc is the number of arguments originally passed to
	// the function.
	argc := len(argv)

	// paramCount is the number of parameters specified in
	// the function's type signature.
	paramCount := len(f.params)

	// If there are fewer arguments than parameters and the
	// first parameter is contextable, insert the evaluation
	// context into the argument list.
	if argc < paramCount && f.params[0].Option == jparse.ParamContextable {
		argv = append([]reflect.Value{f.context}, argv...)
	}

	// If there are still fewer arguments than parameters and
	// the missing arguments correspond to optional parameters,
	// append undefined arguments to the argument list.
	for i := len(argv); i < paramCount; i++ {
		if f.params[i].Option != jparse.ParamOptional {
			break
		}
		argv = append(argv, undefined)
	}

	// argCount is the final number of arguments including
	// any added by this method.
	argCount := len(argv)

	// isVar indicates whether the function is variadic.
	isVar := paramCount > 0 &&
		f.params[paramCount-1].Option == jparse.ParamVariadic

	// If there are a) fewer arguments than parameters or b)
	// extra arguments on a non-variadic function, return an
	// error.
	if argCount < paramCount || (argCount > paramCount && !isVar) {
		return nil, newArgCountError(f, argc)
	}

	return argv, nil
}

func (f *lambdaCallable) validateArgTypes(argv []reflect.Value) ([]reflect.Value, error) {

	paramCount := len(f.params)

	for i, arg := range argv {

		// Don't type check undefined arguments.
		if arg == undefined {
			continue
		}

		var param jparse.Param

		if i < paramCount {
			param = f.params[i]
		} else if paramCount > 0 {
			param = f.params[paramCount-1]
		}

		// If a parameter is an array type, force the
		// corresponding argument to be an array.
		if param.Type == jparse.ParamTypeArray {
			arg = arrayify(arg)
			argv[i] = arg
		}

		if !f.validArgType(arg, param) {
			return nil, newArgTypeError(f, i+1)
		}
	}

	return argv, nil
}

func (f *lambdaCallable) validArgType(arg reflect.Value, p jparse.Param) bool {

	typ := p.Type

	if typ&jparse.ParamTypeAny != 0 {
		return true
	}

	paramTypeJSON := typ&jparse.ParamTypeJSON != 0

	// TODO: Handle ParamTypeNull
	switch {
	case jtypes.IsString(arg):
		return paramTypeJSON || typ&jparse.ParamTypeString != 0
	case jtypes.IsNumber(arg):
		return paramTypeJSON || typ&jparse.ParamTypeNumber != 0
	case jtypes.IsBool(arg):
		return paramTypeJSON || typ&jparse.ParamTypeBool != 0
	case jtypes.IsCallable(arg):
		return typ&jparse.ParamTypeFunc != 0
	case jtypes.IsArray(arg):
		if paramTypeJSON {
			return true
		}
		if typ&jparse.ParamTypeArray != 0 {
			if len(p.SubParams) == 0 {
				return true
			}
			return jtypes.IsArrayOf(arg, func(v reflect.Value) bool {
				return f.validArgType(v, p.SubParams[0])
			})
		}
		return false
	case jtypes.IsMap(arg), jtypes.IsStruct(arg):
		return paramTypeJSON || typ&jparse.ParamTypeObject != 0
	}

	return false
}

func (f *lambdaCallable) wrapVariadicArgs(argv []reflect.Value) []reflect.Value {

	paramCount := len(f.params)

	if paramCount < 1 ||
		f.params[paramCount-1].Option != jparse.ParamVariadic {
		return argv
	}

	n := len(argv) - paramCount + 1
	vars := reflect.MakeSlice(typeInterfaceSlice, n, n)

	for i := 0; i < n; i++ {
		vars.Index(i).Set(argv[paramCount-1+i])
	}

	return append(argv[:paramCount-1], vars)
}

// A partialCallable represents the partial application of
// a Callable.
type partialCallable struct {
	callableName
	callableMarshaler
	fn      jtypes.Callable
	args    []jparse.Node
	env     *environment
	context reflect.Value
}

func (f *partialCallable) ParamCount() int {

	var count int
	for _, arg := range f.args {
		if _, ok := arg.(*jparse.PlaceholderNode); ok {
			count++
		}
	}

	return count
}

func (f *partialCallable) Call(argv []reflect.Value) (reflect.Value, error) {

	var err error
	args := make([]reflect.Value, len(f.args))

	for i, arg := range f.args {

		var v reflect.Value

		switch arg.(type) {
		case *jparse.PlaceholderNode:
			if len(argv) > 0 {
				v = argv[0]
				argv = argv[1:]
			}
		default:
			v, err = eval(arg, f.context, f.env)
			if err != nil {
				return undefined, err
			}
		}

		args[i] = v
	}

	return f.fn.Call(args)
}

// A transformationCallable represents JSONata's object
// transformation operator. It's a function that takes an
// object and updates and/or removes the specified keys.
type transformationCallable struct {
	callableName
	callableMarshaler
	pattern jparse.Node
	updates jparse.Node
	deletes jparse.Node
	env     *environment
}

func (f *transformationCallable) ParamCount() int {
	return 1
}

func (f *transformationCallable) Call(argv []reflect.Value) (reflect.Value, error) {

	err := f.validateArgs(argv)
	if err != nil {
		return undefined, err
	}

	obj, err := f.clone(argv[0])
	if err != nil {
		return undefined, newEvalError(ErrClone, nil, nil)
	}

	if obj == undefined {
		return undefined, nil
	}

	items, err := eval(f.pattern, obj, f.env)
	if err != nil {
		return undefined, err
	}

	items = arrayify(items)

	for i := 0; i < items.Len(); i++ {

		item := jtypes.Resolve(items.Index(i))
		if !jtypes.IsMap(item) {
			continue
		}

		if err := f.updateEntries(item); err != nil {
			return undefined, err
		}

		if f.deletes != nil {
			if err := f.deleteEntries(item); err != nil {
				return undefined, err
			}
		}
	}

	return obj, nil
}

func (f *transformationCallable) validateArgs(argv []reflect.Value) error {

	if argc := len(argv); argc != 1 {
		return newArgCountError(f, argc)
	}

	if obj := argv[0]; obj.IsValid() &&
		!jtypes.IsMap(obj) && !jtypes.IsStruct(obj) && !jtypes.IsArray(obj) {
		return newArgTypeError(f, 1)
	}

	return nil
}

func (f *transformationCallable) updateEntries(item reflect.Value) error {

	updates, err := eval(f.updates, item, f.env)
	if err != nil || updates == undefined {
		return err
	}

	if !jtypes.IsMap(updates) {
		return newEvalError(ErrIllegalUpdate, f.updates, nil)
	}

	for _, key := range updates.MapKeys() {
		item.SetMapIndex(key, updates.MapIndex(key))
	}

	return nil
}

func (f *transformationCallable) deleteEntries(item reflect.Value) error {

	deletes, err := eval(f.deletes, item, f.env)
	if err != nil || deletes == undefined {
		return err
	}

	deletes = arrayify(deletes)

	if !jtypes.IsArrayOf(deletes, jtypes.IsString) {
		return newEvalError(ErrIllegalDelete, f.deletes, nil)
	}

	for i := 0; i < deletes.Len(); i++ {
		key := jtypes.Resolve(deletes.Index(i))
		item.SetMapIndex(key, undefined)
	}

	return nil
}

func (f *transformationCallable) clone(v reflect.Value) (reflect.Value, error) {

	if v == undefined {
		return undefined, nil
	}

	s, err := jlib.String(v.Interface())
	if err != nil {
		return undefined, err
	}

	var dest interface{}
	d := json.NewDecoder(strings.NewReader(s))
	if err = d.Decode(&dest); err != nil {
		return undefined, err
	}

	return reflect.ValueOf(dest), nil
}

// A regexCallable represents a JSONata regular expression. It's
// a function that takes a string argument and returns an object
// that describes the leftmost match. The object also contains
// a Callable that returns the next leftmost match (and so on).
// A return value of undefined signifies no more matches.
type regexCallable struct {
	callableName
	callableMarshaler
	re *regexp.Regexp
}

func newRegexCallable(re *regexp.Regexp) *regexCallable {
	return &regexCallable{
		callableName: callableName{
			name: re.String(),
		},
		re: re,
	}
}

func (f *regexCallable) ParamCount() int {
	return 1
}

func (f *regexCallable) Call(argv []reflect.Value) (reflect.Value, error) {

	if len(argv) < 1 {
		return undefined, nil
	}

	s, ok := jtypes.AsString(argv[0])
	if !ok {
		return undefined, nil
	}

	matches, indexes := f.findMatches(s)
	return newMatchCallable(f.Name(), matches, indexes).Call(nil)
}

var typeRegexPtr = reflect.TypeOf((*regexp.Regexp)(nil))

func (f *regexCallable) ConvertTo(t reflect.Type) (reflect.Value, bool) {
	switch t {
	case typeRegexPtr:
		return reflect.ValueOf(f.re), true
	default:
		return undefined, false
	}
}

func (f *regexCallable) findMatches(s string) ([][]string, [][]int) {

	indexes := f.re.FindAllStringSubmatchIndex(s, -1)
	if indexes == nil {
		return nil, nil
	}

	matches := make([][]string, len(indexes))

	for i, index := range indexes {

		matches[i] = make([]string, len(index)/2)

		for j := range matches[i] {

			if index[j*2] < 0 {
				// Negative indexes indicate capturing groups
				// that don't match any text. Skip them.
				continue
			}
			matches[i][j] = s[index[j*2]:index[j*2+1]]
		}
	}

	return matches, indexes
}

// A matchCallable represents a regular expression match. Its
// Call method returns an object containing the details of the
// match, plus a Callable that returns the details of the next
// match.
type matchCallable struct {
	callableName
	callableMarshaler
	match  string
	start  int
	end    int
	groups []string
	next   jtypes.Callable
}

func newMatchCallable(name string, matches [][]string, indexes [][]int) jtypes.Callable {

	if len(matches) < 1 {
		return &undefinedCallable{
			callableName: callableName{
				name: name,
			},
		}
	}

	return &matchCallable{
		callableName: callableName{
			name: name,
		},
		match:  matches[0][0],
		start:  indexes[0][0],
		end:    indexes[0][1],
		groups: matches[0][1:],
		next:   newMatchCallable("next", matches[1:], indexes[1:]),
	}
}

func (f *matchCallable) Call([]reflect.Value) (reflect.Value, error) {
	return reflect.ValueOf(map[string]interface{}{
		"match":  f.match,
		"start":  f.start,
		"end":    f.end,
		"groups": f.groups,
		"next":   f.next,
	}), nil
}

func (*matchCallable) ParamCount() int {
	return 0
}

// An undefinedCallable is a Callable that always returns undefined.
type undefinedCallable struct {
	callableName
	callableMarshaler
}

func (*undefinedCallable) Call([]reflect.Value) (reflect.Value, error) {
	return undefined, nil
}

func (*undefinedCallable) ParamCount() int {
	return 0
}

// A chainCallable provides function composition.
type chainCallable struct {
	callableName
	callableMarshaler
	callables []jtypes.Callable
}

func (f *chainCallable) ParamCount() int {
	return 1
}

func (f *chainCallable) Call(argv []reflect.Value) (reflect.Value, error) {

	var err error
	var v reflect.Value

	if len(argv) > 0 {
		v = argv[0]
	}

	for _, fn := range f.callables {

		v, err = fn.Call([]reflect.Value{v})
		if err != nil {
			return undefined, err
		}
	}

	return v, nil
}
