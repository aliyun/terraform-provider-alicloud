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

func TestAccAliCloudGaEndpointGroup_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AliCloudGaEndpointGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaEndpointGroupBasicDependence)
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
					"accelerator_id":        "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":           "${alicloud_ga_listener.default.id}",
					"endpoint_group_region": defaultRegionToTest,
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "true",
							"enable_clientip_preservation": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":            CHECKSET,
						"listener_id":               CHECKSET,
						"endpoint_group_region":     defaultRegionToTest,
						"endpoint_configurations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
						{
							"endpoint":                     "${alicloud_eip_address.default.1.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "www.alicloud-provider.cn",
							"type":                         "Domain",
							"weight":                       "30",
							"enable_proxy_protocol":        "true",
							"enable_clientip_preservation": "false",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_request_protocol": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_request_protocol": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "EndpointGroup_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "EndpointGroup_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval_seconds": `5`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval_seconds": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_path": "/healthcheckupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_path": "/healthcheckupdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_port": `30`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_port": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_protocol": "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_protocol": "http",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "60",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_overrides.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_count": `5`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_count": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_percentage": `30`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_percentage": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "false",
							"enable_clientip_preservation": "false",
						},
					},
					"description":                   "EndpointGroup",
					"health_check_interval_seconds": `3`,
					"health_check_path":             "/healthcheck",
					"health_check_port":             `20`,
					"health_check_protocol":         "tcp",
					"name":                          name,
					"port_overrides": []map[string]interface{}{
						{
							"endpoint_port": "10",
							"listener_port": "60",
						},
					},
					"threshold_count":    `3`,
					"traffic_percentage": `20`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_configurations.#":     "1",
						"description":                   "EndpointGroup",
						"health_check_interval_seconds": "3",
						"health_check_path":             "/healthcheck",
						"health_check_port":             "20",
						"health_check_protocol":         "tcp",
						"name":                          name,
						"port_overrides.#":              "1",
						"threshold_count":               "3",
						"traffic_percentage":            "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "EndpointGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "EndpointGroup",
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

