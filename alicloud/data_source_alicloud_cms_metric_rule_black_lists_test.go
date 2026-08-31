package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsMetricRuleBlackListsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCmsMetricRuleBlackList%d", defaultRegionToTest, rand)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMetricRuleBlackListsDataSourceName(name, map[string]string{
			"category": `"ecs"`,
		}),
		fakeConfig: "",
	}
	var existAlicloudCmsMetricRuleBlackListsDataSourceNameMapFunc = func(name int) map[string]string {
		return map[string]string{
			"lists.#": CHECKSET,
		}
	}
	var fakeAlicloudCmsMetricRuleBlackListsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"lists.#": "0",
		}
	}
	var alicloudCmsMetricRuleBlackListsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_metric_rule_black_lists.default",
		existMapFunc: existAlicloudCmsMetricRuleBlackListsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsMetricRuleBlackListsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCmsMetricRuleBlackListsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckAlicloudCmsMetricRuleBlackListsDataSourceName(rand string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
  availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_instance_types" "new" {
	eni_amount = 2
    availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners = "system"
}

data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.1
	zone_id = data.alicloud_slb_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  	name = "${var.name}"
	vpc_id = data.alicloud_vpcs.default.ids.1
}
resource "alicloud_instance" "instance" {
  	image_id = "${data.alicloud_images.default.images.0.id}"
  	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
  	security_groups = "${alicloud_security_group.default.*.id}"
  	internet_charge_type = "PayByTraffic"
  	internet_max_bandwidth_out = "10"
  	availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  	instance_charge_type = "PostPaid"
  	system_disk_category = "cloud_efficiency"
  	vswitch_id = data.alicloud_vswitches.default.ids[0]
}
resource "alicloud_cms_metric_rule_black_list" "default" {
  instances = [
  "{\"instancceId\":\"${alicloud_instance.instance.id}\"}"]
  metrics {
    metric_name = "disk_utilization"
  }
  category                    = "ecs"
  enable_end_time             = 1640608200000
  namespace                   = "acs_ecs_dashboard"
  enable_start_time           = 1640237400000
  metric_rule_black_list_name = "${var.name}"
}
data "alicloud_cms_metric_rule_black_lists" "default" {
	ids        = ["${alicloud_cms_metric_rule_black_list.default.id}"]
  	namespace  = "acs_ecs_dashboard"
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
