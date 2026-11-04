package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudArmsPrometheusMonitoringsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_arms_prometheus_monitorings.default"
	name := fmt.Sprintf("tf-testacc-ArmsPM%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsPrometheusMonitoringsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"ids":        []string{"${alicloud_arms_prometheus_monitoring.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"ids":        []string{"${alicloud_arms_prometheus_monitoring.default.id}_fake"},
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"name_regex": "${alicloud_arms_prometheus_monitoring.default.monitoring_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"name_regex": "${alicloud_arms_prometheus_monitoring.default.monitoring_name}_fake",
		}),
	}
	typeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"type":       "customJob",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"type":       "probe",
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"status":     "run",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"status":     "stop",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"cluster_id": "${alicloud_arms_prometheus_monitoring.default.cluster_id}",
			"ids":        []string{"${alicloud_arms_prometheus_monitoring.default.id}"},
			"name_regex": "${alicloud_arms_prometheus_monitoring.default.monitoring_name}",
			"type":       "customJob",
			"status":     "run",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"prometheus_monitorings.0.id":              acctest.CHECKSET,
			"prometheus_monitorings.0.cluster_id":      acctest.CHECKSET,
			"prometheus_monitorings.0.monitoring_name": acctest.CHECKSET,
			"prometheus_monitorings.0.type":            "customJob",
			"prometheus_monitorings.0.config_yaml":     acctest.CHECKSET,
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
	var alicloudArmsPrometheusMonitoringsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_arms_prometheus_monitorings.default",
		ExistMapFunc: existAliCloudArmsPrometheusMonitoringsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudArmsPrometheusMonitoringsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudArmsPrometheusMonitoringsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, typeConf, statusConf, allConf)
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
