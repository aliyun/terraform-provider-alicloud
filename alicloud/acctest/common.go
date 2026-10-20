package acctest

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// This file provides the flat-suite ACC test harness machinery (resourceCheck /
// resourceAttr / resourceConfig / dataSourceAttr reflection engine and its
// sentinels), sunk from alicloud/service_alicloud_common_test.go so split
// packages consume one framework instead of copying it. The runner is
// terraform-plugin-testing: one test binary links at most one helper/resource,
// and SDKv2 helper/resource already occupies the root binary.

// be used to check attribute map value
const (
	NOSET      = "#NOSET"       // be equivalent to method "TestCheckNoResourceAttrSet"
	CHECKSET   = "#CHECKSET"    // "TestCheckResourceAttrSet"
	REMOVEKEY  = "#REMOVEKEY"   // remove checkMap key
	CLEARLIST  = "#CLEARLIST"   // clear checkMap key and set empty list
	CLEARMAP   = "#CLEARMAP"    // clear checkMap key and set empty map
	REGEXMATCH = "#REGEXMATCH:" // "TestMatchResourceAttr" ,the map name/key like `"attribute" : REGEXMATCH + "attributeString"`
	ForceSleep = "force_sleep"
)

const (
	// indentation symbol
	INDENTATIONSYMBOL = " "

	// child field indend number
	CHILDINDEND = 2
)

// get a function that change checkMap pairs for a series test step
type ResourceAttrMapUpdate func(map[string]string) resource.TestCheckFunc

// get a function that change attributeMap pairs for a series test step
type ResourceTestAccConfigFunc func(map[string]interface{}) string

// check the existence of resource
type ResourceCheck struct {
	// IDRefreshName, like "alicloud_instance.foo"
	ResourceId string

	// The response of the service method DescribeXXX
	resourceObject interface{}

	// The resource service client type, like DnsService, VpcService
	serviceFunc func() interface{}

	// service describe method name
	DescribeMethod string

	// additional attributes
	additionalAttrs []string

	// additional attributes type
	additionalAttrsType map[string]schema.ValueType
}

func ResourceCheckInit(resourceId string, resourceObject interface{}, serviceFunc func() interface{}, additionalAttrs ...string) *ResourceCheck {
	rc := &ResourceCheck{
		ResourceId:      resourceId,
		resourceObject:  resourceObject,
		serviceFunc:     serviceFunc,
		additionalAttrs: additionalAttrs,
	}
	if len(rc.additionalAttrs) > 0 {
		rc.setAdditionalAttrsType()
	}
	return rc
}

func ResourceCheckInitWithDescribeMethod(resourceId string, resourceObject interface{}, serviceFunc func() interface{}, describeMethod string, additionalAttrs ...string) *ResourceCheck {
	rc := &ResourceCheck{
		ResourceId:      resourceId,
		resourceObject:  resourceObject,
		serviceFunc:     serviceFunc,
		DescribeMethod:  describeMethod,
		additionalAttrs: additionalAttrs,
	}
	if len(rc.additionalAttrs) > 0 {
		rc.setAdditionalAttrsType()
	}
	return rc
}

// caching the additional attribute type used to convert the addition attribute value type before calling Get method
func (rc *ResourceCheck) setAdditionalAttrsType() {
	provider := Provider
	resourceType, ok := provider.ResourcesMap[strings.Split(rc.ResourceId, ".")[0]]
	if !ok {
		log.Panicf("invalid resource type: %s", strings.Split(rc.ResourceId, ".")[0])
	}
	if rc.additionalAttrsType == nil {
		rc.additionalAttrsType = make(map[string]schema.ValueType)
	}
	for _, attr := range rc.additionalAttrs {
		if s, ok := resourceType.Schema[attr]; !ok {
			log.Panicf("invalid resource attribute: %s", attr)
		} else {
			rc.additionalAttrsType[attr] = s.Type
		}
	}
	return
}

// check attribute only
type ResourceAttr struct {
	resourceId string
	checkMap   map[string]string
}

func ResourceAttrInit(resourceId string, checkMap map[string]string) *ResourceAttr {
	if checkMap == nil {
		checkMap = make(map[string]string)
	}
	return &ResourceAttr{
		resourceId: resourceId,
		checkMap:   checkMap,
	}
}

// check the existence and attribute of the resource at the same time
type ResourceAttrCheck struct {
	*ResourceCheck
	*ResourceAttr
}

