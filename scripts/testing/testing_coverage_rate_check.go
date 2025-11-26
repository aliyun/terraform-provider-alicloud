//nolint:all
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	resourceName = flag.String("resource", "", "directly specify a resource name to check (e.g., alicloud_realtime_compute_deployment)")

	resourceFileRegex     = regexp.MustCompile("alicloud/(resource)[0-9a-zA-Z_]*")
	resourceFileTestRegex = regexp.MustCompile("alicloud/(resource)[0-9a-zA-Z_]*_test.go")
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
		log.Infof("the diff file is empty, shipped!")
		return
	}

	byt, _ := os.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	resourceNameMap := make(map[string]struct{})
	for _, file := range diff.Files {
		isNameCorrect = true
		resName := ""
		isResource := true
		fileType := "resource"
		if resourceFileTestRegex.MatchString(file.NewName) {
			resName = strings.TrimSuffix(strings.Split(file.NewName, "/")[1], "_test.go")
		} else if resourceFileRegex.MatchString(file.NewName) {
			resName = strings.TrimSuffix(strings.Split(file.NewName, "/")[1], ".go")
		} else {
			continue
		}

		if strings.HasPrefix(resName, "data_source_") {
			isResource = false
			fileType = "data source"
			resName = strings.TrimPrefix(resName, "data_source_")
		} else {
			resName = strings.TrimPrefix(resName, "resource_")
		}

		// Remove alicloud_ prefix to get the clean resource name
		resName = strings.TrimPrefix(resName, "alicloud_")

		if _, ok := resourceNameMap[resName]; ok {
			continue
		} else {
			resourceNameMap[resName] = struct{}{}
		}

		if !checkResourceByName(resName, isResource, fileType) {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

// checkSingleResource checks a single resource by its full name (e.g., "alicloud_ecs_instance")
func checkSingleResource(fullName string) bool {
	isResource := true
	fileType := "resource"
	resName := fullName

	// Determine if it's a data source or resource
	if strings.HasPrefix(fullName, "data.") {
		isResource = false
		fileType = "data source"
		resName = strings.TrimPrefix(fullName, "data.")
	}

	// Remove alicloud_ prefix if present
	resName = strings.TrimPrefix(resName, "alicloud_")

	return checkResourceByName(resName, isResource, fileType)
}

// checkResourceByName checks a resource by name and type
func checkResourceByName(resName string, isResource bool, fileType string) bool {
	isNameCorrect = true

	log.Infof("==> Getting %s %s attributes...", fileType, resName)
	var resource = &schema.Resource{}
	var ok bool
	if isResource {
		resource, ok = alicloud.Provider().(*schema.Provider).ResourcesMap["alicloud_"+resName]
		if !ok || resource == nil {
			log.Errorf("resource %s is not found in the provider ResourceMap\n\n", resName)
			return false
		}
	} else {
		resource, ok = alicloud.Provider().(*schema.Provider).DataSourcesMap["alicloud_"+resName]
		if !ok || resource == nil {
			log.Errorf("data source %s is not found in the provider DataSourcesMap\n\n", resName)
			return false
		}
	}

	schemaAllSet, schemaMustSet, schemaModifySet, schemaForceNewSet :=
		mapset.NewSet(), mapset.NewSet(), mapset.NewSet(), mapset.NewSet()
	getSchemaAttr(false, resource.Schema, &schemaAllSet, &schemaMustSet, &schemaModifySet, &schemaForceNewSet)

	log.Infof("==> Getting %s %s attributes in test cases...", fileType, resName)
	testMustSet, testModifySet, testIgnoreSet :=
		mapset.NewSet(), mapset.NewSet(), mapset.NewSet()
	filePath := "alicloud/"
	if isResource {
		filePath += "resource_alicloud_" + resName + "_test.go"
	} else {
		filePath += "data_source_alicloud_" + resName + "_test.go"
	}
	check := getTestCaseAttr(filePath, resName, &testMustSet, &testModifySet, &testIgnoreSet)

	// "check" denotes the test code is using a standard template
	if check {
		log.Infof("==> checking %s %s attributes' coverage rate", fileType, resName)
		if checkAttributeSet(resName, fileType, schemaMustSet, testMustSet,
			schemaModifySet, testModifySet, schemaForceNewSet, schemaAllSet, testIgnoreSet) && isNameCorrect {
			log.Infof("--- PASS!\n\n")
			return true
		}
	}

	log.Errorf("--- Failed!\n\n")
	return false
}

// get the schema
func getSchemaAttr(isResource bool, schema map[string]*schema.Schema,
	schemaAllSet, schemaMustSet, schemaModifySet, schemaForceNewSet *mapset.Set) {

	schemaAttributes := make(map[string]SchemaAttribute)

	getSchemaAttributes("", schemaAttributes, schema)

	for key, value := range schemaAttributes {
		// "dry_run" or deperacated
		if key == "dry_run" || value.Deprecated != "" {
			continue
		}
		(*schemaAllSet).Add(key)
		if value.Optional || value.Required {
			(*schemaMustSet).Add(key)
			if !value.ForceNew {
				(*schemaModifySet).Add(key)
			}
		}
		if value.ForceNew {
			(*schemaForceNewSet).Add(key)
		}

	}
}

func getSchemaAttributes(rootName string, schemaAttributes map[string]SchemaAttribute,
	resourceSchema map[string]*schema.Schema) {
	for key, value := range resourceSchema {
		if len(value.Removed) != 0 {
			continue
		}

		if rootName != "" {
			key = rootName + "." + key
		}

		if _, ok := schemaAttributes[key]; !ok {
			schemaAttributes[key] = SchemaAttribute{
				Name:       key,
				Type:       value.Type.String(),
				Optional:   value.Optional,
				Required:   value.Required,
				ForceNew:   value.ForceNew,
				Default:    fmt.Sprint(value.Default),
				Deprecated: value.Deprecated,
			}
		}
		if value.Type == schema.TypeSet || value.Type == schema.TypeList {
			if v, ok := value.Elem.(*schema.Schema); ok {
				vv := schemaAttributes[key]
				vv.ElemType = v.Type.String()
				schemaAttributes[key] = vv
			} else {
				vv := schemaAttributes[key]
				vv.ElemType = "Object"
				schemaAttributes[key] = vv
				getSchemaAttributes(key, schemaAttributes, value.Elem.(*schema.Resource).Schema)
			}
		}
	}
}

type SchemaAttribute struct {
	Name        string
	Type        string
	Optional    bool
	Required    bool
	ForceNew    bool
	Default     string
	ElemType    string
	Deprecated  string
	DocsLineNum int
	Removed     string
}

// get the attribute which have been tested
func getTestCaseAttr(filePath string, resourceName string,
	testMustSet, testModifySet, testIgnoreSet *mapset.Set) bool {
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("fail to open test file %s. Error: %s", filePath, err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	resourceTest := ResourceTest{
		resourceName: filePath,
		funcs:        map[string]FuncTest{},
	}
	line := 0
	funcName := ""
	stepNumber := 0
	configStr := ""
	inConfig := false
	inFunc := false
	ignoreStr := ""
	inIgnore := false
	for scanner.Scan() {
		line += 1
		text := scanner.Text()
		// commented line
		if commentedRegex.MatchString(text) {
			continue
		} else if text == "}" {
			inFunc = false

		} else if normalFuncRegex.MatchString(text) {
			if unitFuncRegex.MatchString(text) {
				continue
			}
			if !standardFuncRegex.MatchString(text) {
				name := text[strings.Index(text, "T"):strings.Index(text, "(")]
				log.Errorf("testcase %s should start with TestAccAliCloud", name)
				isNameCorrect = false
			}
			inFunc = true
			funcName = text[strings.Index(text, "T"):strings.Index(text, "(")]
			stepNumber = 0
			resourceTest.funcs[funcName] = FuncTest{
				stepStr:        map[int]string{},
				stepAttributes: map[int]map[string]interface{}{},
			}
		}
		if inFunc {
			if configRegex.MatchString(text) {
				configStr += text
				inConfig = true
				continue
			}
			if inConfig {
				if checkRegex.MatchString(text) {
					resourceTest.funcs[funcName].stepStr[stepNumber] = configStr
					inConfig = false
					configStr = ""
					stepNumber++
				} else {
					// remove extra spaces and '\'
					text = symbolRegex.ReplaceAllString(text, "")
					if len(text) == 0 || strings.HasPrefix(text, "//") {
						continue
					}
					configStr += text + "\n"
				}
			}
			if ignoreRegex.MatchString(text) {
				inIgnore = true
				ignoreStr += text
				continue
			}
			if inIgnore {
				ignoreStr += text
				if strings.Contains(text, "}") {
					inIgnore = false
					ignoreStr = strings.ReplaceAll(ignoreStr, "\"", "")
					ignoreStr = symbolRegex.ReplaceAllString(ignoreStr, "")
					attrStr := ignoreStr[strings.Index(ignoreStr, "{")+1 : strings.Index(ignoreStr, "}")]
					if len(attrStr) != 0 {
						attrSlice := strings.Split(attrStr, ",")
						for _, v := range attrSlice {
							if len(v) != 0 && v != "dry_run" {
								(*testIgnoreSet).Add(v)
							}
						}
					}
					ignoreStr = ""
				}
			}
		}

	}

	return parseConfig(resourceTest, testMustSet, testModifySet)

}

func parseConfig(resourceTest ResourceTest,
	testMustSet, testModifySet *mapset.Set) (toCheck bool) {
	for funcName, f := range resourceTest.funcs {
		// attribute-value map in a test func
		attributeValueMap := map[string]string{}
		for configIndex := 0; configIndex < len(f.stepStr); configIndex++ {
			configStr := f.stepStr[configIndex]
			// the test code is not using a standard template
			if !strings.Contains(configStr, "{") {
				log.Infof("the test case in [%s] does not use a standard template, need manual check", resourceTest.resourceName)
				return false
			}
			configStr = configStr[strings.Index(configStr, "{")+2 : strings.LastIndex(configStr, "}")+1]
			if configStr == "{}" {
				continue
			}
			configSlice := strings.Split(configStr, "\n")

			// traverse each line of config
			for i, v := range configSlice {
				valueIndex := strings.Index(v, ":")

				// "xxx:xxx",
				if strings.Contains(v, ":") {
					splitIndex := strings.Index(v, ":")
					beforeCount := strings.Count(v[:splitIndex], "\"")
					if beforeCount%2 != 0 {
						valueIndex = -1
					}
				}

				valueSuffix := string(v[len(v)-1])
				valueStr := v[valueIndex+1 : len(v)-1]
				bracketCount := strings.Count(v, "]") + strings.Count(v, "[")
				if v[valueIndex+1] == '"' && strings.HasSuffix(v, "},") ||
					strings.HasSuffix(v, "],") && bracketCount%2 == 1 {
					valueSuffix = v[len(v)-2:]
					valueStr = v[valueIndex+1 : len(v)-2]
				}

				// Check if this is a nested JSON string (contains multi-level escaped quotes like \\\")
				// This indicates the value is itself a JSON string with internal JSON structure
				// We need to skip all transformations to preserve its original structure
				isNestedJSONString := strings.HasPrefix(valueStr, "\"") &&
					strings.Contains(valueStr, `\\\"`) &&
					(strings.Contains(valueStr, "{") || strings.Contains(valueStr, "["))

				if !isNestedJSONString {
					// "xxx"+xx, "xxx"+xxx+"xxx, `xxx`+"xxx"+`xxx`
					if strings.Contains(valueStr, "+") {
						v := valueStr
						index := strings.Index(v, "+")
						for index != -1 {
							beforeStr := v[:index]
							beforeStr = strings.TrimSpace(beforeStr)
							if strings.Count(beforeStr, "\"")%2 == 0 ||
								(beforeStr[0] == '`' && beforeStr[len(beforeStr)-1] == '`') {
								v = v[index+1:]
								index = strings.Index(v, "+")
							} else {
								break
							}
						}
						if index == -1 {
							valueStr = strings.ReplaceAll(valueStr, "\"", "")
							valueStr = strings.ReplaceAll(valueStr, "`", "")
							valueStr = "\"" + valueStr + "\""
						}
					}

					// `xxx`
					if strings.HasPrefix(valueStr, "`") && strings.HasSuffix(valueStr, "`") {
						valueStr = "\"" + valueStr[:len(valueStr)-1] + "\"" + valueSuffix
					} else {
						valueStr += valueSuffix
					}

					// value is a map or slice
					for i := 0; i < len(matchMap); i++ {
						k := matchMap[i][0]
						v := matchMap[i][1]
						if strings.Contains(valueStr, k) {
							valueStr = strings.ReplaceAll(valueStr, k, v)
						}
					}

					// value with variable or func
					if variableRegex.MatchString(valueStr) || valueFuncRegex.MatchString(valueStr) {
						valueStr = strings.ReplaceAll(valueStr, "\\\"", "*")
						valueStr = strings.ReplaceAll(valueStr, "\"", "")
						valueStr = "\"" + valueStr + "\"" + valueSuffix
					}

					// "xxx/xxx", "xxx/"xxx"
					if valueOnlySymbol.MatchString(valueStr) {
						valueStr = strings.ReplaceAll(valueStr, "\\", "*")
						valueStr = strings.ReplaceAll(valueStr, "\"", "*")
						valueStr = "\"" + valueStr[1:len(valueStr)-1] + "\"" + valueSuffix
					}
				} else {
					// For nested JSON strings, skip all transformations and just add the suffix
					valueStr += valueSuffix
				}

				configSlice[i] = v[:valueIndex+1] + valueStr
			}

			configStr = strings.Join(configSlice, "")
			configRune := []rune(configStr)

			// match the bracket
			if toCheck = bracketMatch(funcName, configIndex, configRune); !toCheck {
				return toCheck
			}

			configStr = string(configRune)
			jsonData := []byte(configStr)
			var v interface{}
			err := json.Unmarshal(jsonData, &v)
			if err != nil {
				log.Errorf("fail to unmarshal func %v's number %v config: \n%s\n%s", funcName, configIndex, configStr, err)
				return false
			}
			data := v.(map[string]interface{})

			f.stepAttributes[configIndex] = data

			parseAttr(configIndex, "", data, attributeValueMap, testMustSet, testModifySet)
		}
	}
	return true
}

func bracketMatch(funcName string, configIndex int, s []rune) bool {
	var stack []string

	for i := 0; i < len(s); i++ {
		r := string(s[i])
		switch r {
		// remove the extra comma
		case ",":
			if i != len(s)-1 && (s[i+1] == '}' || s[i+1] == ']') {
				s[i] = rune(' ')
			}
		case "{", "[":
			stack = append(stack, r)
		case "}", "]":
			{
				if len(stack) == 0 {
					log.Errorf("fail to math bracket [ ] in func %s's number %d config:%s\n", funcName, configIndex, string(s))
					return false
				}
				temp := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if temp == "[" && r == "}" {
					s[i] = ']'
				}
				if bracket[temp] == r {
					break
				}
			}
		}
	}
	return true
}

func parseAttr(configIndex int, rootName string, data interface{}, attributeValueMap map[string]string,
	testMustSet, testModifySet *mapset.Set) {
	if d, ok := data.(map[string]interface{}); ok {
		for key, value := range d {
			if rootName != "" {
				key = rootName + "." + key
			}
			(*testMustSet).Add(key)
			// check if the attribute has been updated
			if v, ok := attributeValueMap[key]; ok {
				if fmt.Sprintf("%v", value) != v {
					(*testModifySet).Add(key)
				}
			} else if configIndex > 0 {
				(*testModifySet).Add(key)
			}
			attributeValueMap[key] = fmt.Sprintf("%v", value)
			parseAttr(configIndex, key, value, attributeValueMap, testMustSet, testModifySet)
		}
	} else if d, ok := data.([]interface{}); ok {
		for _, v := range d {
			attributeValueMap[rootName] = fmt.Sprintf("%v", data)
			parseAttr(configIndex, rootName, v, attributeValueMap, testMustSet, testModifySet)
		}
	}
}

var (
	commentedRegex    = regexp.MustCompile("^[\t]*//")
	normalFuncRegex   = regexp.MustCompile("^func Test(.*)")
	unitFuncRegex     = regexp.MustCompile("^func TestUnit(.*)")
	standardFuncRegex = regexp.MustCompile("^func TestAccAliCloud(.*)")
	configRegex       = regexp.MustCompile("(^[\t]*)Config:(.*)")
	checkRegex        = regexp.MustCompile("(.*)Check:(.*)")
	ignoreRegex       = regexp.MustCompile("(.*)ImportStateVerifyIgnore:(.*)")
	hasNumRegex       = regexp.MustCompile(`[0-9]+`)
	attrRegex         = regexp.MustCompile("^([{]*)\"([a-zA-Z_0-9-.]+)\":(.*)")
	symbolRegex       = regexp.MustCompile(`\s`)
	variableRegex     = regexp.MustCompile("(^[a-zA-Z_0-9]+)|(\"[+]\")")
	valueFuncRegex    = regexp.MustCompile("[(].*[\"].*[\"].*[)]")
	valueOnlySymbol   = regexp.MustCompile(`.*([^\\\"])(\\)([^\\\"]).*`)
	bracket           = map[string]string{
		"{": "}",
		"[": "]",
	}
	matchMap = map[int][]string{
		0: []string{"[]string{", "["},
		1: []string{"[]map[string]interface{}{", "["},
		2: []string{"[]map[string]string{", "["},
		3: []string{"[]interface{}{", "["},
		4: []string{"map[string]interface{}{", "{"},
		5: []string{"map[string]string{", "{"},
		6: []string{"\\n", ""},
		7: []string{"`${map(", "["},
		8: []string{")}`", "]"},
	}

	isNameCorrect = true
)

type ResourceTest struct {
	resourceName string
	funcs        map[string]FuncTest
}

type FuncTest struct {
	stepStr        map[int]string
	stepAttributes map[int]map[string]interface{}
}

// TreeNode represents a node in the attribute tree
type TreeNode struct {
	Name     string
	Children map[string]*TreeNode
}

// buildAttributeTree builds a tree structure from attribute paths
func buildAttributeTree(attributes []string) *TreeNode {
	root := &TreeNode{
		Name:     "root",
		Children: make(map[string]*TreeNode),
	}

	for _, attr := range attributes {
		parts := strings.Split(attr, ".")
		current := root
		for _, part := range parts {
			if _, exists := current.Children[part]; !exists {
				current.Children[part] = &TreeNode{
					Name:     part,
					Children: make(map[string]*TreeNode),
				}
			}
			current = current.Children[part]
		}
	}

	return root
}

// printTree prints the tree structure with proper formatting
func printTree(node *TreeNode, prefix string, isLast bool, buffer *strings.Builder) {
	if node.Name != "root" {
		// Determine the branch symbol
		branch := "├── "
		if isLast {
			branch = "└── "
		}
		buffer.WriteString(prefix + branch + node.Name + "\n")

		// Update prefix for children
		if isLast {
			prefix += "    "
		} else {
			prefix += "│   "
		}
	}

	// Sort children for consistent output
	childNames := make([]string, 0, len(node.Children))
	for name := range node.Children {
		childNames = append(childNames, name)
	}
	sort.Strings(childNames)

	// Print children
	for i, name := range childNames {
		child := node.Children[name]
		isLastChild := i == len(childNames)-1
		printTree(child, prefix, isLastChild, buffer)
	}
}

// formatAttributesAsTree converts attribute list to tree format string
func formatAttributesAsTree(attributes []interface{}) string {
	if len(attributes) == 0 {
		return ""
	}

	// Convert interface{} slice to string slice
	attrStrings := make([]string, 0, len(attributes))
	for _, attr := range attributes {
		if attrStr, ok := attr.(string); ok {
			attrStrings = append(attrStrings, attrStr)
		}
	}

	// Sort for consistent output
	sort.Strings(attrStrings)

	// Build tree
	tree := buildAttributeTree(attrStrings)

	// Print tree
	var buffer strings.Builder
	printTree(tree, "", true, &buffer)

	return buffer.String()
}

func checkAttributeSet(resourceName string, fileType string, schemaMustSet, testMustSet,
	schemaModifySet, testModifySet, schemaForceNewSet, schemaAllSet, testIgnoreSet mapset.Set) bool {

	isFullCover, isIgnoreLegal, isAllModified := true, true, true

	notCoverSlice := schemaMustSet.Difference(testMustSet).ToSlice()
	if len(notCoverSlice) != 0 {
		isFullCover = false
		schemaCount := float64(len(schemaMustSet.ToSlice()))
		notCoverCount := float64(len(notCoverSlice))
		coverageRate := 1 - (notCoverCount / schemaCount)
		log.Infof("resource %s attributes has %.2f%% testing coverage rate ", resourceName, coverageRate*100)

		// Format as tree structure
		treeStr := formatAttributesAsTree(notCoverSlice)
		log.Errorf("resource %s attributes missing test cases (%d attributes):\n%s", resourceName, len(notCoverSlice), treeStr)
	} else {
		log.Infof("resource %s attributes has 100%% testing coverage rate ", resourceName)
	}

	forceNewButIgnore := schemaForceNewSet.Intersect(testIgnoreSet).ToSlice()
	if len(forceNewButIgnore) != 0 {
		isIgnoreLegal = false
		forceNewButIgnoreStr, _ := json.Marshal(forceNewButIgnore)
		// TODO: 从READ方法区分是否是私有属性，从而区分应该修改ignore数组还是应该修改资源属性
		log.Errorf("resource %s [ForceNew] attributes %v are in ImportStateVerifyIgnore array ", resourceName, string(forceNewButIgnoreStr))
	}
	redundantAttr := testIgnoreSet.Difference(schemaAllSet).ToSlice()
	redundantAttrFinal := []string{}
	for _, v := range redundantAttr {
		vStr := v.(string)
		if hasNumRegex.MatchString(vStr) && strings.Contains(vStr, ".") {
			parts := strings.Split(vStr, ".")
			attrStr := parts[0]
			for _, subAttr := range parts[1:] {
				if _, err := strconv.Atoi(subAttr); err == nil {
					continue
				} else {
					attrStr += "." + subAttr
				}
			}
			if !schemaAllSet.Contains(attrStr) {
				redundantAttrFinal = append(redundantAttrFinal, vStr)
			}
		} else {
			redundantAttrFinal = append(redundantAttrFinal, vStr)
		}
	}
	if len(redundantAttrFinal) != 0 {
		isIgnoreLegal = false
		redundantAttrFinalStr, _ := json.Marshal(redundantAttrFinal)
		log.Errorf("resource %s attributes %v should not in ImportStateVerifyIgnore array", resourceName, string(redundantAttrFinalStr))
	}
	schemaModifySet = schemaModifySet.Difference(testIgnoreSet)

	// Get immutable attributes from documentation
	immutableAttrs := getImmutableAttributesFromDoc(resourceName, fileType)
	if immutableAttrs.Cardinality() > 0 {
		log.Infof("Found %d immutable attributes in documentation: %v", immutableAttrs.Cardinality(), immutableAttrs.ToSlice())

		// Build a set that includes both the immutable attributes and all their nested fields
		immutableWithNested := mapset.NewSet()
		for _, attr := range immutableAttrs.ToSlice() {
			attrStr := attr.(string)
			immutableWithNested.Add(attrStr)

			// Also add any nested fields (e.g., if "resources" is immutable, also exclude "resources.resource_id")
			for _, schemaAttr := range schemaModifySet.ToSlice() {
				schemaAttrStr := schemaAttr.(string)
				if strings.HasPrefix(schemaAttrStr, attrStr+".") {
					immutableWithNested.Add(schemaAttrStr)
					log.Debugf("Excluding nested immutable attribute: %s", schemaAttrStr)
				}
			}
		}

		schemaModifySet = schemaModifySet.Difference(immutableWithNested)
	}

	notModifySlice := schemaModifySet.Difference(testModifySet).ToSlice()
	if len(notModifySlice) != 0 {
		isAllModified = false
		schemaCount := float64(len(schemaModifySet.ToSlice()))
		notCoverCount := float64(len(notModifySlice))
		coverageRate := 1 - (notCoverCount / schemaCount)
		log.Infof("resource %s attributes has %.2f%% modified coverage rate ", resourceName, coverageRate*100)

		// Format as tree structure
		treeStr := formatAttributesAsTree(notModifySlice)
		log.Errorf("resource %s attributes missing modification in test cases (%d attributes):\n%s", resourceName, len(notModifySlice), treeStr)
	} else {
		log.Infof("resource %s attributes has 100%% modified coverage rate ", resourceName)
	}

	return isFullCover && isIgnoreLegal && isAllModified
}

// getImmutableAttributesFromDoc reads the documentation and finds attributes marked as immutable
// These attributes have descriptions like "Note: The parameter is immutable after resource creation"
func getImmutableAttributesFromDoc(resourceName string, fileType string) mapset.Set {
	immutableSet := mapset.NewSet()

	// Remove "alicloud_" prefix if present for documentation path
	docResourceName := strings.TrimPrefix(resourceName, "alicloud_")

	// Construct documentation path
	var docPath string
	if fileType == "resource" {
		docPath = "website/docs/r/" + docResourceName + ".html.markdown"
	} else {
		docPath = "website/docs/d/" + docResourceName + ".html.markdown"
	}

	// Read documentation file
	file, err := os.Open(docPath)
	if err != nil {
		log.Debugf("Cannot open documentation file %s: %v", docPath, err)
		return immutableSet
	}
	defer file.Close()

	// Regex patterns to match immutable declarations
	immutablePattern := regexp.MustCompile(`\*\*Note:\s*The parameter is immutable after resource creation`)
	// Pattern to extract attribute name: * `attr_name` - (Optional/Required...) Description
	attrPattern := regexp.MustCompile(`^\*\s+\x60([a-zA-Z_0-9]+)\x60\s+-\s+\(`)

	scanner := bufio.NewScanner(file)
	var currentAttr string
	var continuedLine string

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this line defines an attribute
		if matches := attrPattern.FindStringSubmatch(line); matches != nil {
			// Save previous attribute if it was immutable
			if currentAttr != "" && strings.Contains(continuedLine, "immutable after resource creation") {
				// Add the attribute and all its nested fields
				immutableSet.Add(currentAttr)
				log.Debugf("Found immutable attribute: %s", currentAttr)
			}

			// Start tracking new attribute
			currentAttr = matches[1]
			continuedLine = line
		} else if currentAttr != "" {
			// Continue reading the description (might span multiple lines)
			continuedLine += " " + strings.TrimSpace(line)

			// Stop if we hit a new section or empty line
			if strings.HasPrefix(line, "##") || strings.HasPrefix(line, "###") {
				// Check if current attribute is immutable before resetting
				if strings.Contains(continuedLine, "immutable after resource creation") {
					immutableSet.Add(currentAttr)
					log.Debugf("Found immutable attribute: %s", currentAttr)
				}
				currentAttr = ""
				continuedLine = ""
			}
		}

		// Also check for inline immutable declarations
		if immutablePattern.MatchString(line) && currentAttr != "" {
			immutableSet.Add(currentAttr)
			log.Debugf("Found immutable attribute: %s", currentAttr)
		}
	}

	// Check the last attribute
	if currentAttr != "" && strings.Contains(continuedLine, "immutable after resource creation") {
		immutableSet.Add(currentAttr)
		log.Debugf("Found immutable attribute: %s", currentAttr)
	}

	return immutableSet
}

func testAll(diff *diffparser.Diff) {
	file, err := os.Open("filename.txt")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	diff.Files = make([]*diffparser.DiffFile, 0, 800)
	line := 0
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "resource_alicloud") {
			diff.Files = append(diff.Files, &diffparser.DiffFile{NewName: "alicloud/" + text})
		}
		line++
	}
}
