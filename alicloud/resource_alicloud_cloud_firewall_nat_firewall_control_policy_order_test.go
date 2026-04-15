package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall NatFirewallControlPolicyOrder. >>> Resource test cases, automatically generated.
// Case test 12703
func TestAccAliCloudCloudFirewallNatFirewallControlPolicyOrder_basic12703(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_nat_firewall_control_policy_order.test.0"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallNatFirewallControlPolicyOrderMap12703)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallNatFirewallControlPolicyOrder")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	resourceTestAccConfigFunc(resourceId, name, alicloudCloudFirewallNatFirewallControlPolicyOrderBasicDependence12703)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: alicloudCloudFirewallNatFirewallControlPolicyOrderBasicDependence12703(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_uuid":       CHECKSET,
						"nat_gateway_id": CHECKSET,
						"direction":      "out",
						"order":          CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudCloudFirewallNatFirewallControlPolicyOrderMap12703 = map[string]string{}

func alicloudCloudFirewallNatFirewallControlPolicyOrderBasicDependence12703(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "test" {}

resource "alicloud_vpc" "test" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "test" {
  vswitch_name = var.name
  zone_id      = data.alicloud_enhanced_nat_available_zones.test.zones.0.zone_id
  cidr_block   = "10.10.0.0/20"
  vpc_id       = alicloud_vpc.test.id
}

resource "alicloud_nat_gateway" "test" {
  vpc_id           = alicloud_vpc.test.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.test.id
  nat_type         = "Enhanced"
}

resource "alicloud_cloud_firewall_address_book" "test" {
  description  = "${var.name}test"
  group_name   = "${var.name}test"
  group_type   = "port"
  address_list = ["22/22", "23/23", "24/24"]
}

resource "alicloud_cloud_firewall_nat_firewall_control_policy" "test" {
  count                 = 3
  nat_gateway_id        = alicloud_nat_gateway.test.id
  application_name_list = ["ANY"]
  release               = "false"
  ip_version            = "4"
  repeat_days           = ["1"]
  repeat_start_time     = "21:00"
  acl_action            = "log"
  dest_port_group       = alicloud_cloud_firewall_address_book.test.group_name
  repeat_type           = "Weekly"
  source                = "1.1.1.1/32"
  direction             = "out"
  repeat_end_time       = "21:30"
  start_time            = "1699156800"
  destination           = "1.1.1.1/32"
  end_time              = "1888545600"
  source_type           = "net"
  proto                 = "TCP"
  new_order             = "-1"
  destination_type      = "net"
  dest_port_type        = "group"
  domain_resolve_type   = "0"
  description           = var.name
  depends_on            = [alicloud_nat_gateway.test]
  lifecycle {
    ignore_changes = [new_order]
  }
}

resource "alicloud_cloud_firewall_nat_firewall_control_policy_order" "test" {
  count          = 3
  order          = "${count.index + 1}"
  acl_uuid       = "${alicloud_cloud_firewall_nat_firewall_control_policy.test[count.index].acl_uuid}"
  nat_gateway_id = "${alicloud_nat_gateway.test.id}"
  direction      = "out"
}`, name)
}

// Test CloudFirewall NatFirewallControlPolicyOrder. <<< Resource test cases, automatically generated.
