package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/arms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Arms EnvPodMonitor. >>> Resource test cases, automatically generated.
// Case 4553
func TestAccAliCloudArmsEnvPodMonitor_basic4553(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_pod_monitor.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsEnvPodMonitorMap4553)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsEnvPodMonitor")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvpodmonitor%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsEnvPodMonitorBasicDependence4553)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id": "${alicloud_arms_environment.environment-cs.id}",
					"config_yaml":    `apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: tf-test-pm1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  selector:\n    matchLabels:\n      app: arms-prometheus-ack-arms-prometheus\n      release: arms-prometheus\n  namespaceSelector:\n    any: true    \n  podMetricsEndpoints:\n  - interval: 30s\n    targetPort: 9335\n    path: /metrics\n  - interval: 10s\n    targetPort: 9335\n    path: /metrics1`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id": acctest.CHECKSET,
						"config_yaml":    acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: tf-test-pm1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  selector:\n    matchLabels:\n      app: arms-prometheus-ack-arms-prometheus\n      release: arms-prometheus\n  namespaceSelector:\n    any: true    \n  podMetricsEndpoints:\n  - interval: 30s\n    targetPort: 9335\n    path: /metrics\n  - interval: 10s\n    targetPort: 9335\n    path: /metrics1`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_yaml": `apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: tf-test-pm1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  selector:\n    matchLabels:\n      app: arms-prometheus-ack-arms-prometheus\n      release: arms-prometheus\n  namespaceSelector:\n    any: true    \n  podMetricsEndpoints:\n  - interval: 30s\n    targetPort: 9335\n    path: /metrics\n  - interval: 11s\n    targetPort: 9335\n    path: /metrics1`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_yaml": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id": "${alicloud_arms_environment.environment-cs.id}",
					"config_yaml":    `apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: tf-test-pm1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  selector:\n    matchLabels:\n      app: arms-prometheus-ack-arms-prometheus\n      release: arms-prometheus\n  namespaceSelector:\n    any: true    \n  podMetricsEndpoints:\n  - interval: 30s\n    targetPort: 9335\n    path: /metrics\n  - interval: 10s\n    targetPort: 9335\n    path: /metrics1`,
					"aliyun_lang":    "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id": acctest.CHECKSET,
						"config_yaml":    acctest.CHECKSET,
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

var AlicloudArmsEnvPodMonitorMap4553 = map[string]string{
	"status":               acctest.CHECKSET,
	"env_pod_monitor_name": acctest.CHECKSET,
	"namespace":            acctest.CHECKSET,
}

func AlicloudArmsEnvPodMonitorBasicDependence4553(name string) string {
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
  pod_cidr             = "10.126.0.0/16"
  service_cidr         = "192.168.0.0/16"
  slb_internet_enabled = true
  is_enterprise_security_group = true
}

locals {
  cluster_id = length(data.alicloud_cs_managed_kubernetes_clusters.default.ids) > 0 ? data.alicloud_cs_managed_kubernetes_clusters.default.ids.0 : alicloud_cs_managed_kubernetes.default.0.id
}

resource "alicloud_arms_environment" "environment-cs" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = local.cluster_id
  environment_sub_type = "ManagedKubernetes"
}


`, name)
}

// Case 4553  twin
func TestAccAliCloudArmsEnvPodMonitor_basic4553_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_env_pod_monitor.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudArmsEnvPodMonitorMap4553)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return arms.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeArmsEnvPodMonitor")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsenvpodmonitor%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudArmsEnvPodMonitorBasicDependence4553)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"environment_id": "${alicloud_arms_environment.environment-cs.id}",
					"config_yaml":    `apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: tf-test-pm1\n  namespace: arms-prom\n  annotations:\n    arms.prometheus.io/discovery: 'true'\nspec:\n  selector:\n    matchLabels:\n      app: arms-prometheus-ack-arms-prometheus\n      release: arms-prometheus\n  namespaceSelector:\n    any: true    \n  podMetricsEndpoints:\n  - interval: 30s\n    targetPort: 9335\n    path: /metrics\n  - interval: 11s\n    targetPort: 9335\n    path: /metrics1`,
					"aliyun_lang":    "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"environment_id": acctest.CHECKSET,
						"config_yaml":    acctest.CHECKSET,
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

// Test Arms EnvPodMonitor. <<< Resource test cases, automatically generated.
