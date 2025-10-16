//nolint:all
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/waigani/diffparser"
)

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.DisableTimestamp = false
	customFormatter.DisableColors = false
	customFormatter.ForceColors = true
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var (
	fileNames = flag.String("fileNames", "", "the files to check diff")
)

func main() {
	exitCode := 0
	flag.Parse()
	if fileNames != nil && len(*fileNames) == 0 {
		log.Warningf("the diff file is empty")
		return
	}
	byt, _ := ioutil.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	fileRegex := regexp.MustCompile("alicloud/resource[0-9a-zA-Z_]*.go")
	fileTestRegex := regexp.MustCompile("alicloud/resource[0-9a-zA-Z_]*_test.go")
	for _, file := range diff.Files {
		fmt.Println()
		if fileRegex.MatchString(file.NewName) {
			if fileTestRegex.MatchString(file.NewName) {
				continue
			}
			resourceName := strings.TrimPrefix(strings.TrimSuffix(strings.Split(file.NewName, "/")[1], ".go"), "resource_")
			log.Infof("==> Checking resource %s breaking change...", resourceName)

			// Get file content from git for both versions
			oldContent, err := getFileContentFromGit(file.NewName, "HEAD^")
			if err != nil {
				log.Warningf("Cannot get old version of %s: %v", file.NewName, err)
				continue
			}
			newContent, err := getFileContentFromGit(file.NewName, "HEAD")
			if err != nil {
				log.Warningf("Cannot get new version of %s: %v", file.NewName, err)
				continue
			}

			// Parse schemas using AST
			oldAttrs := ParseSchemaFromAST(oldContent)
			newAttrs := ParseSchemaFromAST(newContent)
			schemaBreaking := IsBreakingChange(oldAttrs, newAttrs)

			// Check retry error codes changes
			oldRetryCodes := make(map[string]map[string]struct{})
			newRetryCodes := make(map[string]map[string]struct{})
			for _, hunk := range file.Hunks {
				if hunk != nil {
					ParseRetryErrorCodesFromHunk(hunk, oldRetryCodes, newRetryCodes)
				}
			}
			retryBreaking := IsRetryCodeBreaking(oldRetryCodes, newRetryCodes)

			if !schemaBreaking && !retryBreaking {
				log.Infof("--- PASS")
			} else {
				log.Errorf("--- FAIL")
				exitCode = 1
			}
		}
	}

	os.Exit(exitCode)
}

func IsBreakingChange(oldAttrs, newAttrs map[string]map[string]interface{}) (res bool) {
	// First check for removed attributes
	for filedName, oldAttr := range oldAttrs {
		// Check if attribute was deleted
		newAttr, attributeExists := newAttrs[filedName]
		if !attributeExists {
			res = true
			log.Errorf("[Breaking Change]: Attribute '%v' should not been removed!", filedName)
			continue
		}

		// Optional -> Required
		oldOptional, oldHasOptional := oldAttr["Optional"]
		_, oldHasRequired := oldAttr["Required"]
		newRequired, newHasRequired := newAttr["Required"]

		// Check if field changed from optional to required
		isOldOptional := (oldHasOptional && oldOptional.(bool)) || (!oldHasOptional && !oldHasRequired)
		isNewRequired := newHasRequired && newRequired.(bool)

		if isOldOptional && isNewRequired {
			res = true
			log.Errorf("[Breaking Change]: '%v' should not been changed from optional to required!", filedName)
		}

		// Check if field changed from non-required to required (when neither was explicitly optional)
		if !oldHasRequired && newHasRequired && newRequired.(bool) {
			// Only if it wasn't explicitly optional before
			if !oldHasOptional || !oldOptional.(bool) {
				res = true
				log.Errorf("[Breaking Change]: '%v' should not been changed to required!", filedName)
			}
		}

		// Type changed
		typPrev, exist1 := oldAttr["Type"]
		typCurr, exist2 := newAttr["Type"]

		if exist1 && exist2 && typPrev != typCurr {
			res = true
			log.Errorf("[Breaking Change]: '%v' type should not been changed from %v to %v!", filedName, typPrev, typCurr)
		}

		// Non-ForceNew -> ForceNew
		_, exist1 = oldAttr["ForceNew"]
		_, exist2 = newAttr["ForceNew"]
		if !exist1 && exist2 {
			res = true
			log.Errorf("[Breaking Change]: '%v' should not been changed to ForceNew!", filedName)
		}

		// Type string/int: valid values
		validateValuesOld, exist1 := oldAttr["ValidateFuncValues"]
		validateValuesNew, exist2 := newAttr["ValidateFuncValues"]
		if exist1 {
			if !exist2 {
				log.Warningf("[Warning]: '%v' ValidateFunc should not been removed!", filedName)
			} else {
				for key, _ := range validateValuesOld.(map[string]struct{}) {
					if _, ok := validateValuesNew.(map[string]struct{})[key]; !ok {
						res = true
						log.Errorf("[Breaking Change]: '%v' valid value %s should not been removed!", filedName, key)
					}
				}
			}
		}
	}

	// Check for newly added required attributes (Breaking Change)
	for fieldName, newAttr := range newAttrs {
		_, existedBefore := oldAttrs[fieldName]
		if !existedBefore {
			// This is a new field
			newRequired, hasRequired := newAttr["Required"]
			if hasRequired && newRequired.(bool) {
				res = true
				log.Errorf("[Breaking Change]: New required attribute '%v' should not been added! New fields must be Optional.", fieldName)
			}
		}
	}

	return
}

