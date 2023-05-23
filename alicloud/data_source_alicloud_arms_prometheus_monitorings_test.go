package alicloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudArmsPrometheusMonitoringsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_arms_prometheus_monitorings.default"
	name := fmt.Sprintf("tf-testacc-ArmsPM%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsPrometheusMonitoringsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"ids":        []string{"${alicloud_arms_prometheus_monitoring.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"ids":        []string{"${alicloud_arms_prometheus_monitoring.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"name_regex": "${alicloud_arms_prometheus_monitoring.default.monitoring_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"name_regex": "${alicloud_arms_prometheus_monitoring.default.monitoring_name}_fake",
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"type":       "customJob",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"type":       "probe",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"status":     "run",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"status":     "stop",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"ids":        []string{"${alicloud_arms_prometheus_monitoring.default.id}"},
			"name_regex": "${alicloud_arms_prometheus_monitoring.default.monitoring_name}",
			"type":       "customJob",
			"status":     "run",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"ids":        []string{"${alicloud_arms_prometheus_monitoring.default.id}_fake"},
			"name_regex": "${alicloud_arms_prometheus_monitoring.default.monitoring_name}_fake",
			"type":       "probe",
			"status":     "stop",
		}),
	}
	var existAliCloudArmsPrometheusMonitoringsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"names.#":                                  "1",
			"prometheus_monitorings.#":                 "1",
			"prometheus_monitorings.0.id":              CHECKSET,
			"prometheus_monitorings.0.cluster_id":      CHECKSET,
			"prometheus_monitorings.0.monitoring_name": CHECKSET,
			"prometheus_monitorings.0.type":            "customJob",
			"prometheus_monitorings.0.config_yaml":     CHECKSET,
			"prometheus_monitorings.0.status":          "run",
		}
	}
	var fakeAliCloudArmsPrometheusMonitoringsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "0",
			"names.#":                  "0",
			"prometheus_monitorings.#": "0",
		}
	}
	var alicloudArmsPrometheusMonitoringsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_prometheus_monitorings.default",
		existMapFunc: existAliCloudArmsPrometheusMonitoringsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsPrometheusMonitoringsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsPrometheusMonitoringsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, typeConf, statusConf, allConf)
}

func dataSourceArmsPrometheusMonitoringsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_arms_prometheus" "default" {
  		cluster_type        = "ecs"
  		grafana_instance_id = "free"
  		vpc_id              = data.alicloud_vpcs.default.ids.0
  		vswitch_id          = data.alicloud_vswitches.default.ids.0
  		security_group_id   = alicloud_security_group.default.id
  		cluster_name        = "${var.name}-${data.alicloud_vpcs.default.ids.0}"
  		resource_group_id   = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	}

	resource "alicloud_arms_prometheus_monitoring" "default" {
  		cluster_id  = alicloud_arms_prometheus.default.id
  		type        = "customJob"
  		config_yaml = "scrape_configs:\n- job_name: prometheus\n  honor_timestamps: false\n  honor_labels: false\n  scheme: http\n  metrics_path: /metric\n  static_configs:\n  - targets:\n    - 127.0.0.1:9090\n"
  		status      = "run"
	}
`, name)
}
