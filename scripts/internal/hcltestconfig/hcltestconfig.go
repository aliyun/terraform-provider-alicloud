// Package hcltestconfig parses split-layout Framework test configs written as
// fmt.Sprintf HCL helper functions — the plugin-testing dialect the standard
// template parser in scripts/testing cannot read, because a step's Config is a
// helper-function call (return fmt.Sprintf(`resource "..." { ... }`, args))
// rather than an inline map[string]interface{}.
//
// It extracts the attribute names from each helper's HCL resource block, maps
// fmt %[N]s placeholders to their 1-based argument position, and resolves the
// Config: helper(args) calls in TestAcc functions to determine which
// attributes are set (must) and which take different values across steps
// (modify). This lets the attribute-coverage checker measure a split-layout
// Framework resource against its real synthesized schema instead of silently
// passing it as "unmeasurable".
//
// It lives under scripts/internal so scripts/testing can share one
// implementation and so its parsing is unit-testable in memory (no provider
// link), mirroring scripts/internal/fwschema. Unknown syntax a test actually
// uses fails closed (an error), and a test with no resolvable HCL helper
// reports unmeasurable (measured=false) so the caller fails closed rather than
// silently passing.
package hcltestconfig

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Result holds the measured attribute coverage from a split-layout Framework
// test's fmt.Sprintf HCL helper configs.
type Result struct {
	// Must lists the attributes set in at least one config step.
	Must []string
	// Modify lists the attributes whose resolved value differs across steps.
	Modify []string
}

// ResolveFromFile reads a split-layout test file and resolves its HCL helper
// configs. See Resolve for the measurement contract.
func ResolveFromFile(filePath, typeName string) (measured bool, result Result, err error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return false, Result{}, fmt.Errorf("reading %s: %w", filePath, err)
	}
	return Resolve(src, typeName)
}

// Resolve parses split-layout Framework test configs written as fmt.Sprintf
// HCL helper functions and resolves which attributes are set (must) and which
// take different values across steps (modify), for the given resource type name
// (e.g. "alicloud_security_group_snapshot_policy").
//
// measured is true when at least one attribute was extracted from a resolvable
// HCL helper the test references via a Config: call. When false (nil error),
// the test carries no such helper: coverage is genuinely unmeasurable and the
// caller must fail closed rather than silently pass. A malformed helper the
// test uses is returned as an error so the caller fail-closes on a parse
// defect instead of skipping.
func Resolve(testSrc []byte, typeName string) (measured bool, result Result, err error) {
	content := string(testSrc)

	helpers, err := parseHelpers(content, typeName)
	if err != nil {
		return false, Result{}, err
	}
	if len(helpers) == 0 {
		return false, Result{}, nil
	}

	calls, err := parseConfigCalls(content, helpers)
	if err != nil {
		return false, Result{}, err
	}

	mustSet := map[string]bool{}
	modifySet := map[string]bool{}
	for helperName, steps := range calls {
		h := helpers[helperName]
		if h == nil {
			continue
		}
		// Every attribute declared in a helper's HCL block that the test
		// references via a Config call is set in at least one step.
		for attr := range h.attrs {
			mustSet[attr] = true
		}
		// An attribute whose placeholder-resolved value differs across the
		// steps that call this helper is modified. A literal value (argIdx 0)
		// is constant across steps and so is never a modify.
		for attr, argIdx := range h.attrs {
			if argIdx == 0 {
				continue
			}
			seen := map[string]bool{}
			for _, args := range steps {
				if argIdx-1 >= 0 && argIdx-1 < len(args) {
					seen[args[argIdx-1]] = true
				}
			}
			if len(seen) > 1 {
				modifySet[attr] = true
			}
		}
	}

	if len(mustSet) == 0 {
		// Helpers existed but no Config call referenced them: the test does
		// not exercise them, so coverage is unmeasurable.
		return false, Result{}, nil
	}
	return true, Result{Must: setToSlice(mustSet), Modify: setToSlice(modifySet)}, nil
}

