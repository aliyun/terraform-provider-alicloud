package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"log"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudConfigDeliveryChannel_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_delivery_channel.default"
	ra := resourceAttrInit(resourceId, ConfigDeliveryChannelMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDeliveryChannel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigDeliveryChannel%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigDeliveryChannelBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.bucket}",
					"delivery_channel_type":            "OSS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "OSS",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change the role arn must using resource manager master account.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"delivery_channel_assume_role_arn": "acs:ram::118272523xxxxxxx:role/aliyunserviceroleforconfig",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"delivery_channel_assume_role_arn": "acs:ram::118272523xxxxxxx:role/aliyunserviceroleforconfig",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_target_arn": "${local.bucket_change}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_target_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.bucket}",
					"delivery_channel_type":            "OSS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "OSS",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudConfigDeliveryChannel_MNS(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_delivery_channel.default"
	ra := resourceAttrInit(resourceId, ConfigDeliveryChannelMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDeliveryChannel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigDeliveryChannel%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigDeliveryChannelBasicdependence_MNS)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.mns}",
					"delivery_channel_type":            "MNS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "MNS",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_condition": deliveryChannelCondition,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_condition": strings.Replace(strings.Replace(deliveryChannelCondition, `\n`, "\n", -1), `\"`, "\"", -1),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_target_arn": "${local.mns_change}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_target_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.mns}",
					"delivery_channel_type":            "MNS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "MNS",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudConfigDeliveryChannel_SLS(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_delivery_channel.default"
	ra := resourceAttrInit(resourceId, ConfigDeliveryChannelMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDeliveryChannel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigDeliveryChannel%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigDeliveryChannelBasicdependence_SLS)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.sls}",
					"delivery_channel_type":            "SLS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "SLS",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_target_arn": "${local.sls_change}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_target_arn": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_assume_role_arn": "${local.role_arn}",
					"delivery_channel_target_arn":      "${local.sls}",
					"delivery_channel_type":            "SLS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_assume_role_arn": CHECKSET,
						"delivery_channel_target_arn":      CHECKSET,
						"delivery_channel_type":            "SLS",
					}),
				),
			},
		},
	})
}

var ConfigDeliveryChannelMap = map[string]string{}

// Because the bucket cannot be deleted after being used by the delivery channel.
// Use pre-created Oss bucket in this test.
func ConfigDeliveryChannelBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

locals {
  uid          	   = data.alicloud_account.this.id
  role_arn         = data.alicloud_ram_roles.this.roles.0.arn
  bucket	       = format("acs:oss:cn-beijing:%%s:ci-test-bucket-for-config",local.uid)
  bucket_change	   = format("acs:oss:cn-beijing:%%s:ci-test-bucket-for-config-update",local.uid)
}

data "alicloud_account" "this" {}

data "alicloud_ram_roles" "this" {
  name_regex = "^AliyunServiceRoleForConfig$"
}

`, name)
}

func ConfigDeliveryChannelBasicdependence_MNS(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

locals {
  uid          	   = data.alicloud_account.this.id
  role_arn         = data.alicloud_ram_roles.this.roles.0.arn
  mns	       	   = format("acs:mns:%[2]s:%%s:/topics/%%s",local.uid,alicloud_mns_topic.default.name)
  mns_change	   = format("acs:mns:%[2]s:%%s:/topics/%%s",local.uid,alicloud_mns_topic.change.name)
}

resource "alicloud_mns_topic" "default" {
  name = var.name
}
resource "alicloud_mns_topic" "change" {
  name = format("%%s-change",var.name)
}

data "alicloud_account" "this" {}

data "alicloud_ram_roles" "this" {
  name_regex = "^AliyunServiceRoleForConfig$"
}

`, name, defaultRegionToTest)
}

