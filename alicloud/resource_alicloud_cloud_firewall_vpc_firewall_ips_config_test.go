package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall VpcFirewallIpsConfig. >>> Resource test cases, automatically generated.
// Case 入侵防护设置 7786
func TestAccAliCloudCloudFirewallVpcFirewallIpsConfig_basic7786(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_ips_config.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcFirewallIpsConfigMap7786)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcFirewallIpsConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcFirewallIpsConfigBasicDependence7786)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_all_patch": "0",
					"basic_rules":      "0",
					"run_mode":         "0",
					"vpc_firewall_id":  "${alicloud_cloud_firewall_vpc_firewall.test.id}",
					"rule_class":       "0",
					"lang":             "cn-shenzhen",
					"member_uid":       "${data.alicloud_account.test.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_all_patch": "0",
						"basic_rules":      "0",
						"run_mode":         "0",
						"vpc_firewall_id":  CHECKSET,
						"rule_class":       CHECKSET,
						"lang":             "cn-shenzhen",
						"member_uid":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_all_patch": "1",
					"basic_rules":      "1",
					"run_mode":         "1",
					"rule_class":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_all_patch": "1",
						"basic_rules":      "1",
						"run_mode":         "1",
						"rule_class":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "member_uid"},
			},
		},
	})
}

var AlicloudCloudFirewallVpcFirewallIpsConfigMap7786 = map[string]string{}

func AlicloudCloudFirewallVpcFirewallIpsConfigBasicDependence7786(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "test" {
}

data "alicloud_regions" "test" {
  current = true
}
data "alicloud_zones" "test" {
  available_instance_type     = "ecs.sn1ne.large"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "peer" {
  vpc_name   = "${var.name}-peer"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "local" {
  vpc_name   = "${var.name}-local"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "peer01" {
  vpc_id       = alicloud_vpc.peer.id
  cidr_block   = "192.168.10.0/24"
  zone_id      = data.alicloud_zones.test.zones.0.id
  vswitch_name = "${var.name}-peer-01"
}
resource "alicloud_vswitch" "peer02" {
  vpc_id       = alicloud_vpc.peer.id
  cidr_block   = "192.168.20.0/24"
  zone_id      = data.alicloud_zones.test.zones.1.id
  vswitch_name = "${var.name}-peer-02"
}
resource "alicloud_vswitch" "local01" {
  vpc_id       = alicloud_vpc.local.id
  cidr_block   = "172.16.10.0/24"
  zone_id      = data.alicloud_zones.test.zones.0.id
  vswitch_name = "${var.name}-local-01"
}
resource "alicloud_vswitch" "local02" {
  vpc_id       = alicloud_vpc.local.id
  cidr_block   = "172.16.20.0/24"
  zone_id      = data.alicloud_zones.test.zones.1.id
  vswitch_name = "${var.name}-local-02"
}
resource "alicloud_vpc_peer_connection" "test" {
  peer_connection_name = var.name
  vpc_id               = alicloud_vpc.local.id
  accepting_ali_uid    = data.alicloud_account.test.id
  accepting_region_id  = data.alicloud_regions.test.ids.0
  accepting_vpc_id     = alicloud_vpc.peer.id
  description          = "terraform-example"
  force_delete         = true
}

resource "alicloud_vpc_peer_connection_accepter" "test" {
  instance_id  = alicloud_vpc_peer_connection.test.id
  force_delete = true
}

resource "alicloud_route_entry" "local" {
  route_table_id        = alicloud_vpc.local.route_table_id
  destination_cidrblock = "1.2.3.4/32"
  nexthop_type          = "VpcPeer"
  nexthop_id            = alicloud_vpc_peer_connection.test.id
}

resource "alicloud_route_entry" "peer" {
  route_table_id        = alicloud_vpc.peer.route_table_id
  destination_cidrblock = "4.3.2.1/32"
  nexthop_type          = "VpcPeer"
  nexthop_id            = alicloud_vpc_peer_connection.test.id
}

resource "alicloud_cloud_firewall_vpc_firewall" "test" {
  vpc_firewall_name = var.name
  member_uid        = data.alicloud_account.test.id
  peer_vpc {
    vpc_id    = alicloud_vpc.peer.id
    region_no = data.alicloud_regions.test.ids.0
    peer_vpc_cidr_table_list {
      peer_route_table_id = alicloud_vpc.peer.route_table_id
      peer_route_entry_list {
        peer_destination_cidr     = alicloud_route_entry.peer.destination_cidrblock
        peer_next_hop_instance_id = alicloud_vpc_peer_connection.test.id
      }
    }
  }
  local_vpc {
    vpc_id    = alicloud_vpc.local.id
    region_no = data.alicloud_regions.test.ids.0
    local_vpc_cidr_table_list {
      local_route_table_id = alicloud_vpc.local.route_table_id
      local_route_entry_list {
        local_next_hop_instance_id = alicloud_vpc_peer_connection.test.id
        local_destination_cidr     = alicloud_route_entry.local.destination_cidrblock
      }
    }
  }
  status = "open"
}`, name)
}
