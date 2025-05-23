package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	cs "github.com/alibabacloud-go/cs-20151215/v5/client"
)

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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
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
	availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count       = 4
	memory_size          = 8
	kubernetes_node_role = "Worker"
	instance_type_family = "ecs.sn1ne"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING-ACK$"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [local.vswitch_id]
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 37)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

`, name)
}