func ResourceAttrCheckInit(rc *ResourceCheck, ra *ResourceAttr) *ResourceAttrCheck {
	return &ResourceAttrCheck{
		ResourceCheck: rc,
		ResourceAttr:  ra,
	}
}

// check the resource existence by invoking DescribeXXX method of service and assign *ResourceCheck.resourceObject value,
// the service is returned by invoking *ResourceCheck.serviceFunc
func (rc *ResourceCheck) CheckResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		rs, ok := s.RootModule().Resources[rc.ResourceId]
		if !ok {
			return errs.WrapError(fmt.Errorf("can't find resource by id: %s", rc.ResourceId))

		}
		outValue, err := rc.callDescribeMethod(rs)
		if err != nil {
			return errs.WrapError(err)
		}
		errorValue := outValue[1]
		if !errorValue.IsNil() {
			return errs.WrapError(fmt.Errorf("Checking resource %s %s exists error:%s ", rc.ResourceId, rs.Primary.ID, errorValue.Interface().(error).Error()))
		}
		/*if reflect.TypeOf(rc.resourceObject).Elem().String() == outValue[0].Type().String() {
			reflect.ValueOf(rc.resourceObject).Elem().Set(outValue[0])
			return nil
		} else {
			return errs.WrapError(fmt.Errorf("The response object type expected *%s, got %s ", outValue[0].Type().String(), reflect.TypeOf(rc.resourceObject).String()))
		}*/
		return nil
	}
}

// check the resource destroy
func (rc *ResourceCheck) CheckResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		strs := strings.Split(rc.ResourceId, ".")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "alicloud_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}

		if resourceType == "" {
			return errs.WrapError(errs.Error("The resourceId %s is not correct and it should prefix with alicloud_", rc.ResourceId))
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			outValue, err := rc.callDescribeMethod(rs)
			errorValue := outValue[1]
			if !errorValue.IsNil() {
				err = errorValue.Interface().(error)
				if err != nil {
					if errs.NotFoundError(err) {
						continue
					}
					return errs.WrapError(err)
				}
			} else {
				return errs.WrapError(errs.Error("the resource %s %s was not destroyed ! ", rc.ResourceId, rs.Primary.ID))
			}
		}
		return nil
	}
}

// invoking DescribeXXX method of service
func (rc *ResourceCheck) callDescribeMethod(rs *terraform.ResourceState) ([]reflect.Value, error) {
	var err error
	if rs.Primary.ID == "" {
		return nil, errs.WrapError(fmt.Errorf("resource ID is not set"))
	}
	serviceP := rc.serviceFunc()
	if rc.DescribeMethod == "" {
		rc.DescribeMethod, err = getResourceDescribeMethod(rc.ResourceId)
		if err != nil {
			return nil, errs.WrapError(err)
		}
	}
	value := reflect.ValueOf(serviceP)
	typeName := value.Type().String()
	value = value.MethodByName(rc.DescribeMethod)
	if !value.IsValid() {
		return nil, errs.WrapError(errs.Error("The service type %s does not have method %s", typeName, rc.DescribeMethod))
	}
	inValue := []reflect.Value{reflect.ValueOf(rs.Primary.ID)}
	for _, attr := range rc.additionalAttrs {
		if attrValue, ok := rs.Primary.Attributes[attr]; ok {
			if attrType, o := rc.additionalAttrsType[attr]; o {
				switch attrType {
				case schema.TypeBool:
					v, _ := strconv.ParseBool(attrValue)
					inValue = append(inValue, reflect.ValueOf(v))
					continue
				case schema.TypeInt:
					v, _ := strconv.ParseInt(attrValue, 10, 64)
					inValue = append(inValue, reflect.ValueOf(v))
					continue
				}
			}
			inValue = append(inValue, reflect.ValueOf(attrValue))
		}
	}
	return value.Call(inValue), nil
}

func getResourceDescribeMethod(resourceId string) (string, error) {
	start := strings.Index(resourceId, "alicloud_")
	if start < 0 {
		return "", errs.WrapError(fmt.Errorf("the parameter \"name\" don't contain string \"alicloud_\""))
	}
	start += len("alicloud_")
	end := strings.Index(resourceId[start:], ".") + start
	if end < 0 {
		return "", errs.WrapError(fmt.Errorf("the parameter \"name\" don't contain string \".\""))
	}
	strs := strings.Split(resourceId[start:end], "_")
	describeName := "Describe"
	for _, str := range strs {
		describeName = describeName + strings.Title(str)
	}
	return describeName, nil
}

