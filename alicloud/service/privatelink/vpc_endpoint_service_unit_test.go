package privatelink

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitVpcEndpointServiceConnectBandwidthReadByResourceType(t *testing.T) {
	p := unitTestProviderResources()

	cases := []struct {
		name                string
		serviceResourceType string
		response            map[string]interface{}
		expected            int
	}{
		{
			name:                "slb missing connect bandwidth is authoritative",
			serviceResourceType: "slb",
			response: map[string]interface{}{
				"ServiceResourceType": "slb",
			},
			expected: 0,
		},
		{
			name:                "slb returned connect bandwidth is authoritative",
			serviceResourceType: "slb",
			response: map[string]interface{}{
				"ServiceResourceType": "slb",
				"ConnectBandwidth":    100,
			},
			expected: 100,
		},
		{
			name:                "nlb missing connect bandwidth preserves config",
			serviceResourceType: "nlb",
			response: map[string]interface{}{
				"ServiceResourceType": "nlb",
			},
			expected: 3072,
		},
		{
			name:                "alb zero connect bandwidth preserves config",
			serviceResourceType: "alb",
			response: map[string]interface{}{
				"ServiceResourceType": "alb",
				"ConnectBandwidth":    0,
			},
			expected: 3072,
		},
		{
			name:                "gwlb empty connect bandwidth preserves config",
			serviceResourceType: "gwlb",
			response: map[string]interface{}{
				"ServiceResourceType": "gwlb",
				"ConnectBandwidth":    "",
			},
			expected: 3072,
		},
	}

	for _, tc := range cases {
		d, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(nil, nil)
		assert.Nil(t, d.Set("service_resource_type", tc.serviceResourceType))
		assert.Nil(t, d.Set("connect_bandwidth", 3072))
		setVpcEndpointServiceConnectBandwidth(d, tc.response)
		assert.Equal(t, tc.expected, d.Get("connect_bandwidth"), tc.name)
	}
}

