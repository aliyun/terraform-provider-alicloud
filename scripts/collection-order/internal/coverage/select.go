package coverage

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
)

// frameworkRegistered reports whether a type name is registered
// framework-native in any service package. It defaults to
// splitlayout.FrameworkRegistered but is a package-level variable so tests
// can substitute a stub without touching the filesystem.
var frameworkRegistered = splitlayout.FrameworkRegistered

// frameworkFindByTypeName reverse-resolves a framework type name to its
// split-layout paths. It defaults to splitlayout.FrameworkFindByTypeName
// but is a package-level variable so tests can substitute a stub.
var frameworkFindByTypeName = splitlayout.FrameworkFindByTypeName

// frameworkOrderedCollections reports whether a Framework schema source
// declares any configurable ordered collection: ListAttribute,
// ListNestedAttribute (Attributes map) or ListNestedBlock (Blocks map) from
// terraform-plugin-framework/resource/schema. SetAttribute and
// SetNestedAttribute are unordered by definition and excluded, but their
// NestedObject may itself declare ordered collections and is descended. The
// function inspects the Go AST of the implementation file; it does not
// execute the schema, so it is safe on any source.
//
// It recurses into nested schema holders — SingleNestedAttribute,
// NestedAttributeObject, NestedBlockObject, and the NestedObject field of
// ListNestedAttribute/SetNestedAttribute/MapNestedAttribute/ListNestedBlock/
// SetNestedBlock/SingleNestedBlock — so an ordered collection or an
// unresolved cross-file helper buried one or more levels deep is detected,
// not just the top-level Attributes/Blocks map. A top-level scan alone would
// classify schema.Schema{Attributes: map{"x": schema.SingleNestedAttribute{
// Attributes: helperFromAnotherFile()}}} as having no ordered collections,
// which is unsound: the helper may return ListAttribute/ListNestedAttribute.
//
// It is fail-closed: a parse error, a missing schema import, a schema with no
// schema.Schema composite literal in this file, or any Attributes/Blocks/
// NestedObject field (at any nesting depth) built from an unresolved call or
// variable returns an error rather than treating the source as having no
// ordered collections — "cannot prove" never means "proven absent". A resource
// whose schema is defined in this file and has no ordered collections has no
// configurable order, so the SDKv2 order-coverage checker does not apply —
// proven by the schema, not a blanket exemption. A resource that does have
// them cannot be verified by the SDKv2 checker (which requires
// *schema.Resource); Select returns a clear error so the developer must
// decide rather than the resource passing silently.
func frameworkOrderedCollections(src []byte) (bool, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "resource.go", src, parser.SkipObjectResolution)
	if err != nil {
		return false, fmt.Errorf("cannot parse framework schema source: %w", err)
	}
	pkg := frameworkSchemaImport(f)
	if pkg == "" {
		return false, fmt.Errorf("framework schema source has no terraform-plugin-framework/resource/schema import; cannot determine ordered collections")
	}
	orderedSelectors := map[string]bool{
		"ListAttribute":       true,
		"ListNestedAttribute": true,
		"ListNestedBlock":     true,
	}
	hasOrdered := false
	unresolved := false
	schemaDefined := false

	// scanSchemaLit walks a schema composite literal and recurses into its
	// Attributes/Blocks map values and NestedObject field. It descends through
	// SingleNestedAttribute, NestedAttributeObject and NestedBlockObject so a
	// cross-file helper or an ordered collection nested any depth below the
	// root schema.Schema is detected — not just the top-level map.
	var scanSchemaLit func(lit *ast.CompositeLit)
	scanSchemaLit = func(lit *ast.CompositeLit) {
		if hasOrdered {
			return
		}
		sel, ok := lit.Type.(*ast.SelectorExpr)
		if !ok || ident(sel.X) != pkg {
			return
		}
		if sel.Sel.Name == "Schema" {
			schemaDefined = true
		}
		for _, el := range lit.Elts {
			if hasOrdered {
				return
			}
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			switch ident(kv.Key) {
			case "Attributes", "Blocks":
				// The field value must be an inline map literal. A call or
				// variable (cross-file helper) cannot be proven free of ordered
				// collections: fail closed.
				inner, ok := kv.Value.(*ast.CompositeLit)
				if !ok {
					unresolved = true
					continue
				}
				for _, item := range inner.Elts {
					if hasOrdered {
						return
					}
					entry, ok := item.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					attrLit, ok := entry.Value.(*ast.CompositeLit)
					if !ok {
						// An attribute built from a call/variable is unresolved.
						unresolved = true
						continue
					}
					attrSel, ok := attrLit.Type.(*ast.SelectorExpr)
					if !ok || ident(attrSel.X) != pkg {
						unresolved = true
						continue
					}
					if orderedSelectors[attrSel.Sel.Name] {
						hasOrdered = true
						return
					}
					// Recurse into any nested schema holder (e.g.
					// SingleNestedAttribute, NestedAttributeObject,
					// NestedBlockObject) whose own Attributes/Blocks/NestedObject
					// may declare further collections or unresolved helpers.
					scanSchemaLit(attrLit)
				}
			case "NestedObject":
				// The NestedObject field of ListNestedAttribute,
				// SetNestedAttribute, MapNestedAttribute, ListNestedBlock,
				// SetNestedBlock and SingleNestedBlock points at another schema
				// literal (NestedAttributeObject/NestedBlockObject). An
				// unresolved value is a cross-file helper; otherwise descend.
				nestedLit, ok := kv.Value.(*ast.CompositeLit)
				if !ok {
					unresolved = true
					continue
				}
				scanSchemaLit(nestedLit)
			}
		}
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if hasOrdered {
			return false
		}
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		sel, ok := lit.Type.(*ast.SelectorExpr)
		if !ok || ident(sel.X) != pkg {
			return true
		}
		if sel.Sel.Name == "Schema" {
			scanSchemaLit(lit)
			return false // scanSchemaLit descended recursively into children
		}
		// An ordered selector appearing outside a schema.Schema literal (e.g.
		// assigned to a package-level var) still means the resource has ordered
		// collections.
		if orderedSelectors[sel.Sel.Name] {
			hasOrdered = true
			return false
		}
		return true
	})

	if hasOrdered {
		return true, nil
	}
	// A schema provably free of configurable ordered collections must be built
	// entirely from inline composite literals. Any post-construction write to a
	// schema's Attributes or Blocks map — an index-map assignment
	// obj.Attributes[key] = value / obj.Blocks[key] = value, or a field
	// replacement obj.Attributes = value / obj.Blocks = value — adds schema
	// elements that the inline scan above does not see, so the absence of
	// ListAttribute/ListNestedAttribute/ListNestedBlock cannot be proven. This
	// is bounded static analysis: it pattern-matches assignment forms and
	// rejects every unsupported post-construction construction by name; it does
	// not interpret Go data flow, resolve cross-file helpers, or chase
	// receivers, pointers, aliases or closures individually. "Cannot prove"
	// returns an error rather than (false, nil), so a schema built as
	//
	//	s := schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
	//	s.Attributes["ordered"] = orderedField() // returns ListAttribute elsewhere
	//
	// is flagged instead of passing silently as ordered-free.
	if desc := postConstructionSchemaMutation(f); desc != "" {
		return false, fmt.Errorf("framework schema source declares schema.Schema but mutates Attributes/Blocks after construction via %s; the SDKv2 order-coverage checker can only verify schemas built fully inline, cannot determine ordered collections", desc)
	}
	if !schemaDefined {
		return false, fmt.Errorf("framework schema source imports terraform-plugin-framework/resource/schema but declares no schema.Schema composite literal in this file; cannot determine ordered collections (schema may be constructed by a cross-file helper)")
	}
	if unresolved {
		return false, fmt.Errorf("framework schema source declares schema.Schema but an Attributes/Blocks/NestedObject field (possibly nested inside SingleNestedAttribute/NestedAttributeObject/NestedBlockObject) is built from an unresolved call or variable; cannot determine ordered collections (a cross-file helper may introduce ListAttribute/ListNestedAttribute/ListNestedBlock)")
	}
	return false, nil
}

