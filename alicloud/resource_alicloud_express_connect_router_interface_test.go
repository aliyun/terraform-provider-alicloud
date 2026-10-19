package alicloud

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

func TestAccAliCloudExpressConnectRouterInterface_crossRegionDelete(t *testing.T) {
	resourceId := "alicloud_express_connect_router_interface.default"
	name := fmt.Sprintf("tf-testAccExpressConnect%d", acctest.RandIntRange(10000, 99999))
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, expressConnectRouterInterfaceCrossRegionDeleteDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		ProviderFactories: testAccProviderFactoriesAlternate(),
		CheckDestroy:      testAccCheckCrossRegionRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"fast_link_mode":              "false",
					"payment_type":                "PayAsYouGo",
					"opposite_region_id":          "cn-beijing",
					"opposite_router_type":        "VRouter",
					"opposite_router_id":          "${alicloud_vpc.remote.router_id}",
					"opposite_interface_owner_id": "${data.alicloud_account.cross_region.id}",
					"router_type":                 "VRouter",
					"router_id":                   "${alicloud_vpc.local.router_id}",
					"role":                        "InitiatingSide",
					"spec":                        "Mini.2",
					"router_interface_name":       name,
					"description":                 "cross region deletion test",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "opposite_region_id", "cn-beijing"),
					resource.TestCheckResourceAttr("alicloud_express_connect_router_interface.remote", "opposite_region_id", "cn-hangzhou"),
					testAccCheckCrossRegionRouterInterfacesActive,
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"router_interface_name": name + "-updated",
					"description":           "updated cross region deletion test",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "router_interface_name", name+"-updated"),
					testAccCheckCrossRegionRouterInterfacesActive,
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

func expressConnectRouterInterfaceCrossRegionDeleteDependence(name string) string {
	return fmt.Sprintf(`
data "alicloud_account" "cross_region" {}
provider "alicloud" {
 region = "cn-hangzhou"
}
provider "alicloudalt" {
 region = "cn-beijing"
}
resource "alicloud_vpc" "local" {
 vpc_name = "%[1]s-local"
 cidr_block = "172.16.0.0/16"
}
resource "alicloud_vpc" "remote" {
 provider = alicloudalt
 vpc_name = "%[1]s-remote"
 cidr_block = "172.17.0.0/16"
}
resource "alicloud_express_connect_router_interface" "remote" {
 provider = alicloudalt
 router_id = alicloud_vpc.remote.router_id
 router_type = "VRouter"
 role = "AcceptingSide"
 spec = "Negative"
 fast_link_mode = false
 payment_type = "PayAsYouGo"
 opposite_router_type = "VRouter"
 opposite_router_id = alicloud_vpc.local.router_id
 opposite_interface_owner_id = data.alicloud_account.cross_region.id
 opposite_region_id = "cn-hangzhou"
 router_interface_name = "%[1]s-remote"
}
resource "alicloud_router_interface_connection" "local" {
 interface_id = alicloud_express_connect_router_interface.default.id
 opposite_interface_id = alicloud_express_connect_router_interface.remote.id
 depends_on = [alicloud_router_interface_connection.remote]
}
resource "alicloud_router_interface_connection" "remote" {
 provider = alicloudalt
 interface_id = alicloud_express_connect_router_interface.remote.id
 opposite_interface_id = alicloud_express_connect_router_interface.default.id
}
`, name)
}

func TestUnitExpressConnectRouterInterfaceDeleteWait(t *testing.T) {
	for _, tc := range []struct {
		name        string
		status      string
		queryError  string
		deleteError string
		timeout     time.Duration
		wantError   string
		wantReads   int
	}{
		{name: "wait_for_absence", status: "Inactive", timeout: 20 * time.Second, wantReads: 2},
		{name: "idle_is_not_deleted", status: "Idle", timeout: 6 * time.Second, wantError: "timeout", wantReads: 1},
		{name: "query_error", queryError: "Forbidden", timeout: 10 * time.Second, wantError: "Forbidden", wantReads: 1},
		{name: "query_not_found", queryError: "InvalidInstanceId.NotFound", timeout: 10 * time.Second, wantReads: 1},
		{name: "delete_not_found", deleteError: "InvalidInstanceId.NotFound", timeout: 10 * time.Second, wantReads: 0},
		{name: "delete_error", deleteError: "Forbidden", timeout: 10 * time.Second, wantError: "Forbidden", wantReads: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reads := 0
			client := routerInterfaceDeleteTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Form.Get("Action") {
				case "DeleteRouterInterface":
					if r.Form.Get("RouterInterfaceId") != "ri-test" {
						t.Error("wrong router interface")
					}
					if tc.deleteError != "" {
						routerInterfaceDeleteTestError(t, w, tc.deleteError)
						return
					}
					routerInterfaceDeleteTestResponse(t, w, map[string]interface{}{})
				case "DescribeRouterInterfaceAttribute":
					reads++
					if r.Form.Get("InstanceId") != "ri-test" {
						t.Error("wrong waiter instance")
					}
					if tc.queryError != "" {
						routerInterfaceDeleteTestError(t, w, tc.queryError)
						return
					}
					if reads == 1 || tc.name == "idle_is_not_deleted" {
						routerInterfaceDeleteTestResponse(t, w, map[string]interface{}{"RouterInterfaceId": "ri-test", "Status": tc.status, "ChargeType": "AfterPay"})
						return
					}
					routerInterfaceDeleteTestError(t, w, "InvalidInstanceId.NotFound")
				default:
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
				}
			})
			res := resourceAliCloudExpressConnectRouterInterface()
			res.Timeouts = &schema.ResourceTimeout{Delete: schema.DefaultTimeout(tc.timeout)}
			d := res.Data(nil)
			d.SetId("ri-test")
			if err := d.Set("fast_link_mode", false); err != nil {
				t.Fatal(err)
			}
			if err := d.Set("payment_type", "PayAsYouGo"); err != nil {
				t.Fatal(err)
			}
			err := res.Delete(d, client)
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Errorf("wanted %s, got %v", tc.wantError, err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if tc.name == "idle_is_not_deleted" {
				if reads < 1 || reads > 3 {
					t.Errorf("expected bounded polling of persistent Idle, got %d reads", reads)
				}
			} else if reads != tc.wantReads {
				t.Errorf("expected %d reads before Delete returns, got %d", tc.wantReads, reads)
			}
		})
	}
}
