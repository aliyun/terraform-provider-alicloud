package alicloud

import (
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCSManagedKubernetesClustersDataSource_Empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudCSManagedKubernetesClustersDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.#"),
				),
			},
		},
	})
}

const testAccAlicloudCSManagedKubernetesClustersDataSourceEmpty = `
data "alicloud_cs_managed_kubernetes_clusters" "k8s_clusters" {
}
`

func TestAccAlicloudCSManagedKubernetesClustersDataSource_AutoVpc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudCSManagedKubernetesClustersDataSourceAutoVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.#"),
					resource.TestMatchResourceAttr("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.name", regexp.MustCompile("^tf-testAccManagedKubernetes-autoVpc*")),
					resource.TestCheckResourceAttr("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.worker_numbers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.worker_numbers.0", "2"),
					resource.TestCheckResourceAttr("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.worker_disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.worker_disk_size", "40"),
					resource.TestCheckResourceAttr("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.connections.%", "4"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.connections.master_public_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.connections.api_server_internet"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.connections.api_server_intranet"),
					resource.TestCheckResourceAttrSet("data.alicloud_cs_managed_kubernetes_clusters.k8s_clusters", "clusters.0.connections.service_domain"),
				),
			},
		},
	})
}

const testAccAlicloudCSManagedKubernetesClustersDataSourceAutoVpc = `
variable "name" {
	default = "tf-testAccManagedKubernetes-autoVpc"
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

resource "alicloud_cs_managed_kubernetes" "k8s" {
  name_prefix = "${var.name}"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  new_nat_gateway = true
  worker_instance_types = ["${data.alicloud_instance_types.default.instance_types.0.id}"]
  worker_numbers = [2]
  password = "Yourpassword1234"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  install_cloud_monitor = true
  slb_internet_enabled = true
  worker_disk_category  = "cloud_efficiency"
}

data "alicloud_cs_managed_kubernetes_clusters" "k8s_clusters" {
  name_regex = "${alicloud_cs_managed_kubernetes.k8s.name}"
  enable_details = true
}
`