// postConstructionSchemaMutation reports the first statement in the schema
// source that mutates a schema's Attributes or Blocks map after the
// schema.Schema composite literal is constructed, or that escapes such a map
// into a value whose later use the single-file analysis cannot track. It
// returns "" when no such statement exists. See frameworkOrderedCollections
// for the rationale: a schema provably free of ordered collections must be
// built fully inline, so any post-construction write or escape cannot be
// proven and is rejected by name rather than treated as ordered-free. The
// analysis pattern-matches assignment and variable-declaration forms only; it
// does not interpret Go data flow, resolve cross-file helpers, or chase
// receivers, pointers, aliases or closures individually — it rejects every
// unsupported construction uniformly, so a schema built as
//
//	s := schema.Schema{Attributes: map[string]schema.Attribute{"name": schema.StringAttribute{}}}
//	s.Attributes["ordered"] = orderedField() // returns ListAttribute elsewhere
//
// or via an alias
//
//	m := s.Attributes
//	m["ordered"] = orderedField()
//
// is flagged instead of passing silently as ordered-free.
func postConstructionSchemaMutation(f *ast.File) string {
	var first string
	ast.Inspect(f, func(n ast.Node) bool {
		if first != "" {
			return false
		}
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range stmt.Lhs {
				if desc := schemaCollectionWrite(lhs); desc != "" {
					first = desc
					break
				}
			}
			if first == "" {
				if desc := schemaCollectionEscape(stmt.Rhs); desc != "" {
					first = desc
				}
			}
		case *ast.GenDecl:
			for _, spec := range stmt.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				if desc := schemaCollectionEscape(vs.Values); desc != "" {
					first = desc
					break
				}
			}
		}
		return first == ""
	})
	return first
}

