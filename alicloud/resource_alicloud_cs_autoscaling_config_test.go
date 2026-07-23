package alicloud

import (
	"fmt"
	"testing"

	cs "github.com/alibabacloud-go/cs-20151215/v7/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudCSAutoscalingConfig_basic(t *testing.T) {
	var v *cs.CreateAutoscalingConfigRequest
	resourceId := "alicloud_cs_autoscaling_config.default"
	serviceFunc := func() interface{} {
		client, _ := testAccProvider.Meta().(*connectivity.AliyunClient).NewRoaCsClient()
		return &CsClient{client}
	}

	ra := resourceAttrInit(resourceId, baseMap)
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rc.describeMethod = "DescribeCsAutoscalingConfig"
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSAutoscalingConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSAuoscalingConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
		},
		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
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
						"cluster_id":                CHECKSET,
						"cool_down_duration":        CHECKSET,
						"unneeded_duration":         CHECKSET,
						"utilization_threshold":     CHECKSET,
						"gpu_utilization_threshold": CHECKSET,
						"scan_interval":             CHECKSET,
						"scale_down_enabled":        CHECKSET,
						"expander":                  CHECKSET,
						"scaler_type":               CHECKSET,
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
						"skip_nodes_with_system_pods":   CHECKSET,
						"skip_nodes_with_local_storage": CHECKSET,
						"daemonset_eviction_for_nodes":  CHECKSET,
						"max_graceful_termination_sec":  CHECKSET,
						"min_replica_count":             CHECKSET,
						"recycle_node_deletion_enabled": CHECKSET,
						"scale_up_from_zero":            CHECKSET,
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
						"cluster_id":                CHECKSET,
						"cool_down_duration":        CHECKSET,
						"unneeded_duration":         CHECKSET,
						"utilization_threshold":     CHECKSET,
						"gpu_utilization_threshold": CHECKSET,
						"scan_interval":             CHECKSET,
						"scale_down_enabled":        CHECKSET,
						"expander":                  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expander": "priority",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expander": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expander": "least-waste",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expander": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scaler_type": "goatscaler",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaler_type": CHECKSET,
					}),
				),
			},
		},
	})
}

var baseMap = map[string]string{
	"cluster_id": CHECKSET,
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
