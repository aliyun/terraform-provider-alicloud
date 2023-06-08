package scripts

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	set "github.com/deckarep/golang-set"
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
	resourceNames = flag.String("resourceNames", "", "the names of the terraform resources to diff")
	fileNames     = flag.String("fileNames", "", "the files to check diff")
	filterList    = map[string][]string{
		"alicloud_amqp_instance":            {"logistics"},
		"alicloud_cms_alarm":                {"notify_type"},
		"alicloud_cs_serverless_kubernetes": {"private_zone", "create_v2_cluster"},
		"alicloud_slb_listener":             {"lb_protocol", "instance_port", "lb_port"},
		"alicloud_kvstore_instance":         {"connection_string"},
		"alicloud_instance":                 {"subnet_id"},
		"alicloud_hbr_ots_backup_plan":      {"vault_id"},
		"alicloud_nat_gateway":              {"vswitch_id"},
		"alicloud_ecs_disk":                 {"advanced_features", "encrypt_algorithm", "dedicated_block_storage_cluster_id"},
	}
)

type Resource struct {
	Name       string
	Arguments  map[string]interface{}
	Attributes map[string]interface{}
}

type ResourceAttribute struct {
	Name        string
	Type        string
	Optional    string
	Required    string
	ForceNew    bool
	Default     string
	ElemType    string
	Deprecated  string
	DocsLineNum int
}

func TestConsistencyWithDocument(t *testing.T) {
	exitCode := 0
	flag.Parse()
	if fileNames != nil && len(*fileNames) == 0 {
		log.Infof("the diff file is empty, shipped!")
		return
	}

	byt, _ := ioutil.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	//fileRegex := regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*.go")
	//fileTestRegex := regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*_test.go")
	//fileDocsRegex := regexp.MustCompile("website/docs/(d|r)/[0-9a-zA-Z_]*.html.markdown")
	fileRegex := regexp.MustCompile("alicloud/(resource)[0-9a-zA-Z_]*.go")
	fileTestRegex := regexp.MustCompile("alicloud/(resource)[0-9a-zA-Z_]*_test.go")
	fileDocsRegex := regexp.MustCompile("website/docs/(r)/[0-9a-zA-Z_]*.html.markdown")
	resourceNameMap := make(map[string]struct{})
	for _, file := range diff.Files {
		resourceName := ""
		if fileRegex.MatchString(file.NewName) {
			if fileTestRegex.MatchString(file.NewName) {
				continue
			}
			resourceName = strings.TrimPrefix(strings.TrimSuffix(strings.Split(file.NewName, "/")[1], ".go"), "resource_")
		} else if fileDocsRegex.MatchString(file.NewName) {
			resourceName = "alicloud_" + strings.TrimSuffix(strings.Split(file.NewName, "/")[3], ".html.markdown")
		} else {
			continue
		}
		if _, ok := resourceNameMap[resourceName]; ok {
			continue
		} else {
			resourceNameMap[resourceName] = struct{}{}
		}

		log.Infof("==> Checking resource %s attributes consistency...", resourceName)
		resource, ok := alicloud.Provider().(*schema.Provider).ResourcesMap[resourceName]
		if !ok || resource == nil {
			//resourceName = strings.TrimPrefix(resourceName, "data_source_")
			//resource, ok = alicloud.Provider().(*schema.Provider).DataSourcesMap[resourceName]
			//if !ok || resource == nil {
			log.Errorf("resource %s is not found in the provider ResourceMap\n\n", resourceName)
			exitCode = 1
			continue
			//}
		}
		resourceSchema := resource.Schema
		resourceSchemaFromDocs := make(map[string]ResourceAttribute)
		if err := parseResourceDocs(resourceName, resourceSchemaFromDocs); err != nil {
			fmt.Println(err)
			t.Fatal()
		}

		if consistencyCheck(t, resourceName, resourceSchemaFromDocs, resourceSchema) {
			log.Infof("--- PASS!\n\n")
			continue
		}
		log.Errorf("--- Failed!\n\n")
		exitCode = 1
	}
	if exitCode > 0 {
		os.Exit(exitCode)
	}
	return
}