// check attribute func and check resource exist
func (rac *ResourceAttrCheck) resourceAttrMapCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := rac.ResourceCheck.CheckResourceExists()(s)
		if err != nil {
			return errs.WrapError(err)
		}
		return rac.ResourceAttr.resourceAttrMapCheck()(s)
	}
}

// execute the callback before check attribute and check resource exist
func (rac *ResourceAttrCheck) resourceAttrMapCheckWithCallback(callback func()) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := rac.ResourceCheck.CheckResourceExists()(s)
		if err != nil {
			return errs.WrapError(err)
		}
		return rac.ResourceAttr.resourceAttrMapCheckWithCallback(callback)(s)
	}
}

// get ResourceAttrMapUpdate for a series test step and check resource exist
func (rac *ResourceAttrCheck) ResourceAttrMapUpdateSet() ResourceAttrMapUpdate {
	return func(changeMap map[string]string) resource.TestCheckFunc {
		callback := func() {
			rac.updateCheckMapPair(changeMap)
		}
		return rac.resourceAttrMapCheckWithCallback(callback)
	}
}

// make a new map and copy from the old field checkMap, then update it according to the changeMap
func (ra *ResourceAttr) updateCheckMapPair(changeMap map[string]string) {
	if interval, ok := changeMap[ForceSleep]; ok {
		intervalInt, err := strconv.Atoi(interval)
		if err == nil {
			time.Sleep(time.Duration(intervalInt) * time.Second)
			delete(changeMap, ForceSleep)
		}
	}
	newCheckMap := make(map[string]string, len(ra.checkMap))
	for k, v := range ra.checkMap {
		newCheckMap[k] = v
	}
	ra.checkMap = newCheckMap
	if changeMap != nil && len(changeMap) > 0 {
		for rk, rv := range changeMap {
			_, ok := ra.checkMap[rk]
			if rv == REMOVEKEY && ok {
				delete(ra.checkMap, rk)
			} else if ok {
				delete(ra.checkMap, rk)
				ra.checkMap[rk] = rv
			} else {
				ra.checkMap[rk] = rv
			}
		}
	}
}

// check attribute func
func (ra *ResourceAttr) resourceAttrMapCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[ra.resourceId]
		if !ok {
			return errs.WrapError(fmt.Errorf("can't find resource by id: %s", ra.resourceId))
		}
		if rs.Primary.ID == "" {
			return errs.WrapError(fmt.Errorf("resource ID is not set"))
		}

		if ra.checkMap == nil || len(ra.checkMap) == 0 {
			return errs.WrapError(fmt.Errorf("the parameter \"checkMap\" is nil or empty"))
		}

		var errorStrSlice []string
		errorStrSlice = append(errorStrSlice, "")
		for key, value := range ra.checkMap {
			var err error
			if strings.HasPrefix(value, REGEXMATCH) {
				var regex *regexp.Regexp
				regex, err = regexp.Compile(value[len(REGEXMATCH):])
				if err == nil {
					err = resource.TestMatchResourceAttr(ra.resourceId, key, regex)(s)
				} else {
					err = nil
				}
			} else if value == NOSET {
				err = resource.TestCheckNoResourceAttr(ra.resourceId, key)(s)
			} else if value == CHECKSET {
				err = resource.TestCheckResourceAttrSet(ra.resourceId, key)(s)
			} else {
				err = resource.TestCheckResourceAttr(ra.resourceId, key, value)(s)
			}
			if err != nil {
				errorStrSlice = append(errorStrSlice, err.Error())
			}
		}
		if len(errorStrSlice) == 1 {
			return nil
		}
		return errs.WrapError(fmt.Errorf("%s", strings.Join(errorStrSlice, "\n")))
	}
}

// execute the callback before check attribute
func (ra *ResourceAttr) resourceAttrMapCheckWithCallback(callback func()) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		callback()
		return ra.resourceAttrMapCheck()(s)
	}
}

// get ResourceAttrMapUpdate for a series test step
func (ra *ResourceAttr) ResourceAttrMapUpdateSet() ResourceAttrMapUpdate {
	return func(changeMap map[string]string) resource.TestCheckFunc {
		callback := func() {
			ra.updateCheckMapPair(changeMap)
		}
		return ra.resourceAttrMapCheckWithCallback(callback)
	}
}

// ResourceTestAccConfig is the exported form of the flat suite's
// resourceTestAccConfigFunc: the direct ResourceTestAccConfigFunc spelling is
// already taken by the type this constructor returns.
func ResourceTestAccConfig(resourceId string,
	name string,
	configDependence func(name string) string) ResourceTestAccConfigFunc {
	basicInfo := resourceConfig{
		name:             name,
		resourceId:       resourceId,
		attributeMap:     make(map[string]interface{}),
		configDependence: configDependence,
		number:           0,
	}
	return basicInfo.configBuild(false)
}

