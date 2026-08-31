package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDdosbgpIpsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.DdosBgpSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ddosbgp_ip.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ddosbgp_ip.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ddosbgp_ip.default.id}"]`,
			"status": `"normal"`,
		}),
		fakeConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ddosbgp_ip.default.id}"]`,
			"status": `"hole_begin"`,
		}),
	}
	productNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ddosbgp_ip.default.id}"]`,
			"product_name": `"EIP"`,
		}),
		fakeConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ddosbgp_ip.default.id}"]`,
			"product_name": `"WAF"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ddosbgp_ip.default.id}"]`,
			"status":       `"normal"`,
			"product_name": `"EIP"`,
		}),
		fakeConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ddosbgp_ip.default.id}_fake"]`,
			"status":       `"hole_begin"`,
			"product_name": `"WAF"`,
		}),
	}
	var existAlicloudDdosbgpIpsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "1",
			"ips.#":             "1",
			"ips.0.instance_id": CHECKSET,
			"ips.0.ip":          CHECKSET,
			"ips.0.id":          CHECKSET,
			"ips.0.product":     "EIP",
			"ips.0.status":      CHECKSET,
		}
	}
	var fakeAlicloudDdosbgpIpsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudDdosbgpIpsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ddosbgp_ips.default",
		existMapFunc: existAlicloudDdosbgpIpsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDdosbgpIpsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudDdosbgpIpsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, productNameConf, allConf)
}
func testAccCheckAlicloudDdosbgpIpsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccIp-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eip_address" "default" {
	address_name = "${var.name}"
}

data "alicloud_ddosbgp_instances" default {}

resource "alicloud_ddosbgp_ip" "default" {
	instance_id = data.alicloud_ddosbgp_instances.default.ids.0
	ip = alicloud_eip_address.default.ip_address
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
data "alicloud_ddosbgp_ips" "default" {	
	instance_id = data.alicloud_ddosbgp_instances.default.ids.0
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