// schemaCollectionWrite describes a post-construction write to a schema
// Attributes or Blocks map on the given left-hand side, or "" when the LHS
// is not such a write. Recognized forms are index-map assignment
// obj.Attributes[key] = value / obj.Blocks[key] = value and field
// replacement obj.Attributes = value / obj.Blocks = value. The base object
// is not constrained to a tracked schema variable: a receiver, pointer,
// alias or closure-captured value writing to Attributes/Blocks is equally a
// post-construction mutation that cannot be proven, so every such write is
// rejected uniformly instead of enumerating escape forms one by one.
func schemaCollectionWrite(lhs ast.Expr) string {
	switch e := lhs.(type) {
	case *ast.IndexExpr:
		sel, ok := e.X.(*ast.SelectorExpr)
		if !ok {
			return ""
		}
		if sel.Sel.Name == "Attributes" || sel.Sel.Name == "Blocks" {
			return fmt.Sprintf("index-map assignment to %s (obj.%s[key] = value); define attributes inline in the schema.Schema composite literal", sel.Sel.Name, sel.Sel.Name)
		}
	case *ast.SelectorExpr:
		if e.Sel.Name == "Attributes" || e.Sel.Name == "Blocks" {
			return fmt.Sprintf("field replacement of %s (obj.%s = value); define attributes inline in the schema.Schema composite literal", e.Sel.Name, e.Sel.Name)
		}
	}
	return ""
}

// schemaCollectionEscape describes a right-hand side that escapes a schema
// Attributes or Blocks map into a value the single-file analysis cannot
// track, or "" when the RHS is not such an escape. Recognized forms are
// aliasing m := obj.Attributes / m := obj.Blocks and pointer aliasing
// p := &obj.Attributes / p := &obj.Blocks. An escaped map may be mutated
// later through the alias or pointer, so the schema can no longer be proven
// ordered-free; the escape is rejected rather than tracking every later use.
func schemaCollectionEscape(rhs []ast.Expr) string {
	for _, expr := range rhs {
		switch e := expr.(type) {
		case *ast.SelectorExpr:
			if e.Sel.Name == "Attributes" || e.Sel.Name == "Blocks" {
				return fmt.Sprintf("alias of %s (name = obj.%s); define attributes inline in the schema.Schema composite literal", e.Sel.Name, e.Sel.Name)
			}
		case *ast.UnaryExpr:
			sel, ok := e.X.(*ast.SelectorExpr)
			if ok && (sel.Sel.Name == "Attributes" || sel.Sel.Name == "Blocks") {
				return fmt.Sprintf("pointer alias of %s (name = &obj.%s); define attributes inline in the schema.Schema composite literal", sel.Sel.Name, sel.Sel.Name)
			}
		}
	}
	return ""
}