func DataSourceTestAccConfigFunc(resourceId string,
	name string,
	configDependence func(name string) string) ResourceTestAccConfigFunc {
	basicInfo := resourceConfig{
		name:             name,
		resourceId:       resourceId,
		attributeMap:     make(map[string]interface{}),
		configDependence: configDependence,
	}
	return basicInfo.configBuild(true)
}

// be used for generate testcase step config
type resourceConfig struct {
	// the resource name
	name string

	resourceId string

	// store attribute value that primary resource
	attributeMap map[string]interface{}

	// generate assistant test config
	configDependence func(name string) string

	// step number
	number int
}

// according to changeMap to change the attributeMap value
func (b *resourceConfig) configUpdate(changeMap map[string]interface{}) {
	newMap := make(map[string]interface{}, len(b.attributeMap))
	for k, v := range b.attributeMap {
		newMap[k] = v
	}
	b.attributeMap = newMap
	if changeMap != nil && len(changeMap) > 0 {
		for rk, rv := range changeMap {
			_, ok := b.attributeMap[rk]
			// Sentinels (REMOVEKEY/CLEARMAP/CLEARLIST) must be handled regardless of
			// whether the key already exists in attributeMap. Previously the `ok &&`
			// guard skipped the sentinel branch for absent keys, leaking the literal
			// sentinel string into TypeInt fields and producing invalid configs.
			if strValue, isCost := rv.(string); isCost && (strValue == REMOVEKEY || strValue == CLEARMAP || strValue == CLEARLIST) {
				delete(b.attributeMap, rk)
				if strValue == CLEARMAP {
					b.attributeMap[rk] = make(map[string]interface{})
				} else if strValue == CLEARLIST {
					b.attributeMap[rk] = make([]interface{}, 0)
				}
			} else if ok {
				delete(b.attributeMap, rk)
				b.attributeMap[rk] = rv
			} else {
				b.attributeMap[rk] = rv
			}
		}
	}
}

// get BasicConfigFunc for resource a series test step
// overwrite: if true ,the attributeMap will be replace by changMap , other will be update
func (b *resourceConfig) configBuild(overwrite bool) ResourceTestAccConfigFunc {
	return func(changeMap map[string]interface{}) string {
		if overwrite {
			b.attributeMap = changeMap
		} else {
			b.configUpdate(changeMap)
		}
		strs := strings.Split(b.resourceId, ".")
		assistantConfig := b.configDependence(b.name)
		var primaryConfig string
		if strings.Compare("data", strs[0]) == 0 {
			primaryConfig = fmt.Sprintf("\n\ndata \"%s\" \"%s\" ", strs[1], strs[2])
		} else {
			primaryConfig = fmt.Sprintf("\n\nresource \"%s\" \"%s\" ", strs[0], strs[1])
		}
		config := assistantConfig + primaryConfig + fmt.Sprint(valueConvert(0, reflect.ValueOf(b.attributeMap)))
		for _, part := range strings.Split(os.Getenv("DEBUG"), ",") {
			if strings.TrimSpace(part) == "terraform_test" {
				log.Printf("###### (step %d) terraform test configuration ###### %s \n###### (END) ######\n\n", b.number, config)
				b.number += 1
			}
		}
		return config
	}
}

// deal with the parameter common method
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
	case reflect.Int:
		return val.Int()
	case reflect.Int64:
		return val.Int()
	default:
		log.Panicf("invalid attribute value type: %#v", val)
	}
	return ""
}

// deal with list parameter
func listValue(indentation int, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation+CHILDINDEND)+
			fmt.Sprint(valueConvert(indentation+CHILDINDEND, val.Index(i))))
	}

	return fmt.Sprintf("[\n%s\n%s]", strings.Join(valList, ",\n"), addIndentation(indentation))
}

