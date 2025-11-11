package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ExpressConnect RouterInterface. >>> Resource test cases, automatically generated.
// Case RouterInterface-InitiatingSide-Subscription 11729
func TestAccAliCloudExpressConnectRouterInterface_basic11729(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectRouterInterfaceMap11729)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectRouterInterfaceBasicDependence11729)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fast_link_mode":     "false",
					"router_id":          "${alicloud_express_connect_virtual_border_router.default.id}",
					"opposite_region_id": "cn-beijing",
					"role":               "InitiatingSide",
					"router_type":        "VBR",
					"access_point_id":    "ap-cn-hangzhou-jg-B",
					"spec":               "Mini.2",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"auto_renew":         "false",
					"payment_type":       "Subscription",
					"period":             "1",
					"pricing_cycle":      "Month",
					"status":             "Idle",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fast_link_mode":     "false",
						"router_id":          CHECKSET,
						"opposite_region_id": "cn-beijing",
						"role":               "InitiatingSide",
						"router_type":        "VBR",
						"access_point_id":    "ap-cn-hangzhou-jg-B",
						"spec":               "Mini.2",
						"resource_group_id":  CHECKSET,
						"auto_renew":         "false",
						"payment_type":       "Subscription",
						"period":             "1",
						"pricing_cycle":      "Month",
						"status":             "Idle",
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
				ImportStateVerifyIgnore: []string{"auto_renew", "delete_health_check_ip", "period", "pricing_cycle"},
			},
		},
	})
}

var AlicloudExpressConnectRouterInterfaceMap11729 = map[string]string{}

func AlicloudExpressConnectRouterInterfaceBasicDependence11729(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
}
data "alicloud_account" "default" {
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING-JG"
}
resource "alicloud_express_connect_virtual_border_router" "default" {
  physical_connection_id = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  vlan_id                = "1001"
  peer_gateway_ip        = "192.168.254.2"
  peering_subnet_mask    = "255.255.255.0"
  local_gateway_ip       = "192.168.254.1"
}
`, name)
}

// Case RouterInterface-AcceptingSide 11648
func TestAccAliCloudExpressConnectRouterInterface_basic11648(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectRouterInterfaceMap11648)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectRouterInterfaceBasicDependence11648)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fast_link_mode":              "false",
					"opposite_router_id":          "${alicloud_express_connect_virtual_border_router.default.id}",
					"router_id":                   "${alicloud_vpc.default.router_id}",
					"opposite_router_type":        "VBR",
					"opposite_region_id":          "cn-hangzhou",
					"role":                        "AcceptingSide",
					"router_type":                 "VRouter",
					"spec":                        "Negative",
					"opposite_interface_owner_id": "${data.alicloud_account.this.id}",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"opposite_access_point_id":    "ap-cn-hangzhou-jg-B",
					"hc_threshold":                "7",
					"health_check_source_ip":      "172.16.0.1",
					"hc_rate":                     "2000",
					"health_check_target_ip":      "192.168.0.1",
					"status":                      "Idle",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fast_link_mode":              "false",
						"opposite_router_id":          CHECKSET,
						"router_id":                   CHECKSET,
						"opposite_router_type":        "VBR",
						"opposite_region_id":          "cn-hangzhou",
						"role":                        "AcceptingSide",
						"router_type":                 "VRouter",
						"spec":                        "Negative",
						"opposite_interface_owner_id": CHECKSET,
						"resource_group_id":           CHECKSET,
						"opposite_access_point_id":    "ap-cn-hangzhou-jg-B",
						"hc_threshold":                "7",
						"health_check_source_ip":      "172.16.0.1",
						"hc_rate":                     "2000",
						"health_check_target_ip":      "192.168.0.1",
						"status":                      "Idle",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hc_threshold":           "8",
					"health_check_source_ip": "172.16.0.2",
					"hc_rate":                "3000",
					"health_check_target_ip": "192.168.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hc_threshold":           "8",
						"health_check_source_ip": "172.16.0.2",
						"hc_rate":                "3000",
						"health_check_target_ip": "192.168.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_health_check_ip": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_health_check_ip": "true",
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
				ImportStateVerifyIgnore: []string{"auto_renew", "delete_health_check_ip", "period", "pricing_cycle"},
			},
		},
	})
}

var AlicloudExpressConnectRouterInterfaceMap11648 = map[string]string{}

func AlicloudExpressConnectRouterInterfaceBasicDependence11648(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
}
data "alicloud_account" "this" {
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING-JG"
}
data "alicloud_alb_zones" "default" {
}
resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  cidr_block  = "172.16.0.0/16"
  enable_ipv6 = "true"
}
resource "alicloud_vswitch" "zone_a" {
  vswitch_name         = var.name
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "172.16.0.0/24"
  zone_id              = data.alicloud_alb_zones.default.zones.0.id
  ipv6_cidr_block_mask = "6"
}
resource "alicloud_express_connect_virtual_border_router" "default" {
  physical_connection_id = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  vlan_id                = "1001"
  peer_gateway_ip        = "192.168.254.2"
  peering_subnet_mask    = "255.255.255.0"
  local_gateway_ip       = "192.168.254.1"
  depends_on = [alicloud_vswitch.zone_a]
}
`, name)
}

