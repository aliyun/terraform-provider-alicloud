package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// providerRoot is the directory holding the provider source tree. The checker
// runs from the repository root in CI, so "." is the default; tests point it
// at a fixture tree.
var providerRoot = "."

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: document_check <doc> [doc ...]")
		os.Exit(1)
	}
	// One process checks every document on the command line so the
	// registration index built for the first page serves the rest; each
	// failing page still reports its own exit code.
	exitCode := 0
	for _, doc := range os.Args[1:] {
		if code := checkDocument(strings.TrimSpace(doc)); code != 0 {
			exitCode = code
		}
	}
	os.Exit(exitCode)
}

// registeredType describes a type registered by a service package.
type registeredType struct {
	// definitionDir is the package directory owning the register.go file; the
	// registered factory function is defined in a file of the same package.
	definitionDir string
	// factory is the name of the registered factory function.
	factory string
	// section is the register.go section the entry was found in, e.g.
	// SDKResources or SDKDataSources.
	section string
}

// docKind identifies which provider form a document describes. The docs
// directory decides it: r/ pages document resources, d/ pages data sources,
// e/ or ephemeral-resources/ pages ephemeral resources. The kind selects
// which registration sections may match, so a TypeName registered both as a
// resource and as a data source (about 20 names today) reconciles against
// the form the page actually documents instead of the first registration
// found.
type docKind int

const (
	// kindUnknown matches any section; it keeps bare file names (e.g. the
	// checker's own tests) working and never newly skips a reconciliation.
	kindUnknown docKind = iota
	kindResource
	kindDataSource
	kindEphemeral
)

// sectionsForKind returns the register.go section names matching the kind,
// or nil when any section may match.
func sectionsForKind(kind docKind) map[string]bool {
	switch kind {
	case kindResource:
		return map[string]bool{"SDKResources": true, "Resources": true}
	case kindDataSource:
		return map[string]bool{"SDKDataSources": true, "DataSources": true}
	case kindEphemeral:
		return map[string]bool{"EphemeralResources": true}
	}
	return nil
}

var registerEntryPattern = regexp.MustCompile(`TypeName:\s*"([^"]+)"\s*,\s*Factory:\s*([A-Za-z0-9_]+)`)

// registerSectionPattern matches the section headers of a
// conns.ServicePackage literal. Longer names come first so Resources does
// not swallow SDKResources.
var registerSectionPattern = regexp.MustCompile(`(?m)^[ \t]*(SDKResources|SDKDataSources|EphemeralResources|Resources|DataSources|Functions|ListResources|Actions)[ \t]*:`)

// registrationIndex maps a TypeName to all of its registrations in the
// service packages. One WalkDir of alicloud/service (about 170 register.go
// files) fills it, so checking a batch of documents rescans the tree once
// per process instead of once per page. registrationIndexRoot records the
// provider root the index was built from; the tests point providerRoot at
// different fixture trees and must never see entries of a previous root.
var (
	registrationIndex     map[string][]registeredType
	registrationIndexRoot string
)

// buildRegistrationIndex walks every register.go under alicloud/service once
// and records each TypeName entry together with the section it was declared
// in. It is called through findRegisteredType; the returned map is also kept
// in registrationIndex for reuse.
func buildRegistrationIndex() map[string][]registeredType {
	registrationIndex = map[string][]registeredType{}
	registrationIndexRoot = providerRoot
	serviceRoot := filepath.Join(providerRoot, "alicloud", "service")
	if info, err := os.Stat(serviceRoot); err != nil || !info.IsDir() {
		return registrationIndex
	}
	_ = filepath.WalkDir(serviceRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || entry.Name() != "register.go" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		headerIndexes := registerSectionPattern.FindAllStringSubmatchIndex(string(content), -1)
		sectionAt := func(offset int) string {
			current := ""
			for _, match := range headerIndexes {
				if match[0] > offset {
					break
				}
				current = string(content[match[2]:match[3]])
			}
			return current
		}
		for _, match := range registerEntryPattern.FindAllStringSubmatchIndex(string(content), -1) {
			typeName := string(content[match[2]:match[3]])
			registrationIndex[typeName] = append(registrationIndex[typeName], registeredType{
				definitionDir: filepath.Dir(path),
				factory:       string(content[match[4]:match[5]]),
				section:       sectionAt(match[0]),
			})
		}
		return nil
	})
	return registrationIndex
}

