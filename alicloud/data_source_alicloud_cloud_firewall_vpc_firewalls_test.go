package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudFirewallVpcFirewallDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallVpcFirewallSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_firewall.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_firewall.default.id}_fake"]`,
		}),
	}
	StatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallSourceConfig(rand, map[string]string{
			"vpc_firewall_id": `"${alicloud_cloud_firewall_vpc_firewall.default.id}"`,
			"status":          `"opened"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallSourceConfig(rand, map[string]string{
			"vpc_firewall_id": `"${alicloud_cloud_firewall_vpc_firewall.default.id}"`,
			"status":          `"closed"`,
		}),
	}
	VpcFirewallNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallSourceConfig(rand, map[string]string{
			"vpc_firewall_id":   `"${alicloud_cloud_firewall_vpc_firewall.default.id}"`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall.default.vpc_firewall_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallSourceConfig(rand, map[string]string{
			"vpc_firewall_id":   `"${alicloud_cloud_firewall_vpc_firewall.default.id}"`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall.default.vpc_firewall_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall.default.id}"]`,
			"status":            `"opened"`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall.default.vpc_firewall_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall.default.id}_fake"]`,
			"status":            `"closed"`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall.default.vpc_firewall_name}_fake"`,
		}),
	}

	CloudFirewallVpcFirewallCheckInfo.dataSourceTestCheck(t, rand, idsConf, StatusConf, VpcFirewallNameConf, allConf)
}

var existCloudFirewallVpcFirewallMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"firewalls.#":    "1",
		"firewalls.0.id": CHECKSET,
	}
}

var fakeCloudFirewallVpcFirewallMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"firewalls.#": "0",
	}
}

var CloudFirewallVpcFirewallCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_firewall_vpc_firewalls.default",
	existMapFunc: existCloudFirewallVpcFirewallMapFunc,
	fakeMapFunc:  fakeCloudFirewallVpcFirewallMapFunc,
}

func testAccCheckAlicloudCloudFirewallVpcFirewallSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-%d"
}
data "alicloud_account" "current" {
}
data "alicloud_vpcs" "vpcs_ds" {
  name_regex = "^cfw-1-default-NODELETING"
}
data "alicloud_route_tables" "local_vpc" {
  vpc_id = "${data.alicloud_vpcs.vpcs_ds.vpcs.0.id}"
}
data "alicloud_vpcs" "vpcs_ds_peer" {
  name_regex = "^cfw-2-default-NODELETING"
}
data "alicloud_route_tables" "local_peer" {
  vpc_id = "${data.alicloud_vpcs.vpcs_ds_peer.vpcs.0.id}"
}
data "alicloud_vpc_peer_connections" "cfw_vpc_peer" {
  name_regex = "^cfw-default-NODELETING"
}

resource "alicloud_cloud_firewall_vpc_firewall" "default" {
  vpc_firewall_name = "${var.name}"
  member_uid        = "${data.alicloud_account.current.id}"
  peer_vpc {
    vpc_id    = "${data.alicloud_vpcs.vpcs_ds_peer.vpcs.0.id}"
    region_no = "%s"
	peer_vpc_cidr_table_list {
      peer_route_table_id = "${data.alicloud_route_tables.local_peer.tables.0.id}"
      peer_route_entry_list {
        peer_destination_cidr     = "${data.alicloud_vpcs.vpcs_ds.vpcs.0.cidr_block}"
        peer_next_hop_instance_id = "${data.alicloud_vpc_peer_connections.cfw_vpc_peer.connections.0.id}"
      }
    }
  }
  local_vpc {
    vpc_id    = "${data.alicloud_vpcs.vpcs_ds.vpcs.0.id}"
    region_no = "%s"
    local_vpc_cidr_table_list {
      local_route_table_id = "${data.alicloud_route_tables.local_vpc.tables.0.id}"
      local_route_entry_list {
        local_next_hop_instance_id = "${data.alicloud_vpc_peer_connections.cfw_vpc_peer.connections.0.id}"
        local_destination_cidr     = "${data.alicloud_vpcs.vpcs_ds_peer.vpcs.0.cidr_block}"
      }
    }
  }
  status = "open"
}

data "alicloud_cloud_firewall_vpc_firewalls" "default" {
%s
}

`, rand, os.Getenv("ALICLOUD_REGION"), os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, "\n   "))
	return config
}
