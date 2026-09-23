package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_vpc_bgp_group",
		&resource.Sweeper{
			Name: "alicloud_vpc_bgp_group",
			F:    testSweepVpcBgpGroup,
		})
}

func testSweepVpcBgpGroup(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.VPCBgpGroupSupportRegions) {
		log.Printf("[INFO] Skipping Vpc Bgp Group unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		log.Printf("[ERROR] getting Alicloud client: %s", err)
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeBgpGroups"
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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

		resp, err := jsonpath.Get("$.BgpGroups.BgpGroup", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.BgpGroups.BgpGroup", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Name"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Vpc Bgp Group: %s", item["Name"].(string))
				continue
			}
			action := "DeleteBgpGroup"
			request := map[string]interface{}{
				"BgpGroupId": item["BgpGroupId"],
				"RegionId":   client.RegionId,
			}
			_, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Vpc Bgp Group (%s): %s", item["Name"].(string), err)
			}
			log.Printf("[INFO] Delete Vpc Bgp Group success: %s ", item["Name"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudVPCBgpGroup_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_vpc_bgp_group.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCBgpGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcBgpGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcbgpgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCBgpGroupBasicDependence0)
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
					"router_id":      "${alicloud_express_connect_virtual_border_router.default.id}",
					"peer_asn":       "1111",
					"bgp_group_name": "${var.name}",
					"route_limit":    "110",
					"clear_auth_key": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"router_id":      CHECKSET,
						"peer_asn":       "1111",
						"bgp_group_name": name,
						"route_limit":    "110",
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
					"bgp_group_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bgp_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_limit": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_limit": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_asn": "1112",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_asn": "1112",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_asn": "64513",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_asn": "64513",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_key": "YourPassword+123456789",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_key": "YourPassword+123456789",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "${var.name}",
					"peer_asn":       "1111",
					"local_asn":      "64512",
					"bgp_group_name": "${var.name}",
					"auth_key":       "YourPassword+12345678",
					"clear_auth_key": "false",
					"is_fake_asn":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    name,
						"peer_asn":       "1111",
						"local_asn":      "64512",
						"bgp_group_name": name,
						"auth_key":       "YourPassword+12345678",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"clear_auth_key"},
			},
		},
	})
}
func TestAccAliCloudVPCBgpGroup_basic1(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_vpc_bgp_group.default"
	ra := resourceAttrInit(resourceId, AlicloudVPCBgpGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcBgpGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcbgpgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVPCBgpGroupBasicDependence0)
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
					"description":    "${var.name}",
					"router_id":      "${alicloud_express_connect_virtual_border_router.default.id}",
					"peer_asn":       "1111",
					"bgp_group_name": "${var.name}",
					"local_asn":      "64512",
					"auth_key":       "YourPassword+12345678",
					"is_fake_asn":    "true",
					"ip_version":     "IPv4",
					"clear_auth_key": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    name,
						"router_id":      CHECKSET,
						"peer_asn":       "1111",
						"bgp_group_name": name,
						"local_asn":      "64512",
						"auth_key":       "YourPassword+12345678",
						"is_fake_asn":    "true",
						"ip_version":     "IPv4",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"clear_auth_key"},
			},
		},
	})
}

var AlicloudVPCBgpGroupMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudVPCBgpGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_express_connect_physical_connections" "default" {
	name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.default.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
`, name, acctest.RandIntRange(1, 2999))
}
