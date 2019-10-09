package alicloud

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudCSManagedKubernetesClustersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_managed_kubernetes_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testaccmanagedk8s-%d", rand),
		dataSourceCSManagedKubernetesClustersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_managed_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_managed_kubernetes.default.name}-fake",
		}),
	}

	var existCSManagedKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"ids.0":                                CHECKSET,
			"names.#":                              "1",
			"names.0":                              REGEXMATCH + fmt.Sprintf("tf-testaccmanagedk8s-%d", rand),
			"clusters.#":                           "1",
			"clusters.0.id":                        CHECKSET,
			"clusters.0.name":                      REGEXMATCH + fmt.Sprintf("tf-testaccmanagedk8s-%d", rand),
			"clusters.0.availability_zone":         CHECKSET,
			"clusters.0.security_group_id":         CHECKSET,
			"clusters.0.nat_gateway_id":            CHECKSET,
			"clusters.0.vpc_id":                    CHECKSET,
			"clusters.0.worker_nodes.#":            "2",
			"clusters.0.connections.%":             CHECKSET,
			"clusters.0.log_config.#":              "1",
			"clusters.0.log_config.0.type":         "",  // Because the API does not return  field 'type', the default value of empty is used
			"clusters.0.log_config.0.project":      "",  // Because the API does not return  field 'project', the default value of empty is used
			"clusters.0.worker_data_disk_category": "",  // Because the API does not return  field 'worker_data_disk_category', the default value of empty is used
			"clusters.0.worker_data_disk_size":     "0", // Because the API does not return  field 'worker_data_disk_size', the default value of 0 is used
		}
	}

	var fakeCSManagedKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var csManagedKubernetesClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCSManagedKubernetesClustersMapFunc,
		fakeMapFunc:  fakeCSManagedKubernetesClustersMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
	}
	csManagedKubernetesClustersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceCSManagedKubernetesClustersConfigDependence(name string) string {
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

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix = "${var.name}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_ids = ["${alicloud_vswitch.default.id}"]
  new_nat_gateway = true
  worker_instance_types = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
  worker_number = 2
  password = "Yourpassword1234"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  install_cloud_monitor = true
  slb_internet_enabled = true
  worker_disk_category  = "cloud_efficiency"
  worker_data_disk_category = "cloud_ssd"
  worker_data_disk_size =  200
  log_config {
    type = "SLS"
    project = "${var.name}-managed-sls"
  }
}
`, name)
}
