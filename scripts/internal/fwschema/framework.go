// Package fwschema synthesizes an SDKv2 *schema.Resource from a
// framework-native resource's Schema() method, so the SDKv2-based attribute
// coverage and collection-order checkers apply to framework resources with the
// same fail-closed semantics as SDKv2 ones.
//
// It lives under scripts/internal so both scripts/collection-order and
// scripts/testing can share one implementation instead of each duplicating the
// framework AST walk (which would drift) or skipping framework resources as
// "outside scope" (which lets an uncovered ordered collection or an
// unregistered attribute through silently).
package fwschema

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FrameworkResource synthesizes an SDKv2 *schema.Resource from a
// framework-native resource's Schema() method, reading the implementation
// file from disk. A parse failure, a missing Schema method or an unrecognized
// attribute/block kind is reported as an error: the caller fail-closes rather
// than silently treating a framework resource as having nothing to cover (which
// would let an uncovered ordered collection or an untested attribute through).
func FrameworkResource(implPath string) (*schema.Resource, error) {
	src, err := os.ReadFile(implPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", implPath, err)
	}
	return FrameworkResourceSource(implPath, src)
}

// FrameworkResourceSource is the file-read-free core of FrameworkResource,
// exposed so tests can feed in-memory source.
func FrameworkResourceSource(implPath string, src []byte) (*schema.Resource, error) {
	f, err := parser.ParseFile(token.NewFileSet(), implPath, src, parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", implPath, err)
	}
	schemaLit, err := findSchemaLiteral(f)
	if err != nil {
		return nil, err
	}
	r := &schema.Resource{Schema: map[string]*schema.Schema{}}
	if err := frameworkFields(schemaLit, r.Schema); err != nil {
		return nil, err
	}
	return r, nil
}

// findSchemaLiteral locates the schema.Schema composite literal assigned in the
// resource's Schema() method. It does not assume a particular receiver or
// response variable name, and resolves the framework schema package alias
// rather than hard-coding "schema".
func findSchemaLiteral(f *ast.File) (*ast.CompositeLit, error) {
	schemaPkg := importName(f, "github.com/hashicorp/terraform-plugin-framework/resource/schema")
	if schemaPkg == "" {
		return nil, fmt.Errorf("resource does not import terraform-plugin-framework/resource/schema by name")
	}
	var found *ast.CompositeLit
	var missing bool
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "Schema" || fn.Recv == nil || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if found != nil {
				return false
			}
			lit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			if selectorIs(lit.Type, schemaPkg, "Schema") {
				found = lit
				return false
			}
			return true
		})
		if found != nil {
			return found, nil
		}
		missing = true
	}
	if missing {
		return nil, fmt.Errorf("Schema() method declared no schema.Schema literal")
	}
	return nil, fmt.Errorf("no Schema() method with a receiver was found")
}

// frameworkFields fills a schema map from a schema.Schema (or NestedBlockObject
// / SingleNestedBlock) composite literal's Attributes and Blocks entries.
func frameworkFields(lit *ast.CompositeLit, out map[string]*schema.Schema) error {
	attrs := fields(lit)
	if attrMap, ok := attrs["Attributes"]; ok {
		if err := frameworkAttributeMap(attrMap, out); err != nil {
			return err
		}
	}
	if blockMap, ok := attrs["Blocks"]; ok {
		if err := frameworkBlockMap(blockMap, out); err != nil {
			return err
		}
	}
	return nil
}

func frameworkAttributeMap(expr ast.Expr, out map[string]*schema.Schema) error {
	pairs, err := keyedLiterals(expr)
	if err != nil {
		return fmt.Errorf("Attributes: %w", err)
	}
	for name, lit := range pairs {
		s, err := frameworkAttribute(lit)
		if err != nil {
			return fmt.Errorf("attribute %q: %w", name, err)
		}
		out[name] = s
	}
	return nil
}

func frameworkBlockMap(expr ast.Expr, out map[string]*schema.Schema) error {
	pairs, err := keyedLiterals(expr)
	if err != nil {
		return fmt.Errorf("Blocks: %w", err)
	}
	for name, lit := range pairs {
		s, err := frameworkBlock(lit)
		if err != nil {
			return fmt.Errorf("block %q: %w", name, err)
		}
		out[name] = s
	}
	return nil
}

