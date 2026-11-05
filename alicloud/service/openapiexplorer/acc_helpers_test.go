package openapiexplorer_test

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/openapiexplorer"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Package-local counterparts of the flat alicloud test helpers
// (service_alicloud_common_test.go), which are not reachable from a service
// package test.

const (
	NOSET      = "#NOSET"       // TestCheckNoResourceAttr
	CHECKSET   = "#CHECKSET"    // TestCheckResourceAttrSet
	REMOVEKEY  = "#REMOVEKEY"   // remove the key from the check/attribute map
	CLEARLIST  = "#CLEARLIST"   // set the attribute to an empty list
	CLEARMAP   = "#CLEARMAP"    // set the attribute to an empty map
	REGEXMATCH = "#REGEXMATCH:" // TestMatchResourceAttr

	indentationSymbol = " "
	childIndent       = 2
)

type describeFunc func(*openapiexplorer.OpenApiExplorerServiceV2, string) (map[string]interface{}, error)

func testAccService() *openapiexplorer.OpenApiExplorerServiceV2 {
	return openapiexplorer.NewOpenApiExplorerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
}

// resourceAttrCheck checks the existence of a resource through its Describe
// method and the attributes listed in checkMap.
type resourceAttrCheck struct {
	resourceId string
	checkMap   map[string]string
	describe   describeFunc
}

func resourceAttrCheckInit(resourceId string, checkMap map[string]string, describe describeFunc) *resourceAttrCheck {
	if checkMap == nil {
		checkMap = make(map[string]string)
	}
	return &resourceAttrCheck{
		resourceId: resourceId,
		checkMap:   checkMap,
		describe:   describe,
	}
}

func (rac *resourceAttrCheck) checkResourceExists(s *terraform.State) error {
	rs, ok := s.RootModule().Resources[rac.resourceId]
	if !ok {
		return fmt.Errorf("can't find resource by id: %s", rac.resourceId)
	}
	if rs.Primary.ID == "" {
		return fmt.Errorf("resource ID is not set")
	}
	if _, err := rac.describe(testAccService(), rs.Primary.ID); err != nil {
		return fmt.Errorf("checking resource %s %s exists error: %w", rac.resourceId, rs.Primary.ID, err)
	}
	return nil
}

// checkResourceDestroy verifies that every resource of the checked type is gone.
func (rac *resourceAttrCheck) checkResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := strings.Split(rac.resourceId, ".")[0]
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			_, err := rac.describe(testAccService(), rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("the resource %s %s was not destroyed", rac.resourceId, rs.Primary.ID)
			}
			if !errs.NotFoundError(err) {
				return err
			}
		}
		return nil
	}
}

// resourceAttrMapUpdateSet returns a function that merges changeMap into the
// accumulated checkMap and checks existence and attributes for a test step.
func (rac *resourceAttrCheck) resourceAttrMapUpdateSet() func(map[string]string) resource.TestCheckFunc {
	return func(changeMap map[string]string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			if err := rac.checkResourceExists(s); err != nil {
				return err
			}
			rac.updateCheckMapPair(changeMap)
			return rac.checkAttributes(s)
		}
	}
}

func (rac *resourceAttrCheck) updateCheckMapPair(changeMap map[string]string) {
	newCheckMap := make(map[string]string, len(rac.checkMap))
	for k, v := range rac.checkMap {
		newCheckMap[k] = v
	}
	for k, v := range changeMap {
		if v == REMOVEKEY {
			delete(newCheckMap, k)
			continue
		}
		newCheckMap[k] = v
	}
	rac.checkMap = newCheckMap
}

func (rac *resourceAttrCheck) checkAttributes(s *terraform.State) error {
	if len(rac.checkMap) == 0 {
		return fmt.Errorf("the parameter \"checkMap\" is nil or empty")
	}
	var errorStrSlice []string
	for key, value := range rac.checkMap {
		var err error
		switch {
		case strings.HasPrefix(value, REGEXMATCH):
			if regex, rErr := regexp.Compile(value[len(REGEXMATCH):]); rErr == nil {
				err = resource.TestMatchResourceAttr(rac.resourceId, key, regex)(s)
			}
		case value == NOSET:
			err = resource.TestCheckNoResourceAttr(rac.resourceId, key)(s)
		case value == CHECKSET:
			err = resource.TestCheckResourceAttrSet(rac.resourceId, key)(s)
		default:
			err = resource.TestCheckResourceAttr(rac.resourceId, key, value)(s)
		}
		if err != nil {
			errorStrSlice = append(errorStrSlice, err.Error())
		}
	}
	if len(errorStrSlice) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errorStrSlice, "\n"))
}

