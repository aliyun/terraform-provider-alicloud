package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudFirewallVpcFirewallCenDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallVpcFirewallCenSupportRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
		}),
	}

	CenIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
			"cen_id": `"${data.alicloud_cen_instances.cen_instances_ds.instances.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
			"cen_id": `"${data.alicloud_cen_instances.cen_instances_ds.instances.0.id}_fake"`,
		}),
	}
	StatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
			"status": `"opened"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
			"status": `"closed"`,
		}),
	}
	VpcFirewallNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall_cen.default.vpc_firewall_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall_cen.default.vpc_firewall_name}fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
			"cen_id":            `"${data.alicloud_cen_instances.cen_instances_ds.instances.0.id}"`,
			"status":            `"opened"`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall_cen.default.vpc_firewall_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
			"cen_id":            `"${data.alicloud_cen_instances.cen_instances_ds.instances.0.id}_fake"`,
			"status":            `"closed"`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall_cen.default.vpc_firewall_name}_fake"`,
		}),
	}

	CloudFirewallVpcFirewallCenCheckInfo.dataSourceTestCheck(t, rand, idsConf, CenIdConf, StatusConf, VpcFirewallNameConf, allConf)
}

var existCloudFirewallVpcFirewallCenMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"cens.#":    "1",
		"cens.0.id": CHECKSET,
	}
}

var fakeCloudFirewallVpcFirewallCenMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"cens.#": "0",
	}
}

var CloudFirewallVpcFirewallCenCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_firewall_vpc_firewall_cens.default",
	existMapFunc: existCloudFirewallVpcFirewallCenMapFunc,
	fakeMapFunc:  fakeCloudFirewallVpcFirewallCenMapFunc,
}

func testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand int, attrMap map[string]string) string {
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

data "alicloud_cen_instances" "cen_instances_ds" {
  name_regex = "^cfw-test-no-deleting"
}

data "alicloud_vpcs" "vpcs_ds" {
  name_regex = "^cfw-test-no-delete1"
}

data "alicloud_vpcs" "vpcs_self" {
  name_regex = "^default-NODELETING"
}

resource "alicloud_cloud_firewall_vpc_firewall_cen" "default" {
  cen_id = "${data.alicloud_cen_instances.cen_instances_ds.instances.0.id}"
  local_vpc {
    network_instance_id = "${data.alicloud_vpcs.vpcs_ds.vpcs.0.id}"
  }
  status            = "open"
  member_uid        = "${data.alicloud_account.current.id}"
  vpc_region        = "%s"
  vpc_firewall_name ="${var.name}"
}

data "alicloud_cloud_firewall_vpc_firewall_cens" "default" {
%s
}
`, rand, os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, "\n   "))
	return config
}