func TestFieldCompatibilityCheck(t *testing.T) {
	flag.Parse()
	if fileNames != nil && len(*fileNames) == 0 {
		log.Warningf("the diff file is empty")
		return
	}
	byt, _ := ioutil.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	res := false
	fileRegex := regexp.MustCompile("alicloud/resource[0-9a-zA-Z_]*.go")
	fileTestRegex := regexp.MustCompile("alicloud/resource[0-9a-zA-Z_]*_test.go")
	for _, file := range diff.Files {
		if fileRegex.MatchString(file.NewName) {
			if fileTestRegex.MatchString(file.NewName) {
				continue
			}
			log.Debugf("checking compatibility of the file %s", file.NewName)
			for _, hunk := range file.Hunks {
				if hunk != nil {
					prev := ParseField(hunk.OrigRange, hunk.OrigRange.Length)
					current := ParseField(hunk.NewRange, hunk.NewRange.Length)
					res = CompatibilityRule(prev, current, file.NewName)
				}
			}
		}
	}
	if res {
		t.Fatal("incompatible changes occurred")
		os.Exit(1)
	}
}

func CompatibilityRule(prev, current map[string]map[string]interface{}, fileName string) (res bool) {
	for filedName, obj := range prev {
		// Optional -> Required
		_, exist1 := obj["Optional"]
		_, exist2 := current[filedName]["Required"]
		if exist1 && exist2 {
			res = true
			log.Errorf("[Incompatible Change]: there should not change attribute %v to required from optional in the file %v!", fileName, filedName)
		}
		// Type changed
		typPrev, exist1 := obj["Type"]
		typCurr, exist2 := current[filedName]["Type"]

		if exist1 && exist2 && typPrev != typCurr {
			res = true
			log.Errorf("[Incompatible Change]: there should not to change the type of attribute %v in the file %v!", fileName, filedName)
		}

		_, exist1 = obj["ForceNew"]
		_, exist2 = current[filedName]["ForceNew"]
		if !exist1 && exist2 {
			res = true
			log.Errorf("[Incompatible Change]: there should not to change attribute %v to ForceNew from normal in the file %v!", fileName, filedName)
		}

		// type string: valid values
		validateArrPrev, exist1 := obj["ValidateFuncString"]
		validateArrCurrent, exist2 := current[filedName]["ValidateFuncString"]
		if exist1 && exist2 && len(validateArrPrev.([]string)) > len(validateArrCurrent.([]string)) {
			res = true
			log.Errorf("[Incompatible Change]: attribute %v enum values should not less than before in the file %v!", fileName, filedName)
		}

	}
	return
}

func ParseField(hunk diffparser.DiffRange, length int) map[string]map[string]interface{} {
	schemaRegex := regexp.MustCompile("^\\t*\"([a-zA-Z_]*)\"")
	typeRegex := regexp.MustCompile("^\\t*Type:\\s+schema.([a-zA-Z]*)")
	optionRegex := regexp.MustCompile("^\\t*Optional:\\s+([a-z]*),")
	forceNewRegex := regexp.MustCompile("^\\t*ForceNew:\\s+([a-z]*),")
	requiredRegex := regexp.MustCompile("^\\t*Required:\\s+([a-z]*),")
	validateStringRegex := regexp.MustCompile("^\\t*ValidateFunc: ?validation.StringInSlice\\(\\[\\]string\\{([a-z\\-A-Z_,\"\\s]*)")

	temp := map[string]interface{}{}
	schemaName := ""
	raw := make(map[string]map[string]interface{}, 0)
	for i := 0; i < length; i++ {
		currentLine := hunk.Lines[i]
		content := currentLine.Content
		fieldNameMatched := schemaRegex.FindAllStringSubmatch(content, -1)
		if fieldNameMatched != nil && fieldNameMatched[0] != nil {
			if len(schemaName) != 0 && schemaName != fieldNameMatched[0][1] {
				temp["Name"] = schemaName
				raw[schemaName] = temp
				temp = map[string]interface{}{}
			}
			schemaName = fieldNameMatched[0][1]
		}

		if !schemaRegex.MatchString(currentLine.Content) && currentLine.Mode == diffparser.UNCHANGED {
			continue
		}

		typeMatched := typeRegex.FindAllStringSubmatch(content, -1)
		typeValue := ""
		if typeMatched != nil && typeMatched[0] != nil {
			typeValue = typeMatched[0][1]
			temp["Type"] = typeValue
		}

		optionalMatched := optionRegex.FindAllStringSubmatch(content, -1)
		optionValue := ""
		if optionalMatched != nil && optionalMatched[0] != nil {
			optionValue = optionalMatched[0][1]
			op, _ := strconv.ParseBool(optionValue)
			temp["Optional"] = op
		}

		forceNewMatched := forceNewRegex.FindAllStringSubmatch(content, -1)
		forceNewValue := ""
		if forceNewMatched != nil && forceNewMatched[0] != nil {
			forceNewValue = forceNewMatched[0][1]
			fc, _ := strconv.ParseBool(forceNewValue)
			temp["ForceNew"] = fc
		}

		requiredMatched := requiredRegex.FindAllStringSubmatch(content, -1)
		requiredValue := ""
		if requiredMatched != nil && requiredMatched[0] != nil {
			requiredValue = requiredMatched[0][1]
			rq, _ := strconv.ParseBool(requiredValue)
			temp["Required"] = rq
		}

		validateStringMatched := validateStringRegex.FindAllStringSubmatch(content, -1)
		validateStringValue := ""
		if validateStringMatched != nil && validateStringMatched[0] != nil {
			validateStringValue = validateStringMatched[0][1]
			temp["ValidateFuncString"] = strings.Split(validateStringValue, ",")
		}

	}
	if _, exist := raw[schemaName]; !exist && len(temp) >= 1 {
		temp["Name"] = schemaName
		raw[schemaName] = temp
	}
	return raw
}

