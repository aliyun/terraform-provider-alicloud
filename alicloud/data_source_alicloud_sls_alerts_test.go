// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSlsAlertDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	ProjectNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlsAlertSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_alert.default.id}"]`,
			"project_name": `"${alicloud_log_project.defaultINsMgl.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlsAlertSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_alert.default.id}_fake"]`,
			"project_name": `"${alicloud_log_project.defaultINsMgl.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlsAlertSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_alert.default.id}"]`,
			"project_name": `"${alicloud_log_project.defaultINsMgl.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlsAlertSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_sls_alert.default.id}_fake"]`,
			"project_name": `"${alicloud_log_project.defaultINsMgl.id}"`,
		}),
	}

	SlsAlertCheckInfo.dataSourceTestCheck(t, rand, ProjectNameConf, allConf)
}

var existSlsAlertMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"alerts.#":                 "1",
		"alerts.0.description":     CHECKSET,
		"alerts.0.alert_name":      CHECKSET,
		"alerts.0.configuration.#": CHECKSET,
		"alerts.0.schedule.#":      CHECKSET,
		"alerts.0.display_name":    CHECKSET,
	}
}

var fakeSlsAlertMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"alerts.#": "0",
	}
}

var SlsAlertCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_sls_alerts.default",
	existMapFunc: existSlsAlertMapFunc,
	fakeMapFunc:  fakeSlsAlertMapFunc,
}

func testAccCheckAlicloudSlsAlertSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-test%d"
}
variable "alert_name" {
  default = "openapi-terraform-alert"
}

variable "project_name" {
  default = "terraform-alert-test"
}

resource "alicloud_log_project" "defaultINsMgl" {
  description = "terraform-alert-test"
  name        = var.name
}



resource "alicloud_sls_alert" "default" {
  configuration {
    type    = "tpl"
    version = "2"
    query_list {
      query          = "* | select *"
      time_span_type = "Relative"
      start          = "-15m"
      end            = "now"
      store_type     = "log"
      project        = alicloud_log_project.defaultINsMgl.id
      store          = "alert"
      region         = "cn-hangzhou"
      power_sql_mode = "disable"
      chart_title    = "wkb-chart"
      dashboard_id   = "wkb-dashboard"
      ui             = "{}"
      role_arn       = "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole"
    }
    query_list {
      store_type = "meta"
      store      = "user.rds_ip_whitelist"
    }
    query_list {
      store_type = "meta"
      store      = "mytest1"
    }
    group_configuration {
      type   = "no_group"
      fields = ["a", "b"]
    }
    join_configurations {
      type      = "no_join"
      condition = "aa"
    }
    join_configurations {
      type      = "cross_join"
      condition = "qqq"
    }
    join_configurations {
      type      = "inner_join"
      condition = "fefefe"
    }
    severity_configurations {
      severity = "6"
      eval_condition {
        condition       = "__count__ > 1"
        count_condition = "cnt > 0"
      }
    }
    labels {
      key   = "a"
      value = "b"
    }
    annotations {
      key   = "x"
      value = "y"
    }
    auto_annotation = true
    send_resolved   = false
    threshold       = "1"
    no_data_fire    = false
    sink_event_store {
      enabled     = true
      endpoint    = "cn-shanghai-intranet.log.aliyuncs.com"
      project     = "wkb-wangren"
      event_store = "alert"
      role_arn    = "acs:ram::1654218965343050:role/aliyunlogetlrole"
    }
    sink_cms {
      enabled = false
    }
    sink_alerthub {
      enabled = false
    }
    template_configuration {
      template_id = "sls.app.ack.autoscaler.cluster_unhealthy"
      type        = "sys"
      version     = "1.0"
      lang        = "cn"
    }
    condition_configuration {
      condition       = "cnt > 3"
      count_condition = "__count__ < 3"
    }
    policy_configuration {
      alert_policy_id  = "sls.builtin.dynamic"
      action_policy_id = "wkb-action"
      repeat_interval  = "1m"
    }
    dashboard        = "internal-alert"
    mute_until       = "0"
    no_data_severity = "6"
    tags             = ["wkb", "wangren", "sls"]
  }
  alert_name   = var.alert_name
  project_name = alicloud_log_project.defaultINsMgl.id
  schedule {
    type            = "Cron"
    run_immdiately  = true
    time_zone       = "+0800"
    delay           = "10"
    cron_expression = "0/5 * * * *"
  }
  display_name = "openapi-terraform"
  description  = "create alert"
}

data "alicloud_sls_alerts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
