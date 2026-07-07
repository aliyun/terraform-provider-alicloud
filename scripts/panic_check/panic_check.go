//nolint:all
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var (
	fileNames    = flag.String("fileNames", "", "path to a file containing a list of Go files to check (one per line); if empty, scans all .go files under alicloud/")
	includeTests = flag.Bool("includeTests", false, "include _test.go files in the check")
	exclude      = flag.String("exclude", "", "comma-separated list of file path substrings to exclude (e.g., common.go)")
)

func main() {
	flag.Parse()

	var files []string

	if fileNames != nil && len(*fileNames) > 0 {
		// Read file list from the given file
		data, err := os.ReadFile(*fileNames)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file list %s: %v\n", *fileNames, err)
			os.Exit(2)
		}
		for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
			line = strings.TrimSpace(line)
			if line != "" && strings.HasSuffix(line, ".go") {
				files = append(files, line)
			}
		}
	} else {
		// Scan all .go files under alicloud/
		var err error
		files, err = findGoFiles("alicloud")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error scanning alicloud/: %v\n", err)
			os.Exit(2)
		}
	}

	// Build exclude list
	var excludePatterns []string
	if exclude != nil && len(*exclude) > 0 {
		excludePatterns = strings.Split(*exclude, ",")
		for i := range excludePatterns {
			excludePatterns[i] = strings.TrimSpace(excludePatterns[i])
		}
	}

	exitCode := 0
	checkedCount := 0
	panicCount := 0

	for _, file := range files {
		// Skip test files unless -includeTests is set
		if !*includeTests && strings.HasSuffix(file, "_test.go") {
			continue
		}

		// Skip excluded files
		if isExcluded(file, excludePatterns) {
			continue
		}

		// Skip files that don't exist (e.g., deleted in diff)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		checkedCount++
		panics := checkFileForPanics(file)
		for _, p := range panics {
			fmt.Printf("%s:%d:%d: panic() call in function %s\n", p.file, p.line, p.col, p.funcName)
			panicCount++
			exitCode = 1
		}
	}

	if exitCode == 0 {
		fmt.Printf("panic-check: OK (checked %d files, 0 panic calls found)\n", checkedCount)
	} else {
		fmt.Printf("panic-check: FAIL (checked %d files, %d panic call(s) found)\n", checkedCount, panicCount)
	}

	os.Exit(exitCode)
}

type panicLocation struct {
	file     string
	line     int
	col      int
	funcName string
}

func checkFileForPanics(filename string) []panicLocation {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not parse %s: %v\n", filename, err)
		return nil
	}

	var panics []panicLocation
	currentFunc := "<package-level>"

	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			// Track the enclosing function name
			currentFunc = node.Name.Name
			// Continue inspecting the function body
			return true
		case *ast.CallExpr:
			if ident, ok := node.Fun.(*ast.Ident); ok && ident.Name == "panic" {
				pos := fset.Position(node.Pos())
				panics = append(panics, panicLocation{
					file:     filename,
					line:     pos.Line,
					col:      pos.Column,
					funcName: currentFunc,
				})
			}
			return true
		}
		return true
	})

	return panics
}

func findGoFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func isExcluded(file string, patterns []string) bool {
	for _, p := range patterns {
		if p != "" && strings.Contains(file, p) {
			return true
		}
	}
	return false
}
