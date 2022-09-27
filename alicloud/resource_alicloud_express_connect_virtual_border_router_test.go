package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_express_connect_virtual_border_router", &resource.Sweeper{
		Name: "alicloud_express_connect_virtual_border_router",
		F:    testSweepExpressConnectVirtualBorderRouters,
		Dependencies: []string{
			"alicloud_cen_instance",
		},
	})
}

func testSweepExpressConnectVirtualBorderRouters(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	var response interface{}
	for {
		action := "DescribeVirtualBorderRouters"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			break
		}
		resp, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		if err != nil {
			log.Printf("[ERROR] parsing %s response got an error: %s", action, err)
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			vbrName := fmt.Sprint(item["Name"])
			vbrId := fmt.Sprint(item["VbrId"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(vbrName), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping VirtualBorderRouter: %s (%s)", vbrName, vbrId)
				continue
			}
			action = "DeleteVirtualBorderRouter"
			request := map[string]interface{}{
				"VbrId":       vbrId,
				"RegionId":    client.RegionId,
				"ClientToken": buildClientToken("DeleteVirtualBorderRouter"),
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(1*time.Minute, func() *resource.RetryError {
				_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"DependencyViolation.BgpGroup"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] %s got an error: %v", action, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudExpressConnectVirtualBorderRouter_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_virtual_border_router.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectVirtualBorderRouterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectVirtualBorderRouter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 2999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectvirtualborderrouter%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectVirtualBorderRouterBasicDependence0)
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
					"physical_connection_id":     "${data.alicloud_express_connect_physical_connections.default.ids.0}",
					"vlan_id":                    fmt.Sprint(rand),
					"local_gateway_ip":           "10.0.0.1",
					"peer_gateway_ip":            "10.0.0.2",
					"peering_subnet_mask":        "255.255.255.252",
					"virtual_border_router_name": "tf-testAcc-PrT1AqAjKvGgLQpbygetjH6f",
					"description":                "tf-testAcc-llZJhorzazsS81mf2PVyFEAA",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_border_router_name": "tf-testAcc-PrT1AqAjKvGgLQpbygetjH6f",
						"description":                "tf-testAcc-llZJhorzazsS81mf2PVyFEAA",
						"physical_connection_id":     CHECKSET,
						"vlan_id":                    fmt.Sprint(rand),
						"local_gateway_ip":           "10.0.0.1",
						"peer_gateway_ip":            "10.0.0.2",
						"peering_subnet_mask":        "255.255.255.252",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_border_router_name": "tf-testAcc-1n8AGD0BcJcReSrQUAxTqaXC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_border_router_name": "tf-testAcc-1n8AGD0BcJcReSrQUAxTqaXC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"circuit_code": "tf-testAcc-m6VI39qqUEn76tiS06q862Jk",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"circuit_code": "tf-testAcc-m6VI39qqUEn76tiS06q862Jk",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-testAcc-ZwDPyqNDkTOoXueyCaUAL6Kj",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAcc-ZwDPyqNDkTOoXueyCaUAL6Kj",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detect_multiplier": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detect_multiplier": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_rx_interval": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_rx_interval": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_tx_interval": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_tx_interval": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "terminated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "terminated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peering_subnet_mask": "255.255.255.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peering_subnet_mask": "255.255.255.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"local_gateway_ip": "10.0.0.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_gateway_ip": "10.0.0.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_gateway_ip": "10.0.0.4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_gateway_ip": "10.0.0.4",
					}),
				),
			},
			// Currently, the product does not support ipv6
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"enable_ipv6": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"enable_ipv6": "true",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"local_ipv6_gateway_ip": "2001:4004:3c4d:0015:0000:0000:0000:1a2b",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"local_ipv6_gateway_ip": "2001:4004:3c4d:0015:0000:0000:0000:1a2b",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"peer_ipv6_gateway_ip": "2001:4004:3c4d:0015:0000:0000:0000:1a2b",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"peer_ipv6_gateway_ip": "2001:4004:3c4d:0015:0000:0000:0000:1a2b",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"peering_ipv6_subnet_mask": "2408:4004:cc:400::/56",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"peering_ipv6_subnet_mask": "2408:4004:cc:400::/56",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"vlan_id": fmt.Sprint(rand + 1),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vlan_id": fmt.Sprint(rand + 1),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_border_router_name": "tf-testAcc-MImKNETo3qwDBwnHVW3UUB8Y",
					"status":                     "active",
					"circuit_code":               "tf-testAcc-hM7XVPPmgiQkbPNgaQtqqGzX",
					"description":                "tf-testAcc-aoMEQnZ9PgEgzHjEV69O21rp",
					"detect_multiplier":          "10",
					"enable_ipv6":                "false",
					"min_rx_interval":            "300",
					"local_gateway_ip":           "192.168.0.11",
					"min_tx_interval":            "300",
					"peer_gateway_ip":            "192.168.0.12",
					"peering_subnet_mask":        "255.255.255.0",
					"vlan_id":                    fmt.Sprint(rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_border_router_name": "tf-testAcc-MImKNETo3qwDBwnHVW3UUB8Y",
						"status":                     "active",
						"circuit_code":               "tf-testAcc-hM7XVPPmgiQkbPNgaQtqqGzX",
						"description":                "tf-testAcc-aoMEQnZ9PgEgzHjEV69O21rp",
						"detect_multiplier":          "10",
						"enable_ipv6":                "false",
						"min_rx_interval":            "300",
						"local_gateway_ip":           "192.168.0.11",
						"min_tx_interval":            "300",
						"peer_gateway_ip":            "192.168.0.12",
						"peering_subnet_mask":        "255.255.255.0",
						"vlan_id":                    fmt.Sprint(rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"vbr_owner_id", "bandwidth"},
			},
		},
	})
}

