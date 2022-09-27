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

func TestAccAlicloudExpressConnectPhysicalConnection_domesic(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectPhysicalConnectionBasicDependence0)
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
					"access_point_id":          getAccessPointId(),
					"type":                     "VPC",
					"peer_location":            "testacc12345",
					"physical_connection_name": "${var.name}",
					"description":              "${var.name}",
					"line_operator":            "CU",
					"port_type":                "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id":          CHECKSET,
						"type":                     "VPC",
						"peer_location":            "testacc12345",
						"physical_connection_name": name,
						"description":              name,
						"line_operator":            "CU",
						"port_type":                "1000Base-LX",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name + "_Update",
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
					"description": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"line_operator": "CU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"line_operator": "CU",
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
					"line_operator": "CO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"line_operator": "CO",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_location": "浙江省---vfjdbg_21e",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_location": "浙江省---vfjdbg_21e",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"port_type": "10GBase-LR",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"port_type": "10GBase-LR",
			//		}),
			//	),
			//},
			// Only confirmed connection can be enabled.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "Enabled",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "Enabled",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name,
					"status":                   "Canceled",
					"bandwidth":                "15",
					"circuit_code":             "longtel002",
					"description":              name,
					"line_operator":            "CT",
					"peer_location":            "testacc12345",
					"port_type":                "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name,
						"status":                   "Canceled",
						"bandwidth":                "15",
						"circuit_code":             "longtel002",
						"description":              name,
						"line_operator":            "CT",
						"peer_location":            "testacc12345",
						"port_type":                "1000Base-LX",
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

func TestAccAlicloudExpressConnectPhysicalConnection_intl(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectPhysicalConnectionBasicDependence0)
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
					//"access_point_id": "${data.alicloud_express_connect_access_points.default.ids.0}",
					"access_point_id":          getAccessPointId(),
					"type":                     "VPC",
					"peer_location":            "testacc12345",
					"physical_connection_name": "${var.name}",
					"description":              "${var.name}",
					"line_operator":            "Other",
					"port_type":                "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id":          CHECKSET,
						"type":                     "VPC",
						"peer_location":            "testacc12345",
						"physical_connection_name": name,
						"description":              name,
						"line_operator":            "Other",
						"port_type":                "1000Base-LX",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name + "_Update",
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
					"description": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_Update",
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
					"peer_location": "国际---vfjdbg_21e",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_location": "国际---vfjdbg_21e",
					}),
				),
			},
			// Currently, the internal region does not support 10G
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"port_type": "10GBase-LR",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"port_type": "10GBase-LR",
			//		}),
			//	),
			//},
			// Only confirmed connection can be enabled.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "Enabled",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "Enabled",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name,
					"status":                   "Canceled",
					"bandwidth":                "15",
					"circuit_code":             "longtel002",
					"description":              name,
					"line_operator":            "Other",
					"peer_location":            "testacc12345",
					"port_type":                "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name,
						"status":                   "Canceled",
						"bandwidth":                "15",
						"circuit_code":             "longtel002",
						"description":              name,
						"line_operator":            "Other",
						"peer_location":            "testacc12345",
						"port_type":                "1000Base-LX",
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

func TestAccAlicloudExpressConnectPhysicalConnection_domesic1(t *testing.T) {
	t.Skipf("There is an api bug that its describe response does not return CircuitCode. If the bug fixed, reopen this case")
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectPhysicalConnectionBasicDependence1)
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
					"redundant_physical_connection_id": "${data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id}",
					"type":                             "VPC",
					"peer_location":                    "testacc12345",
					"physical_connection_name":         name,
					"description":                      "${var.name}",
					"line_operator":                    "CU",
					"port_type":                        "10GBase-LR",
					"bandwidth":                        "10",
					"circuit_code":                     "longtel001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id":                  CHECKSET,
						"redundant_physical_connection_id": CHECKSET,
						"type":                             "VPC",
						"peer_location":                    "testacc12345",
						"physical_connection_name":         name,
						"description":                      name,
						"line_operator":                    "CU",
						"port_type":                        "10GBase-LR",
						"bandwidth":                        "10",
						"circuit_code":                     "longtel001",
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

var AlicloudExpressConnectPhysicalConnectionMap0 = map[string]string{
	"status":                           CHECKSET,
	"redundant_physical_connection_id": "",
	"bandwidth":                        CHECKSET,
}

func AlicloudExpressConnectPhysicalConnectionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_express_connect_access_points" "default" {
	status = "recommended"
}
`, name)
}

func AlicloudExpressConnectPhysicalConnectionBasicDependence1(name string) string {
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

func TestUnitAlicloudExpressConnectPhysicalConnection(t *testing.T) {
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
	err = resourceAlicloudExpressConnectPhysicalConnectionCreate(dInit, rawClient)
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
		err := resourceAlicloudExpressConnectPhysicalConnectionCreate(dInit, rawClient)
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
	err = resourceAlicloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
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
		err := resourceAlicloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
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
		err := resourceAlicloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
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
		err := resourceAlicloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
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
		err := resourceAlicloudExpressConnectPhysicalConnectionUpdate(dExisted, rawClient)
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
		err := resourceAlicloudExpressConnectPhysicalConnectionRead(dExisted, rawClient)
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
	err = resourceAlicloudExpressConnectPhysicalConnectionDelete(dExisted, rawClient)
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
		err := resourceAlicloudExpressConnectPhysicalConnectionDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
