package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCddcDedicatedHostGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 200)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_group.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_cddc_dedicated_host_group.default.id}"]`,
			"engine": `"MySQL"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_cddc_dedicated_host_group.default.id}_fake"]`,
			"engine": `"MySQL_fake"`,
		}),
	}
	var existAlicloudCddcDedicatedHostGroupsNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"groups.#":                           "1",
			"groups.0.engine":                    "mysql",
			"groups.0.dedicated_host_group_desc": fmt.Sprintf("tf-testAccName-%d",rand),
			"groups.0.allocation_policy":         "Evenly",
			"groups.0.cpu_allocation_ratio":      "101",
			"groups.0.mem_allocation_ratio":      "50",
			"groups.0.disk_allocation_ratio":     "200",
			"groups.0.host_replace_policy":       "Manual",
			"groups.0.create_time":               CHECKSET,
		}
	}
	var fakeAlicloudCddcDedicatedHostGroupsNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}
	var alicloudSaeNamespaceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cddc_dedicated_host_groups.default",
		existMapFunc: existAlicloudCddcDedicatedHostGroupsNameMapFunc,
		fakeMapFunc:  fakeAlicloudCddcDedicatedHostGroupsNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CddcSupportRegions)
	}
	alicloudSaeNamespaceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck,idsConf, allConf)
}
func testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccName-%d"
}
resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
	engine = "MySQL"
	vpc_id = alicloud_vpc.vpc.id
	cpu_allocation_ratio = 101
	mem_allocation_ratio = 50
	disk_allocation_ratio = 200
	allocation_policy = "Evenly"
	host_replace_policy = "Manual"
	dedicated_host_group_desc = var.name
}

data "alicloud_cddc_dedicated_host_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
