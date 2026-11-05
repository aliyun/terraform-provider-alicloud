package nlb

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_nlb_load_balancer",
		&resource.Sweeper{
			Name: "alicloud_nlb_load_balancer",
			F:    testSweepNlbLoadBalancer,
		})
}

func testSweepNlbLoadBalancer(region string) error {
	rawClient, err := nlbTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-example",
		"tf_example",
		"terraform-example",
	}
	action := "ListLoadBalancers"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["MaxResults"] = helper.PageSizeLarge

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = aliyunClient.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.LoadBalancers", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.LoadBalancers", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			lbName := ""
			if item["LoadBalancerName"] != nil {
				lbName = fmt.Sprint(item["LoadBalancerName"])
			}

			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(lbName), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Nlb Load Balancer: %s", lbName)
					continue
				}
			}

			action = "UpdateLoadBalancerProtection"
			if err != nil {
				return errs.WrapError(err)
			}
			request = make(map[string]interface{})
			request["LoadBalancerId"] = item["LoadBalancerId"]
			request["RegionId"] = aliyunClient.RegionId
			request["ClientToken"] = helper.BuildClientToken(action)
			request["DeletionProtectionEnabled"] = "false"
			_, err = aliyunClient.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
			if err != nil {
				log.Printf("[ERROR] Failed to UpdateLoadBalancerProtection Nlb Load Balancer (%s): %s", lbName, err)
			}
			log.Printf("[INFO] UpdateLoadBalancerProtection Nlb Load Balancer success: %s ", lbName)

			action = "DeleteLoadBalancer"
			request := map[string]interface{}{
				"LoadBalancerId": item["LoadBalancerId"],
				"RegionId":       aliyunClient.RegionId,
			}
			_, err = aliyunClient.RpcPost("Nlb", "2022-04-30", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Nlb Load Balancer (%s): %s", lbName, err)
			}
			log.Printf("[INFO] Delete Nlb Load Balancer success: %s ", lbName)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

