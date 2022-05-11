package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSAEApplicationScalingRulesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 1000)
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeApplicationScalingRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_application_scaling_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSaeApplicationScalingRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_application_scaling_rule.default.id}_fake"]`,
		}),
	}
	var existAlicloudSaeApplicationScalingRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                          "1",
			"rules.#":                                        "1",
			"rules.0.id":                                     CHECKSET,
			"rules.0.app_id":                                 CHECKSET,
			"rules.0.create_time":                            CHECKSET,
			"rules.0.scaling_rule_enable":                    "true",
			"rules.0.scaling_rule_name":                      fmt.Sprintf("tftestacc%d", rand),
			"rules.0.scaling_rule_type":                      "mix",
			"rules.0.scaling_rule_metric.#":                  "1",
			"rules.0.scaling_rule_metric.0.max_replicas":     "50",
			"rules.0.scaling_rule_metric.0.min_replicas":     "3",
			"rules.0.scaling_rule_metric.0.metrics.#":        "3",
			"rules.0.scaling_rule_metric.0.metrics_status.#": "1",
			"rules.0.scaling_rule_metric.0.metrics_status.0.desired_replicas":       CHECKSET,
			"rules.0.scaling_rule_metric.0.metrics_status.0.next_scale_time_period": CHECKSET,
			"rules.0.scaling_rule_metric.0.metrics_status.0.current_replicas":       CHECKSET,
			"rules.0.scaling_rule_metric.0.metrics_status.0.max_replicas":           CHECKSET,
			"rules.0.scaling_rule_metric.0.metrics_status.0.min_replicas":           CHECKSET,
			"rules.0.scaling_rule_metric.0.metrics_status.0.last_scale_time":        "",
			"rules.0.scaling_rule_metric.0.metrics_status.0.current_metrics.#":      CHECKSET,
			"rules.0.scaling_rule_metric.0.metrics_status.0.next_scale_metrics.#":   CHECKSET,
			"rules.0.scaling_rule_timer.#":                                          "1",
			"rules.0.scaling_rule_timer.0.begin_date":                               "2022-02-25",
			"rules.0.scaling_rule_timer.0.end_date":                                 "2022-03-25",
			"rules.0.scaling_rule_timer.0.period":                                   "* * *",
			"rules.0.scaling_rule_timer.0.schedules.#":                              "2",
		}
	}
	var fakeAlicloudSaeApplicationScalingRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudSaeApplicationScalingRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sae_application_scaling_rules.default",
		existMapFunc: existAlicloudSaeApplicationScalingRulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSaeApplicationScalingRulesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudSaeApplicationScalingRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}
func testAccCheckAlicloudSaeApplicationScalingRulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tftestacc%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_sae_namespace" "default" {
  namespace_description = var.name
  namespace_id          = join(":",["%s",var.name])
  namespace_name        = var.name
}
resource "alicloud_sae_application" "default" {
  app_description = var.name
  app_name        = var.name
  namespace_id    = alicloud_sae_namespace.default.namespace_id
  image_url       = "registry-vpc.cn-hangzhou.aliyuncs.com/lxepoo/apache-php5"
  package_type    = "Image"
  jdk             = "Open JDK 8"
  vswitch_id      = data.alicloud_vswitches.default.ids.0
  vpc_id          = data.alicloud_vpcs.default.ids.0
  timezone        = "Asia/Shanghai"
  replicas        = "5"
  cpu             = "500"
  memory          = "2048"
}

resource "alicloud_sae_application_scaling_rule" "default" {
  app_id = alicloud_sae_application.default.id
  scaling_rule_name = var.name
  scaling_rule_type = "mix"
  min_ready_instances      = 2
  min_ready_instance_ratio = -1
  scaling_rule_enable = true
  scaling_rule_timer {
    begin_date = "2022-02-25"
    end_date   = "2022-03-25"
    period     = "* * *"
    schedules {
      at_time      = "08:00"
      max_replicas = 10
      min_replicas = 3
    }
    schedules {
      at_time      = "20:00"
      max_replicas = 50
      min_replicas = 3
    }
  }
  scaling_rule_metric {
    max_replicas = 50
    min_replicas = 3
    metrics {
      metric_type                       = "CPU"
      metric_target_average_utilization = 20
    }
    metrics {
      metric_type                       = "MEMORY"
      metric_target_average_utilization = 30
    }
    metrics {
      metric_type                       = "tcpActiveConn"
      metric_target_average_utilization = 20
    }
    scale_up_rules {
      step                         = 10
      disabled                     = false
      stabilization_window_seconds = 0
    }
    scale_down_rules {
      step                         = 10
      disabled                     = false
      stabilization_window_seconds = 0
    }
  }
}

data "alicloud_sae_application_scaling_rules" "default" {	
	app_id = alicloud_sae_application.default.id
	%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
