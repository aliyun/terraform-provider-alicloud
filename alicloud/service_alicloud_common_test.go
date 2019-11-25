package alicloud

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

/**
	This file aims to provide some const test cases and applied them for several specified resource or data source's test cases.
These common test cases are used to creating some dependence resources, like vpc, vswitch and security group.
*/

// be used to check attribute map value
const (
	NOSET      = "#NOSET"       // be equivalent to method "TestCheckNoResourceAttrSet"
	CHECKSET   = "#CHECKSET"    // "TestCheckResourceAttrSet"
	REMOVEKEY  = "#REMOVEKEY"   // remove checkMap key
	REGEXMATCH = "#REGEXMATCH:" // "TestMatchResourceAttr" ,the map name/key like `"attribute" : REGEXMATCH + "attributeString"`
)

const (
	// indentation symbol
	INDENTATIONSYMBOL = " "

	// child field indend number
	CHILDINDEND = 2
)

// get a function that change checkMap pairs for a series test step
type resourceAttrMapUpdate func(map[string]string) resource.TestCheckFunc

// get a function that change attributeMap pairs for a series test step
type ResourceTestAccConfigFunc func(map[string]interface{}) string

// check the existence of resource
type resourceCheck struct {
	// IDRefreshName, like "alicloud_instance.foo"
	resourceId string

	// The response of the service method DescribeXXX
	resourceObject interface{}

	// The resource service client type, like DnsService, VpcService
	serviceFunc func() interface{}

	// service describe method name
	describeMethod string
}

func resourceCheckInit(resourceId string, resourceObject interface{}, serviceFunc func() interface{}) *resourceCheck {
	return &resourceCheck{
		resourceId:     resourceId,
		resourceObject: resourceObject,
		serviceFunc:    serviceFunc,
	}
}

func resourceCheckInitWithDescribeMethod(resourceId string, resourceObject interface{}, serviceFunc func() interface{}, describeMethod string) *resourceCheck {
	return &resourceCheck{
		resourceId:     resourceId,
		resourceObject: resourceObject,
		serviceFunc:    serviceFunc,
		describeMethod: describeMethod,
	}
}

// check attribute only
type resourceAttr struct {
	resourceId string
	checkMap   map[string]string
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

// check the existence and attribute of the resource at the same time
type resourceAttrCheck struct {
	*resourceCheck
	*resourceAttr
}

func resourceAttrCheckInit(rc *resourceCheck, ra *resourceAttr) *resourceAttrCheck {
	return &resourceAttrCheck{
		resourceCheck: rc,
		resourceAttr:  ra,
	}
}

// check the resource existence by invoking DescribeXXX method of service and assign *resourceCheck.resourceObject value,
// the service is returned by invoking *resourceCheck.serviceFunc
func (rc *resourceCheck) checkResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		rs, ok := s.RootModule().Resources[rc.resourceId]
		if !ok {
			return WrapError(fmt.Errorf("can't find resource by id: %s", rc.resourceId))

		}
		outValue, err := rc.callDescribeMethod(rs)
		if err != nil {
			return WrapError(err)
		}
		errorValue := outValue[1]
		if !errorValue.IsNil() {
			return WrapError(fmt.Errorf("Checking resource %s %s exists error:%s ", rc.resourceId, rs.Primary.ID, errorValue.Interface().(error).Error()))
		}
		if reflect.TypeOf(rc.resourceObject).Elem().String() == outValue[0].Type().String() {
			reflect.ValueOf(rc.resourceObject).Elem().Set(outValue[0])
			return nil
		} else {
			return WrapError(fmt.Errorf("The response object type expected *%s, got %s ", outValue[0].Type().String(), reflect.TypeOf(rc.resourceObject).String()))
		}
	}
}

// check the resource destroy
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
			return WrapError(Error("The resourceId %s is not correct and it should prefix with alicloud_", rc.resourceId))
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
					if NotFoundError(err) {
						continue
					}
					return WrapError(err)
				}
			} else {
				return WrapError(Error("the resource %s %s was not destroyed ! ", rc.resourceId, rs.Primary.ID))
			}
		}
		return nil
	}
}

