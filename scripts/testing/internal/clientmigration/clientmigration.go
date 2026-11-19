// Package clientmigration recognizes the connectivity.Config client constructor
// migration in test files without excluding changes to the tests themselves.
package clientmigration

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
	"regexp"
	"strings"

	"github.com/waigani/diffparser"
)

const connectivityImport = `"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"`

var blobIndex = regexp.MustCompile(`(?m)^index ([0-9a-f]{4,64})\.\.([0-9a-f]{4,64})(?: [0-7]+)?$`)

// OnlyChange reports whether reversing just the added NewClient calls
// reproduces the entire original blob. Missing objects retain full checking.
func OnlyChange(file *diffparser.DiffFile) bool {
	if file.Mode != diffparser.MODIFIED || file.OrigName != file.NewName {
		return false
	}
	lines := make(map[int]bool)
	for _, h := range file.Hunks {
		for _, line := range h.NewRange.Lines {
			if line.Mode == diffparser.ADDED && strings.Contains(line.Content, "connectivity.NewClient(") {
				lines[line.Number] = true
			}
		}
	}
	if len(lines) == 0 {
		return false
	}
	ids := blobIndex.FindStringSubmatch(file.DiffHeader)
	if ids == nil {
		return false
	}
	old, err := exec.Command("git", "cat-file", "blob", ids[1]).Output()
	if err != nil {
		return false
	}
	current, err := exec.Command("git", "cat-file", "blob", ids[2]).Output()
	return err == nil && equivalent(old, current, lines)
}

func equivalent(old, current []byte, addedLines map[int]bool) bool {
	positions := token.NewFileSet()
	file, err := parser.ParseFile(positions, "", current, 0)
	if err != nil {
		return false
	}
	imported := false
	for _, spec := range file.Imports {
		if spec.Path.Value == connectivityImport && (spec.Name == nil || spec.Name.Name == "connectivity") {
			imported = true
		}
	}
	if !imported {
		return false
	}
	var restored bytes.Buffer
	cursor := 0
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || !selector(call.Fun, "NewClient") || len(call.Args) != 1 {
			return true
		}
		start := positions.PositionFor(call.Pos(), false)
		if !addedLines[start.Line] {
			return true
		}
		end := positions.PositionFor(call.End(), false).Offset
		arg := call.Args[0]
		argText := string(current[positions.PositionFor(arg.Pos(), false).Offset:positions.PositionFor(arg.End(), false).Offset])
		previous := previousCall(arg, argText)
		if previous == "" || string(current[start.Offset:end]) != "connectivity.NewClient("+argText+")" {
			return true
		}
		restored.Write(current[cursor:start.Offset])
		restored.WriteString(previous)
		cursor = end
		return false
	})
	restored.Write(current[cursor:])
	return cursor != 0 && bytes.Equal(restored.Bytes(), old)
}

func selector(expr ast.Expr, name string) bool {
	s, ok := expr.(*ast.SelectorExpr)
	if !ok || s.Sel.Name != name {
		return false
	}
	pkg, ok := s.X.(*ast.Ident)
	return ok && pkg.Name == "connectivity" && pkg.Obj == nil
}

func addressOperand(expr ast.Expr) (ast.Expr, bool) {
	if addr, ok := expr.(*ast.UnaryExpr); ok && addr.Op == token.AND {
		return addr.X, true
	}
	return expr, false
}

func previousCall(arg ast.Expr, text string) string {
	operand, addressed := addressOperand(arg)
	if literal, ok := operand.(*ast.CompositeLit); ok && addressed && selector(literal.Type, "Config") {
		return "(" + text + ").Client()"
	}
	id, ok := operand.(*ast.Ident)
	if !ok || id.Obj == nil {
		return ""
	}
	decl, ok := id.Obj.Decl.(*ast.AssignStmt)
	if !ok || decl.Tok != token.DEFINE || len(decl.Lhs) != 1 || len(decl.Rhs) != 1 {
		return ""
	}
	value, pointer := addressOperand(decl.Rhs[0])
	lit, ok := value.(*ast.CompositeLit)
	if !ok || !selector(lit.Type, "Config") || addressed == pointer {
		return ""
	}
	if text != id.Name && text != "&"+id.Name {
		return ""
	}
	return id.Name + ".Client()"
}