// lintignore: R001
func TestUnitAlicloudNlbLoadBalancer(t *testing.T) {
	dInit, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"load_balancer_name": "CreateLoadBalancerValue",
		"resource_group_id":  "CreateLoadBalancerValue",
		"load_balancer_type": "CreateLoadBalancerValue",
		"address_type":       "CreateLoadBalancerValue",
		"address_ip_version": "CreateLoadBalancerValue",
		"tags": map[string]string{
			"Created": "CreateLoadBalancerValue",
			"For":     "CreateLoadBalancerValue",
		},
		"vpc_id": "CreateLoadBalancerValue",
		"zone_mappings": []map[string]interface{}{
			{
				"vswitch_id":           "CreateLoadBalancerValue1",
				"zone_id":              "CreateLoadBalancerValue1",
				"allocation_id":        "CreateLoadBalancerValue1",
				"private_ipv4_address": "CreateLoadBalancerValue1",
			},
			{
				"vswitch_id":           "CreateLoadBalancerValue2",
				"zone_id":              "CreateLoadBalancerValue2",
				"allocation_id":        "CreateLoadBalancerValue2",
				"private_ipv4_address": "CreateLoadBalancerValue2",
			},
		},
		"cross_zone_enabled":   false,
		"bandwidth_package_id": "CreateLoadBalancerValue",
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
	rawClient, err := nlbTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// GetLoadBalancerAttribute
		"LoadBalancerName":           "CreateLoadBalancerValue",
		"CrossZoneEnabled":           false,
		"AddressIpVersion":           "CreateLoadBalancerValue",
		"AddressType":                "CreateLoadBalancerValue",
		"BandwidthPackageId":         "CreateLoadBalancerValue",
		"CreateTime":                 "DefaultValue",
		"DNSName":                    "DefaultValue",
		"Ipv6AddressType":            "DefaultValue",
		"LoadBalancerBusinessStatus": "DefaultValue",
		"LoadBalancerId":             "CreateLoadBalancerValue",
		"LoadBalancerType":           "CreateLoadBalancerValue",
		"ResourceGroupId":            "CreateLoadBalancerValue",
		"LoadBalancerStatus":         "Active",
		"VpcId":                      "CreateLoadBalancerValue",
		"ZoneMappings": []interface{}{
			map[string]interface{}{
				"EipType": "Common",
				"LoadBalancerAddresses": []interface{}{
					map[string]interface{}{
						"AllocationId":       "CreateLoadBalancerValue1",
						"PrivateIPv4Address": "CreateLoadBalancerValue1",
					},
				},
				"VSwitchId": "CreateLoadBalancerValue1",
				"ZoneId":    "CreateLoadBalancerValue1",
			},
			map[string]interface{}{
				"EipType": "Common",
				"LoadBalancerAddresses": []interface{}{
					map[string]interface{}{
						"AllocationId":       "CreateLoadBalancerValue2",
						"PrivateIPv4Address": "CreateLoadBalancerValue2",
					},
				},
				"VSwitchId": "CreateLoadBalancerValue2",
				"ZoneId":    "CreateLoadBalancerValue2",
			},
		},
		// ListTagResources
		"TagResources": []interface{}{
			map[string]interface{}{
				"TagKey":   "Created",
				"TagValue": "CreateLoadBalancerValue",
			},
			map[string]interface{}{
				"TagKey":   "For",
				"TagValue": "CreateLoadBalancerValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateLoadBalancer
		"LoadbalancerId": "CreateLoadBalancerValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_nlb_load_balancer", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceLoadBalancerCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateLoadBalancer" {
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
		err := resourceLoadBalancerCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceLoadBalancerUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateLoadBalancerAttribute
	attributesDiff := map[string]interface{}{
		"cross_zone_enabled": true,
		"load_balancer_name": "UpdateLoadBalancerAttributeValue",
	}
	diff, err := newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceLoadBalancer().Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetLoadBalancerAttribute Response
		"CrossZoneEnabled": true,
		"LoadBalancerName": "UpdateLoadBalancerAttributeValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateLoadBalancerAttribute" {
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
		err := resourceLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpdateLoadBalancerZones
	attributesDiff = map[string]interface{}{
		"zone_mappings": []interface{}{
			map[string]interface{}{
				"allocation_id":        "UpdateLoadBalancerZonesValue1",
				"private_ipv4_address": "UpdateLoadBalancerZonesValue1",
				"vswitch_id":           "UpdateLoadBalancerZonesValue1",
				"zone_id":              "UpdateLoadBalancerZonesValue1",
			},
			map[string]interface{}{
				"allocation_id":        "UpdateLoadBalancerZonesValue2",
				"private_ipv4_address": "UpdateLoadBalancerZonesValue2",
				"vswitch_id":           "UpdateLoadBalancerZonesValue2",
				"zone_id":              "UpdateLoadBalancerZonesValue2",
			},
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetLoadBalancerAttribute Response
		"ZoneMappings": []interface{}{
			map[string]interface{}{
				"EipType": "Common",
				"LoadBalancerAddresses": []interface{}{
					map[string]interface{}{
						"AllocationId":       "UpdateLoadBalancerZonesValue1",
						"PrivateIPv4Address": "UpdateLoadBalancerZonesValue1",
					},
				},
				"VSwitchId": "UpdateLoadBalancerZonesValue1",
				"ZoneId":    "UpdateLoadBalancerZonesValue1",
			},
			map[string]interface{}{
				"EipType": "Common",
				"LoadBalancerAddresses": []interface{}{
					map[string]interface{}{
						"AllocationId":       "UpdateLoadBalancerZonesValue2",
						"PrivateIPv4Address": "UpdateLoadBalancerZonesValue2",
					},
				},
				"VSwitchId": "UpdateLoadBalancerZonesValue2",
				"ZoneId":    "UpdateLoadBalancerZonesValue2",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateLoadBalancerZones" {
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
		err := resourceLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpdateLoadBalancerAddressTypeConfig
	attributesDiff = map[string]interface{}{
		"address_type": "UpdateLoadBalancerValue",
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// UpdateLoadBalancerAddressTypeConfig Response
		"AddressType": "UpdateLoadBalancerValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateLoadBalancerAddressTypeConfig" {
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
		err := resourceLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// AttachCommonBandwidthPackageToLoadBalancer
	attributesDiff = map[string]interface{}{
		"bandwidth_package_id": "UpdateLoadBalancerValue",
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// AttachCommonBandwidthPackageToLoadBalancer Response
		"BandwidthPackageId": "UpdateLoadBalancerValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AttachCommonBandwidthPackageToLoadBalancer" {
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
		err := resourceLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// DetachCommonBandwidthPackageFromLoadBalancer
	attributesDiff = map[string]interface{}{
		"bandwidth_package_id": "UpdateLoadBalancerValue1",
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DetachCommonBandwidthPackageFromLoadBalancer Response
		"BandwidthPackageId": "UpdateLoadBalancerValue1",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DetachCommonBandwidthPackageFromLoadBalancer" {
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
		err := resourceLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// TagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"TagResourcesValue_1": "TagResourcesValue_1",
			"TagResourcesValue_2": "TagResourcesValue_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListTagResources Response
		"TagResources": []interface{}{
			map[string]interface{}{
				"TagKey":   "TagResourcesValue_1",
				"TagValue": "TagResourcesValue_1",
			},
			map[string]interface{}{
				"TagKey":   "TagResourcesValue_2",
				"TagValue": "TagResourcesValue_2",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "TagResources" {
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
		err := resourceLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UntagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"UntagResourcesValue3_1": "UnTagResourcesValue3_1",
			"UntagResourcesValue3_2": "UnTagResourcesValue3_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListTagResources Response
		"TagResources": []interface{}{
			map[string]interface{}{
				"TagKey":   "UntagResourcesValue3_1",
				"TagValue": "UnTagResourcesValue3_1",
			},
			map[string]interface{}{
				"TagKey":   "UntagResourcesValue3_2",
				"TagValue": "UnTagResourcesValue3_2",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UntagResources" {
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
		err := resourceLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceLoadBalancer().Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceLoadBalancer().Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetLoadBalancerAttribute" {
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
		err := resourceLoadBalancerRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceLoadBalancerDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceLoadBalancer().Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteLoadBalancer" {
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
			if *action == "GetLoadBalancerAttribute" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceLoadBalancerDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test Nlb LoadBalancer. >>> Resource test cases, automatically generated.
// Case 3678