// invoking DescribeXXX method of service
func (rc *resourceCheck) callDescribeMethod(rs *terraform.ResourceState) ([]reflect.Value, error) {
	var err error
	if rs.Primary.ID == "" {
		return nil, WrapError(fmt.Errorf("resource ID is not set"))
	}
	serviceP := rc.serviceFunc()
	if rc.describeMethod == "" {
		rc.describeMethod, err = getResourceDescribeMethod(rc.resourceId)
		if err != nil {
			return nil, WrapError(err)
		}
	}
	value := reflect.ValueOf(serviceP)
	typeName := value.Type().String()
	value = value.MethodByName(rc.describeMethod)
	if !value.IsValid() {
		return nil, WrapError(Error("The service type %s does not have method %s", typeName, rc.describeMethod))
	}
	inValue := []reflect.Value{reflect.ValueOf(rs.Primary.ID)}
	return value.Call(inValue), nil
}

func getResourceDescribeMethod(resourceId string) (string, error) {
	start := strings.Index(resourceId, "alicloud_")
	if start < 0 {
		return "", WrapError(fmt.Errorf("the parameter \"name\" don't contain string \"alicloud_\""))
	}
	start += len("alicloud_")
	end := strings.Index(resourceId[start:], ".") + start
	if end < 0 {
		return "", WrapError(fmt.Errorf("the parameter \"name\" don't contain string \".\""))
	}
	strs := strings.Split(resourceId[start:end], "_")
	describeName := "Describe"
	for _, str := range strs {
		describeName = describeName + strings.Title(str)
	}
	return describeName, nil
}

// check attribute func and check resource exist
func (rac *resourceAttrCheck) resourceAttrMapCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := rac.resourceCheck.checkResourceExists()(s)
		if err != nil {
			return WrapError(err)
		}
		return rac.resourceAttr.resourceAttrMapCheck()(s)
	}
}

// execute the callback before check attribute and check resource exist
func (rac *resourceAttrCheck) resourceAttrMapCheckWithCallback(callback func()) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := rac.resourceCheck.checkResourceExists()(s)
		if err != nil {
			return WrapError(err)
		}
		return rac.resourceAttr.resourceAttrMapCheckWithCallback(callback)(s)
	}
}

// get resourceAttrMapUpdate for a series test step and check resource exist
func (rac *resourceAttrCheck) resourceAttrMapUpdateSet() resourceAttrMapUpdate {
	return func(changeMap map[string]string) resource.TestCheckFunc {
		callback := func() {
			rac.updateCheckMapPair(changeMap)
		}
		return rac.resourceAttrMapCheckWithCallback(callback)
	}
}

// make a new map and copy from the old field checkMap, then update it according to the changeMap
func (ra *resourceAttr) updateCheckMapPair(changeMap map[string]string) {
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
func (ra *resourceAttr) resourceAttrMapCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[ra.resourceId]
		if !ok {
			return WrapError(fmt.Errorf("can't find resource by id: %s", ra.resourceId))
		}
		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("resource ID is not set"))
		}

		if ra.checkMap == nil || len(ra.checkMap) == 0 {
			return WrapError(fmt.Errorf("the parameter \"checkMap\" is nil or empty"))
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
		return WrapError(fmt.Errorf(strings.Join(errorStrSlice, "\n")))
	}
}

// execute the callback before check attribute
func (ra *resourceAttr) resourceAttrMapCheckWithCallback(callback func()) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		callback()
		return ra.resourceAttrMapCheck()(s)
	}
}

// get resourceAttrMapUpdate for a series test step
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

// be used for generate testcase step config
type resourceConfig struct {
	// the resource name
	name string

	resourceId string

	// store attribute value that primary resource
	attributeMap map[string]interface{}

	// generate assistant test config
	configDependence func(name string) string
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
			if strValue, isCost := rv.(string); ok && isCost && strValue == REMOVEKEY {
				delete(b.attributeMap, rk)
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
		return assistantConfig + primaryConfig + valueConvert(0, reflect.ValueOf(b.attributeMap))
	}
}

