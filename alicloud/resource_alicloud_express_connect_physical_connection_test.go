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

func TestAccAliCloudExpressConnectPhysicalConnection_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudExpressConnectPhysicalConnectionBasicDependence0)
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
					// currently， not all access points are available
					//"access_point_id":          "${data.alicloud_express_connect_access_points.default.ids.0}",
					"access_point_id": getAccessPointId(),
					"line_operator":   "CU",
					"port_type":       "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id": CHECKSET,
						"line_operator":   "CU",
						"port_type":       "1000Base-LX",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line_operator": "CM",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"line_operator": "CM",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"circuit_code": "longtel001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"circuit_code": "longtel001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_location": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_location": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Confirmed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Confirmed",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
					"period": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":   "Enabled",
						"order_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Terminated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Terminated",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "pricing_cycle", "order_id"},
			},
		},
	})
}

func TestAccAliCloudExpressConnectPhysicalConnection_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudExpressConnectPhysicalConnectionBasicDependence0)
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
					// currently， not all access points are available
					//"access_point_id":          "${data.alicloud_express_connect_access_points.default.ids.0}",
					"access_point_id":                  getAccessPointId(),
					"line_operator":                    "CU",
					"type":                             "VPC",
					"port_type":                        "1000Base-LX",
					"bandwidth":                        "10",
					"circuit_code":                     "longtel001",
					"peer_location":                    name,
					"redundant_physical_connection_id": "${data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id}",
					"physical_connection_name":         name,
					"description":                      name,
					"status":                           "Enabled",
					"period":                           "1",
					"pricing_cycle":                    "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id":                  CHECKSET,
						"line_operator":                    "CU",
						"type":                             "VPC",
						"port_type":                        "1000Base-LX",
						"bandwidth":                        "10",
						"circuit_code":                     "longtel001",
						"peer_location":                    name,
						"redundant_physical_connection_id": CHECKSET,
						"physical_connection_name":         name,
						"description":                      name,
						"status":                           "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Terminated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Terminated",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "pricing_cycle", "order_id"},
			},
		},
	})
}

func TestAccAliCloudExpressConnectPhysicalConnection_basic0_intl(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudExpressConnectPhysicalConnectionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// currently， not all access points are available
					//"access_point_id":          "${data.alicloud_express_connect_access_points.default.ids.0}",
					"access_point_id": getAccessPointId(),
					"line_operator":   "Other",
					"port_type":       "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id": CHECKSET,
						"line_operator":   "Other",
						"port_type":       "1000Base-LX",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line_operator": "Equinix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"line_operator": "Equinix",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"circuit_code": "longtel001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"circuit_code": "longtel001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_location": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_location": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Confirmed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Confirmed",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Canceled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Canceled",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "pricing_cycle", "order_id"},
			},
		},
	})
}

func TestAccAliCloudExpressConnectPhysicalConnection_basic0_intl_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AliCloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudExpressConnectPhysicalConnectionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// currently， not all access points are available
					//"access_point_id":          "${data.alicloud_express_connect_access_points.default.ids.0}",
					"access_point_id":                  getAccessPointId(),
					"line_operator":                    "Other",
					"type":                             "VPC",
					"port_type":                        "1000Base-LX",
					"bandwidth":                        "10",
					"circuit_code":                     "longtel001",
					"peer_location":                    name,
					"redundant_physical_connection_id": "${data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id}",
					"physical_connection_name":         name,
					"description":                      name,
					"status":                           "Enabled",
					"period":                           "1",
					"pricing_cycle":                    "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id":                  CHECKSET,
						"line_operator":                    "Other",
						"type":                             "VPC",
						"port_type":                        "1000Base-LX",
						"bandwidth":                        "10",
						"circuit_code":                     "longtel001",
						"peer_location":                    name,
						"redundant_physical_connection_id": CHECKSET,
						"physical_connection_name":         name,
						"description":                      name,
						"status":                           "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Terminated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Terminated",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "pricing_cycle", "order_id"},
			},
		},
	})
}

var AliCloudExpressConnectPhysicalConnectionMap0 = map[string]string{
	"type":          CHECKSET,
	"bandwidth":     CHECKSET,
	"peer_location": CHECKSET,
	"status":        CHECKSET,
}

func AliCloudExpressConnectPhysicalConnectionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_express_connect_access_points" "default" {
  		status = "recommended"
	}

	data "alicloud_express_connect_physical_connections" "nameRegex" {
  		name_regex = "^preserved-NODELETING"
	}