func parseResourceDocs(resourceName string, resourceAttributes map[string]ResourceAttribute) error {
	splitRes := strings.Split(resourceName, "alicloud_")
	if len(splitRes) < 2 {
		log.Errorf("parsing resource name %s failed.", resourceName)
		return fmt.Errorf(fmt.Sprintf("parsing resource name %s failed.", resourceName))
	}
	basePath := "../website/docs/r/"
	filePath := strings.Join([]string{basePath, splitRes[1], ".html.markdown"}, "")

	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("open resource %s docs failed. Error: %s", filePath, err)
		return err
	}
	defer file.Close()

	argsRegex := regexp.MustCompile("## Argument Reference")
	attribRegex := regexp.MustCompile("## Attributes Reference")
	secondLevelRegex := regexp.MustCompile("^### `([a-zA-Z_\\.0-9]*)`")
	argumentsFieldRegex := regexp.MustCompile("^\\* `([a-zA-Z_0-9]*)`[ ]*-? ?(\\(.*\\)) ?(.*)")
	attributeFieldRegex := regexp.MustCompile("^\\* `([a-zA-Z_0-9]*)`[ ]*-?(.*)")

	name := filepath.Base(filePath)
	re := regexp.MustCompile("[a-z0-9A-Z_]*")
	resourceName = "alicloud_" + re.FindString(name)
	//result := &Resource{Name: resourceName, Arguments: map[string]interface{}{}, Attributes: map[string]interface{}{}}
	//log.Infof("the resourceName = %s\n", resourceName)

	scanner := bufio.NewScanner(file)
	phase := "Argument"
	record := false
	subAttributeName := ""
	line := 0
	for scanner.Scan() {
		line += 1
		text := scanner.Text()
		if strings.HasPrefix(text, "#") && (strings.HasSuffix(text, "Timeouts") || strings.HasSuffix(text, "Import")) {
			break
		}
		if argsRegex.MatchString(text) {
			record = true
			phase = "Argument"
			continue
		}
		if attribRegex.MatchString(text) {
			record = true
			phase = "Attribute"
			subAttributeName = ""
			continue
		}
		if secondLevelRegex.MatchString(text) {
			record = true
			parts := strings.Split(text, " ")
			subAttributeName = strings.Trim(parts[len(parts)-1], "`")
			continue
		}
		if record {
			var matched [][]string
			if phase == "Argument" {
				matched = argumentsFieldRegex.FindAllStringSubmatch(text, 1)
			} else if phase == "Attribute" {
				matched = attributeFieldRegex.FindAllStringSubmatch(text, 1)
			}

			for _, m := range matched {
				attribute := parseMatchLine(m, phase, subAttributeName)
				if attribute == nil {
					continue
				}
				attribute.DocsLineNum = line
				resourceAttributes[attribute.Name] = *attribute
			}
		}
	}
	return nil
}

