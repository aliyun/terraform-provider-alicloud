package ddosbgp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDdosbgpIpsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.DdosBgpSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ddosbgp_ip.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ddosbgp_ip.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ddosbgp_ip.default.id}"]`,
			"status": `"normal"`,
		}),
		FakeConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ddosbgp_ip.default.id}"]`,
			"status": `"hole_begin"`,
		}),
	}
	productNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ddosbgp_ip.default.id}"]`,
			"product_name": `"EIP"`,
		}),
		FakeConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ddosbgp_ip.default.id}"]`,
			"product_name": `"WAF"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ddosbgp_ip.default.id}"]`,
			"status":       `"normal"`,
			"product_name": `"EIP"`,
		}),
		FakeConfig: testAccCheckAlicloudDdosbgpIpsDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_ddosbgp_ip.default.id}_fake"]`,
			"status":       `"hole_begin"`,
			"product_name": `"WAF"`,
		}),
	}
	var existAlicloudDdosbgpIpsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "1",
			"ips.#":             "1",
			"ips.0.instance_id": acctest.CHECKSET,
			"ips.0.ip":          acctest.CHECKSET,
			"ips.0.id":          acctest.CHECKSET,
			"ips.0.product":     acctest.CHECKSET,
			"ips.0.status":      acctest.CHECKSET,
		}
	}
	var fakeAlicloudDdosbgpIpsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudDdosbgpIpsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ddosbgp_ips.default",
		ExistMapFunc: existAlicloudDdosbgpIpsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDdosbgpIpsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudDdosbgpIpsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, productNameConf, allConf)
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
