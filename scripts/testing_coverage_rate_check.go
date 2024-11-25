//nolint:all
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
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
	fileNames = flag.String("fileNames", "", "the files to check diff")

	resourceFileRegex     = regexp.MustCompile("alicloud/(resource)[0-9a-zA-Z_]*")
	resourceFileTestRegex = regexp.MustCompile("alicloud/(resource)[0-9a-zA-Z_]*_test.go")
)

func main() {
	exitCode := 0
	flag.Parse()
	if fileNames != nil && len(*fileNames) == 0 {
		log.Infof("the diff file is empty, shipped!")
		return
	}

	byt, _ := os.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	resourceNameMap := make(map[string]struct{})
	for _, file := range diff.Files {
		isNameCorrect = true
		resourceName := ""
		isResource := true
		fileType := "resource"
		if resourceFileTestRegex.MatchString(file.NewName) {
			resourceName = strings.TrimSuffix(strings.Split(file.NewName, "/")[1], "_test.go")
		} else if resourceFileRegex.MatchString(file.NewName) {
			resourceName = strings.TrimSuffix(strings.Split(file.NewName, "/")[1], ".go")
		} else {
			continue
		}

		if strings.HasPrefix(resourceName, "data_source_") {
			isResource = false
			fileType = "data source"
			resourceName = strings.TrimPrefix(resourceName, "data_source_")
		} else {
			resourceName = strings.TrimPrefix(resourceName, "resource_")
		}

		if _, ok := resourceNameMap[resourceName]; ok {
			continue
		} else {
			resourceNameMap[resourceName] = struct{}{}
		}

		log.Infof("==> Getting %s %s attributes...", fileType, resourceName)
		var resource = &schema.Resource{}
		var ok bool
		if isResource {
			resource, ok = alicloud.Provider().(*schema.Provider).ResourcesMap[resourceName]
			if !ok || resource == nil {
				log.Errorf("resource %s is not found in the provider ResourceMap\n\n", resourceName)
				exitCode = 1
				continue
			}
		} else {
			resource, ok = alicloud.Provider().(*schema.Provider).DataSourcesMap[resourceName]
			if !ok || resource == nil {
				log.Errorf("data source %s is not found in the provider DataSourcesMap\n\n", resourceName)
				exitCode = 1
				continue
			}
		}
		schemaAllSet, schemaMustSet, schemaModifySet, schemaForceNewSet :=
			mapset.NewSet(), mapset.NewSet(), mapset.NewSet(), mapset.NewSet()
		getSchemaAttr(false, resource.Schema, &schemaAllSet, &schemaMustSet, &schemaModifySet, &schemaForceNewSet)

		log.Infof("==> Getting %s %s attributes in test cases...", fileType, resourceName)
		testMustSet, testModifySet, testIgnoreSet :=
			mapset.NewSet(), mapset.NewSet(), mapset.NewSet()
		filePath := "alicloud/"
		if isResource {
			filePath += "resource_" + resourceName + "_test.go"
		} else {
			filePath += "data_source_" + resourceName + "_test.go"
		}
		check := getTestCaseAttr(filePath, resourceName, &testMustSet, &testModifySet, &testIgnoreSet)

		// "check" denotes the test code is using a standard template
		if check {
			log.Infof("==> checking %s %s attributes' coverage rate", fileType, resourceName)
			if checkAttributeSet(resourceName, fileType, schemaMustSet, testMustSet,
				schemaModifySet, testModifySet, schemaForceNewSet, schemaAllSet, testIgnoreSet) && isNameCorrect {
				log.Infof("--- PASS!\n\n")
				continue
			}
		}

		log.Errorf("--- Failed!\n\n")
		exitCode = 1

	}

	os.Exit(exitCode)

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
		notCoverStr, _ := json.Marshal(notCoverSlice)
		log.Errorf("resource %s attributes %v missing test cases", resourceName, string(notCoverStr))
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

	notModifySlice := schemaModifySet.Difference(testModifySet).ToSlice()
	if len(notModifySlice) != 0 {
		isAllModified = false
		schemaCount := float64(len(schemaModifySet.ToSlice()))
		notCoverCount := float64(len(notModifySlice))
		coverageRate := 1 - (notCoverCount / schemaCount)
		log.Infof("resource %s attributes has %.2f%% modified coverage rate ", resourceName, coverageRate*100)
		notModifyStr, _ := json.Marshal(notModifySlice)
		log.Errorf("resource %s attributes %v missing modification in test cases", resourceName, string(notModifyStr))
	} else {
		log.Infof("resource %s attributes has 100%% modified coverage rate ", resourceName)
	}

	return isFullCover && isIgnoreLegal && isAllModified
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
