package alicloud

import (
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCSManagedKubernetesClustersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKubernetesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cs_managed_kubernetes_clusters.default"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.default", "clusters.#"),
					resource.TestMatchResourceAttr("data.alicloud_cs_managed_kubernetes_clusters.default", "clusters.0.name", regexp.MustCompile("^tf-testAccManagedK8s-datasource*")),
					resource.TestCheckResourceAttr("data.alicloud_cs_managed_kubernetes_clusters.default", "clusters.0.worker_nodes.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_cs_managed_kubernetes_clusters.default", "clusters.0.connections.%", "4"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.default", "clusters.0.connections.master_public_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.default", "clusters.0.connections.api_server_internet"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.default", "clusters.0.connections.api_server_intranet"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.default", "clusters.0.connections.service_domain"),
				),
			},
		},
	})
}

const testAccManagedKubernetesDataSource = `
variable "name" {
	default = "tf-testAccManagedK8s-datasource"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_managed_kubernetes" "k8s" {
  name_prefix = "${var.name}"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  vswitch_ids = ["${alicloud_vswitch.foo.id}"]
  new_nat_gateway = true
  worker_instance_types = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
  worker_number = 2
  password = "Yourpassword1234"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  install_cloud_monitor = true
  slb_internet_enabled = true
  worker_disk_category  = "cloud_efficiency"
}
data "alicloud_cs_managed_kubernetes_clusters" "default" {
  name_regex = "${alicloud_cs_managed_kubernetes.k8s.name}"
  enable_details = true
}
`