`, name)
}

func TestUnitAliCloudExpressConnectPhysicalConnection(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"access_point_id":                  "CreatePhysicalConnectionValue",
		"bandwidth":                        "10",
		"circuit_code":                     "CreatePhysicalConnectionValue",
		"description":                      "CreatePhysicalConnectionValue",
		"line_operator":                    "CreatePhysicalConnectionValue",
		"peer_location":                    "CreatePhysicalConnectionValue",
		"physical_connection_name":         "CreatePhysicalConnectionValue",
		"port_type":                        "CreatePhysicalConnectionValue",
		"redundant_physical_connection_id": "CreatePhysicalConnectionValue",
		"type":                             "CreatePhysicalConnectionValue",
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
		// DescribePhysicalConnections
		"PhysicalConnectionSet": map[string]interface{}{
			"PhysicalConnectionType": []interface{}{
				map[string]interface{}{
					"AccessPointId":                 "CreatePhysicalConnectionValue",
					"Bandwidth":                     "10",
					"CircuitCode":                   "CreatePhysicalConnectionValue",
					"Description":                   "CreatePhysicalConnectionValue",
					"LineOperator":                  "CreatePhysicalConnectionValue",
					"PeerLocation":                  "CreatePhysicalConnectionValue",
					"Name":                          "CreatePhysicalConnectionValue",
					"PortType":                      "CreatePhysicalConnectionValue",
					"RedundantPhysicalConnectionId": "CreatePhysicalConnectionValue",
					"Status":                        "Allocated",
					"Type":                          "CreatePhysicalConnectionValue",
					"PhysicalConnectionId":          "CreatePhysicalConnectionValue",
				},
			},
		},
		"PhysicalConnectionId": "CreatePhysicalConnectionValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreatePhysicalConnection
		"PhysicalConnectionId": "CreatePhysicalConnectionValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_express_connect_physical_connection", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudExpressConnectPhysicalConnectionCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribePhysicalConnections Response
		"PhysicalConnectionSet": map[string]interface{}{
			"PhysicalConnectionType": []interface{}{
				map[string]interface{}{
					"PhysicalConnectionId": "CreatePhysicalConnectionValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreatePhysicalConnection" {
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
		err := resourceAliCloudExpressConnectPhysicalConnectionCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyPhysicalConnectionAttribute
	attributesDiff := map[string]interface{}{
		"bandwidth":                        "20",
		"circuit_code":                     "ModifyPhysicalConnectionAttributeValue",
		"description":                      "ModifyPhysicalConnectionAttributeValue",
		"line_operator":                    "ModifyPhysicalConnectionAttributeValue",
		"peer_location":                    "ModifyPhysicalConnectionAttributeValue",
		"physical_connection_name":         "ModifyPhysicalConnectionAttributeValue",
		"port_type":                        "ModifyPhysicalConnectionAttributeValue",
		"redundant_physical_connection_id": "ModifyPhysicalConnectionAttributeValue",
	}
	diff, err := newInstanceDiff("alicloud_express_connect_physical_connection", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribePhysicalConnections Response
		"PhysicalConnectionSet": map[string]interface{}{
			"PhysicalConnectionType": []interface{}{
				map[string]interface{}{
					"Bandwidth":                     "20",
					"CircuitCode":                   "ModifyPhysicalConnectionAttributeValue",
					"Description":                   "ModifyPhysicalConnectionAttributeValue",
					"LineOperator":                  "ModifyPhysicalConnectionAttributeValue",
					"PeerLocation":                  "ModifyPhysicalConnectionAttributeValue",
					"Name":                          "ModifyPhysicalConnectionAttributeValue",
					"PortType":                      "ModifyPhysicalConnectionAttributeValue",
					"RedundantPhysicalConnectionId": "ModifyPhysicalConnectionAttributeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyPhysicalConnectionAttribute" {
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
		err := resourceAliCloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// CancelPhysicalConnection
	attributesDiff = map[string]interface{}{
		"status": "Canceled",
	}
	diff, err = newInstanceDiff("alicloud_express_connect_physical_connection", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribePhysicalConnections Response
		"PhysicalConnectionSet": map[string]interface{}{
			"PhysicalConnectionType": []interface{}{
				map[string]interface{}{
					"Status": "Canceled",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CancelPhysicalConnection" {
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
		err := resourceAliCloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
		patches.Reset()

		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// EnablePhysicalConnection
	attributesDiff = map[string]interface{}{
		"status": "Enabled",
	}
	diff, err = newInstanceDiff("alicloud_express_connect_physical_connection", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribePhysicalConnections Response
		"PhysicalConnectionSet": map[string]interface{}{
			"PhysicalConnectionType": []interface{}{
				map[string]interface{}{
					"Status": "Enabled",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "EnablePhysicalConnection" {
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
		err := resourceAliCloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// TerminatePhysicalConnection
	attributesDiff = map[string]interface{}{
		"status": "Terminated",
	}
	diff, err = newInstanceDiff("alicloud_express_connect_physical_connection", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribePhysicalConnections Response
		"PhysicalConnectionSet": map[string]interface{}{
			"PhysicalConnectionType": []interface{}{
				map[string]interface{}{
					"Status": "Terminated",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "TerminatePhysicalConnection" {
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
		err := resourceAliCloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_express_connect_physical_connection"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribePhysicalConnections" {
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
		err := resourceAliCloudExpressConnectPhysicalConnectionRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudExpressConnectPhysicalConnectionDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeletePhysicalConnection" {
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
		err := resourceAliCloudExpressConnectPhysicalConnectionDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
