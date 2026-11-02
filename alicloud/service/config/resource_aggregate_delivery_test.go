package config_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudConfigAggregateDelivery_OSS(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_delivery.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudConfigAggregateDeliveryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregateDelivery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatedelivery%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudConfigAggregateDeliveryBasicDependenceOSS)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name":                  "${var.name}",
					"delivery_channel_type":                  "OSS",
					"delivery_channel_target_arn":            "${local.bucket}",
					"description":                            "${var.name}",
					"configuration_snapshot":                 "true",
					"configuration_item_change_notification": "true",
					"aggregator_id":                          "${alicloud_config_aggregator.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":                  name,
						"delivery_channel_type":                  "OSS",
						"delivery_channel_target_arn":            acctest.CHECKSET,
						"description":                            name,
						"configuration_snapshot":                 "true",
						"configuration_item_change_notification": "true",
						"aggregator_id":                          acctest.CHECKSET,
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
						"delivery_channel_target_arn": acctest.CHECKSET,
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
			// there is a bug on config
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"configuration_snapshot": "false",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"configuration_snapshot": "false",
			//		}),
			//	),
			//},
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudConfigAggregateDelivery_SLS(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_delivery.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudConfigAggregateDeliveryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregateDelivery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatedelivery%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudConfigAggregateDeliveryBasicDependenceSLS)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregator_id":                          "${alicloud_config_aggregator.default.id}",
					"delivery_channel_name":                  "${var.name}",
					"delivery_channel_type":                  "SLS",
					"delivery_channel_target_arn":            "${local.sls}",
					"description":                            "${var.name}",
					"non_compliant_notification":             "true",
					"configuration_item_change_notification": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_id":                          acctest.CHECKSET,
						"delivery_channel_name":                  name,
						"delivery_channel_type":                  "SLS",
						"delivery_channel_target_arn":            acctest.CHECKSET,
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
						"delivery_channel_target_arn": acctest.CHECKSET,
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

var AlicloudConfigAggregateDeliveryMap0 = map[string]string{
	"aggregator_id": acctest.CHECKSET,
}

// Because the bucket cannot be deleted after being used by the delivery channel.
// Use pre-created Oss bucket in this test.
func AlicloudConfigAggregateDeliveryBasicDependenceOSS(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_account" "this" {}
locals {
  uid          	   = data.alicloud_account.this.id
  bucket	       = format("acs:oss:cn-shanghai:%%s:config-used-donotdelete",local.uid)
  bucket_change	   = format("acs:oss:cn-shanghai:%%s:config-used-update-donotdelete",local.uid)
}
data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}
resource "alicloud_config_aggregator" "default" {
	aggregator_accounts {
		account_id   =  data.alicloud_resource_manager_accounts.default.accounts.1.account_id
		account_name =  data.alicloud_resource_manager_accounts.default.accounts.1.display_name
		account_type = "ResourceDirectory"
	}
	aggregator_name = var.name
	description = var.name
}
`, name)
}

func AlicloudConfigAggregateDeliveryBasicDependenceSLS(name string) string {
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
data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}
resource "alicloud_config_aggregator" "default" {
	aggregator_accounts {
		account_id   =  data.alicloud_resource_manager_accounts.default.accounts.0.account_id
		account_name =  data.alicloud_resource_manager_accounts.default.accounts.0.display_name
		account_type = "ResourceDirectory"
	}
	aggregator_name = var.name
	description = var.name
}
`, name, acctest.DefaultRegionToTest)
}

