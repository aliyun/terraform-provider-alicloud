package alicloud

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_nlb_load_balancer",
		&resource.Sweeper{
			Name: "alicloud_nlb_load_balancer",
			F:    testSweepNlbLoadBalancer,
		})
}

func testSweepNlbLoadBalancer(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListLoadBalancers"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = aliyunClient.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.LoadBalancers", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.LoadBalancers", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["LoadBalancerName"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Nlb Load Balancer: %s", item["LoadBalancerName"].(string))
					continue
				}
			}
			action := "DeleteLoadBalancer"
			request := map[string]interface{}{
				"LoadBalancerId": item["LoadBalancerId"],
				"RegionId":       aliyunClient.RegionId,
			}
			_, err = aliyunClient.RpcPost("Nlb", "2022-04-30", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Nlb Load Balancer (%s): %s", item["LoadBalancerName"].(string), err)
			}
			log.Printf("[INFO] Delete Nlb Load Balancer success: %s ", item["LoadBalancerName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAliCloudNlbLoadBalancer_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	checkoutSupportedRegions(t, true, connectivity.NLBSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudNLBLoadBalancerMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNLBLoadBalancerBasicDependence0)
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
					"load_balancer_name":             "${var.name}",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"load_balancer_type":             "Network",
					"address_type":                   "Internet",
					"address_ip_version":             "Ipv4",
					"vpc_id":                         "${alicloud_vpc.default.id}",
					"deletion_protection_enabled":    "true",
					"modification_protection_status": "ConsoleProtection",
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_1}",
							"zone_id":    "${local.zone_id_1}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
					},
					"cross_zone_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":             name,
						"resource_group_id":              CHECKSET,
						"load_balancer_type":             "Network",
						"address_type":                   "Internet",
						"address_ip_version":             "Ipv4",
						"vpc_id":                         CHECKSET,
						"deletion_protection_enabled":    "true",
						"modification_protection_status": "ConsoleProtection",
						"tags.%":                         "2",
						"tags.Created":                   "tfTestAcc0",
						"tags.For":                       "Tftestacc 0",
						"zone_mappings.#":                "2",
						"cross_zone_enabled":             "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_reason":     "tf-open",
					"modification_protection_reason": "tf-open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_reason":     "tf-open",
						"modification_protection_reason": "tf-open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_reason":     "tf-open-update",
					"modification_protection_reason": "tf-open-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_reason":     "tf-open-update",
						"modification_protection_reason": "tf-open-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_enabled":    "false",
					"modification_protection_status": "NonProtection",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_enabled":    "false",
						"modification_protection_status": "NonProtection",
						"deletion_protection_reason":     "",
						"modification_protection_reason": "",
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

func TestAccAliCloudNlbLoadBalancer_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	checkoutSupportedRegions(t, true, connectivity.NLBSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudNLBLoadBalancerMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNLBLoadBalancerBasicDependence0)
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
					"load_balancer_name": "${var.name}",
					"address_type":       "Internet",
					"vpc_id":             "${alicloud_vpc.default.id}",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_1}",
							"zone_id":    "${local.zone_id_1}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
						"address_type":       "Internet",
						"vpc_id":             CHECKSET,
						"zone_mappings.#":    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "tfTestAcc0",
						"For":     "Tftestacc 0",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "tfTestAcc0",
						"tags.For":     "Tftestacc 0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_1}",
							"zone_id":    "${local.zone_id_1}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
						{
							"vswitch_id": "${local.vswitch_id_3}",
							"zone_id":    "${local.zone_id_3}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_type": "Intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type": "Intranet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_zone_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_zone_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": "${var.name}",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_1}",
							"zone_id":    "${local.zone_id_1}",
						},
						{
							"vswitch_id": "${local.vswitch_id_2}",
							"zone_id":    "${local.zone_id_2}",
						},
					},
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
						"zone_mappings.#":    "2",
						"cross_zone_enabled": "true",
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

