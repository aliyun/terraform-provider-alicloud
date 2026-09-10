//nolint:all
package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Common regex patterns
var (
	commentedRegex    = regexp.MustCompile("^[\t]*//")
	normalFuncRegex   = regexp.MustCompile("^func Test(.*)")
	unitFuncRegex     = regexp.MustCompile("^func TestUnit(.*)")
	standardFuncRegex = regexp.MustCompile("^func TestAccAliCloud(.*)")
	configRegex       = regexp.MustCompile("(^[\t]*)Config:(.*)")
	checkRegex        = regexp.MustCompile("(.*)Check:(.*)")
	ignoreRegex       = regexp.MustCompile("(.*)ImportStateVerifyIgnore:(.*)")
	symbolRegex       = regexp.MustCompile(`\s`)
	variableRegex     = regexp.MustCompile("(^[a-zA-Z_0-9]+)|(\"[+]\")")
	valueFuncRegex    = regexp.MustCompile("[(].*[\"].*[\"].*[)]")
	valueOnlySymbol   = regexp.MustCompile(`.*([^\\\"])(\\)([^\\\"]).*`)

	bracket = map[string]string{
		"{": "}",
		"[": "]",
	}
	matchMap = map[int][]string{
		0: {"[]string{", "["},
		1: {"[]map[string]interface{}{", "["},
		2: {"[]map[string]string{", "["},
		3: {"[]interface{}{", "["},
		4: {"map[string]interface{}{", "{"},
		5: {"map[string]string{", "{"},
		6: {"\\n", ""},
		7: {"`${map(", "["},
		8: {")}`", "]"},
	}
)

// SchemaAttribute represents a schema attribute
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

// ResourceTest represents test cases for a resource
type ResourceTest struct {
	FilePath string
	Funcs    map[string]FuncTest
}

// FuncTest represents a single test function
type FuncTest struct {
	StepStr        map[int]string
	StepAttributes map[int]map[string]interface{}
}

// GetSchemaAttributes gets all schema attributes for a resource
func GetSchemaAttributes(resName string, isResource bool) (mapset.Set, mapset.Set, mapset.Set, mapset.Set) {
	schemaAllSet := mapset.NewSet()
	schemaMustSet := mapset.NewSet()
	schemaModifySet := mapset.NewSet()
	schemaForceNewSet := mapset.NewSet()

	var resource *schema.Resource
	var ok bool

	if isResource {
		resource, ok = alicloud.Provider().(*schema.Provider).ResourcesMap["alicloud_"+resName]
	} else {
		resource, ok = alicloud.Provider().(*schema.Provider).DataSourcesMap["alicloud_"+resName]
	}

	if !ok || resource == nil {
		return schemaAllSet, schemaMustSet, schemaModifySet, schemaForceNewSet
	}

	schemaAttributes := make(map[string]SchemaAttribute)
	getSchemaAttributesRecursive("", schemaAttributes, resource.Schema)

	for key, value := range schemaAttributes {
		if key == "dry_run" || value.Deprecated != "" {
			continue
		}
		schemaAllSet.Add(key)
		if value.Optional || value.Required {
			schemaMustSet.Add(key)
			if !value.ForceNew {
				schemaModifySet.Add(key)
			}
		}
		if value.ForceNew {
			schemaForceNewSet.Add(key)
		}
	}

	return schemaAllSet, schemaMustSet, schemaModifySet, schemaForceNewSet
}

func getSchemaAttributesRecursive(rootName string, schemaAttributes map[string]SchemaAttribute,
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
			} else if v, ok := value.Elem.(*schema.Resource); ok {
				vv := schemaAttributes[key]
				vv.ElemType = "Object"
				schemaAttributes[key] = vv
				getSchemaAttributesRecursive(key, schemaAttributes, v.Schema)
			}
		}
	}
}