// deal with map parameter
func mapValue(indentation int, val reflect.Value) string {
	var valList []string
	for _, keyV := range val.MapKeys() {
		mapVal := getRealValueType(val.MapIndex(keyV))
		var line string
		if mapVal.Kind() == reflect.Slice && mapVal.Len() > 0 {
			eleVal := getRealValueType(mapVal.Index(0))
			if eleVal.Kind() == reflect.Map {
				line = fmt.Sprintf(`%s%s`, addIndentation(indentation),
					listValueMapChild(indentation+CHILDINDEND, keyV.String(), mapVal))
				valList = append(valList, line)
				continue
			}
		}
		value := valueConvert(indentation+len(keyV.String())+CHILDINDEND+3, val.MapIndex(keyV))
		switch value.(type) {
		case bool:
			line = fmt.Sprintf(`%s%s = %t`, addIndentation(indentation+CHILDINDEND), keyV.String(), value)
		case int:
			line = fmt.Sprintf(`%s%s = %d`, addIndentation(indentation+CHILDINDEND), keyV.String(), value)
		case int64:
			line = fmt.Sprintf(`%s%s = %d`, addIndentation(indentation+CHILDINDEND), keyV.String(), value)
		default:
			line = fmt.Sprintf(`%s%s = %s`, addIndentation(indentation+CHILDINDEND), keyV.String(), value)
		}

		valList = append(valList, line)
	}
	return fmt.Sprintf("{\n%s\n%s}", strings.Join(valList, "\n"), addIndentation(indentation))
}

// deal with list parameter that child element is map
func listValueMapChild(indentation int, key string, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation)+key+" "+
			mapValue(indentation, getRealValueType(val.Index(i))))
	}

	return fmt.Sprintf("%s\n%s", strings.Join(valList, "\n"), addIndentation(indentation))
}

func getRealValueType(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Interface:
		return getRealValueType(reflect.ValueOf(value.Interface()))
	default:
		return value
	}
}

func addIndentation(indentation int) string {
	return strings.Repeat(INDENTATIONSYMBOL, indentation)
}

// in most cases, the TestCheckFunc list of dataSource test case is repeated，so we make an abstract in
// order to reduce redundant code.
// DataSourceAttr has 3 field ,incloud ResourceId  ExistMapFunc FakeMapFunc, every dataSource test can use only one
type DataSourceAttr struct {
	// IDRefreshName, like "data.alicloud_dns_records.record"
	ResourceId string

	// get existMap function
	ExistMapFunc func(rand int) map[string]string

	// get fakeMap function
	FakeMapFunc func(rand int) map[string]string
}

// get exist and empty ResourceAttrMapUpdate function
func (dsa *DataSourceAttr) checkDataSourceAttr(rand int) (exist, empty ResourceAttrMapUpdate) {
	exist = ResourceAttrInit(dsa.ResourceId, dsa.ExistMapFunc(rand)).ResourceAttrMapUpdateSet()
	empty = ResourceAttrInit(dsa.ResourceId, dsa.FakeMapFunc(rand)).ResourceAttrMapUpdateSet()
	return
}

// according to configs generate step list and execute the test
// lintignore: AT001
func (dsa *DataSourceAttr) DataSourceTestCheck(t *testing.T, rand int, configs ...DataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			PreCheck(t)
		},
		ProviderFactories: SDKProviderFactories,
		Steps:             steps,
	})
}

// according to configs generate step list and execute the test with preCheck
// lintignore: AT001
func (dsa *DataSourceAttr) DataSourceTestCheckWithPreCheck(t *testing.T, rand int, preCheck func(), configs ...DataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          preCheck,
		ProviderFactories: SDKProviderFactories,
		Steps:             steps,
	})
}

// per schema attribute test config
type DataSourceTestAccConfig struct {
	// be equal to testCase config string,but the result has only one record
	ExistConfig string

	// if the DataSourceAttr.ExistMapFunc returned map value not match we want, ExistChangMap can alter checkMap for ExistConfig
	ExistChangMap map[string]string

	// be equal to testCase config string,but the result is empty
	FakeConfig string

	// if the DataSourceAttr.FakeMapFunc returned map value not match we want, FakeChangMap can alter checkMap for FakeConfig
	FakeChangMap map[string]string
}

// build test cases for each attribute
func (conf *DataSourceTestAccConfig) buildDataSourceSteps(t *testing.T, info *DataSourceAttr, rand int) []resource.TestStep {
	testAccCheckExist, testAccCheckEmpty := info.checkDataSourceAttr(rand)
	var steps []resource.TestStep
	if conf.ExistConfig != "" {
		step := resource.TestStep{
			Config: conf.ExistConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckExist(conf.ExistChangMap),
			),
		}
		steps = append(steps, step)
	}
	if conf.FakeConfig != "" {
		step := resource.TestStep{
			Config: conf.FakeConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckEmpty(conf.FakeChangMap),
			),
		}
		steps = append(steps, step)
	}
	return steps
}
