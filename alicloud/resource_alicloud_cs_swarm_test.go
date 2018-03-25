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
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSSwarm_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "size", "2"),
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
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSSwarm_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "size", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", "ClusterOfTestFromTerraform"),
				),
			},

			resource.TestStep{
				Config: testAccCSSwarm_updateAfter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "size", "3"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", "ClusterOfTestFromTerraformUpdate"),
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
			return fmt.Errorf("DB Instance not found")
		}

		*d = attr
		return nil
	}
}

func testAccCheckContainerClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient).csconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cs_swarm" {
			continue
		}

		cluster, err := client.DescribeCluster(rs.Primary.ID)

		if err != nil {
			if IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return err
		}

		if cluster.ClusterID == "" {
			return nil
		}
	}

	return nil
}

const testAccCSSwarm_basic = `
data "alicloud_images" main {
  most_recent = true
  name_regex = "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "tf_test_swarm"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "ecs.n4.small"
  name_prefix = "ClusterOfVpcTest"
  size = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
`

const testAccCSSwarm_update = `
data "alicloud_images" main {
  most_recent = true
  name_regex = "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "tf_test_swarm"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "ecs.n4.small"
  name = "ClusterOfTestFromTerraform"
  size = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
`

const testAccCSSwarm_updateAfter = `
data "alicloud_images" main {
  most_recent = true
  name_regex = "^centos_6\\w{1,5}[64].*"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name = "tf_test_swarm"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Just$test"
  instance_type = "ecs.n4.small"
  name = "ClusterOfTestFromTerraformUpdate"
  size = 3
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
}
`
