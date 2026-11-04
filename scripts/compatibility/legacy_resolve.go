//nolint:all
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
)

// legacyKind maps a flat-layout stem to the registration kind it implies. An
// empty stem names an entity born in the split layout (ephemeral, action or
// list resource) that has no flat-layout predecessor to resolve.
func legacyKind(legacyStem string) (splitlayout.Kind, bool) {
	switch legacyStem {
	case "resource_":
		return splitlayout.KindResource, true
	case "data_source_":
		return splitlayout.KindDataSource, true
	}
	return "", false
}

// providerFactories parses the root alicloud/provider.go ResourcesMap and
// DataSourcesMap literals and returns one entry per registered public type,
// carrying the bare constructor identity (the map value's call-expression
// function name) and the kind its map literal belongs to.
// splitlayout.ProviderRegisteredNames leaves the Factory blank because the
// title gate only needs the type name and kind; the breaking-change resolver
// must instead follow the registered constructor to its declaring file, so it
// fills the Factory from the call expression. A value that is not a
// zero-argument call on a plain identifier — a helper, a variable, or a call
// with arguments — cannot be followed to a declaring file by name and is
// omitted; the service registers then carry the registration for a split type.
func providerFactories(providerSrc []byte) []splitlayout.Entry {
	f, err := parser.ParseFile(token.NewFileSet(), "provider.go", providerSrc, parser.SkipObjectResolution)
	if err != nil {
		return nil
	}
	var entries []splitlayout.Entry
	ast.Inspect(f, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "DataSourcesMap" && key.Name != "ResourcesMap" {
			return true
		}
		kind := splitlayout.KindResource
		if key.Name == "DataSourcesMap" {
			kind = splitlayout.KindDataSource
		}
		lit, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			return true
		}
		for _, elt := range lit.Elts {
			pair, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			name := stringLit(pair.Key)
			if name == "" {
				continue
			}
			call, ok := pair.Value.(*ast.CallExpr)
			if !ok || len(call.Args) != 0 {
				continue
			}
			factory := funcCallName(call.Fun)
			if factory == "" {
				continue
			}
			entries = append(entries, splitlayout.Entry{TypeName: name, Factory: factory, Kind: kind})
		}
		return true
	})
	return entries
}

// registerFactories collects the SDKv2 and framework registrations declared by
// a set of service/core register.go sources. Constructors carry bare
// identifiers in the register source; a split package qualifies them at
// lookup time, but the declaring file is found by the bare name.
func registerFactories(registers map[string][]byte) []splitlayout.Entry {
	var entries []splitlayout.Entry
	for _, src := range registers {
		entries = append(entries, splitlayout.SDKRegisterEntries(src)...)
		entries = append(entries, splitlayout.FrameworkRegisterEntries(src)...)
	}
	return entries
}

// resolveFactory finds the bare constructor registered for typeName under kind
// at the baseline, searching the provider root maps then the service/core
// registers. It returns the constructor identity, whether the type was
// registered, and an error only when the same type+kind resolves to more than
// one constructor — an ambiguous registration the checker must not guess from.
func resolveFactory(typeName string, kind splitlayout.Kind, providerSrc []byte, registers map[string][]byte) (string, bool, error) {
	var found string
	for _, e := range providerFactories(providerSrc) {
		if e.TypeName != typeName || e.Kind != kind || e.Factory == "" {
			continue
		}
		if found != "" && found != e.Factory {
			return "", false, fmt.Errorf("ambiguous factory for %s %s: %s and %s", kind, typeName, found, e.Factory)
		}
		found = e.Factory
	}
	for _, e := range registerFactories(registers) {
		if e.TypeName != typeName || e.Kind != kind || e.Factory == "" {
			continue
		}
		if found != "" && found != e.Factory {
			return "", false, fmt.Errorf("ambiguous factory for %s %s: %s and %s", kind, typeName, found, e.Factory)
		}
		found = e.Factory
	}
	return found, found != "", nil
}

// findDeclaringFile returns the single file whose top-level declarations define
// func factory. It AST-verifies the declaration so a comment, string literal,
// call site or method receiver cannot satisfy a text match. Zero matches mean
// a registered constructor with no implementation at the baseline; more than
// one means the constructor is declared in multiple files and cannot be
// uniquely attributed. Either case fails closed rather than guess a file.
func findDeclaringFile(factory string, impls map[string][]byte) (string, error) {
	var matches []string
	for path, src := range impls {
		if declaresFunc(src, factory) {
			matches = append(matches, path)
		}
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("registered factory %s has no declaring implementation file", factory)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("factory %s declared in multiple files: %s", factory, strings.Join(matches, ", "))
	}
}

