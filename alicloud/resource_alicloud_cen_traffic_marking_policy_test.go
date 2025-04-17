package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCenTrafficMarkingPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_traffic_marking_policy.default"
	checkoutSupportedRegions(t, true, connectivity.CenTRSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCENTrafficMarkingPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTrafficMarkingPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentrafficmarkingpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCENTrafficMarkingPolicyBasicDependence0)
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
					"priority":                    "5",
					"description":                 "${var.name}",
					"traffic_marking_policy_name": "${var.name}",
					"marking_dscp":                "5",
					"transit_router_id":           "${alicloud_cen_transit_router.default.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority":                    "5",
						"description":                 name,
						"traffic_marking_policy_name": name,
						"marking_dscp":                "5",
						"transit_router_id":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_marking_policy_name": "${var.name}_update",
					"force":                       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_marking_policy_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "force"},
			},
		},
	})
}

var AlicloudCENTrafficMarkingPolicyMap0 = map[string]string{
	"dry_run": NOSET,
	"status":  CHECKSET,
}

func AlicloudCENTrafficMarkingPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  cen_id              = alicloud_cen_instance.default.id
  transit_router_name = var.name
}
`, name)
}

func TestUnitAccAlicloudCENTrafficMarkingPolicy(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cen_traffic_marking_policy"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cen_traffic_marking_policy"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"priority":                    5,
		"dry_run":                     true,
		"description":                 "CreateCenTrafficMarkingPolicyValue",
		"traffic_marking_policy_name": "CreateCenTrafficMarkingPolicyValue",
		"marking_dscp":                5,
		"transit_router_id":           "CreateCenTrafficMarkingPolicyValue",
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
		"TrafficMarkingPolicies": []interface{}{
			map[string]interface{}{
				"TrafficMarkingPolicyStatus":      "Active",
				"TrafficMarkingPolicyId":          "CreateCenTrafficMarkingPolicyValue",
				"MarkingDscp":                     5,
				"TrafficMarkingPolicyName":        "CreateCenTrafficMarkingPolicyValue",
				"Priority":                        5,
				"TrafficMarkingPolicyDescription": "CreateCenTrafficMarkingPolicyValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"TrafficMarkingPolicyId": "CreateCenTrafficMarkingPolicyValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cen_traffic_marking_policy", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCenTrafficMarkingPolicyCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateTrafficMarkingPolicy" {
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
		err := resourceAliCloudCenTrafficMarkingPolicyCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cen_traffic_marking_policy"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCenTrafficMarkingPolicyUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"description":                 "UpdateCenTrafficMarkingPolicyValue",
		"traffic_marking_policy_name": "UpdateCenTrafficMarkingPolicyValue",
	}
	diff, err := newInstanceDiff("alicloud_cen_traffic_marking_policy", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cen_traffic_marking_policy"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"TrafficMarkingPolicies": []interface{}{
			map[string]interface{}{
				"TrafficMarkingPolicyName":        "UpdateCenTrafficMarkingPolicyValue",
				"TrafficMarkingPolicyDescription": "UpdateCenTrafficMarkingPolicyValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateTrafficMarkingPolicyAttribute" {
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
		err := resourceAliCloudCenTrafficMarkingPolicyUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_cen_traffic_marking_policy"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_cen_traffic_marking_policy", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cen_traffic_marking_policy"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListTrafficMarkingPolicies" {
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
		err := resourceAliCloudCenTrafficMarkingPolicyRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCbnClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudCenTrafficMarkingPolicyDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_cen_traffic_marking_policy", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_cen_traffic_marking_policy"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteTrafficMarkingPolicy" {
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
			if *action == "ListTrafficMarkingPolicies" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudCenTrafficMarkingPolicyDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test Cen TrafficMarkingPolicy. >>> Resource test cases, automatically generated.
// Case TR支持IPv6-善问-线上 7854
func TestAccAliCloudCenTrafficMarkingPolicy_basic7854(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_traffic_marking_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTrafficMarkingPolicyMap7854)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTrafficMarkingPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentrafficmarkingpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTrafficMarkingPolicyBasicDependence7854)
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
					"traffic_marking_policy_name": name,
					"marking_dscp":                "11",
					"priority":                    "11",
					"transit_router_id":           "${alicloud_cen_transit_router.default8JJJSl.transit_router_id}",
					"traffic_match_rules": []map[string]interface{}{
						{
							"match_dscp":                     "11",
							"dst_cidr":                       "::/0",
							"traffic_match_rule_description": "xxx",
							"protocol":                       "UDP",
							"src_cidr":                       "::/0",
							"traffic_match_rule_name":        "ttt",
							"address_family":                 "IPv6",
							"dst_port_range": []string{
								"1", "2"},
							"src_port_range": []string{
								"1", "2"},
						},
						{
							"match_dscp":                     "13",
							"dst_cidr":                       "192.169.6.6/32",
							"traffic_match_rule_description": "xxxhg",
							"protocol":                       "UDP",
							"src_cidr":                       "192.166.3.3/32",
							"traffic_match_rule_name":        "ttt信息发达",
							"address_family":                 "IPv4",
							"dst_port_range": []string{
								"1", "105"},
							"src_port_range": []string{
								"1", "101"},
						},
						{
							"match_dscp":                     "15",
							"traffic_match_rule_description": "765432",
							"traffic_match_rule_name":        "09876",
							"dst_port_range": []string{
								"1", "105"},
							"src_port_range": []string{
								"1", "101"},
						},
					},
					"description": "ttt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_marking_policy_name": name,
						"marking_dscp":                "11",
						"priority":                    "11",
						"transit_router_id":           CHECKSET,
						"traffic_match_rules.#":       "3",
						"description":                 "ttt",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "xxx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "xxx",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_marking_policy_name": name + "_update",
					"traffic_match_rules": []map[string]interface{}{
						{
							"address_family":                 "IPv4",
							"match_dscp":                     "12",
							"dst_cidr":                       "10.0.4.0/24",
							"traffic_match_rule_description": "xxxxxtt",
							"protocol":                       "SSH",
							"traffic_match_rule_name":        "gdafdax",
							"src_cidr":                       "10.0.2.0/24",
							"dst_port_range": []string{
								"22", "22"},
							"src_port_range": []string{
								"22", "22"},
						},
						{
							"address_family":                 "IPv6",
							"match_dscp":                     "18",
							"traffic_match_rule_description": "zzxxxxxtt",
							"protocol":                       "SSH",
							"traffic_match_rule_name":        "gdafdaxd",
							"src_cidr":                       "::/0",
							"dst_port_range": []string{
								"22", "22"},
							"src_port_range": []string{
								"22", "22"},
						},
						{
							"match_dscp":                     "55",
							"traffic_match_rule_description": "xx",
							"traffic_match_rule_name":        "yy",
							"dst_port_range": []string{
								"1", "101"},
							"src_port_range": []string{
								"1", "105"},
						},
					},
					"description": "hhh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_marking_policy_name": name + "_update",
						"traffic_match_rules.#":       "3",
						"description":                 "hhh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_marking_policy_name": name + "_update",
					"description":                 "oooo",
					"traffic_match_rules":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_marking_policy_name": name + "_update",
						"description":                 "oooo",
						"traffic_match_rules":         NOSET,
						"traffic_match_rules.#":       "0",
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

var AlicloudCenTrafficMarkingPolicyMap7854 = map[string]string{
	"status":                    CHECKSET,
	"traffic_marking_policy_id": CHECKSET,
}

func AlicloudCenTrafficMarkingPolicyBasicDependence7854(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cen_instance" "defaultcIz05m" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default8JJJSl" {
  support_multicast          = true
  cen_id                     = alicloud_cen_instance.defaultcIz05m.id
  transit_router_name        = format("%%s1", var.name)
  transit_router_description = "tr"
}


`, name)
}

// Test Cen TrafficMarkingPolicy. <<< Resource test cases, automatically generated.
