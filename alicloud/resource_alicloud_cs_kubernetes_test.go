package alicloud

import (
	"testing"

	"fmt"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudCSKubernetes_basic(t *testing.T) {
	var k8s cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cs_kubernetes.k8s",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerKubernetes_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_kubernetes.k8s", &k8s),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_number", "1"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "name", "testAccContainerKubernetes-basic"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_nodes.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_disk_category", "cloud_ssd"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_disk_size", "50"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_disk_size", "40"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "connections.%", "4"),
				),
			},
		},
	})
}

func TestAccAlicloudCSKubernetes_autoVpc(t *testing.T) {
	var k8s cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cs_kubernetes.k8s",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerKubernetes_autoVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_kubernetes.k8s", &k8s),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_number", "1"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_nodes.#", "3"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_disk_category", "cloud_ssd"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_disk_size", "50"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_disk_size", "40"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "connections.%", "4"),
				),
			},
		},
	})
}

func testAccCheckKubernetesClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient).csconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cs_kubernetes" {
			continue
		}

		cluster, err := client.DescribeCluster(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
				continue
			}
			return err
		}

		if cluster.ClusterID != "" {
			return fmt.Errorf("Error container cluster %s still exists.", rs.Primary.ID)
		}
	}

	return nil
}

const testAccContainerKubernetes_basic = `
variable "name" {
	default = "testAccContainerKubernetes-basic"
}
data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_kubernetes" "k8s" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  new_nat_gateway = true
  master_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  worker_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  worker_number = 1
  master_disk_category  = "cloud_ssd"
  worker_disk_size = 50
  password = "Test12345"
  pod_cidr = "192.168.1.0/24"
  service_cidr = "192.168.2.0/24"
  enable_ssh = true
  install_cloud_monitor = true
}
`

const testAccContainerKubernetes_autoVpc = `
provider "alicloud" {
	region="cn-hangzhou"
}
variable "name" {
	default = "testAccContainerKubernetes-autoVpc"
}
data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_cs_kubernetes" "k8s" {
  name_prefix = "${var.name}"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  new_nat_gateway = true
  master_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  worker_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  worker_number = 1
  password = "Test12345"
  pod_cidr = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
  enable_ssh = true
  install_cloud_monitor = true
  worker_disk_category  = "cloud_ssd"
  master_disk_size = 50
}
`