func parseMatchLine(words []string, phase, rootName string) *ResourceAttribute {
	result := ResourceAttribute{}
	if phase == "Argument" && len(words) >= 4 {
		if rootName != "" {
			result.Name = rootName + "." + words[1]
		} else {
			result.Name = words[1]
		}
		//result["Description"] = words[3]
		if strings.Contains(words[2], "Optional") {
			result.Optional = "true"
		}
		if strings.Contains(words[2], "Required") {
			result.Required = "true"
		}
		if strings.Contains(words[2], "ForceNew") {
			result.ForceNew = true
		}
		if strings.Contains(words[2], "Deprecated") {
			result.Deprecated = "Deprecated since"
		}
		return &result
	}
	if phase == "Attribute" && len(words) >= 3 {
		if words[1] == "id" {
			return nil
		}
		if rootName != "" {
			result.Name = rootName + "." + words[1]
		} else {
			result.Name = words[1]
		}
		//result["Description"] = words[2]
		return &result
	}
	return nil
}

func consistencyCheck(t *testing.T, resourceName string, resourceAttributeFromDocs map[string]ResourceAttribute, resourceSchemaDefined map[string]*schema.Schema) bool {
	isConsistent := true
	filteredList := set.NewSet()
	if val, ok := filterList[resourceName]; ok {
		for _, v := range val {
			filteredList.Add(v)
		}
	}
	//defer func() {
	//	if r := recover(); r != nil {
	//		res = true
	//		log.Errorf("internal error: Please email terraform@alibabacloud.com to report the issue with the related resource")
	//		t.Fatal()
	//	}
	//}()

	// the number of the schema field + 1(id) should equal to the number defined in document
	resourceAttributes := make(map[string]ResourceAttribute)
	getResourceAttributes("", resourceAttributes, resourceSchemaDefined)

	for attributeKey, attributeValue := range resourceAttributes {
		attributeDocsValue, ok := resourceAttributeFromDocs[attributeKey]
		if !ok {
			isConsistent = false
			log.Errorf("'%v' which described in the docs not found in the resource schema", attributeKey)
		}
		if attributeValue.Optional == "true" && attributeDocsValue.Optional != attributeValue.Optional {
			isConsistent = false
			log.Errorf("'%v' should be marked as Optional in the document", attributeKey)
		}
		if attributeValue.Required == "true" && attributeDocsValue.Required != attributeValue.Required {
			isConsistent = false
			log.Errorf("'%v' should be marked as Required in the document", attributeKey)
		}
		if attributeValue.ForceNew && !attributeDocsValue.ForceNew {
			isConsistent = false
			log.Errorf("'%v' should be marked as ForceNew in the document description", attributeKey)
		}
		if attributeValue.Deprecated != "" && attributeDocsValue.Deprecated == "" {
			isConsistent = false
			log.Errorf("'%v' should be marked as Deprecated in the document description", attributeKey)
		}
	}
	for attributeKey, _ := range resourceAttributeFromDocs {
		if _, ok := resourceAttributes[attributeKey]; !ok {
			isConsistent = false
			log.Errorf("'%v' which described in the docs not found in the resource schema", attributeKey)
		}
	}
	return isConsistent
}

func getResourceAttributes(rootName string, resourceAttributeMap map[string]ResourceAttribute, resourceSchema map[string]*schema.Schema) {
	for key, value := range resourceSchema {
		if rootName != "" {
			key = rootName + "." + key
		}

		if _, ok := resourceAttributeMap[key]; !ok {
			resourceAttributeMap[key] = ResourceAttribute{
				Name:       key,
				Type:       value.Type.String(),
				Optional:   fmt.Sprint(value.Optional),
				Required:   fmt.Sprint(value.Required),
				ForceNew:   value.ForceNew,
				Default:    fmt.Sprint(value.Default),
				Deprecated: value.Deprecated,
			}
		}
		if value.Type == schema.TypeSet || value.Type == schema.TypeList {
			if v, ok := value.Elem.(schema.Schema); ok {
				vv := resourceAttributeMap[key]
				vv.ElemType = v.Type.String()
				resourceAttributeMap[key] = vv
			} else {
				vv := resourceAttributeMap[key]
				vv.ElemType = "Object"
				resourceAttributeMap[key] = vv
				getResourceAttributes(key, resourceAttributeMap, value.Elem.(*schema.Resource).Schema)
			}
		}
	}
}
