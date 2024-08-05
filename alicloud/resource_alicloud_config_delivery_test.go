package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudConfigDelivery_OSS(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_delivery.default"
	checkoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	ra := resourceAttrInit(resourceId, AlicloudConfigDeliveryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDelivery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigdelivery%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigDeliveryBasicDependenceOSS)
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
					"delivery_channel_name":                  "${var.name}",
					"delivery_channel_type":                  "OSS",
					"delivery_channel_target_arn":            "${local.bucket}",
					"description":                            "${var.name}",
					"configuration_snapshot":                 "true",
					"configuration_item_change_notification": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":                  name,
						"delivery_channel_type":                  "OSS",
						"delivery_channel_target_arn":            CHECKSET,
						"description":                            name,
						"configuration_snapshot":                 "true",
						"configuration_item_change_notification": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
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
					"configuration_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration_snapshot":                 "true",
					"configuration_item_change_notification": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":                 "true",
						"configuration_item_change_notification": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name":                  "${var.name}",
					"delivery_channel_target_arn":            "${local.bucket}",
					"description":                            "${var.name}",
					"configuration_snapshot":                 "true",
					"configuration_item_change_notification": "true",
					"status":                                 "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":                  name,
						"delivery_channel_target_arn":            CHECKSET,
						"description":                            name,
						"configuration_snapshot":                 "true",
						"configuration_item_change_notification": "true",
						"status":                                 "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudConfigDelivery_MNS(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_delivery.default"
	checkoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	ra := resourceAttrInit(resourceId, AlicloudConfigDeliveryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDelivery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigdelivery%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigDeliveryBasicDependenceMNS)
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
					"delivery_channel_name":                  "${var.name}",
					"delivery_channel_type":                  "MNS",
					"delivery_channel_target_arn":            "${local.mns}",
					"description":                            "${var.name}",
					"configuration_item_change_notification": "true",
					"non_compliant_notification":             "true",
					"delivery_channel_condition":             configDeliveryChannelCondition,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":                  name,
						"delivery_channel_type":                  "MNS",
						"delivery_channel_target_arn":            CHECKSET,
						"description":                            name,
						"configuration_item_change_notification": "true",
						"non_compliant_notification":             "true",
						"delivery_channel_condition":             strings.Replace(strings.Replace(configDeliveryChannelCondition, `\n`, "\n", -1), `\"`, "\"", -1),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
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
					"non_compliant_notification": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"non_compliant_notification": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"non_compliant_notification":             "true",
					"configuration_item_change_notification": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"non_compliant_notification":             "true",
						"configuration_item_change_notification": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name":                  "${var.name}",
					"delivery_channel_target_arn":            "${local.mns}",
					"description":                            "${var.name}",
					"configuration_item_change_notification": "true",
					"non_compliant_notification":             "true",
					"status":                                 "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":                  name,
						"delivery_channel_target_arn":            CHECKSET,
						"description":                            name,
						"configuration_item_change_notification": "true",
						"non_compliant_notification":             "true",
						"status":                                 "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudConfigDelivery_SLS(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_delivery.default"
	checkoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	ra := resourceAttrInit(resourceId, AlicloudConfigDeliveryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDelivery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigdelivery%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigDeliveryBasicDependenceSLS)
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
					"delivery_channel_name":                  "${var.name}",
					"delivery_channel_type":                  "SLS",
					"delivery_channel_target_arn":            "${local.sls}",
					"description":                            "${var.name}",
					"non_compliant_notification":             "true",
					"configuration_item_change_notification": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":                  name,
						"delivery_channel_type":                  "SLS",
						"delivery_channel_target_arn":            CHECKSET,
						"description":                            name,
						"non_compliant_notification":             "true",
						"configuration_item_change_notification": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
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
					"non_compliant_notification": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"non_compliant_notification": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"non_compliant_notification":             "true",
					"configuration_item_change_notification": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"non_compliant_notification":             "true",
						"configuration_item_change_notification": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudConfigDeliveryMap0 = map[string]string{
	"configuration_item_change_notification": CHECKSET,
}

// Because the bucket cannot be deleted after being used by the delivery channel.
// Use pre-created Oss bucket in this test.
func AlicloudConfigDeliveryBasicDependenceOSS(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_account" "this" {}
locals {
  uid          	   = data.alicloud_account.this.id
  bucket	       = format("acs:oss:cn-shanghai:%%s:tf-test-bucket-for-config",local.uid)
  bucket_change	   = format("acs:oss:cn-shanghai:%%s:tf-test-bucket-for-config-update",local.uid)
}
`, name)
}

func AlicloudConfigDeliveryBasicDependenceMNS(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_account" "this" {}
locals {
  uid          = data.alicloud_account.this.id
  mns	       = format("acs:mns:%[2]s:%%s:/topics/%%s",local.uid,alicloud_mns_topic.default.name)
  mns_change   = format("acs:mns:%[2]s:%%s:/topics/%%s",local.uid,alicloud_mns_topic.change.name)
}
resource "alicloud_mns_topic" "default" {
  name = var.name
}
resource "alicloud_mns_topic" "change" {
  name = format("%%s-change",var.name)
}

`, name, defaultRegionToTest)
}

