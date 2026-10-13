// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudCmsEventNotifyPolicyDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}_fake"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
	}

	WorkspaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}_fake"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
			"name":      `"${var.name}_fake"`,
		}),
	}
	NameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}"]`,
			"name":      `"${var.name}"`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}_fake"]`,
			"name":      `"${var.name}_fake"`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,

			"name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}_fake"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
			"name":      `"${var.name}_fake"`,
		}),
	}

	CmsEventNotifyPolicyCheckInfo.dataSourceTestCheck(t, rand, idsConf, WorkspaceConf, NameConf, allConf)
}

var existCmsEventNotifyPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#":                   "1",
		"policies.0.description":       CHECKSET,
		"policies.0.create_time":       CHECKSET,
		"policies.0.enabled":           CHECKSET,
		"policies.0.notify_strategy.#": CHECKSET,
		"policies.0.name":              CHECKSET,
		"policies.0.uuid":              CHECKSET,
		"policies.0.version":           CHECKSET,
		"policies.0.user_id":           CHECKSET,
		"policies.0.update_time":       CHECKSET,
		"policies.0.workspace":         CHECKSET,
	}
}

var fakeCmsEventNotifyPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#": "0",
	}
}

var CmsEventNotifyPolicyCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cms_event_notify_policies.default",
	existMapFunc: existCmsEventNotifyPolicyMapFunc,
	fakeMapFunc:  fakeCmsEventNotifyPolicyMapFunc,
}

func testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testacccmsenp%d"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_event_notify_policy" "default" {
  description = "Event notification policy managed by Terraform"
  response_plan {
    repeat_notify_setting {
      end_incident_state = "resolved"
      repeat_interval    = "30"
    }
    auto_recover_seconds = "600"
  }
  notify_strategy {
    description                  = "Notify strategy for list test"
    ignore_restored_notification = false
    grouping_setting {
      grouping_keys = ["severity"]
      period_min    = "5"
      times         = "1"
      silence_sec   = "300"
    }
    routes {
      effect_time_range {
        time_zone            = "Asia/Shanghai"
        start_time_in_minute = "0"
        end_time_in_minute   = "1439"
        day_in_week          = ["1", "2", "3", "4", "5"]
      }
      channels {
        receivers            = ["tf-test-group"]
        channel_type         = "DING"
        enabled_sub_channels = []
      }
    }
  }
  subscription {
    subscribe_legacy_event = true
    filter_setting {
      relation = "AND"
      conditions {
        field = "severity"
        op    = "EQ"
        value = "CRITICAL"
      }
    }
  }
  workspace = alicloud_cms_workspace.default.id
  name      = var.name
}

data "alicloud_cms_event_notify_policies" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
