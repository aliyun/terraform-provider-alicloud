//nolint:all
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
	skippedSchemaKeys = map[string]struct{}{
		"page_size":   {},
		"page_number": {},
		"total_count": {},
		"max_results": {},
	}
)

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
	Removed     string
}

func main() {
	exitCode := 0
	flag.Parse()
	if fileNames != nil && len(*fileNames) == 0 {
		log.Infof("the diff file is empty, shipped!")
		return
	}

	byt, _ := ioutil.ReadFile(*fileNames)
	diff, _ := diffparser.Parse(string(byt))
	fileRegex := regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*.go")
	fileTestRegex := regexp.MustCompile("alicloud/(resource|data_source)[0-9a-zA-Z_]*_test.go")
	fileDocsRegex := regexp.MustCompile("website/docs/(r|d)/[0-9a-zA-Z_]*.html.markdown")
	resourceNameMap := make(map[string]struct{})
	for _, file := range diff.Files {
		resourceName := ""
		isResource := true
		docsPath := "website/docs/r/"
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

		log.Infof("==> Checking resource or data-source %s attributes consistency...", resourceName)
		resource, ok := alicloud.Provider().(*schema.Provider).ResourcesMap[resourceName]
		if !ok || resource == nil {
			resourceName = strings.TrimPrefix(resourceName, "data_source_")
			resource, ok = alicloud.Provider().(*schema.Provider).DataSourcesMap[resourceName]
			if !ok || resource == nil {
				log.Errorf("resource %s is not found in the provider ResourceMap\n\n", resourceName)
				exitCode = 1
				continue
			}
			docsPath = "website/docs/d/"
			isResource = false
		}
		resourceSchema := resource.Schema
		resourceSchemaFromDocs := make(map[string]ResourceAttribute)
		if err := parseResourceDocs(resourceName, docsPath, isResource, resourceSchemaFromDocs); err != nil {
			log.Errorf("parsing the resource %s docs failed. error: %s", resourceName, err)
			continue
		}

		if consistencyCheck(resourceName, resourceSchemaFromDocs, resourceSchema) {
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

func parseResourceDocs(resourceName, docsPath string, isResource bool, resourceAttributes map[string]ResourceAttribute) error {
	splitRes := strings.Split(resourceName, "alicloud_")
	if len(splitRes) < 2 {
		log.Errorf("parsing resource name %s failed.", resourceName)
		return fmt.Errorf(fmt.Sprintf("parsing resource name %s failed.", resourceName))
	}

	filePath := strings.Join([]string{docsPath, splitRes[1], ".html.markdown"}, "")

	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("open resource %s docs failed. Error: %s", filePath, err)
		return err
	}
	defer file.Close()

	argsRegex := regexp.MustCompile("## Argument Reference")
	attribRegex := regexp.MustCompile("## Attributes Reference")
	secondLevelRegex := regexp.MustCompile("^### `([a-zA-Z_\\.0-9-]*)`")
	argumentsFieldRegex := regexp.MustCompile("^\\* `([a-zA-Z_0-9]*)`[ ]*-? ?(\\(.*\\)) ?(.*)")
	attributeFieldRegex := regexp.MustCompile("^\\s*\\* `([a-zA-Z_0-9]*)`[ ]*-?(.*)")

	name := filepath.Base(filePath)
	re := regexp.MustCompile("[a-z0-9A-Z_]*")
	resourceName = "alicloud_" + re.FindString(name)

	scanner := bufio.NewScanner(file)
	phase := "Argument"
	record := false
	subAttributeName := ""
	line := 0
	rootPrefixLen := 0
	rootName := ""
	extraResourceAttribute := map[string]struct{}{}
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
		if secondLevelRegex.MatchString(strings.TrimSpace(text)) {
			record = true
			parts := strings.Split(strings.TrimSpace(text), " ")
			subAttributeName = strings.Replace(strings.Trim(parts[len(parts)-1], "`"), "-", ".", -1)
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
				if phase == "Attribute" {
					thisLen := len(strings.Split(m[0], "*")[0])

					if rootPrefixLen < thisLen {
						if subAttributeName != "" {
							subAttributeName += "." + rootName
						} else {
							subAttributeName = rootName
						}
						rootName = m[1]
					} else if rootPrefixLen > thisLen {
						parts := strings.Split(subAttributeName, ".")
						backIndex := rootPrefixLen - thisLen
						if backIndex%2 == 0 {
							backIndex /= 2
						} else {
							return fmt.Errorf("resource %s docs %s have not been formatted.", resourceName, docsPath)
						}
						if len(parts) > 0 {
							for backIndex > 0 {
								backIndex--
								if strings.Contains(subAttributeName, ".") {
									subAttributeName = strings.TrimSuffix(strings.TrimSuffix(subAttributeName, parts[len(parts)-1]), ".")
									parts = parts[:len(parts)-1]
								} else {
									subAttributeName = ""
								}
								rootName = m[1]
							}
						}
					} else {
						rootName = m[1]
					}
					rootPrefixLen = thisLen
				}
				attribute := parseMatchLine(m, phase, subAttributeName)
				if attribute == nil {
					continue
				}
				attribute.DocsLineNum = line
				if _, ok := resourceAttributes[attribute.Name]; !ok {
					resourceAttributes[attribute.Name] = *attribute
				} else if isResource {
					extraResourceAttribute[attribute.Name] = struct{}{}
				}
				if phase == "Attribute" {
					for key := range extraResourceAttribute {
						if strings.HasPrefix(attribute.Name, key+".") {
							delete(extraResourceAttribute, key)
						}
					}
				}
			}
		}
	}

	for key := range extraResourceAttribute {
		log.Errorf("'%v' has been set in the `## Argument Reference` and it should be removed from `## Attributes Reference`", key)

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
		if rootName != "" {
			result.Name = rootName + "." + words[1]
		} else {
			result.Name = words[1]
		}
		if strings.Contains(words[2], "Deprecated") {
			result.Deprecated = "Deprecated since"
		}
		//result["Description"] = words[2]
		return &result
	}
	return nil
}

