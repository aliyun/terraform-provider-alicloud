package expressconnect_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnect"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
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
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1
	var response interface{}
	for {
		action := "DescribeVirtualBorderRouters"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
			if !sweep.SweepAll() {
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
			}
			action = "DeleteVirtualBorderRouter"
			request := map[string]interface{}{
				"VbrId":       vbrId,
				"RegionId":    client.RegionId,
				"ClientToken": helper.BuildClientToken("DeleteVirtualBorderRouter"),
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(1*time.Minute, func() *retry.RetryError {
				_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
				if err != nil {
					if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"DependencyViolation.BgpGroup"}) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] %s got an error: %v", action, err)
			}
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// lintignore: R001
func TestUnitAliCloudExpressConnectVirtualBorderRouter(t *testing.T) {
	p := acctest.Provider.ResourcesMap
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_express_connect_virtual_border_router", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewVpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = expressconnect.ResourceAliCloudExpressConnectVirtualBorderRouterCreate(dInit, rawClient)
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
		err := expressconnect.ResourceAliCloudExpressConnectVirtualBorderRouterCreate(dInit, rawClient)
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = expressconnect.ResourceAliCloudExpressConnectVirtualBorderRouterUpdate(dExisted, rawClient)
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
	diff, err := acctest.NewInstanceDiff("alicloud_express_connect_virtual_border_router", attributes, attributesDiff, dInit.State())
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
		err := expressconnect.ResourceAliCloudExpressConnectVirtualBorderRouterUpdate(dExisted, rawClient)
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
	diff, err = acctest.NewInstanceDiff("alicloud_express_connect_virtual_border_router", attributes, attributesDiff, dExisted.State())
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
		err := expressconnect.ResourceAliCloudExpressConnectVirtualBorderRouterUpdate(dExisted, rawClient)
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
	diff, err = acctest.NewInstanceDiff("alicloud_express_connect_virtual_border_router", attributes, attributesDiff, dExisted.State())
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
		err := expressconnect.ResourceAliCloudExpressConnectVirtualBorderRouterUpdate(dExisted, rawClient)
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
		err := expressconnect.ResourceAliCloudExpressConnectVirtualBorderRouterRead(dExisted, rawClient)
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = expressconnect.ResourceAliCloudExpressConnectVirtualBorderRouterDelete(dExisted, rawClient)
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
		err := expressconnect.ResourceAliCloudExpressConnectVirtualBorderRouterDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}

// Test ExpressConnect VirtualBorderRouter. >>> Resource test cases, automatically generated.
// Case VirtualBorderRouter_v3 11568
func TestAccAliCloudExpressConnectVirtualBorderRouter_basic11568(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_virtual_border_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudExpressConnectVirtualBorderRouterMap11568)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectVirtualBorderRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 2999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudExpressConnectVirtualBorderRouterBasicDependence11568)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"local_gateway_ip":       "10.0.0.1",
					"peer_gateway_ip":        "10.0.0.2",
					"peering_subnet_mask":    "255.255.255.252",
					"physical_connection_id": "${data.alicloud_express_connect_physical_connections.default.ids.0}",
					"vlan_id":                fmt.Sprint(rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_gateway_ip":       "10.0.0.1",
						"peer_gateway_ip":        "10.0.0.2",
						"peering_subnet_mask":    "255.255.255.252",
						"physical_connection_id": acctest.CHECKSET,
						"vlan_id":                fmt.Sprint(rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"circuit_code": "1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"circuit_code": "1024",
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
					"detect_multiplier": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detect_multiplier": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6":              "true",
					"peer_ipv6_gateway_ip":     "2408:4004:cc:400::2",
					"peering_ipv6_subnet_mask": "2408:4004:cc:400::/56",
					"local_ipv6_gateway_ip":    "2408:4004:cc:400::1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":              "true",
						"peer_ipv6_gateway_ip":     "2408:4004:cc:400::2",
						"peering_ipv6_subnet_mask": "2408:4004:cc:400::/56",
						"local_ipv6_gateway_ip":    "2408:4004:cc:400::1",
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
					"local_ipv6_gateway_ip": "2408:4004:cc:400::6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_ipv6_gateway_ip": "2408:4004:cc:400::6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_rx_interval": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_rx_interval": "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_tx_interval": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_tx_interval": "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mtu": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mtu": "1500",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_ipv6_gateway_ip": "2408:4004:cc:400::8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_ipv6_gateway_ip": "2408:4004:cc:400::8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peering_ipv6_subnet_mask": "2408:4004:cc:400::/55",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peering_ipv6_subnet_mask": "2408:4004:cc:400::/55",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sitelink_enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sitelink_enable": "true",
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
					"virtual_border_router_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_border_router_name": name,
					}),
				),
			},
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vbr_owner_id"},
			},
		},
	})
}

