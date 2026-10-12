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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudAlidnsGtmInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_gtm_instance.default"
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudAlidnsGtmInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsGtmInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%salidnsgtminstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsGtmInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":           name,
					"payment_type":            "Subscription",
					"period":                  "1",
					"renewal_status":          "ManualRenewal",
					"package_edition":         "ultimate",
					"health_check_task_count": "100",
					"sms_notification_count":  "1000",
					"public_cname_mode":       "SYSTEM_ASSIGN",
					"ttl":                     "60",
					"cname_type":              "PUBLIC",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"alert_group":             []string{"${alicloud_cms_alarm_contact_group.default.0.alarm_contact_group_name}"},
					"public_user_domain_name": "${var.domain_name}",
					"strategy_mode":           "GEO",
					"alert_config": []map[string]interface{}{
						{
							"sms_notice":      "true",
							"notice_type":     "ADDR_ALERT",
							"email_notice":    "true",
							"dingtalk_notice": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":           name,
						"cname_type":              "PUBLIC",
						"ttl":                     "60",
						"alert_group.#":           "1",
						"alert_config.#":          "1",
						"resource_group_id":       CHECKSET,
						"public_cname_mode":       "SYSTEM_ASSIGN",
						"strategy_mode":           "GEO",
						"public_user_domain_name": CHECKSET,
						"public_rr":               CHECKSET,
						"public_zone_name":        CHECKSET,
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
					"ttl": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ttl": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"strategy_mode": "LATENCY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"strategy_mode": "LATENCY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_group": []string{"${alicloud_cms_alarm_contact_group.default.0.alarm_contact_group_name}", "${alicloud_cms_alarm_contact_group.default.1.alarm_contact_group_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_group.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_config": []map[string]interface{}{
						{
							"sms_notice":      "true",
							"notice_type":     "ADDR_RESUME",
							"email_notice":    "true",
							"dingtalk_notice": "true",
						},
						{
							"sms_notice":      "true",
							"notice_type":     "ADDR_POOL_GROUP_UNAVAILABLE",
							"email_notice":    "true",
							"dingtalk_notice": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_config.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     name,
					"public_cname_mode": "SYSTEM_ASSIGN",
					"ttl":               "60",
					"alert_group":       []string{"${alicloud_cms_alarm_contact_group.default.0.alarm_contact_group_name}"},
					"strategy_mode":     "GEO",
					"alert_config": []map[string]interface{}{
						{
							"sms_notice":      "true",
							"notice_type":     "ADDR_ALERT",
							"email_notice":    "true",
							"dingtalk_notice": "true",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"ttl":               "60",
						"alert_config.#":    "1",
						"alert_group.#":     "1",
						"strategy_mode":     "GEO",
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang":              "zh",
					"force_update":      true,
					"renew_period":      1,
					"public_cname_mode": "CUSTOM",
					"public_rr":         "www",
					"public_zone_name":  "${var.domain_name}",
					"alert_config": []map[string]interface{}{
						{
							"sms_notice":      "true",
							"notice_type":     "ADDR_ALERT",
							"email_notice":    "false",
							"dingtalk_notice": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang":              "zh",
						"force_update":      "true",
						"public_cname_mode": "CUSTOM",
						"public_rr":         "www",
						"alert_config.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"public_cname_mode":       "SYSTEM_ASSIGN",
					"public_rr":               "api",
					"public_zone_name":        "test.${var.domain_name}",
					"public_user_domain_name": "test.${var.domain_name}",
					"alert_config": []map[string]interface{}{
						{
							"sms_notice":      "false",
							"notice_type":     "ADDR_ALERT",
							"email_notice":    "true",
							"dingtalk_notice": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_rr": "api",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// health_check_task_count/sms_notification_count/renewal_status are
				// ForceNew BSS order parameters now written back via
				// QueryAvailableInstances in Read, so they no longer need to be
				// ignored. lang/force_update/period are non-ForceNew attributes
				// without a readback path, so import verification skips them.
				ImportStateVerifyIgnore: []string{"lang", "force_update", "period"},
			},
		},
	})
}

var AlicloudAlidnsGtmInstanceMap0 = map[string]string{}

func AlicloudAlidnsGtmInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
variable "domain_name" {
  default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_cms_alarm_contact_group" "default" {
  count                    = 2
  alarm_contact_group_name = join("-", [var.name, count.index])
}
`, name, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"))
}

// lintignore: R001
func TestUnitAlicloudAlidnsGtmInstance(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"payment_type":            "CreateInstanceValue",
		"period":                  1,
		"renew_period":            1,
		"renewal_status":          "CreateInstanceValue",
		"package_edition":         "CreateInstanceValue",
		"health_check_task_count": 1,
		"sms_notification_count":  1,
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
		// DescribeAlidnsGtmInstance
		"ResourceGroupId": "CreateInstanceValue",
		"PaymentType":     "CreateInstanceValue",
		"VersionCode":     "CreateInstanceValue",
		"Config": map[string]interface{}{
			"CnameType":            "CreateInstanceValue",
			"InstanceName":         "CreateInstanceValue",
			"StrategyMode":         "CreateInstanceValue",
			"PublicCnameMode":      "CreateInstanceValue",
			"PublicRr":             "CreateInstanceValue",
			"PublicUserDomainName": "CreateInstanceValue",
			"PubicZoneName":        "CreateInstanceValue",
			"Ttl":                  1,
			"AlertGroup":           "[\"1\"]",
			"AlertConfig": map[string]interface{}{
				"AlertConfig": []interface{}{
					map[string]interface{}{
						"SmsNotice":      "CreateInstanceValue",
						"NoticeType":     "CreateInstanceValue",
						"EmailNotice":    "CreateInstanceValue",
						"DingtalkNotice": "CreateInstanceValue",
					},
				},
			},
		},
		"Data": map[string]interface{}{
			"InstanceId": "CreateInstanceValue",
		},
		"Code": "Success",
	}
	CreateMockResponse := map[string]interface{}{
		"Data": map[string]interface{}{
			"InstanceId": "CreateInstanceValue",
		},
		"Code": "Success",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_alidns_gtm_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBssopenapiClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudAlidnsGtmInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		"Data": map[string]interface{}{
			"InstanceId": "CreateInstanceValue",
		},
		"Code": "Success",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateInstance" {
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
		err := resourceAlicloudAlidnsGtmInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidnsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudAlidnsGtmInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	//MoveGtmResourceGroup
	attributesDiff := map[string]interface{}{
		"resource_group_id": "MoveGtmResourceGroupValue",
		"lang":              "MoveGtmResourceGroupValue",
	}
	diff, err := newInstanceDiff("alicloud_alidns_gtm_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDnsGtmInstance Response
		"ResourceGroupId": "MoveGtmResourceGroupValue",
		"Lang":            "MoveGtmResourceGroupValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "MoveGtmResourceGroup" {
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
		err := resourceAlicloudAlidnsGtmInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	//SwitchDnsGtmInstanceStrategyMode
	attributesDiff = map[string]interface{}{
		"strategy_mode": "MoveGtmResourceGroupValue",
		"lang":          "MoveGtmResourceGroupValue",
	}
	diff, err = newInstanceDiff("alicloud_alidns_gtm_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDnsGtmInstance Response
		"Config": map[string]interface{}{
			"StrategyMode": "MoveGtmResourceGroupValue",
		},
		"Lang": "MoveGtmResourceGroupValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "SwitchDnsGtmInstanceStrategyMode" {
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
		err := resourceAlicloudAlidnsGtmInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	//UpdateDnsGtmInstanceGlobalConfig
	attributesDiff = map[string]interface{}{
		"alert_config": []map[string]interface{}{
			{
				"sms_notice":      true,
				"notice_type":     "ADDR_RESUME",
				"email_notice":    true,
				"dingtalk_notice": true,
			},
		},
		"alert_group":             []string{"2"},
		"instance_name":           "MoveGtmResourceGroupValue",
		"ttl":                     2,
		"public_cname_mode":       "MoveGtmResourceGroupValue",
		"public_rr":               "MoveGtmResourceGroupValue",
		"public_user_domain_name": "MoveGtmResourceGroupValue",
		"public_zone_name":        "MoveGtmResourceGroupValue",
		"cname_type":              "MoveGtmResourceGroupValue",
		"force_update":            true,
		"lang":                    "MoveGtmResourceGroupValue",
	}
	diff, err = newInstanceDiff("alicloud_alidns_gtm_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Config": map[string]interface{}{
			"CnameType":            "MoveGtmResourceGroupValue",
			"InstanceName":         "MoveGtmResourceGroupValue",
			"PublicCnameMode":      "MoveGtmResourceGroupValue",
			"PublicRr":             "MoveGtmResourceGroupValue",
			"PublicUserDomainName": "MoveGtmResourceGroupValue",
			"PubicZoneName":        "MoveGtmResourceGroupValue",
			"Ttl":                  2,
			"AlertGroup":           "[\"2\"]",
			"AlertConfig": map[string]interface{}{
				"AlertConfig": []interface{}{
					map[string]interface{}{
						"SmsNotice":      true,
						"NoticeType":     "ADDR_RESUME",
						"EmailNotice":    true,
						"DingtalkNotice": true,
					},
				},
			},
		},
		"Lang":        "MoveGtmResourceGroupValue",
		"ForceUpdate": true,
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateDnsGtmInstanceGlobalConfig" {
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
		err := resourceAlicloudAlidnsGtmInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDnsGtmInstance" {
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
		err := resourceAlicloudAlidnsGtmInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	err = resourceAlicloudAlidnsGtmInstanceDelete(dExisted, rawClient)
	assert.Nil(t, err)
}

// TestUnitAlicloudAlidnsGtmInstanceStrategyModeLang asserts that when
// strategy_mode is updated and lang is set, the Lang attribute is forwarded
// onto the SwitchDnsGtmInstanceStrategyMode request body. Regression guard for
// a bug that wrote Lang onto the stale MoveGtmResourceGroup request map instead,
// so the RPC was invoked without Lang and the language selection was dropped.
func TestUnitAlicloudAlidnsGtmInstanceStrategyModeLang(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"payment_type":    "Subscription",
		"period":          1,
		"package_edition": "ultimate",
		"strategy_mode":   "GEO",
	}
	// Set attributes with string-literal keys (tfproviderlint R001 forbids variable keys).
	assert.Nil(t, dInit.Set("payment_type", "Subscription"))
	assert.Nil(t, dInit.Set("period", 1))
	assert.Nil(t, dInit.Set("package_edition", "ultimate"))
	assert.Nil(t, dInit.Set("strategy_mode", "GEO"))
	// State() returns nil without an ID, which would make newInstanceDiff panic;
	// give the existing resource an ID so the diff/state machinery is well-formed.
	dInit.SetId("gtm-test-instance-id")

	attributesDiff := map[string]interface{}{
		"strategy_mode": "LATENCY",
		"lang":          "zh",
	}
	diff, err := newInstanceDiff("alicloud_alidns_gtm_instance", attributes, attributesDiff, dInit.State())
	assert.Nil(t, err)
	dExisted, _ := schema.InternalMap(p["alicloud_alidns_gtm_instance"].Schema).Data(dInit.State(), diff)

	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	aliClient := rawClient.(*connectivity.AliyunClient)

	var strategyModeBody map[string]interface{}
	var strategyModeCalled bool
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "RpcPost", func(_ *connectivity.AliyunClient, apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
		if apiName == "SwitchDnsGtmInstanceStrategyMode" {
			strategyModeCalled = true
			strategyModeBody = body
		}
		return map[string]interface{}{"Code": "Success"}, nil
	})
	defer patches.Reset()

	err = resourceAlicloudAlidnsGtmInstanceUpdate(dExisted, aliClient)
	assert.Nil(t, err)

	assert.True(t, strategyModeCalled, "strategy_mode change must invoke SwitchDnsGtmInstanceStrategyMode")
	assert.NotNil(t, strategyModeBody, "SwitchDnsGtmInstanceStrategyMode request body must be captured")
	langValue, langOk := strategyModeBody["Lang"]
	assert.True(t, langOk, "SwitchDnsGtmInstanceStrategyMode request must carry Lang when lang is set")
	assert.Equal(t, "zh", langValue, "Lang value must match the configured lang")
}
