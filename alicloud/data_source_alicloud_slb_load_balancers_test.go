package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSlbLoadBalancersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_slb_load_balancer.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_slb_load_balancer.default.id}_fake"]`,
		}),
	}
	vpcIDConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_slb_load_balancer.default.id}"]`,
			"vpc_id": `"${alicloud_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_slb_load_balancer.default.id}"]`,
			"vpc_id": `"${alicloud_vpc.default.id}_fake"`,
		}),
	}
	vswitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_slb_load_balancer.default.id}"]`,
			"vswitch_id": `"${alicloud_slb_load_balancer.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_slb_load_balancer.default.id}"]`,
			"vswitch_id": `"${alicloud_slb_load_balancer.default.vswitch_id}_fake"`,
		}),
	}
	netWorkTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_slb_load_balancer.default.id}"]`,
			"network_type": `"vpc"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_slb_load_balancer.default.id}"]`,
			"network_type": `"classic"`,
		}),
	}
	masterZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_slb_load_balancer.default.id}"]`,
			"master_zone_id": `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_slb_load_balancer.default.id}"]`,
			"master_zone_id": `"${data.alicloud_zones.default.zones.0.id}_fake"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_slb_load_balancer.default.id}"]`,
			"resource_group_id": `"${alicloud_slb_load_balancer.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_slb_load_balancer.default.id}"]`,
			"resource_group_id": `"${alicloud_slb_load_balancer.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"tags": fmt.Sprintf(`{
				Created = "TF-%d"
		}`, rand),
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"tags": fmt.Sprintf(`{
				Created = "fake-%d"
		}`, rand),
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_slb_load_balancer.default.load_balancer_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_slb_load_balancer.default.load_balancer_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_slb_load_balancer.default.id}"]`,
			"status": `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_slb_load_balancer.default.id}"]`,
			"status": `"inactive"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_slb_load_balancer.default.id}"]`,
			"name_regex":        `"${alicloud_slb_load_balancer.default.load_balancer_name}"`,
			"vpc_id":            `"${alicloud_vpc.default.id}"`,
			"vswitch_id":        `"${alicloud_slb_load_balancer.default.vswitch_id}"`,
			"network_type":      `"vpc"`,
			"master_zone_id":    `"${data.alicloud_zones.default.zones.0.id}"`,
			"resource_group_id": `"${alicloud_slb_load_balancer.default.resource_group_id}"`,
			"status":            `"active"`,
			"tags": fmt.Sprintf(`{
				Created = "TF-%d"
		}`, rand),
		}),
		fakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_slb_load_balancer.default.id}_fake"]`,
			"name_regex":        `"${alicloud_slb_load_balancer.default.load_balancer_name}_fake"`,
			"vpc_id":            `"${alicloud_vpc.default.id}_fake"`,
			"vswitch_id":        `"${alicloud_slb_load_balancer.default.vswitch_id}_fake"`,
			"network_type":      `"classic"`,
			"master_zone_id":    `"${data.alicloud_zones.default.zones.0.id}_fake"`,
			"resource_group_id": `"${alicloud_slb_load_balancer.default.resource_group_id}_fake"`,
			"status":            `"inactive"`,
			"tags": fmt.Sprintf(`{
				Created = "fake-%d"
		}`, rand),
		}),
	}
	var existAlicloudSlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                      "1",
			"names.#":                                    "1",
			"balancers.#":                                "1",
			"balancers.0.address":                        CHECKSET,
			"balancers.0.address_ip_version":             `ipv4`,
			"balancers.0.address_type":                   `intranet`,
			"balancers.0.bandwidth":                      CHECKSET,
			"balancers.0.internet_charge_type":           `PayByTraffic`,
			"balancers.0.delete_protection":              `off`,
			"balancers.0.load_balancer_name":             CHECKSET,
			"balancers.0.master_zone_id":                 CHECKSET,
			"balancers.0.modification_protection_reason": "",
			"balancers.0.modification_protection_status": CHECKSET,
			"balancers.0.payment_type":                   `PayAsYouGo`,
			"balancers.0.resource_group_id":              CHECKSET,
			"balancers.0.slave_zone_id":                  CHECKSET,
			"balancers.0.load_balancer_spec":             `slb.s1.small`,
			"balancers.0.status":                         `active`,
			"balancers.0.tags.%":                         `1`,
			"balancers.0.vswitch_id":                     CHECKSET,
		}
	}
	var fakeAlicloudSlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSlbLoadBalancersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_load_balancers.default",
		existMapFunc: existAlicloudSlbLoadBalancersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSlbLoadBalancersDataSourceNameMapFunc,
	}
	alicloudSlbLoadBalancersCheckInfo.dataSourceTestCheck(t, rand, idsConf, vpcIDConf, vswitchConf, netWorkTypeConf, masterZoneConf, resourceGroupIdConf, tagsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccLoadBalancer-%[1]d"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {
	name_regex = "^default$"
}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_slb_load_balancer" "default" {
	address_type = "intranet"
	vswitch_id = alicloud_vswitch.default.id
	load_balancer_name = var.name
	load_balancer_spec = "slb.s1.small"
    master_zone_id = "${data.alicloud_zones.default.zones.0.id}"
  	resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
	tags =  {
		Created = "TF-%[1]d"
	}
}

data "alicloud_slb_load_balancers" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
