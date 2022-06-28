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

func TestAccAlicloudECPInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecp_instance.default"
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECPInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECPInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name,
					"description":       name,
					"key_pair_name":     "${alicloud_ecp_key_pair.default.key_pair_name}",
					"security_group_id": "${alicloud_security_group.group.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"image_id":          "android_9_0_0_release_2851157_20211201.vhd",
					"instance_type":     "${local.instance_type}",
					"payment_type":      "PayAsYouGo",
					"vnc_password":      "Cp1234",
					"force":             "true",
					"depends_on":        []string{"alicloud_ecp_key_pair.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"description":       name,
						"key_pair_name":     CHECKSET,
						"security_group_id": CHECKSET,
						"vswitch_id":        CHECKSET,
						"image_id":          CHECKSET,
						"payment_type":      CHECKSET,
						"instance_type":     CHECKSET,
						"resolution":        CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "auto_renew", "force", "vnc_password", "auto_pay", "payment_type", "eip_bandwidth", "period_unit"},
			},
		},
	})
}

var AlicloudECPInstanceMap0 = map[string]string{}

func AlicloudECPInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_ecp_zones" "default" {
}

data "alicloud_ecp_instance_types" "default" {
}

locals {
  count_size               = length(data.alicloud_ecp_zones.default.zones)
  zone_id                  = data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id
  instance_type_count_size = length(data.alicloud_ecp_instance_types.default.instance_types)
  instance_type            = data.alicloud_ecp_instance_types.default.instance_types[local.instance_type_count_size - 1].instance_type
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecp_key_pair" "default" {
  key_pair_name   = var.name
  public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}
`, name)
}
func TestAccAlicloudECPInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecp_instance.default"
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECPInstanceMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudphoneService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcpInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECPInstanceBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name,
					"description":       name,
					"force":             "true",
					"payment_type":      "PayAsYouGo",
					"key_pair_name":     "${alicloud_ecp_key_pair.default.0.key_pair_name}",
					"security_group_id": "${alicloud_security_group.group.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"image_id":          "android_9_0_0_release_2851157_20211201.vhd",
					"instance_type":     "${local.instance_type}",
					"status":            "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"description":       name,
						"payment_type":      "PayAsYouGo",
						"key_pair_name":     CHECKSET,
						"security_group_id": CHECKSET,
						"vswitch_id":        CHECKSET,
						"image_id":          CHECKSET,
						"instance_type":     CHECKSET,
						"status":            "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Stopped",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Stopped",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": "${alicloud_ecp_key_pair.default.1.key_pair_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vnc_password": "Cp1232",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vnc_password": "Cp1232",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name + "update",
					"description":       name + "update",
					"key_pair_name":     "${alicloud_ecp_key_pair.default.0.key_pair_name}",
					"vnc_password":      "Cp1234",
					"payment_type":      "PayAsYouGo",
					"force":             "true",
					"security_group_id": "${alicloud_security_group.group.id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"image_id":          "android_9_0_0_release_2851157_20211201.vhd",
					"instance_type":     "${local.instance_type}",
					"status":            "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name + "update",
						"description":       name + "update",
						"key_pair_name":     CHECKSET,
						"vnc_password":      "Cp1234",
						"security_group_id": CHECKSET,
						"payment_type":      "PayAsYouGo",
						"vswitch_id":        CHECKSET,
						"image_id":          CHECKSET,
						"instance_type":     CHECKSET,
						"status":            "Running",
					}),
				),
			},
		},
	})
}

var AlicloudECPInstanceMap1 = map[string]string{
	"period_unit":   NOSET,
	"period":        NOSET,
	"status":        CHECKSET,
	"auto_renew":    NOSET,
	"force":         CHECKSET,
	"eip_bandwidth": NOSET,
	"amount":        NOSET,
	"auto_pay":      NOSET,
}

func AlicloudECPInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_ecp_zones" "default" {
}

data "alicloud_ecp_instance_types" "default" {
}

locals {
  count_size               = length(data.alicloud_ecp_zones.default.zones)
  zone_id                  = data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id
  instance_type_count_size = length(data.alicloud_ecp_instance_types.default.instance_types)
  instance_type            = data.alicloud_ecp_instance_types.default.instance_types[local.instance_type_count_size - 1].instance_type
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecp_key_pair" "default" {
  count           = 2
  key_pair_name   = join("", [var.name, count.index])
  public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}
`, name)
}

