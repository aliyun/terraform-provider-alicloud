package alicloud

import (
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCSKubernetesClustersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudCSKubernetesClustersDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cs_kubernetes_clusters.k8s_clusters"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.#"),
					resource.TestMatchResourceAttr("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.name", regexp.MustCompile("^tf-testAccKubernetes-datasource*")),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.worker_numbers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.worker_numbers.0", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.master_nodes.#", "3"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.worker_disk_category", "cloud_ssd"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.master_disk_size", "50"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.master_disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.worker_disk_size", "40"),
					resource.TestCheckResourceAttr("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.connections.%", "4"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.connections.master_public_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.connections.api_server_internet"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.connections.api_server_intranet"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_kubernetes_clusters.k8s_clusters", "clusters.0.connections.service_domain"),
				),
			},
		},
	})
}

const testAccAlicloudCSKubernetesClustersDataSource = `
variable "name" {
	default = "tf-testAccKubernetes-datasource"
}
data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "master" {
	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
	kubernetes_node_role = "Master"
}

data "alicloud_instance_types" "worker" {
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

resource "alicloud_cs_kubernetes" "k8s" {
  name_prefix = "${var.name}"
  vswitch_ids = ["${alicloud_vswitch.foo.id}"]
  new_nat_gateway = true
  master_instance_types = ["${data.alicloud_instance_types.master.instance_types.0.id}"]
  worker_instance_types = ["${data.alicloud_instance_types.worker.instance_types.0.id}"]
  worker_numbers = [1]
  password = "Yourpassword1234"
  pod_cidr = "192.168.1.0/24"
  service_cidr = "192.168.2.0/24"
  enable_ssh = true
  install_cloud_monitor = true
  worker_disk_category  = "cloud_ssd"
  master_disk_size = 50
}

data "alicloud_cs_kubernetes_clusters" "k8s_clusters" {
  name_regex = "${alicloud_cs_kubernetes.k8s.name}"
  enable_details = true
}
`
