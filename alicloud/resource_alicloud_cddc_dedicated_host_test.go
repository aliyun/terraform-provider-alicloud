package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCddcDedicatedHost_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cddcdedicatedhost%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostBasicDependence0)
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
					"dedicated_host_group_id": "${local.dedicated_host_group_id}",
					"host_class":              "${data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code}",
					"zone_id":                 "${data.alicloud_cddc_zones.default.ids.0}",
					"vswitch_id":              "${data.alicloud_vswitches.default.ids.0}",
					"payment_type":            "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_group_id": CHECKSET,
						"host_class":              CHECKSET,
						"zone_id":                 CHECKSET,
						"vswitch_id":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allocation_status": "Allocatable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_status": "Allocatable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_class": "${data.alicloud_cddc_host_ecs_level_infos.default.infos.1.res_class_code}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_class": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "CDDC_DEDICATED",
						"For":     "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "CDDC_DEDICATED",
						"tags.For":     "TF",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name":         "${var.name}_update",
					"allocation_status": "Suspended",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "CDDC_DEDICATED",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name":         name + "_update",
						"allocation_status": "Suspended",
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.For":          "CDDC_DEDICATED",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_category", "payment_type", "used_time", "period", "auto_renew", "os_password"},
			},
		},
	})
}

func TestAccAlicloudCddcDedicatedHost_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cddc_dedicated_host.default"
	ra := resourceAttrInit(resourceId, AlicloudCDDCDedicatedHostMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CddcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCddcDedicatedHost")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cddcdedicatedhost%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCDDCDedicatedHostBasicDependence0)
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
					"dedicated_host_group_id": "${local.dedicated_host_group_id}",
					"host_class":              "${data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code}",
					"auto_renew":              "false",
					"zone_id":                 "${data.alicloud_cddc_zones.default.ids.0}",
					"vswitch_id":              "${data.alicloud_vswitches.default.ids.0}",
					"payment_type":            "Subscription",
					"host_name":               "${var.name}",
					"period":                  "Month",
					"used_time":               "1",
					"image_category":          "AliLinux",
					"os_password":             "Password1234.",
					"allocation_status":       "Allocatable",
					"tags": map[string]string{
						"Created": "CDDC_DEDICATED",
						"For":     "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dedicated_host_group_id": CHECKSET,
						"host_name":               name,
						"host_class":              CHECKSET,
						"zone_id":                 CHECKSET,
						"vswitch_id":              CHECKSET,
						"allocation_status":       "Allocatable",
						"tags.%":                  "2",
						"tags.Created":            "CDDC_DEDICATED",
						"tags.For":                "TF",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_category", "payment_type", "used_time", "period", "auto_renew", "os_password"},
			},
		},
	})
}

var AlicloudCDDCDedicatedHostMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudCDDCDedicatedHostBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

data "alicloud_cddc_dedicated_host_groups" "default" {
  engine     = "MySQL"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
	count = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? 0 : 1
	engine = "MySQL"
	vpc_id = data.alicloud_vpcs.default.ids.0
	cpu_allocation_ratio = 101
	mem_allocation_ratio = 50
	disk_allocation_ratio = 200
	allocation_policy = "Evenly"
	host_replace_policy = "Manual"
	dedicated_host_group_desc = var.name
	open_permission = true
}
locals {
	dedicated_host_group_id = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.ids.0 : concat(alicloud_cddc_dedicated_host_group.default[*].id, [""])[0]
}
`, name)
}

func TestUnitAlicloudCddcDedicatedHost(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cddc_dedicated_host"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cddc_dedicated_host"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"auto_renew":              true,
		"dedicated_host_group_id": "CreateDedicatedHostValue",
		"host_class":              "CreateDedicatedHostValue",
		"host_name":               "CreateDedicatedHostValue",
		"image_category":          "CreateDedicatedHostValue",
		"os_password":             "CreateDedicatedHostValue",
		"payment_type":            "Subscription",
		"period":                  "CreateDedicatedHostValue",
		"used_time":               60,
		"zone_id":                 "CreateDedicatedHostValue",
		"vswitch_id":              "CreateDedicatedHostValue",
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
		// DescribeDedicatedHostAttribute
		"DedicatedHostGroupId": "CreateDedicatedHostValue",
		"HostClass":            "CreateDedicatedHostValue",
		"HostName":             "CreateDedicatedHostValue",
		"HostStatus":           "1",
		"VSwitchId":            "CreateDedicatedHostValue",
		"ZoneId":               "CreateDedicatedHostValue",
		"AllocationStatus":     "1",
		"DedicatedHostId":      "CreateDedicatedHostValue",
		//ListTagResources
		"TagResources": []interface{}{
			map[string]interface{}{
				"TagKey":   "tag_key",
				"TagValue": "tag_value",
			},
		},
		"Tags": map[string]interface{}{
			"key": "value",
		},
		"DedicateHostList": map[string]interface{}{
			"DedicateHostList": []interface{}{
				map[string]interface{}{
					"DedicatedHostId": "CreateDedicatedHostValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDedicatedHost
		"DedicateHostList": map[string]interface{}{
			"DedicateHostList": []interface{}{
				map[string]interface{}{
					"DedicatedHostId": "CreateDedicatedHostValue",
				},
			},
		},
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cddc_dedicated_host", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCddcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCddcDedicatedHostCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDedicatedHostAttribute Response
		"DedicateHostList": map[string]interface{}{
			"DedicateHostList": []interface{}{
				map[string]interface{}{
					"DedicatedHostId": "CreateDedicatedHostValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDedicatedHost" {
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
		err := resourceAlicloudCddcDedicatedHostCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cddc_dedicated_host"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCddcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCddcDedicatedHostUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDedicatedHostAttribute
	attributesDiff := map[string]interface{}{
		"tags": map[string]string{
			"key": "value",
		},
		"host_name":         "ModifyDedicatedHostAttributeValue",
		"allocation_status": "Allocatable",
	}
	diff, err := newInstanceDiff("alicloud_cddc_dedicated_host", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cddc_dedicated_host"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHostAttribute Response
		"Tags": map[string]interface{}{
			"key": "value",
		},
		"TagResources": []interface{}{
			map[string]interface{}{
				"TagKey":   "key",
				"TagValue": "value",
			},
		},
		"HostName":         "ModifyDedicatedHostAttributeValue",
		"AllocationStatus": "Allocatable",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDedicatedHostAttribute" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudCddcDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cddc_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyDedicatedHostClass
	attributesDiff = map[string]interface{}{
		"host_class": "ModifyDedicatedHostAttributeValue",
	}
	diff, err = newInstanceDiff("alicloud_cddc_dedicated_host", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cddc_dedicated_host"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHostAttribute Response
		"HostClass": "ModifyDedicatedHostAttributeValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDedicatedHostClass" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudCddcDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cddc_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	diff, err = newInstanceDiff("alicloud_cddc_dedicated_host", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cddc_dedicated_host"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDedicatedHostAttribute" {
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
		err := resourceAlicloudCddcDedicatedHostRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudCddcDedicatedHostDelete(dExisted, rawClient)
	assert.Nil(t, err)

}
