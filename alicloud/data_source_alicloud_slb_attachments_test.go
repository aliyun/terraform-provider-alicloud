package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbAttachmentsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbAttachmentsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_attachments.filtered_attachments"),
					resource.TestCheckResourceAttr("data.alicloud_slb_attachments.filtered_attachments", "slb_attachments.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_slb_attachments.filtered_attachments", "slb_attachments.0.instance_id"),
					resource.TestCheckResourceAttr("data.alicloud_slb_attachments.filtered_attachments", "slb_attachments.0.weight", "42"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbAttachmentsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbAttachmentsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slb_attachments.filtered_attachments"),
					resource.TestCheckResourceAttr("data.alicloud_slb_attachments.filtered_attachments", "slb_attachments.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_attachments.filtered_attachments", "slb_attachments.0.instance_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_slb_attachments.filtered_attachments", "slb_attachments.0.weight"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSlbAttachmentsDataSourceBasic = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbAttachmentsDataSourceBasic"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}
data "alicloud_images" "images" {
  name_regex = "^ubuntu_16.*_64"
  most_recent = true
  owners = "system"
}
data "alicloud_instance_types" "instance_types" {
 	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	cpu_core_count = 2
	memory_size = 4
}

resource "alicloud_vpc" "sample_vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "sample_vswitch" {
  vpc_id = "${alicloud_vpc.sample_vpc.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.az.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

resource "alicloud_security_group" "sample_security_group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.sample_vpc.id}"
}

resource "alicloud_instance" "sample_instance" {
  image_id = "${data.alicloud_images.images.images.0.id}"

  instance_type = "${data.alicloud_instance_types.instance_types.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${alicloud_security_group.sample_security_group.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

resource "alicloud_slb_attachment" "sample_slb_attachment" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
  instance_ids = ["${alicloud_instance.sample_instance.id}"]
  weight = 42
}

data "alicloud_slb_attachments" "filtered_attachments" {
  load_balancer_id = "${alicloud_slb_attachment.sample_slb_attachment.load_balancer_id}"
  instance_ids = ["${alicloud_instance.sample_instance.id}"]
}
`

const testAccCheckAlicloudSlbAttachmentsDataSourceEmpty = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbAttachmentsDataSourceBasic"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "sample_vpc" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "sample_vswitch" {
  vpc_id = "${alicloud_vpc.sample_vpc.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.az.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

data "alicloud_slb_attachments" "filtered_attachments" {
  load_balancer_id = "${alicloud_slb.sample_slb.id}"
}
`
