// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSlsMachineGroupDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	ProjectNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlsMachineGroupSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_machine_group.default.id}"]`,
			"project_name": `"${var.project_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlsMachineGroupSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_machine_group.default.id}_fake"]`,
			"project_name": `"${var.project_name}"`,
		}),
	}
	GroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlsMachineGroupSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_machine_group.default.id}"]`,
			"project_name": `"${var.project_name}"`,
			"group_name":   `"group1"`,
		}),
		fakeConfig: testAccCheckAlicloudSlsMachineGroupSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_machine_group.default.id}_fake"]`,
			"project_name": `"${var.project_name}"`,
			"group_name":   `"group1_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlsMachineGroupSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_machine_group.default.id}"]`,
			"project_name": `"${var.project_name}"`,
			"group_name":   `"group1"`,
		}),
		fakeConfig: testAccCheckAlicloudSlsMachineGroupSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_machine_group.default.id}_fake"]`,
			"project_name": `"${var.project_name}"`,
			"group_name":   `"group1_fake"`,
		}),
	}

	SlsMachineGroupCheckInfo.dataSourceTestCheck(t, rand, ProjectNameConf, GroupNameConf, allConf)
}

var existSlsMachineGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":            "1",
		"groups.0.group_name": CHECKSET,
	}
}

var fakeSlsMachineGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var SlsMachineGroupCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_sls_machine_groups.default",
	existMapFunc: existSlsMachineGroupMapFunc,
	fakeMapFunc:  fakeSlsMachineGroupMapFunc,
}

func testAccCheckAlicloudSlsMachineGroupSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlsMachineGroup%d"
}
variable "project_name" {
  default = "project-for-machine-group-terraform"
}

resource "alicloud_log_project" "defaultyJqrue" {
  description = "for terraform test"
  name        = var.project_name
}

resource "alicloud_sls_machine_group" "default2" {
  group_name            = "group2"
  project_name          = alicloud_log_project.defaultyJqrue.project_name
  machine_identify_type = "ip"
  group_attribute {
    group_topic   = "test"
    external_name = "test"
  }
  machine_list = ["192.168.1.1"]
}


resource "alicloud_sls_machine_group" "default3" {
  group_name            = "group3"
  project_name          = alicloud_sls_machine_group.default2.project_name
  machine_identify_type = "ip"
  group_attribute {
    group_topic   = "test"
    external_name = "test"
  }
  machine_list = ["192.168.1.1"]
}

resource "alicloud_sls_machine_group" "default" {
  group_name            = "group1"
  project_name          = alicloud_sls_machine_group.default3.project_name
  machine_identify_type = "ip"
  group_attribute {
    group_topic   = "test"
    external_name = "test"
  }
  machine_list = ["192.168.1.1"]
}

data "alicloud_sls_machine_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
