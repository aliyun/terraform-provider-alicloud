package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudCSSwarm_vpc(t *testing.T) {
	var container cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSSwarm_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_size", "20"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.0.eip", ""),
				),
			},
		},
	})
}

func TestAccAlicloudCSSwarm_update(t *testing.T) {
	var container cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSSwarm_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", "testAccCSSwarm-update"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "2"),
				),
			},

			resource.TestStep{
				Config: testAccCSSwarm_updateAfter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "3"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", "testAccCSSwarm-updateafter"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "3"),
				),
			},
		},
	})
}

func testAccCheckContainerClusterExists(n string, d *cs.ClusterType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cluster, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found:%s", n)
		}

		if cluster.Primary.ID == "" {
			return fmt.Errorf("No Container cluster ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient).csconn
		attr, err := client.DescribeCluster(cluster.Primary.ID)
		log.Printf("[DEBUG] check cluster %s attribute %#v", cluster.Primary.ID, attr)

		if err != nil {
			return err
		}

		if attr.ClusterID == "" {
			return fmt.Errorf("Container cluster not found")
		}

		*d = attr
		return nil
	}
}

func testAccCheckSwarmClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient).csconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cs_swarm" {
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

const testAccCSSwarm_basic = `
variable "name" {
	default = "testAccCSSwarm-basic"
}
data "alicloud_images" main {
	most_recent = true
	name_regex = "^centos_6\\w{1,5}[64].*"
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

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  release_eip = "true"
}
`

const testAccCSSwarm_update = `
variable "name" {
	default = "testAccCSSwarm-update"
}
data "alicloud_images" main {
	most_recent = true
	name_regex = "^centos_6\\w{1,5}[64].*"
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

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
`

const testAccCSSwarm_updateAfter = `
variable "name" {
	default = "testAccCSSwarm-updateafter"
}
data "alicloud_images" main {
	most_recent = true
	name_regex = "^centos_6\\w{1,5}[64].*"
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

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name = "${var.name}"
  node_number = 3
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
`
