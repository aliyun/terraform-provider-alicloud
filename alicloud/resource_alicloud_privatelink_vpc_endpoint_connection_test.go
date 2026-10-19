package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_privatelink_vpc_endpoint_connection",
		&resource.Sweeper{
			Name: "alicloud_privatelink_vpc_endpoint_connection",
			F:    testSweepPrivatelinkVpcEndpointConnection,
		})
}

// testSweepPrivatelinkVpcEndpointConnection disables every endpoint connection
// belonging to test-owned VPC endpoint services so that both the service sweeper
// and the endpoint sweeper can delete their parents cleanly. The disable mirrors
// resourceAliCloudPrivateLinkVpcEndpointConnectionDelete
// (DisableVpcEndpointConnection).
func testSweepPrivatelinkVpcEndpointConnection(region string) error {
	if testSweepPreCheckWithRegions(region, false, connectivity.PrivateLinkRegions) {
		log.Printf("[INFO] Skipping privatelink unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	services, err := listAllPrivatelinkVpcEndpointServices(client)
	if err != nil {
		log.Printf("[ERROR] Failed to list privatelink vpc endpoint services: %s", err)
		return nil
	}
	for _, svc := range services {
		serviceId := fmt.Sprint(svc["ServiceId"])
		if !privatelinkSweepParentShouldDelete(privatelinkServiceSweepName(svc)) {
			log.Printf("[INFO] Skipping privatelink vpc endpoint connection under protected service: %s", serviceId)
			continue
		}
		listAction := "ListVpcEndpointConnections"
		listReq := map[string]interface{}{
			"ServiceId":  serviceId,
			"MaxResults": PageSizeLarge,
		}
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			response, err := client.RpcPost("Privatelink", "2020-04-15", listAction, nil, listReq, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
					break
				}
				log.Printf("[ERROR] Failed to list privatelink vpc endpoint connections for service %s: %s", serviceId, err)
				break
			}
			resp, err := jsonpath.Get("$.Connections", response)
			if err != nil {
				log.Printf("[ERROR] Failed to parse $.Connections for service %s: %s", serviceId, err)
				break
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				endpointId := fmt.Sprint(item["EndpointId"])
				deleteAction := "DisableVpcEndpointConnection"
				deleteReq := map[string]interface{}{
					"ServiceId":   serviceId,
					"EndpointId":  endpointId,
					"RegionId":    client.RegionId,
					"ClientToken": buildClientToken(deleteAction),
				}
				if _, err := client.RpcPost("Privatelink", "2020-04-15", deleteAction, nil, deleteReq, true); err != nil {
					if IsExpectedErrors(err, []string{"EndpointServiceNotFound", "EndpointNotFound"}) {
						continue
					}
					log.Printf("[ERROR] Failed to disable privatelink vpc endpoint connection %s:%s: %s", serviceId, endpointId, err)
					continue
				}
				log.Printf("[INFO] Delete privatelink vpc endpoint connection success: %s:%s", serviceId, endpointId)
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				listReq["NextToken"] = nextToken
			} else {
				break
			}
		}
	}
	return nil
}

func TestAccAliCloudPrivatelinkVpcEndpointConnection_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivatelinkVpcEndpointConnectionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivatelinkService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivatelinkVpcEndpointConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudPrivatelinkVpcEndpointConnectionBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.default.id}",
					"service_id":  "${alicloud_privatelink_vpc_endpoint_service.default.id}",
					"bandwidth":   "1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id": CHECKSET,
						"service_id":  CHECKSET,
						"bandwidth":   "1024",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "1000",
					}),
				),
			},
		},
	})
}

var AlicloudPrivatelinkVpcEndpointConnectionMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudPrivatelinkVpcEndpointConnectionBasicDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	     name_regex = "^default-NODELETING$"
	}
	resource "alicloud_security_group" "default" {
	 name = "tf-testAcc-for-privatelink"
	 description = "privatelink test security group"
	 vpc_id = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	service_description = "test for privatelink connection"
	connect_bandwidth = 103
	auto_accept_connection = false
	}
	resource "alicloud_privatelink_vpc_endpoint" "default" {
	 service_id = alicloud_privatelink_vpc_endpoint_service.default.id
	 vpc_id = data.alicloud_vpcs.default.ids.0
	 security_group_ids = [alicloud_security_group.default.id]
	 vpc_endpoint_name = "testformaintf"
	 depends_on = [alicloud_privatelink_vpc_endpoint_service.default]
	}
