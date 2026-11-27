// Package frameworkschema parses the terraform-plugin-framework schema shape
// from Go source.
//
// The framework declares its schema as
//
//	resp.Schema = schema.Schema{
//	    Attributes: map[string]schema.Attribute{...},
//	    Blocks:     map[string]schema.Block{...},
//	}
//
// A runtime-based gate cannot read it: framework entities are registered in
// conns.ServicePackage.Resources and friends, not in the SDKv2 provider maps.
// Both the breaking-change checker and the docs consistency checker therefore
// read the declaration from source, and they have to read it the same way —
// hence this shared package rather than two parsers that drift apart.
package frameworkschema

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Parse returns every attribute and block declared in the framework schema
// shape, keyed by its dotted path ("instance_endpoints.acl_enable").
//
// Attribute and block type names normalize to the SDKv2 schema.TypeXxx
// constants they play the role of, and a RequiresReplace plan modifier
// normalizes to ForceNew, so the result speaks the same dialect as the SDKv2
// parser and can feed the same comparison code.
//
// An unparsable source yields an empty map. Callers that must not treat "no
// fields" as "nothing to check" — the consistency gate does not — are expected
// to reject an empty result themselves.
func Parse(source string) map[string]map[string]interface{} {
	fields := make(map[string]map[string]interface{})

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, 0)
	if err != nil {
		return fields
	}

	e := &extractor{fields: fields, decls: packageDecls(file)}
	ast.PreorderStack(file, nil, func(n ast.Node, stack []ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		if isFieldMap(lit, "Attribute") || isFieldMap(lit, "Block") {
			e.extractFields(lit, pathPrefix(stack))
		}
		return true
	})

	return fields
}

// extractor turns schema declarations into field entries. The package-level
// declaration index is what lets a nested object type named by a variable be
// expanded where it is used.
type extractor struct {
	fields map[string]map[string]interface{}
	decls  map[string]ast.Expr
}

// packageDecls indexes package-level var and const declarations by name.
// Framework schemas routinely factor the field list of an object into a named
// variable — either the schema map itself or the AttrTypes of an ObjectType —
// so the declaration alone does not say where those fields belong.
func packageDecls(file *ast.File) map[string]ast.Expr {
	decls := make(map[string]ast.Expr)
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || (genDecl.Tok != token.VAR && genDecl.Tok != token.CONST) {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Values) != len(valueSpec.Names) {
				continue
			}
			for i, name := range valueSpec.Names {
				decls[name.Name] = valueSpec.Values[i]
			}
		}
	}
	return decls
}

// maxObjectDepth bounds the descent through nested object declarations. A
// provider schema nests a handful of levels at most; the bound only keeps a
// pathological source from recursing without end.
const maxObjectDepth = 8

// Framework schema packages can use different import aliases, so the map value
// type is matched by its selector name only. Both dialects are matched this
// way: map[string]schema.Attribute and map[string]schema.Block name a field
// list, and map[string]attr.Type names the same field list from the types side.
func isFieldMap(lit *ast.CompositeLit, fieldType string) bool {
	mapType, ok := lit.Type.(*ast.MapType)
	if !ok {
		return false
	}
	selExpr, ok := mapType.Value.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	return selExpr.Sel.Name == fieldType
}

// pathPrefix walks the enclosing nodes of a field map and joins the attribute
// or block names that contain it, so a nested "acl_enable" under
// "instance_endpoints" is keyed as "instance_endpoints.acl_enable" — the same
// dotted form the SDKv2 flattener and the docs parser produce.
func pathPrefix(stack []ast.Node) string {
	var prefix string
	for i := 1; i < len(stack); i++ {
		field, ok := stack[i].(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		container, ok := stack[i-1].(*ast.CompositeLit)
		if !ok || !(isFieldMap(container, "Attribute") || isFieldMap(container, "Block")) {
			continue
		}
		if key, ok := field.Key.(*ast.BasicLit); ok {
			if name := strings.Trim(key.Value, `"`); name != "" {
				prefix += name + "."
			}
		}
	}
	return prefix
}

func (e *extractor) extractFields(lit *ast.CompositeLit, prefix string) {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.BasicLit)
		if !ok {
			continue
		}
		name := strings.Trim(key.Value, `"`)
		if name == "" {
			continue
		}

		props := map[string]interface{}{"Name": name}
		if body, ok := kv.Value.(*ast.CompositeLit); ok {
			extractProperties(body, props)
		}
		e.fields[prefix+name] = props
		// The field's own value carries its nested object: a nested attribute's
		// Attributes map, or a List/Set/Map attribute whose ElementType is an
		// object. Those are read here rather than by the walk in Parse because
		// the object is often named by a variable, which the walk cannot place
		// under this field.
		e.extractObject(kv.Value, prefix+name+".", 0)
	}
}