func TestAccAliCloudExpressConnectVirtualBorderRouter_basic11568_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_virtual_border_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudExpressConnectVirtualBorderRouterMap11568)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectVirtualBorderRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 2999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudExpressConnectVirtualBorderRouterBasicDependence11568)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"local_gateway_ip":           "10.0.0.1",
					"peer_gateway_ip":            "10.0.0.2",
					"peering_subnet_mask":        "255.255.255.252",
					"physical_connection_id":     "${data.alicloud_express_connect_physical_connections.default.ids.0}",
					"vlan_id":                    fmt.Sprint(rand),
					"bandwidth":                  "100",
					"circuit_code":               "1024",
					"description":                name,
					"detect_multiplier":          "6",
					"enable_ipv6":                "true",
					"peer_ipv6_gateway_ip":       "2408:4004:cc:400::2",
					"peering_ipv6_subnet_mask":   "2408:4004:cc:400::/56",
					"local_ipv6_gateway_ip":      "2408:4004:cc:400::1",
					"min_rx_interval":            "600",
					"min_tx_interval":            "600",
					"mtu":                        "1500",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"sitelink_enable":            "true",
					"status":                     "active",
					"virtual_border_router_name": name,
					"vbr_owner_id":               "${data.alicloud_account.default.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_gateway_ip":           "10.0.0.1",
						"peer_gateway_ip":            "10.0.0.2",
						"peering_subnet_mask":        "255.255.255.252",
						"physical_connection_id":     acctest.CHECKSET,
						"vlan_id":                    fmt.Sprint(rand),
						"bandwidth":                  "100",
						"circuit_code":               "1024",
						"description":                name,
						"detect_multiplier":          "6",
						"enable_ipv6":                "true",
						"peer_ipv6_gateway_ip":       "2408:4004:cc:400::2",
						"peering_ipv6_subnet_mask":   "2408:4004:cc:400::/56",
						"local_ipv6_gateway_ip":      "2408:4004:cc:400::1",
						"min_rx_interval":            "600",
						"min_tx_interval":            "600",
						"mtu":                        "1500",
						"resource_group_id":          acctest.CHECKSET,
						"sitelink_enable":            "true",
						"status":                     "active",
						"virtual_border_router_name": name,
						"tags.%":                     "2",
						"tags.Created":               "TF",
						"tags.For":                   "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vbr_owner_id"},
			},
		},
	})
}

var AliCloudExpressConnectVirtualBorderRouterMap11568 = map[string]string{
	"detect_multiplier": acctest.CHECKSET,
	"min_rx_interval":   acctest.CHECKSET,
	"min_tx_interval":   acctest.CHECKSET,
	"mtu":               acctest.CHECKSET,
	"resource_group_id": acctest.CHECKSET,
	"status":            acctest.CHECKSET,
	"route_table_id":    acctest.CHECKSET,
	"create_time":       acctest.CHECKSET,
}