func ConfigDeliveryChannelBasicdependence_SLS(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

locals {
  uid          	   = data.alicloud_account.this.id
  role_arn         = data.alicloud_ram_roles.this.roles.0.arn
  sls	       	   = format("acs:log:%[2]s:%%s:project/%%s/logstore/%%s",local.uid,alicloud_log_project.this.name,alicloud_log_store.this.name)
  sls_change	   = format("acs:log:%[2]s:%%s:project/%%s/logstore/%%s",local.uid,alicloud_log_project.this.name,alicloud_log_store.change.name)
}

resource "alicloud_log_project" "this" {
  name = var.name
}
resource "alicloud_log_store" "this" {
  name = var.name
  project = alicloud_log_project.this.name
}
resource "alicloud_log_store" "change" {
  name = format("%%s-change",var.name)
  project = alicloud_log_project.this.name
}

data "alicloud_account" "this" {}

data "alicloud_ram_roles" "this" {
  name_regex = "^AliyunServiceRoleForConfig$"
}

`, strings.ToLower(name), defaultRegionToTest)
}

const deliveryChannelCondition = `[\n{\n\"filterType\":\"ResourceType\",\n\"values\":[\n\"ACS::CEN::CenInstance\",\n],\n\"multiple\":true\n}\n]\n`

func TestUnitAlicloudConfigDeliveryChannel(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_delivery_channel"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_delivery_channel"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"delivery_channel_assume_role_arn": "CreateDeliveryChannelValue",
		"delivery_channel_condition":       "CreateDeliveryChannelValue",
		"delivery_channel_name":            "CreateDeliveryChannelValue",
		"description":                      "CreateDeliveryChannelValue",
		"delivery_channel_target_arn":      "CreateDeliveryChannelValue",
		"delivery_channel_type":            "CreateDeliveryChannelValue",
		"status":                           1,
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
		// DescribeDeliveryChannels
		"DeliveryChannels": []interface{}{
			map[string]interface{}{
				"DeliveryChannelAssumeRoleArn": "CreateDeliveryChannelValue",
				"DeliveryChannelCondition":     "CreateDeliveryChannelValue",
				"DeliveryChannelName":          "CreateDeliveryChannelValue",
				"DeliveryChannelTargetArn":     "CreateDeliveryChannelValue",
				"DeliveryChannelType":          "CreateDeliveryChannelValue",
				"Description":                  "CreateDeliveryChannelValue",
				"Status":                       1,
				"DeliveryChannelId":            "CreateDeliveryChannelValue",
			},
		},
		"DeliveryChannelId": "CreateDeliveryChannelValue",
	}
	CreateMockResponse := map[string]interface{}{
		// PutDeliveryChannel
		"DeliveryChannels": []interface{}{
			map[string]interface{}{
				"DeliveryChannelId": "CreateDeliveryChannelValue",
			},
		},
		"DeliveryChannelId": "CreateDeliveryChannelValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_config_delivery_channel", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigDeliveryChannelCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDeliveryChannels Response
		"DeliveryChannels": []interface{}{
			map[string]interface{}{
				"DeliveryChannelId": "CreateDeliveryChannelValue",
			},
		},
		"DeliveryChannelId": "CreateDeliveryChannelValue",
	}
	errorCodes := []string{"NonRetryableError", "DeliveryChannelSlsUnreachableError", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "PutDeliveryChannel" {
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
		err := resourceAlicloudConfigDeliveryChannelCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_delivery_channel"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigDeliveryChannelUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// PutDeliveryChannel
	attributesDiff := map[string]interface{}{
		"delivery_channel_assume_role_arn": "UpdateDeliveryChannelValue",
		"delivery_channel_target_arn":      "UpdateDeliveryChannelValue",
		"delivery_channel_type":            "UpdateDeliveryChannelValue",
		"delivery_channel_condition":       "UpdateDeliveryChannelValue",
		"delivery_channel_name":            "UpdateDeliveryChannelValue",
		"description":                      "UpdateDeliveryChannelValue",
		"status":                           0,
	}
	diff, err := newInstanceDiff("alicloud_config_delivery_channel", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_delivery_channel"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDeliveryChannels Response
		"DeliveryChannels": []interface{}{
			map[string]interface{}{
				"DeliveryChannelAssumeRoleArn": "UpdateDeliveryChannelValue",
				"DeliveryChannelCondition":     "UpdateDeliveryChannelValue",
				"DeliveryChannelName":          "UpdateDeliveryChannelValue",
				"DeliveryChannelTargetArn":     "UpdateDeliveryChannelValue",
				"DeliveryChannelType":          "UpdateDeliveryChannelValue",
				"Description":                  "UpdateDeliveryChannelValue",
				"Status":                       0,
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "DeliveryChannelSlsUnreachableError", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "PutDeliveryChannel" {
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
		err := resourceAlicloudConfigDeliveryChannelUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_delivery_channel"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribeDeliveryChannels" {
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
		err := resourceAlicloudConfigDeliveryChannelRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudConfigDeliveryChannelDelete(dExisted, rawClient)
	assert.Nil(t, err)
}
