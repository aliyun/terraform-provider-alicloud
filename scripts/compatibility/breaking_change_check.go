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
					ParseRetryErrorCodes(hunk.OrigRange, hunk.OrigRange.Length, oldRetryCodes)
					ParseRetryErrorCodes(hunk.NewRange, hunk.NewRange.Length, newRetryCodes)
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

// ParseRetryErrorCodes extracts retry error codes from IsExpectedErrors calls
// Format: IsExpectedErrors(err, []string{"ErrorCode1", "ErrorCode2"})
// Map structure: map[apiContext]map[errorCode]struct{}
func ParseRetryErrorCodes(hunk diffparser.DiffRange, length int, retryCodesMap map[string]map[string]struct{}) {
	// Regex to match: IsExpectedErrors(err, []string{"code1", "code2", ...})
	expectedErrorsRegex := regexp.MustCompile(`IsExpectedErrors\(err,\s*\[\]string\{([^}]*)\}`)
	// Regex to extract action name from previous lines
	actionRegex := regexp.MustCompile(`action\s*:?=\s*"([^"]*)"`)

	currentAction := ""
	for i := 0; i < length; i++ {
		currentLine := hunk.Lines[i]
		content := currentLine.Content

		// Try to find action definition
		actionMatched := actionRegex.FindStringSubmatch(content)
		if actionMatched != nil && len(actionMatched) > 1 {
			currentAction = actionMatched[1]
		}

		// Try to find IsExpectedErrors call
		expectedErrorsMatched := expectedErrorsRegex.FindStringSubmatch(content)
		if expectedErrorsMatched != nil && len(expectedErrorsMatched) > 1 {
			// Extract error codes from the matched string
			errorCodesStr := expectedErrorsMatched[1]
			// Parse individual error codes: "Code1", "Code2", ...
			codeRegex := regexp.MustCompile(`"([^"]+)"`)
			codeMatches := codeRegex.FindAllStringSubmatch(errorCodesStr, -1)

			if len(codeMatches) > 0 {
				// Use line number + action as context if action is not found
				context := currentAction
				if context == "" {
					context = fmt.Sprintf("line_%d", currentLine.Number)
				}

				if _, exist := retryCodesMap[context]; !exist {
					retryCodesMap[context] = make(map[string]struct{})
				}

				for _, match := range codeMatches {
					if len(match) > 1 {
						retryCodesMap[context][match[1]] = struct{}{}
					}
				}
			}
		}
	}
}

// IsRetryCodeBreaking checks if retry error codes have been reduced
func IsRetryCodeBreaking(oldRetryCodes, newRetryCodes map[string]map[string]struct{}) bool {
	res := false

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

		// Check if any error codes were removed
		for oldCode := range oldCodes {
			if _, stillExists := newCodes[oldCode]; !stillExists {
				res = true
				log.Errorf("[Breaking Change]: Retry error code '%s' for '%s' should not been removed!", oldCode, context)
			}
		}
	}

	return res
}
