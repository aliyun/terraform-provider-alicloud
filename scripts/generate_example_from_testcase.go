package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/zclconf/go-cty/cty"

	roa "github.com/alibabacloud-go/tea-roa/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	credential "github.com/aliyun/credentials-go/credentials"
	mapset "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
)

var (
	resourceMap   = alicloud.Provider().(*schema.Provider).ResourcesMap
	dataSourceMap = alicloud.Provider().(*schema.Provider).DataSourcesMap

	// tempFilePath          = "./testcases"
	// examplePath           = "./testcase-example/"
	ossSchemaPath         = "ossSchema.json"
	ossNotFoundSchemaPath = "ossNotFound.json"

	resourceBlocksKeyWords = []string{"depends_on", "count",
		"for_each", "provider", "lifecycle", "provisioner"}
	resourceBlocksSet = mapset.NewSet()

	alicloudProvider = "terraform {\n" +
		"  required_providers {\n" +
		"    alicloud = {\n" +
		"      source  = \"aliyun/alicloud\"\n" +
		"    }\n" +
		"  }\n" +
		"}"

	tempFilePath = flag.String("testCasePath", "testcases", "testcase path")
	examplePath  = flag.String("outputPath", "testcase-example", "example output path")
	updateSchema = flag.Bool("updateSchema", false, "Update schema.json file to outputDir")
	enableDocs   = flag.Bool("enableDocs", false, "Enable document generate")
)

func main() {

	flag.Parse()
	if tempFilePath != nil && len(*tempFilePath) == 0 {
		log.Infof("the testCasePath is empty!")
		return
	}
	if examplePath != nil && len(*examplePath) == 0 {
		log.Infof("the example outputPath is empty!")
		return
	}

	for _, v := range resourceBlocksKeyWords {
		resourceBlocksSet.Add(v)
	}

	tfFileList := GetFiles(*tempFilePath)

	resourceTestCasesMap := map[string]map[string]*ResourceTestExp{}
	resourceTestCasesMap["resource"] = map[string]*ResourceTestExp{}
	resourceTestCasesMap["data-source"] = map[string]*ResourceTestExp{}

	for _, filePath := range tfFileList {
		// fmt.Println(filePath)
		parseFile(filePath, resourceTestCasesMap)
	}

	ossResourceSchema := map[string]interface{}{}
	ossResourceSchema["resource"] = map[string]interface{}{}
	ossResourceSchema["data-source"] = map[string]interface{}{}

	if *updateSchema {
		getOSSResourceInfo(ossSchemaPath, ossNotFoundSchemaPath, resourceTestCasesMap, ossResourceSchema)
	}

	mergeSteps(ossSchemaPath, resourceTestCasesMap)

	mergeTestCases(resourceTestCasesMap)

	processExp(resourceTestCasesMap)

	printExpToFile(resourceTestCasesMap)
}

func GetFiles(folder string) []string {
	tfFileList := []string{}
	files, _ := os.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			tfFileList = append(tfFileList, GetFiles(folder+"/"+file.Name())...)
		} else {
			// fmt.Println(folder + "/" + file.Name())
			tfFileList = append(tfFileList, folder+"/"+file.Name())
		}
	}

	return tfFileList
}