func TestUnitAlicloudECPInstance(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"instance_name":     "RunInstancesValue",
		"description":       "RunInstancesValue",
		"force":             true,
		"payment_type":      "PayAsYouGo",
		"key_pair_name":     "RunInstancesValue",
		"security_group_id": "RunInstancesValue",
		"vswitch_id":        "RunInstancesValue",
		"image_id":          "RunInstancesValue",
		"instance_type":     "RunInstancesValue",
		"status":            "Running",
		"auto_pay":          true,
		"auto_renew":        true,
		"eip_bandwidth":     100,
		"period":            "1",
		"period_unit":       "Month",
		"resolution":        "RunInstancesValue",
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
		// ListInstances
		"Instances": map[string]interface{}{
			"Instance": []interface{}{
				map[string]interface{}{
					"InstanceId":      "RunInstancesValue",
					"Status":          "Running",
					"Description":     "RunInstancesValue",
					"ImageId":         "RunInstancesValue",
					"InstanceName":    "RunInstancesValue",
					"InstanceType":    "RunInstancesValue",
					"ChargeType":      "PostPaid",
					"KeyPairName":     "RunInstancesValue",
					"Resolution":      "RunInstancesValue",
					"SecurityGroupId": "RunInstancesValue",
					"VpcAttributes": map[string]interface{}{
						"VSwitchId": "RunInstancesValue",
					},
				},
			},
		},
		"InstanceIds": map[string]interface{}{
			"InstanceId": []interface{}{
				"RunInstancesValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// RunInstances
		"InstanceIds": map[string]interface{}{
			"InstanceId": []interface{}{
				"RunInstancesValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecp_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudphoneClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcpInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// ListInstances Response
		"InstanceIds": map[string]interface{}{
			"InstanceId": []interface{}{
				"RunInstancesValue",
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RunInstances" {
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
		err := resourceAlicloudEcpInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudphoneClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcpInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateInstanceAttribute
	attributesDiff := map[string]interface{}{
		"instance_name": "UpdateInstanceAttributeValue",
		"description":   "UpdateInstanceAttributeValue",
		"key_pair_name": "UpdateInstanceAttributeValue",
		"vnc_password":  "UpdateInstanceAttributeValue",
	}
	diff, err := newInstanceDiff("alicloud_ecp_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListInstances Response
		"Instances": map[string]interface{}{
			"Instance": []interface{}{
				map[string]interface{}{
					"Description":  "UpdateInstanceAttributeValue",
					"InstanceName": "UpdateInstanceAttributeValue",
					"KeyPairName":  "UpdateInstanceAttributeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateInstanceAttribute" {
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
		err := resourceAlicloudEcpInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// StopInstances
	attributesDiff = map[string]interface{}{
		"status": "Stopped",
	}
	diff, err = newInstanceDiff("alicloud_ecp_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListInstances Response
		"Instances": map[string]interface{}{
			"Instance": []interface{}{
				map[string]interface{}{
					"Status": "Stopped",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "StopInstances" {
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
		err := resourceAlicloudEcpInstanceUpdate(dExisted, rawClient)
		patches.Reset()

		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecp_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListInstances" {
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
		err := resourceAlicloudEcpInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudphoneClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEcpInstanceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "CloudPhoneInstances.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteInstances" {
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
		err := resourceAlicloudEcpInstanceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil", "CloudPhoneInstances.NotFound":
			assert.Nil(t, err)
		}
	}

}