// findRegisteredType locates the registration of the given type in the
// service packages under alicloud/service, restricted to the sections of the
// given doc kind. It returns nil for unregistered names, e.g. the provider
// index page or the checker's own test documents, in which case no
// deprecation reconciliation applies. The registration index is shared per
// process and rebuilt when providerRoot changes.
func findRegisteredType(typeName string, kind docKind) *registeredType {
	if registrationIndex == nil || registrationIndexRoot != providerRoot {
		buildRegistrationIndex()
	}
	sections := sectionsForKind(kind)
	for _, registered := range registrationIndex[typeName] {
		if sections == nil || sections[registered.section] {
			return &registered
		}
	}
	return nil
}

// declaresTypeLevelDeprecation reports whether the definition of the given
// registered type declares a type-level deprecation, i.e. a
// DeprecationMessage on the resource or data source schema itself (the
// schema.Resource literal of an SDKv2 factory, or the schema of the type a
// Framework factory returns, both in the factory's definition file).
// Attribute-level deprecations use the Deprecated: field and do not match.
func declaresTypeLevelDeprecation(typeName string, kind docKind) bool {
	registered := findRegisteredType(typeName, kind)
	if registered == nil {
		return false
	}
	files, err := os.ReadDir(registered.definitionDir)
	if err != nil {
		return false
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") || strings.HasSuffix(file.Name(), "_test.go") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(registered.definitionDir, file.Name()))
		if err != nil {
			continue
		}
		if strings.Contains(string(content), "func "+registered.factory+"(") {
			return declaresTypeLevelDeprecationIn(string(content))
		}
	}
	return false
}

// typeLevelDeprecationLiterals lists the composite literal types whose
// DeprecationMessage field is the type-level deprecation this checker
// reconciles: the schema.Resource an SDKv2 factory returns, and the
// schema.Schema a Framework (or ephemeral) Schema method assigns. A
// DeprecationMessage on any other literal, e.g. schema.StringAttribute
// inside an attribute map, deprecates only that attribute and must not
// demand a DEPRECATED note for the whole page.
var typeLevelDeprecationLiterals = map[string]bool{
	"schema.Schema":   true,
	"schema.Resource": true,
}

// declaresTypeLevelDeprecationIn reports whether content contains a
// DeprecationMessage field of a type-level composite literal. A plain
// substring search is not enough: the same field name exists on
// attribute-level literals (e.g. schema.StringAttribute), and the token can
// also appear inside strings and comments. The scan walks the source as
// text, tracks the type path of the innermost open composite literal, and
// skips strings, rune literals and comments, so the field counts only when
// its enclosing literal is one of the whitelisted types.
func declaresTypeLevelDeprecationIn(content string) bool {
	// path accumulates the identifier bytes since the last separator; a
	// qualified name such as schema.Schema is kept with its dots.
	path := make([]byte, 0, 32)
	// literalStack holds the type path of each open composite literal,
	// innermost last. Whitespace continues a path so `schema.Schema {`
	// still records its type; any other separator resets it.
	var literalStack []string
	for i := 0; i < len(content); i++ {
		switch c := content[i]; c {
		case '"':
			i = skipString(content, i)
		case '`':
			i = skipRawString(content, i)
		case '\'':
			i = skipRune(content, i)
		case '/':
			if i+1 < len(content) && content[i+1] == '/' {
				i = skipLineComment(content, i)
			} else if i+1 < len(content) && content[i+1] == '*' {
				i = skipBlockComment(content, i)
			} else {
				path = path[:0]
			}
		case '{':
			literalStack = append(literalStack, string(path))
			path = path[:0]
		case '}':
			if len(literalStack) > 0 {
				literalStack = literalStack[:len(literalStack)-1]
			}
			path = path[:0]
		case ':':
			if len(literalStack) > 0 && string(path) == "DeprecationMessage" &&
				typeLevelDeprecationLiterals[literalStack[len(literalStack)-1]] {
				return true
			}
			path = path[:0]
		default:
			if isIdentByte(c) || c == '.' {
				path = append(path, c)
			} else if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
				path = path[:0]
			}
		}
	}
	return false
}

// isIdentByte reports whether c can appear in a Go identifier.
func isIdentByte(c byte) bool {
	return c == '_' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9'
}

