package alicloud

import (
	"testing"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCSKubernetes_basic(t *testing.T) {
	var k8s cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cs_kubernetes.k8s",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerKubernetes_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_kubernetes.k8s", &k8s),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_number", "3"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "name", "terraform-test-for-k8s-basic"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "nodes.#", "6"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_disk_category", "cloud_ssd"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_disk_size", "50"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_disk_size", "40"),
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
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerKubernetes_autoVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_kubernetes.k8s", &k8s),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_number", "3"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "nodes.#", "6"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_disk_category", "cloud_ssd"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_disk_size", "50"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "master_disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_kubernetes.k8s", "worker_disk_size", "40"),
				),
			},
		},
	})
}

const testAccContainerKubernetes_basic = `

provider "alicloud" {
	region="cn-shanghai"
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
  name = "tf_test_image"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_kubernetes" "k8s" {
  name = "terraform-test-for-k8s-basic"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  new_nat_gateway = true
  master_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  worker_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  worker_number = 3
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
	region="us-west-1"
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
  name_prefix = "terraform-test-for-k8s"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  new_nat_gateway = true
  master_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  worker_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  worker_number = 3
  password = "Test12345"
  pod_cidr = "192.168.1.0/24"
  service_cidr = "192.168.2.0/24"
  enable_ssh = true
  install_cloud_monitor = true
  worker_disk_category  = "cloud_ssd"
  master_disk_size = 50
}
`