// frameworkAttribute maps one framework attribute composite literal to an SDKv2
// schema entry. It reads the real Required/Optional/Computed flags from the
// literal instead of forcing Optional+Computed on every kind: a Required
// attribute stays Required, and a Computed-only (output) attribute stays
// Computed without Optional, so the attribute-coverage and collection-order
// checkers see the same configurability the framework declares rather than a
// synthesized Optional+Computed that hides a Required or a read-only field.
func frameworkAttribute(lit *ast.CompositeLit) (*schema.Schema, error) {
	kind := selectorSel(lit.Type)
	fl := fields(lit)
	required, err := boolField(fl, "Required")
	if err != nil {
		return nil, err
	}
	optional, err := boolField(fl, "Optional")
	if err != nil {
		return nil, err
	}
	computed, err := boolField(fl, "Computed")
	if err != nil {
		return nil, err
	}
	switch kind {
	case "StringAttribute":
		return &schema.Schema{Type: schema.TypeString, Required: required, Optional: optional, Computed: computed}, nil
	case "Int64Attribute":
		return &schema.Schema{Type: schema.TypeInt, Required: required, Optional: optional, Computed: computed}, nil
	case "BoolAttribute":
		return &schema.Schema{Type: schema.TypeBool, Required: required, Optional: optional, Computed: computed}, nil
	case "Float64Attribute":
		return &schema.Schema{Type: schema.TypeFloat, Required: required, Optional: optional, Computed: computed}, nil
	case "DynamicAttribute":
		return &schema.Schema{Type: schema.TypeString, Required: required, Optional: optional, Computed: computed}, nil
	case "NumberAttribute":
		return &schema.Schema{Type: schema.TypeFloat, Required: required, Optional: optional, Computed: computed}, nil
	case "ListAttribute", "SetAttribute", "MapAttribute":
		// Primitive collections without a nested object carry no re-orderable
		// members; they are never candidates (ListAttribute maps to TypeList
		// with no Elem).
		t := schema.TypeList
		if kind == "SetAttribute" {
			t = schema.TypeSet
		} else if kind == "MapAttribute" {
			t = schema.TypeMap
		}
		return &schema.Schema{Type: t, Required: required, Optional: optional, Computed: computed}, nil
	case "ListNestedAttribute", "SetNestedAttribute", "SingleNestedAttribute":
		return nestedAttribute(kind, lit, fl, required, optional, computed)
	case "":
		return nil, fmt.Errorf("attribute is not a framework schema.<Kind> composite literal")
	default:
		return nil, fmt.Errorf("unrecognized framework attribute kind: schema.%s", kind)
	}
}

func nestedAttribute(kind string, lit *ast.CompositeLit, fl map[string]ast.Expr, required, optional, computed bool) (*schema.Schema, error) {
	var nested *schema.Resource
	switch kind {
	case "ListNestedAttribute", "SetNestedAttribute":
		// ListNestedAttribute and SetNestedAttribute wrap a NestedObject
		// (schema.NestedAttribute) carrying the nested Attributes/Blocks; the
		// absence of that literal is a parse failure, not an empty schema.
		objLit, hasObj := compositeLitOf(fl["NestedObject"])
		if !hasObj {
			return nil, fmt.Errorf("schema.%s has no NestedObject composite literal", kind)
		}
		nested = &schema.Resource{Schema: map[string]*schema.Schema{}}
		if err := frameworkFields(objLit, nested.Schema); err != nil {
			return nil, err
		}
	case "SingleNestedAttribute":
		// SingleNestedAttribute has Attributes (and Blocks) directly, with no
		// NestedObject field — it mirrors SingleNestedBlock's shape. Reading the
		// nested fields off the attribute literal itself keeps a valid
		// SingleNestedAttribute from failing "has no NestedObject composite
		// literal", while a ListNestedAttribute that drops NestedObject still
		// fails closed.
		nested = &schema.Resource{Schema: map[string]*schema.Schema{}}
		if err := frameworkFields(lit, nested.Schema); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("schema.%s", kind)
	}
	switch kind {
	case "ListNestedAttribute":
		return &schema.Schema{Type: schema.TypeList, Required: required, Optional: optional, Computed: computed, Elem: nested, MaxItems: 0}, nil
	case "SetNestedAttribute":
		return &schema.Schema{Type: schema.TypeSet, Required: required, Optional: optional, Computed: computed, Elem: nested}, nil
	case "SingleNestedAttribute":
		return &schema.Schema{Type: schema.TypeList, Required: required, Optional: optional, Computed: computed, Elem: nested, MaxItems: 1}, nil
	}
	return nil, fmt.Errorf("schema.%s", kind)
}