// declaresFunc reports whether src declares a top-level function named factory
// (no receiver). It guards findDeclaringFile against a grep matching a
// comment, string, method, or a similarly-named function in an unrelated file.
func declaresFunc(src []byte, factory string) bool {
	f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.SkipObjectResolution)
	if err != nil {
		return false
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil {
			continue
		}
		if fn.Name.Name == factory {
			return true
		}
	}
	return false
}

// funcCallName returns the identifier a call expression is invoked on, or the
// empty string when the function is not a plain identifier (a selector or
// parenthesised expression cannot name a flat-layout constructor).
func funcCallName(expr ast.Expr) string {
	if i, ok := expr.(*ast.Ident); ok {
		return i.Name
	}
	return ""
}

// stringLit unquotes a string literal, returning the empty string for any other
// expression. It mirrors splitlayout.stringLit so this file stays free of that
// unexported helper.
func stringLit(expr ast.Expr) string {
	if s, ok := expr.(*ast.BasicLit); ok && s.Kind == token.STRING {
		if v, err := strconv.Unquote(s.Value); err == nil {
			return v
		}
	}
	return ""
}

// resolveLegacyImplGit resolves the flat-layout file that declared the
// constructor registered for typeName under kind at oldRevision. It reads
// provider.go and the service/core registers at oldRevision, resolves the
// registered constructor, then locates its unique declaring file through git
// grep and AST verification. It returns the declaring file path and content;
// an empty path with no error means the type was not registered at oldRevision
// (a genuinely new entity). Ambiguous, registered-but-undeclared or
// unresolvable git state fail closed rather than fall back to a filename guess.
func resolveLegacyImplGit(typeName string, kind splitlayout.Kind, oldRevision string) (string, string, error) {
	providerSrc, err := getFileContentFromGit("alicloud/provider.go", oldRevision)
	if err != nil {
		return "", "", fmt.Errorf("cannot read baseline provider.go at %s: %w", oldRevision, err)
	}
	registers, err := registersAt(oldRevision)
	if err != nil {
		return "", "", err
	}
	factory, registered, err := resolveFactory(typeName, kind, []byte(providerSrc), registers)
	if err != nil {
		return "", "", err
	}
	if !registered {
		// resolveFactory follows only constructors expressible as a
		// zero-argument call on a plain identifier (providerFactories omits a
		// helper, a variable, or a call with arguments). A type registered at
		// the baseline with such a value is nonetheless registered by name —
		// ProviderRegisteredNames records every root ResourcesMap/
		// DataSourcesMap entry by type name and kind regardless of the value
		// expression — so it cannot be waived as a genuinely-new entity and
		// exempted from comparison. Fail closed rather than guess a factory
		// expression or declaring file; do not widen the resolver to support
		// arbitrary expressions. A type not registered by name at the baseline
		// is a genuinely new entity.
		if registeredByName([]byte(providerSrc), typeName, kind) {
			return "", "", fmt.Errorf("%s %s is registered at baseline %s but its constructor is not a resolvable zero-argument factory call", kind, typeName, oldRevision)
		}
		return "", "", nil
	}
	candidates, err := grepDeclaringFiles(factory, oldRevision)
	if err != nil {
		return "", "", err
	}
	impls := map[string][]byte{}
	for _, path := range candidates {
		src, gerr := getFileContentFromGit(path, oldRevision)
		if gerr != nil {
			return "", "", fmt.Errorf("cannot read candidate %s at %s: %w", path, oldRevision, gerr)
		}
		impls[path] = []byte(src)
	}
	path, err := findDeclaringFile(factory, impls)
	if err != nil {
		return "", "", err
	}
	return path, string(impls[path]), nil
}

// registersAt reads the supported service/core register.go sources at
// oldRevision. Products that did not exist at oldRevision simply contribute no
// registrations there, the same way the collection-order gate reads the base
// registers. A register.go absent at oldRevision (a product not yet present at
// the baseline) is distinguished from a real git read error: absence is decided
// by whether the path is present in the baseline tree (git ls-tree), not by
// git cat-file -e, which exits 1 for a missing blob object exactly as it does
// for an absent path. A register whose tree entry still exists but whose blob
// object is unreadable is a broken git state, not an expected absence, and
// fails closed so it cannot be silently swallowed into an empty register set.
func registersAt(oldRevision string) (map[string][]byte, error) {
	paths, err := splitlayout.RegisterPaths()
	if err != nil {
		return nil, err
	}
	registers := map[string][]byte{}
	for _, path := range paths {
		exists, err := pathExistsInTree(path, oldRevision)
		if err != nil {
			// The tree query itself failed (bad revision, corrupt tree): a
			// real git state error, not an expected absence. Fail closed.
			return nil, fmt.Errorf("cannot determine register %s at %s: %w", path, oldRevision, err)
		}
		if !exists {
			// The register.go is not in the baseline tree: its product is not
			// present at the baseline and contributes no registrations,
			// matching the collection-order gate. Skip without error.
			continue
		}
		src, err := getFileContentFromGit(path, oldRevision)
		if err != nil {
			// The file is in the tree but git could not read its blob (e.g. a
			// missing object): a broken git state, not an expected absence.
			// Fail closed rather than silently swallow it into an empty
			// register set, which could hide a real registration and let its
			// type be waived as a genuinely-new entity.
			return nil, fmt.Errorf("cannot read register %s at %s: %w", path, oldRevision, err)
		}
		registers[path] = []byte(src)
	}
	return registers, nil
}