// deal with the parameter common method
func valueConvert(indentation int, val reflect.Value) string {
	switch val.Kind() {
	case reflect.Interface:
		return valueConvert(indentation, reflect.ValueOf(val.Interface()))
	case reflect.String:
		return fmt.Sprintf("\"%s\"", val.String())
	case reflect.Slice:
		return listValue(indentation, val)
	case reflect.Map:
		return mapValue(indentation, val)
	default:
		log.Panicf("the map value must be string  map or slice type! %s", val)
	}
	return ""
}

// deal with list parameter
func listValue(indentation int, val reflect.Value) string {
	var valList []string
	for i := 0; i < val.Len(); i++ {
		valList = append(valList, addIndentation(indentation+CHILDINDEND)+
			valueConvert(indentation+CHILDINDEND, val.Index(i)))
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
		line = fmt.Sprintf(`%s%s = %s`, addIndentation(indentation+CHILDINDEND), keyV.String(),
			valueConvert(indentation+len(keyV.String())+CHILDINDEND+3, val.MapIndex(keyV)))
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
// dataSourceAttr has 3 field ,incloud resourceId  existMapFunc fakeMapFunc, every dataSource test can use only one
type dataSourceAttr struct {
	// IDRefreshName, like "data.alicloud_dns_records.record"
	resourceId string

	// get existMap function
	existMapFunc func(rand int) map[string]string

	// get fakeMap function
	fakeMapFunc func(rand int) map[string]string
}

// get exist and empty resourceAttrMapUpdate function
func (dsa *dataSourceAttr) checkDataSourceAttr(rand int) (exist, empty resourceAttrMapUpdate) {
	exist = resourceAttrInit(dsa.resourceId, dsa.existMapFunc(rand)).resourceAttrMapUpdateSet()
	empty = resourceAttrInit(dsa.resourceId, dsa.fakeMapFunc(rand)).resourceAttrMapUpdateSet()
	return
}

// according to configs generate step list and execute the test
func (dsa *dataSourceAttr) dataSourceTestCheck(t *testing.T, rand int, configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps:     steps,
	})
}

// according to configs generate step list and execute the test with preCheck
func (dsa *dataSourceAttr) dataSourceTestCheckWithPreCheck(t *testing.T, rand int, preCheck func(), configs ...dataSourceTestAccConfig) {
	var steps []resource.TestStep
	for _, conf := range configs {
		steps = append(steps, conf.buildDataSourceSteps(t, dsa, rand)...)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps:     steps,
	})
}

// per schema attribute test config
type dataSourceTestAccConfig struct {
	// be equal to testCase config string,but the result has only one record
	existConfig string

	// if the dataSourceAttr.existMapFunc returned map value not match we want, existChangMap can alter checkMap for existConfig
	existChangMap map[string]string

	// be equal to testCase config string,but the result is empty
	fakeConfig string

	// if the dataSourceAttr.fakeMapFunc returned map value not match we want, fakeChangMap can alter checkMap for fakeConfig
	fakeChangMap map[string]string
}

// build test cases for each attribute
func (conf *dataSourceTestAccConfig) buildDataSourceSteps(t *testing.T, info *dataSourceAttr, rand int) []resource.TestStep {
	testAccCheckExist, testAccCheckEmpty := info.checkDataSourceAttr(rand)
	var steps []resource.TestStep
	if conf.existConfig != "" {
		step := resource.TestStep{
			Config: conf.existConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckExist(conf.existChangMap),
			),
		}
		steps = append(steps, step)
	}
	if conf.fakeConfig != "" {
		step := resource.TestStep{
			Config: conf.fakeConfig,
			Check: resource.ComposeTestCheckFunc(
				testAccCheckEmpty(conf.fakeChangMap),
			),
		}
		steps = append(steps, step)
	}
	return steps
}