func setToSlice(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// helper captures one fmt.Sprintf HCL config helper function. attrs maps an
// attribute name to its 1-based fmt argument index (the N in %[N]s, or the
// sequential position of a plain %s); 0 means the attribute is a literal value
// (constant across steps, never a modify).
type helper struct {
	name  string
	attrs map[string]int
}

var (
	// funcSigRegex matches a helper function signature: func <name>(<params>) string {
	// The params group [^)]* is sufficient for the plain identifier/string
	// parameter lists test config helpers use; a function-typed parameter
	// (with a nested ')') would not match, and such a helper falls through to
	// "unmeasurable" so the caller fail-closes rather than mis-parsing.
	funcSigRegex = regexp.MustCompile(`(?m)^func\s+(\w+)\s*\(([^)]*)\)\s*string\s*\{`)
	// hclAttrLineRegex matches one HCL attribute assignment: <name> = <value>.
	hclAttrLineRegex = regexp.MustCompile(`^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*(.+?)\s*$`)
	// hclIndexedVerbRegex matches an explicit fmt verb %[N]s and captures N.
	hclIndexedVerbRegex = regexp.MustCompile(`%\[(\d+)\]s`)
	// configLineRegex matches a Config: <value> line (multiline-aware so it
	// anchors at the start of each line).
	configLineRegex = regexp.MustCompile(`(?m)^[\t ]*Config:\s*(.+)$`)
	// funcStartRegex matches the start of a TestAcc function.
	funcStartRegex = regexp.MustCompile(`^func\s+TestAcc`)
	// commentLineRegex matches a commented-out line.
	commentLineRegex = regexp.MustCompile(`^\s*//`)
)

// parseHelpers finds fmt.Sprintf HCL config helper functions in the test source
// that declare the given resource type. A helper is a function returning string
// whose body calls fmt.Sprintf(`<HCL resource block>`, args). Only helpers
// whose HCL block declares the target resource type are kept, so a test file
// mixing several resources does not cross-contaminate coverage. A malformed
// helper (e.g. an invalid explicit index) is returned as an error so the caller
// fail-closes on a parse defect.
func parseHelpers(content, typeName string) (map[string]*helper, error) {
	helpers := map[string]*helper{}
	for _, loc := range funcSigRegex.FindAllStringSubmatchIndex(content, -1) {
		name := content[loc[2]:loc[3]]
		sprintfIdx := strings.Index(content[loc[1]:], "fmt.Sprintf(")
		if sprintfIdx < 0 {
			continue
		}
		sprintfIdx += loc[1]
		template, _, ok := extractRawString(content, sprintfIdx+len("fmt.Sprintf("))
		if !ok {
			continue
		}
		if !strings.Contains(template, "resource \""+typeName+"\"") {
			continue
		}
		h, err := parseHCLTemplate(template)
		if err != nil {
			return nil, fmt.Errorf("helper %s: %w", name, err)
		}
		if h != nil {
			h.name = name
			helpers[name] = h
		}
	}
	return helpers, nil
}

// extractRawString extracts a Go raw string literal (backtick-delimited)
// starting at content[idx] (skipping leading whitespace). Returns the content
// (without backticks), the index past the closing backtick, and ok=false when
// no raw string literal is present (e.g. an interpreted string or a non-string
// argument), which means the call is not an HCL config helper.
func extractRawString(content string, idx int) (text string, end int, ok bool) {
	for idx < len(content) && (content[idx] == ' ' || content[idx] == '\t' || content[idx] == '\n' || content[idx] == '\r') {
		idx++
	}
	if idx >= len(content) || content[idx] != '`' {
		return "", 0, false
	}
	start := idx + 1
	end = strings.IndexByte(content[start:], '`')
	if end < 0 {
		return "", 0, false
	}
	return content[start : start+end], start + end + 1, true
}

// parseHCLTemplate extracts attribute names and their fmt placeholder argument
// indices from an HCL resource block. An attribute line is `<name> = <value>`.
// If the value contains %[N]s, the attribute maps to argument N; a plain %s
// maps to the next sequential position; otherwise the attribute is a literal
// (argIdx 0, constant across steps). An invalid explicit index is reported as
// an error so the caller fail-closes on a malformed verb.
func parseHCLTemplate(template string) (*helper, error) {
	h := &helper{attrs: map[string]int{}}
	lines := strings.Split(template, "\n")
	inBlock := false
	depth := 0
	seq := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			continue
		}
		opens := strings.Count(trimmed, "{")
		closes := strings.Count(trimmed, "}")
		// Block headers (resource/data/variable/locals/provider/terraform)
		// open a block; only descend into a resource block's attributes.
		if strings.HasPrefix(trimmed, "resource ") || strings.HasPrefix(trimmed, "data ") ||
			strings.HasPrefix(trimmed, "variable ") || strings.HasPrefix(trimmed, "locals ") ||
			strings.HasPrefix(trimmed, "provider ") || strings.HasPrefix(trimmed, "terraform ") {
			if opens > 0 {
				inBlock = true
				depth += opens - closes
			}
			continue
		}
		if !inBlock {
			continue
		}
		depth += opens - closes
		if depth <= 0 {
			inBlock = false
			continue
		}
		// Attribute assignment: name = value.
		m := hclAttrLineRegex.FindStringSubmatch(trimmed)
		if m == nil {
			continue
		}
		attrName := m[1]
		value := m[2]
		if im := hclIndexedVerbRegex.FindStringSubmatch(value); im != nil {
			n, err := strconv.Atoi(im[1])
			if err != nil {
				return nil, fmt.Errorf("attribute %s: invalid placeholder %q: %w", attrName, im[0], err)
			}
			h.attrs[attrName] = n
		} else if strings.Contains(value, "%s") {
			seq++
			h.attrs[attrName] = seq
		} else {
			h.attrs[attrName] = 0 // literal value
		}
	}
	if len(h.attrs) == 0 {
		return nil, nil
	}
	return h, nil
}

