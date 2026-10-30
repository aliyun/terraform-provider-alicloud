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

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
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
	fileNames    = flag.String("fileNames", "", "the files to check diff")
	resourceName = flag.String("resource", "", "directly specify a resource name to check (e.g., alicloud_ecs_instance)")
)

func main() {
	exitCode := 0
	flag.Parse()

	// Mode 1: Direct resource name check (for local use)
	if resourceName != nil && len(*resourceName) > 0 {
		log.Infof("==> Checking resource: %s", *resourceName)
		if !checkSingleResource(*resourceName) {
			exitCode = 1
		}
		os.Exit(exitCode)
	}

	// Mode 2: Diff file check (for CI/CD)
	if fileNames != nil && len(*fileNames) == 0 {
		log.Warningf("the diff file is empty")
		return
	}
	byt, _ := ioutil.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	fileRegex := regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*.go")
	fileTestRegex := regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*_test.go")
	for _, file := range diff.Files {
		fmt.Println()
		if fileRegex.MatchString(file.NewName) {
			if fileTestRegex.MatchString(file.NewName) {
				continue
			}
			base := strings.TrimSuffix(strings.Split(file.NewName, "/")[1], ".go")
			// The pre-split flat layout also holds data sources, and both
			// kinds carry a schema the AST parsers can compare.
			kind, legacyStem := "resource", "resource_"
			resName := strings.TrimPrefix(base, "resource_")
			if strings.HasPrefix(base, "data_source_") {
				kind, legacyStem = "data source", "data_source_"
				resName = strings.TrimPrefix(base, "data_source_")
			}
			log.Infof("==> Checking %s %s breaking change...", kind, resName)

			if !checkBreakingChange(file.NewName, resName, legacyStem, "HEAD^", "HEAD") {
				exitCode = 1
			}
		} else if product, stem, isTest, ok := splitlayout.IsSplitSource(file.NewName); ok {
			if isTest {
				continue
			}
			resolved, err := splitlayout.Resolve(product, stem)
			if err != nil {
				log.Warnf("fail to resolve split-layout file %s. Error: %s", file.NewName, err)
				continue
			}
			// SDKv2 resources and data sources alike are schema-bearing and
			// checked; the legacy stem locates a moved file's flat-layout
			// predecessor.
			for _, r := range resolved {
				legacyStem := "data_source_"
				if r.IsResource() {
					legacyStem = "resource_"
				}
				log.Infof("==> Checking %s %s breaking change...", r.Kind, r.TypeName)
				if !checkBreakingChange(r.ImplPath, r.TypeName, legacyStem, "HEAD^", "HEAD") {
					exitCode = 1
				}
			}
			// Framework-native registrations — managed resources, data
			// sources, ephemerals, list resources and actions — are checked
			// the same way through the framework AST parser.
			fwResolved, err := splitlayout.FrameworkResolve(product, stem)
			if err != nil {
				log.Warnf("fail to resolve framework registrations of %s. Error: %s", file.NewName, err)
				continue
			}
			for _, r := range fwResolved {
				legacyStem := ""
				switch r.Kind {
				case splitlayout.KindResource:
					legacyStem = "resource_"
				case splitlayout.KindDataSource:
					legacyStem = "data_source_"
				}
				log.Infof("==> Checking %s %s breaking change...", r.Kind, r.TypeName)
				if !checkBreakingChange(r.ImplPath, r.TypeName, legacyStem, "HEAD^", "HEAD") {
					exitCode = 1
				}
			}
			if len(resolved) == 0 && len(fwResolved) == 0 {
				log.Warnf("==> %s registers no resource, data source, ephemeral, action or list resource", file.NewName)
			}
		}
	}

	os.Exit(exitCode)
}