// getFileContentFromGit retrieves file content from git at specified revision
func getFileContentFromGit(filePath, revision string) (string, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", revision, filePath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// ParseSchemaFromAST uses Go AST to parse schema definition from source code
func ParseSchemaFromAST(source string) map[string]map[string]interface{} {
	attributeMap := make(map[string]map[string]interface{})

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, 0)
	if err != nil {
		log.Warningf("Failed to parse file: %v", err)
		return attributeMap
	}

	// Find the resource function and extract Schema
	ast.Inspect(file, func(n ast.Node) bool {
		// Look for composite literals (map[string]*schema.Schema{...})
		if compLit, ok := n.(*ast.CompositeLit); ok {
			// Check if this is a Schema map
			if isSchemaMap(compLit) {
				extractSchemaFields(compLit, attributeMap)
			}
		}
		return true
	})

	return attributeMap
}

// isSchemaMap checks if a composite literal is a schema.Schema map
func isSchemaMap(compLit *ast.CompositeLit) bool {
	if mapType, ok := compLit.Type.(*ast.MapType); ok {
		if starExpr, ok := mapType.Value.(*ast.StarExpr); ok {
			if selExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
				if ident, ok := selExpr.X.(*ast.Ident); ok {
					return ident.Name == "schema" && selExpr.Sel.Name == "Schema"
				}
			}
		}
	}
	return false
}

// extractSchemaFields extracts field information from schema map
func extractSchemaFields(compLit *ast.CompositeLit, attributeMap map[string]map[string]interface{}) {
	for _, elt := range compLit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			// Get field name
			var fieldName string
			if basicLit, ok := kv.Key.(*ast.BasicLit); ok {
				fieldName = strings.Trim(basicLit.Value, `"`)
			}

			if fieldName == "" {
				continue
			}

			// Extract field properties
			fieldProps := make(map[string]interface{})
			fieldProps["Name"] = fieldName

			// Parse field definition (should be another composite literal)
			if fieldComp, ok := kv.Value.(*ast.CompositeLit); ok {
				extractFieldProperties(fieldComp, fieldProps)
			} else if unary, ok := kv.Value.(*ast.UnaryExpr); ok {
				// Handle &schema.Schema{...}
				if fieldComp, ok := unary.X.(*ast.CompositeLit); ok {
					extractFieldProperties(fieldComp, fieldProps)
				}
			}

			attributeMap[fieldName] = fieldProps
		}
	}
}

