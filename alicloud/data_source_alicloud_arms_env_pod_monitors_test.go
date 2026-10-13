package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudArmsEnvPodMonitorsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_arms_env_pod_monitors.default"
	name := fmt.Sprintf("tf-testacc%sarmsenvpodmonitor%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsEnvPodMonitorsConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_pod_monitor.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_pod_monitor.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_pod_monitor.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_pod_monitor.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_pod_monitor.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_pod_monitor.default.env_pod_monitor_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_pod_monitor.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_pod_monitor.default.env_pod_monitor_name}_fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_pod_monitor.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_pod_monitor.default.id}"},
			"name_regex":     "${alicloud_arms_env_pod_monitor.default.env_pod_monitor_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_pod_monitor.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_pod_monitor.default.id}_fake"},
			"name_regex":     "${alicloud_arms_env_pod_monitor.default.env_pod_monitor_name}_fake",
		}),
	}
	var existAliCloudArmsEnvPodMonitorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"monitors.#":                      "1",
			"monitors.0.id":                   CHECKSET,
			"monitors.0.config_yaml":          CHECKSET,
			"monitors.0.env_pod_monitor_name": CHECKSET,
			"monitors.0.environment_id":       CHECKSET,
			"monitors.0.namespace":            CHECKSET,
			"monitors.0.region_id":            CHECKSET,
			"monitors.0.status":               CHECKSET,
		}
	}
	var fakeAliCloudArmsEnvPodMonitorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"monitors.#": "0",
		}
	}
	var alicloudArmsEnvPodMonitorsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_env_pod_monitors.default",
		existMapFunc: existAliCloudArmsEnvPodMonitorsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsEnvPodMonitorsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsEnvPodMonitorsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceArmsEnvPodMonitorsConfig0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	data "alicloud_cs_managed_kubernetes_clusters" "default" {
  		name_regex = "^Default"
	}

	resource "alicloud_cs_managed_kubernetes" "default" {
  		count                = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? 0 : 1
  		name                 = var.name
  		cluster_spec         = "ack.pro.small"
  		worker_vswitch_ids   = [data.alicloud_vswitches.default.ids.0]
  		new_nat_gateway      = false
  		pod_cidr             = "10.132.0.0/16"
  		service_cidr         = "192.168.0.0/16"
  		slb_internet_enabled = true
  		is_enterprise_security_group = true
	}

	locals {
  		cluster_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
	}

	resource "alicloud_arms_environment" "environment-cs" {
  		environment_type     = "CS"
  		environment_name     = var.name
  		bind_resource_id     = local.cluster_id
  		environment_sub_type = "ManagedKubernetes"
	}

	resource "alicloud_arms_env_pod_monitor" "default" {
  		aliyun_lang    = "en"
  		environment_id = alicloud_arms_environment.environment-cs.id
  		config_yaml    = <<EOF
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: tf-test-pm1
  namespace: arms-prom
  annotations:
    arms.prometheus.io/discovery: 'true'
    o11y.aliyun.com/addon-name: mysql
    o11y.aliyun.com/addon-version: 1.0.2
    o11y.aliyun.com/release-name: mysql2
spec:
  selector:
    matchLabels:
      app: arms-prometheus-ack-arms-prometheus
      release: arms-prometheus
  namespaceSelector:
    any: true    
  podMetricsEndpoints:
  - interval: 30s
    targetPort: 9335
    path: /metrics
  - interval: 11s
    targetPort: 9335
    path: /metric
EOF
	}
`, name)
}