func consistencyCheck(resourceName string, resourceAttributeFromDocs map[string]ResourceAttribute, resourceSchemaDefined map[string]*schema.Schema) bool {
	isConsistent := true
	filteredList := set.NewSet()
	if val, ok := filterList[resourceName]; ok {
		for _, v := range val {
			filteredList.Add(v)
		}
	}

	// the number of the schema field + 1(id) should equal to the number defined in document
	resourceAttributes := make(map[string]ResourceAttribute)
	getResourceAttributes("", resourceAttributes, resourceSchemaDefined, "")

	for attributeKey, attributeValue := range resourceAttributes {
		if _, ok := skippedSchemaKeys[attributeKey]; ok {
			continue
		}
		attributeDocsValue, ok := resourceAttributeFromDocs[attributeKey]
		if attributeValue.Removed != "" {
			continue
		}
		if !ok {
			isConsistent = false
			log.Errorf("'%v' is not found in the docs", attributeKey)
		}
		if attributeValue.Deprecated != "" {
			if attributeDocsValue.Deprecated == "" {
				isConsistent = false
				log.Errorf("'%v' should be marked as Deprecated in the document description", attributeKey)
			}
			continue
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
	}
	for attributeKey, _ := range resourceAttributeFromDocs {
		if _, ok := resourceAttributes[attributeKey]; !ok && attributeKey != "id" {
			isConsistent = false
			log.Errorf("'%v' which described in the docs not found in the resource schema", attributeKey)
		}
	}
	return isConsistent
}

func getResourceAttributes(rootName string, resourceAttributeMap map[string]ResourceAttribute, resourceSchema map[string]*schema.Schema, rootRemoved string) {
	for key, value := range resourceSchema {
		if rootName != "" {
			key = rootName + "." + key
		}

		var thisRemoved = value.Removed
		if len(thisRemoved) == 0 && len(rootRemoved) != 0 {
			thisRemoved = rootRemoved
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
				Removed:    thisRemoved,
			}
		}
		if value.Type == schema.TypeSet || value.Type == schema.TypeList {
			if v, ok := value.Elem.(*schema.Schema); ok {
				vv := resourceAttributeMap[key]
				vv.ElemType = v.Type.String()
				resourceAttributeMap[key] = vv
			} else {
				vv := resourceAttributeMap[key]
				vv.ElemType = "Object"
				resourceAttributeMap[key] = vv
				getResourceAttributes(key, resourceAttributeMap, value.Elem.(*schema.Resource).Schema, thisRemoved)
			}
		}
		if value.Type == schema.TypeMap {
			if _, ok := value.Elem.(*schema.Resource); ok {
				vv := resourceAttributeMap[key]
				vv.ElemType = "Object"
				resourceAttributeMap[key] = vv
				getResourceAttributes(key, resourceAttributeMap, value.Elem.(*schema.Resource).Schema, thisRemoved)

			}
		}
	}
}
