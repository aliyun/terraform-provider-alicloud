package cloudfirewall_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudfirewall "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	resource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudFirewall NatFirewall. >>> Resource test cases, automatically generated.
// Case 云防火墙资源测试-v1 6822
func TestAccAliCloudCloudFirewallNatFirewall_basic6822(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallNatFirewall")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"region_no":              acctest.DefaultRegionToTest,
						"vpc_id":                 acctest.CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"nat_gateway_id":         acctest.CHECKSET,
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
						"region_no":              acctest.DefaultRegionToTest,
						"vswitch_auto":           "true",
						"strict_mode":            "0",
						"vpc_id":                 acctest.CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"lang":                   "zh",
						"nat_gateway_id":         acctest.CHECKSET,
						"nat_route_entry_list.#": "1",
						"firewall_switch":        "close",
						"vswitch_cidr":           "172.16.5.0/24",
						"status":                 "closed",
						"vswitch_id":             acctest.CHECKSET,
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
	"status": acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallNatFirewall")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			//acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"region_no":              acctest.DefaultRegionToTest,
						"vswitch_auto":           "true",
						"strict_mode":            "0",
						"vpc_id":                 acctest.CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"lang":                   "zh",
						"nat_gateway_id":         acctest.CHECKSET,
						"nat_route_entry_list.#": "1",
						"firewall_switch":        "close",
						"vswitch_cidr":           "172.16.5.0/24",
						"status":                 "closed",
						"vswitch_id":             acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallNatFirewall")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			//acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"region_no":              acctest.CHECKSET,
						"vswitch_auto":           "true",
						"strict_mode":            "0",
						"vpc_id":                 acctest.CHECKSET,
						"proxy_name":             "YQC-防火墙测试",
						"lang":                   "zh",
						"nat_gateway_id":         acctest.CHECKSET,
						"nat_route_entry_list.#": "1",
						"firewall_switch":        "close",
						"vswitch_cidr":           "172.16.5.0/24",
						"status":                 "closed",
						"vswitch_id":             acctest.CHECKSET,
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

// TestAccAliCloudCloudFirewallNatFirewallNatRouteEntryListOrder verifies that
// reordering the multi-member nat_route_entry_list (a plain TypeList of nested
// resources with no DiffSuppressFunc) converges to an empty post-apply plan.
func TestAccAliCloudCloudFirewallNatFirewallNatRouteEntryListOrder(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_nat_firewall.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallMap6822)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallNatFirewall")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallnatfirewall%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudFirewallNatFirewallBasicDependence6822)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "192.168.0.0/16",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
					"firewall_switch": "close",
					"vswitch_cidr":    "172.16.5.0/24",
					"status":          "closed",
					"vswitch_id":      "${alicloud_snat_entry.defaultAKE43g.source_vswitch_id}",
				}),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nat_route_entry_list": []map[string]interface{}{
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "192.168.0.0/16",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
						{
							"nexthop_id":       "${alicloud_nat_gateway.default2iRZpC.id}",
							"destination_cidr": "0.0.0.0/0",
							"nexthop_type":     "NatGateway",
							"route_table_id":   "${alicloud_vpc.defaultikZ0gD.route_table_id}",
						},
					},
				}),
			},
		},
	})
}
