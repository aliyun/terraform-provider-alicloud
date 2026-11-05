// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsEventNotifyPolicyDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}_fake"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
	}

	WorkspaceConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}_fake"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
			"name":      `"${var.name}_fake"`,
		}),
	}
	NameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}"]`,
			"name":      `"${var.name}"`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}_fake"]`,
			"name":      `"${var.name}_fake"`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,

			"name": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsEventNotifyPolicySourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_event_notify_policy.default.id}_fake"]`,
			"workspace": `"${alicloud_cms_workspace.default.id}"`,
			"name":      `"${var.name}_fake"`,
		}),
	}

	CmsEventNotifyPolicyCheckInfo.DataSourceTestCheck(t, rand, idsConf, WorkspaceConf, NameConf, allConf)
}

var existCmsEventNotifyPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#":                   "1",
		"policies.0.description":       acctest.CHECKSET,
		"policies.0.create_time":       acctest.CHECKSET,
		"policies.0.enabled":           acctest.CHECKSET,
		"policies.0.notify_strategy.#": acctest.CHECKSET,
		"policies.0.name":              acctest.CHECKSET,
		"policies.0.uuid":              acctest.CHECKSET,
		"policies.0.version":           acctest.CHECKSET,
		"policies.0.user_id":           acctest.CHECKSET,
		"policies.0.update_time":       acctest.CHECKSET,
		"policies.0.workspace":         acctest.CHECKSET,
	}
}

var fakeCmsEventNotifyPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#": "0",
	}
}

var CmsEventNotifyPolicyCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cms_event_notify_policies.default",
	ExistMapFunc: existCmsEventNotifyPolicyMapFunc,
	FakeMapFunc:  fakeCmsEventNotifyPolicyMapFunc,
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