func TestAccAliCloudGaEndpointGroup_basic01(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_endpoint_group.default"
	ra := resourceAttrInit(resourceId, AliCloudGaEndpointGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaEndpointGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaEndpointGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaEndpointGroupBasicDependence)
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
					"accelerator_id":        "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":           "${alicloud_ga_listener.default.id}",
					"endpoint_group_region": defaultRegionToTest,
					"endpoint_group_type":   "virtual",
					"endpoint_configurations": []map[string]interface{}{
						{
							"endpoint":                     "${alicloud_eip_address.default.0.ip_address}",
							"type":                         "PublicIp",
							"weight":                       "20",
							"enable_proxy_protocol":        "true",
							"enable_clientip_preservation": "false",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "EndpointGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":            CHECKSET,
						"listener_id":               CHECKSET,
						"endpoint_group_region":     defaultRegionToTest,
						"endpoint_group_type":       "virtual",
						"endpoint_configurations.#": "1",
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "EndpointGroup",
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

var AliCloudGaEndpointGroupMap = map[string]string{
	"endpoint_group_type":       CHECKSET,
	"endpoint_request_protocol": CHECKSET,
	"threshold_count":           CHECKSET,
	"endpoint_group_ip_list.#":  CHECKSET,
	"status":                    CHECKSET,
}

func AliCloudGaEndpointGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth              = 100
  		type                   = "Basic"
  		bandwidth_type         = "Enhanced"
  		payment_type           = "PayAsYouGo"
  		billing_type           = "PayBy95"
  		ratio                  = 30
  		bandwidth_package_name = var.name
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
  		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}

	resource "alicloud_ga_listener" "default" {
  		port_ranges {
    		from_port = "60"
    		to_port   = "60"
  		}
  		accelerator_id  = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		client_affinity = "SOURCE_IP"
  		protocol        = "HTTP"
  		name            = var.name
	}

	resource "alicloud_eip_address" "default" {
  		count                = 2
  		bandwidth            = "10"
  		internet_charge_type = "PayByBandwidth"
  		address_name         = var.name
	}
`, name)
}

func TestUnitAliCloudGaEndpointGroup(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"accelerator_id":        "CreateEndpointGroupValue",
		"description":           "CreateEndpointGroupValue",
		"endpoint_group_region": "CreateEndpointGroupValue",
		"endpoint_configurations": []map[string]interface{}{
			{
				"enable_clientip_preservation": true,
				"endpoint":                     "CreateEndpointGroupValue",
				"type":                         "PublicIp",
				"weight":                       20,
			},
		},
		"endpoint_group_type":           "CreateEndpointGroupValue",
		"listener_id":                   "CreateEndpointGroupValue",
		"endpoint_request_protocol":     "CreateEndpointGroupValue",
		"health_check_interval_seconds": 3,
		"health_check_path":             "CreateEndpointGroupValue",
		"health_check_port":             20,
		"health_check_protocol":         "CreateEndpointGroupValue",
		"name":                          "CreateEndpointGroupValue",
		"port_overrides": []map[string]interface{}{
			{
				"endpoint_port": 10,
				"listener_port": 60,
			},
		},
		"threshold_count":    3,
		"traffic_percentage": 20,
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
		// DescribeEndpointGroup
		"AcceleratorId": "CreateEndpointGroupValue",
		"Description":   "CreateEndpointGroupValue",
		"EndpointConfigurations": []interface{}{
			map[string]interface{}{
				"EnableClientIPPreservation": true,
				"Endpoint":                   "CreateEndpointGroupValue",
				"Type":                       "PublicIp",
				"Weight":                     20,
			},
		},
		"EndpointGroupRegion":        "CreateEndpointGroupValue",
		"EndpointGroupType":          "CreateEndpointGroupValue",
		"HealthCheckIntervalSeconds": 3,
		"HealthCheckPath":            "CreateEndpointGroupValue",
		"HealthCheckPort":            20,
		"HealthCheckProtocol":        "CreateEndpointGroupValue",
		"ListenerId":                 "CreateEndpointGroupValue",
		"EndpointRequestProtocol":    "CreateEndpointGroupValue",
		"Name":                       "CreateEndpointGroupValue",
		"PortOverrides": []interface{}{
			map[string]interface{}{
				"EndpointPort": 10,
				"ListenerPort": 60,
			},
		},
		"State":             "active",
		"ThresholdCount":    3,
		"TrafficPercentage": 20,
	}
	CreateMockResponse := map[string]interface{}{
		// CreateEndpointGroup
		"EndpointGroupId": "CreateEndpointGroupValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_endpoint_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaEndpointGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeEndpointGroup Response
		"EndpointGroupId": "CreateEndpointGroupValue",
	}
	errorCodes := []string{"NonRetryableError", "GA_NOT_STEADY", "StateError.Accelerator", "StateError.EndPointGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateEndpointGroup" {
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
		err := resourceAliCloudGaEndpointGroupCreate(dInit, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(dInit.State(), nil)
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
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaEndpointGroupUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateEndpointGroup
	attributesDiff := map[string]interface{}{
		"description": "UpdateEndpointGroup",
		"endpoint_configurations": []map[string]interface{}{
			{
				"enable_clientip_preservation": false,
				"endpoint":                     "UpdateEndpointGroup",
				"type":                         "UpdateEndpointGroup",
				"weight":                       30,
			},
		},
		"endpoint_request_protocol":     "UpdateEndpointGroup",
		"health_check_interval_seconds": 4,
		"health_check_path":             "UpdateEndpointGroup",
		"health_check_port":             30,
		"health_check_protocol":         "UpdateEndpointGroup",
		"name":                          "UpdateEndpointGroup",
		"port_overrides": []map[string]interface{}{
			{
				"endpoint_port": 20,
				"listener_port": 70,
			},
		},
		"threshold_count":    4,
		"traffic_percentage": 30,
	}
	diff, err := newInstanceDiff("alicloud_ga_endpoint_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeEndpointGroup Response
		"Description": "UpdateEndpointGroup",
		"EndpointConfigurations": []interface{}{
			map[string]interface{}{
				"EnableClientIPPreservation": false,
				"Endpoint":                   "UpdateEndpointGroup",
				"Type":                       "UpdateEndpointGroup",
				"Weight":                     30,
			},
		},
		"HealthCheckIntervalSeconds": 4,
		"HealthCheckPath":            "UpdateEndpointGroup",
		"HealthCheckPort":            30,
		"HealthCheckProtocol":        "UpdateEndpointGroup",
		"EndpointRequestProtocol":    "UpdateEndpointGroup",
		"Name":                       "UpdateEndpointGroup",
		"PortOverrides": []interface{}{
			map[string]interface{}{
				"EndpointPort": 20,
				"ListenerPort": 70,
			},
		},
		"State":             "active",
		"ThresholdCount":    4,
		"TrafficPercentage": 30,
	}
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "StateError.EndPointGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateEndpointGroup" {
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
		err := resourceAliCloudGaEndpointGroupUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribeEndpointGroup" {
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
		err := resourceAliCloudGaEndpointGroupRead(dExisted, rawClient)
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
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaEndpointGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "StateError.EndPointGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteEndpointGroup" {
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
			if *action == "DeleteEndpointGroup" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaEndpointGroupDelete(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
