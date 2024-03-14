package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECSInvocation_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_invocation.default"
	ra := resourceAttrInit(resourceId, AlicloudECSInvocationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsinvocation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSInvocationBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": []string{"${alicloud_instance.default.id}"},
					"command_id":  "${alicloud_ecs_command.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id.#": "1",
						"command_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"windows_password_name"},
			},
		},
	})
}

func TestAccAlicloudECSInvocation_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_invocation.default"
	ra := resourceAttrInit(resourceId, AlicloudECSInvocationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsinvocation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSInvocationBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": []string{"${alicloud_instance.default.id}"},
					"command_id":  "${alicloud_ecs_command.default.id}",
					"repeat_mode": "Period",
					"timed":       "true",
					"frequency":   "0 15 10 ? * *",
					"parameters": map[string]string{
						"name": "exampleName",
					},
					"username":              "root",
					"windows_password_name": "axtSecretPassword",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id.#":   "1",
						"command_id":      CHECKSET,
						"repeat_mode":     "Period",
						"timed":           "true",
						"frequency":       "0 15 10 ? * *",
						"parameters.%":    "1",
						"parameters.name": "exampleName",
						"username":        "root",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"windows_password_name"},
			},
		},
	})
}

var AlicloudECSInvocationMap0 = map[string]string{
	"instance_id.#": CHECKSET,
}

func AlicloudECSInvocationBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_ecs_command" "default" {
	name              = var.name
	command_content   = "bHMK"
	description       = "For Terraform Test"
	type              = "RunShellScript"
	working_dir       = "/root"
}
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name              = "${var.name}"
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

resource "alicloud_instance" "default" {
	vswitch_id = "${alicloud_vswitch.default.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
}
`, name)
}

func AlicloudECSInvocationBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name              = "${var.name}"
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

resource "alicloud_instance" "default" {
	vswitch_id = "${alicloud_vswitch.default.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
}

resource "alicloud_ecs_command" "default" {
	name              = var.name
	command_content   = "ZWNobyBoZWxsbyx7e25hbWV9fQ=="
	description       = "For Terraform Test"
	type              = "RunShellScript"
	working_dir       = "/root"
	enable_parameter  = true
}
`, name)
}

func TestUnitAlicloudECSInvocation(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_id": []string{"CreateECSInvocationValue"},
		"command_id":  "CreateECSInvocationValue",
		"repeat_mode": "CreateECSInvocationValue",
		"timed":       true,
		"frequency":   "CreateECSInvocationValue",
		"parameters": map[string]string{
			"name": "CreateECSInvocationValue",
		},
		"username":              "CreateECSInvocationValue",
		"windows_password_name": "CreateECSInvocationValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"TotalCount": 1,
		"PageSize":   10,
		"PageNumber": 1,
		"Invocations": map[string]interface{}{
			"Invocation": []interface{}{
				map[string]interface{}{
					"InvocationStatus": "Scheduled",
					"Parameters":       "{\"name\":\"CreateECSInvocationValue\"}",
					"Timed":            true,
					"InvokeInstances": map[string]interface{}{
						"InvokeInstance": []interface{}{
							map[string]interface{}{
								"Dropped":              0,
								"InvocationStatus":     "Scheduled",
								"InstanceId":           "CreateECSInvocationValue",
								"Timed":                true,
								"InstanceInvokeStatus": "Running",
								"ErrorInfo":            "",
								"StartTime":            "",
								"Repeats":              0,
								"FinishTime":           "",
								"Output":               "",
								"CreationTime":         "2022-05-17T10:28:05Z",
								"UpdateTime":           "2022-05-17T10:28:05Z",
								"ErrorCode":            "",
								"StopTime":             "",
							},
						},
					},
					"CommandContent": "CreateECSInvocationValue",
					"RepeatMode":     "CreateECSInvocationValue",
					"InvokeStatus":   "Running",
					"CommandType":    "RunShellScript",
					"Username":       "CreateECSInvocationValue",
					"CreationTime":   "2022-05-17T10:28:05Z",
					"Frequency":      "CreateECSInvocationValue",
					"CommandId":      "CreateECSInvocationValue",
					"InvokeId":       "CreateECSInvocationValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"InvokeId": "CreateECSInvocationValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecs_invocation", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcsInvocationCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "InvokeCommand" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEcsInvocationCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			//assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_ecs_invocation", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_invocation"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeInvocations" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEcsInvocationRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcsInvocationDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "StopInvocation" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEcsInvocationDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