// frameworkSchemaImport returns the identifier used in src for the
// terraform-plugin-framework/resource/schema package, or "" when absent.
func frameworkSchemaImport(f *ast.File) string {
	for _, imp := range f.Imports {
		path, _ := strconv.Unquote(imp.Path.Value)
		if path != "github.com/hashicorp/terraform-plugin-framework/resource/schema" {
			continue
		}
		if imp.Name != nil {
			return imp.Name.Name
		}
		return "schema"
	}
	return ""
}

// Registry reads constructor identities, not an inferred resource name from a filename.
func Registry(src []byte) (map[string]string, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "provider.go", src, 0)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	var parseErr error
	ast.Inspect(f, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok || ident(kv.Key) != "ResourcesMap" {
			return true
		}
		entries := fields(kv.Value)
		if entries == nil {
			parseErr = fmt.Errorf("ResourcesMap must be a literal map")
			return false
		}
		for name, expr := range entries {
			call, ok := expr.(*ast.CallExpr)
			if !ok || ident(call.Fun) == "" || len(call.Args) != 0 {
				parseErr = fmt.Errorf("unresolved resource constructor for %s", name)
				continue
			}
			result[name] = ident(call.Fun)
		}
		return false
	})
	if parseErr != nil {
		return nil, parseErr
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("ResourcesMap is missing or empty")
	}
	return result, nil
}

// RegistryUnion merges the provider.go ResourcesMap literal with the SDKv2
// registrations declared by the split service packages. Split constructors
// carry a product qualifier ("cloudauth.ResourceFaceConfig") because
// top-level names are only unique per package; the flat literal keeps bare
// names. A declaration appearing in both sources is the same startup error
// assembleProvider reports.
func RegistryUnion(providerSrc []byte, registers map[string][]byte) (map[string]string, error) {
	reg, err := Registry(providerSrc)
	if err != nil {
		return nil, err
	}
	for path, src := range registers {
		product, ok := splitlayout.IsServiceRegister(path)
		if !ok {
			continue
		}
		for _, entry := range splitlayout.SDKRegisterEntries(src) {
			if entry.Kind != splitlayout.KindResource {
				continue
			}
			qualified := product + "." + entry.Factory
			if previous, dup := reg[entry.TypeName]; dup {
				return nil, fmt.Errorf("duplicate resource declaration for %s: %s and %s", entry.TypeName, previous, qualified)
			}
			reg[entry.TypeName] = qualified
		}
	}
	return reg, nil
}

// qualify keys a file's top-level names the way the registry references
// them: bare in the flat alicloud package, product-qualified in the split
// service packages, where names are only unique per package.
func qualify(path string, names map[string]bool) map[string]bool {
	product, ok := splitlayout.ServiceProduct(path)
	if !ok {
		return names
	}
	out := make(map[string]bool, len(names))
	for name := range names {
		out[product+"."+name] = true
	}
	return out
}

