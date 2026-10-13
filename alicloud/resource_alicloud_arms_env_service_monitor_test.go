package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id": "${alicloud_arms_environment.env-cs.id}",
					"config_yaml":    `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 10s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus\n`,
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
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 10s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus\n`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 11s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus`,
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
					"config_yaml":    `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 10s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus\n`,
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
  pod_cidr             = "10.200.0.0/16"
  service_cidr         = "192.168.0.0/16"
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

locals {
  cluster_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}

resource "alicloud_arms_environment" "env-cs" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = local.cluster_id
  environment_sub_type = "ManagedKubernetes"
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id": "${alicloud_arms_environment.env-cs.id}",
					"config_yaml":    `apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: arms-admin1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  endpoints:\n  - interval: 30s\n    port: operator\n    path: /metrics\n  - interval: 11s\n    port: operator1\n    path: /metrics\n  namespaceSelector:\n    any: true\n  selector:\n    matchLabels:\n     app: arms-prometheus-ack-arms-prometheus`,
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
