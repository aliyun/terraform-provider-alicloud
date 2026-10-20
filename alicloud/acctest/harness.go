package acctest

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
			rac.ResourceAttr.updateCheckMapPair(changeMap)
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
