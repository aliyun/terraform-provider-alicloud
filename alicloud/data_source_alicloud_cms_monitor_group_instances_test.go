package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsMonitorGroupInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMonitorGroupInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_monitor_group_instances.default.id}"]`,
		}),
	}
	var existAlicloudCmsMonitorGroupInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"instances.#":                           "1",
			"instances.0.instances.0.category":      "slb",
			"instances.0.instances.0.instance_id":   CHECKSET,
			"instances.0.instances.0.instance_name": fmt.Sprintf("tf-testAccMonitorGroupInstances-%d", rand),
			"instances.0.instances.0.region_id":     defaultRegionToTest,
		}
	}
	var fakeAlicloudCmsMonitorGroupInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCmsMonitorGroupInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_monitor_group_instances.default",
		existMapFunc: existAlicloudCmsMonitorGroupInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsMonitorGroupInstancesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.FnfSupportRegions)
	}
	alicloudCmsMonitorGroupInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}
func testAccCheckAlicloudCmsMonitorGroupInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccMonitorGroupInstances-%d"
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
  load_balancer_name = var.name
  load_balancer_spec = "slb.s2.small"
  vswitch_id = data.alicloud_vswitches.default.ids.0
}
resource "alicloud_cms_monitor_group" "default" {
monitor_group_name = var.name
}
resource "alicloud_cms_monitor_group_instances" "default" {
  group_id = alicloud_cms_monitor_group.default.id
  instances {
    instance_id = alicloud_slb_load_balancer.default.id
    instance_name = alicloud_slb_load_balancer.default.name
    region_id = "%s"
    category = "slb"
  }
}

data "alicloud_cms_monitor_group_instances" "default" {	
	%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
