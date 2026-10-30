package cs_test

import (
	"fmt"
	"testing"

	roacs "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudCSAutoscalingConfig_basic(t *testing.T) {
	var v *roacs.CreateAutoscalingConfigRequest
	resourceId := "alicloud_cs_autoscaling_config.default"
	serviceFunc := func() interface{} {
		client, _ := acctest.Provider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return cs.NewCsClient(client)
	}

	ra := acctest.ResourceAttrInit(resourceId, baseMap)
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rc.DescribeMethod = "DescribeCsAutoscalingConfig"
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSAutoscalingConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAuoscalingConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":                "${alicloud_cs_managed_kubernetes.default.id}",
					"cool_down_duration":        "10m",
					"unneeded_duration":         "10m",
					"utilization_threshold":     "0.5",
					"gpu_utilization_threshold": "0.5",
					"scan_interval":             "30s",
					"scale_down_enabled":        "true",
					"expander":                  "random",
					"scaler_type":               "cluster-autoscaler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                acctest.CHECKSET,
						"cool_down_duration":        acctest.CHECKSET,
						"unneeded_duration":         acctest.CHECKSET,
						"utilization_threshold":     acctest.CHECKSET,
						"gpu_utilization_threshold": acctest.CHECKSET,
						"scan_interval":             acctest.CHECKSET,
						"scale_down_enabled":        acctest.CHECKSET,
						"expander":                  acctest.CHECKSET,
						"scaler_type":               acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"skip_nodes_with_system_pods":   "true",
					"skip_nodes_with_local_storage": "false",
					"daemonset_eviction_for_nodes":  "false",
					"max_graceful_termination_sec":  "14400",
					"min_replica_count":             "0",
					"recycle_node_deletion_enabled": "false",
					"scale_up_from_zero":            "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"skip_nodes_with_system_pods":   acctest.CHECKSET,
						"skip_nodes_with_local_storage": acctest.CHECKSET,
						"daemonset_eviction_for_nodes":  acctest.CHECKSET,
						"max_graceful_termination_sec":  acctest.CHECKSET,
						"min_replica_count":             acctest.CHECKSET,
						"recycle_node_deletion_enabled": acctest.CHECKSET,
						"scale_up_from_zero":            acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cool_down_duration":        "5m",
					"unneeded_duration":         "5m",
					"utilization_threshold":     "0.6",
					"gpu_utilization_threshold": "0.6",
					"scan_interval":             "40s",
					"scale_down_enabled":        "false",
					"expander":                  "least-waste",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                acctest.CHECKSET,
						"cool_down_duration":        acctest.CHECKSET,
						"unneeded_duration":         acctest.CHECKSET,
						"utilization_threshold":     acctest.CHECKSET,
						"gpu_utilization_threshold": acctest.CHECKSET,
						"scan_interval":             acctest.CHECKSET,
						"scale_down_enabled":        acctest.CHECKSET,
						"expander":                  acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expander": "priority",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expander": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expander": "least-waste",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expander": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaler_type": "goatscaler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaler_type": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var baseMap = map[string]string{
	"cluster_id": acctest.CHECKSET,
}

func resourceCSAuoscalingConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id           = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitch.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  vswitch_ids   	   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 37)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

`, name)
}