// skipString returns the index of the closing quote of the interpreted
// string literal starting at content[start]; the caller advances past it.
func skipString(content string, start int) int {
	for i := start + 1; i < len(content); i++ {
		switch content[i] {
		case '\\':
			i++
		case '"':
			return i
		}
	}
	return len(content) - 1
}

// skipRawString returns the index of the closing backquote of the raw
// string literal starting at content[start].
func skipRawString(content string, start int) int {
	if end := strings.IndexByte(content[start+1:], '`'); end >= 0 {
		return start + 1 + end
	}
	return len(content) - 1
}

// skipRune returns the index of the closing quote of the rune literal
// starting at content[start].
func skipRune(content string, start int) int {
	for i := start + 1; i < len(content); i++ {
		switch content[i] {
		case '\\':
			i++
		case '\'':
			return i
		}
	}
	return len(content) - 1
}

// skipLineComment returns the index of the newline ending the line comment
// starting at content[start].
func skipLineComment(content string, start int) int {
	if end := strings.IndexByte(content[start:], '\n'); end >= 0 {
		return start + end
	}
	return len(content) - 1
}

// skipBlockComment returns the index of the slash closing the block comment
// starting at content[start].
func skipBlockComment(content string, start int) int {
	if end := strings.Index(content[start+2:], "*/"); end >= 0 {
		return start + 2 + end + 1
	}
	return len(content) - 1
}

// Page categories under website/docs. Only pages documenting a Terraform
// object (resource, data source, action, ephemeral resource or provider
// function) follow a generated-name title template; guides and the provider
// index carry business titles and must not be forced into the alicloud_
// resource-name template.
type pageCategory int

const (
	pageCategoryObject    pageCategory = iota // website/docs/r, /d, /actions and unknown object directories
	pageCategoryEphemeral                     // website/docs/e, /ephemeral-resources
	pageCategoryFunction                      // website/docs/functions
	pageCategoryGuide                         // website/docs/guides and bare file names without a directory
)

// classifyPage derives the page category from the directory that owns the
// page. Both published suffixes (.html.markdown and .html.md) resolve to the
// same page name so that guides written as .html.md are not turned into a
// broken resource name such as alicloud_xxx.html.md.
func classifyPage(docsFileName string) (pageCategory, string) {
	docsFileNameParts := strings.Split(docsFileName, "/")
	fileName := docsFileNameParts[len(docsFileNameParts)-1]
	pageName := strings.TrimSuffix(strings.TrimSuffix(fileName, ".html.markdown"), ".html.md")
	if len(docsFileNameParts) < 2 {
		// A bare file name carries no owning directory, so it cannot be
		// classified as an object page. The CI file loop always passes
		// website/docs/... paths; this branch only covers defensive
		// invocations, while pages really under website/docs keep the
		// object classification below.
		return pageCategoryGuide, pageName
	}
	switch docsFileNameParts[len(docsFileNameParts)-2] {
	case "e", "ephemeral-resources":
		return pageCategoryEphemeral, pageName
	case "functions":
		return pageCategoryFunction, pageName
	case "guides":
		return pageCategoryGuide, pageName
	default:
		return pageCategoryObject, pageName
	}
}