// extractObject records the fields of every object declared inside expr under
// prefix, following package-level variables so that a named object type or
// field list expands where it is used.
func (e *extractor) extractObject(expr ast.Expr, prefix string, depth int) {
	if depth > maxObjectDepth {
		return
	}
	switch node := expr.(type) {
	case *ast.Ident:
		if value, ok := e.decls[node.Name]; ok {
			e.extractObject(value, prefix, depth+1)
		}
		return
	case *ast.ParenExpr:
		e.extractObject(node.X, prefix, depth+1)
		return
	case *ast.UnaryExpr:
		e.extractObject(node.X, prefix, depth+1)
		return
	}

	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return
	}
	// A schema map is an object's field list in the schema dialect.
	if isFieldMap(lit, "Attribute") || isFieldMap(lit, "Block") {
		e.extractFields(lit, prefix)
		return
	}
	// map[string]attr.Type is the same field list in the types dialect.
	if isFieldMap(lit, "Type") {
		e.extractTypeFields(lit, prefix)
		return
	}
	// Otherwise this is an attribute declaration or an object type wrapping
	// one: its elements may still carry the field list (Attributes, AttrTypes,
	// ElementType, NestedObject).
	for _, elt := range lit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			e.extractObject(kv.Value, prefix, depth+1)
		}
	}
}

// extractTypeFields records the fields of a map[string]attr.Type, the framework
// spelling of an object's field list. Its values are types rather than
// declarations, so a field known only from here has a type and nothing else —
// no Optional, Required, ForceNew or Deprecated, which only a schema attribute
// declares.
func (e *extractor) extractTypeFields(lit *ast.CompositeLit, prefix string) {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.BasicLit)
		if !ok {
			continue
		}
		name := strings.Trim(key.Value, `"`)
		if name == "" {
			continue
		}
		props := map[string]interface{}{"Name": name}
		// A field known only from an attr.Type list is marked: its field list is
		// commonly a package-level variable, and a caller that compares two
		// single-file revisions cannot always resolve that variable. The marker
		// lets such a caller tell "the object no longer resolves" apart from
		// "the field was removed".
		props["FromTypeList"] = true
		if sel, ok := kv.Value.(*ast.SelectorExpr); ok {
			if equivalent, ok := typeEquivalents[sel.Sel.Name]; ok {
				props["Type"] = equivalent
			} else {
				props["Type"] = sel.Sel.Name
			}
		}
		e.fields[prefix+name] = props
		e.extractObject(kv.Value, prefix+name+".", 0)
	}
}

// typeEquivalents maps framework attribute and block type names to the SDKv2
// schema.TypeXxx constants they play the role of, so a framework entity and an
// SDKv2 one describe the same field identically.
var typeEquivalents = map[string]string{
	"StringAttribute":       "TypeString",
	"BoolAttribute":         "TypeBool",
	"Int64Attribute":        "TypeInt",
	"Float64Attribute":      "TypeFloat",
	"NumberAttribute":       "TypeFloat",
	"ListAttribute":         "TypeList",
	"SetAttribute":          "TypeSet",
	"MapAttribute":          "TypeMap",
	"SingleNestedAttribute": "TypeMap",
	"ObjectNestedAttribute": "TypeMap",
	"ListNestedAttribute":   "TypeList",
	"SetNestedAttribute":    "TypeSet",
	"ListNestedBlock":       "TypeList",
	"SetNestedBlock":        "TypeSet",
	"SingleNestedBlock":     "TypeMap",
	// The types dialect spells the same roles as attr.Type values inside a
	// map[string]attr.Type. Normalizing them keeps one vocabulary across both
	// spellings, so a nested field reports the same type whichever way its
	// object is declared.
	"StringType":  "TypeString",
	"BoolType":    "TypeBool",
	"Int64Type":   "TypeInt",
	"Float64Type": "TypeFloat",
	"NumberType":  "TypeFloat",
	"ListType":    "TypeList",
	"SetType":     "TypeSet",
	"MapType":     "TypeMap",
	"ObjectType":  "TypeMap",
}

func extractProperties(lit *ast.CompositeLit, props map[string]interface{}) {
	if selExpr, ok := lit.Type.(*ast.SelectorExpr); ok {
		if equivalent, ok := typeEquivalents[selExpr.Sel.Name]; ok {
			props["Type"] = equivalent
		} else {
			props["Type"] = selExpr.Sel.Name
		}
	}
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch ident.Name {
		case "Optional", "Required", "Computed", "Sensitive":
			if value, ok := kv.Value.(*ast.Ident); ok {
				props[ident.Name] = value.Name == "true"
			}
		case "Deprecated":
			// The framework spells deprecation as a free-form message.
			if msg, ok := stringValue(kv.Value); ok {
				props["Deprecated"] = msg
			}
		case "PlanModifiers":
			// RequiresReplace is the framework counterpart of ForceNew.
			if requiresReplace(kv.Value) {
				props["ForceNew"] = true
			}
		}
	}
}

// requiresReplace reports whether a PlanModifiers list contains a
// RequiresReplace plan modifier.
func requiresReplace(expr ast.Expr) bool {
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return false
	}
	for _, elt := range compLit.Elts {
		switch e := elt.(type) {
		case *ast.CallExpr:
			if selExpr, ok := e.Fun.(*ast.SelectorExpr); ok && selExpr.Sel.Name == "RequiresReplace" {
				return true
			}
		case *ast.Ident:
			if e.Name == "RequiresReplace" {
				return true
			}
		case *ast.SelectorExpr:
			if e.Sel.Name == "RequiresReplace" {
				return true
			}
		}
	}
	return false
}

// stringValue resolves a string literal or a simple "+"-joined literal
// expression, which is how multi-line deprecation messages are written.
func stringValue(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", false
		}
		return strings.Trim(e.Value, `"`), true
	case *ast.BinaryExpr:
		left, ok := stringValue(e.X)
		if !ok {
			return "", false
		}
		right, ok := stringValue(e.Y)
		if !ok {
			return "", false
		}
		return left + right, true
	}
	return "", false
}
