package ecs

// White-box unit test for alicloud_ecs_dedicated_host, split out of the external
// ecs_test file. This test calls unexported CRUD funcs (resourceEcsDedicatedHost
// Create/Read/Update/Delete), newInstanceDiff and the mapMerge/String aliases, so it
// must live in the internal package ecs test binary (mirroring disk_unit_test.go).
// Provider().ResourcesMap / sharedClientForRegion are replaced with the
// package-internal unitTestProviderResources / unitTestSharedClientForRegion
// copies to avoid importing acctest (import cycle: alicloud -> service/ecs).

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// lintignore: R001
func TestUnitAlicloudECSDedicatedHost(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"dedicated_host_type":   "AllocateDedicatedHostsValue",
		"description":           "AllocateDedicatedHostsValue",
		"dedicated_host_name":   "AllocateDedicatedHostsValue",
		"action_on_maintenance": "AllocateDedicatedHostsValue",
		"auto_placement":        "AllocateDedicatedHostsValue",
		"min_quantity":          1,
		"network_attributes": []map[string]interface{}{
			{
				"udp_timeout":     70,
				"slb_udp_timeout": 70,
			},
		},
		"resource_group_id": "AllocateDedicatedHostsValue",
		"tags": map[string]string{
			"Created": "TF",
		},
		"zone_id":                   "AllocateDedicatedHostsValue",
		"auto_release_time":         "AllocateDedicatedHostsValue",
		"auto_renew":                false,
		"auto_renew_period":         1,
		"cpu_over_commit_ratio":     1.2,
		"dedicated_host_cluster_id": "AllocateDedicatedHostsValue",
		"expired_time":              "AllocateDedicatedHostsValue",
		"payment_type":              "AllocateDedicatedHostsValue",
		"sale_cycle":                "AllocateDedicatedHostsValue",
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
		// DescribeDedicatedHosts
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"DedicatedHostId":     "AllocateDedicatedHostsValue",
					"ActionOnMaintenance": "AllocateDedicatedHostsValue",
					"AutoPlacement":       "AllocateDedicatedHostsValue",
					"AutoReleaseTime":     "AllocateDedicatedHostsValue",
					"CpuOverCommitRatio":  "AllocateDedicatedHostsValue",
					"DedicatedHostName":   "AllocateDedicatedHostsValue",
					"DedicatedHostType":   "AllocateDedicatedHostsValue",
					"Description":         "AllocateDedicatedHostsValue",
					"NetworkAttributes": map[string]interface{}{
						"SlbUdpTimeout": 70,
						"UdpTimeout":    70,
					},
					"ChargeType":      "AllocateDedicatedHostsValue",
					"ResourceGroupId": "AllocateDedicatedHostsValue",
					"SaleCycle":       "AllocateDedicatedHostsValue",
					"Status":          "Available",
					"ZoneId":          "AllocateDedicatedHostsValue",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "Created",
								"Value": "TF",
							},
						},
					},
				},
			},
		},
		"DedicatedHostIdSets": map[string]interface{}{
			"DedicatedHostId": []interface{}{
				"AllocateDedicatedHostsValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// AllocateDedicatedHosts
		"DedicatedHostIdSets": map[string]interface{}{
			"DedicatedHostId": []interface{}{
				"AllocateDedicatedHostsValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecs_dedicated_host", errorCode))
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
	err = resourceEcsDedicatedHostCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHostIdSets": map[string]interface{}{
			"DedicatedHostId": []interface{}{
				"AllocateDedicatedHostsValue",
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AllocateDedicatedHosts" {
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
		err := resourceEcsDedicatedHostCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceEcsDedicatedHostUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDedicatedHostAutoReleaseTime
	attributesDiff := map[string]interface{}{
		"auto_release_time": "ModifyDedicatedHostAutoReleaseTimeValue",
	}
	diff, err := newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"AutoReleaseTime": "ModifyDedicatedHostAutoReleaseTimeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDedicatedHostAutoReleaseTime" {
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
		err := resourceEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// JoinResourceGroup
	attributesDiff = map[string]interface{}{
		"resource_group_id": "JoinResourceGroupValue",
	}
	diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"ResourceGroupId": "JoinResourceGroupValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "JoinResourceGroup" {
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
		err := resourceEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// RenewDedicatedHosts
	attributesDiff = map[string]interface{}{
		"expired_time": "RenewDedicatedHostsValue",
		"sale_cycle":   "RenewDedicatedHostsValue",
	}
	diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"SaleCycle": "RenewDedicatedHostsValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RenewDedicatedHosts" {
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
		err := resourceEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// ModifyDedicatedHostsChargeType
	attributesDiff = map[string]interface{}{
		"expired_time": "ModifyDedicatedHostsChargeTypeValue",
		"payment_type": "ModifyDedicatedHostsChargeTypeValue",
		"sale_cycle":   "ModifyDedicatedHostsChargeTypeValue",
		"detail_fee":   false,
		"dry_run":      false,
	}
	diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"ChargeType": "ModifyDedicatedHostsChargeTypeValue",
					"SaleCycle":  "ModifyDedicatedHostsChargeTypeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDedicatedHostsChargeType" {
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
		err := resourceEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// ModifyDedicatedHostAttribute
	attributesDiff = map[string]interface{}{
		"action_on_maintenance": "ModifyDedicatedHostAttributeValue",
		"auto_placement":        "ModifyDedicatedHostAttributeValue",
		"cpu_over_commit_ratio": 1.6,
		"dedicated_host_name":   "ModifyDedicatedHostAttributeValue",
		"description":           "ModifyDedicatedHostAttributeValue",
		"network_attributes": []map[string]interface{}{
			{
				"udp_timeout":     80,
				"slb_udp_timeout": 80,
			},
		},
		"dedicated_host_cluster_id": "ModifyDedicatedHostAttributeValue",
	}
	diff, err = newInstanceDiff("alicloud_ecs_dedicated_host", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDedicatedHosts Response
		"DedicatedHosts": map[string]interface{}{
			"DedicatedHost": []interface{}{
				map[string]interface{}{
					"ActionOnMaintenance": "ModifyDedicatedHostAttributeValue",
					"AutoPlacement":       "ModifyDedicatedHostAttributeValue",
					"CpuOverCommitRatio":  "ModifyDedicatedHostAttributeValue",
					"DedicatedHostName":   "ModifyDedicatedHostAttributeValue",
					"Description":         "ModifyDedicatedHostAttributeValue",
					"NetworkAttributes": map[string]interface{}{
						"SlbUdpTimeout": 80,
						"UdpTimeout":    80,
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
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
		err := resourceEcsDedicatedHostUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_dedicated_host"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribeDedicatedHosts" {
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
		err := resourceEcsDedicatedHostRead(dExisted, rawClient)
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
	err = resourceEcsDedicatedHostDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "IncorrectHostStatus.Initializing", "nil", "InvalidDedicatedHostId.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ReleaseDedicatedHost" {
				switch errorCode {
				case "NonRetryableError", "InvalidDedicatedHostId.NotFound":
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
		err := resourceEcsDedicatedHostDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidDedicatedHostId.NotFound":
			assert.Nil(t, err)
		}
	}
}
