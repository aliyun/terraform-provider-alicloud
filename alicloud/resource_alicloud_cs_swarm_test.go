package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func SkipTestAccAlicloudCSSwarm_vpc(t *testing.T) {
	var container cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
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

func SkipTestAccAlicloudCSSwarm_vpc_zero_node(t *testing.T) {
	var container cs.ClusterType
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions)
		},

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSwarm_basic_zero_node(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "0"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", fmt.Sprintf("tf-testAccCSSwarm-basic-zero-node-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "0"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_category", ""),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_size", "0"),
				),
			},

			{
				Config: testAccCSSwarm_basic_zero_node_update(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", fmt.Sprintf("tf-testAccCSSwarm-basic-zero-node-update-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_size", "20"),
				),
			},
		},
	})
}

func SkipTestAccAlicloudCSSwarm_vpc_noslb(t *testing.T) {
	var container cs.ClusterType

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSwarm_basic_noslb,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "disk_size", "20"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.0.eip", ""),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "slb_id", ""),
				),
			},
		},
	})
}

func SkipTestAccAlicloudCSSwarm_update(t *testing.T) {
	var container cs.ClusterType
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },

		IDRefreshName: "alicloud_cs_swarm.cs_vpc",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSwarm_update(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "2"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", fmt.Sprintf("tf-testAccCSSwarm-update-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "nodes.#", "2"),
				),
			},

			{
				Config: testAccCSSwarm_updateAfter(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerClusterExists("alicloud_cs_swarm.cs_vpc", &container),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "node_number", "3"),
					resource.TestCheckResourceAttr("alicloud_cs_swarm.cs_vpc", "name", fmt.Sprintf("tf-testAccCSSwarm-updateafter-%d", rand)),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeCluster(cluster.Primary.ID)
		})
		log.Printf("[DEBUG] check cluster %s attribute %#v", cluster.Primary.ID, raw)

		if err != nil {
			return err
		}
		attr, _ := raw.(cs.ClusterType)
		if attr.ClusterID == "" {
			return fmt.Errorf("Container cluster not found")
		}

		*d = attr
		return nil
	}
}

func testAccCheckSwarmClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cs_swarm" {
			continue
		}

		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeCluster(rs.Primary.ID)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
				continue
			}
			return err
		}
		cluster, _ := raw.(cs.ClusterType)
		if cluster.ClusterID != "" {
			return fmt.Errorf("Error container cluster %s still exists.", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCSSwarm_basic = `
variable "name" {
	default = "tf-testAccCSSwarm-basic"
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
  vpc_name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Yourpassword1234"
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

func testAccCSSwarm_basic_zero_node(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccCSSwarm-basic-zero-node-%d"
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
	  vpc_name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_vswitch" "foo" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.foo.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	}

	resource "alicloud_cs_swarm" "cs_vpc" {
	  password = "Yourpassword1234"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  name = "${var.name}"
	  node_number = 0
	  cidr_block = "172.20.0.0/24"
	  image_id = "${data.alicloud_images.main.images.0.id}"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	  release_eip = "true"
	}
	`, rand)
}

func testAccCSSwarm_basic_zero_node_update(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccCSSwarm-basic-zero-node-update-%d"
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
	  vpc_name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_vswitch" "foo" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.foo.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	}

	resource "alicloud_cs_swarm" "cs_vpc" {
	  password = "Yourpassword1234"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  name = "${var.name}"
	  node_number = 2
	  disk_category = "cloud_efficiency"
	  disk_size = 20
	  cidr_block = "172.20.0.0/24"
	  image_id = "${data.alicloud_images.main.images.0.id}"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	`, rand)
}

const testAccCSSwarm_basic_noslb = `
variable "name" {
	default = "tf-testAccCSSwarm-basic"
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
  vpc_name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "Yourpassword1234"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  name_prefix = "${var.name}"
  node_number = 2
  disk_category = "cloud_efficiency"
  disk_size = 20
  cidr_block = "172.20.0.0/24"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.foo.id}"
  release_eip = "true"
  need_slb = "false"
}
`

func testAccCSSwarm_update(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccCSSwarm-update-%d"
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
	  vpc_name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_vswitch" "foo" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.foo.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	}

	resource "alicloud_cs_swarm" "cs_vpc" {
	  password = "Yourpassword1234"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  name = "${var.name}"
	  node_number = 2
	  disk_category = "cloud_efficiency"
	  disk_size = 20
	  cidr_block = "172.20.0.0/24"
	  image_id = "${data.alicloud_images.main.images.0.id}"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	`, rand)
}

func testAccCSSwarm_updateAfter(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccCSSwarm-updateafter-%d"
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
	  vpc_name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_vswitch" "foo" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.foo.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
	}

	resource "alicloud_cs_swarm" "cs_vpc" {
	  password = "Yourpassword1234"
	  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	  name = "${var.name}"
	  node_number = 3
	  disk_category = "cloud_efficiency"
	  disk_size = 20
	  cidr_block = "172.20.0.0/24"
	  image_id = "${data.alicloud_images.main.images.0.id}"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	`, rand)
}
