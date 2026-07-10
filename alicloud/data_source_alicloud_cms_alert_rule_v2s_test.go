// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudCmsAlertRuleV2DataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsAlertRuleV2SourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_alert_rule_v2.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlertRuleV2SourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_alert_rule_v2.default.id}_fake"]`,
		}),
	}

	WorkspaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsAlertRuleV2SourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_alert_rule_v2.default.id}"]`,
			"workspace": `"default-cms-1511928242963727-cn-hangzhou"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlertRuleV2SourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_alert_rule_v2.default.id}_fake"]`,
			"workspace": `"default-cms-1511928242963727-cn-hangzhou"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsAlertRuleV2SourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_alert_rule_v2.default.id}"]`,
			"workspace": `"default-cms-1511928242963727-cn-hangzhou"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlertRuleV2SourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_alert_rule_v2.default.id}_fake"]`,
			"workspace": `"default-cms-1511928242963727-cn-hangzhou"`,
		}),
	}

	CmsAlertRuleV2CheckInfo.dataSourceTestCheck(t, rand, idsConf, WorkspaceConf, allConf)
}

var existCmsAlertRuleV2MapFunc = func(rand int) map[string]string {
	return map[string]string{
		"v2s.#":                               "1",
		"v2s.0.created_at":                    CHECKSET,
		"v2s.0.observe_resource_type":         CHECKSET,
		"v2s.0.datasource_type":               CHECKSET,
		"v2s.0.datasource_config.#":           CHECKSET,
		"v2s.0.display_name":                  CHECKSET,
		"v2s.0.notify_config.#":               CHECKSET,
		"v2s.0.status":                        CHECKSET,
		"v2s.0.observe_resource_global_scope": CHECKSET,
		"v2s.0.action_integration_config.#":   CHECKSET,
		"v2s.0.severity_levels":               CHECKSET,
		"v2s.0.query_config.#":                CHECKSET,
		"v2s.0.enabled":                       CHECKSET,
		"v2s.0.alert_rule_v2_id":              CHECKSET,
		"v2s.0.updated_at":                    CHECKSET,
		"v2s.0.content_template":              CHECKSET,
		"v2s.0.schedule_config.#":             CHECKSET,
		"v2s.0.arms_integration_config.#":     CHECKSET,
		"v2s.0.condition_config.#":            CHECKSET,
		"v2s.0.workspace":                     CHECKSET,
	}
}

var fakeCmsAlertRuleV2MapFunc = func(rand int) map[string]string {
	return map[string]string{
		"v2s.#": "0",
	}
}

var CmsAlertRuleV2CheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cms_alert_rule_v2s.default",
	existMapFunc: existCmsAlertRuleV2MapFunc,
	fakeMapFunc:  fakeCmsAlertRuleV2MapFunc,
}

func testAccCheckAlicloudCmsAlertRuleV2SourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCmsAlertRuleV2%d"
}


resource "alicloud_cms_alert_rule_v2" "default" {
  content_template = "umodel test alert on $${metric}"
  schedule_config {
    type          = "FIXED"
    interval_secs = "60"
  }
  datasource_config {
    type = "UMODEL"
  }
  action_integration_config {
    enabled = false
  }
  arms_integration_config {
    enabled = false
  }
  query_config {
    entity_type   = "instance"
    type          = "UMODEL_METRICSET_QUERY"
    entity_domain = "ecs"
    metric        = "CPUUtilization"
    label_filters {
      operator = "="
      value    = "web-server"
      name     = "app"
    }
    label_filters {
      operator = "="
      value    = "production"
      name     = "env"
    }
    metric_set = "acs_ecs_dashboard"
  }
  display_name = "regression-umodel-10"
  enabled      = true
  notify_config {
    type = "DIRECT_NOTIFY"
    channels {
      type        = "GROUP"
      identifiers = ["regression-test"]
    }
  }
  condition_config {
    operator      = "GT"
    type          = "UMODEL_METRICSET_CONDITION"
    severity      = "CRITICAL"
    duration_secs = "60"
    threshold     = 90
  }
  workspace = "default-cms-1511928242963727-cn-hangzhou"
}

data "alicloud_cms_alert_rule_v2s" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