var AlicloudExpressConnectVirtualBorderRouterMap0 = map[string]string{
	"enable_ipv6":  CHECKSET,
	"vbr_owner_id": NOSET,
	"bandwidth":    NOSET,
	"status":       "active",
}

func AlicloudExpressConnectVirtualBorderRouterBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}
`, name)
}

func TestUnitAlicloudExpressConnectVirtualBorderRouter(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_express_connect_virtual_border_router"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_express_connect_virtual_border_router"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"bandwidth":                  10,
		"circuit_code":               "CreateVirtualBorderRouterValue",
		"description":                "CreateVirtualBorderRouterValue",
		"enable_ipv6":                false,
		"local_gateway_ip":           "CreateVirtualBorderRouterValue",
		"local_ipv6_gateway_ip":      "CreateVirtualBorderRouterValue",
		"peer_gateway_ip":            "CreateVirtualBorderRouterValue",
		"peering_ipv6_subnet_mask":   "CreateVirtualBorderRouterValue",
		"peering_subnet_mask":        "CreateVirtualBorderRouterValue",
		"physical_connection_id":     "CreateVirtualBorderRouterValue",
		"vbr_owner_id":               "CreateVirtualBorderRouterValue",
		"virtual_border_router_name": "CreateVirtualBorderRouterValue",
		"vlan_id":                    1,
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
		// DescribeVirtualBorderRouters
		"VirtualBorderRouterSet": map[string]interface{}{
			"VirtualBorderRouterType": []interface{}{
				map[string]interface{}{
					"CircuitCode":           "CreateVirtualBorderRouterValue",
					"Description":           "CreateVirtualBorderRouterValue",
					"DetectMultiplier":      3,
					"EnableIpv6":            false,
					"LocalGatewayIp":        "CreateVirtualBorderRouterValue",
					"LocalIpv6GatewayIp":    "CreateVirtualBorderRouterValue",
					"MinRxInterval":         200,
					"MinTxInterval":         200,
					"PeerGatewayIp":         "CreateVirtualBorderRouterValue",
					"PeerIpv6GatewayIp":     "CreateVirtualBorderRouterValue",
					"PeeringIpv6SubnetMask": "CreateVirtualBorderRouterValue",
					"PeeringSubnetMask":     "CreateVirtualBorderRouterValue",
					"PhysicalConnectionId":  "CreateVirtualBorderRouterValue",
					"Status":                "CreateVirtualBorderRouterValue",
					"Name":                  "CreateVirtualBorderRouterValue",
					"VlanId":                1,
					"VbrId":                 "CreateVirtualBorderRouterValue",
				},
			},
		},
		"VbrId": "CreateVirtualBorderRouterValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateVirtualBorderRouter
		"VbrId": "CreateVirtualBorderRouterValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_express_connect_virtual_border_router", errorCode))
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
	err = resourceAlicloudExpressConnectVirtualBorderRouterCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeVirtualBorderRouters Response
		"VirtualBorderRouterSet": map[string]interface{}{
			"VirtualBorderRouterType": []interface{}{
				map[string]interface{}{
					"VbrId": "CreateVirtualBorderRouterValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateVirtualBorderRouter" {
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
		err := resourceAlicloudExpressConnectVirtualBorderRouterCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_express_connect_virtual_border_router"].Schema).Data(dInit.State(), nil)
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
	err = resourceAlicloudExpressConnectVirtualBorderRouterUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyVirtualBorderRouterAttribute
	attributesDiff := map[string]interface{}{
		"circuit_code":                    "ModifyVirtualBorderRouterAttributeValue",
		"description":                     "ModifyVirtualBorderRouterAttributeValue",
		"detect_multiplier":               5,
		"min_rx_interval":                 300,
		"min_tx_interval":                 300,
		"enable_ipv6":                     true,
		"local_gateway_ip":                "ModifyVirtualBorderRouterAttributeValue",
		"local_ipv6_gateway_ip":           "ModifyVirtualBorderRouterAttributeValue",
		"peer_gateway_ip":                 "ModifyVirtualBorderRouterAttributeValue",
		"peer_ipv6_gateway_ip":            "ModifyVirtualBorderRouterAttributeValue",
		"peering_ipv6_subnet_mask":        "ModifyVirtualBorderRouterAttributeValue",
		"peering_subnet_mask":             "ModifyVirtualBorderRouterAttributeValue",
		"virtual_border_router_name":      "ModifyVirtualBorderRouterAttributeValue",
		"vlan_id":                         2,
		"associated_physical_connections": "ModifyVirtualBorderRouterAttributeValue",
		"bandwidth":                       20,
	}
	diff, err := newInstanceDiff("alicloud_express_connect_virtual_border_router", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_express_connect_virtual_border_router"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeVirtualBorderRouters Response
		"VirtualBorderRouterSet": map[string]interface{}{
			"VirtualBorderRouterType": []interface{}{
				map[string]interface{}{
					"CircuitCode":           "ModifyVirtualBorderRouterAttributeValue",
					"Description":           "ModifyVirtualBorderRouterAttributeValue",
					"DetectMultiplier":      5,
					"EnableIpv6":            true,
					"LocalGatewayIp":        "ModifyVirtualBorderRouterAttributeValue",
					"LocalIpv6GatewayIp":    "ModifyVirtualBorderRouterAttributeValue",
					"MinRxInterval":         300,
					"MinTxInterval":         300,
					"PeerGatewayIp":         "ModifyVirtualBorderRouterAttributeValue",
					"PeerIpv6GatewayIp":     "ModifyVirtualBorderRouterAttributeValue",
					"PeeringIpv6SubnetMask": "ModifyVirtualBorderRouterAttributeValue",
					"PeeringSubnetMask":     "ModifyVirtualBorderRouterAttributeValue",
					"Status":                "ModifyVirtualBorderRouterAttributeValue",
					"Name":                  "ModifyVirtualBorderRouterAttributeValue",
					"VlanId":                2,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyVirtualBorderRouterAttribute" {
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
		err := resourceAlicloudExpressConnectVirtualBorderRouterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_express_connect_virtual_border_router"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// RecoverVirtualBorderRouter
	attributesDiff = map[string]interface{}{
		"status": "active",
	}
	diff, err = newInstanceDiff("alicloud_express_connect_virtual_border_router", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_express_connect_virtual_border_router"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeVirtualBorderRouters Response
		"VirtualBorderRouterSet": map[string]interface{}{
			"VirtualBorderRouterType": []interface{}{
				map[string]interface{}{
					"Status": "active",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RecoverVirtualBorderRouter" {
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
		err := resourceAlicloudExpressConnectVirtualBorderRouterUpdate(dExisted, rawClient)
		patches.Reset()

		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_express_connect_virtual_border_router"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// TerminateVirtualBorderRouter
	attributesDiff = map[string]interface{}{
		"status": "terminated",
	}
	diff, err = newInstanceDiff("alicloud_express_connect_virtual_border_router", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_express_connect_virtual_border_router"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeVirtualBorderRouters Response
		"VirtualBorderRouterSet": map[string]interface{}{
			"VirtualBorderRouterType": []interface{}{
				map[string]interface{}{
					"Status": "terminated",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "TerminateVirtualBorderRouter" {
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
		err := resourceAlicloudExpressConnectVirtualBorderRouterUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_express_connect_virtual_border_router"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribeVirtualBorderRouters" {
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
		err := resourceAlicloudExpressConnectVirtualBorderRouterRead(dExisted, rawClient)
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
	err = resourceAlicloudExpressConnectVirtualBorderRouterDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "DependencyViolation.BgpGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteVirtualBorderRouter" {
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
		err := resourceAlicloudExpressConnectVirtualBorderRouterDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
