// Command sweep_package locates the test binary that owns a literal sweeper
// registration. It only parses source files; it never runs a sweeper.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func sweepPackages(root, resource string) ([]string, error) {
	packages := map[string]bool{}
	err := filepath.WalkDir(filepath.Join(root, "alicloud"), func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return err
		}
		aliases := map[string]bool{}
		for _, im := range f.Imports {
			p, err := strconv.Unquote(im.Path.Value)
			if err != nil {
				return err
			}
			if p != "github.com/hashicorp/terraform-plugin-testing/helper/resource" {
				continue
			}
			name := "resource"
			if im.Name != nil {
				name = im.Name.Name
			}
			aliases[name] = true
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "AddTestSweepers" {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if !ok || !aliases[id.Name] {
				return true
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			name, err := strconv.Unquote(lit.Value)
			if err != nil || name != resource {
				return true
			}
			rel, err := filepath.Rel(root, filepath.Dir(path))
			if err == nil {
				packages["./"+filepath.ToSlash(rel)] = true
			}
			return true
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(packages))
	for p := range packages {
		result = append(result, p)
	}
	sort.Strings(result)
	return result, nil
}

func selectSweepPackage(packages []string, explicit string) (string, error) {
	if explicit != "" {
		for _, p := range packages {
			if p == explicit {
				return p, nil
			}
		}
		return "", fmt.Errorf("SWEEP_PACKAGE %q does not register the requested resource", explicit)
	}
	if len(packages) == 0 {
		return "", fmt.Errorf("no literal sweeper registration found for RESOURCE")
	}
	if len(packages) != 1 {
		return "", fmt.Errorf("multiple sweeper packages %v; select one with SWEEP_PACKAGE", packages)
	}
	return packages[0], nil
}

func main() {
	resource := flag.String("resource", "", "exact alicloud resource type name")
	explicit := flag.String("package", "", "optional explicit owning package")
	flag.Parse()
	if !strings.HasPrefix(*resource, "alicloud_") || strings.ContainsAny(*resource, "/\\ \t\r\n") {
		fmt.Fprintln(os.Stderr, "RESOURCE must be an exact alicloud resource type name")
		os.Exit(1)
	}
	packages, err := sweepPackages(".", *resource)
	var selected string
	if err == nil {
		selected, err = selectSweepPackage(packages, *explicit)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(selected)
}