// resourceTestAccConfigFunc returns a function that renders the dependence
// config plus the primary resource built from the accumulated attribute map.
func resourceTestAccConfigFunc(resourceId string, name string, configDependence func(name string) string) func(map[string]interface{}) string {
	attributeMap := make(map[string]interface{})
	return func(changeMap map[string]interface{}) string {
		newMap := make(map[string]interface{}, len(attributeMap))
		for k, v := range attributeMap {
			newMap[k] = v
		}
		for k, v := range changeMap {
			if sentinel, ok := v.(string); ok {
				switch sentinel {
				case REMOVEKEY:
					delete(newMap, k)
					continue
				case CLEARMAP:
					newMap[k] = make(map[string]interface{})
					continue
				case CLEARLIST:
					newMap[k] = make([]interface{}, 0)
					continue
				}
			}
			newMap[k] = v
		}
		attributeMap = newMap

		strs := strings.Split(resourceId, ".")
		primaryConfig := fmt.Sprintf("\n\nresource \"%s\" \"%s\" ", strs[0], strs[1])
		return configDependence(name) + primaryConfig + fmt.Sprint(valueConvert(0, reflect.ValueOf(attributeMap)))
	}
}

func valueConvert(indentation int, val reflect.Value) interface{} {
	switch val.Kind() {
	case reflect.Interface:
		return valueConvert(indentation, reflect.ValueOf(val.Interface()))
	case reflect.String:
		return fmt.Sprintf("\"%s\"", val.String())
	case reflect.Slice:
		return listValue(indentation, val)
	case reflect.Map:
		return mapValue(indentation, val)
	case reflect.Bool:
		return val.Bool()
	case reflect.Int, reflect.Int64:
		return val.Int()
	default:
		log.Panicf("invalid attribute value type: %#v", val)
	}
	return ""
}

func listValue(indentation int, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation+childIndent)+
			fmt.Sprint(valueConvert(indentation+childIndent, val.Index(i))))
	}
	return fmt.Sprintf("[\n%s\n%s]", strings.Join(valList, ",\n"), addIndentation(indentation))
}

func mapValue(indentation int, val reflect.Value) string {
	var valList []string
	for _, keyV := range val.MapKeys() {
		mapVal := getRealValueType(val.MapIndex(keyV))
		if mapVal.Kind() == reflect.Slice && mapVal.Len() > 0 {
			if getRealValueType(mapVal.Index(0)).Kind() == reflect.Map {
				valList = append(valList, addIndentation(indentation)+
					listValueMapChild(indentation+childIndent, keyV.String(), mapVal))
				continue
			}
		}
		value := valueConvert(indentation+len(keyV.String())+childIndent+3, val.MapIndex(keyV))
		var line string
		switch value.(type) {
		case bool:
			line = fmt.Sprintf(`%s%s = %t`, addIndentation(indentation+childIndent), keyV.String(), value)
		case int64:
			line = fmt.Sprintf(`%s%s = %d`, addIndentation(indentation+childIndent), keyV.String(), value)
		default:
			line = fmt.Sprintf(`%s%s = %s`, addIndentation(indentation+childIndent), keyV.String(), value)
		}
		valList = append(valList, line)
	}
	return fmt.Sprintf("{\n%s\n%s}", strings.Join(valList, "\n"), addIndentation(indentation))
}

func listValueMapChild(indentation int, key string, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation)+key+" "+
			mapValue(indentation, getRealValueType(val.Index(i))))
	}
	return fmt.Sprintf("%s\n%s", strings.Join(valList, "\n"), addIndentation(indentation))
}

func getRealValueType(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Interface {
		return getRealValueType(reflect.ValueOf(value.Interface()))
	}
	return value
}

func addIndentation(indentation int) string {
	return strings.Repeat(indentationSymbol, indentation)
}
