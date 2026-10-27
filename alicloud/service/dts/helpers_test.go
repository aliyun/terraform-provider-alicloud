package dts_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ---------------------------------------------------------------------------
// Constants and types replicated from the root package test infrastructure
// (service_alicloud_common_test.go). The leaf test package cannot import the
// root package's unexported test helpers, so they are duplicated here with
// errs.* replacing the root package's WrapError/Error/NotFoundError aliases.
// ---------------------------------------------------------------------------

const (
	NOSET      = "#NOSET"
	CHECKSET   = "#CHECKSET"
	REMOVEKEY  = "#REMOVEKEY"
	CLEARLIST  = "#CLEARLIST"
	CLEARMAP   = "#CLEARMAP"
	REGEXMATCH = "#REGEXMATCH:"
	ForceSleep = "force_sleep"
)

const (
	INDENTATIONSYMBOL = " "
	CHILDINDEND       = 2
)

type resourceAttrMapUpdate func(map[string]string) resource.TestCheckFunc
type ResourceTestAccConfigFunc func(map[string]interface{}) string

type resourceCheck struct {
	resourceId          string
	resourceObject      interface{}
	serviceFunc         func() interface{}
	describeMethod      string
	additionalAttrs     []string
	additionalAttrsType map[string]schema.ValueType
}

type resourceAttr struct {
	resourceId string
	checkMap   map[string]string
}

type resourceAttrCheck struct {
	*resourceCheck
	*resourceAttr
}

type resourceConfig struct {
	name             string
	resourceId       string
	attributeMap     map[string]interface{}
	configDependence func(name string) string
	number           int
}

func resourceCheckInit(resourceId string, resourceObject interface{}, serviceFunc func() interface{}, additionalAttrs ...string) *resourceCheck {
	rc := &resourceCheck{
		resourceId:      resourceId,
		resourceObject:  resourceObject,
		serviceFunc:     serviceFunc,
		additionalAttrs: additionalAttrs,
	}
	if len(rc.additionalAttrs) > 0 {
		rc.setAdditionalAttrsType()
	}
	return rc
}

func resourceCheckInitWithDescribeMethod(resourceId string, resourceObject interface{}, serviceFunc func() interface{}, describeMethod string, additionalAttrs ...string) *resourceCheck {
	rc := &resourceCheck{
		resourceId:      resourceId,
		resourceObject:  resourceObject,
		serviceFunc:     serviceFunc,
		describeMethod:  describeMethod,
		additionalAttrs: additionalAttrs,
	}
	if len(rc.additionalAttrs) > 0 {
		rc.setAdditionalAttrsType()
	}
	return rc
}

func (rc *resourceCheck) setAdditionalAttrsType() {
	provider := acctest.Provider
	resourceType, ok := provider.ResourcesMap[strings.Split(rc.resourceId, ".")[0]]
	if !ok {
		log.Panicf("invalid resource type: %s", strings.Split(rc.resourceId, ".")[0])
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

func resourceAttrInit(resourceId string, checkMap map[string]string) *resourceAttr {
	if checkMap == nil {
		checkMap = make(map[string]string)
	}
	return &resourceAttr{
		resourceId: resourceId,
		checkMap:   checkMap,
	}
}

func resourceAttrCheckInit(rc *resourceCheck, ra *resourceAttr) *resourceAttrCheck {
	return &resourceAttrCheck{
		resourceCheck: rc,
		resourceAttr:  ra,
	}
}

func (rc *resourceCheck) checkResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		rs, ok := s.RootModule().Resources[rc.resourceId]
		if !ok {
			return errs.WrapError(fmt.Errorf("can't find resource by id: %s", rc.resourceId))
		}
		outValue, err := rc.callDescribeMethod(rs)
		if err != nil {
			return errs.WrapError(err)
		}
		errorValue := outValue[1]
		if !errorValue.IsNil() {
			return errs.WrapError(fmt.Errorf("Checking resource %s %s exists error:%s ", rc.resourceId, rs.Primary.ID, errorValue.Interface().(error).Error()))
		}
		return nil
	}
}

