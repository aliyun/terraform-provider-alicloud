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
	resourceName = flag.String("resource", "", "the name of the terraform resource to diff")
	fileName     = flag.String("file_name", "", "the file to check diff")
	filterList   = map[string][]string{
		"alicloud_amqp_instance":            []string{"logistics"},
		"alicloud_cms_alarm":                []string{"notify_type"},
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

func TestConsistencyWithDocument(t *testing.T) {
	flag.Parse()
	if resourceName != nil && len(*resourceName) == 0 {
		log.Warningf("the resource name is empty")
		return
	}
	obj := alicloud.Provider().(*schema.Provider).ResourcesMap[*resourceName].Schema
	objSchema := make(map[string]interface{}, 0)
	objMd, err := parseResource(*resourceName)
	if err != nil {
		log.Error(err)
		t.Fatal()
	}
	mergeMaps(objSchema, objMd.Arguments, objMd.Attributes)

	if consistencyCheck(t, *resourceName, objSchema, obj) {
		t.Fatal("the consistency with document has occurred")
		os.Exit(1)
	}
}

func TestFieldCompatibilityCheck(t *testing.T) {
	flag.Parse()
	if fileName != nil && len(*fileName) == 0 {
		log.Warningf("the diff file is empty")
		return
	}
	byt, _ := ioutil.ReadFile(*fileName)
	diff, _ := diffparser.Parse(string(byt))
	res := false
	fileRegex := regexp.MustCompile("alicloud/resource[0-9a-zA-Z_]*.go")
	fileTestRegex := regexp.MustCompile("alicloud/resource[0-9a-zA-Z_]*_test.go")
	for _, file := range diff.Files {
		if fileRegex.MatchString(file.NewName) {
			if fileTestRegex.MatchString(file.NewName) {
				continue
			}
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

func parseResource(resourceName string) (*Resource, error) {
	splitRes := strings.Split(resourceName, "alicloud_")
	if len(splitRes) < 2 {
		log.Errorf("the resource name parsed failed")
		return nil, fmt.Errorf("the resource name parsed failed")
	}
	basePath := "../website/docs/r/"
	filePath := strings.Join([]string{basePath, splitRes[1], ".html.markdown"}, "")

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("cannot open text file: %s, err: [%v]", filePath, err)
		return nil, err
	}
	defer file.Close()

	argsRegex := regexp.MustCompile("## Argument Reference")
	attribRegex := regexp.MustCompile("## Attributes Reference")
	secondLevelRegex := regexp.MustCompile("^#+")
	argumentsFieldRegex := regexp.MustCompile("^\\* `([a-zA-Z_0-9]*)`[ ]*-? ?(\\(.*\\)) ?(.*)")
	attributeFieldRegex := regexp.MustCompile("^\\* `([a-zA-Z_0-9]*)`[ ]*-?(.*)")

	name := filepath.Base(filePath)
	re := regexp.MustCompile("[a-z0-9A-Z_]*")
	resourceName = "alicloud_" + re.FindString(name)
	result := &Resource{Name: resourceName, Arguments: map[string]interface{}{}, Attributes: map[string]interface{}{}}
	log.Infof("the resourceName = %s\n", resourceName)

	scanner := bufio.NewScanner(file)
	phase := "Argument"
	record := false
	for scanner.Scan() {
		line := scanner.Text()
		if argsRegex.MatchString(line) {
			record = true
			phase = "Argument"
			continue
		}
		if attribRegex.MatchString(line) {
			record = true
			phase = "Attribute"
			continue
		}
		if secondLevelRegex.MatchString(line) && strings.HasSuffix(line, "params") {
			record = true
			continue
		}
		if record {
			if secondLevelRegex.MatchString(line) && !strings.HasSuffix(line, "params") {
				record = false
				continue
			}
			var matched [][]string
			if phase == "Argument" {
				matched = argumentsFieldRegex.FindAllStringSubmatch(line, 1)
			} else if phase == "Attribute" {
				matched = attributeFieldRegex.FindAllStringSubmatch(line, 1)
			}

			for _, m := range matched {
				Field := parseMatchLine(m, phase)
				Field["Type"] = phase
				if v, exist := Field["Name"]; exist {
					result.Arguments[v.(string)] = Field
				}
			}
		}
	}
	return result, nil
}

func parseMatchLine(words []string, phase string) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	if phase == "Argument" && len(words) >= 4 {
		result["Name"] = words[1]
		result["Description"] = words[3]
		if strings.Contains(words[2], "Optional") {
			result["Optional"] = true
		}
		if strings.Contains(words[2], "Required") {
			result["Required"] = true
		}
		if strings.Contains(words[2], "ForceNew") {
			result["ForceNew"] = true
		}
		return result
	}
	if phase == "Attribute" && len(words) >= 3 {
		result["Name"] = words[1]
		result["Description"] = words[2]
		return result
	}
	return nil
}

func consistencyCheck(t *testing.T, resourceName string, doc map[string]interface{}, resource map[string]*schema.Schema) bool {
	res := false
	fileteredList := set.NewSet()
	if val, ok := filterList[resourceName]; ok {
		for _, v := range val {
			fileteredList.Add(v)
		}
	}
	defer func() {
		if r := recover(); r != nil {
			res = true
			log.Errorf("internal error: Please email terraform@alibabacloud.com to report the issue with the related resource")
			t.Fatal()
		}
	}()

	// the number of the schema field + 1(id) should equal to the number defined in document
	if len(resource)+1 != len(doc) {
		record := set.NewSet()
		for field, _ := range doc {
			if field == "id" || fileteredList.Contains(field) {
				delete(doc, field)
				continue
			}
			if _, exist := resource[field]; exist {
				delete(doc, field)
				delete(resource, field)
			} else if !exist {
				// the field existed in Document,but not existed in resource
				record.Add(field)
			}
		}
		if len(resource) != 0 {
			for field, _ := range resource {
				if fileteredList.Contains(field) {
					record.Remove(field)
					continue
				}
				// the field existed in resource,but not existed in document
				record.Add(field)
			}
		}
		if record.Cardinality() != 0 {
			log.Errorf("there is missing attribute %v description in the document", record)
			return true
		}
	}
	for field, docFieldObj := range doc {
		docObj := docFieldObj.(map[string]interface{})
		if docObj["Type"] == "Attribute" || fileteredList.Contains(field) {
			continue
		}
		resourceFieldObj := resource[field]
		if resourceFieldObj == nil {
			res = true
			panic(field)
		}
		if _, exist1 := docObj["Optional"]; exist1 && !resourceFieldObj.Optional {
			res = true
			log.Errorf("attribute %v should be marked as Optional in the in the document description", field)
		}
		if _, exist1 := docObj["Required"]; exist1 && !resourceFieldObj.Required {
			res = true
			log.Errorf("attribute %v should be marked as Required in the in the document description", field)
		}
		if _, exist1 := docObj["ForceNew"]; exist1 && !resourceFieldObj.ForceNew {
			res = true
			log.Errorf("attribute %v should be marked as ForceNew in the document description", field)
		}
	}
	return res
}

func mergeMaps(Dst map[string]interface{}, arr ...map[string]interface{}) map[string]interface{} {
	for _, m := range arr {
		for k, v := range m {
			if _, exist := Dst[k]; exist {
				continue
			}
			Dst[k] = v
		}
	}
	return Dst
}
