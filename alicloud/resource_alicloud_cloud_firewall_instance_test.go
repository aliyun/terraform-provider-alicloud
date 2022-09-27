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

func TestAccAlicloudCloudFirewallInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BssOpenApiService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "QueryAvailableInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddpinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{30})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":    "Subscription",
					"spec":            "premium_version",
					"ip_number":       "20",
					"band_width":      "10",
					"cfw_log":         "false",
					"cfw_log_storage": "1000",
					"cfw_service":     "false",
					"period":          "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":    "Subscription",
						"spec":            "premium_version",
						"ip_number":       "20",
						"band_width":      "10",
						"cfw_log":         "false",
						"cfw_log_storage": "1000",
						"cfw_service":     "false",
						"period":          "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_service": "true",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_service": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fw_vpc_number": "3",
					"modify_type":   "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fw_vpc_number": "3",
						"modify_type":   "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"band_width":  "20",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"band_width":  "20",
						"modify_type": "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log_storage": "1200",
					"modify_type":     "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log_storage": "1200",
						"modify_type":     "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cfw_log":     "true",
					"modify_type": "Upgrade",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cfw_log":     "true",
						"modify_type": "Upgrade",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_period": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_period": "6",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"band_width", "cfw_log", "cfw_log_storage", "cfw_service", "ip_number", "payment_type", "period", "modify_type", "spec"},
			},
		},
	})
}

var AlicloudCloudFirewallInstanceMap0 = map[string]string{}

func AlicloudCloudFirewallInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestUnitAlicloudCloudFirewallInstance(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cloud_firewall_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cloud_firewall_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"payment_type":    "CreateInstanceValue",
		"spec":            "CreateInstanceValue",
		"renewal_status":  "CreateInstanceValue",
		"ip_number":       20,
		"band_width":      10,
		"cfw_log":         false,
		"cfw_log_storage": 1000,
		"cfw_service":     false,
		"period":          6,
		"fw_vpc_number":   10,
		"instance_count":  10,
		"logistics":       "CreateInstanceValue",
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
		// QueryAvailableInstances
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{
				map[string]interface{}{
					"InstanceID":          "CreateInstanceValue",
					"CreateTime":          "CreateInstanceValue",
					"RenewStatus":         "CreateInstanceValue",
					"RenewalDurationUnit": "M",
					"Status":              "CreateInstanceValue",
					"SubscriptionType":    "CreateInstanceValue",
					"EndTime":             "CreateInstanceValue",
				},
			},
			"InstanceId": "CreateInstanceValue",
		},
		"Code": "Success",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateInstance
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cloud_firewall_instance", errorCode))
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
	err = resourceAlicloudCloudFirewallInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// QueryAvailableInstances Response
		"Data": map[string]interface{}{
			"InstanceId": "CreateInstanceValue",
		},
		"Code": "Success",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "NotApplicable", "nil"}
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
		err := resourceAlicloudCloudFirewallInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cloud_firewall_instance"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewBssopenapiClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudCloudFirewallInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// RenewInstance
	attributesDiff := map[string]interface{}{
		"renew_period": 1,
	}
	diff, err := newInstanceDiff("alicloud_cloud_firewall_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cloud_firewall_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// QueryAvailableInstances Response
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{
				map[string]interface{}{
					"RenewPeriod": 1,
				},
			},
		},
		"Code": "Success",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "NotApplicable", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RenewInstance" {
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
		err := resourceAlicloudCloudFirewallInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cloud_firewall_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyInstance
	attributesDiff = map[string]interface{}{
		"cfw_service":     true,
		"fw_vpc_number":   20,
		"ip_number":       30,
		"cfw_log_storage": 2000,
		"cfw_log":         true,
		"band_width":      20,
		"spec":            "enterprise_version",
		"instance_count":  20,
		"modify_type":     "ModifyInstanceValue",
	}
	diff, err = newInstanceDiff("alicloud_cloud_firewall_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cloud_firewall_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// QueryAvailableInstances Response
		"Data": map[string]interface{}{
			"InstanceList": []interface{}{
				map[string]interface{}{
					"CfwService":    true,
					"FwVpcNumber":   20,
					"IpNumber":      30,
					"CfwLogStorage": 2000,
					"CfwLog":        true,
					"BandWidth":     20,
					"Spec":          "enterprise_version",
					"InstanceCount": 20,
					"ModifyType":    "ModifyInstanceValue",
				},
			},
		},
		"Code": "Success",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "NotApplicable", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyInstance" {
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
		err := resourceAlicloudCloudFirewallInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cloud_firewall_instance"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "NotApplicable", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "QueryAvailableInstances" {
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
		err := resourceAlicloudCloudFirewallInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudCloudFirewallInstanceDelete(dExisted, rawClient)
	assert.Nil(t, err)

}