// ParseTestFile parses a test file and returns test cases with covered and modified attributes
func ParseTestFile(filePath string, testMustSet, testModifySet, testIgnoreSet *mapset.Set) (map[string]*TestCaseInfo, bool) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("fail to open test file %s. Error: %s", filePath, err)
		return nil, false
	}
	defer file.Close()

	testCases := make(map[string]*TestCaseInfo)
	scanner := bufio.NewScanner(file)

	resourceTest := ResourceTest{
		FilePath: filePath,
		Funcs:    map[string]FuncTest{},
	}

	line := 0
	funcName := ""
	stepNumber := 0
	configStr := ""
	inConfig := false
	inFunc := false
	ignoreStr := ""
	inIgnore := false
	funcStartLine := 0

	for scanner.Scan() {
		line++
		text := scanner.Text()

		if commentedRegex.MatchString(text) {
			continue
		}

		if text == "}" && inFunc {
			inFunc = false
			funcName = ""
		}

		if normalFuncRegex.MatchString(text) && !unitFuncRegex.MatchString(text) {
			inFunc = true
			funcName = text[strings.Index(text, "T"):strings.Index(text, "(")]
			funcStartLine = line
			stepNumber = 0
			resourceTest.Funcs[funcName] = FuncTest{
				StepStr:        map[int]string{},
				StepAttributes: map[int]map[string]interface{}{},
			}

			// Initialize test case info
			testCases[funcName] = &TestCaseInfo{
				Name:               funcName,
				CoveredAttributes:  mapset.NewSet(),
				ModifiedAttributes: mapset.NewSet(),
				HasImportTest:      false,
				Steps:              0,
				LineNumber:         funcStartLine,
			}
		}

		if inFunc {
			// Check for ImportState
			if strings.Contains(text, "ImportState:") && strings.Contains(text, "true") {
				if testCases[funcName] != nil {
					testCases[funcName].HasImportTest = true
				}
			}

			if configRegex.MatchString(text) {
				configStr += text
				inConfig = true
				continue
			}

			if inConfig {
				if checkRegex.MatchString(text) {
					resourceTest.Funcs[funcName].StepStr[stepNumber] = configStr
					inConfig = false
					configStr = ""
					stepNumber++
				} else {
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

	// Parse configs and populate test case info
	success := parseConfigForTestCases(resourceTest, testMustSet, testModifySet, testCases)

	// Update step counts
	for funcName, testCase := range testCases {
		if f, ok := resourceTest.Funcs[funcName]; ok {
			testCase.Steps = len(f.StepStr)
		}
	}

	return testCases, success
}

// TestCaseInfo contains information about a test case
type TestCaseInfo struct {
	Name               string
	CoveredAttributes  mapset.Set
	ModifiedAttributes mapset.Set
	HasImportTest      bool
	Steps              int
	LineNumber         int
	EstimatedDuration  int
	ComplexityScore    float64
}

func parseConfigForTestCases(resourceTest ResourceTest, testMustSet, testModifySet *mapset.Set,
	testCases map[string]*TestCaseInfo) bool {

	for funcName, f := range resourceTest.Funcs {
		attributeValueMap := map[string]string{}

		for configIndex := 0; configIndex < len(f.StepStr); configIndex++ {
			configStr := f.StepStr[configIndex]

			if !strings.Contains(configStr, "{") {
				log.Debugf("the test case in [%s] does not use a standard template", funcName)
				continue
			}

			configStr = configStr[strings.Index(configStr, "{")+2 : strings.LastIndex(configStr, "}")+1]
			if configStr == "{}" {
				continue
			}

			configSlice := strings.Split(configStr, "\n")

			for i, v := range configSlice {
				valueIndex := strings.Index(v, ":")

				if strings.Contains(v, ":") {
					splitIndex := strings.Index(v, ":")
					beforeCount := strings.Count(v[:splitIndex], "\"")
					if beforeCount%2 != 0 {
						valueIndex = -1
					}
				}

				if valueIndex < 0 {
					continue
				}

				valueSuffix := string(v[len(v)-1])
				valueStr := v[valueIndex+1 : len(v)-1]
				bracketCount := strings.Count(v, "]") + strings.Count(v, "[")
				if v[valueIndex+1] == '"' && strings.HasSuffix(v, "},") ||
					strings.HasSuffix(v, "],") && bracketCount%2 == 1 {
					valueSuffix = v[len(v)-2:]
					valueStr = v[valueIndex+1 : len(v)-2]
				}

				isNestedJSONString := strings.HasPrefix(valueStr, "\"") &&
					strings.Contains(valueStr, `\\\"`) &&
					(strings.Contains(valueStr, "{") || strings.Contains(valueStr, "["))

				if !isNestedJSONString {
					if strings.Contains(valueStr, "+") {
						vv := valueStr
						index := strings.Index(vv, "+")
						for index != -1 {
							beforeStr := vv[:index]
							beforeStr = strings.TrimSpace(beforeStr)
							if strings.Count(beforeStr, "\"")%2 == 0 ||
								(beforeStr[0] == '`' && beforeStr[len(beforeStr)-1] == '`') {
								vv = vv[index+1:]
								index = strings.Index(vv, "+")
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

					if strings.HasPrefix(valueStr, "`") && strings.HasSuffix(valueStr, "`") {
						valueStr = "\"" + valueStr[:len(valueStr)-1] + "\"" + valueSuffix
					} else {
						valueStr += valueSuffix
					}

					for i := 0; i < len(matchMap); i++ {
						k := matchMap[i][0]
						vv := matchMap[i][1]
						if strings.Contains(valueStr, k) {
							valueStr = strings.ReplaceAll(valueStr, k, vv)
						}
					}

					if variableRegex.MatchString(valueStr) || valueFuncRegex.MatchString(valueStr) {
						valueStr = strings.ReplaceAll(valueStr, "\\\"", "*")
						valueStr = strings.ReplaceAll(valueStr, "\"", "")
						valueStr = "\"" + valueStr + "\"" + valueSuffix
					}

					if valueOnlySymbol.MatchString(valueStr) {
						valueStr = strings.ReplaceAll(valueStr, "\\", "*")
						valueStr = strings.ReplaceAll(valueStr, "\"", "*")
						valueStr = "\"" + valueStr[1:len(valueStr)-1] + "\"" + valueSuffix
					}
				} else {
					valueStr += valueSuffix
				}

				configSlice[i] = v[:valueIndex+1] + valueStr
			}

			configStr = strings.Join(configSlice, "")
			configRune := []rune(configStr)

			if !BracketMatch(configRune) {
				continue
			}

			configStr = string(configRune)
			jsonData := []byte(configStr)
			var data interface{}
			err := json.Unmarshal(jsonData, &data)
			if err != nil {
				log.Debugf("fail to unmarshal func %v's number %v config: %s", funcName, configIndex, err)
				continue
			}

			dataMap, ok := data.(map[string]interface{})
			if !ok {
				continue
			}

			f.StepAttributes[configIndex] = dataMap

			// Get test case info
			testCase := testCases[funcName]
			if testCase != nil {
				parseAttrForTestCase(configIndex, "", dataMap, attributeValueMap, testMustSet, testModifySet, testCase)
			}
		}
	}

	return true
}

func parseAttrForTestCase(configIndex int, rootName string, data interface{},
	attributeValueMap map[string]string, testMustSet, testModifySet *mapset.Set, testCase *TestCaseInfo) {

	if d, ok := data.(map[string]interface{}); ok {
		for key, value := range d {
			if rootName != "" {
				key = rootName + "." + key
			}

			// Add to covered attributes
			(*testMustSet).Add(key)
			testCase.CoveredAttributes.Add(key)

			// Check if modified
			if v, ok := attributeValueMap[key]; ok {
				if fmt.Sprintf("%v", value) != v {
					(*testModifySet).Add(key)
					testCase.ModifiedAttributes.Add(key)
				}
			} else if configIndex > 0 {
				(*testModifySet).Add(key)
				testCase.ModifiedAttributes.Add(key)
			}
			attributeValueMap[key] = fmt.Sprintf("%v", value)

			parseAttrForTestCase(configIndex, key, value, attributeValueMap, testMustSet, testModifySet, testCase)
		}
	} else if d, ok := data.([]interface{}); ok {
		for _, v := range d {
			attributeValueMap[rootName] = fmt.Sprintf("%v", data)
			parseAttrForTestCase(configIndex, rootName, v, attributeValueMap, testMustSet, testModifySet, testCase)
		}
	}
}

// BracketMatch checks if brackets are balanced and fixes mismatches
func BracketMatch(s []rune) bool {
	var stack []string

	for i := 0; i < len(s); i++ {
		r := string(s[i])
		switch r {
		case ",":
			if i != len(s)-1 && (s[i+1] == '}' || s[i+1] == ']') {
				s[i] = rune(' ')
			}
		case "{", "[":
			stack = append(stack, r)
		case "}", "]":
			if len(stack) == 0 {
				return false
			}
			temp := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if temp == "[" && r == "}" {
				s[i] = ']'
			}
			if bracket[temp] != r {
				return false
			}
		}
	}
	return true
}
