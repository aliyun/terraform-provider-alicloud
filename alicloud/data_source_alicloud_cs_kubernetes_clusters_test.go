package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCSKubernetesClustersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_kubernetes_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacckubernetes-%d", rand),
		dataSourceCSKubernetesClustersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{alicloud_cs_kubernetes.default.id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     alicloud_cs_kubernetes.default.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_kubernetes.default.name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{alicloud_cs_kubernetes.default.id},
			"name_regex":     alicloud_cs_kubernetes.default.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{alicloud_cs_kubernetes.default.id},
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
			"clusters.0.log_config.#":                    "1",
			"clusters.0.log_config.0.type":               "logtail-ds",
			"clusters.0.log_config.0.project":            fmt.Sprintf("tf-testacckubernetes-%d-delicate-sls", rand),
			"clusters.0.worker_data_disk_category":       "",  // Because the API does not return  field 'worker_data_disk_category', the default value of empty is used
			"clusters.0.worker_data_disk_size":           "0", // Because the API does not return  field 'worker_data_disk_size', the default value of 0 is used
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
data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default_m" {
	availability_zone = data.alicloud_zones.default.zones.0.id
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Master"
}

data "alicloud_instance_types" "default_w" {
	availability_zone = data.alicloud_zones.default.zones.0.id
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "default" {
  name = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  name = var.name
  vpc_id = alicloud_vpc.default.id
  cidr_block = "10.1.1.0/24"
  availability_zone = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cs_kubernetes" "default" {
  name_prefix = var.name
  vswitch_ids = [alicloud_vswitch.default.id]
  new_nat_gateway = true
  master_instance_types = [data.alicloud_instance_types.default_m.instance_types.0.id]
  worker_instance_types = [data.alicloud_instance_types.default_w.instance_types.0.id]
  worker_numbers = [1]
  password = "Yourpassword1234"
  pod_cidr = "192.168.1.0/24"
  service_cidr = "192.168.2.0/24"
  enable_ssh = true
  install_cloud_monitor = true
  worker_disk_category  = "cloud_ssd"
  worker_data_disk_category = "cloud_ssd"
  worker_data_disk_size =  200
  master_disk_size = 50	 
  log_config {
    type = "SLS"
    project = "${var.name}-delicate-sls"
  }
}
`, name)
}
