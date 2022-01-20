package alicloud

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlicloudEssScalingrulesDataSource(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"ids": `["${alicloud_ess_scaling_rule.default.id}_fake"]`,
		}),
	}

	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"SimpleScalingRule"`,
		}),
		fakeConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"TargetTrackingScalingRule"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"name_regex": `"${alicloud_ess_scaling_rule.default.scaling_rule_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssScalingrulesDataSourceConfig(map[string]string{
			"ids":              `["${alicloud_ess_scaling_rule.default.id}"]`,
			"scaling_group_id": `"${alicloud_ess_scaling_rule.default.scaling_group_id}_fake"`,
			"type":             `"SimpleScalingRule"`,
			"name_regex":       `"${alicloud_ess_scaling_rule.default.scaling_rule_name}"`,
		}),
	}

	var existEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                          "1",
			"ids.#":                            "1",
			"names.#":                          "1",
			"rules.0.id":                       CHECKSET,
			"rules.0.scaling_group_id":         CHECKSET,
			"rules.0.name":                     CHECKSET,
			"rules.0.type":                     CHECKSET,
			"rules.0.cooldown":                 CHECKSET,
			"rules.0.adjustment_type":          "QuantityChangeInCapacity",
			"rules.0.adjustment_value":         "1",
			"rules.0.min_adjustment_magnitude": CHECKSET,
			"rules.0.scaling_rule_ari":         CHECKSET,
		}
	}

	var fakeEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var EssScalingrulesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scaling_rules.default",
		existMapFunc: existEssRecordsMapFunc,
		fakeMapFunc:  fakeEssRecordsMapFunc,
	}

	EssScalingrulesRecordsCheckInfo.dataSourceTestCheck(t, -1, idsConf, scalingGroupIdConf, typeConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudEssScalingrulesDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccDataSourceEssScalingRules"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alicloud_security_group.default.id}"
  	cidr_ip = "172.16.0.0/24"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${local.vswitch_id}"]
}

resource "alicloud_ess_scaling_rule" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	scaling_rule_name = "${var.name}"
	adjustment_type = "QuantityChangeInCapacity"
	adjustment_value = 1
}

data "alicloud_ess_scaling_rules" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}