// pathExistsInTree reports whether path is present in the tree of revision,
// distinguishing a genuinely absent file from a git read or object error. It
// inspects the tree entry (git ls-tree) rather than the blob object: ls-tree
// reads the tree object and lists a path whose blob object is missing exactly
// like one whose blob is present, whereas git cat-file -e exits 1 for a missing
// blob and would misclassify broken git state as an expected absence. A tree
// query error (bad revision, corrupt tree) propagates; only an empty ls-tree
// output proves the path is not in the tree.
func pathExistsInTree(path, revision string) (bool, error) {
	cmd := exec.Command("git", "ls-tree", "-z", revision, "--", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git ls-tree %s at %s: %w: %s", path, revision, err, output)
	}
	return strings.TrimSpace(string(output)) != "", nil
}

// registeredByName reports whether typeName is registered under kind in the
// root alicloud/provider.go source, by name only and regardless of the map
// value's expression form. splitlayout.ProviderRegisteredNames records every
// root ResourcesMap/DataSourcesMap entry by type name and kind (its Factory
// field is blank); providerFactories only follows values that are
// zero-argument calls on a plain identifier. A type registered with any other
// value (a helper, a variable, or a call with arguments) is still registered
// by name and must not be waived as a genuinely-new entity merely because its
// constructor cannot be followed to a declaring file.
func registeredByName(providerSrc []byte, typeName string, kind splitlayout.Kind) bool {
	for _, e := range splitlayout.ProviderRegisteredNames(providerSrc) {
		if e.TypeName == typeName && e.Kind == kind {
			return true
		}
	}
	return false
}

// grepDeclaringFiles lists the alicloud/ files at oldRevision whose text
// contains "func <factory>", narrowing the candidate set before AST
// verification. The "func " prefix excludes call sites, comments and method
// receivers (a method has a receiver between func and the name); declaresFunc
// then confirms a real top-level declaration. git grep exits 1 when no file
// matches, which means a registered constructor with no implementation at the
// baseline — fail closed rather than guess a path.
func grepDeclaringFiles(factory, oldRevision string) ([]string, error) {
	cmd := exec.Command("git", "grep", "-l", "-z", "-F", "func "+factory, oldRevision, "--", "alicloud/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// git grep exits 1 when no file matches; that is a registered
		// constructor with no implementation at the baseline — fail closed.
		return nil, fmt.Errorf("git grep %s at %s: %w: %s", factory, oldRevision, err, output)
	}
	// When git grep scans a tree revision, each NUL-delimited result is
	// prefixed "<revision>:<path>" (e.g. "84b723f9:alicloud/resource_x.go").
	// The previous parser treated that whole string as the path, so
	// getFileContentFromGit ran "git show <rev>:<rev>:<path>" and exited 128.
	// Strip the known revision prefix and validate the resulting path; any
	// entry that does not match the expected prefixed form or yields an
	// illegal path fails closed rather than feeding a malformed path downstream.
	prefix := oldRevision + ":"
	var paths []string
	for _, raw := range strings.Split(string(output), "\x00") {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}
		path := strings.TrimPrefix(entry, prefix)
		if path == entry {
			return nil, fmt.Errorf("git grep %s at %s: entry without expected revision prefix %q: %q", factory, oldRevision, prefix, entry)
		}
		if !isLegalDeclaringPath(path) {
			return nil, fmt.Errorf("git grep %s at %s: illegal candidate path: %q", factory, oldRevision, path)
		}
		paths = append(paths, path)
	}
	return paths, nil
}

// isLegalDeclaringPath reports whether path is a plausible alicloud/ Go source
// file that could declare a flat-layout constructor. It rejects absolute paths,
// backslashes, "." / ".." segments and non-.go files so a malformed or injected
// grep entry cannot reach getFileContentFromGit as a bogus path.
func isLegalDeclaringPath(path string) bool {
	if !strings.HasPrefix(path, "alicloud/") {
		return false
	}
	if strings.ContainsRune(path, '\\') {
		return false
	}
	rest := strings.TrimPrefix(path, "alicloud/")
	if rest == "" {
		return false
	}
	for _, seg := range strings.Split(rest, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return false
		}
	}
	return strings.HasSuffix(path, ".go")
}