func parseFile(filePath string, resourceTestCasesMap map[string]map[string]*ResourceTestExp) {

	path := strings.Split(filePath, "/")
	resourceName := ""
	resoureType := path[len(path)-4]
	fileName := path[len(path)-1]

	if strings.HasSuffix(fileName, "dependence.tf") {
		// dependence
		dependence, err := os.ReadFile(filePath)
		if err != nil {
			log.Errorf("fail to open file %s. Error: %s", filePath, err)
		}

		dependenceStr := string(dependence)
		dependenceSlice := strings.SplitN(dependenceStr, "\n", 3)
		resourceName = strings.Split(dependenceSlice[0], " ")[1]
		primaryConfig := dependenceSlice[1][2:]
		dependenceStr = dependenceSlice[2]

		if _, ok := resourceTestCasesMap[resoureType][resourceName]; !ok {
			resourceTestCasesMap[resoureType][resourceName] =
				&ResourceTestExp{
					SchemaAttributes: map[string]*SchemaAttr{},
					TestCases:        map[string]*TestCase{},
					Examples:         map[string]string{},
					ExpDescriptions:  map[string]string{},
				}
		}

		testCaseName := fileName[:strings.Index(fileName, "-dependence.tf")]

		if testCase, ok := resourceTestCasesMap[resoureType][resourceName].
			TestCases[testCaseName]; !ok {
			resourceTestCasesMap[resoureType][resourceName].TestCases[testCaseName] =
				&TestCase{
					Dependences:   dependenceStr,
					PrimaryConfig: primaryConfig,
				}
		} else {
			testCase.Dependences = dependenceStr
			testCase.PrimaryConfig = primaryConfig
			resourceTestCasesMap[resoureType][resourceName].TestCases[testCaseName] = testCase
		}

	} else {
		// step
		file, err := os.Open(filePath)
		if err != nil {
			log.Errorf("fail to open file %s. Error: %s", filePath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		line := 0
		steps := []map[string]interface{}{}
		for scanner.Scan() {
			line++
			text := scanner.Text()
			if strings.HasPrefix(text, "#") {
				if resourceName == "" {
					resourceName = strings.Split(text, " ")[1]
				}
				continue
			}
			jsonData := []byte(text)
			var stepMap interface{}
			json.Unmarshal(jsonData, &stepMap)
			if stepMap == nil {
				log.Errorf("fail to unmarshal file %v's  %v step", filePath, line)
				continue
			}
			data := stepMap.(map[string]interface{})
			steps = append(steps, data)
		}

		if _, ok := resourceTestCasesMap[resoureType][resourceName]; !ok {
			resourceTestCasesMap[resoureType][resourceName] =
				&ResourceTestExp{
					SchemaAttributes: map[string]*SchemaAttr{},
					TestCases:        map[string]*TestCase{},
					Examples:         map[string]string{},
					ExpDescriptions:  map[string]string{},
				}
		}

		testCaseName := fileName[:strings.Index(fileName, ".txt")]

		if testCase, ok := resourceTestCasesMap[resoureType][resourceName].
			TestCases[testCaseName]; !ok {
			resourceTestCasesMap[resoureType][resourceName].TestCases[testCaseName] =
				&TestCase{
					Steps: steps,
				}
		} else {
			testCase.Steps = steps
			resourceTestCasesMap[resoureType][resourceName].TestCases[testCaseName] = testCase
		}

	}
}

func getOSSResourceInfo(ossSchemaPath, ossNotFoundSchemaPath string, resourceTestCasesMap map[string]map[string]*ResourceTestExp,
	ossResourceSchema map[string]interface{}) {

	ossNotFound := map[string]struct{}{}

	for resourceType, resources := range resourceTestCasesMap {
		for resourceName := range resources {

			resourceInfo := sendRequest(resourceName)
			if resourceInfo == nil {
				ossNotFound[resourceName] = struct{}{}
				continue
			}
			var attributeMap map[string]interface{}
			if v, ok := resourceInfo["resourceType"]; ok {
				if value, ok := v.(map[string]interface{})["properties"]; ok {
					attributeMap = value.(map[string]interface{})
				}
			}
			if attributeMap == nil {
				return
			}
			ossResourceSchema[resourceType].(map[string]interface{})[resourceName] = attributeMap

			fmt.Println("resource: ", resourceName)
		}
	}

	ossFile, err := os.OpenFile(ossSchemaPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ossFile.Close()

	data, _ := json.Marshal(ossResourceSchema)
	ossFile.WriteString(string(data) + "\n")

	ossNotFoundFile, err := os.OpenFile(ossNotFoundSchemaPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ossNotFoundFile.Close()

	notFoundData, _ := json.Marshal(ossNotFound)
	ossNotFoundFile.WriteString(string(notFoundData) + "\n")

}

func sendRequest(resourceName string) map[string]interface{} {
	ALICLOUD_ACCESS_KEY := os.Getenv("ALICLOUD_ACCESS_KEY")
	ALICLOUD_SECRET_KEY := os.Getenv("ALICLOUD_SECRET_KEY")
	if len(ALICLOUD_ACCESS_KEY) == 0 || len(ALICLOUD_SECRET_KEY) == 0 {
		fmt.Println("Ak not set")
		return nil
	}

	path := fmt.Sprintf("/terraformResourceType/%s", resourceName)

	conn, err := NewIaCClient(ALICLOUD_ACCESS_KEY, ALICLOUD_SECRET_KEY, "", "iac.aliyuncs.com")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	deadline := time.Now().Add(time.Duration(2) * time.Second)
	retryCount := 0
	for {
		response, err := conn.DoRequest(tea.String("2021-08-06"), nil, tea.String("GET"), tea.String("AK"), tea.String(path), nil, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			// fmt.Println("[ERROR] invoking job failed:", err, ". Retry", retryCount, " times.")
			fmt.Println("[ERROR] invoking job failed:", resourceName)
		} else {
			if value, ok := response["body"]; ok {
				return value.(map[string]interface{})
			} else {
				return nil
			}
		}
		retryCount++
		if time.Now().After(deadline) {
			fmt.Println("timeout", err)
			return nil
		}
		time.Sleep(1 * time.Second)
	}

}

func NewIaCClient(accessKey, secretKey, securityToken, endpoint string) (*roa.Client, error) {
	credsConfig := &credential.Config{
		Type:            tea.String("access_key"),
		AccessKeyId:     tea.String(accessKey),
		AccessKeySecret: tea.String(secretKey),
	}
	if securityToken != "" {
		credsConfig.Type = tea.String("sts")
		credsConfig.SecurityToken = tea.String(securityToken)
	}
	credential, err := credential.NewCredential(credsConfig)
	if err != nil {
		return nil, err
	}
	config := roa.Config{
		Endpoint:       tea.String(endpoint),
		RegionId:       tea.String("cn-hangzhou"),
		UserAgent:      tea.String("AlibabaCloud (Linux; amd64) Java/1.8.0_152-b187 Core/4.6.2 HTTPClient/ApacheHttpClient"),
		Protocol:       tea.String("HTTPS"),
		ReadTimeout:    tea.Int(30000),
		ConnectTimeout: tea.Int(30000),
		MaxIdleConns:   tea.Int(500),
		Credential:     credential,
	}

	//if c.SourceIp != "" {
	//	config.SetSourceIp(c.SourceIp)
	//}
	client, err := roa.NewClient(&config)
	if err == nil {
		client.UserAgent = config.UserAgent
	}

	return client, nil
}

// merge the steps in one testcase
func mergeSteps(ossSchemaPath string, resourceTestCasesMap map[string]map[string]*ResourceTestExp) {
	var ossResourceSchema map[string]interface{}
	jsonFile, err := os.Open(ossSchemaPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &ossResourceSchema)
	if err != nil {
		fmt.Println("fail to unmarshal ossSchema")
	}

	for resourceType, resources := range resourceTestCasesMap {

		for resourceName, resourceTest := range resources {

			if strings.EqualFold(resourceName, "alicloud_nlb_listener") {
				fmt.Println()
			}

			// get the resource schema
			var resourceSchema map[string]*schema.Schema
			if resourceType == "resource" {
				resource, ok := resourceMap[resourceName]
				if !ok || resource == nil {
					log.Errorf("resource %s is not found in the provider ResourceMap\n\n", resourceName)
					continue
				}
				resourceSchema = resource.Schema
			} else {
				resource, ok := dataSourceMap[resourceName]
				if !ok || resource == nil {
					log.Errorf("data source %s is not found in the provider DataSourceMap\n\n", resourceName)
					continue
				}
				resourceSchema = resource.Schema
			}
			schemaAttributes := map[string]*SchemaAttr{}
			getResourceSchemaAttributes("", schemaAttributes, resourceSchema)

			if ossResourceSchema != nil {
				ossSchema := ossResourceSchema[resourceType].(map[string]interface{})
				if value, ok := ossSchema[resourceName]; ok {
					modifySchemaAttribute("", value.(map[string]interface{}), schemaAttributes)
				}
			}

			for k, v := range schemaAttributes {
				resourceTest.SchemaAttributes[k] = v
			}

			// merge steps
			for testCaseName, testCase := range resourceTest.TestCases {
				finalStep := []map[string]interface{}{}
				if len(testCase.Steps) == 1 {
					resourceTestCasesMap[resourceType][resourceName].
						TestCases[testCaseName].Tags = []map[string]interface{}{{}}
					continue
				} else {

					currentFinalStep := map[string]interface{}{}

					for _, step := range testCase.Steps {
						mergedStep := compareStep(step, currentFinalStep, schemaAttributes)
						// if the attributes are conflicted, split into two
						if mergedStep == nil {
							finalStep = append(finalStep, copyMap(currentFinalStep))
							currentFinalStep = copyMap(step)
						} else {
							currentFinalStep = mergedStep
						}
					}
					finalStep = append(finalStep, copyMap(currentFinalStep))

					resourceTestCasesMap[resourceType][resourceName].
						TestCases[testCaseName].Steps = finalStep
					resourceTestCasesMap[resourceType][resourceName].
						TestCases[testCaseName].Tags = make([]map[string]interface{}, len(finalStep))
				}
			}
		}
	}

}

func modifySchemaAttribute(rootName string, attributeMap map[string]interface{}, schemaAttributes map[string]*SchemaAttr) {

	for _, value := range attributeMap {
		attrValue := value.(map[string]interface{})
		tfKey := attrValue["terraformKey"].(string)
		if rootName != "" {
			tfKey = rootName + "." + tfKey
		}
		if _, ok := schemaAttributes[tfKey]; !ok {
			continue
		}
		// fmt.Println("attribute:", rootName, tfKey, reflect.TypeOf(value))
		if value, ok := attrValue["description"]; ok {
			schemaAttributes[tfKey].Description = value.(string)
		}
		if value, ok := attrValue["descriptionEn"]; ok {
			schemaAttributes[tfKey].DescriptionEn = value.(string)
		}
		_, schemaAttributes[tfKey].HasEnum = attrValue["enum"]
		if value, ok := attrValue["items"]; ok {
			if childMap, ok := value.(map[string]interface{})["properties"]; ok {
				modifySchemaAttribute(tfKey, childMap.(map[string]interface{}), schemaAttributes)
			}
		}
	}

}

func compareStep(currentStep, preStep map[string]interface{},
	schemaAttributes map[string]*SchemaAttr) map[string]interface{} {

	// get attributes of map and abandon the 'Removed' and 'Deprecated'
	currentAttrSet := mapset.NewSet()
	getAttributeSet(schemaAttributes, "", currentStep, currentAttrSet)

	preAttrSet := mapset.NewSet()
	getAttributeSet(schemaAttributes, "", preStep, preAttrSet)

	// one set can cover another
	if len(preStep) == 0 || preAttrSet.IsSubset(currentAttrSet) {
		return copyMap(currentStep)
	} else if len(currentStep) == 0 || currentAttrSet.IsSubset(preAttrSet) {
		return copyMap(preStep)
	}

	// some attributes appear in A but not in B, so as B
	currentDiff := currentAttrSet.Difference(preAttrSet)
	preDiff := preAttrSet.Difference(currentAttrSet)

	// Check for attribute conflicts
	conflicSet := mapset.NewSet()
	for attr := range currentDiff.Iter() {
		if attrStr, okString := attr.(string); okString {
			if value, okMap := schemaAttributes[attrStr]; okMap && value.ConflictsWith != nil {
				for _, v := range value.ConflictsWith {
					conflicSet.Add(v)
				}
			}
		}
	}
	for attr := range preDiff.Iter() {
		if conflicSet.Contains(attr) {
			return nil
		}
	}

	preStep = mergeMapNew("", preStep, currentStep, schemaAttributes)

	return copyMap(preStep)
}

func mergeMapNew(rootName string, mA, mB map[string]interface{},
	schemaAttributes map[string]*SchemaAttr) map[string]interface{} {
	for key, value := range mB {
		fullKey := key
		if rootName != "" {
			fullKey = rootName + "." + key
		}
		maxItems := -1
		if attr, ok := schemaAttributes[fullKey]; ok {
			maxItems = attr.MaxItems
		}
		if _, ok := mA[key]; !ok {
			mA[key] = value
		} else {
			if v, ok := value.(map[string]interface{}); ok {
				mA[key] = mergeMapNew(fullKey, mA[key].(map[string]interface{}), v, schemaAttributes)
			} else if v, ok := value.([]interface{}); ok {
				if sliceA, ok := mA[key].([]interface{}); ok {
					for _, vSliceB := range v {
						hasSame := false
						for _, vSliceA := range sliceA {
							if reflect.DeepEqual(vSliceA, vSliceB) {
								hasSame = true
							}
						}
						if maxItems != -1 && len(sliceA) >= maxItems {
							break
						} else if !hasSame {
							sliceA = append(sliceA, vSliceB)
						}
					}
					mA[key] = sliceA
				}
			}
		}
	}

	return mA
}

// merge diffB into mergedMap
func mergeMap(diffB, attrSetA mapset.Set, stepB, mergedMap map[string]interface{}) {
	for attr := range diffB.Iter() {
		attrStr, ok := attr.(string)
		if !ok {
			continue
		}
		attrSlice := strings.Split(attrStr, ".")
		root := ""
		path := ""
		for _, key := range attrSlice {
			if root != "" {
				path = root + "." + key
			} else {
				path = key
			}

			if attrSetA.Contains(path) {
				root = path
				continue
			} else {
				parentMap := findAttr(mergedMap, root)
				insertValue := findAttr(stepB, path)
				if parentMap == nil || insertValue == nil {
					continue
				}
				if data, ok := parentMap.(map[string]interface{}); ok {
					data[key] = insertValue
				}
				attrSetA.Add(path)
			}
			root = path
		}
	}
}

func findAttr(step map[string]interface{}, path string) interface{} {
	if path == "" {
		return step
	} else if !strings.Contains(path, ".") {
		return step[path]
	} else {
		firstKeyIndex := strings.Index(path, ".")
		if data, ok := step[path[:firstKeyIndex]].(map[string]interface{}); ok {
			return findAttr(data, path[firstKeyIndex+1:])
		}
	}
	return nil
}

func copyMap(old map[string]interface{}) map[string]interface{} {
	new := map[string]interface{}{}
	for k, v := range old {
		new[k] = v
	}

	return new
}

func getAttributeSet(schemaAttributes map[string]*SchemaAttr, rootName string,
	step interface{}, attibuteSet mapset.Set) {
	if d, ok := step.(map[string]interface{}); ok {
		for key, value := range d {
			attrKey := key
			if rootName != "" {
				attrKey = rootName + "." + attrKey
			}
			if strings.HasPrefix(attrKey, "tags.") ||
				resourceBlocksSet.Contains(attrKey) {
				continue
			}
			schema := schemaAttributes[attrKey]
			if schema != nil && len(schema.Deprecated) == 0 && len(schema.Removed) == 0 {
				attibuteSet.Add(attrKey)
				getAttributeSet(schemaAttributes, attrKey, value, attibuteSet)
			} else {
				delete(d, key)
			}
		}
	} else if d, ok := step.([]interface{}); ok {
		for _, v := range d {
			getAttributeSet(schemaAttributes, rootName, v, attibuteSet)
		}
	}

}

// merge the testcases' example of the resource
func mergeTestCases(resourceTestCasesMap map[string]map[string]*ResourceTestExp) {
	for resourceType, resources := range resourceTestCasesMap {

		for resourceName, resourceTest := range resources {

			// get the resource schema
			resourceSchema := resourceTest.SchemaAttributes

			enumForceNewAttrSet := mapset.NewSet()
			requiredAttrSet := mapset.NewSet()
			for k, v := range resourceSchema {
				if v.ForceNew && v.HasEnum {
					enumForceNewAttrSet.Add(k)
				}
				if v.Required {
					requiredAttrSet.Add(k)
				}
			}

			// hasMergedExp := true
			for {
				testCaseKey := []string{}
				for k := range resourceTest.TestCases {
					testCaseKey = append(testCaseKey, k)
				}
				if len(testCaseKey) == 1 {
					break
				}
				sort.Strings(testCaseKey)
				testCaseCount := len(testCaseKey)
				hasMerged := false
				for ti := 0; ti < testCaseCount; ti++ {
					for tj := ti + 1; tj < testCaseCount; tj++ {
						testCaseA := resourceTest.TestCases[testCaseKey[ti]]
						testCaseB := resourceTest.TestCases[testCaseKey[tj]]
						mergedTestCase := compareTestCase(testCaseA, testCaseB,
							resourceSchema, enumForceNewAttrSet)
						if mergedTestCase != nil {
							hasMerged = true
							resourceTest.TestCases[testCaseKey[ti]+" + "+testCaseKey[tj]] = mergedTestCase
						}
					}
				}

				for k, v := range resourceTest.TestCases {
					newSteps := []map[string]interface{}{}
					for _, s := range v.Steps {
						if s != nil {
							newSteps = append(newSteps, s)
						}
					}
					if len(newSteps) == 0 {
						delete(resourceTest.TestCases, k)
					} else {
						resourceTest.TestCases[k].Steps = newSteps
					}
				}
				if !hasMerged {
					break
				}
				// hasMergedExp = hasMerged
			}

			// if strings.EqualFold(resourceName, "alicloud_mongodb_instance") {
			// 	fmt.Println()
			// }

			if *enableDocs {
				extractDocExample(resourceTest)
				resourceTestCasesMap[resourceType][resourceName] = resourceTest
			}
		}
	}
}

func compareTestCase(testCaseA, testCaseB *TestCase,
	schemaAttributes map[string]*SchemaAttr, enumForceNewAttrSet mapset.Set) *TestCase {

	testCaseA.Tags = make([]map[string]interface{}, len(testCaseA.Steps))
	for i := range testCaseA.Tags {
		testCaseA.Tags[i] = map[string]interface{}{}
	}
	testCaseB.Tags = make([]map[string]interface{}, len(testCaseB.Steps))
	for i := range testCaseB.Tags {
		testCaseB.Tags[i] = map[string]interface{}{}
	}

	mergedTestCase := []map[string]interface{}{}
	mergedTag := []map[string]interface{}{}
	for si, stepA := range testCaseA.Steps {
		for sj, stepB := range testCaseB.Steps {

			if stepA == nil || stepB == nil {
				continue
			}

			attrSetA := mapset.NewSet()
			getAttributeSet(schemaAttributes, "", stepA, attrSetA)
			attrSetB := mapset.NewSet()
			getAttributeSet(schemaAttributes, "", stepB, attrSetB)

			isDiffForceNew := false
			forceNewSet := attrSetA.Intersect(attrSetB).Intersect(enumForceNewAttrSet)
			for attr := range forceNewSet.Iter() {
				if attrStr, okString := attr.(string); okString {
					valueA := findAttr(stepA, attrStr)
					valueB := findAttr(stepB, attrStr)
					if !reflect.DeepEqual(valueA, valueB) {
						testCaseA.Tags[si][attrStr] = valueA.(string)
						testCaseB.Tags[sj][attrStr] = valueB.(string)
						isDiffForceNew = true
						// fmt.Println(testCaseA.PrimaryConfig)
						// break
					}
				}
			}
			if isDiffForceNew {
				continue
			}

			// one set can cover another
			if len(stepA) == 0 || attrSetA.IsSubset(attrSetB) {
				testCaseA.Steps[si] = nil
			} else if len(stepB) == 0 || attrSetB.IsSubset(attrSetA) {
				testCaseB.Steps[sj] = nil
			}

			// some attributes appear in A but not in B, so as B
			diffA := attrSetA.Difference(attrSetB)
			diffB := attrSetB.Difference(attrSetA)

			// Check for attribute conflicts
			conflicSet := mapset.NewSet()
			for attr := range diffA.Iter() {
				if attrStr, okString := attr.(string); okString {
					if value, okMap := schemaAttributes[attrStr]; okMap && value.ConflictsWith != nil {
						for _, v := range value.ConflictsWith {
							conflicSet.Add(v)
						}
					}
				}
			}
			isConflict := false
			for attr := range diffB.Iter() {
				if conflicSet.Contains(attr) {
					isConflict = true
					break
				}
			}
			if isConflict {
				continue
			}

			// merge the map
			// mergedMap := copyMap(stepA)
			// mergeMap(diffB, attrSetA, stepB, mergedMap)
			mergedMap := mergeMapNew("", stepA, stepB, schemaAttributes)

			tags := copyMap(testCaseA.Tags[si])
			for k, v := range testCaseB.Tags[sj] {
				tags[k] = v
			}

			mergedTestCase = append(mergedTestCase, mergedMap)
			mergedTag = append(mergedTag, tags)
			testCaseA.Steps[si] = nil
			testCaseB.Steps[sj] = nil
		}
	}

	if len(mergedTestCase) == 0 {
		return nil
	}

	finalDependences := mergeDependence(testCaseA.Dependences, testCaseB.Dependences)

	return &TestCase{
		Dependences:   finalDependences,
		PrimaryConfig: testCaseA.PrimaryConfig,
		Steps:         mergedTestCase,
		Tags:          mergedTag,
	}
}

func mergeDependence(dependenceA, dependenceB string) string {
	if len(dependenceA) == 0 {
		return dependenceB
	} else if len(dependenceB) == 0 {
		return dependenceA
	}

	attrA, attrB := parseHCL(dependenceA), parseHCL(dependenceB)

	if attrA.Equal(attrB) {
		return dependenceA
	} else if attrA.IsProperSubset(attrB) {
		return dependenceB
	} else if attrB.IsProperSubset(attrA) {
		return dependenceA
	} else {
		return dependenceA + dependenceB
	}

}

// parse HCL string and return the primaryConfig set
func parseHCL(hclStr string) mapset.Set {
	hclSet := mapset.NewSet()

	file, diags := hclsyntax.ParseConfig([]byte(hclStr), "example.hcl", hcl.InitialPos)
	if diags.HasErrors() {
		return hclSet
	}
	blocks := file.Body.(*hclsyntax.Body).Blocks

	for _, block := range blocks {
		name := block.Type
		for _, v := range block.Labels {
			name += " " + v
		}
		hclSet.Add(name)
	}

	return hclSet
}

func extractDocExample(resourceTest *ResourceTestExp) {
	allExpAttr := mapset.NewSet()

	minCoverage := 1.1
	docExampleTestCaseName := ""
	docExampleIndex := 0

	for testCaseName, testCase := range resourceTest.TestCases {
		for j, exp := range testCase.Steps {
			attrSet := mapset.NewSet()
			getAttributeSet(resourceTest.SchemaAttributes, "", exp, attrSet)
			for v := range attrSet.Iter() {
				allExpAttr.Add(v)
			}

			coverage := (float64)(len(attrSet.ToSlice())) /
				(float64)(len(resourceTest.SchemaAttributes))
			if coverage < minCoverage {
				docExampleTestCaseName = testCaseName
				docExampleIndex = j
				minCoverage = coverage
			}

		}
	}

	// fmt.Println()

	resourceTest.TestCases["docExample"] = &TestCase{
		Dependences:   resourceTest.TestCases[docExampleTestCaseName].Dependences,
		PrimaryConfig: resourceTest.TestCases[docExampleTestCaseName].PrimaryConfig,
		Steps:         []map[string]interface{}{copyMap(resourceTest.TestCases[docExampleTestCaseName].Steps[docExampleIndex])},
		Tags:          []map[string]interface{}{map[string]interface{}{}},
	}

	resourceTest.TestCases[docExampleTestCaseName].Steps[docExampleIndex] = nil

	coverage := (float64)(len(allExpAttr.ToSlice())) /
		(float64)(len(resourceTest.SchemaAttributes))

	resourceTest.Coverage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", coverage), 64)

}

func extractExp(resourceTest *ResourceTestExp, requiredAttrSet mapset.Set) {

	allExpAttr := mapset.NewSet()

	for testCaseName, testCase := range resourceTest.TestCases {

		if strings.Contains(testCaseName, "docExample") {
			continue
		}

		docExamples := &TestCase{
			Dependences:   testCase.Dependences,
			PrimaryConfig: testCase.PrimaryConfig,
			Steps:         []map[string]interface{}{},
			Tags:          []map[string]interface{}{},
		}

		for j, exp := range testCase.Steps {
			attrSet := mapset.NewSet()
			getAttributeSet(resourceTest.SchemaAttributes, "", exp, attrSet)
			for v := range attrSet.Iter() {
				allExpAttr.Add(v)
			}

			if requiredAttrSet.Equal(attrSet) {
				docExamples.Steps = append(docExamples.Steps, exp)
				resourceTest.TestCases[testCaseName].Steps[j] = nil

			} else if requiredAttrSet.IsProperSubset(attrSet) {
				hasSet := mapset.NewSet()
				requiredMap := extractBasic("", requiredAttrSet, hasSet, exp)
				if len(hasSet.ToSlice()) == 0 {
					continue
				}
				docExamples.Steps = append(docExamples.Steps, requiredMap)

				// needToExact = true
				// baseMap = exp
				// expSetCount = float64(len(attrSet.ToSlice()))
				// baseTestCase = testCase
			}

		}

		if len(docExamples.Steps) != 0 {
			docExamples.Tags = make([]map[string]interface{}, len(docExamples.Steps))
			resourceTest.TestCases[testCaseName+"docExample"] = docExamples
		}
	}

	coverage := (float64)(len(allExpAttr.ToSlice())) /
		(float64)(len(resourceTest.SchemaAttributes))

	resourceTest.Coverage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", coverage), 64)

}

func extractBasic(rootName string, requiredAttrSet, hasSet mapset.Set, baseMap map[string]interface{}) map[string]interface{} {
	requiredMap := map[string]interface{}{}
	for key, value := range baseMap {
		fullName := key
		if rootName != "" {
			fullName = rootName + "." + key
		}
		if requiredAttrSet.Contains(fullName) || resourceBlocksSet.Contains(fullName) {
			hasSet.Add(fullName)
			if baseValueMap, ok := value.(map[string]interface{}); ok {
				requiredMap[key] = extractBasic(fullName, requiredAttrSet, hasSet, baseValueMap)
			} else {
				requiredMap[key] = baseMap[key]
			}
		}
	}

	return requiredMap
}

func getResourceSchemaAttributes(rootName string, schemaAttributes map[string]*SchemaAttr,
	resourceSchema map[string]*schema.Schema) {
	for key, value := range resourceSchema {
		if rootName != "" {
			key = rootName + "." + key
		}

		if _, ok := schemaAttributes[key]; !ok {
			schemeAttr := &SchemaAttr{
				Name:            key,
				Type:            value.Type.String(),
				MinItems:        value.MinItems,
				MaxItems:        value.MaxItems,
				Required:        value.Required,
				ForceNew:        value.ForceNew,
				HasValidateFunc: false,
				Deprecated:      value.Deprecated,
				Removed:         value.Removed,
				ConflictsWith:   value.ConflictsWith,
				Description:     "",
				DescriptionEn:   "",
				HasEnum:         false,
			}
			if value.ValidateFunc != nil {
				// funcType := reflect.ValueOf(&value.ValidateFunc).Elem()
				// px := funcType.Addr().Interface()
				// fmt.Println(reflect.TypeOf(px))
				schemeAttr.HasValidateFunc = true
			}
			schemaAttributes[key] = schemeAttr
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
				getResourceSchemaAttributes(key, schemaAttributes, value.Elem.(*schema.Resource).Schema)
			}
		}
	}
}

type SchemaAttr struct {
	Name            string
	Type            string
	Required        bool
	ForceNew        bool
	HasValidateFunc bool
	Deprecated      string
	Removed         string
	ElemType        string
	MinItems        int
	MaxItems        int
	ConflictsWith   []string
	Description     string
	DescriptionEn   string
	HasEnum         bool
}

type ResourceTestExp struct {
	SchemaAttributes map[string]*SchemaAttr
	TestCases        map[string]*TestCase
	Examples         map[string]string
	ExpDescriptions  map[string]string
	Coverage         float64
}

type TestCase struct {
	Dependences   string
	PrimaryConfig string
	Steps         []map[string]interface{}
	Tags          []map[string]interface{}
}

func processExp(resourceTestCasesMap map[string]map[string]*ResourceTestExp) {
	for resourceType, resources := range resourceTestCasesMap {
		for resourceName, resourceTest := range resources {

			// if strings.EqualFold(resourceName, "alicloud_ack_one_cluster") {
			// 	fmt.Println()
			// }

			// fmt.Println(resourceName)

			nonForceNewAttr := mapset.NewSet()
			for k, v := range resourceTest.SchemaAttributes {
				if !v.ForceNew {
					nonForceNewAttr.Add(k)
				}
			}

			// remove the redundant dependence and extract the variables
			for testCaseName, testCase := range resourceTest.TestCases {

				// remove the redundant dependence
				dependenceStr := testCase.Dependences
				// fmt.Println(dependenceStr)
				dependenceSet := mapset.NewSet()
				dpBlock := map[string]*hclwrite.Block{}
				file, diags := hclwrite.ParseConfig([]byte(dependenceStr), "example.hcl", hcl.InitialPos)
				if diags.HasErrors() {
					continue
				}
				blocks := file.Body().Blocks()

				newFile := hclwrite.NewEmptyFile()
				newBody := newFile.Body()
				for _, block := range blocks {
					rName := block.Type()
					quoteName := ""
					if block.Type() == "variable" {
						quoteName += "var"
					} else if block.Type() == "locals" {
						quoteName += "local"
					} else if block.Type() == "resource" {
						quoteName += ""
					} else if block.Type() == "data" {
						quoteName += ""
					}

					for _, v := range block.Labels() {
						rName += " " + v
						quoteName += "." + v
					}

					quoteName = strings.Trim(quoteName, ".")

					if !dependenceSet.Contains(rName) ||
						strings.EqualFold(rName, "locals") {
						dependenceSet.Add(rName)
						newBody.AppendBlock(block)
						newBody.AppendNewline()
						dpBlock[quoteName] = block
						if block.Type() == "locals" {
							for k := range block.Body().Attributes() {
								dpBlock[quoteName+"."+k] = block
							}
							delete(dpBlock, quoteName)
						}
					}
				}

				// if resourceName == "alicloud_cms_namespace" {
				// 	fmt.Println()
				// }

				zonesRegex := regexp.MustCompile("^alicloud(.*)_zones$")
				for j, step := range testCase.Steps {

					if step == nil {
						continue
					}

					privateDPBlock := map[string]*hclwrite.Block{}
					for k, v := range dpBlock {
						privateDPBlock[k] = v
					}

					// remove redundant dependence of "required" example
					for {
						resourceStr := string(newFile.Bytes()) + testCase.PrimaryConfig +
							fmt.Sprint(alicloud.ValueConvert(0, reflect.ValueOf(step)))
						hasNotQuote := false
						for k, v := range privateDPBlock {
							if !strings.Contains(resourceStr, k) {
								if strings.HasPrefix(k, "local") {
									attr := strings.Split(k, ".")[1]
									v.Body().RemoveAttribute(attr)
									delete(privateDPBlock, k)
									if len(v.Body().Attributes()) == 0 {
										newBody.RemoveBlock(v)
									}
								} else {
									hasNotQuote = true
									newBody.RemoveBlock(v)
									delete(privateDPBlock, k)
									// rName := strings.Split(k, ".")
									// delete(exampleDataSourceDependence, rName[0])
								}
							}
						}
						if !hasNotQuote {
							break
						}
					}

					// extract Requirements
					dataSourceDependence := map[string]struct{}{}
					newBlocks := newBody.Blocks()
					for _, block := range newBlocks {
						if block.Type() == "data" {
							dataSourceDependence[block.Labels()[0]] = struct{}{}
						}
					}
					for k := range dataSourceDependence {
						if match := zonesRegex.MatchString(k); match {
							delete(dataSourceDependence, k)
						}
					}
					dataSourceDescription := ""
					if len(dataSourceDependence) > 0 {
						dataSourceDescription = "\n## Requirements\n  \n  Before using this example, you first need to create the following dependency resources.\n"
						for k := range dataSourceDependence {
							dataSourceDescription += fmt.Sprintf("- `%s`\n", k)
						}
					}

					// modify variable dependence
					for k, v := range privateDPBlock {
						if strings.Contains(k, "var.") {
							varBody := v.Body()
							varBody.SetAttributeValue("description", cty.StringVal("This variable can be used in all resources in this example."))
						}
					}

					// extract the vairable
					varMap := map[string]string{}
					for k, v := range step {
						if vStr, ok := v.(string); ok && nonForceNewAttr.Contains(k) &&
							!strings.HasPrefix(vStr, "${") {
							step[k] = "var." + k
							vStr = strings.TrimSuffix(vStr, "-update")
							vStr = strings.TrimSuffix(vStr, "_update")
							vStr = strings.TrimSuffix(vStr, "update")
							varMap[k] = vStr
						}
					}

					varBlocks := []*hclwrite.Block{}
					for varKey, varValue := range varMap {
						varBlock := newBody.AppendNewBlock("variable", []string{varKey})
						varBody := varBlock.Body()
						if value, ok := resourceTest.SchemaAttributes[varKey]; ok {
							if len(value.DescriptionEn) != 0 {
								varBody.SetAttributeValue("description", cty.StringVal(value.DescriptionEn))
							} else {
								varBody.SetAttributeValue("description", cty.StringVal("This variable can be used in all resources in this example."))
							}
							typeValue, defaultValue := getHclType(value.Type, varValue)
							varBody.SetAttributeRaw("type", typeValue)
							varBody.SetAttributeValue("default", defaultValue)
						} else {
							varBody.SetAttributeValue("default", cty.StringVal(varValue))
						}

						varBlocks = append(varBlocks, varBlock)
					}

					// get the final example
					exp := string(newFile.Bytes()) + "\n" + testCase.PrimaryConfig +
						fmt.Sprint(alicloud.ValueConvert(0, reflect.ValueOf(step)))

					for _, block := range varBlocks {
						newBody.RemoveBlock(block)
					}
					rName := strings.ReplaceAll(resourceName, "_", "-")

					suffix := "-complete"
					if strings.EqualFold(testCaseName, "docExample") {
						suffix = "-docExample"
					}
					expName := rName + suffix
					forceNewDescription := ""
					if testCase.Tags[j] != nil && len(testCase.Tags[j]) > 0 {
						tagKey := []string{}
						tagSlice := []string{}
						for k := range testCase.Tags[j] {
							tagKey = append(tagKey, k)
						}
						sort.Strings(tagKey)
						for _, k := range tagKey {
							v := testCase.Tags[j][k]
							tagSlice = append(tagSlice, v.(string))
							forceNewDescription += fmt.Sprintf("- `%s` = \"%s\"\n", k, v)
						}

						tags := strings.Join(tagSlice, "-")
						expName += "-" + tags
					}

					expDescription := fmt.Sprintf("\n" +
						"## Introduction\n" +
						"\n")

					if forceNewDescription == "" {
						expDescription += fmt.Sprintf(
							"This example is used to create a `%s` resource.\n", resourceName)
					} else {
						expDescription += fmt.Sprintf(
							"This example is used to create a `%s` resource, where\n%s"+
								"\n", resourceName, forceNewDescription)
					}
					if dataSourceDescription != "" {
						expDescription += fmt.Sprintf(
							"%s"+
								"\n", dataSourceDescription)
					}
					resourceTest.ExpDescriptions[expName] = expDescription
					resourceTest.Examples[expName] = exp
				}

			}
			resourceTestCasesMap[resourceType][resourceName] = resourceTest
		}
	}
}

func getHclType(schemaType string, value string) (hclType hclwrite.Tokens, val cty.Value) {
	switch schemaType {
	case "TypeBool":
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return hclwrite.TokensForIdentifier("bool"), cty.BoolVal(boolValue)
		}
	case "TypeInt":
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return hclwrite.TokensForIdentifier("number"), cty.NumberIntVal(intValue)
		}
	case "TypeFloat":
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return hclwrite.TokensForIdentifier("number"), cty.NumberFloatVal(floatValue)
		}
	case "TypeString":
		return hclwrite.TokensForIdentifier("string"), cty.StringVal(value)
	default:
		return nil, cty.Value{}

	}
	return hclwrite.TokensForIdentifier(schemaType), cty.StringVal(value)

}