`)
}

// lintignore: R001
func TestUnitAlicloudPrivatelinkVpcEndpointConnection(t *testing.T) {
	p := Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_connection"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_connection"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"endpoint_id": "EnableVpcEndpointConnectionValue",
		"service_id":  "EnableVpcEndpointConnectionValue",
		"bandwidth":   1000,
		"dry_run":     false,
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
		// ListVpcEndpointConnections
		"Connections": []interface{}{
			map[string]interface{}{
				"ConnectionStatus": "Connected",
				"Bandwidth":        1000,
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// EnableVpcEndpointConnection
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_privatelink_vpc_endpoint_connection", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPrivatelinkClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointConnectionCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// ListVpcEndpointConnections Response
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "EndpointConnectionOperationDenied", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "EnableVpcEndpointConnection" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointConnectionCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_connection"].Schema).Data(dInit.State(), nil)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointConnectionUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateVpcEndpointConnectionAttribute
	attributesDiff := map[string]interface{}{
		"bandwidth": 500,
		"dry_run":   true,
	}
	diff, err := newInstanceDiff("alicloud_privatelink_vpc_endpoint_connection", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_connection"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListVpcEndpointConnections Response
		"Connections": []interface{}{
			map[string]interface{}{
				"Bandwidth": 500,
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "EndpointConnectionOperationDenied", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateVpcEndpointConnectionAttribute" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointConnectionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_connection"].Schema).Data(dExisted.State(), nil)
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
			if *action == "ListVpcEndpointConnections" {
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
		err := resourceAliCloudPrivateLinkVpcEndpointConnectionRead(dExisted, rawClient)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPrivateLinkVpcEndpointConnectionDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_privatelink_vpc_endpoint_connection", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_privatelink_vpc_endpoint_connection"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "EndpointConnectionOperationDenied", "nil", "EndpointConnectionNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DisableVpcEndpointConnection" {
				switch errorCode {
				case "NonRetryableError", "EndpointConnectionNotFound":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Connections": []interface{}{
								map[string]interface{}{
									"ConnectionStatus": "Disconnected",
								},
							},
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudPrivateLinkVpcEndpointConnectionDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "EndpointConnectionNotFound":
			assert.Nil(t, err)
		}
	}

}

// Test PrivateLink VpcEndpointConnection. >>> Resource test cases, automatically generated.
// Case 4902
func TestAccAliCloudPrivateLinkVpcEndpointConnection_basic4902(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointConnectionMap4902)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpointConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointConnectionBasicDependence4902)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.defaultEndpoint.id}",
					"service_id":  "${alicloud_privatelink_vpc_endpoint_service.defaultService.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id": CHECKSET,
						"service_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "3072",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "3072",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "4096",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "4096",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "3072",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "3072",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.defaultEndpoint.id}",
					"service_id":  "${alicloud_privatelink_vpc_endpoint_service.defaultService.id}",
					"bandwidth":   "3072",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id": CHECKSET,
						"service_id":  CHECKSET,
						"bandwidth":   "3072",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudPrivateLinkVpcEndpointConnectionMap4902 = map[string]string{
	"status": CHECKSET,
}

func AlicloudPrivateLinkVpcEndpointConnectionBasicDependence4902(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "eu-west-1"
}

resource "alicloud_vpc" "defaultVpc" {
  description = "test-zejun"
  cidr_block  = "10.0.0.0/8"
  vpc_name    = var.name

}

resource "alicloud_security_group" "defaultSg" {
  name = var.name

  vpc_id = alicloud_vpc.defaultVpc.id
}

resource "alicloud_privatelink_vpc_endpoint_service" "defaultService" {
  service_description   = "test-zejun"
  connect_bandwidth     = "3072"
  service_resource_type = "slb"
}

resource "alicloud_privatelink_vpc_endpoint" "defaultEndpoint" {
  vpc_id     = alicloud_vpc.defaultVpc.id
  service_id = alicloud_privatelink_vpc_endpoint_service.defaultService.id
  security_group_ids = [alicloud_security_group.defaultSg.id]
}


`, name)
}

// Case 4902  twin
func TestAccAliCloudPrivateLinkVpcEndpointConnection_basic4902_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_privatelink_vpc_endpoint_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateLinkVpcEndpointConnectionMap4902)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrivateLinkServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateLinkVpcEndpointConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sprivatelinkvpcendpointconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateLinkVpcEndpointConnectionBasicDependence4902)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint_id": "${alicloud_privatelink_vpc_endpoint.defaultEndpoint.id}",
					"service_id":  "${alicloud_privatelink_vpc_endpoint_service.defaultService.id}",
					"bandwidth":   "3072",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id": CHECKSET,
						"service_id":  CHECKSET,
						"bandwidth":   "3072",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

// Test PrivateLink VpcEndpointConnection. <<< Resource test cases, automatically generated.