func frameworkBlock(lit *ast.CompositeLit) (*schema.Schema, error) {
	kind := selectorSel(lit.Type)
	switch kind {
	case "ListNestedBlock":
		objLit, err := blockNestedObject(lit, kind)
		if err != nil {
			return nil, err
		}
		nested := &schema.Resource{Schema: map[string]*schema.Schema{}}
		if err := frameworkFields(objLit, nested.Schema); err != nil {
			return nil, err
		}
		return &schema.Schema{Type: schema.TypeList, Optional: true, Elem: nested, MaxItems: 0}, nil
	case "SetNestedBlock":
		objLit, err := blockNestedObject(lit, kind)
		if err != nil {
			return nil, err
		}
		nested := &schema.Resource{Schema: map[string]*schema.Schema{}}
		if err := frameworkFields(objLit, nested.Schema); err != nil {
			return nil, err
		}
		return &schema.Schema{Type: schema.TypeSet, Optional: true, Elem: nested}, nil
	case "SingleNestedBlock":
		nested := &schema.Resource{Schema: map[string]*schema.Schema{}}
		if err := frameworkFields(lit, nested.Schema); err != nil {
			return nil, err
		}
		return &schema.Schema{Type: schema.TypeList, Optional: true, Elem: nested, MaxItems: 1}, nil
	case "":
		return nil, fmt.Errorf("block is not a framework schema.<Kind> composite literal")
	default:
		return nil, fmt.Errorf("unrecognized framework block kind: schema.%s", kind)
	}
}

func blockNestedObject(lit *ast.CompositeLit, kind string) (*ast.CompositeLit, error) {
	nestedFields := fields(lit)
	objLit, ok := compositeLitOf(nestedFields["NestedObject"])
	if !ok {
		return nil, fmt.Errorf("schema.%s has no NestedObject composite literal", kind)
	}
	return objLit, nil
}

// keyedLiterals parses a map literal expression into ordered name→composite-literal
// pairs, preserving declaration order so diagnostics match the source.
func keyedLiterals(expr ast.Expr) (map[string]*ast.CompositeLit, error) {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("expected a composite literal, got %T", expr)
	}
	out := map[string]*ast.CompositeLit{}
	for _, el := range lit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			return nil, fmt.Errorf("expected a key/value element, got %T", el)
		}
		key := text(kv.Key)
		if key == "" {
			key = ident(kv.Key)
		}
		if key == "" {
			return nil, fmt.Errorf("attribute/block key is not a literal")
		}
		value, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			return nil, fmt.Errorf("value of %q is not a composite literal", key)
		}
		if _, exists := out[key]; exists {
			return nil, fmt.Errorf("duplicate key %q", key)
		}
		out[key] = value
	}
	return out, nil
}

func compositeLitOf(expr ast.Expr) (*ast.CompositeLit, bool) {
	lit, ok := expr.(*ast.CompositeLit)
	return lit, ok
}

// boolField reads a bool-valued composite-literal field (e.g. Required: true).
// A missing field defaults to false — the framework's zero value — so an
// attribute omitting Required/Optional/Computed is reported with that flag
// false rather than forced. An explicit true/false ident is read literally.
// Any other expression — a call like configurableFromAnotherFile(), a non-bool
// ident, a selector, a composite — is unresolved: the synthesizer cannot prove
// the flag is false, so it fails closed rather than defaulting to false and
// hiding a possibly-configurable ordered collection from Candidates. The error
// is propagated by frameworkAttribute so the whole resource fails to synthesize
// rather than silently dropping the attribute.
func boolField(flds map[string]ast.Expr, name string) (bool, error) {
	e, ok := flds[name]
	if !ok {
		return false, nil
	}
	i, ok := e.(*ast.Ident)
	if !ok {
		return false, fmt.Errorf("%s is %T, not a bool literal; an unresolved configurability expression must fail closed, not default to false", name, e)
	}
	switch i.Name {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("%s is ident %q, not a bool literal; an unresolved configurability expression must fail closed, not default to false", name, i.Name)
	}
}

func selectorIs(expr ast.Expr, pkg, name string) bool {
	s, ok := expr.(*ast.SelectorExpr)
	return ok && ident(s.X) == pkg && s.Sel.Name == name
}

func selectorSel(expr ast.Expr) string {
	if s, ok := expr.(*ast.SelectorExpr); ok {
		return s.Sel.Name
	}
	return ""
}

func importName(f *ast.File, path string) string {
	for _, imp := range f.Imports {
		if text(imp.Path) == path {
			if imp.Name != nil {
				if imp.Name.Name == "." || imp.Name.Name == "_" {
					return ""
				}
				return imp.Name.Name
			}
			parts := strings.Split(path, "/")
			return parts[len(parts)-1]
		}
	}
	return ""
}

func ident(expr ast.Expr) string {
	if i, ok := expr.(*ast.Ident); ok {
		return i.Name
	}
	return ""
}

func text(expr ast.Expr) string {
	if s, ok := expr.(*ast.BasicLit); ok && s.Kind == token.STRING {
		v, _ := strconv.Unquote(s.Value)
		return v
	}
	return ""
}

func fields(expr ast.Expr) map[string]ast.Expr {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	result := map[string]ast.Expr{}
	for _, e := range lit.Elts {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			return nil
		}
		key := ident(kv.Key)
		if key == "" {
			key = text(kv.Key)
		}
		if key == "" {
			return nil
		}
		if _, exists := result[key]; exists {
			return nil
		}
		result[key] = kv.Value
	}
	return result
}