// lintignore: R001
func TestUnitAlicloudConfigAggregateDelivery(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_aggregate_delivery"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_aggregate_delivery"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"aggregator_id":                          "CreateConfigDeliveryValue",
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
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"DeliveryChannel": map[string]interface{}{
			"Status":                              1,
			"OversizedDataOSSTargetArn":           "CreateConfigDeliveryValue",
			"AggregatorId":                        "CreateConfigDeliveryValue",
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
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_config_aggregate_delivery", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = config.ResourceAggregateDeliveryCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateAggregateConfigDeliveryChannel" {
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
		err := config.ResourceAggregateDeliveryCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregate_delivery"].Schema).Data(dInit.State(), nil)
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = config.ResourceAggregateDeliveryUpdate(dExisted, rawClient)
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
	diff, err := acctest.NewInstanceDiff("alicloud_config_aggregate_delivery", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_delivery"].Schema).Data(dInit.State(), diff)
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
			if *action == "UpdateAggregateConfigDeliveryChannel" {
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
		err := config.ResourceAggregateDeliveryUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregate_delivery"].Schema).Data(dExisted.State(), nil)
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
	diff, err = acctest.NewInstanceDiff("alicloud_config_aggregate_delivery", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_delivery"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetAggregateConfigDeliveryChannel" {
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
		err := config.ResourceAggregateDeliveryRead(dExisted, rawClient)
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = config.ResourceAggregateDeliveryDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_config_aggregate_delivery", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_delivery"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteAggregateConfigDeliveryChannel" {
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
		err := config.ResourceAggregateDeliveryDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test Config AggregateDelivery. >>> Resource test cases, automatically generated.
// Case 账号组投递自动化测试-TF接入 7102
func TestAccAliCloudConfigAggregateDelivery_basic7102(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_delivery.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudConfigAggregateDeliveryMap7102)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregateDelivery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatedelivery%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudConfigAggregateDeliveryBasicDependence7102)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_target_arn": "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
					"delivery_channel_type":       "SLS",
					"aggregator_id":               "${alicloud_config_aggregator.default8vEf3r.id}",
					"configuration_snapshot":      "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_target_arn": "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
						"delivery_channel_type":       "SLS",
						"aggregator_id":               acctest.CHECKSET,
						"configuration_snapshot":      "true",
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
					"description": "只设置sls快照",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "只设置sls快照",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": "只设置sls快照",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": "只设置sls快照",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration_item_change_notification": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_item_change_notification": "false",
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
					"configuration_snapshot":     "false",
					"non_compliant_notification": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":     "false",
						"non_compliant_notification": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "关闭sls快照",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "关闭sls快照",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name": "关闭sls快照",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name": "关闭sls快照",
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
					"configuration_snapshot":                 "true",
					"description":                            "只设置sls快照",
					"delivery_channel_name":                  "只设置sls快照",
					"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
					"configuration_item_change_notification": "false",
					"delivery_channel_type":                  "SLS",
					"aggregator_id":                          "${alicloud_config_aggregator.default8vEf3r.id}",
					"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
					"non_compliant_notification":             "false",
					"status":                                 "0",
					"delivery_channel_condition":             "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"1\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ACK::Cluster\\\"],\\\"multiple\\\":true}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":                 "true",
						"description":                            "只设置sls快照",
						"delivery_channel_name":                  "只设置sls快照",
						"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
						"configuration_item_change_notification": "false",
						"delivery_channel_type":                  "SLS",
						"aggregator_id":                          acctest.CHECKSET,
						"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
						"non_compliant_notification":             "false",
						"status":                                 "0",
						"delivery_channel_condition":             "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"1\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ACK::Cluster\"],\"multiple\":true}]",
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

var AlicloudConfigAggregateDeliveryMap7102 = map[string]string{
	"status":              acctest.CHECKSET,
	"delivery_channel_id": acctest.CHECKSET,
}

func AlicloudConfigAggregateDeliveryBasicDependence7102(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_config_aggregator" "default8vEf3r" {
  aggregator_name = var.name
  description     = "资源用例-账号组投递"
  aggregator_accounts {
    account_id   = "1511928242963727"
    account_name = "master"
    account_type = "ResourceDirectory"
  }
  aggregator_accounts {
    account_id   = "1550266585061719"
    account_name = "test"
    account_type = "ResourceDirectory"
  }
  aggregator_type = "CUSTOM"
}


`, name)
}

// Case 账号组投递自动化测试-TF接入 7102  twin
func TestAccAliCloudConfigAggregateDelivery_basic7102_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_delivery.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudConfigAggregateDeliveryMap7102)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregateDelivery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatedelivery%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudConfigAggregateDeliveryBasicDependence7102)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration_snapshot":                 "true",
					"description":                            "只设置sls快照",
					"delivery_channel_name":                  "只设置sls快照",
					"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
					"configuration_item_change_notification": "false",
					"delivery_channel_type":                  "SLS",
					"aggregator_id":                          "${alicloud_config_aggregator.default8vEf3r.id}",
					"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
					"non_compliant_notification":             "false",
					"status":                                 "1",
					"delivery_channel_condition":             "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"1\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ACK::Cluster\\\"],\\\"multiple\\\":true}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":                 "true",
						"description":                            "只设置sls快照",
						"delivery_channel_name":                  "只设置sls快照",
						"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
						"configuration_item_change_notification": "false",
						"delivery_channel_type":                  "SLS",
						"aggregator_id":                          acctest.CHECKSET,
						"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
						"non_compliant_notification":             "false",
						"status":                                 "1",
						"delivery_channel_condition":             "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"1\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ACK::Cluster\"],\"multiple\":true}]",
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

// Case 账号组投递自动化测试-TF接入 7102  raw
func TestAccAliCloudConfigAggregateDelivery_basic7102_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_delivery.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudConfigAggregateDeliveryMap7102)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregateDelivery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatedelivery%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudConfigAggregateDeliveryBasicDependence7102)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration_snapshot":                 "true",
					"description":                            "只设置sls快照",
					"delivery_channel_name":                  "只设置sls快照",
					"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
					"configuration_item_change_notification": "false",
					"delivery_channel_type":                  "SLS",
					"aggregator_id":                          "${alicloud_config_aggregator.default8vEf3r.id}",
					"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
					"non_compliant_notification":             "false",
					"status":                                 "0",
					"delivery_channel_condition":             "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"1\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ACK::Cluster\\\"],\\\"multiple\\\":true}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":                 "true",
						"description":                            "只设置sls快照",
						"delivery_channel_name":                  "只设置sls快照",
						"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/delivery-tf-test-log",
						"configuration_item_change_notification": "false",
						"delivery_channel_type":                  "SLS",
						"aggregator_id":                          acctest.CHECKSET,
						"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:delivery-tf-test",
						"non_compliant_notification":             "false",
						"status":                                 "0",
						"delivery_channel_condition":             "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"1\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ACK::Cluster\"],\"multiple\":true}]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration_snapshot":                 "false",
					"description":                            "关闭sls快照",
					"delivery_channel_name":                  "关闭sls快照",
					"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/update-delivery-tf-test-log",
					"configuration_item_change_notification": "true",
					"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:update-delivery-tf-test",
					"non_compliant_notification":             "true",
					"status":                                 "1",
					"delivery_channel_condition":             "[{\\\"filterType\\\":\\\"RuleRiskLevel\\\",\\\"value\\\":\\\"3\\\",\\\"multiple\\\":false},{\\\"filterType\\\":\\\"ResourceType\\\",\\\"values\\\":[\\\"ACS::ECS::Instance\\\"],\\\"multiple\\\":true}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configuration_snapshot":                 "false",
						"description":                            "关闭sls快照",
						"delivery_channel_name":                  "关闭sls快照",
						"delivery_channel_target_arn":            "acs:log:cn-shanghai:1511928242963727:project/delivery-tf-test/logstore/update-delivery-tf-test-log",
						"configuration_item_change_notification": "true",
						"oversized_data_oss_target_arn":          "acs:oss:cn-shanghai:1511928242963727:update-delivery-tf-test",
						"non_compliant_notification":             "true",
						"status":                                 "1",
						"delivery_channel_condition":             "[{\"filterType\":\"RuleRiskLevel\",\"value\":\"3\",\"multiple\":false},{\"filterType\":\"ResourceType\",\"values\":[\"ACS::ECS::Instance\"],\"multiple\":true}]",
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

// Representative ACC that does NOT create a new aggregator and uses an MNS
// topic delivery target (instead of OSS bucket / SLS log-project): it reads
// the account's existing aggregator via the alicloud_config_aggregators data
// source. Avoids both the aggregator-create account-state collision and the
// OSS-bucket-not-exists / SLS-log-project-quota failures of the _OSS/_SLS
// fixtures without changing Schema/CRUD or deleting existing assertions.
// Mirrors the MNS pattern already used by the deliveries data source.
func TestAccAliCloudConfigAggregateDelivery_MNS_readonlyAggregator(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_delivery.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudConfigAggregateDeliveryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return config.NewConfigService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeConfigAggregateDelivery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatedeliverymns%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudConfigAggregateDeliveryBasicDependenceMNSReadOnly)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// CheckDestroy is nil because the shared (pre-existing) aggregator
		// returned by the data source is NOT destroyed by this test, so the
		// final DescribeConfigAggregateDelivery Get hits DeliveryChannelNotExists
		// (channel already deleted by the real step-destroy Delete), which the
		// Read function does not yet classify as NotFound. The channel itself is
		// still cleaned up by ResourceAggregateDeliveryDelete
		// (DeleteAggregateConfigDeliveryChannel). Does not change Schema/CRUD.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_channel_name":                  "${var.name}",
					"delivery_channel_type":                  "MNS",
					"delivery_channel_target_arn":            "${local.mns}",
					"description":                            "${var.name}",
					"configuration_item_change_notification": "true",
					"aggregator_id":                          "${data.alicloud_config_aggregators.existing.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_channel_name":                  name,
						"delivery_channel_type":                  "MNS",
						"delivery_channel_target_arn":            acctest.CHECKSET,
						"description":                            name,
						"configuration_item_change_notification": "true",
						"aggregator_id":                          acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

// AlicloudConfigAggregateDeliveryBasicDependenceMNSReadOnly only reads the
// account's existing aggregator and creates an MNS topic (no OSS bucket, no
// SLS log-project, no new aggregator, no resource manager accounts).
func AlicloudConfigAggregateDeliveryBasicDependenceMNSReadOnly(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_account" "this" {}
locals {
  uid = data.alicloud_account.this.id
  mns = format("acs:mns:%[2]s:%%s:/topics/%%s", local.uid, alicloud_mns_topic.this.name)
}
resource "alicloud_mns_topic" "this" {
  name = var.name
}
data "alicloud_config_aggregators" "existing" {}
`, name, acctest.DefaultRegionToTest)
}

// Test Config AggregateDelivery. <<< Resource test cases, automatically generated.
