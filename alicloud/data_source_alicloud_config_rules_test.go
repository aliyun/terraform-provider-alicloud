package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudConfigRulesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_config_rules.example"

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_rule.example.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_rule.example.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_config_rule.example.id}"]`,
			"status": `"ACTIVE"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_config_rule.example.id}"]`,
			"status": `"INACTIVE"`,
		}),
	}

	riskLevelConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}"`,
			"risk_level": `1`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}"`,
			"risk_level": `2`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}"`,
			"ids":        `["${alicloud_config_rule.example.id}"]`,
			"risk_level": `1`,
			"status":     `"ACTIVE"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}_fake"`,
			"ids":        `["${alicloud_config_rule.example.id}_fake"]`,
			"risk_level": `2`,
			"status":     `"INACTIVE"`,
		}),
	}

	var existConfigRulesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                                    "1",
			"names.#":                                    "1",
			"ids.#":                                      "1",
			"rules.0.account_id":                         CHECKSET,
			"rules.0.config_rule_arn":                    CHECKSET,
			"rules.0.id":                                 CHECKSET,
			"rules.0.config_rule_id":                     CHECKSET,
			"rules.0.config_rule_state":                  CHECKSET,
			"rules.0.status":                             CHECKSET,
			"rules.0.compliance.#":                       CHECKSET,
			"rules.0.description":                        fmt.Sprintf("tf-testAccConfigRule%d", rand),
			"rules.0.input_parameters.%":                 "1",
			"rules.0.modified_timestamp":                 CHECKSET,
			"rules.0.risk_level":                         "1",
			"rules.0.rule_name":                          fmt.Sprintf("tf-testAccConfigRule%d", rand),
			"rules.0.event_source":                       "aliyun.config",
			"rules.0.source_maximum_execution_frequency": "",
			"rules.0.scope_compliance_resource_types.#":  "1",
			"rules.0.source_detail_message_type":         "ConfigurationItemChangeNotification",
			"rules.0.source_identifier":                  "ecs-instances-in-vpc",
			"rules.0.source_owner":                       "ALIYUN",
			"rules.0.tag_key_scope":                      "tfTest",
			"rules.0.tag_value_scope":                    "tfTest 123",
			"rules.0.resource_types_scope.#":             "1",
			"rules.0.resource_group_ids_scope":           CHECKSET,
			"rules.0.region_ids_scope":                   "cn-hangzhou",
			"rules.0.maximum_execution_frequency":        "",
			"rules.0.exclude_resource_ids_scope":         CHECKSET,
			"rules.0.config_rule_trigger_types":          "ConfigurationItemChangeNotification",
			"rules.0.compliance_pack_id":                 "",
		}
	}

	var fakeConfigRulesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var rolesRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existConfigRulesRecordsMapFunc,
		fakeMapFunc:  fakeConfigRulesRecordsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
	}

	rolesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, riskLevelConf, allConf)

}

func testAccCheckAlicloudConfigRulesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccConfigRule%d"
}

resource "alicloud_resource_manager_resource_group" "example" {
  count = 2
  resource_group_name = join("-", [var.name, count.index])
  display_name        = join("-", [var.name, count.index])
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
   availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count = 1
  memory_size = 2
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu"
  most_recent = true
  owners = "system"
}

resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
  name = "${var.name}"
  description = "foo"
  vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  image_id = "${data.alicloud_images.default.images.0.id}"
  # series III
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
  instance_name = "${var.name}"
  user_data = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

data "alicloud_instances" "default" {
 name_regex = "${alicloud_instance.default.instance_name}"
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "${alicloud_instance.default.instance_name}"
}

resource "alicloud_config_rule" "example" {
  rule_name                       = var.name
  description                     = var.name
  source_identifier               = "ecs-instances-in-vpc"
  source_owner                    = "ALIYUN"
  resource_types_scope 			  = ["ACS::ECS::Instance"]
  risk_level                      = 1
  config_rule_trigger_types       = "ConfigurationItemChangeNotification"
  tag_key_scope 				  = "tfTest"
  tag_value_scope 				  = "tfTest 123"
  resource_group_ids_scope 		  = data.alicloud_resource_manager_resource_groups.default.ids.0
  exclude_resource_ids_scope      = data.alicloud_instances.default.instances[0].id
  region_ids_scope 				  = "cn-hangzhou"
  input_parameters  = {
		vpcIds= data.alicloud_instances.default.instances[0].vpc_id
  }
}

data "alicloud_config_rules" "example"{
 enable_details = true
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
