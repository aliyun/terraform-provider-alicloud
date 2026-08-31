package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Arms PrometheusMonitoring. >>> Resource test cases, automatically generated.
// Case 3657
func TestAccAlicloudArmsPrometheusMonitoring_basic3657(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_monitoring.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsPrometheusMonitoringMap3657)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsPrometheusMonitoring")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusmonitoring%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsPrometheusMonitoringBasicDependence3657)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":      "run",
					"type":        "serviceMonitor",
					"cluster_id":  "${alicloud_arms_prometheus.default.cluster_id}",
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: tomcat-demo\n  namespace: default\nspec:\n  endpoints:\n  - bearerTokenSecret:\n      key: ''\n    interval: 30s\n    path: /metrics\n    port: tomcat-monitor\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n      app: tomcat`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":      "run",
						"type":        "serviceMonitor",
						"cluster_id":  CHECKSET,
						"config_yaml": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "run",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "run",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: tomcat-demo\n  namespace: default\nspec:\n  endpoints:\n  - bearerTokenSecret:\n      key: ''\n    interval: 30s\n    path: /metrics\n    port: tomcat-monitor\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n      app: tomcat`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "stop",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "stop",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: tomcat-demo\n  namespace: default\nspec:\n  endpoints:\n  - bearerTokenSecret:\n      key: ''\n    interval: 31s\n    path: /metrics\n    port: tomcat-monitor\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n      app: tomcat`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":      "run",
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: tomcat-demo\n  namespace: default\nspec:\n  endpoints:\n  - bearerTokenSecret:\n      key: ''\n    interval: 30s\n    path: /metrics\n    port: tomcat-monitor\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n      app: tomcat`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":      "run",
						"config_yaml": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudArmsPrometheusMonitoringMap3657 = map[string]string{
	"monitoring_name": CHECKSET,
}

func AlicloudArmsPrometheusMonitoringBasicDependence3657(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
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
  key_name             = alicloud_key_pair.default.key_name
  desired_size         = 2
}

resource "alicloud_arms_prometheus" "default" {
  cluster_type        = "aliyun-cs"
  grafana_instance_id = "free"
  cluster_id          = alicloud_cs_kubernetes_node_pool.default.cluster_id
}

`, name)
}

// Case 3657  twin
func TestAccAlicloudArmsPrometheusMonitoring_basic3657_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_prometheus_monitoring.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsPrometheusMonitoringMap3657)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsPrometheusMonitoring")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsprometheusmonitoring%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsPrometheusMonitoringBasicDependence3657)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":      "stop",
					"type":        "serviceMonitor",
					"cluster_id":  "${alicloud_arms_prometheus.default.cluster_id}",
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: tomcat-demo\n  namespace: default\nspec:\n  endpoints:\n  - bearerTokenSecret:\n      key: ''\n    interval: 31s\n    path: /metrics\n    port: tomcat-monitor\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n      app: tomcat`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":      "stop",
						"type":        "serviceMonitor",
						"cluster_id":  CHECKSET,
						"config_yaml": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Arms PrometheusMonitoring. <<< Resource test cases, automatically generated.
