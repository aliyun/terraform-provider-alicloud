package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Arms EnvServiceMonitor. >>> Resource test cases, automatically generated.
// Case 4551
func TestAccAliCloudArmsEnvServiceMonitor_basic4551(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_service_monitor.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvServiceMonitorMap4551)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvServiceMonitor")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvservicemonitor%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvServiceMonitorBasicDependence4551)
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
					"environment_id": "${alicloud_arms_environment.env-cs.id}",
					"config_yaml":    `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\n    o11y.aliyun.com/addon-name: mysql\n    o11y.aliyun.com/addon-version: 1.0.1\n    o11y.aliyun.com/release-name: mysql1\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 10s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus\n`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id": CHECKSET,
						"config_yaml":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\n    o11y.aliyun.com/addon-name: mysql\n    o11y.aliyun.com/addon-version: 1.0.1\n    o11y.aliyun.com/release-name: mysql1\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 10s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus\n`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\n    o11y.aliyun.com/addon-name: mysql\n    o11y.aliyun.com/addon-version: 1.0.1\n    o11y.aliyun.com/release-name: mysql1\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 11s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id": "${alicloud_arms_environment.env-cs.id}",
					"config_yaml":    `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\n    o11y.aliyun.com/addon-name: mysql\n    o11y.aliyun.com/addon-version: 1.0.1\n    o11y.aliyun.com/release-name: mysql1\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 10s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus\n`,
					"aliyun_lang":    "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id": CHECKSET,
						"config_yaml":    CHECKSET,
						"aliyun_lang":    "zh",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

var AlicloudArmsEnvServiceMonitorMap4551 = map[string]string{
	"status":                   CHECKSET,
	"env_service_monitor_name": CHECKSET,
	"namespace":                CHECKSET,
}

func AlicloudArmsEnvServiceMonitorBasicDependence4551(name string) string {
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
  key_name             = alicloud_key_pair.default.key_name
  desired_size         = 2
}

resource "alicloud_arms_environment" "env-cs" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = alicloud_cs_kubernetes_node_pool.default.cluster_id
  environment_sub_type = "ACK"
}

`, name)
}

// Case 4551  twin
func TestAccAliCloudArmsEnvServiceMonitor_basic4551_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_service_monitor.default"
	ra := resourceAttrInit(resourceId, AlicloudArmsEnvServiceMonitorMap4551)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsEnvServiceMonitor")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvservicemonitor%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudArmsEnvServiceMonitorBasicDependence4551)
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
					"environment_id": "${alicloud_arms_environment.env-cs.id}",
					"config_yaml":    `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\n    o11y.aliyun.com/addon-name: mysql\n    o11y.aliyun.com/addon-version: 1.0.1\n    o11y.aliyun.com/release-name: mysql1\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 11s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus`,
					"aliyun_lang":    "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id": CHECKSET,
						"config_yaml":    CHECKSET,
						"aliyun_lang":    "en",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aliyun_lang"},
			},
		},
	})
}

// Test Arms EnvServiceMonitor. <<< Resource test cases, automatically generated.
