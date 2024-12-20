package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCloudFirewallVpcFirewallDataSource(t *testing.T) {
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
# 获取当前阿里云uid
data "alicloud_account" "current" {
}
data "alicloud_regions" "current" {
  current = true
}
data "alicloud_zones" "zone" {
  available_instance_type = "ecs.sn1ne.large"
  available_resource_creation = "VSwitch"
}
# 创建VPC 1
resource "alicloud_vpc" "peer" {
  vpc_name   = "${var.name}-peer"
  cidr_block = "192.168.0.0/16"
}
# 创建VPC 2
resource "alicloud_vpc" "local" {
  vpc_name   = "${var.name}-local"
  cidr_block = "172.16.0.0/12"
}
# 创建一个Vswitch CIDR 块为 192.168.10.0/24
resource "alicloud_vswitch" "peer01" {
  vpc_id       = alicloud_vpc.peer.id
  cidr_block   = "192.168.10.0/24"
  zone_id      = data.alicloud_zones.zone.zones.0.id
  vswitch_name = "${var.name}-peer-01"
}
# 创建另一个Vswitch CIDR 块为 192.168.20.0/24
resource "alicloud_vswitch" "peer02" {
  vpc_id       = alicloud_vpc.peer.id
  cidr_block   = "192.168.20.0/24"
  zone_id      = data.alicloud_zones.zone.zones.1.id
  vswitch_name = "${var.name}-peer-02"
}
# 创建一个Vswitch CIDR 块为 172.16.10.0/24
resource "alicloud_vswitch" "local01" {
  vpc_id       = alicloud_vpc.local.id
  cidr_block   = "172.16.10.0/24"
  zone_id      = data.alicloud_zones.zone.zones.0.id
  vswitch_name = "${var.name}-local-01"
}
# 创建另一个Vswitch CIDR 块为 172.16.20.0/24
resource "alicloud_vswitch" "local02" {
  vpc_id       = alicloud_vpc.local.id
  cidr_block   = "172.16.20.0/24"
  zone_id      = data.alicloud_zones.zone.zones.1.id
  vswitch_name = "${var.name}-local-02"
}
# 创建VPC对等连接
resource "alicloud_vpc_peer_connection" "default" {
  # 对等连接名称
  peer_connection_name = var.name
  # 发起方VPC_ID
  vpc_id = alicloud_vpc.local.id
  # 接收方 VPC 对等连接的 Alibaba Cloud 账号 ID
  accepting_ali_uid = data.alicloud_account.current.id
  # 接收方 VPC 对等连接的区域 ID。同区域创建时，输入与发起方相同的区域 ID；跨区域创建时，输入不同的区域 ID。
  accepting_region_id = data.alicloud_regions.current.ids.0
  # 接收端VPC_ID
  accepting_vpc_id = alicloud_vpc.peer.id
  # 描述
  description = "terraform-example"
  # 是否强制删除
  force_delete = true
}
# 接收端
resource "alicloud_vpc_peer_connection_accepter" "default" {
  instance_id = alicloud_vpc_peer_connection.default.id
  # 是否强制删除
  force_delete = true
}
# 配置路由条目-vpc-A
resource "alicloud_route_entry" "local" {
  # VPC-A 路由表ID
  route_table_id = alicloud_vpc.local.route_table_id
  # 目标网段，自定义
  destination_cidrblock = "1.2.3.4/32"
  # 下一跳类型
  nexthop_type = "VpcPeer"
  # 下一跳id
  nexthop_id = alicloud_vpc_peer_connection.default.id
}
# 配置路由条目2 -vpc-B
resource "alicloud_route_entry" "peer" {
  # VPC-A 路由表id
  route_table_id = alicloud_vpc.peer.route_table_id
  # 目标网段，自定义
  destination_cidrblock = "4.3.2.1/32"
  # 下一跳类型
  nexthop_type = "VpcPeer"
  # 下一跳id
  nexthop_id = alicloud_vpc_peer_connection.default.id
}

resource "alicloud_cloud_firewall_vpc_firewall" "default" {
  vpc_firewall_name = var.name
  member_uid        = data.alicloud_account.current.id
  peer_vpc {
    vpc_id    = alicloud_vpc.peer.id
    region_no = data.alicloud_regions.current.ids.0
	peer_vpc_cidr_table_list {
      peer_route_table_id = alicloud_vpc.peer.route_table_id
      peer_route_entry_list {
        peer_destination_cidr     = alicloud_route_entry.peer.destination_cidrblock
        peer_next_hop_instance_id = alicloud_vpc_peer_connection.default.id
      }
    }
  }
  local_vpc {
    vpc_id    = alicloud_vpc.local.id
    region_no = data.alicloud_regions.current.ids.0
    local_vpc_cidr_table_list {
      local_route_table_id = alicloud_vpc.local.route_table_id
      local_route_entry_list {
        local_next_hop_instance_id = alicloud_vpc_peer_connection.default.id
        local_destination_cidr     = alicloud_route_entry.local.destination_cidrblock
      }
    }
  }
  status = "open"
}

data "alicloud_cloud_firewall_vpc_firewalls" "default" {
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}