func AliCloudExpressConnectVirtualBorderRouterBasicDependence11568(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "default" {
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}
`, name)
}

// Case VirtualBorderRouter_v3 11569
func TestAccAliCloudExpressConnectVirtualBorderRouter_basic11569(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_virtual_border_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudExpressConnectVirtualBorderRouterMap11569)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectVirtualBorderRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 2999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudExpressConnectVirtualBorderRouterBasicDependence11569)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"local_gateway_ip":       "10.0.0.1",
					"peer_gateway_ip":        "10.0.0.2",
					"peering_subnet_mask":    "255.255.255.252",
					"physical_connection_id": "${data.alicloud_express_connect_physical_connections.default.ids.0}",
					"vlan_id":                fmt.Sprint(rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_gateway_ip":       "10.0.0.1",
						"peer_gateway_ip":        "10.0.0.2",
						"peering_subnet_mask":    "255.255.255.252",
						"physical_connection_id": acctest.CHECKSET,
						"vlan_id":                fmt.Sprint(rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"circuit_code": "1024",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"circuit_code": "1024",
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
					"detect_multiplier": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detect_multiplier": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6":              "true",
					"peer_ipv6_gateway_ip":     "2408:4004:cc:400::2",
					"peering_ipv6_subnet_mask": "2408:4004:cc:400::/56",
					"local_ipv6_gateway_ip":    "2408:4004:cc:400::1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6":              "true",
						"peer_ipv6_gateway_ip":     "2408:4004:cc:400::2",
						"peering_ipv6_subnet_mask": "2408:4004:cc:400::/56",
						"local_ipv6_gateway_ip":    "2408:4004:cc:400::1",
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
					"local_ipv6_gateway_ip": "2408:4004:cc:400::6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_ipv6_gateway_ip": "2408:4004:cc:400::6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_rx_interval": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_rx_interval": "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_tx_interval": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_tx_interval": "600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mtu": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mtu": "1500",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_ipv6_gateway_ip": "2408:4004:cc:400::8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_ipv6_gateway_ip": "2408:4004:cc:400::8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peering_ipv6_subnet_mask": "2408:4004:cc:400::/55",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peering_ipv6_subnet_mask": "2408:4004:cc:400::/55",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sitelink_enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sitelink_enable": "true",
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
					"virtual_border_router_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_border_router_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vbr_owner_id"},
			},
		},
	})
}

func TestAccAliCloudExpressConnectVirtualBorderRouter_basic11569_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_virtual_border_router.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudExpressConnectVirtualBorderRouterMap11569)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return expressconnect.NewExpressConnectServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeExpressConnectVirtualBorderRouter")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 2999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudExpressConnectVirtualBorderRouterBasicDependence11569)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"local_gateway_ip":           "10.0.0.1",
					"peer_gateway_ip":            "10.0.0.2",
					"peering_subnet_mask":        "255.255.255.252",
					"physical_connection_id":     "${data.alicloud_express_connect_physical_connections.default.ids.0}",
					"vlan_id":                    fmt.Sprint(rand),
					"bandwidth":                  "100",
					"circuit_code":               "1024",
					"description":                name,
					"detect_multiplier":          "6",
					"enable_ipv6":                "true",
					"peer_ipv6_gateway_ip":       "2408:4004:cc:400::2",
					"peering_ipv6_subnet_mask":   "2408:4004:cc:400::/56",
					"local_ipv6_gateway_ip":      "2408:4004:cc:400::1",
					"min_rx_interval":            "600",
					"min_tx_interval":            "600",
					"mtu":                        "1500",
					"resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"sitelink_enable":            "true",
					"status":                     "active",
					"virtual_border_router_name": name,
					"vbr_owner_id":               "${data.alicloud_account.default.id}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_gateway_ip":           "10.0.0.1",
						"peer_gateway_ip":            "10.0.0.2",
						"peering_subnet_mask":        "255.255.255.252",
						"physical_connection_id":     acctest.CHECKSET,
						"vlan_id":                    acctest.CHECKSET,
						"bandwidth":                  acctest.CHECKSET,
						"circuit_code":               "1024",
						"description":                name,
						"detect_multiplier":          "6",
						"enable_ipv6":                "true",
						"peer_ipv6_gateway_ip":       "2408:4004:cc:400::2",
						"peering_ipv6_subnet_mask":   "2408:4004:cc:400::/56",
						"local_ipv6_gateway_ip":      "2408:4004:cc:400::1",
						"min_rx_interval":            "600",
						"min_tx_interval":            "600",
						"mtu":                        "1500",
						"resource_group_id":          acctest.CHECKSET,
						"sitelink_enable":            "true",
						"status":                     "active",
						"virtual_border_router_name": name,
						"tags.%":                     "2",
						"tags.Created":               "TF",
						"tags.For":                   "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vbr_owner_id"},
			},
		},
	})
}

var AliCloudExpressConnectVirtualBorderRouterMap11569 = map[string]string{
	"detect_multiplier": acctest.CHECKSET,
	"min_rx_interval":   acctest.CHECKSET,
	"min_tx_interval":   acctest.CHECKSET,
	"mtu":               acctest.CHECKSET,
	"resource_group_id": acctest.CHECKSET,
	"status":            acctest.CHECKSET,
	"route_table_id":    acctest.CHECKSET,
	"create_time":       acctest.CHECKSET,
}

func AliCloudExpressConnectVirtualBorderRouterBasicDependence11569(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "default" {
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}
`, name)
}

// Test ExpressConnect VirtualBorderRouter. <<< Resource test cases, automatically generated.