// checkSingleResource checks a single registered entity by its full type
// name (e.g., "alicloud_ecs_instance", "alicloud_ecs_images" or
// "alicloud_kms_plaintext"), whatever its kind: resource, data source,
// ephemeral, action or list resource.
func checkSingleResource(fullName string) bool {
	// Remove alicloud_ prefix if present
	resName := strings.TrimPrefix(fullName, "alicloud_")
	typeName := "alicloud_" + resName

	// Flat layout first: a type name can exist as a resource, a data source
	// or both, and both are checked.
	passed, checked := true, false
	for _, legacyStem := range []string{"resource_", "data_source_"} {
		filePath := "alicloud/" + legacyStem + typeName + ".go"
		if _, err := os.Stat(filePath); err != nil {
			continue
		}
		checked = true
		kind := "resource"
		if legacyStem == "data_source_" {
			kind = "data source"
		}
		log.Infof("==> Checking %s %s breaking change...", kind, typeName)
		passed = checkBreakingChange(filePath, typeName, legacyStem, "HEAD^", "HEAD") && passed
	}
	if checked {
		return passed
	}

	// Migrated entities live under alicloud/service/<product>/; fall back to
	// the registrations when the flat-layout files are gone.
	if r := splitlayout.FindByTypeName(typeName); r != nil {
		legacyStem := "data_source_"
		if r.IsResource() {
			legacyStem = "resource_"
		}
		log.Infof("==> Checking %s %s breaking change...", r.Kind, typeName)
		return checkBreakingChange(r.ImplPath, typeName, legacyStem, "HEAD^", "HEAD")
	}
	// Framework-native entities are checked through the same AST parsers.
	if r := splitlayout.FrameworkFindByTypeName(typeName); r != nil {
		legacyStem := ""
		switch r.Kind {
		case splitlayout.KindResource:
			legacyStem = "resource_"
		case splitlayout.KindDataSource:
			legacyStem = "data_source_"
		}
		log.Infof("==> Checking %s %s breaking change...", r.Kind, typeName)
		return checkBreakingChange(r.ImplPath, typeName, legacyStem, "HEAD^", "HEAD")
	}
	log.Errorf("No implementation file found for %s in the flat or split layouts", typeName)
	return false
}

// checkBreakingChange checks breaking changes for one implementation file
// between two git revisions. fullTypeName (e.g. "alicloud_cloudauth_face_config")
// together with legacyStem ("resource_" or "data_source_") locates the
// pre-split flat-layout file: a split-layout file that is new at oldRevision
// is usually a moved entity, not a new one, and must be compared against its
// legacy path. An empty legacyStem means the kind has no flat-layout history
// (ephemerals, actions and list resources were born in the split layout).
func checkBreakingChange(filePath, fullTypeName, legacyStem, oldRevision, newRevision string) bool {
	// Get file content from git for both versions
	oldContent, err := getFileContentFromGit(filePath, oldRevision)
	if err != nil {
		// Check if this is a new file (doesn't exist in old revision)
		if isFileNew(filePath, oldRevision) {
			// A split-layout migration leaves the old content at the
			// flat-layout path; compare against that before treating the
			// entity as new.
			if fullTypeName != "" && legacyStem != "" {
				legacyPath := "alicloud/" + legacyStem + fullTypeName + ".go"
				if !isFileNew(legacyPath, oldRevision) {
					legacyContent, lerr := getFileContentFromGit(legacyPath, oldRevision)
					if lerr != nil {
						log.Warningf("Cannot get old version of %s: %v", legacyPath, lerr)
						return false
					}
					newContent, nerr := getFileContentFromGit(filePath, newRevision)
					if nerr != nil {
						log.Warningf("Cannot get new version of %s: %v", filePath, nerr)
						return false
					}
					return compareBreakingChange(legacyContent, newContent)
				}
			}
			log.Infof("This is a new file, checking for new entity requirements...")
			// For new entities, we need to check if any fields are required
			newContent, err := getFileContentFromGit(filePath, newRevision)
			if err != nil {
				log.Warningf("Cannot get new version of %s: %v", filePath, err)
				return false
			}
			return checkNewEntityRequiredFields(newContent)
		}
		log.Warningf("Cannot get old version of %s: %v", filePath, err)
		return false
	}
	newContent, err := getFileContentFromGit(filePath, newRevision)
	if err != nil {
		log.Warningf("Cannot get new version of %s: %v", filePath, err)
		return false
	}

	return compareBreakingChange(oldContent, newContent)
}

