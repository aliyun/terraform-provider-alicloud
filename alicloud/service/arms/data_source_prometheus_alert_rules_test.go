package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudARMSPrometheusAlertRulesDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_arms_prometheus_alert_rules.default"
	name := fmt.Sprintf("tf-testacc-ArmsPrometheusAlertRules%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsPrometheusAlertRulesConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"name_regex": "${alicloud_arms_prometheus_alert_rule.default.prometheus_alert_rule_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"name_regex": "${alicloud_arms_prometheus_alert_rule.default.prometheus_alert_rule_name}_fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}_fake"},
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
			"status":     "1",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}_fake"},
			"status":     "0",
		}),
	}
	typeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
			"type":       "${alicloud_arms_prometheus_alert_rule.default.type}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":        []string{"${alicloud_arms_prometheus_alert_rule.default.id}_fake"},
			"type":       "${alicloud_arms_prometheus_alert_rule.default.type}_fake",
		}),
	}
	match_expressions := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id":        "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":               []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
			"match_expressions": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id":        "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"ids":               []string{"${alicloud_arms_prometheus_alert_rule.default.id}_fake"},
			"match_expressions": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 20",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id":        "${data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id}",
			"name_regex":        "${alicloud_arms_prometheus_alert_rule.default.prometheus_alert_rule_name}",
			"ids":               []string{"${alicloud_arms_prometheus_alert_rule.default.id}"},
			"status":            "1",
			"type":              "${alicloud_arms_prometheus_alert_rule.default.type}",
			"match_expressions": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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

	var ArmsPrometheusAlertRulesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existArmsPrometheusAlertRuleMapFunc,
		FakeMapFunc:  fakeArmsPrometheusAlertRuleMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckPrePaidResources(t)
	}
	ArmsPrometheusAlertRulesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, match_expressions, statusConf, typeConf, allConf)

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
