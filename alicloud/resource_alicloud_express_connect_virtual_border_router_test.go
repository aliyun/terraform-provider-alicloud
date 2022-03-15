package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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
	checkoutAccount(t, true)
	defer checkoutAccount(t, false)
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_virtual_border_router.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectVirtualBorderRouterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectVirtualBorderRouter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectvirtualborderrouter%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectVirtualBorderRouterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VbrSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_id":     "${data.alicloud_express_connect_physical_connections.default.ids.0}",
					"vlan_id":                    "1000",
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
						"vlan_id":                    "1000",
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
					"vlan_id": "290",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vlan_id": "290",
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
					"vlan_id":                    "639",
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
						"vlan_id":                    "639",
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