// compareBreakingChange compares two resource file versions for schema and
// retry-code breaking changes.
func compareBreakingChange(oldContent, newContent string) bool {
	// Parse schemas using AST; the SDKv2 and framework shapes are disjoint,
	// so parsing both covers migrated files on either side of the split.
	oldAttrs := parseSchemas(oldContent)
	newAttrs := parseSchemas(newContent)
	schemaBreaking := IsBreakingChange(oldAttrs, newAttrs)

	// Check retry error codes changes - parse from full file content
	oldRetryCodes := ParseRetryErrorCodesFromContent(oldContent)
	newRetryCodes := ParseRetryErrorCodesFromContent(newContent)
	retryBreaking := IsRetryCodeBreaking(oldRetryCodes, newRetryCodes)

	if !schemaBreaking && !retryBreaking {
		log.Infof("--- PASS")
		return true
	} else {
		log.Errorf("--- FAIL")
		return false
	}
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
				// A Required subfield under a newly added Optional parent
				// block is not a breaking change: configurations written
				// before the block existed do not reference it, so they
				// remain valid. Only waive when the immediate parent block is
				// itself new and Optional; a Required subfield added under a
				// block that already existed still invalidates old configs
				// that use the block, and a top-level new Required is always
				// breaking.
				if isNewOptionalNestedRequired(fieldName, oldAttrs, newAttrs) {
					log.Infof("[Info]: New required attribute '%v' is nested under a new optional block; not a breaking change for existing configurations.", fieldName)
					continue
				}
				res = true
				log.Errorf("[Breaking Change]: New required attribute '%v' should not been added! New fields must be Optional.", fieldName)
			}
		}
	}

	return
}

// isNewOptionalNestedRequired reports whether fieldName is a nested subfield
// (its key carries a parent-path separator) whose immediate parent block is
// newly added in this change and declared Optional. Top-level fields and
// subfields under a parent that already existed are never waived: the former
// is a real new top-level required attribute, and the latter invalidates
// configurations that already use the block. The parent is not modeled only
// when the path is deeper than the schemas captured here; in that case the
// violation is left to surface rather than waived on incomplete evidence.
func isNewOptionalNestedRequired(fieldName string, oldAttrs, newAttrs map[string]map[string]interface{}) bool {
	idx := strings.LastIndex(fieldName, ".")
	if idx < 0 {
		return false // top-level new required: always breaking
	}
	parent := fieldName[:idx]
	if _, existedBefore := oldAttrs[parent]; existedBefore {
		return false // parent block already existed: required subfield breaks old configs
	}
	parentAttr, ok := newAttrs[parent]
	if !ok {
		return false // parent not modeled: do not waive on incomplete evidence
	}
	optional, hasOptional := parentAttr["Optional"]
	return hasOptional && optional.(bool)
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

// isFileNew checks if a file is new (doesn't exist in the specified revision)
func isFileNew(filePath, revision string) bool {
	cmd := exec.Command("git", "cat-file", "-e", fmt.Sprintf("%s:%s", revision, filePath))
	err := cmd.Run()
	return err != nil // If error, file doesn't exist in that revision
}

// checkNewEntityRequiredFields checks if a new resource, data source,
// ephemeral or action has required fields
// For new entities, having required fields is acceptable, we just log them
func checkNewEntityRequiredFields(content string) bool {
	attrs := parseSchemas(content)
	hasRequiredFields := false

	for fieldName, attr := range attrs {
		if required, ok := attr["Required"]; ok && required.(bool) {
			hasRequiredFields = true
			log.Infof("  New entity has required field: '%s'", fieldName)
		}
	}

	if hasRequiredFields {
		log.Infof("Note: New entities can have required fields. This is acceptable for new entities.")
	}
	log.Infof("--- PASS (new entity)")
	return true
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
				// Extract top-level fields with an empty path prefix and stop
				// descending: nested schema maps live inside Elem and are
				// recursed by extractSchemaFields with their parent path, so
				// they are never flattened back onto bare field names and a
				// Required subfield is never mistaken for a top-level field.
				extractSchemaFields(compLit, attributeMap, "")
				return false
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

// extractSchemaFields extracts field information from a schema map. prefix is
// the dotted path of the enclosing block ("" at the resource root), so a
// subfield declared inside a nested block is stored under
// "parent_block.sub_field" rather than the bare "sub_field": that keeps
// same-named subfields under different parent blocks from overwriting each
// other and lets IsBreakingChange judge a Required subfield against its real
// parent context instead of as a phantom top-level field.
func extractSchemaFields(compLit *ast.CompositeLit, attributeMap map[string]map[string]interface{}, prefix string) {
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

			fullKey := prefix + fieldName

			// Extract field properties
			fieldProps := make(map[string]interface{})
			fieldProps["Name"] = fieldName

			// Parse field definition (should be another composite literal)
			var schemaComp *ast.CompositeLit
			if fieldComp, ok := kv.Value.(*ast.CompositeLit); ok {
				schemaComp = fieldComp
			} else if unary, ok := kv.Value.(*ast.UnaryExpr); ok {
				// Handle &schema.Schema{...}
				schemaComp, _ = unary.X.(*ast.CompositeLit)
			}
			if schemaComp != nil {
				extractFieldProperties(schemaComp, fieldProps)
			}

			attributeMap[fullKey] = fieldProps

			// Recurse into a nested schema map declared under Elem (e.g.
			// &schema.Resource{Schema: map[string]*schema.Schema{...}}) so its
			// subfields keep the parent block path instead of being flattened.
			if schemaComp != nil {
				if nestedMap := nestedSchemaMap(schemaComp); nestedMap != nil {
					extractSchemaFields(nestedMap, attributeMap, fullKey+".")
				}
			}
		}
	}
}

