package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall NatFirewall. >>> Resource test cases, automatically generated.
// Case 云防火墙资源测试-v1 6822
func TestAccAliCloudCloudFirewallNatFirewall_basic6822(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
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
					"status":         "normal",
					"region_no":      "${data.alicloud_regions.current.ids.0}",
					"vswitch_auto":   "true",
					"strict_mode":    "0",
					"vpc_id":         "${alicloud_vpc.defaultikZ0gD.id}",
					"proxy_name":     "YQC-防火墙测试",
					"lang":           "zh",
					"nat_gateway_id": "${alicloud_nat_gateway.default2iRZpC.id}",
					"nat_route_entry_list": []map[string]interface{}{
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "0.0.0.0/0",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
					"firewall_switch": "close",
					"vswitch_cidr":    "172.16.5.0/24",
					"vswitch_id":      "${alicloud_snat_entry.defaultAKE43g.source_vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_no":              defaultRegionToTest,
						"vpc_id":                 CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"nat_gateway_id":         CHECKSET,
						"nat_route_entry_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "closed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "closed",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"region_no":      "${data.alicloud_regions.current.ids.0}",
					"vswitch_auto":   "true",
					"strict_mode":    "0",
					"vpc_id":         "${alicloud_vpc.defaultikZ0gD.id}",
					"proxy_name":     "YQC-防火墙测试",
					"lang":           "zh",
					"nat_gateway_id": "${alicloud_nat_gateway.default2iRZpC.id}",
					"nat_route_entry_list": []map[string]interface{}{
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "0.0.0.0/0",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
					"firewall_switch": "close",
					"vswitch_cidr":    "172.16.5.0/24",
					"status":          "closed",
					"vswitch_id":      "${alicloud_vswitch.defaultp4O7qi.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_no":              defaultRegionToTest,
						"vswitch_auto":           "true",
						"strict_mode":            "0",
						"vpc_id":                 CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"lang":                   "zh",
						"nat_gateway_id":         CHECKSET,
						"nat_route_entry_list.#": "1",
						"firewall_switch":        "close",
						"vswitch_cidr":           "172.16.5.0/24",
						"status":                 "closed",
						"vswitch_id":             CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"firewall_switch", "lang", "vswitch_id", "vswitch_auto", "vswitch_cidr"},
			},
		},
	})
}

var AlicloudCloudFirewallNatFirewallMap6822 = map[string]string{
	"status": CHECKSET,
}

func AlicloudCloudFirewallNatFirewallBasicDependence6822(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_regions" "current" {
  current = true
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultikZ0gD" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultp4O7qi" {
  vpc_id       = alicloud_vpc.defaultikZ0gD.id
  cidr_block   = "172.16.6.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_nat_gateway" "default2iRZpC" {
  eip_bind_mode    = "MULTI_BINDED"
  vpc_id           = alicloud_vpc.defaultikZ0gD.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.defaultp4O7qi.id
  nat_type         = "Enhanced"
  network_type = "internet"
}

resource "alicloud_eip_address" "defaultyiRwgs" {
  address_name = var.name
}

resource "alicloud_eip_association" "defaults2MTuO" {
  instance_id   = alicloud_nat_gateway.default2iRZpC.id
  allocation_id = alicloud_eip_address.defaultyiRwgs.id
  mode          = "NAT"
  instance_type = "Nat"
}

resource "alicloud_snat_entry" "defaultAKE43g" {
  snat_ip           = alicloud_eip_address.defaultyiRwgs.ip_address
  snat_table_id     = alicloud_nat_gateway.default2iRZpC.snat_table_ids
  source_vswitch_id = alicloud_vswitch.defaultp4O7qi.id
}


`, name)
}

// Case 云防火墙资源测试-v1 6822  twin
func TestAccAliCloudCloudFirewallNatFirewall_basic6822_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"region_no":      "${data.alicloud_regions.current.ids.0}",
					"vswitch_auto":   "true",
					"strict_mode":    "0",
					"vpc_id":         "${alicloud_vpc.defaultikZ0gD.id}",
					"proxy_name":     "YQC-防火墙测试",
					"lang":           "zh",
					"nat_gateway_id": "${alicloud_nat_gateway.default2iRZpC.id}",
					"nat_route_entry_list": []map[string]interface{}{
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "0.0.0.0/0",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
					"firewall_switch": "close",
					"vswitch_cidr":    "172.16.5.0/24",
					"status":          "closed",
					"vswitch_id":      "${alicloud_snat_entry.defaultAKE43g.source_vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_no":              defaultRegionToTest,
						"vswitch_auto":           "true",
						"strict_mode":            "0",
						"vpc_id":                 CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"lang":                   "zh",
						"nat_gateway_id":         CHECKSET,
						"nat_route_entry_list.#": "1",
						"firewall_switch":        "close",
						"vswitch_cidr":           "172.16.5.0/24",
						"status":                 "closed",
						"vswitch_id":             CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"firewall_switch", "lang", "vswitch_id", "vswitch_auto", "vswitch_cidr"},
			},
		},
	})
}

// Case 云防火墙资源测试-v1 6822  raw
func TestAccAliCloudCloudFirewallNatFirewall_basic6822_raw(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"region_no":      "${data.alicloud_regions.current.ids.0}",
					"vswitch_auto":   "true",
					"strict_mode":    "0",
					"vpc_id":         "${alicloud_vpc.defaultikZ0gD.id}",
					"proxy_name":     "YQC-防火墙测试",
					"lang":           "zh",
					"nat_gateway_id": "${alicloud_nat_gateway.default2iRZpC.id}",
					"nat_route_entry_list": []map[string]interface{}{
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "0.0.0.0/0",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
					"firewall_switch": "close",
					"vswitch_cidr":    "172.16.5.0/24",
					"status":          "closed",
					"vswitch_id":      "${alicloud_snat_entry.defaultAKE43g.source_vswitch_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_no":              CHECKSET,
						"vswitch_auto":           "true",
						"strict_mode":            "0",
						"vpc_id":                 CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"lang":                   "zh",
						"nat_gateway_id":         CHECKSET,
						"nat_route_entry_list.#": "1",
						"firewall_switch":        "close",
						"vswitch_cidr":           "172.16.5.0/24",
						"status":                 "closed",
						"vswitch_id":             CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "normal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"firewall_switch", "lang", "vswitch_id", "vswitch_auto", "vswitch_cidr"},
			},
		},
	})
}

// Test CloudFirewall NatFirewall. <<< Resource test cases, automatically generated.
