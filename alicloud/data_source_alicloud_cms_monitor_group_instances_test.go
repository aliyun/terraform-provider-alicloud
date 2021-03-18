package alicloud

import (
	"fmt"
	"os"
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
			"instances.0.instances.0.region_id":     os.Getenv("ALICLOUD_REGION"),
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
data "alicloud_vpcs" "default" {
  is_default = true
}
data "alicloud_vswitches" "default" {
  ids = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
}
resource "alicloud_slb" "default" {
  name = var.name
  specification = "slb.s2.small"
  vswitch_id = data.alicloud_vswitches.default.ids.0
}
resource "alicloud_cms_monitor_group" "default" {
monitor_group_name = var.name
}
resource "alicloud_cms_monitor_group_instances" "default" {
  group_id = alicloud_cms_monitor_group.default.id
  instances {
    instance_id = alicloud_slb.default.id
    instance_name = alicloud_slb.default.name
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