var AlicloudNLBLoadBalancerMap0 = map[string]string{
	"cross_zone_enabled":            CHECKSET,
	"load_balancer_type":            CHECKSET,
	"status":                        CHECKSET,
	"address_ip_version":            CHECKSET,
	"load_balancer_name":            CHECKSET,
	"vpc_id":                        CHECKSET,
	"zone_mappings.#":               CHECKSET,
	"address_type":                  CHECKSET,
	"resource_group_id":             CHECKSET,
	"dns_name":                      CHECKSET,
	"load_balancer_business_status": CHECKSET,
	"tags.%":                        CHECKSET,
}

func AlicloudNLBLoadBalancerBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_nlb_zones" "default" {
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_vpc" "default" {
	  name       = var.name
	  cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
	  count      = length(data.alicloud_nlb_zones.default.zones)
	  vpc_id     = alicloud_vpc.default.id
	  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 3, count.index)
	  zone_id    = data.alicloud_nlb_zones.default.zones[count.index].id
	}

	locals {
  		zone_id_1    = data.alicloud_nlb_zones.default.zones.0.id
  		vswitch_id_1 = alicloud_vswitch.default.0.id
  		zone_id_2    = data.alicloud_nlb_zones.default.zones.1.id
  		vswitch_id_2 = alicloud_vswitch.default.1.id
  		zone_id_3    = data.alicloud_nlb_zones.default.zones.2.id
  		vswitch_id_3 = alicloud_vswitch.default.2.id
	}

	resource "alicloud_common_bandwidth_package" "default" {
  		bandwidth            = 2
  		internet_charge_type = "PayByBandwidth"
  		name                 = "${var.name}"
  		description          = "${var.name}_description"
	}