// parseConfigCalls finds Config: helper(args) lines inside TestAcc functions
// and resolves the positional call arguments per helper, preserving step
// order. Returns a map of helperName -> [][]string (one arg slice per step).
// Only helpers present in the helpers map are tracked; an inline
// map[string]interface{} Config (not a helper call) is ignored here and left
// to the standard template parser. An unterminated argument list is an error so
// the caller fail-closes rather than silently passing on a malformed call.
func parseConfigCalls(content string, helpers map[string]*helper) (map[string][][]string, error) {
	calls := map[string][][]string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	inFunc := false
	for scanner.Scan() {
		text := scanner.Text()
		if commentLineRegex.MatchString(text) {
			continue
		}
		if funcStartRegex.MatchString(text) {
			inFunc = true
			continue
		}
		if !inFunc {
			continue
		}
		m := configLineRegex.FindStringSubmatch(text)
		if m == nil {
			continue
		}
		configValue := strings.TrimSpace(m[1])
		configValue = strings.TrimSuffix(configValue, ",")
		openParen := strings.Index(configValue, "(")
		if openParen <= 0 {
			continue
		}
		helperName := strings.TrimSpace(configValue[:openParen])
		if helpers[helperName] == nil {
			continue
		}
		args, err := splitCallArgs(configValue[openParen+1:])
		if err != nil {
			return nil, fmt.Errorf("Config call to %s: %w", helperName, err)
		}
		calls[helperName] = append(calls[helperName], args)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return calls, nil
}

// splitCallArgs splits comma-separated Go call arguments — the content from
// just after the call's opening parenthesis to its matching close — into
// individual argument strings, respecting string literals, raw strings,
// parentheses, brackets and braces so a comma inside a nested expression is
// not mistaken for a separator. The closing ')' at depth 0 terminates parsing
// and is not included.
func splitCallArgs(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	depth := 1 // inside the call's opening paren
	var cur strings.Builder
	var args []string
	inStr := false
	inRaw := false
	flush := func() {
		if v := strings.TrimSpace(cur.String()); v != "" {
			args = append(args, v)
		}
		cur.Reset()
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case inRaw:
			if c == '`' {
				inRaw = false
			}
			cur.WriteByte(c)
		case inStr:
			switch c {
			case '\\':
				cur.WriteByte(c)
				if i+1 < len(s) {
					i++
					cur.WriteByte(s[i])
				}
			case '"':
				inStr = false
				cur.WriteByte(c)
			default:
				cur.WriteByte(c)
			}
		default:
			switch c {
			case '`':
				inRaw = true
				cur.WriteByte(c)
			case '"':
				inStr = true
				cur.WriteByte(c)
			case '(', '[', '{':
				depth++
				cur.WriteByte(c)
			case ')', ']', '}':
				depth--
				if depth == 0 && c == ')' {
					flush()
					return args, nil
				}
				cur.WriteByte(c)
			case ',':
				if depth == 1 {
					flush()
				} else {
					cur.WriteByte(c)
				}
			default:
				cur.WriteByte(c)
			}
		}
	}
	// No matching close paren: the call is unterminated (possibly a multi-line
	// call this single-line scan did not fully capture). Fail closed with a
	// clear message rather than returning partial args that would mislead
	// coverage measurement.
	return nil, fmt.Errorf("unterminated argument list (no matching ')')")
}