func (rc *resourceCheck) checkResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceId, ".")
		var resourceType string
		for _, str := range strs {
			if strings.Contains(str, "alicloud_") {
				resourceType = strings.Trim(str, " ")
				break
			}
		}

		if resourceType == "" {
			return errs.WrapError(errs.Error("The resourceId %s is not correct and it should prefix with alicloud_", rc.resourceId))
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
				return errs.WrapError(errs.Error("the resource %s %s was not destroyed ! ", rc.resourceId, rs.Primary.ID))
			}
		}
		return nil
	}
}

func (rc *resourceCheck) callDescribeMethod(rs *terraform.ResourceState) ([]reflect.Value, error) {
	var err error
	if rs.Primary.ID == "" {
		return nil, errs.WrapError(fmt.Errorf("resource ID is not set"))
	}
	serviceP := rc.serviceFunc()
	if rc.describeMethod == "" {
		rc.describeMethod, err = getResourceDescribeMethod(rc.resourceId)
		if err != nil {
			return nil, errs.WrapError(err)
		}
	}
	value := reflect.ValueOf(serviceP)
	typeName := value.Type().String()
	value = value.MethodByName(rc.describeMethod)
	if !value.IsValid() {
		return nil, errs.WrapError(errs.Error("The service type %s does not have method %s", typeName, rc.describeMethod))
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

func (rac *resourceAttrCheck) resourceAttrMapCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := rac.resourceCheck.checkResourceExists()(s)
		if err != nil {
			return errs.WrapError(err)
		}
		return rac.resourceAttr.resourceAttrMapCheck()(s)
	}
}

func (rac *resourceAttrCheck) resourceAttrMapCheckWithCallback(callback func()) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := rac.resourceCheck.checkResourceExists()(s)
		if err != nil {
			return errs.WrapError(err)
		}
		return rac.resourceAttr.resourceAttrMapCheckWithCallback(callback)(s)
	}
}

func (rac *resourceAttrCheck) resourceAttrMapUpdateSet() resourceAttrMapUpdate {
	return func(changeMap map[string]string) resource.TestCheckFunc {
		callback := func() {
			rac.updateCheckMapPair(changeMap)
		}
		return rac.resourceAttrMapCheckWithCallback(callback)
	}
}

func (ra *resourceAttr) updateCheckMapPair(changeMap map[string]string) {
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

func (ra *resourceAttr) resourceAttrMapCheck() resource.TestCheckFunc {
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

func (ra *resourceAttr) resourceAttrMapCheckWithCallback(callback func()) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		callback()
		return ra.resourceAttrMapCheck()(s)
	}
}

func (ra *resourceAttr) resourceAttrMapUpdateSet() resourceAttrMapUpdate {
	return func(changeMap map[string]string) resource.TestCheckFunc {
		callback := func() {
			ra.updateCheckMapPair(changeMap)
		}
		return ra.resourceAttrMapCheckWithCallback(callback)
	}
}

