package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudARMSPrometheusAlertRulesDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_arms_prometheus_alert_rules.default"
	name := fmt.Sprintf("tf-testacc-ArmsPrometheusAlertRules%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsPrometheusAlertRulesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"name_regex": "${alicloud_arms_prometheus_alert_rule.default.prometheus_alert_rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"name_regex": "${alicloud_arms_prometheus_alert_rule.default.prometheus_alert_rule_name}_fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}_fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
			"status":     "1",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}_fake"},
			"status":     "0",
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
			"type":       "${alicloud_arms_prometheus_alert_rule.default.type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}_fake"},
			"type":       "${alicloud_arms_prometheus_alert_rule.default.type}_fake",
		}),
	}
	match_expressions := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id":        "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":               []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
			"match_expressions": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id":        "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":               []string{"${alicloud_arms_prometheus_alert_rule.default.id}_fake"},
			"match_expressions": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 20",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id":        "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"name_regex":        "${alicloud_arms_prometheus_alert_rule.default.prometheus_alert_rule_name}",
			"ids":               []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
			"status":            "1",
			"type":              "${alicloud_arms_prometheus_alert_rule.default.type}",
			"match_expressions": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id":        "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"name_regex":        "${alicloud_arms_prometheus_alert_rule.default.prometheus_alert_rule_name}_fake",
			"ids":               []string{"${alicloud_arms_prometheus_alert_rule.default.id}_fake"},
			"status":            "0",
			"type":              "${alicloud_arms_prometheus_alert_rule.default.type}_fake",
			"match_expressions": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 20",
		}),
	}
	var existArmsPrometheusAlertRuleMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"rules.#":                            "1",
			"rules.0.prometheus_alert_rule_name": name,
		}
	}

	var fakeArmsPrometheusAlertRuleMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"names.#": "0",
			"ids.#":   "0",
		}
	}

	var ArmsPrometheusAlertRulesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existArmsPrometheusAlertRuleMapFunc,
		fakeMapFunc:  fakeArmsPrometheusAlertRuleMapFunc,
	}
	preCheck := func() {
		testAccPreCheckPrePaidResources(t)
	}
	ArmsPrometheusAlertRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, match_expressions, statusConf, typeConf, allConf)

}

func dataSourceArmsPrometheusAlertRulesConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		data "alicloud_cs_managed_kubernetes_clusters" "default" {
		  name_regex = "Default"
		}
		resource "alicloud_arms_prometheus_alert_rule" "default" {
		  prometheus_alert_rule_name = var.name
		  cluster_id                 = data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id
		  expression                 = "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10"
		  message                    = "node available memory is less than 10"
		  duration                   = 1
		  notify_type                = "ALERT_MANAGER"
		  type                       = var.name
		  labels  {
			name = "TF"
			value = "Test"
		  }
		}
		`, name)
}
