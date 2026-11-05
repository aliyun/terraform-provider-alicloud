package cs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCSKubernetesClustersDataSource(t *testing.T) {
	acctest.CheckoutAccount(t, true)
	defer acctest.CheckoutAccount(t, false)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cs_kubernetes_clusters.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacckubernetes-%d", rand),
		dataSourceCSKubernetesClustersConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_kubernetes.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_kubernetes.default.id}-fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_kubernetes.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_cs_kubernetes.default.name}-fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_kubernetes.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_cs_kubernetes.default.id}"},
			"name_regex":     "${alicloud_cs_kubernetes.default.name}-fake",
		}),
	}
	var existCSKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"ids.0":                           acctest.CHECKSET,
			"names.#":                         "1",
			"names.0":                         acctest.REGEXMATCH + fmt.Sprintf("tf-testacckubernetes-%d", rand),
			"clusters.#":                      "1",
			"clusters.0.id":                   acctest.CHECKSET,
			"clusters.0.name":                 acctest.REGEXMATCH + fmt.Sprintf("tf-testacckubernetes-%d", rand),
			"clusters.0.availability_zone":    acctest.CHECKSET,
			"clusters.0.security_group_id":    acctest.CHECKSET,
			"clusters.0.vpc_id":               acctest.CHECKSET,
			"clusters.0.master_nodes.#":       "1",
			"clusters.0.master_disk_size":     "50",
			"clusters.0.master_disk_category": "cloud_essd",
			"clusters.0.connections.#":        "1",
			"clusters.0.connections.0.master_public_ip":    acctest.CHECKSET,
			"clusters.0.connections.0.api_server_internet": acctest.CHECKSET,
			"clusters.0.connections.0.api_server_intranet": acctest.CHECKSET,
			"clusters.0.connections.0.service_domain":      acctest.CHECKSET,
		}
	}

	var fakeCSKubernetesClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var csKubernetesClustersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCSKubernetesClustersMapFunc,
		FakeMapFunc:  fakeCSKubernetesClustersMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
	}
	csKubernetesClustersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Master"
  system_disk_category = "cloud_essd"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_instance_types.default_m.instance_types.0.availability_zones.0
  vswitch_name = var.name
}

locals {
  vswitch_id = alicloud_vswitch.vswitch.id
}

resource "alicloud_cs_kubernetes" "default" {
  name                  = "${var.name}"
  master_vswitch_ids    = ["${local.vswitch_id}", "${local.vswitch_id}", "${local.vswitch_id}"]
  new_nat_gateway       = false
  master_instance_types = ["${data.alicloud_instance_types.default_m.instance_types.0.id}", "${data.alicloud_instance_types.default_m.instance_types.0.id}", "${data.alicloud_instance_types.default_m.instance_types.0.id}"]
  password              = "Yourpassword1234"
  pod_cidr              = cidrsubnet("10.0.0.0/8", 8, 33)
  service_cidr          = cidrsubnet("172.18.0.0/16", 4, 4)
  install_cloud_monitor = true
  master_disk_size      = 50
  master_disk_category  = "cloud_essd"
  proxy_mode            = "ipvs"
}
`, name)
}
