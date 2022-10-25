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

func TestAccAlicloudAlidnsMonitorConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alidns_monitor_config.default"
	checkoutSupportedRegions(t, true, connectivity.AlidnsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudAlidnsMonitorConfigMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlidnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlidnsMonitorConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlidnsMonitorConfigBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_ICP_DOMAIN_NAME")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"addr_pool_id":        "${alicloud_alidns_address_pool.default.id}",
					"evaluation_count":    "1",
					"interval":            "60",
					"timeout":             "5000",
					"protocol_type":       "TCP",
					"monitor_extend_info": `{\"failureRate\":50,\"port\":80}`,
					"isp_city_node": []map[string]interface{}{
						{
							"city_code": "503",
							"isp_code":  "465",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"addr_pool_id":        CHECKSET,
						"protocol_type":       "TCP",
						"evaluation_count":    "1",
						"interval":            "60",
						"timeout":             "5000",
						"monitor_extend_info": CHECKSET,
						"isp_city_node.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"evaluation_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"evaluation_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"timeout": "10000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"timeout": "10000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":       "PING",
					"monitor_extend_info": `{\"packetNum\":20,\"packetLossRate\":10,\"failureRate\":50}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":       "PING",
						"monitor_extend_info": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"isp_city_node": []map[string]interface{}{
						{
							"city_code": "569",
							"isp_code":  "465",
						},
						{
							"city_code": "491",
							"isp_code":  "232",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"isp_city_node.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol_type":       "TCP",
					"monitor_extend_info": `{\"failureRate\":50,\"port\":80}`,
					"evaluation_count":    "1",
					"timeout":             "5000",
					"isp_city_node": []map[string]interface{}{
						{
							"city_code": "503",
							"isp_code":  "465",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol_type":       "TCP",
						"evaluation_count":    "1",
						"timeout":             "5000",
						"monitor_extend_info": CHECKSET,
						"isp_city_node.#":     "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"addr_pool_id"},
			},
		},
	})
}

var AlicloudAlidnsMonitorConfigMap0 = map[string]string{}

func AlicloudAlidnsMonitorConfigBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "domain_name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

resource "alicloud_alidns_gtm_instance" "default" {
  instance_name           = var.name
  payment_type            = "Subscription"
  period                  = 1
  renewal_status          = "ManualRenewal"
  package_edition         = "ultimate"
  health_check_task_count = 100
  sms_notification_count  = 1000
  public_cname_mode       = "SYSTEM_ASSIGN"
  ttl                     = 60
  cname_type              = "PUBLIC"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  alert_group             = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  public_user_domain_name = var.domain_name
  alert_config {
    sms_notice      = true
    notice_type     = "ADDR_ALERT"
    email_notice    = true
    dingtalk_notice = true
  }
}

resource "alicloud_alidns_address_pool" "default" {
  address_pool_name = var.name
  instance_id       = alicloud_alidns_gtm_instance.default.id
  lba_strategy      = "RATIO"
  type              = "IPV4"
  address {
    attribute_info = "{\"lineCodeRectifyType\":\"RECTIFIED\",\"lineCodes\":[\"os_namerica_us\"]}"
    remark         = "address_remark"
    address        = "1.1.1.1"
    mode           = "SMART"
    lba_weight     = 1
  }
}
`, name, os.Getenv("ALICLOUD_ICP_DOMAIN_NAME"))
}

func TestUnitAlicloudAlidnsMonitorConfig(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_alidns_monitor_config"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_alidns_monitor_config"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"addr_pool_id":        "AddDnsGtmMonitorValue",
		"evaluation_count":    1,
		"interval":            60,
		"timeout":             5000,
		"protocol_type":       "AddDnsGtmMonitorValue",
		"monitor_extend_info": "AddDnsGtmMonitorValue",
		"isp_city_node": []map[string]interface{}{
			{
				"city_code": "AddDnsGtmMonitorValue",
				"isp_code":  "AddDnsGtmMonitorValue",
			},
		},
		"lang": "AddDnsGtmMonitorValue",
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
		// DescribeDnsGtmMonitorConfig
		"EvaluationCount":   1,
		"Interval":          60,
		"Timeout":           5000,
		"MonitorExtendInfo": "AddDnsGtmMonitorValue",
		"ProtocolType":      "AddDnsGtmMonitorValue",
		"IspCityNodes": map[string]interface{}{
			"IspCityNode": []interface{}{
				map[string]interface{}{
					"CityCode": "AddDnsGtmMonitorValue",
					"IspCode":  "AddDnsGtmMonitorValue",
				},
			},
		},
		"MonitorConfigId": "AddDnsGtmMonitorValue",
	}
	CreateMockResponse := map[string]interface{}{
		// AddDnsGtmMonitor
		"MonitorConfigId": "AddDnsGtmMonitorValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_alidns_monitor_config", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidnsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudAlidnsMonitorConfigCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDnsGtmMonitorConfig Response
		"MonitorConfigId": "AddDnsGtmMonitorValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AddDnsGtmMonitor" {
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
		err := resourceAlicloudAlidnsMonitorConfigCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_alidns_monitor_config"].Schema).Data(dInit.State(), nil)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudAlidnsMonitorConfigUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	//UpdateDnsGtmMonitor
	attributesDiff := map[string]interface{}{
		"isp_city_node": []map[string]interface{}{
			{
				"city_code": "UpdateDnsGtmMonitorValue",
				"isp_code":  "UpdateDnsGtmMonitorValue",
			},
		},
		"monitor_extend_info": "UpdateDnsGtmMonitorValue",
		"protocol_type":       "UpdateDnsGtmMonitorValue",
		"evaluation_count":    2,
		"interval":            120,
		"timeout":             6000,
	}
	diff, err := newInstanceDiff("alicloud_alidns_monitor_config", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_alidns_monitor_config"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDnsGtmMonitorConfig
		"EvaluationCount":   2,
		"Interval":          120,
		"Timeout":           6000,
		"MonitorExtendInfo": "UpdateDnsGtmMonitorValue",
		"ProtocolType":      "UpdateDnsGtmMonitorValue",
		"IspCityNodes": map[string]interface{}{
			"IspCityNode": []interface{}{
				map[string]interface{}{
					"CityCode": "UpdateDnsGtmMonitorValue",
					"IspCode":  "UpdateDnsGtmMonitorValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateDnsGtmMonitor" {
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
		err := resourceAlicloudAlidnsMonitorConfigUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_alidns_monitor_config"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDnsGtmMonitorConfig" {
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
		err := resourceAlicloudAlidnsMonitorConfigRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidnsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudAlidnsMonitorConfigDelete(dExisted, rawClient)
	patches.Reset()
	assert.Nil(t, err)

}
