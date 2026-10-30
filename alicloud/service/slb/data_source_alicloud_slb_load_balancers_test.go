package slb_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSLBLoadBalancersDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_slb_load_balancer.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_slb_load_balancer.default.id}_fake"]`,
		}),
	}
	vpcIDConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_slb_load_balancer.default.id}"]`,
			"vpc_id": `"${data.alicloud_vpcs.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_slb_load_balancer.default.id}"]`,
			"vpc_id": `"${data.alicloud_vpcs.default.ids.0}_fake"`,
		}),
	}
	vswitchConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_slb_load_balancer.default.id}"]`,
			"vswitch_id": `"${alicloud_slb_load_balancer.default.vswitch_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_slb_load_balancer.default.id}"]`,
			"vswitch_id": `"${alicloud_slb_load_balancer.default.vswitch_id}_fake"`,
		}),
	}
	netWorkTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_slb_load_balancer.default.id}"]`,
			"network_type": `"vpc"`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_slb_load_balancer.default.id}"]`,
			"network_type": `"classic"`,
		}),
	}
	masterZoneConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_slb_load_balancer.default.id}"]`,
			"master_zone_id": `"${data.alicloud_slb_zones.default.zones.0.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_slb_load_balancer.default.id}"]`,
			"master_zone_id": `"${data.alicloud_slb_zones.default.zones.0.id}_fake"`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_slb_load_balancer.default.id}"]`,
			"resource_group_id": `"${alicloud_slb_load_balancer.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_slb_load_balancer.default.id}"]`,
			"resource_group_id": `"${alicloud_slb_load_balancer.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"tags": fmt.Sprintf(`{
				Created = "TF-%d"
		}`, rand),
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"tags": fmt.Sprintf(`{
				Created = "fake-%d"
		}`, rand),
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_slb_load_balancer.default.load_balancer_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_slb_load_balancer.default.load_balancer_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_slb_load_balancer.default.id}"]`,
			"status": `"active"`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_slb_load_balancer.default.id}"]`,
			"status": `"inactive"`,
		}),
	}
	pagingConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_name": `"${alicloud_slb_load_balancer.default.load_balancer_name}"`,
			"page_number":        `1`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"load_balancer_name": `"${alicloud_slb_load_balancer.default.load_balancer_name}"`,
			"page_number":        `2`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_slb_load_balancer.default.id}"]`,
			"name_regex":        `"${alicloud_slb_load_balancer.default.load_balancer_name}"`,
			"vpc_id":            `"${data.alicloud_vpcs.default.ids.0}"`,
			"vswitch_id":        `"${alicloud_slb_load_balancer.default.vswitch_id}"`,
			"network_type":      `"vpc"`,
			"master_zone_id":    `"${data.alicloud_slb_zones.default.zones.0.id}"`,
			"resource_group_id": `"${alicloud_slb_load_balancer.default.resource_group_id}"`,
			"status":            `"active"`,
			"tags": fmt.Sprintf(`{
				Created = "TF-%d"
		}`, rand),
			"page_number": `1`,
		}),
		FakeConfig: testAccCheckAlicloudSlbLoadBalancersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_slb_load_balancer.default.id}_fake"]`,
			"name_regex":        `"${alicloud_slb_load_balancer.default.load_balancer_name}_fake"`,
			"vpc_id":            `"${data.alicloud_vpcs.default.ids.0}_fake"`,
			"vswitch_id":        `"${alicloud_slb_load_balancer.default.vswitch_id}_fake"`,
			"network_type":      `"classic"`,
			"master_zone_id":    `"${data.alicloud_slb_zones.default.zones.0.id}_fake"`,
			"resource_group_id": `"${alicloud_slb_load_balancer.default.resource_group_id}_fake"`,
			"status":            `"inactive"`,
			"tags": fmt.Sprintf(`{
				Created = "fake-%d"
		}`, rand),
			"page_number": `2`,
		}),
	}
	var existAlicloudSlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                      "1",
			"names.#":                                    "1",
			"balancers.#":                                "1",
			"total_count":                                CHECKSET,
			"balancers.0.address":                        CHECKSET,
			"balancers.0.address_ip_version":             `ipv4`,
			"balancers.0.address_type":                   `intranet`,
			"balancers.0.bandwidth":                      CHECKSET,
			"balancers.0.create_time":                    CHECKSET,
			"balancers.0.create_time_stamp":              CHECKSET,
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
			"slbs.0.creation_time":                       CHECKSET,
		}
	}
	var fakeAlicloudSlbLoadBalancersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSlbLoadBalancersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_slb_load_balancers.default",
		ExistMapFunc: existAlicloudSlbLoadBalancersDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudSlbLoadBalancersDataSourceNameMapFunc,
	}
	alicloudSlbLoadBalancersCheckInfo.DataSourceTestCheck(t, rand, idsConf, vpcIDConf, vswitchConf, netWorkTypeConf, masterZoneConf, resourceGroupIdConf, tagsConf, nameRegexConf, statusConf, pagingConf, allConf)
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

data "alicloud_resource_manager_resource_groups" "default" {
	name_regex = "^default$"
}

data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_slb_zones.default.zones.0.id
}

resource "alicloud_slb_load_balancer" "default" {
	address_type = "intranet"
	vswitch_id = data.alicloud_vswitches.default.ids[0]
	load_balancer_name = var.name
	load_balancer_spec = "slb.s1.small"
    master_zone_id = "${data.alicloud_slb_zones.default.zones.0.id}"
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
