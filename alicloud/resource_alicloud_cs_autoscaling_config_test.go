package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	cs "github.com/alibabacloud-go/cs-20151215/v3/client"
)

func TestAccAlicloudCSAutoscalingConfig_basic(t *testing.T) {
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
					"cluster_id":                "${alicloud_cs_managed_kubernetes.default.0.id}",
					"cool_down_duration":        "10m",
					"unneeded_duration":         "10m",
					"utilization_threshold":     "0.5",
					"gpu_utilization_threshold": "0.5",
					"scan_interval":             "30s",
					"scale_down_enabled":        "true",
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
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
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
  name_prefix                 = "${var.name}"
  count                       = 1
  cluster_spec                = "ack.pro.small"
  worker_vswitch_ids          = [local.vswitch_id]
  new_nat_gateway             = true
  worker_instance_types       = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
  worker_number               = 2
  node_port_range             = "30000-32767"
  password                    = "Hello1234"
  pod_cidr                    = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr                = cidrsubnet("172.16.0.0/16", 4, 7)
  install_cloud_monitor       = true
  slb_internet_enabled        = true
  worker_disk_category        = "cloud_efficiency"
  worker_data_disk_category   = "cloud_ssd"
  worker_data_disk_size       = 200
  worker_disk_size            = 40
  worker_instance_charge_type = "PostPaid"
}
`, name)
}