// extractFieldProperties extracts Type, Optional, Required, ForceNew etc from field definition
func extractFieldProperties(compLit *ast.CompositeLit, props map[string]interface{}) {
	for _, elt := range compLit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			var propName string
			if ident, ok := kv.Key.(*ast.Ident); ok {
				propName = ident.Name
			}

			switch propName {
			case "Type":
				// Extract schema.TypeXxx
				if selExpr, ok := kv.Value.(*ast.SelectorExpr); ok {
					if xIdent, ok := selExpr.X.(*ast.Ident); ok && xIdent.Name == "schema" {
						props["Type"] = selExpr.Sel.Name // TypeString, TypeInt, etc.
					}
				}
			case "Optional", "Required", "ForceNew", "Computed":
				// Extract boolean values
				if ident, ok := kv.Value.(*ast.Ident); ok {
					props[propName] = (ident.Name == "true")
				}
			case "ValidateFunc":
				// Extract validation function (simplified for now)
				props[propName] = "present"
			}
		}
	}
}

// ParseRetryErrorCodesFromHunk extracts retry error codes from IsExpectedErrors calls in diff hunks
// Format: IsExpectedErrors(err, []string{"ErrorCode1", "ErrorCode2"})
// Map structure: map[apiContext]map[errorCode]struct{}
func ParseRetryErrorCodesFromHunk(hunk *diffparser.DiffHunk, oldRetryCodes, newRetryCodes map[string]map[string]struct{}) {
	// Regex to match: IsExpectedErrors(err, []string{"code1", "code2", ...})
	expectedErrorsRegex := regexp.MustCompile(`IsExpectedErrors\([^,]+,\s*\[\]string\{([^}]*)\}`)
	// Regex to extract action name from current or previous lines
	actionRegex := regexp.MustCompile(`action\s*:?=\s*"([^"]*)"`)
	// Regex to extract function name context (both regular functions and methods)
	funcRegex := regexp.MustCompile(`func\s+(?:\([^)]*\)\s*)?(\w+)`)

	// Collect all lines first to enable better context searching
	lines := hunk.WholeRange.Lines
	lineContents := make([]string, len(lines))
	for i, line := range lines {
		lineContents[i] = strings.TrimSpace(line.Content)
	}

	for i, line := range lines {
		content := strings.TrimSpace(line.Content)
		if content == "" {
			continue
		}

		// Look for IsExpectedErrors calls
		if expectedErrorsMatch := expectedErrorsRegex.FindStringSubmatch(content); expectedErrorsMatch != nil && len(expectedErrorsMatch) > 1 {
			errorCodesStr := expectedErrorsMatch[1]
			// Parse individual error codes: "Code1", "Code2", ...
			codeRegex := regexp.MustCompile(`"([^"]+)"`)
			codeMatches := codeRegex.FindAllStringSubmatch(errorCodesStr, -1)

			if len(codeMatches) > 0 {
				// Find the best context for this IsExpectedErrors call
				context := findContextForErrorCodes(lineContents, i, actionRegex, funcRegex)

				log.Debugf("Line %d: Found IsExpectedErrors with context '%s'", line.Number, context)

				// Determine which map to update based on line mode
				var targetMap map[string]map[string]struct{}
				switch line.Mode {
				case diffparser.REMOVED:
					targetMap = oldRetryCodes
					log.Debugf("Found REMOVED retry codes in %s: %v", context, codeMatches)
				case diffparser.ADDED:
					targetMap = newRetryCodes
					log.Debugf("Found ADDED retry codes in %s: %v", context, codeMatches)
				default:
					// For unchanged lines, add to both maps
					targetMap = oldRetryCodes
					// Also add to newRetryCodes
					if _, exist := newRetryCodes[context]; !exist {
						newRetryCodes[context] = make(map[string]struct{})
					}
					for _, match := range codeMatches {
						if len(match) > 1 {
							newRetryCodes[context][match[1]] = struct{}{}
						}
					}
				}

				if targetMap != nil {
					if _, exist := targetMap[context]; !exist {
						targetMap[context] = make(map[string]struct{})
					}
					for _, match := range codeMatches {
						if len(match) > 1 {
							targetMap[context][match[1]] = struct{}{}
						}
					}
				}
			}
		}
	}
}