// lintignore: R001
func TestUnitAlicloudPrivatelinkVpcEndpointService(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"service_description":    "CreateVpcEndpointServiceValue",
		"connect_bandwidth":      100,
		"auto_accept_connection": false,
		"dry_run":                false,
		"payer":                  "CreateVpcEndpointServiceValue",
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
	rawClient, err := unitTestSharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// GetVpcEndpointServiceAttribute
		"AutoAcceptEnabled":     false,
		"ConnectBandwidth":      100,
		"ServiceBusinessStatus": "CreateVpcEndpointServiceValue",
		"ServiceDescription":    "CreateVpcEndpointServiceValue",
		"ServiceDomain":         "CreateVpcEndpointServiceValue",
		"ServiceStatus":         "Active",
		// ListVpcEndpointServiceResources — an empty array keeps jsonpath.Get
		// successful and reflects no attached service resources in unit tests.
		"Resources": []interface{}{},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateVpcEndpoint
		"ServiceId": "CreateVpcEndpointServiceValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_privatelink_vpc_endpoint_service", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointServiceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetVpcEndpointServiceAttribute Response
		"ServiceId": "CreateVpcEndpointServiceValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateVpcEndpointService" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointServiceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointServiceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateVpcEndpointServiceAttribute
	attributesDiff := map[string]interface{}{
		"auto_accept_connection": true,
		"connect_bandwidth":      200,
		"service_description":    "UpdateVpcEndpointServiceAttributeValue",
		"dry_run":                true,
	}
	diff, err := newInstanceDiff("alicloud_privatelink_vpc_endpoint_service", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetVpcEndpointServiceAttribute Response
		"AutoAcceptEnabled":  true,
		"ConnectBandwidth":   200,
		"ServiceDescription": "UpdateVpcEndpointServiceAttributeValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateVpcEndpointServiceAttribute" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointServiceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(dExisted.State(), nil)
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
			if *action == "GetVpcEndpointServiceAttribute" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointServiceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointServiceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "EndpointServiceConnectionDependence", "nil", "EndpointServiceNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVpcEndpointService" {
				switch errorCode {
				case "NonRetryableError", "EndpointServiceNotFound":
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
		err := resourceAliCloudPrivateLinkVpcEndpointServiceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "EndpointServiceNotFound":
			assert.Nil(t, err)
		}
	}

}

func TestUnitVpcEndpointServiceResourceZoneIdentity(t *testing.T) {
	p := unitTestProviderResources()
	setSchema := p["alicloud_privatelink_vpc_endpoint_service"].Schema["resource"]
	hashFn := setSchema.Set

	// zone_id participates in the TypeSet element hash for every resource
	// type (the fix from #11104), so two inline resource blocks with the same
	// resource_id + resource_type but a different zone_id are preserved as
	// distinct attachments instead of collapsing into a single block.
	for _, rt := range []string{"alb", "nlb", "gwlb", "slb"} {
		assert.NotEqual(t, hashFn(map[string]interface{}{
			"resource_id":   "lb-bp1xxx",
			"resource_type": rt,
			"zone_id":       "cn-hangzhou-h",
		}), hashFn(map[string]interface{}{
			"resource_id":   "lb-bp1xxx",
			"resource_type": rt,
			"zone_id":       "cn-hangzhou-i",
		}), "%s resource must preserve zone identity in hash", rt)
	}

	// A non-empty zone_id is returned unchanged; a genuinely absent/nil value
	// normalizes to "" so it is never sent as a literal ZoneId upstream. The
	// service-returned value is always preserved, never discarded.
	assert.Equal(t, "cn-hangzhou-h", vpcEndpointServiceResourceZoneId("cn-hangzhou-h"), "non-empty zone must be preserved")
	assert.Equal(t, "", vpcEndpointServiceResourceZoneId(nil), "nil zone must normalize to empty")
	assert.Equal(t, "", vpcEndpointServiceResourceZoneId(""), "empty zone must stay empty")
}

func TestUnitVpcEndpointServiceResourceReadZoneNormalization(t *testing.T) {
	p := unitTestProviderResources()
	d, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_service"].Schema).Data(nil, nil)
	d.SetId("eps-test-zone-norm")

	// The Read path derives a Privatelink client from a *connectivity.AliyunClient,
	// but no real cloud credentials are exercised here: DoRequest is patched below
	// so every RpcPost returns the mock and no live API call is made. Inject dummy
	// credentials via t.Setenv (restored automatically when the test ends) so the
	// shared helper builds the client instead of skipping, keeping this assertion
	// running in CI without secrets.
	t.Setenv("ALICLOUD_ACCESS_KEY", "unit-test-access-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "unit-test-secret-key")
	rawClient, err := unitTestSharedClientForRegion("cn-hangzhou")
	assert.Nil(t, err)
	if err != nil {
		t.Fatalf("failed to build unit test client: %s", err)
	}
	rawClient = rawClient.(*connectivity.AliyunClient)

	ReadMockResponse := map[string]interface{}{
		// GetVpcEndpointServiceAttribute
		"ServiceId":           "eps-test-zone-norm",
		"ServiceDescription":  "zone-norm-test",
		"ServiceResourceType": "slb",
		"ServiceStatus":       "Active",
		"AutoAcceptEnabled":   true,
		// ListVpcEndpointServiceResources — the upstream can return a non-empty
		// ZoneId for an attached backend (e.g. the zone the SLB or NLB was
		// created in). Read must preserve the service-returned value rather
		// than discarding it; zone_id is Optional+Computed, so the observable
		// value must stay available to consuming expressions.
		"Resources": []interface{}{
			map[string]interface{}{
				"ResourceId":   "lb-bp1xxx",
				"ResourceType": "slb",
				"ZoneId":       "cn-hangzhou-h",
			},
			map[string]interface{}{
				"ResourceId":   "lb-bp2xxx",
				"ResourceType": "nlb",
				"ZoneId":       "cn-hangzhou-g",
			},
		},
		// ListTagResources — no tags.
		"TagResources": []interface{}{},
	}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
		return ReadMockResponse, nil
	})
	defer patches.Reset()

	err = resourceAliCloudPrivateLinkVpcEndpointServiceRead(d, rawClient)
	assert.Nil(t, err)

	resources := d.Get("resource").(*schema.Set).List()
	assert.Equal(t, 2, len(resources), "expected exactly two resource items after Read")

	byType := map[string]string{}
	for _, r := range resources {
		item := r.(map[string]interface{})
		byType[item["resource_type"].(string)] = item["zone_id"].(string)
	}
	assert.Equal(t, "cn-hangzhou-h", byType["slb"], "slb zone_id must preserve the service-returned value, not be normalized to empty")
	assert.Equal(t, "cn-hangzhou-g", byType["nlb"], "nlb zone_id must preserve the service-returned read-back zone")
}
