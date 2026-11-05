package ga

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

// White-box (internal package ga) TestUnit cases moved out of the
// external ga_test file so they can call the unexported resource CRUD
// funcs. Shared local helpers (unitTestProviderResources /
// newInstanceDiff / unitTestSharedClientForRegion) live in
// unit_test_helpers_test.go; see that file for why acctest cannot be
// imported here.

// lintignore: R001
func TestUnitAliCloudGaListener(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"accelerator_id": "CreateListenerValue",
		"description":    "CreateListenerValue",
		"name":           "CreateListenerValue",
		"port_ranges": []map[string]interface{}{
			{
				"from_port": 60,
				"to_port":   70,
			},
		},
		"proxy_protocol": true,
		"certificates": []map[string]interface{}{
			{
				"id": "CreateListenerValue",
			},
		},
		"client_affinity": "CreateListenerValue",
		"protocol":        "CreateListenerValue",
		"backend_ports": []map[string]interface{}{
			{
				"from_port": 80,
				"to_port":   80,
			},
		},
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
		// DescribeListener
		"Description": "CreateListenerValue",
		"Certificates": []interface{}{
			map[string]interface{}{
				"Id": "CreateListenerValue",
			},
		},
		"ClientAffinity": "CreateListenerValue",
		"Name":           "CreateListenerValue",
		"PortRanges": []interface{}{
			map[string]interface{}{
				"FromPort": 60,
				"ToPort":   70,
			},
		},
		"BackendPorts": []interface{}{
			map[string]interface{}{
				"FromPort": 80,
				"ToPort":   80,
			},
		},
		"State":      "active",
		"Protocol":   "CreateListenerValue",
		"ListenerId": "CreateListenerValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateListener
		"ListenerId": "CreateListenerValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ga_listener", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaListenerCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeListener Response
		"ListenerId": "CreateListenerValue",
	}
	errorCodes := []string{"NonRetryableError", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateListener" {
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
		err := resourceAliCloudGaListenerCreate(dInit, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaListenerUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateListener
	attributesDiff := map[string]interface{}{
		"accelerator_id": "UpdateListenerValue",
		"description":    "UpdateListenerValue",
		"name":           "UpdateListenerValue",
		"port_ranges": []map[string]interface{}{
			{
				"from_port": 70,
				"to_port":   80,
			},
		},
		"proxy_protocol": true,
		"certificates": []map[string]interface{}{
			{
				"id": "UpdateListenerValue",
			},
		},
		"client_affinity": "UpdateListenerValue",
		"protocol":        "UpdateListenerValue",
		"backend_ports": []map[string]interface{}{
			{
				"from_port": 90,
				"to_port":   90,
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_ga_listener", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeListener Response
		"Description": "UpdateListenerValue",
		"Certificates": []interface{}{
			map[string]interface{}{
				"Id": "UpdateListenerValue",
			},
		},
		"ClientAffinity": "UpdateListenerValue",
		"Name":           "UpdateListenerValue",
		"PortRanges": []interface{}{
			map[string]interface{}{
				"FromPort": 70,
				"ToPort":   80,
			},
		},
		"BackendPorts": []interface{}{
			map[string]interface{}{
				"FromPort": 90,
				"ToPort":   90,
			},
		},
		"Protocol":   "UpdateListenerValue",
		"ListenerId": "UpdateListenerValue",
	}
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateListener" {
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
		err := resourceAliCloudGaListenerUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_listener"].Schema).Data(dExisted.State(), nil)
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
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeListener" {
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
		err := resourceAliCloudGaListenerRead(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaListenerDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteListener" {
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
			if *action == "DescribeListener" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaListenerDelete(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}

func TestUnitGaListenerNestedConfigurationBuilders(t *testing.T) {
	endpointGroups := buildGaListenerEndpointGroupConfigurations([]interface{}{
		map[string]interface{}{
			"endpoint_group_region":                        "cn-hangzhou",
			"endpoint_group_name":                          "egc-name",
			"endpoint_group_description":                   "egc-description",
			"traffic_percentage":                           100,
			"health_check_enabled":                         true,
			"health_check_interval_seconds":                5,
			"health_check_path":                            "/health",
			"health_check_port":                            80,
			"health_check_protocol":                        "HTTP",
			"threshold_count":                              3,
			"endpoint_request_protocol":                    "HTTPS",
			"endpoint_group_type":                          "default",
			"enable_client_ip_preservation_proxy_protocol": true,
			"enable_client_ip_preservation_toa":            true,
			"endpoint_protocol_version":                    "HTTP2",
			"endpoint_ip_version":                          "IPv6",
			"health_check_host":                            "health.example.com",
			"endpoint_configurations": []interface{}{
				map[string]interface{}{
					"type":                          "PrivateSubNet",
					"weight":                        20,
					"endpoint":                      "10.0.0.1",
					"api_keys":                      []interface{}{"key-1"},
					"vpc_id":                        "vpc-123",
					"vswitch_ids":                   []interface{}{"vsw-1", "vsw-2"},
					"sub_address":                   "192.168.0.1",
					"enable_proxy_protocol":         true,
					"provider":                      "custom",
					"enable_client_ip_preservation": true,
				},
			},
			"port_overrides": []interface{}{
				map[string]interface{}{
					"listener_port": 80,
					"endpoint_port": 8080,
				},
			},
		},
	})
	assert.Len(t, endpointGroups, 1)
	egc := endpointGroups[0]
	assert.Equal(t, "cn-hangzhou", egc["EndpointGroupRegion"])
	assert.Equal(t, "egc-name", egc["EndpointGroupName"])
	assert.Equal(t, "egc-description", egc["EndpointGroupDescription"])
	assert.Equal(t, 100, egc["TrafficPercentage"])
	assert.Equal(t, true, egc["HealthCheckEnabled"])
	assert.Equal(t, 5, egc["HealthCheckIntervalSeconds"])
	assert.Equal(t, "/health", egc["HealthCheckPath"])
	assert.Equal(t, 80, egc["HealthCheckPort"])
	assert.Equal(t, "HTTP", egc["HealthCheckProtocol"])
	assert.Equal(t, 3, egc["ThresholdCount"])
	assert.Equal(t, "HTTPS", egc["EndpointRequestProtocol"])
	assert.Equal(t, "default", egc["EndpointGroupType"])
	assert.Equal(t, true, egc["EnableClientIpPreservationProxyProtocol"])
	assert.Equal(t, true, egc["EnableClientIpPreservationToa"])
	assert.Equal(t, "HTTP2", egc["EndpointProtocolVersion"])
	assert.Equal(t, "IPv6", egc["EndpointIpVersion"])
	assert.Equal(t, "health.example.com", egc["HealthCheckHost"])

	endpoints := egc["EndpointConfigurations"].([]map[string]interface{})
	assert.Len(t, endpoints, 1)
	assert.Equal(t, "PrivateSubNet", endpoints[0]["Type"])
	assert.Equal(t, 20, endpoints[0]["Weight"])
	assert.Equal(t, "10.0.0.1", endpoints[0]["Endpoint"])
	assert.Equal(t, []string{"key-1"}, endpoints[0]["ApiKeys"])
	assert.Equal(t, "vpc-123", endpoints[0]["VpcId"])
	assert.Equal(t, []string{"vsw-1", "vsw-2"}, endpoints[0]["VSwitchIds"])
	assert.Equal(t, "192.168.0.1", endpoints[0]["SubAddress"])
	assert.Equal(t, true, endpoints[0]["EnableProxyProtocol"])
	assert.Equal(t, "custom", endpoints[0]["Provider"])
	assert.Equal(t, true, endpoints[0]["EnableClientIPPreservation"])

	portOverrides := egc["PortOverrides"].([]map[string]interface{})
	assert.Len(t, portOverrides, 1)
	assert.Equal(t, 80, portOverrides[0]["ListenerPort"])
	assert.Equal(t, 8080, portOverrides[0]["EndpointPort"])

	customRouting := buildGaListenerCustomRoutingEndpointGroupConfigurations([]interface{}{
		map[string]interface{}{
			"endpoint_group_region": "cn-hangzhou",
			"name":                  "cr-name",
			"description":           "cr-description",
			"destination_configurations": []interface{}{
				map[string]interface{}{
					"protocols": []interface{}{"TCP", "UDP"},
					"from_port": 10000,
					"to_port":   10001,
				},
			},
			"custom_routing_endpoint_configurations": []interface{}{
				map[string]interface{}{
					"type":                       "PrivateSubNet",
					"endpoint":                   "10.0.0.0/24",
					"traffic_to_endpoint_policy": "AllowCustom",
					"policy_configurations": []interface{}{
						map[string]interface{}{
							"address": "10.0.0.1",
							"custom_routing_port_ranges": []interface{}{
								map[string]interface{}{
									"from_port": 10000,
									"to_port":   10005,
								},
							},
						},
					},
				},
			},
		},
	})
	assert.Len(t, customRouting, 1)
	cr := customRouting[0]
	assert.Equal(t, "cn-hangzhou", cr["EndpointGroupRegion"])
	assert.Equal(t, "cr-name", cr["Name"])
	assert.Equal(t, "cr-description", cr["Description"])

	destinations := cr["DestinationConfigurations"].([]map[string]interface{})
	assert.Len(t, destinations, 1)
	assert.Equal(t, []string{"TCP", "UDP"}, destinations[0]["Protocols"])
	assert.Equal(t, 10000, destinations[0]["FromPort"])
	assert.Equal(t, 10001, destinations[0]["ToPort"])

	crEndpoints := cr["CustomRoutingEndpointConfigurations"].([]map[string]interface{})
	assert.Len(t, crEndpoints, 1)
	assert.Equal(t, "PrivateSubNet", crEndpoints[0]["Type"])
	assert.Equal(t, "10.0.0.0/24", crEndpoints[0]["Endpoint"])
	assert.Equal(t, "AllowCustom", crEndpoints[0]["TrafficToEndpointPolicy"])

	policies := crEndpoints[0]["PolicyConfigurations"].([]map[string]interface{})
	assert.Len(t, policies, 1)
	assert.Equal(t, "10.0.0.1", policies[0]["Address"])
	ranges := policies[0]["PortRanges"].([]map[string]interface{})
	assert.Len(t, ranges, 1)
	assert.Equal(t, 10000, ranges[0]["FromPort"])
	assert.Equal(t, 10005, ranges[0]["ToPort"])

	backendPorts := buildGaListenerBackendPorts([]interface{}{
		map[string]interface{}{"from_port": 80, "to_port": 80},
	})
	assert.Len(t, backendPorts, 1)
	assert.Equal(t, 80, backendPorts[0]["FromPort"])
	assert.Equal(t, 80, backendPorts[0]["ToPort"])
}
