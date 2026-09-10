package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudArmsEnvServiceMonitorsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_arms_env_service_monitors.default"
	name := fmt.Sprintf("tf-testacc%sarmsenvservicemonitor%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsEnvServiceMonitorsConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_service_monitor.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_service_monitor.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_service_monitor.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_service_monitor.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_service_monitor.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_service_monitor.default.env_service_monitor_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_service_monitor.default.environment_id}",
			"name_regex":     "${alicloud_arms_env_service_monitor.default.env_service_monitor_name}_fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_service_monitor.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_service_monitor.default.id}"},
			"name_regex":     "${alicloud_arms_env_service_monitor.default.env_service_monitor_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_id": "${alicloud_arms_env_service_monitor.default.environment_id}",
			"ids":            []string{"${alicloud_arms_env_service_monitor.default.id}_fake"},
			"name_regex":     "${alicloud_arms_env_service_monitor.default.env_service_monitor_name}_fake",
		}),
	}
	var existAliCloudArmsEnvServiceMonitorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"monitors.#":                          "1",
			"monitors.0.id":                       CHECKSET,
			"monitors.0.config_yaml":              CHECKSET,
			"monitors.0.env_service_monitor_name": CHECKSET,
			"monitors.0.environment_id":           CHECKSET,
			"monitors.0.namespace":                CHECKSET,
			"monitors.0.region_id":                CHECKSET,
			"monitors.0.status":                   CHECKSET,
		}
	}
	var fakeAliCloudArmsEnvServiceMonitorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"monitors.#": "0",
		}
	}
	var alicloudArmsEnvServiceMonitorsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_env_service_monitors.default",
		existMapFunc: existAliCloudArmsEnvServiceMonitorsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsEnvServiceMonitorsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsEnvServiceMonitorsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceArmsEnvServiceMonitorsConfig0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_enhanced_nat_available_zones" "enhanced" {
	}

	resource "alicloud_vpc" "vpc" {
  		description = "api-resource-test1-hz"
  		cidr_block  = "192.168.0.0/16"
  		vpc_name    = var.name
	}

	resource "alicloud_vswitch" "vswitch" {
  		description  = "api-resource-test1-hz"
  		vpc_id       = alicloud_vpc.vpc.id
  		vswitch_name = var.name
  		zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  		cidr_block = cidrsubnet(alicloud_vpc.vpc.cidr_block, 8, 8)
	}


	resource "alicloud_snapshot_policy" "default" {
  		name            = var.name
  		repeat_weekdays = ["1", "2", "3"]
  		retention_days  = -1
  		time_points     = ["1", "22", "23"]
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = alicloud_vswitch.vswitch.zone_id
  		cpu_core_count       = 2
  		memory_size          = 4
  		kubernetes_node_role = "Worker"
  		instance_type_family = "ecs.sn1ne"
	}

	resource "alicloud_cs_managed_kubernetes" "default" {
  		name               = var.name
  		cluster_spec       = "ack.pro.small"
  		version            = "1.24.6-aliyun.1"
  		new_nat_gateway    = true
  		node_cidr_mask     = 26
  		proxy_mode         = "ipvs"
  		service_cidr       = "172.23.0.0/16"
  		pod_cidr           = "10.95.0.0/16"
  		worker_vswitch_ids = [alicloud_vswitch.vswitch.id]
	}

	resource "alicloud_key_pair" "default" {
  		key_pair_name = var.name
	}

	resource "alicloud_cs_kubernetes_node_pool" "default" {
  		name                 = "desired_size"
  		cluster_id           = alicloud_cs_managed_kubernetes.default.id
  		vswitch_ids          = [alicloud_vswitch.vswitch.id]
  		instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  		system_disk_category = "cloud_efficiency"
  		system_disk_size     = 40
  		key_name             = alicloud_key_pair.default.key_pair_name
  		desired_size         = 2
	}

	resource "alicloud_arms_environment" "default" {
  		environment_type     = "CS"
  		environment_name     = var.name
  		bind_resource_id     = alicloud_cs_kubernetes_node_pool.default.cluster_id
  		environment_sub_type = "ManagedKubernetes"
	}

	resource "alicloud_arms_env_service_monitor" "default" {
  		aliyun_lang    = "en"
  		environment_id = alicloud_arms_environment.default.id
  		config_yaml    = <<EOF
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: arms-admin1
  namespace: arms-prom
  annotations:
    arms.prometheus.io/discovery: 'true'
    o11y.aliyun.com/addon-name: mysql
    o11y.aliyun.com/addon-version: 1.0.1
    o11y.aliyun.com/release-name: mysql1
spec:
  endpoints:
  - interval: 30s
    port: operator
    path: /metrics
  - interval: 10s
    port: operator1
    path: /metrics
  namespaceSelector:
    any: true
  selector:
    matchLabels:
     app: arms-prometheus-ack-arms-prometheus
EOF
	}
`, name)
}