// nestedSchemaMap returns the nested map[string]*schema.Schema declared under
// a field's Elem, if any. A primitive Elem such as &schema.Schema{Type: ...}
// yields no nested map; only &schema.Resource{Schema: map[...]} (and the rare
// direct schema-map Elem) does, which is what makes a block a block.
func nestedSchemaMap(schemaComp *ast.CompositeLit) *ast.CompositeLit {
	for _, elt := range schemaComp.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != "Elem" {
			continue
		}
		var resComp *ast.CompositeLit
		switch e := kv.Value.(type) {
		case *ast.UnaryExpr:
			resComp, _ = e.X.(*ast.CompositeLit)
		case *ast.CompositeLit:
			resComp = e
		}
		if resComp == nil {
			continue
		}
		for _, rElt := range resComp.Elts {
			rKv, ok := rElt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			rKey, ok := rKv.Key.(*ast.Ident)
			if !ok || rKey.Name != "Schema" {
				continue
			}
			if nestedMap, ok := rKv.Value.(*ast.CompositeLit); ok && isSchemaMap(nestedMap) {
				return nestedMap
			}
		}
	}
	return nil
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

// parseSchemas extracts schema attributes from SDKv2 and plugin-framework
// shapes alike. The two shapes are disjoint (map[string]*schema.Schema vs
// map[string]schema.Attribute under any import alias), so merging the two
// parsers' output covers every registered entity kind and compares a
// migrated file cleanly against its pre-split flat-layout content.
func parseSchemas(source string) map[string]map[string]interface{} {
	attributeMap := ParseSchemaFromAST(source)
	for name, props := range ParseFrameworkSchemaFromAST(source) {
		attributeMap[name] = props
	}
	return attributeMap
}

// ParseFrameworkSchemaFromAST parses the plugin-framework schema shape:
// attribute maps declared as map[string]schema.Attribute under any import
// alias (the resource, ephemeral and action schema packages all share the
// Attribute type name), including the nested maps inside blocks. Attribute
// type names normalize to their SDKv2 equivalents so a migrated entity
// compares against its flat-layout predecessor without phantom type breaks.
func ParseFrameworkSchemaFromAST(source string) map[string]map[string]interface{} {
	attributeMap := make(map[string]map[string]interface{})

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, 0)
	if err != nil {
		log.Warningf("Failed to parse file: %v", err)
		return attributeMap
	}

	ast.Inspect(file, func(n ast.Node) bool {
		compLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		if isFrameworkAttributeMap(compLit) {
			extractFrameworkFields(compLit, attributeMap)
		}
		return true
	})

	return attributeMap
}

// isFrameworkAttributeMap checks if a composite literal is a framework
// attribute map (map[string]schema.Attribute). The package identifier is not
// inspected: the resource, ephemeral and action schema packages import under
// different aliases but share the Attribute type name.
func isFrameworkAttributeMap(compLit *ast.CompositeLit) bool {
	mapType, ok := compLit.Type.(*ast.MapType)
	if !ok {
		return false
	}
	selExpr, ok := mapType.Value.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	return selExpr.Sel.Name == "Attribute"
}