func checkDocument(docsFileName string) int {
	docsFile, err := os.OpenFile(docsFileName, os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("open docs file %s failed. Error:%s", docsFileName, err)
		return 1
	}
	defer docsFile.Close()
	scanner := bufio.NewScanner(docsFile)
	category, pageName := classifyPage(docsFileName)
	resourceName := "alicloud_" + pageName
	// expectedTitle and pageTitleExpected stay empty for guide pages: their
	// titles are business titles and are only checked for presence, never
	// against the resource-name template.
	expectedTitle := "# " + resourceName
	pageTitleExpected := "page_title: \"Alicloud: " + resourceName + "\""
	// kind selects which registration sections the deprecation
	// reconciliation may match for this page; guide and function pages stay
	// kindUnknown so any section may match, exactly as pages outside r/ and
	// d/ were matched before the page categories existed.
	kind := kindUnknown
	switch category {
	case pageCategoryEphemeral:
		expectedTitle = "# Ephemeral: " + resourceName
		kind = kindEphemeral
	case pageCategoryFunction:
		expectedTitle = "# Function: " + pageName
		pageTitleExpected = "page_title: \"Alibaba Cloud: " + pageName + "\""
	case pageCategoryGuide:
		expectedTitle = ""
		pageTitleExpected = ""
	default:
		// Object pages document a Terraform type: r/ pages resources and
		// d/ pages data sources; other object directories stay unknown.
		// classifyPage only returns the object category for paths with an
		// owning directory, so the index below is always valid.
		docsFileNameParts := strings.Split(docsFileName, "/")
		switch docsFileNameParts[len(docsFileNameParts)-2] {
		case "r":
			kind = kindResource
		case "d":
			kind = kindDataSource
		}
	}
	exitCode := 0
	fmt.Printf("\n==> Checking docs content of %s ...", docsFileName)
	titleCheck := false
	titleDescriptionCheck := false
	descriptionCheck := false
	versionChecked := false
	deprecatedNoteSeen := false
	exampleCheck := false
	exampleBlockOpen := false
	argumentCheck := false
	attributesCheck := false
	importCheck := false
	importBlockOpen := false
	timeoutCheck := false
	line := 0
	for scanner.Scan() {
		line += 1
		text := scanner.Text()
		if line == 1 {
			titleCheck = true
			if text != "---" {
				fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "---", text)
				exitCode = 1
			}
			continue
		}
		if titleCheck {
			if strings.HasPrefix(text, "layout:") && text != "layout: \"alicloud\"" {
				fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "layout: \"alicloud\"", text)
				exitCode = 1
				continue
			}
			if strings.HasPrefix(text, "page_title:") {
				if pageTitleExpected != "" {
					if text != pageTitleExpected {
						fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, pageTitleExpected, text)
						exitCode = 1
					}
				} else {
					// Guide pages keep their business title; only a
					// non-empty quoted page_title is required.
					value := strings.TrimSpace(strings.TrimPrefix(text, "page_title:"))
					if len(value) < 3 || !strings.HasPrefix(value, "\"") || !strings.HasSuffix(value, "\"") ||
						strings.TrimSpace(value[1:len(value)-1]) == "" {
						fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "page_title: \"<non-empty guide title>\"", text)
						exitCode = 1
					}
				}
				continue
			}
			if strings.HasPrefix(text, "description:") {
				titleDescriptionCheck = true
				if text != "description: |-" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "description: |- ", text)
					exitCode = 1
				}
				continue
			}
			if titleDescriptionCheck {
				titleDescriptionCheck = false
				if !strings.HasPrefix(text, "  ") {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "  "+text, text)
					exitCode = 1
				}
				continue
			}
			if text == "---" {
				titleCheck = false
				descriptionCheck = true
				continue
			}
		}
		if descriptionCheck {
			// Object pages keep the historical title trigger; provider
			// function pages match their own heading; guide pages have no
			// title template and are never compared against expectedTitle.
			isBodyTitle := strings.HasPrefix(text, "# alicloud") || strings.HasSuffix(text, resourceName)
			switch category {
			case pageCategoryFunction:
				isBodyTitle = strings.HasPrefix(text, "# Function:")
			case pageCategoryGuide:
				isBodyTitle = false
			}
			if isBodyTitle {
				if titleCheck {
					fmt.Printf("\n[Error] line %d: docs title has not been closed.", line)
					exitCode = 1
				}
				if strings.Contains(text, "\\") {
					fmt.Printf("\n[WARNING] line %d: please remove the \\", line)
					text = strings.Replace(text, "\\", "", -1)
				}
				if text != expectedTitle {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, expectedTitle, text)
					exitCode = 1
				}
				continue
			}
			if strings.Contains(text, "https://help.aliyun.com/") {
				fmt.Printf("\n[Error] line %d: Expected: an international site link: %s. Got: %s.", line, "https://www.alibabacloud.com/help", "https://help.aliyun.com/")
				exitCode = 1
				continue
			}
			if strings.Contains(text, "Available since v") {
				versionChecked = true
				v := "Available since v"
				if strings.Contains(text, v) && !strings.HasPrefix(text, "-> **NOTE:** "+v) {
					parts := strings.Split(text, v)
					fmt.Printf("\n[Error] line %d: Expected: %s Got: %s", line, "-> **NOTE:** "+v+strings.Replace(parts[1], "+", "", -1), text)
					exitCode = 1
				}
				continue
			}
			if strings.Contains(text, "Available in v") {
				versionChecked = true
				v := "Available in v"
				if strings.Contains(text, v) {
					parts := strings.Split(text, v)
					v = strings.Replace(strings.Replace(v, "in", "since", -1), "from", "since", -1)
					fmt.Printf("\n[Error] line %d: Expected: %s Got: %s", line, parts[0]+v+strings.Replace(parts[1], "+", "", -1), text)
					exitCode = 1
				}
				continue
			}
			if strings.Contains(text, "-> **DEPRECATED:**") || strings.Contains(text, "-> **NOTE:** deprecated since") ||
				strings.Contains(text, "-> **NOTE:** Deprecated since") ||
				(strings.Contains(text, "This resource has been") && (strings.Contains(text, "deprecated") || strings.Contains(text, "DEPRECATED"))) {
				versionChecked = true
				deprecatedNoteSeen = true
				if !strings.HasPrefix(text, "-> **DEPRECATED:**") {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "-> **DEPRECATED:** This [resource|data source] has been deprecated from vxxx", text)
					exitCode = 1
				}
				continue
			}

			if strings.HasPrefix(text, "## Example Usage") || strings.HasSuffix(text, "Example Usage") {
				descriptionCheck = false
				exampleCheck = true
				if !versionChecked {
					fmt.Printf("\n[Error] line %d: missing available or deprecated verison info.", line)
					exitCode = 1
				}
				if text != "## Example Usage" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Example Usage", text)
					exitCode = 1
				}
				continue
			}
		}
		if exampleCheck {
			if strings.HasPrefix(text, "```") {
				if !exampleBlockOpen {
					exampleBlockOpen = true
				} else {
					exampleBlockOpen = false
					if text != "```" {
						fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "```", text)
						exitCode = 1
					}
				}
				if exampleBlockOpen && text != "```terraform" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "```terraform", text)
					exitCode = 1
				}
				continue
			}
			if exampleBlockOpen && strings.Contains(text, "test") {
				text = strings.TrimSpace(text)
				parts := strings.Split(text, "=")
				if key := strings.TrimSpace(parts[0]); key == "default" || strings.HasSuffix(key, "name") {
					if value := strings.TrimSpace(parts[1]); strings.Contains(value, "test") {
						fmt.Printf("\n[Error] line %d: avoid using 'test' as soon as possable in the example", line)
						exitCode = 1
					}
				}
			}
			if exampleBlockOpen && strings.Contains(text, "depends_on") {
				text = strings.TrimSpace(text)
				parts := strings.Split(text, "=")
				if key := strings.TrimSpace(parts[0]); key == "depends_on" {
					fmt.Printf("\n[Error] line %d: avoid using 'depends_on' as soon as possable in the example", line)
					exitCode = 1
				}
			}
			// Provider function pages continue with their own sections
			// instead of the object-page section walk: the Signature,
			// Arguments and Return Type headings terminate the example
			// section, so the ```text signature fence is not mistaken for
			// a Terraform example block.
			if category == pageCategoryFunction &&
				(strings.HasPrefix(text, "## Signature") || strings.HasPrefix(text, "## Arguments") || strings.HasPrefix(text, "## Return Type")) {
				exampleCheck = false
				if exampleBlockOpen {
					fmt.Printf("\n[Error] line %d: example block has not been closed.", line)
					exitCode = 1
				}
				continue
			}
			if strings.HasPrefix(text, "## Argument Reference") || strings.HasSuffix(text, "Argument Reference") {
				exampleCheck = false
				argumentCheck = true
				if exampleBlockOpen {
					fmt.Printf("\n[Error] line %d: example block has not been closed.", line)
					exitCode = 1
				}
				if text != "## Argument Reference" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Argument Reference", text)
					exitCode = 1
				}
				continue
			}
			if !argumentCheck && strings.HasPrefix(text, "## Attributes Reference") || strings.HasSuffix(text, "Attributes Reference") {
				exampleCheck = false
				attributesCheck = true
				if text != "## Attributes Reference" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Attributes Reference", text)
					exitCode = 1
				}
				continue
			}
		}
		if argumentCheck {
			if strings.HasPrefix(text, "* `") {
				parts := strings.Split(strings.Split(text, ")")[0], "(")
				for i, tag := range strings.Split(parts[len(parts)-1], ",") {
					if strings.HasPrefix(tag, " ") {
						if i == 0 {
							fmt.Printf("\n[Error] line %d: please remove redundant space prefix for %s. ", line, tag)
							exitCode = 1
						}
					} else if i > 0 {
						fmt.Printf("\n[Error] line %d: missing space prefix for %s. ", line, tag)
						exitCode = 1
					}
					if strings.HasSuffix(tag, " ") {
						fmt.Printf("\n[Error] line %d: please remove redundant space suffix for %s.", line, tag)
						exitCode = 1
					}
				}
				continue
			}
			if strings.HasPrefix(text, "###") {
				parts := strings.Split(text, " ")
				block := strings.ToLower(parts[len(parts)-1])
				if len(parts) > 2 ||
					(strings.HasPrefix(block, "`") && !strings.HasSuffix(block, "`")) ||
					(!strings.HasPrefix(block, "`") && strings.HasSuffix(block, "`")) {
					block = strings.Trim(block, "`")
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "### `"+block+"`", text)
					exitCode = 1
				}
				block = strings.Trim(block, "`")
				blockParts := strings.Split(block, "-")
				blockLink := fmt.Sprintf("[`%s`](#%s)", blockParts[len(blockParts)-1], block)
				docsContent, err := os.ReadFile(docsFileName)
				if err != nil {
					fmt.Printf("\n[Error] reading docs file %s failed. Error: %s", docsFileName, err)
					exitCode = 1
				} else if !strings.Contains(string(docsContent), blockLink) {
					fmt.Printf("\n[Error] line %d: missing link for block `%s`. Expected link like: See %s below.", line, block, blockLink)
					exitCode = 1
				}
				continue
			}
			if strings.HasPrefix(text, "## Argument Reference") || strings.HasSuffix(text, "Argument Reference") {
				argumentCheck = false
				attributesCheck = true
				if text != "## Attributes Reference" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Attributes Reference", text)
					exitCode = 1
				}
				continue
			}
			if strings.HasPrefix(text, "## Attributes Reference") || strings.HasSuffix(text, "Attributes Reference") {
				argumentCheck = false
				attributesCheck = true
				if text != "## Attributes Reference" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Attributes Reference", text)
					exitCode = 1
				}
				continue
			}
		}

		if attributesCheck {
			if strings.HasPrefix(text, "## Timeout") || strings.HasSuffix(text, "Timeouts") || strings.HasSuffix(text, "Timeout") {
				attributesCheck = false
				timeoutCheck = true
				if text != "## Timeouts" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Timeouts", text)
					exitCode = 1
				}
				continue
			}
			if strings.HasPrefix(text, "## Import") || strings.HasSuffix(text, "Import") {
				attributesCheck = false
				importCheck = true
				if text != "## Import" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Import", text)
					exitCode = 1
				}
				continue
			}
		}

		if timeoutCheck {
			if strings.HasPrefix(text, "## Import") || strings.HasSuffix(text, "Import") {
				timeoutCheck = false
				importCheck = true
				if text != "## Import" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "## Import", text)
					exitCode = 1
				}
				continue
			}
		}

		if importCheck {
			if !importBlockOpen && strings.Contains(text, "terraform import") {
				fmt.Printf("\n[Error] line %d: Expected: %s.", line-1, "```shell")
				exitCode = 1
				break
			}

			if strings.HasPrefix(text, "```") {
				if !importBlockOpen {
					importBlockOpen = true
				} else {
					importBlockOpen = false
					if text != "```" {
						fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "```", text)
						exitCode = 1
					}
				}
				if importBlockOpen && text != "```shell" {
					fmt.Printf("\n[Error] line %d: Expected: %s. Got: %s.", line, "```shell", text)
					exitCode = 1
				}
				continue
			}
		}
	}
	if importBlockOpen {
		fmt.Println("\n[Error] Error: import block has not been closed.")
		exitCode = 1
	}
	// Reconcile the page against the provider schema: a type-level
	// DeprecationMessage must be stated in the description, otherwise a
	// deleted DEPRECATED note would silently pass the format-only checks.
	if !deprecatedNoteSeen && declaresTypeLevelDeprecation(resourceName, kind) {
		fmt.Printf("\n[Error] %s is deprecated in the provider schema (DeprecationMessage); the description must keep a \"-> **DEPRECATED:**\" statement.", resourceName)
		exitCode = 1
	}
	fmt.Println("\n--- Finished!")
	return exitCode
}
