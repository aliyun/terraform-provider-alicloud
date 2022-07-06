package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCSKubernetesClustersDataSource(t *testing.T) {
	checkoutAccount(t, true)
	defer checkoutAccount(t, false)
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_kubernetes_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacckubernetes-%d", rand),
		dataSourceCSKubernetesClustersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_kubernetes.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_kubernetes.default.name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_kubernetes.default.name}-fake",
		}),
	}
	var existCSKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                      "1",
			"ids.0":                                      CHECKSET,
			"names.#":                                    "1",
			"names.0":                                    REGEXMATCH + fmt.Sprintf("tf-testacckubernetes-%d", rand),
			"clusters.#":                                 "1",
			"clusters.0.id":                              CHECKSET,
			"clusters.0.name":                            REGEXMATCH + fmt.Sprintf("tf-testacckubernetes-%d", rand),
			"clusters.0.availability_zone":               CHECKSET,
			"clusters.0.security_group_id":               CHECKSET,
			"clusters.0.nat_gateway_id":                  CHECKSET,
			"clusters.0.vpc_id":                          CHECKSET,
			"clusters.0.worker_numbers.#":                "1",
			"clusters.0.worker_numbers.0":                "1",
			"clusters.0.master_nodes.#":                  "3",
			"clusters.0.worker_disk_category":            "cloud_ssd",
			"clusters.0.master_disk_size":                "50",
			"clusters.0.master_disk_category":            "cloud_efficiency",
			"clusters.0.worker_disk_size":                "40",
			"clusters.0.connections.%":                   "4",
			"clusters.0.connections.master_public_ip":    CHECKSET,
			"clusters.0.connections.api_server_internet": CHECKSET,
			"clusters.0.connections.api_server_intranet": CHECKSET,
			"clusters.0.connections.service_domain":      CHECKSET,
		}
	}

	var fakeCSKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var csKubernetesClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCSKubernetesClustersMapFunc,
		fakeMapFunc:  fakeCSKubernetesClustersMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
	}
	csKubernetesClustersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceCSKubernetesClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default_m" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 4
	memory_size = 8
	kubernetes_node_role = "Master"
}

data "alicloud_instance_types" "default_w" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 4
	memory_size = 8
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

resource "alicloud_cs_kubernetes" "default" {
  name = "${var.name}"
  master_vswitch_ids = ["${local.vswitch_id}","${local.vswitch_id}","${local.vswitch_id}"]
  worker_vswitch_ids = ["${local.vswitch_id}"]
  new_nat_gateway = true
  master_instance_types = ["${data.alicloud_instance_types.default_m.instance_types.0.id}","${data.alicloud_instance_types.default_m.instance_types.0.id}","${data.alicloud_instance_types.default_m.instance_types.0.id}"]
  worker_instance_types = ["${data.alicloud_instance_types.default_w.instance_types.0.id}"]
  worker_number = 1
  password = "Yourpassword1234"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  enable_ssh = true
  install_cloud_monitor = true
  worker_disk_category  = "cloud_ssd"
  worker_data_disk_category = "cloud_ssd"
  worker_data_disk_size =  200
  master_disk_size = 50
}
`, name)
}