`, name)
}

var AlicloudNLBLoadBalancerMap1 = map[string]string{
	"address_type":       CHECKSET,
	"resource_group_id":  CHECKSET,
	"status":             CHECKSET,
	"cross_zone_enabled": CHECKSET,
	"load_balancer_type": CHECKSET,
	"vpc_id":             CHECKSET,
	"zone_mappings.#":    CHECKSET,
	"address_ip_version": CHECKSET,
	"load_balancer_name": CHECKSET,
}

func TestUnitAlicloudNlbLoadBalancer(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"load_balancer_name": "CreateLoadBalancerValue",
		"resource_group_id":  "CreateLoadBalancerValue",
		"load_balancer_type": "CreateLoadBalancerValue",
		"address_type":       "CreateLoadBalancerValue",
		"address_ip_version": "CreateLoadBalancerValue",
		"tags": map[string]string{
			"Created": "CreateLoadBalancerValue",
			"For":     "CreateLoadBalancerValue",
		},
		"vpc_id": "CreateLoadBalancerValue",
		"zone_mappings": []map[string]interface{}{
			{
				"vswitch_id":           "CreateLoadBalancerValue1",
				"zone_id":              "CreateLoadBalancerValue1",
				"allocation_id":        "CreateLoadBalancerValue1",
				"private_ipv4_address": "CreateLoadBalancerValue1",
			},
			{
				"vswitch_id":           "CreateLoadBalancerValue2",
				"zone_id":              "CreateLoadBalancerValue2",
				"allocation_id":        "CreateLoadBalancerValue2",
				"private_ipv4_address": "CreateLoadBalancerValue2",
			},
		},
		"cross_zone_enabled":   false,
		"bandwidth_package_id": "CreateLoadBalancerValue",
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
		// GetLoadBalancerAttribute
		"LoadBalancerName":           "CreateLoadBalancerValue",
		"CrossZoneEnabled":           false,
		"AddressIpVersion":           "CreateLoadBalancerValue",
		"AddressType":                "CreateLoadBalancerValue",
		"BandwidthPackageId":         "CreateLoadBalancerValue",
		"CreateTime":                 "DefaultValue",
		"DNSName":                    "DefaultValue",
		"Ipv6AddressType":            "DefaultValue",
		"LoadBalancerBusinessStatus": "DefaultValue",
		"LoadBalancerId":             "CreateLoadBalancerValue",
		"LoadBalancerType":           "CreateLoadBalancerValue",
		"ResourceGroupId":            "CreateLoadBalancerValue",
		"LoadBalancerStatus":         "Active",
		"VpcId":                      "CreateLoadBalancerValue",
		"ZoneMappings": []interface{}{
			map[string]interface{}{
				"EipType": "Common",
				"LoadBalancerAddresses": []interface{}{
					map[string]interface{}{
						"AllocationId":       "CreateLoadBalancerValue1",
						"PrivateIPv4Address": "CreateLoadBalancerValue1",
					},
				},
				"VSwitchId": "CreateLoadBalancerValue1",
				"ZoneId":    "CreateLoadBalancerValue1",
			},
			map[string]interface{}{
				"EipType": "Common",
				"LoadBalancerAddresses": []interface{}{
					map[string]interface{}{
						"AllocationId":       "CreateLoadBalancerValue2",
						"PrivateIPv4Address": "CreateLoadBalancerValue2",
					},
				},
				"VSwitchId": "CreateLoadBalancerValue2",
				"ZoneId":    "CreateLoadBalancerValue2",
			},
		},
		// ListTagResources
		"TagResources": []interface{}{
			map[string]interface{}{
				"TagKey":   "Created",
				"TagValue": "CreateLoadBalancerValue",
			},
			map[string]interface{}{
				"TagKey":   "For",
				"TagValue": "CreateLoadBalancerValue",
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateLoadBalancer
		"LoadbalancerId": "CreateLoadBalancerValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_nlb_load_balancer", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudNlbLoadBalancerCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateLoadBalancer" {
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
		err := resourceAliCloudNlbLoadBalancerCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudNlbLoadBalancerUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateLoadBalancerAttribute
	attributesDiff := map[string]interface{}{
		"cross_zone_enabled": true,
		"load_balancer_name": "UpdateLoadBalancerAttributeValue",
	}
	diff, err := newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetLoadBalancerAttribute Response
		"CrossZoneEnabled": true,
		"LoadBalancerName": "UpdateLoadBalancerAttributeValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateLoadBalancerAttribute" {
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
		err := resourceAliCloudNlbLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpdateLoadBalancerZones
	attributesDiff = map[string]interface{}{
		"zone_mappings": []interface{}{
			map[string]interface{}{
				"allocation_id":        "UpdateLoadBalancerZonesValue1",
				"private_ipv4_address": "UpdateLoadBalancerZonesValue1",
				"vswitch_id":           "UpdateLoadBalancerZonesValue1",
				"zone_id":              "UpdateLoadBalancerZonesValue1",
			},
			map[string]interface{}{
				"allocation_id":        "UpdateLoadBalancerZonesValue2",
				"private_ipv4_address": "UpdateLoadBalancerZonesValue2",
				"vswitch_id":           "UpdateLoadBalancerZonesValue2",
				"zone_id":              "UpdateLoadBalancerZonesValue2",
			},
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetLoadBalancerAttribute Response
		"ZoneMappings": []interface{}{
			map[string]interface{}{
				"EipType": "Common",
				"LoadBalancerAddresses": []interface{}{
					map[string]interface{}{
						"AllocationId":       "UpdateLoadBalancerZonesValue1",
						"PrivateIPv4Address": "UpdateLoadBalancerZonesValue1",
					},
				},
				"VSwitchId": "UpdateLoadBalancerZonesValue1",
				"ZoneId":    "UpdateLoadBalancerZonesValue1",
			},
			map[string]interface{}{
				"EipType": "Common",
				"LoadBalancerAddresses": []interface{}{
					map[string]interface{}{
						"AllocationId":       "UpdateLoadBalancerZonesValue2",
						"PrivateIPv4Address": "UpdateLoadBalancerZonesValue2",
					},
				},
				"VSwitchId": "UpdateLoadBalancerZonesValue2",
				"ZoneId":    "UpdateLoadBalancerZonesValue2",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateLoadBalancerZones" {
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
		err := resourceAliCloudNlbLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpdateLoadBalancerAddressTypeConfig
	attributesDiff = map[string]interface{}{
		"address_type": "UpdateLoadBalancerValue",
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// UpdateLoadBalancerAddressTypeConfig Response
		"AddressType": "UpdateLoadBalancerValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateLoadBalancerAddressTypeConfig" {
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
		err := resourceAliCloudNlbLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// AttachCommonBandwidthPackageToLoadBalancer
	attributesDiff = map[string]interface{}{
		"bandwidth_package_id": "UpdateLoadBalancerValue",
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// AttachCommonBandwidthPackageToLoadBalancer Response
		"BandwidthPackageId": "UpdateLoadBalancerValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AttachCommonBandwidthPackageToLoadBalancer" {
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
		err := resourceAliCloudNlbLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// DetachCommonBandwidthPackageFromLoadBalancer
	attributesDiff = map[string]interface{}{
		"bandwidth_package_id": "UpdateLoadBalancerValue1",
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DetachCommonBandwidthPackageFromLoadBalancer Response
		"BandwidthPackageId": "UpdateLoadBalancerValue1",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DetachCommonBandwidthPackageFromLoadBalancer" {
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
		err := resourceAliCloudNlbLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// TagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"TagResourcesValue_1": "TagResourcesValue_1",
			"TagResourcesValue_2": "TagResourcesValue_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListTagResources Response
		"TagResources": []interface{}{
			map[string]interface{}{
				"TagKey":   "TagResourcesValue_1",
				"TagValue": "TagResourcesValue_1",
			},
			map[string]interface{}{
				"TagKey":   "TagResourcesValue_2",
				"TagValue": "TagResourcesValue_2",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "TagResources" {
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
		err := resourceAliCloudNlbLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UntagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"UntagResourcesValue3_1": "UnTagResourcesValue3_1",
			"UntagResourcesValue3_2": "UnTagResourcesValue3_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListTagResources Response
		"TagResources": []interface{}{
			map[string]interface{}{
				"TagKey":   "UntagResourcesValue3_1",
				"TagValue": "UnTagResourcesValue3_1",
			},
			map[string]interface{}{
				"TagKey":   "UntagResourcesValue3_2",
				"TagValue": "UnTagResourcesValue3_2",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UntagResources" {
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
		err := resourceAliCloudNlbLoadBalancerUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetLoadBalancerAttribute" {
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
		err := resourceAliCloudNlbLoadBalancerRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudNlbLoadBalancerDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_nlb_load_balancer", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_nlb_load_balancer"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteLoadBalancer" {
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
			if *action == "GetLoadBalancerAttribute" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudNlbLoadBalancerDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test Nlb LoadBalancer. >>> Resource test cases, automatically generated.
// Case 3678
func TestAccAliCloudNlbLoadBalancer_basic3678(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbLoadBalancerMap3678)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbLoadBalancerBasicDependence3678)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.NLBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vsj.id}",
							"zone_id":    "${alicloud_vswitch.vsj.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsk.id}",
							"zone_id":    "${alicloud_vswitch.vsk.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsg.id}",
							"zone_id":    "${alicloud_vswitch.vsg.zone_id}",
						},
					},
					"address_type":       "Intranet",
					"vpc_id":             "${alicloud_vpc.vpc.id}",
					"load_balancer_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#":    "3",
						"address_type":       "Intranet",
						"vpc_id":             CHECKSET,
						"load_balancer_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_config": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_zone_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vsj.id}",
							"zone_id":    "${alicloud_vswitch.vsj.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsk.id}",
							"zone_id":    "${alicloud_vswitch.vsk.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsg.id}",
							"zone_id":    "${alicloud_vswitch.vsg.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_type": "Intranet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type": "Intranet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "ConsoleProtection",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_config": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_zone_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_config": []map[string]interface{}{
						{},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vsj.id}",
							"zone_id":    "${alicloud_vswitch.vsj.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsk.id}",
							"zone_id":    "${alicloud_vswitch.vsk.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsg.id}",
							"zone_id":    "${alicloud_vswitch.vsg.zone_id}",
						},
					},
					"address_type":       "Intranet",
					"address_ip_version": "Ipv4",
					"load_balancer_type": "Network",
					"vpc_id":             "${alicloud_vpc.vpc.id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"security_group_ids": []string{},
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
						},
					},
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name + "_update",
						"zone_mappings.#":      "3",
						"address_type":         "Intranet",
						"address_ip_version":   "Ipv4",
						"load_balancer_type":   "Network",
						"vpc_id":               CHECKSET,
						"resource_group_id":    CHECKSET,
						"security_group_ids.#": "0",
						"cross_zone_enabled":   "true",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudNlbLoadBalancerMap3678 = map[string]string{
	"load_balancer_type": "Network",
	"status":             CHECKSET,
	"create_time":        CHECKSET,
}

func AlicloudNlbLoadBalancerBasicDependence3678(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_nlb_zones" "default" {
}

resource "alicloud_vpc" "vpc" {
  vpc_name = var.name

  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vsj" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
  cidr_block   = "192.168.10.0/24"
  vswitch_name = var.name

}

resource "alicloud_vswitch" "vsk" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
  cidr_block   = "192.168.20.0/24"
  vswitch_name = var.name

}

resource "alicloud_security_group" "defaultLkkjal" {
  vpc_id              = alicloud_vpc.vpc.id
  name = var.name

}

resource "alicloud_security_group" "defaultmlAdy7" {
  vpc_id              = alicloud_vpc.vpc.id
  name = var.name

}

resource "alicloud_security_group" "defaultCr6BU3" {
  vpc_id              = alicloud_vpc.vpc.id
  name = var.name

}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_vswitch" "vsg" {
  vpc_id       = alicloud_vpc.vpc.id
  zone_id      = data.alicloud_nlb_zones.default.zones.2.id
  cidr_block   = "192.168.30.0/24"
  vswitch_name = var.name

}


`, name)
}