// extractFrameworkFields extracts attribute information from a framework
// attribute map.
func extractFrameworkFields(compLit *ast.CompositeLit, attributeMap map[string]map[string]interface{}) {
	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		basicLit, ok := kv.Key.(*ast.BasicLit)
		if !ok {
			continue
		}
		fieldName := strings.Trim(basicLit.Value, `"`)
		if fieldName == "" {
			continue
		}

		fieldProps := make(map[string]interface{})
		fieldProps["Name"] = fieldName

		if fieldComp, ok := kv.Value.(*ast.CompositeLit); ok {
			extractFrameworkFieldProperties(fieldComp, fieldProps)
		}

		attributeMap[fieldName] = fieldProps
	}
}

// frameworkTypeEquivalents maps plugin-framework attribute and block type
// names to the SDKv2 schema.TypeXxx constants they play the role of.
var frameworkTypeEquivalents = map[string]string{
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
}

// extractFrameworkFieldProperties extracts Type, Optional, Required,
// Computed etc from an attribute or block definition.
func extractFrameworkFieldProperties(compLit *ast.CompositeLit, props map[string]interface{}) {
	if selExpr, ok := compLit.Type.(*ast.SelectorExpr); ok {
		if equivalent, ok := frameworkTypeEquivalents[selExpr.Sel.Name]; ok {
			props["Type"] = equivalent
		} else {
			props["Type"] = selExpr.Sel.Name
		}
	}
	for _, elt := range compLit.Elts {
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

// ParseRetryErrorCodesFromContent extracts retry error codes from full file content
// Formats: IsExpectedErrors(err, []string{"Code"}) and rdsErrorHasCode(err, "Code").
// Map structure: map[apiName]map[errorCode]struct{}
func ParseRetryErrorCodesFromContent(content string) map[string]map[string]struct{} {
	retryCodesMap := make(map[string]map[string]struct{})
	// Both helpers may receive a named error such as deleteErr.
	expectedErrorsRegex := regexp.MustCompile(`(?:IsExpectedErrors|rdsErrorHasCode)\(\s*[a-zA-Z_][a-zA-Z0-9_]*,\s*(?:\[\]string\{)?([^})]*)`)
	// Regex to find action variable: action := "ApiName" or action = "ApiName"
	actionRegex := regexp.MustCompile(`action\s*:?=\s*"([^"]*)"`)

	lines := strings.Split(content, "\n")
	currentAction := ""

	for _, line := range lines {
		// Look for action definition
		if actionMatches := actionRegex.FindStringSubmatch(line); actionMatches != nil && len(actionMatches) > 1 {
			currentAction = actionMatches[1]
		}

		// Look for a supported error-code classifier
		if expectedErrorsMatches := expectedErrorsRegex.FindStringSubmatch(line); expectedErrorsMatches != nil && len(expectedErrorsMatches) > 1 {
			if currentAction == "" {
				// Skip if no action found in current scope
				continue
			}

			errorCodesStr := expectedErrorsMatches[1]

			// Extract individual error codes
			codeRegex := regexp.MustCompile(`"([^"]+)"`)
			codeMatches := codeRegex.FindAllStringSubmatch(errorCodesStr, -1)
			if len(codeMatches) == 0 {
				continue
			}

			if _, exists := retryCodesMap[currentAction]; !exists {
				retryCodesMap[currentAction] = make(map[string]struct{})
			}

			for _, match := range codeMatches {
				if len(match) > 1 {
					errorCode := match[1]
					retryCodesMap[currentAction][errorCode] = struct{}{}
				}
			}
		}
	}

	return retryCodesMap
}

// IsRetryCodeBreaking checks if retry error codes have been reduced for each API
func IsRetryCodeBreaking(oldRetryCodes, newRetryCodes map[string]map[string]struct{}) bool {
	res := false

	// Check each API's retry codes
	for apiName, oldCodes := range oldRetryCodes {
		newCodes, exist := newRetryCodes[apiName]
		if !exist {
			// The entire IsExpectedErrors call was removed for this API
			res = true
			log.Errorf("[Breaking Change]: Retry error codes for API '%s' have been completely removed!", apiName)
			for oldCode := range oldCodes {
				log.Errorf("  - Removed error code: '%s'", oldCode)
			}
			continue
		}

		// Check if any error code was removed
		for oldCode := range oldCodes {
			if _, stillExists := newCodes[oldCode]; !stillExists {
				res = true
				log.Errorf("[Breaking Change]: Retry error code '%s' for API '%s' should not been removed!", oldCode, apiName)
			}
		}
	}

	return res
}
