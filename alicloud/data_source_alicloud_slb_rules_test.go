package alicloud

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlicloudSLBRulesDataSource_basic(t *testing.T) {
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbRulesDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_rule.default.load_balancer_id}"`,
			"frontend_port":    `"${alicloud_slb_rule.default.frontend_port}"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbRulesDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_rule.default.load_balancer_id}"`,
			"frontend_port":    `"${alicloud_slb_rule.default.frontend_port}"`,
			"name_regex":       `"${alicloud_slb_rule.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbRulesDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_rule.default.load_balancer_id}"`,
			"frontend_port":    `"${alicloud_slb_rule.default.frontend_port}"`,
			"name_regex":       `"${alicloud_slb_rule.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbRulesDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_rule.default.load_balancer_id}"`,
			"frontend_port":    `"${alicloud_slb_rule.default.frontend_port}"`,
			"ids":              `["${alicloud_slb_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSlbRulesDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_rule.default.load_balancer_id}"`,
			"frontend_port":    `"${alicloud_slb_rule.default.frontend_port}"`,
			"ids":              `["${alicloud_slb_rule.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbRulesDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_rule.default.load_balancer_id}"`,
			"frontend_port":    `"${alicloud_slb_rule.default.frontend_port}"`,
			"ids":              `["${alicloud_slb_rule.default.id}"]`,
			"name_regex":       `"${alicloud_slb_rule.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbRulesDataSourceConfig(map[string]string{
			"load_balancer_id": `"${alicloud_slb_rule.default.load_balancer_id}"`,
			"frontend_port":    `"${alicloud_slb_rule.default.frontend_port}"`,
			"ids":              `["${alicloud_slb_rule.default.id}_fake"]`,
			"name_regex":       `"${alicloud_slb_rule.default.name}"`,
		}),
	}

	var existSLBRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_rules.#":                 "1",
			"ids.#":                       "1",
			"names.#":                     "1",
			"slb_rules.0.id":              CHECKSET,
			"slb_rules.0.name":            "tf-testaccslbrulesdatasourcebasic",
			"slb_rules.0.domain":          "*.aliyun.com",
			"slb_rules.0.url":             "/image",
			"slb_rules.0.server_group_id": CHECKSET,
		}
	}

	var fakeSLBRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"slb_rules.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var slbRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_rules.default",
		existMapFunc: existSLBRulesMapFunc,
		fakeMapFunc:  fakeSLBRulesMapFunc,
	}

	slbRulesCheckInfo.dataSourceTestCheck(t, -1, basicConf, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudSlbRulesDataSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testaccslbrulesdatasourcebasic"
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners = "system"
}
data "alicloud_instance_types" "default" {
 	cpu_core_count = 1
	memory_size = 2
	availability_zone = data.alicloud_slb_zones.default.zones.0.id
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

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "${var.name}"
  cookie_timeout = 86400
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
}

resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  servers {
      server_ids = ["${alicloud_instance.default.id}"]
      port = 80
      weight = 100
    }
}

resource "alicloud_slb_rule" "default" {
  load_balancer_id = "${alicloud_slb_load_balancer.default.id}"
  frontend_port = "${alicloud_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alicloud_slb_server_group.default.id}"
}

data "alicloud_slb_rules" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}