func resourceTestAccConfigFunc(resourceId string,
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

func dataSourceTestAccConfigFunc(resourceId string,
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

func (b *resourceConfig) configUpdate(changeMap map[string]interface{}) {
	newMap := make(map[string]interface{}, len(b.attributeMap))
	for k, v := range b.attributeMap {
		newMap[k] = v
	}
	b.attributeMap = newMap
	if changeMap != nil && len(changeMap) > 0 {
		for rk, rv := range changeMap {
			_, ok := b.attributeMap[rk]
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

func listValue(indentation int, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation+CHILDINDEND)+
			fmt.Sprint(valueConvert(indentation+CHILDINDEND, val.Index(i))))
	}

	return fmt.Sprintf("[\n%s\n%s]", strings.Join(valList, ",\n"), addIndentation(indentation))
}

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

// ---------------------------------------------------------------------------
// sharedClientForRegion — replicated from alicloud_sweeper_test.go for the
// sweeper test functions that need a *connectivity.AliyunClient.
// ---------------------------------------------------------------------------

var dtsEndpoints sync.Map
var dtsSignVersion sync.Map

func sharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey, securityToken string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}

	securityToken = os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &dtsEndpoints,
		SignVersion: &dtsSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	accountId := os.Getenv("ALICLOUD_ACCOUNT_ID")
	if accountId == "" {
		accountId = os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID")
	}
	if accountId != "" {
		conf.AccountId = accountId
	}
	accountType := os.Getenv("ALICLOUD_ACCOUNT_TYPE")
	if accountType == "" {
		accountType = os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")
	}
	if accountType != "" {
		conf.AccountType = accountType
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// testUnitDtsMockMeta builds a *connectivity.AliyunClient with dummy credentials
// so that TestUnit functions can set up gomonkey mocks and run without real
// cloud credentials. Credentials are not validated at construction time, and
// all API calls are intercepted by gomonkey patches on DoRequest, so no actual
// network requests are made. This mirrors the testUnitCdnServiceMockMeta pattern
// in the root package (data_source_alicloud_cdn_service_test.go).
func testUnitDtsMockMeta(t *testing.T) *connectivity.AliyunClient {
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            new(sync.Map),
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	meta, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return meta
}

// ---------------------------------------------------------------------------
// testAccPreCheckWithRegions — replicated from provider_test.go for the
// region-scoped precheck used by DTS acceptance tests.
// ---------------------------------------------------------------------------

var defaultRegionToTest = func() string {
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		return v
	}
	return "cn-beijing"
}()

func testAccPreCheckWithRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	if v := os.Getenv("ALICLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_SECRET_KEY"); v == "" {
		t.Fatal("ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		t.Logf("[WARNING] The region is not set and using cn-beijing as test region")
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
	}
	checkoutSupportedRegions(t, supported, regions)
	defaultRegionToTest = os.Getenv("ALICLOUD_REGION")
}

func checkoutSupportedRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	region := os.Getenv("ALICLOUD_REGION")
	find := false
	backupRegion := string(connectivity.APSouthEast1)
	if region == string(connectivity.APSouthEast1) {
		backupRegion = string(connectivity.EUCentral1)
	}

	checkoutRegion := os.Getenv("CHECKOUT_REGION")
	if checkoutRegion == "true" {
		if region == string(connectivity.Hangzhou) {
			region = string(connectivity.EUCentral1)
			os.Setenv("ALICLOUD_REGION", region)
		}
	}
	backupRegionFind := false
	hangzhouRegionFind := false
	for _, r := range regions {
		if region == string(r) {
			find = true
		}
		if string(r) == backupRegion {
			backupRegionFind = true
		}
		if string(connectivity.Hangzhou) == string(r) {
			hangzhouRegionFind = true
		}
	}

	if (find && !supported) || (!find && supported) {
		if supported {
			if backupRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, backupRegion)
				os.Setenv("ALICLOUD_REGION", backupRegion)
				defaultRegionToTest = backupRegion
				return
			}
			if hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				defaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			if len(regions) > 0 {
				targetRegion := string(regions[0])
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, targetRegion)
				os.Setenv("ALICLOUD_REGION", targetRegion)
				defaultRegionToTest = targetRegion
				return
			}
			t.Skipf("Skipping unsupported region %s. Supported regions: %s.", region, regions)
		} else {
			if !backupRegionFind {
				t.Logf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, backupRegion)
				os.Setenv("ALICLOUD_REGION", backupRegion)
				defaultRegionToTest = backupRegion
				return
			}
			if !hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				defaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			t.Skipf("Skipping unsupported region %s. Unsupported regions: %s.", region, regions)
		}
		t.Skipped()
	}
}