func AlicloudConfigDeliveryBasicDependenceSLS(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_account" "this" {}
locals {
  uid          	   = data.alicloud_account.this.id
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
`, name, defaultRegionToTest)
}

const configDeliveryChannelCondition = `[\n{\n\"filterType\":\"ResourceType\",\n\"values\":[\n\"ACS::CEN::CenInstance\",\n],\n\"multiple\":true\n}\n]\n`

func TestUnitAlicloudConfigDelivery(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_delivery"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_delivery"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"delivery_channel_name":                  "CreateConfigDeliveryValue",
		"delivery_channel_type":                  "CreateConfigDeliveryValue",
		"delivery_channel_target_arn":            "CreateConfigDeliveryValue",
		"description":                            "CreateConfigDeliveryValue",
		"delivery_channel_condition":             "CreateConfigDeliveryValue",
		"oversized_data_oss_target_arn":          "CreateConfigDeliveryValue",
		"configuration_snapshot":                 true,
		"configuration_item_change_notification": true,
		"non_compliant_notification":             true,
		"status":                                 1,
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
		"DeliveryChannel": map[string]interface{}{
			"Status":                              1,
			"OversizedDataOSSTargetArn":           "CreateConfigDeliveryValue",
			"ConfigurationSnapshot":               true,
			"Description":                         "CreateConfigDeliveryValue",
			"DeliveryChannelId":                   "CreateConfigDeliveryValue",
			"DeliveryChannelName":                 "CreateConfigDeliveryValue",
			"DeliveryChannelTargetArn":            "CreateConfigDeliveryValue",
			"DeliveryChannelAssumeRoleArn":        "CreateConfigDeliveryValue",
			"DeliveryChannelType":                 "CreateConfigDeliveryValue",
			"DeliveryChannelCondition":            "CreateConfigDeliveryValue",
			"NonCompliantNotification":            true,
			"ConfigurationItemChangeNotification": true,
		},
	}
	CreateMockResponse := map[string]interface{}{
		"DeliveryChannelId": "CreateConfigDeliveryValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_config_delivery", errorCode))
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
	err = resourceAliCloudConfigDeliveryCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateConfigDeliveryChannel" {
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
		err := resourceAliCloudConfigDeliveryCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_delivery"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudConfigDeliveryUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"delivery_channel_name":                  "UpdateConfigDeliveryValue",
		"delivery_channel_target_arn":            "UpdateConfigDeliveryValue",
		"description":                            "UpdateConfigDeliveryValue",
		"delivery_channel_condition":             "UpdateConfigDeliveryValue",
		"oversized_data_oss_target_arn":          "UpdateConfigDeliveryValue",
		"configuration_snapshot":                 false,
		"configuration_item_change_notification": false,
		"non_compliant_notification":             false,
		"status":                                 0,
	}
	diff, err := newInstanceDiff("alicloud_config_delivery", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_delivery"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"DeliveryChannel": map[string]interface{}{
			"Status":                              0,
			"OversizedDataOSSTargetArn":           "UpdateConfigDeliveryValue",
			"ConfigurationSnapshot":               false,
			"Description":                         "UpdateConfigDeliveryValue",
			"DeliveryChannelName":                 "UpdateConfigDeliveryValue",
			"DeliveryChannelTargetArn":            "UpdateConfigDeliveryValue",
			"DeliveryChannelAssumeRoleArn":        "UpdateConfigDeliveryValue",
			"DeliveryChannelCondition":            "UpdateConfigDeliveryValue",
			"NonCompliantNotification":            false,
			"ConfigurationItemChangeNotification": false,
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateConfigDeliveryChannel" {
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
		err := resourceAliCloudConfigDeliveryUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_delivery"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_config_delivery", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_delivery"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetConfigDeliveryChannel" {
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
		err := resourceAliCloudConfigDeliveryRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudConfigDeliveryDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_config_delivery", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_delivery"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteConfigDeliveryChannel" {
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
		err := resourceAliCloudConfigDeliveryDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test Config Delivery. >>> Resource test cases, automatically generated.
// Case 投递自动化测试-TF接入 7096
func TestAccAliCloudConfigDelivery_basic7096(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_delivery.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigDeliveryMap7096)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDelivery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigdelivery%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigDeliveryBasicDependence7096)
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
					"configuration_snapshot":                 "true",
					"description":                            "资源用例-创建",
					"delivery_channel_name":                  "资源用例-创建",
					"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
					"configuration_item_change_notification": "true",
					"delivery_channel_type":                  "SLS",
					"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
					"non_compliant_notification":             "false",
					"delivery_channel_condition":             "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"1\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ACK::Cluster\\\"],\\\"multiple\\\":true}]",
					"status":                                 "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_target_arn": "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
						"delivery_channel_type":       "SLS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "资源用例-创建",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "资源用例-创建",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": "资源用例-创建",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": "资源用例-创建",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration_item_change_notification": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_item_change_notification": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oversized_data_oss_target_arn": "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oversized_data_oss_target_arn": "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"non_compliant_notification": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"non_compliant_notification": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_condition": "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"1\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ACK::Cluster\\\"],\\\"multiple\\\":true}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_condition": "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"1\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ACK::Cluster\"],\"multiple\":true}]",
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
					"configuration_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "资源用例-编辑",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "资源用例-编辑",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": "资源用例-编辑",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": "资源用例-编辑",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_target_arn": "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/update-delivery-tf-test-log",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_target_arn": "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/update-delivery-tf-test-log",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration_item_change_notification": "false",
					"configuration_snapshot":                 "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_item_change_notification": "false",
						"configuration_snapshot":                 "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"oversized_data_oss_target_arn": "acs:oss:cn-shanghai:1511928242963727:update-delivery-tf-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oversized_data_oss_target_arn": "acs:oss:cn-shanghai:1511928242963727:update-delivery-tf-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"non_compliant_notification": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"non_compliant_notification": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_condition": "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"3\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ECS::Instance\\\"],\\\"multiple\\\":true}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_condition": "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"3\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ECS::Instance\"],\"multiple\":true}]",
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
					"configuration_snapshot":                 "true",
					"description":                            "资源用例-创建",
					"delivery_channel_name":                  "资源用例-创建",
					"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
					"configuration_item_change_notification": "true",
					"delivery_channel_type":                  "SLS",
					"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
					"non_compliant_notification":             "false",
					"delivery_channel_condition":             "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"1\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ACK::Cluster\\\"],\\\"multiple\\\":true}]",
					"status":                                 "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":                 "true",
						"description":                            "资源用例-创建",
						"delivery_channel_name":                  "资源用例-创建",
						"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
						"configuration_item_change_notification": "true",
						"delivery_channel_type":                  "SLS",
						"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
						"non_compliant_notification":             "false",
						"delivery_channel_condition":             "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"1\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ACK::Cluster\"],\"multiple\":true}]",
						"status":                                 "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudConfigDeliveryMap7096 = map[string]string{
	"status": CHECKSET,
}

func AlicloudConfigDeliveryBasicDependence7096(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 投递自动化测试-TF接入 7096  twin
func TestAccAliCloudConfigDelivery_basic7096_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_delivery.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigDeliveryMap7096)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDelivery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigdelivery%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigDeliveryBasicDependence7096)
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
					"configuration_snapshot":                 "true",
					"description":                            "资源用例-创建",
					"delivery_channel_name":                  "资源用例-创建",
					"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
					"configuration_item_change_notification": "true",
					"delivery_channel_type":                  "SLS",
					"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
					"non_compliant_notification":             "false",
					"delivery_channel_condition":             "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"1\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ACK::Cluster\\\"],\\\"multiple\\\":true}]",
					"status":                                 "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":                 "true",
						"description":                            "资源用例-创建",
						"delivery_channel_name":                  "资源用例-创建",
						"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
						"configuration_item_change_notification": "true",
						"delivery_channel_type":                  "SLS",
						"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
						"non_compliant_notification":             "false",
						"delivery_channel_condition":             "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"1\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ACK::Cluster\"],\"multiple\":true}]",
						"status":                                 "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 投递自动化测试-TF接入 7096  raw
func TestAccAliCloudConfigDelivery_basic7096_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_delivery.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigDeliveryMap7096)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigDelivery")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigdelivery%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigDeliveryBasicDependence7096)
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
					"configuration_snapshot":                 "true",
					"description":                            "资源用例-创建",
					"delivery_channel_name":                  "资源用例-创建",
					"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
					"configuration_item_change_notification": "true",
					"delivery_channel_type":                  "SLS",
					"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
					"non_compliant_notification":             "false",
					"delivery_channel_condition":             "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"1\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ACK::Cluster\\\"],\\\"multiple\\\":true}]",
					"status":                                 "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":                 "true",
						"description":                            "资源用例-创建",
						"delivery_channel_name":                  "资源用例-创建",
						"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
						"configuration_item_change_notification": "true",
						"delivery_channel_type":                  "SLS",
						"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
						"non_compliant_notification":             "false",
						"delivery_channel_condition":             "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"1\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ACK::Cluster\"],\"multiple\":true}]",
						"status":                                 "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration_snapshot":                 "false",
					"description":                            "资源用例-编辑",
					"delivery_channel_name":                  "资源用例-编辑",
					"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/update-delivery-tf-test-log",
					"configuration_item_change_notification": "false",
					"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:update-delivery-tf-test",
					"non_compliant_notification":             "true",
					"delivery_channel_condition":             "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"3\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ECS::Instance\\\"],\\\"multiple\\\":true}]",
					"status":                                 "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":                 "false",
						"description":                            "资源用例-编辑",
						"delivery_channel_name":                  "资源用例-编辑",
						"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/update-delivery-tf-test-log",
						"configuration_item_change_notification": "false",
						"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:update-delivery-tf-test",
						"non_compliant_notification":             "true",
						"delivery_channel_condition":             "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"3\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ECS::Instance\"],\"multiple\":true}]",
						"status":                                 "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Config Delivery. <<< Resource test cases, automatically generated.