func (s *VpcService) needSweepVpc(vpcId, vswitchId string) (bool, error) {
	if vpcId == "" && vswitchId != "" {
		object, err := s.DescribeVSwitch(vswitchId)
		if err != nil && !NotFoundError(err) {
			return false, WrapError(err)
		}
		name := strings.ToLower(object.VSwitchName)
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the vswitch (%s (%s)).", object.VSwitchId, object.VSwitchName)
			return true, nil
		}
		vpcId = object.VpcId
	}
	if vpcId != "" {
		object, err := s.DescribeVpc(vpcId)
		if err != nil {
			if NotFoundError(err) {
				return false, nil
			}
			return false, WrapError(err)
		}
		name := strings.ToLower(object.VpcName)
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the VPC (%s (%s)).", object.VpcId, object.VpcName)
			return true, nil
		}
	}
	return false, nil
}

func (s *VpcService) sweepVpc(id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Vpc %s ...", id)
	request := vpc.CreateDeleteVpcRequest()
	request.VpcId = id
	_, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteVpc(request)
	})

	return WrapError(err)
}

func (s *VpcService) sweepVSwitch(id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Vswitch %s ...", id)
	request := vpc.CreateDeleteVSwitchRequest()
	request.VSwitchId = id
	_, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteVSwitch(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return WrapError(err)
}

func (s *VpcService) sweepNatGateway(id string) error {
	if id == "" {
		return nil
	}

	log.Printf("[INFO] Deleting Nat Gateway %s ...", id)
	request := vpc.CreateDeleteNatGatewayRequest()
	request.NatGatewayId = id
	request.Force = requests.NewBoolean(true)
	_, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DeleteNatGateway(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return WrapError(err)
}

func (s *EcsService) sweepSecurityGroup(id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Deleting Security Group %s ...", id)
	request := ecs.CreateDeleteSecurityGroupRequest()
	request.SecurityGroupId = id
	_, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteSecurityGroup(request)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return WrapError(err)
}

func (s *SlbService) sweepSlb(id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Set SLB DeleteProtection to off before deleting %s ...", id)
	request := slb.CreateSetLoadBalancerDeleteProtectionRequest()
	request.LoadBalancerId = id
	request.DeleteProtection = "off"
	_, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.SetLoadBalancerDeleteProtection(request)
	})
	if err != nil {
		log.Printf("[ERROR] Set SLB %s DeleteProtection to off failed.", id)
	}
	log.Printf("[DEBUG] Deleting SLB %s ...", id)
	delRequest := slb.CreateDeleteLoadBalancerRequest()
	delRequest.LoadBalancerId = id
	_, err = s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DeleteLoadBalancer(delRequest)
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return WrapError(err)
}

const EcsInstanceCommonNoZonesTestCase = `
data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*_64"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  name              = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}
`

const EcsInstanceCommonTestCase = `
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*_64"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}
`

const RdsCommonTestCase = `
data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}
`
const KVStoreCommonTestCase = `
data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${lookup(data.alicloud_zones.default.zones[(length(data.alicloud_zones.default.zones)-1)%length(data.alicloud_zones.default.zones)], "id")}"
  name              = "${var.name}"
}
`

const DBMultiAZCommonTestCase = `
data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
  multi = true
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.multi_zone_ids[0]}"
  name              = "${var.name}"
}
`

const ElasticsearchInstanceCommonTestCase = `
data "alicloud_zones" "default" {
    available_resource_creation = "${var.creation}"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-1], "id")}"
  name              = "${var.name}"
}`

const SlbVpcCommonTestCase = `
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}
`

const EmrCommonTestCase = `
data "alicloud_emr_main_versions" "default" {
}

data "alicloud_emr_instance_types" "default" {
    destination_resource = "InstanceType"
    cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
    support_local_storage = false
    instance_charge_type = "PostPaid"
    support_node_type = ["MASTER", "CORE"]
}

data "alicloud_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_emr_instance_types.default.types.0.zone_id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com", 
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
    description = "this is a role test."
    force = true
}
`