// findContextForErrorCodes searches for the best context (action or function name) for an IsExpectedErrors call
// It uses a multi-strategy approach with different search ranges
func findContextForErrorCodes(lines []string, currentIndex int, actionRegex, funcRegex *regexp.Regexp) string {
	// Strategy 1: Search backwards within the same function scope (up to 5000 lines for extreme cases)
	maxBackwardSearch := 5000
	if currentIndex < maxBackwardSearch {
		maxBackwardSearch = currentIndex
	}

	// Look for action definition in reverse order
	for i := currentIndex - 1; i >= currentIndex-maxBackwardSearch; i-- {
		if actionMatch := actionRegex.FindStringSubmatch(lines[i]); actionMatch != nil && len(actionMatch) > 1 {
			log.Debugf("Found action context '%s' at distance %d from IsExpectedErrors", actionMatch[1], currentIndex-i)
			return actionMatch[1]
		}

		// If we encounter another function definition, stop searching
		if funcRegex.MatchString(lines[i]) && i < currentIndex-5 {
			log.Debugf("Hit function boundary at line %d, stopping context search", i)
			break
		}
	}

	// Strategy 2: Look for function name context
	for i := currentIndex - 1; i >= 0 && i >= currentIndex-20; i-- {
		if funcMatch := funcRegex.FindStringSubmatch(lines[i]); funcMatch != nil && len(funcMatch) > 1 {
			funcName := funcMatch[1]
			log.Debugf("Found function context '%s' at distance %d from IsExpectedErrors", funcName, currentIndex-i)
			return funcName
		}
	}

	// Strategy 3: Search forward for action definition (in case it's defined after IsExpectedErrors, rare but possible)
	maxForwardSearch := 10
	if currentIndex+maxForwardSearch >= len(lines) {
		maxForwardSearch = len(lines) - currentIndex - 1
	}

	for i := currentIndex + 1; i <= currentIndex+maxForwardSearch; i++ {
		if actionMatch := actionRegex.FindStringSubmatch(lines[i]); actionMatch != nil && len(actionMatch) > 1 {
			log.Debugf("Found forward action context '%s' at distance %d from IsExpectedErrors", actionMatch[1], i-currentIndex)
			return actionMatch[1]
		}
	}

	// Strategy 4: Use line number as fallback
	fallbackContext := fmt.Sprintf("line_%d", currentIndex+1)
	log.Debugf("Using fallback context '%s' for IsExpectedErrors", fallbackContext)
	return fallbackContext
}

// IsRetryCodeBreaking checks if retry error codes have been reduced
func IsRetryCodeBreaking(oldRetryCodes, newRetryCodes map[string]map[string]struct{}) bool {
	res := false

	// Log all found retry codes for debugging
	log.Debugf("Old retry codes found: %d contexts", len(oldRetryCodes))
	for context, codes := range oldRetryCodes {
		codeList := make([]string, 0, len(codes))
		for code := range codes {
			codeList = append(codeList, code)
		}
		log.Debugf("  %s: %v", context, codeList)
	}

	log.Debugf("New retry codes found: %d contexts", len(newRetryCodes))
	for context, codes := range newRetryCodes {
		codeList := make([]string, 0, len(codes))
		for code := range codes {
			codeList = append(codeList, code)
		}
		log.Debugf("  %s: %v", context, codeList)
	}

	// Check for removed contexts (entire IsExpectedErrors calls removed)
	for context, oldCodes := range oldRetryCodes {
		newCodes, exist := newRetryCodes[context]

		// If the IsExpectedErrors call was completely removed
		if !exist {
			if len(oldCodes) > 0 {
				res = true
				log.Errorf("[Breaking Change]: Retry error codes for '%s' have been completely removed!", context)
				for code := range oldCodes {
					log.Errorf("  - Removed error code: '%s'", code)
				}
			}
			continue
		}

		// Check if any error codes were removed within the same context
		for oldCode := range oldCodes {
			if _, stillExists := newCodes[oldCode]; !stillExists {
				res = true
				log.Errorf("[Breaking Change]: Retry error code '%s' for '%s' should not be removed!", oldCode, context)
			}
		}
	}

	// Log summary
	if res {
		log.Errorf("Retry error code breaking changes detected!")
	} else {
		log.Infof("No retry error code breaking changes detected.")
	}

	return res
}