// Case 3862
func TestAccAliCloudNlbLoadBalancer_basic3862(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbLoadBalancerMap3862)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbLoadBalancerBasicDependence3862)
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
					"zone_mappings": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.defaultVSwitch.zone_id}",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultkR35um.id}",
							"zone_id":    "${alicloud_vswitch.defaultkR35um.zone_id}",
						},
					},
					"address_type":       "Internet",
					"vpc_id":             "${alicloud_vpc.defaultvVpc.id}",
					"load_balancer_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#":    "2",
						"address_type":       "Internet",
						"vpc_id":             CHECKSET,
						"load_balancer_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_zone_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_zone_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_mappings": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.defaultVSwitch.zone_id}",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultkR35um.id}",
							"zone_id":    "${alicloud_vswitch.defaultkR35um.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_mappings.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address_type": "Internet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address_type": "Internet",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
					"zone_mappings": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.defaultVSwitch.zone_id}",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultkR35um.id}",
							"zone_id":    "${alicloud_vswitch.defaultkR35um.zone_id}",
						},
					},
					"address_type":         "Internet",
					"cross_zone_enabled":   "true",
					"vpc_id":               "${alicloud_vpc.defaultvVpc.id}",
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.cbwp.id}",
					"load_balancer_type":   "Network",
					"address_ip_version":   "Ipv4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name + "_update",
						"zone_mappings.#":      "2",
						"address_type":         "Internet",
						"cross_zone_enabled":   "true",
						"vpc_id":               CHECKSET,
						"bandwidth_package_id": CHECKSET,
						"load_balancer_type":   "Network",
						"address_ip_version":   "Ipv4",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudNlbLoadBalancerMap3862 = map[string]string{
	"load_balancer_type": "Network",
	"status":             CHECKSET,
	"create_time":        CHECKSET,
}

