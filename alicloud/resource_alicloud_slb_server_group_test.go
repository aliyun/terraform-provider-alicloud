package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

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
			{
				Config: testAccSlbServerGroupVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbServerGroupExists("alicloud_slb_server_group.group", &group),
					resource.TestCheckResourceAttr(
						"alicloud_slb_server_group.group", "name", "tf-testAccSlbServerGroupVpc"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_server_group.group", "servers.#", "2"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbServerGroup_empty(t *testing.T) {
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
			{
				Config: testAccSlbServerGroupEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbServerGroupExists("alicloud_slb_server_group.group", &group),
					resource.TestCheckResourceAttr(
						"alicloud_slb_server_group.group", "name", "tf-testAccSlbServerGroupEmpty"),
					resource.TestCheckResourceAttr(
						"alicloud_slb_server_group.group", "servers.#", "0"),
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
			return fmt.Errorf("No SLB Server Group ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		slbService := SlbService{client}
		gr, err := slbService.DescribeSlbVServerGroupAttribute(rs.Primary.ID)
		if err != nil {
			return err
		}

		*group = *gr

		return nil
	}
}

func testAccCheckSlbServerGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	slbService := SlbService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_server_group" {
			continue
		}

		// Try to find the Slb server group
		if _, err := slbService.DescribeSlbVServerGroupAttribute(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("SLB Server Group %s still exist.", rs.Primary.ID)
	}

	return nil
}

const testAccSlbServerGroupVpc = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "image" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "tf-testAccSlbServerGroupVpc"
}

resource "alicloud_vpc" "main" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "alicloud_security_group" "group" {
  name = "${var.name}"
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_instance" "instance" {
  image_id = "${data.alicloud_images.image.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = ["${alicloud_security_group.group.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb" "instance" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb_server_group" "group" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  name = "${var.name}"
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

const testAccSlbServerGroupEmpty = `
data "alicloud_zones" "default" {
  "available_disk_category"     = "cloud_efficiency"
  "available_resource_creation" = "VSwitch"
}

variable "name" {
  default = "tf-testAccSlbServerGroupEmpty"
}

resource "alicloud_vpc" "main" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id            = "${alicloud_vpc.main.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_slb" "instance" {
  name     = "${var.name}"
  internet = true
}

resource "alicloud_slb_server_group" "group" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  name             = "${var.name}"
}
`
