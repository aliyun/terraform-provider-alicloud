package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudSlbServerGroup_classic(t *testing.T) {
	var group slb.DescribeVServerGroupAttributeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_server_group.group",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbServerGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbServerGroupClassic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbServerGroupExists("alicloud_slb_server_group.group", &group),
					resource.TestCheckResourceAttr(
						"alicloud_slb_server_group.group", "name", "tf-server-group"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_server_group.group", "servers.#", "3"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbServerGroup_vpc(t *testing.T) {
	var group slb.DescribeVServerGroupAttributeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_server_group.group",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbServerGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSlbServerGroupVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbServerGroupExists("alicloud_slb_server_group.group", &group),
					resource.TestCheckResourceAttr(
						"alicloud_slb_server_group.group", "name", "tf-server-group"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_server_group.group", "servers.#", "2"),
				),
			},
		},
	})
}

func testAccCheckSlbServerGroupExists(n string, group *slb.DescribeVServerGroupAttributeResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SLB listener ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		gr, err := client.slbconn.DescribeVServerGroupAttribute(&slb.DescribeVServerGroupAttributeArgs{
			RegionId:       client.Region,
			VServerGroupId: rs.Primary.ID,
		})
		if err != nil {
			return fmt.Errorf("DescribeVServerGroupAttribute got an error: %#v", err)
		}
		if gr == nil {
			return fmt.Errorf("Specified VServer Group not found")
		}

		*group = *gr

		return nil
	}
}

func testAccCheckSlbServerGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_server_group" {
			continue
		}

		// Try to find the Slb server group
		group, err := client.slbconn.DescribeVServerGroupAttribute(&slb.DescribeVServerGroupAttributeArgs{
			RegionId:       client.Region,
			VServerGroupId: rs.Primary.ID,
		})
		if err != nil {
			if IsExceptedError(err, VServerGroupNotFoundMessage) || IsExceptedError(err, InvalidParameter) {
				return nil
			}
			return fmt.Errorf("DescribeVServerGroupAttribute got an error: %#v", err)
		}
		if group != nil {

		}
		return fmt.Errorf("SLB listener still exist")
	}

	return nil
}

const testAccSlbServerGroupClassic = `
data "alicloud_images" "image" {
	most_recent = true
	owners = "system"
	name_regex = "^centos_6\\w{1,5}[64]{1}.*"
}

data "alicloud_zones" "zone" {}

resource "alicloud_vpc" "main" {
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
  depends_on = [
    "alicloud_vpc.main"]
}
resource "alicloud_security_group" "group" {
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_instance" "vpc" {
  image_id = "${data.alicloud_images.image.images.0.id}"
  instance_type = "ecs.n4.small"
  count = "2"
  security_groups = ["${alicloud_security_group.group.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_security_group" "classic" {}

resource "alicloud_instance" "classic" {
  image_id = "${data.alicloud_images.image.images.0.id}"
  instance_type = "ecs.n4.small"
  security_groups = ["${alicloud_security_group.classic.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
}

resource "alicloud_slb" "instance" {
  internet = true
}

resource "alicloud_slb_server_group" "group" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  servers = [
    {
      server_ids = ["${alicloud_instance.vpc.*.id}"]
      port = 100
      weight = 10
    },
    {
      server_ids = ["${alicloud_instance.classic.*.id}", "${alicloud_instance.vpc.*.id}"]
      port = 80
      weight = 100
    },
    {
      server_ids = ["${alicloud_instance.classic.*.id}"]
      port = 22
      weight = 100
    }
  ]
}
`

const testAccSlbServerGroupVpc = `
data "alicloud_images" "image" {
	most_recent = true
	owners = "system"
	name_regex = "^centos_6\\w{1,5}[64]{1}.*"
}

data "alicloud_zones" "zone" {}

resource "alicloud_vpc" "main" {
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
  depends_on = [
    "alicloud_vpc.main"]
}
resource "alicloud_security_group" "group" {
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_instance" "instance" {
  image_id = "${data.alicloud_images.image.images.0.id}"
  instance_type = "ecs.n4.small"
  count = "2"
  security_groups = ["${alicloud_security_group.group.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb" "instance" {
  vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb_server_group" "group" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  servers = [
    {
      server_ids = ["${alicloud_instance.instance.0.id}", "${alicloud_instance.instance.1.id}"]
      port = 100
      weight = 10
    },
    {
      server_ids = ["${alicloud_instance.instance.*.id}"]
      port = 80
      weight = 100
    }
  ]
}
`