func AlicloudNlbLoadBalancerBasicDependence3862(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_nlb_zones" "default" {
}

resource "alicloud_vpc" "defaultvVpc" {
  description = "test"
  cidr_block  = "10.0.0.0/8"
  enable_ipv6 = true
  vpc_name    = var.name

}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultvVpc.id
  cidr_block   = "10.0.1.0/24"
  zone_id      = data.alicloud_nlb_zones.default.zones.0.id
  vswitch_name = var.name

}

resource "alicloud_vswitch" "defaultkR35um" {
  description  = "test"
  zone_id      = data.alicloud_nlb_zones.default.zones.1.id
  vpc_id       = alicloud_vpc.defaultvVpc.id
  cidr_block   = "10.0.2.0/24"
  vswitch_name = var.name

  ipv6_cidr_block_mask = "8"
}

resource "alicloud_common_bandwidth_package" "cbwp" {
  bandwidth            = "1000"
  internet_charge_type = "PayByBandwidth"
}


`, name)
}

// Case 3678  twin
func TestAccAliCloudNlbLoadBalancer_basic3678_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbLoadBalancerMap3678)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbLoadBalancerBasicDependence3678)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.NLBSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name,
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.vsj.id}",
							"zone_id":    "${alicloud_vswitch.vsj.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsk.id}",
							"zone_id":    "${alicloud_vswitch.vsk.zone_id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.vsg.id}",
							"zone_id":    "${alicloud_vswitch.vsg.zone_id}",
						},
					},
					"address_type":       "Intranet",
					"address_ip_version": "Ipv4",
					"load_balancer_type": "Network",
					"vpc_id":             "${alicloud_vpc.vpc.id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"security_group_ids": []string{
						"${alicloud_security_group.defaultLkkjal.id}"},
					"deletion_protection_config": []map[string]interface{}{
						{
							"enabled": "false",
							"reason":  "",
						},
					},
					"modification_protection_config": []map[string]interface{}{
						{
							"status": "NonProtection",
							"reason": "",
						},
					},
					"cross_zone_enabled": "true",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name,
						"zone_mappings.#":      "3",
						"address_type":         "Intranet",
						"address_ip_version":   "Ipv4",
						"load_balancer_type":   "Network",
						"vpc_id":               CHECKSET,
						"resource_group_id":    CHECKSET,
						"security_group_ids.#": "1",
						"cross_zone_enabled":   "true",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case 3862  twin
func TestAccAliCloudNlbLoadBalancer_basic3862_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbLoadBalancerMap3862)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snlbloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbLoadBalancerBasicDependence3862)
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
					"load_balancer_name": name,
					"zone_mappings": []map[string]interface{}{
						{
							"zone_id":    "${alicloud_vswitch.defaultVSwitch.zone_id}",
							"vswitch_id": "${alicloud_vswitch.defaultVSwitch.id}",
						},
						{
							"vswitch_id": "${alicloud_vswitch.defaultkR35um.id}",
							"zone_id":    "${alicloud_vswitch.defaultkR35um.zone_id}",
						},
					},
					"address_type":         "Internet",
					"cross_zone_enabled":   "true",
					"vpc_id":               "${alicloud_vpc.defaultvVpc.id}",
					"bandwidth_package_id": "${alicloud_common_bandwidth_package.cbwp.id}",
					"load_balancer_type":   "Network",
					"address_ip_version":   "Ipv4",
					"ipv6_address_type":    "Intranet",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name":   name,
						"zone_mappings.#":      "2",
						"address_type":         "Internet",
						"cross_zone_enabled":   "true",
						"vpc_id":               CHECKSET,
						"bandwidth_package_id": CHECKSET,
						"load_balancer_type":   "Network",
						"address_ip_version":   "Ipv4",
						"ipv6_address_type":    "Intranet",
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Nlb LoadBalancer. <<< Resource test cases, automatically generated.