// Case RouterInterface-FastLinkMode 5930
func TestAccAliCloudExpressConnectRouterInterface_basic5930(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_express_connect_router_interface.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectRouterInterfaceMap5930)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressConnectServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeExpressConnectRouterInterface")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccexpressconnect%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectRouterInterfaceBasicDependence5930)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fast_link_mode":              "true",
					"opposite_router_id":          "${alicloud_vpc.default.router_id}",
					"router_id":                   "${alicloud_express_connect_virtual_border_router.default.id}",
					"opposite_router_type":        "VRouter",
					"opposite_region_id":          "cn-hangzhou",
					"role":                        "InitiatingSide",
					"router_type":                 "VBR",
					"access_point_id":             "ap-cn-hangzhou-jg-B",
					"spec":                        "Mini.2",
					"description":                 "terraform-example",
					"router_interface_name":       name,
					"opposite_interface_owner_id": "${data.alicloud_account.this.id}",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"auto_renew":                  "true",
					"payment_type":                "PayAsYouGo",
					"period":                      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fast_link_mode":              "true",
						"opposite_router_id":          CHECKSET,
						"router_id":                   CHECKSET,
						"opposite_router_type":        "VRouter",
						"opposite_region_id":          "cn-hangzhou",
						"role":                        "InitiatingSide",
						"router_type":                 "VBR",
						"access_point_id":             "ap-cn-hangzhou-jg-B",
						"spec":                        "Mini.2",
						"description":                 "terraform-example",
						"router_interface_name":       name,
						"opposite_interface_owner_id": CHECKSET,
						"resource_group_id":           CHECKSET,
						"auto_renew":                  "true",
						"payment_type":                "PayAsYouGo",
						"period":                      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":           "terraform-example-1",
					"router_interface_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           "terraform-example-1",
						"router_interface_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "Mini.5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "Mini.5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Active",
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
				Config: testAccConfig(map[string]interface{}{}),
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
				ImportStateVerifyIgnore: []string{"auto_renew", "delete_health_check_ip", "period", "pricing_cycle"},
			},
		},
	})
}

var AlicloudExpressConnectRouterInterfaceMap5930 = map[string]string{
	"opposite_interface_id": CHECKSET,
}

func AlicloudExpressConnectRouterInterfaceBasicDependence5930(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_resource_manager_resource_groups" "default" {
}
data "alicloud_account" "this" {
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING-JG"
}
data "alicloud_alb_zones" "default" {
}
resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  cidr_block  = "172.16.0.0/16"
  enable_ipv6 = "true"
}
resource "alicloud_vswitch" "zone_a" {
  vswitch_name         = var.name
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "172.16.0.0/24"
  zone_id              = data.alicloud_alb_zones.default.zones.0.id
  ipv6_cidr_block_mask = "6"
}
resource "alicloud_express_connect_virtual_border_router" "default" {
  physical_connection_id = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  vlan_id                = "1001"
  peer_gateway_ip        = "192.168.254.2"
  peering_subnet_mask    = "255.255.255.0"
  local_gateway_ip       = "192.168.254.1"
  depends_on = [alicloud_vswitch.zone_a]
}
`, name)
}

// Test ExpressConnect RouterInterface. <<< Resource test cases, automatically generated.
