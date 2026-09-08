package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsEscalationPoliciesDataSource_nameRegex(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEscalationPoliciesDataSourceName(rand, map[string]string{
			"workspace":   `"${alicloud_cms_workspace.default.workspace_name}"`,
			"name_regex":  `"${alicloud_cms_escalation_policy.default.name}"`,
			"output_file": `"/tmp/${alicloud_cms_escalation_policy.default.name}-escalation-policies-output.txt"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEscalationPoliciesDataSourceName(rand, map[string]string{
			"workspace":  `"${alicloud_cms_workspace.default.workspace_name}"`,
			"name_regex": `"${alicloud_cms_escalation_policy.default.name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEscalationPoliciesDataSourceName(rand, map[string]string{
			"workspace": `"${alicloud_cms_workspace.default.workspace_name}"`,
			"ids":       `["${alicloud_cms_escalation_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEscalationPoliciesDataSourceName(rand, map[string]string{
			"workspace": `"${alicloud_cms_workspace.default.workspace_name}"`,
			"ids":       `["${alicloud_cms_escalation_policy.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEscalationPoliciesDataSourceName(rand, map[string]string{
			"workspace":  `"${alicloud_cms_workspace.default.workspace_name}"`,
			"name_regex": `"${alicloud_cms_escalation_policy.default.name}"`,
			"ids":        `["${alicloud_cms_escalation_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEscalationPoliciesDataSourceName(rand, map[string]string{
			"workspace":  `"${alicloud_cms_workspace.default.workspace_name}"`,
			"name_regex": `"${alicloud_cms_escalation_policy.default.name}_fake"`,
			"ids":        `["${alicloud_cms_escalation_policy.default.id}_fake"]`,
		}),
	}
	var existAlicloudCmsEscalationPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"policies.#":             "1",
			"policies.0.id":          CHECKSET,
			"policies.0.uuid":        CHECKSET,
			"policies.0.name":        CHECKSET,
			"policies.0.workspace":   CHECKSET,
			"policies.0.create_time": CHECKSET,
			"policies.0.update_time": CHECKSET,
		}
	}
	var fakeAlicloudCmsEscalationPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}
	var alicloudCmsEscalationPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_escalation_policies.default",
		existMapFunc: existAlicloudCmsEscalationPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsEscalationPoliciesDataSourceNameMapFunc,
	}
	alicloudCmsEscalationPoliciesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCmsEscalationPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	name := fmt.Sprintf("tf-data-escalation-%d", rand)
	config := fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = var.name
}

resource "alicloud_cms_escalation_policy" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  name        = var.name
  enable      = true
  description = "Terraform escalation policy for datasource test"
  escalation_stage_list {
    index                  = 1
    cycle_notify_interval = 5
    cycle_notify_count    = 3
    trigger_delay         = 10
    target_incident_state = "acknowledged"
    notify_channels {
      channel_type         = "DING"
      receivers            = ["tf-test-group"]
      enabled_sub_channels = []
    }
    effect_time_range {
      time_zone            = "Asia/Shanghai"
      start_time_in_minute = 0
      end_time_in_minute   = 1439
      day_in_week          = [1, 2, 3, 4, 5]
    }
  }
}

data "alicloud_cms_escalation_policies" "default" {
`, name)

	for k, v := range attrMap {
		config += fmt.Sprintf("\n  %s = %s", k, v)
	}
	config += "\n}"
	return config
}

func TestAccAliCloudCmsEscalationPoliciesDataSource_empty(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	emptyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEscalationPoliciesDataSourceEmpty(rand),
		fakeConfig:  testAccCheckAlicloudCmsEscalationPoliciesDataSourceEmpty(rand),
	}
	var existAlicloudCmsEscalationPoliciesDataSourceEmptyMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET,
			"policies.#": CHECKSET,
		}
	}
	var fakeAlicloudCmsEscalationPoliciesDataSourceEmptyMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET,
			"policies.#": CHECKSET,
		}
	}
	var alicloudCmsEscalationPoliciesEmptyCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_escalation_policies.default",
		existMapFunc: existAlicloudCmsEscalationPoliciesDataSourceEmptyMapFunc,
		fakeMapFunc:  fakeAlicloudCmsEscalationPoliciesDataSourceEmptyMapFunc,
	}
	alicloudCmsEscalationPoliciesEmptyCheckInfo.dataSourceTestCheck(t, rand, emptyConf)
}

func testAccCheckAlicloudCmsEscalationPoliciesDataSourceEmpty(rand int) string {
	name := fmt.Sprintf("tf-data-escalation-empty-%d", rand)
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "%s"
}

data "alicloud_cms_escalation_policies" "default" {
  name_regex = "tf-data-escalation-empty-non-existent"
}
`, name)
}