// Select maps both changed implementations and tests to current registrations. A deleted
// implementation is ignored only when its constructors are no longer registered.
func Select(changed []string, before, after map[string]string, sources map[string][]byte, oldSource func(string) ([]byte, error)) (map[string]string, error) {
	functions := map[string]string{}
	fileFunctions := map[string]map[string]bool{}
	for path, src := range sources {
		names, err := functionNames(path, src)
		if err != nil {
			return nil, err
		}
		names = qualify(path, names)
		fileFunctions[path] = names
		for name := range names {
			if previous, ok := functions[name]; ok {
				return nil, fmt.Errorf("duplicate function %s in %s and %s", name, previous, path)
			}
			functions[name] = path
		}
	}
	selected := map[string]string{}
	add := func(name string) error {
		path := functions[after[name]]
		if path == "" {
			return fmt.Errorf("cannot locate implementation for registered resource %s (%s)", name, after[name])
		}
		selected[name] = strings.TrimSuffix(path, ".go") + "_test.go"
		return nil
	}
	for _, path := range changed {
		if path == "alicloud/provider.go" {
			for name, constructor := range after {
				if before[name] != constructor {
					if err := add(name); err != nil {
						return nil, err
					}
				}
			}
			continue
		}
		// A product register.go is a registration source the same way
		// provider.go is: select every registration the diff touched. The
		// union registry keeps split constructors product-qualified, so
		// untouched products contribute nothing.
		if _, ok := splitlayout.IsServiceRegister(path); ok {
			for name, constructor := range after {
				if before[name] != constructor {
					if err := add(name); err != nil {
						return nil, err
					}
				}
			}
			continue
		}
		flat := strings.HasPrefix(path, "alicloud/resource_alicloud_") && strings.HasSuffix(path, ".go")
		_, stem, _, isSplit := splitlayout.IsSplitSource(path)
		isSplit = isSplit && strings.HasPrefix(stem, "resource_")
		if !flat && !isSplit {
			// Data sources, docs and framework-only service files stay
			// outside the managed-resource scope, as they were in the flat
			// layout.
			continue
		}
		implementation := strings.TrimSuffix(path, "_test.go")
		if implementation != path {
			implementation += ".go"
		}
		names, exists := fileFunctions[implementation]
		if !exists {
			// A flat-layout test file whose implementation moved to the
			// service/ layout resolves to a path that no longer exists here.
			// If the resource is framework-native, resolve its service-layout
			// implementation and inspect the Framework schema for configurable
			// ordered collections. This is fail-closed at every step: a
			// missing path, missing source, parse error or missing schema
			// import returns a clear error rather than silently skipping —
			// "cannot prove" is never treated as "proven absent".
			if flat {
				typeName := "alicloud_" + strings.TrimSuffix(strings.TrimPrefix(implementation, "alicloud/resource_alicloud_"), ".go")
				if frameworkRegistered(typeName) {
					resolved := frameworkFindByTypeName(typeName)
					if resolved == nil {
						return nil, fmt.Errorf("framework-native resource %s is registered but its implementation path cannot be resolved; the SDKv2 order-coverage checker cannot verify Framework schemas", typeName)
					}
					implSrc, ok := sources[resolved.ImplPath]
					if !ok {
						return nil, fmt.Errorf("framework-native resource %s: implementation source %s is not available; cannot determine ordered collections", typeName, resolved.ImplPath)
					}
					hasOrdered, err := frameworkOrderedCollections(implSrc)
					if err != nil {
						return nil, fmt.Errorf("framework-native resource %s: %w", typeName, err)
					}
					if !hasOrdered {
						continue
					}
					return nil, fmt.Errorf("framework-native resource %s has configurable ordered collections (ListAttribute/ListNestedAttribute/ListNestedBlock); the SDKv2 order-coverage checker cannot verify Framework schemas; see scripts/collection-order/README.md", typeName)
				}
			}
			src, err := oldSource(implementation)
			if err != nil {
				return nil, fmt.Errorf("cannot resolve changed resource %s: %w", path, err)
			}
			base, err := functionNames(implementation, src)
			if err != nil {
				return nil, err
			}
			names = qualify(implementation, base)
		}
		matched := false
		for name, constructor := range after {
			if names[constructor] {
				matched = true
				if err := add(name); err != nil {
					return nil, err
				}
			}
		}
		if !matched && exists && flat {
			return nil, fmt.Errorf("changed resource file %s has no current ResourcesMap registration", path)
		}
		if !matched && !exists && flat {
			removed := false
			for _, constructor := range before {
				removed = removed || names[constructor]
			}
			if !removed {
				return nil, fmt.Errorf("cannot establish deletion of registered resource %s", path)
			}
		}
		// A changed split file with no registration is a framework resource,
		// data source or helper outside this checker: the product's
		// register.go is the source of truth for the SDKv2 scope.
	}
	return selected, nil
}

func functionNames(filename string, src []byte) (map[string]bool, error) {
	f, err := parser.ParseFile(token.NewFileSet(), filename, src, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	names := map[string]bool{}
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil {
			names[fn.Name.Name] = true
		}
	}
	return names, nil
}

// Names gives deterministic diagnostics and command output.
func Names(selected map[string]string) []string {
	names := make([]string, 0, len(selected))
	for name := range selected {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
