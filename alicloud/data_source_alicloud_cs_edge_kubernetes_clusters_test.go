package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCSEdgeKubernetesClustersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_edge_kubernetes_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testaccedgek8s-%d", rand),
		dataSourceCSEdgeKubernetesClustersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_edge_kubernetes.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_edge_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_edge_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_edge_kubernetes.default.name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_edge_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_edge_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_edge_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_edge_kubernetes.default.name}-fake",
		}),
	}

	var existCSEdgeKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        CHECKSET,
			"names.#":                      "1",
			"names.0":                      REGEXMATCH + fmt.Sprintf("tf-testaccedgek8s-%d", rand),
			"clusters.#":                   "1",
			"clusters.0.id":                CHECKSET,
			"clusters.0.name":              REGEXMATCH + fmt.Sprintf("tf-testaccedgek8s-%d", rand),
			"clusters.0.availability_zone": CHECKSET,
			"clusters.0.security_group_id": CHECKSET,
			"clusters.0.nat_gateway_id":    CHECKSET,
			"clusters.0.vpc_id":            CHECKSET,
			"clusters.0.worker_nodes.#":    "2",
			"clusters.0.connections.%":     CHECKSET,
		}
	}

	var fakeCSEdgeKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var csEdgeKubernetesClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCSEdgeKubernetesClustersMapFunc,
		fakeMapFunc:  fakeCSEdgeKubernetesClustersMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
	}
	csEdgeKubernetesClustersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceCSEdgeKubernetesClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
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

variable "disks" {
  type = list(object({
    size = string
    category = string
  }))
  default = [
    {
      "size" = "200",
      "category" = "cloud_efficiency",
    }
  ]

}


resource "alicloud_cs_edge_kubernetes" "default" {
  name_prefix = "${var.name}"
  worker_vswitch_ids = ["${local.vswitch_id}"]
  new_nat_gateway = true
  worker_instance_types = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
  worker_number = 2
  password = "Yourpassword1234"
  pod_cidr = "172.31.0.0/16"
  service_cidr = "172.21.0.0/20"
  install_cloud_monitor = true
  slb_internet_enabled = true
  worker_disk_category  = "cloud_efficiency"
  dynamic "worker_data_disks" {
      for_each = var.disks
      content {
        size       = lookup(worker_data_disks.value, "size", var.disks)
        category       = lookup(worker_data_disks.value, "category", var.disks)
      }
  }
}
`, name)
}