const EmrGatewayTestCase = `
data "alicloud_emr_main_versions" "default" {
}

data "alicloud_emr_instance_types" "default" {
    destination_resource = "InstanceType"
    cluster_type = "GATEWAY"
    support_local_storage = false
    instance_charge_type = "PostPaid"
    support_node_type = ["GATEWAY"]
}

data "alicloud_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = "GATEWAY"
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = "GATEWAY"
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.default.types.0.id
	zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_emr_instance_types.default.types.0.zone_id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com", 
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
    description = "this is a role test."
    force = true
}

resource "alicloud_emr_cluster" "default" {
    name = "${var.name}"

    emr_ver = data.alicloud_emr_main_versions.default.main_versions.0.emr_version

    cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0

    host_group {
        host_group_name = "master_group"
        host_group_type = "MASTER"
        node_count = "2"
        instance_type = data.alicloud_emr_instance_types.default.types.0.id
        disk_type = data.alicloud_emr_disk_types.data_disk.types.0.value
        disk_capacity = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
        disk_count = "1"
        sys_disk_type = data.alicloud_emr_disk_types.system_disk.types.0.value
		sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
    }

	host_group {
        host_group_name = "core_group"
        host_group_type = "CORE"
        node_count = "2"
        instance_type = data.alicloud_emr_instance_types.default.types.0.id
        disk_type = data.alicloud_emr_disk_types.data_disk.types.0.value
        disk_capacity = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
        disk_count = "4"
        sys_disk_type = data.alicloud_emr_disk_types.system_disk.types.0.value
        sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
    }

    high_availability_enable = true
    option_software_list = ["HBASE","PRESTO",]
    zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
    security_group_id = alicloud_security_group.default.id
    is_open_public_ip = true
    charge_type = "PostPaid"
    vswitch_id = alicloud_vswitch.default.id
    user_defined_emr_ecs_role = alicloud_ram_role.default.name
    ssh_enable = true
    master_pwd = "ABCtest1234!"
}
`
const EmrLocalStorageTestCase = `
data "alicloud_emr_main_versions" "default" {
}

data "alicloud_emr_instance_types" "local_disk" {
    destination_resource = "InstanceType"
    cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
    support_local_storage = true
    instance_charge_type = "PostPaid"
    support_node_type = ["CORE"]
}

data "alicloud_emr_instance_types" "cloud_disk" {
    destination_resource = "InstanceType"
    cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
    instance_charge_type = "PostPaid"
    support_node_type = ["MASTER"]
}

data "alicloud_emr_disk_types" "data_disk" {
	destination_resource = "DataDisk"
	cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.cloud_disk.types.0.id
	zone_id = data.alicloud_emr_instance_types.cloud_disk.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
	destination_resource = "SystemDisk"
	cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
	instance_charge_type = "PostPaid"
	instance_type = data.alicloud_emr_instance_types.cloud_disk.types.0.id
	zone_id = data.alicloud_emr_instance_types.cloud_disk.types.0.zone_id
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_emr_instance_types.cloud_disk.types.0.zone_id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_ram_role" "default" {
	name = "${var.name}"
	document = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com", 
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
    description = "this is a role test."
    force = true
}
`

const SlbListenerCommonTestCase = `
variable "ip_version" {
  default = "ipv4"
}	
resource "alicloud_slb" "default" {
  name = "${var.name}"
  internet_charge_type = "PayByTraffic"
  address_type = "internet"
  specification = "slb.s1.small"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
  }
  entry_list {
      entry="168.10.10.0/24"
      comment="second"
  }
}
`
const SlbListenerVserverCommonTestCase = `
data "alicloud_zones" "default" {
  available_disk_category = "cloud_efficiency"
  available_resource_creation= "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*_64"
  most_recent = true
  owners = "system"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${alicloud_security_group.default.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
}

resource "alicloud_slb_master_slave_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers {
      server_id = "${alicloud_instance.default.0.id}"
      port = 80
      weight = 100
      server_type = "Master"
  }
  servers {
      server_id = "${alicloud_instance.default.1.id}"
      port = 80
      weight = 100
      server_type = "Slave"
  }
}
`