// distinguish the resource type and print to different file
func printExpToFile(resourceTestCasesMap map[string]map[string]*ResourceTestExp) {

	count := 0
	testAccRegex := regexp.MustCompile(`".*(?i)(testacc|test-acc).*"`)
	tfexpRegex := regexp.MustCompile(`"tf-example.*"`)

	for resourceType, resources := range resourceTestCasesMap {
		for _, resourceTest := range resources {
			for key, text := range resourceTest.Examples {

				expDescription := resourceTest.ExpDescriptions[key]

				key = strings.Replace(key, "alicloud", "101", 1)

				// if !strings.EqualFold(key, "101-nlb-listener-additional-certificate-attachment-complete") {
				// 	continue
				// }

				text = testAccRegex.ReplaceAllString(text, "\"tf-example\"")
				text = tfexpRegex.ReplaceAllString(text, "\"tf-example\"")

				examplePath := *examplePath + "/" + resourceType +
					"/" + key + "/"

				_, err := os.Stat(examplePath)
				if err != nil || os.IsNotExist(err) {
					err := os.MkdirAll(examplePath, 0777)
					if err != nil {
						return
					}
				}

				file, diags := hclwrite.ParseConfig([]byte(text), "example.hcl", hcl.InitialPos)
				if diags.HasErrors() {
					fmt.Println(key)
					fmt.Println(text)
					count++
					continue
				}
				blocks := file.Body().Blocks()

				variableFile := hclwrite.NewEmptyFile()
				variableBody := variableFile.Body()

				resourceFile := hclwrite.NewEmptyFile()
				resourceBody := resourceFile.Body()

				for _, block := range blocks {
					typeName := block.Type()
					if typeName == "variable" {
						variableBody.AppendBlock(block)
						variableBody.AppendNewline()
					} else {
						resourceBody.AppendBlock(block)
						resourceBody.AppendNewline()
					}
				}

				providerFile, err := os.OpenFile(examplePath+"/provider.tf",
					os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
				if err != nil {
					continue
				}
				providerFile.WriteString(alicloudProvider)
				providerFile.Close()

				if len(variableBody.Blocks()) != 0 {
					expVarFile, err := os.OpenFile(examplePath+"/variable.tf",
						os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
					if err != nil {
						continue
					}
					expVarFile.WriteString(string(variableFile.Bytes()))
					expVarFile.Close()
				}

				if len(resourceBody.Blocks()) != 0 {
					expFile, err := os.OpenFile(examplePath+"/main.tf",
						os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
					if err != nil {
						continue
					}
					expFile.WriteString(string(resourceFile.Bytes()))
					expFile.Close()
				}

				expFile, err := os.OpenFile(examplePath+"/header.md",
					os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
				if err != nil {
					continue
				}
				expFile.WriteString(expDescription)
				expFile.Close()

				os.Create(examplePath + "/footer.md")

				dir, _ := os.ReadDir(examplePath)
				if len(dir) == 0 {
					os.Remove(examplePath)
				}

			}
		}
	}
	fmt.Println(count)
}
