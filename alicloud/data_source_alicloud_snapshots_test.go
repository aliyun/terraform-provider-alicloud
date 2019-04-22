package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSnapshotsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceSnapshotsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_snapshots.snapshots"),
					resource.TestCheckResourceAttr("data.alicloud_snapshots.snapshots", "snapshots.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_snapshots.snapshots", "snapshots.0.name", "tf-testAcc-snapshot"),
					resource.TestCheckResourceAttr("data.alicloud_snapshots.snapshots", "snapshots.0.description", "TF Test"),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.status"),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.progress"),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.source_disk_id"),
					resource.TestCheckResourceAttr("data.alicloud_snapshots.snapshots", "snapshots.0.source_disk_size", "20"),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.source_disk_type"),
					resource.TestCheckResourceAttr("data.alicloud_snapshots.snapshots", "snapshots.0.product_code", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.retention_days"),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.remain_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.status"),
					resource.TestCheckResourceAttrSet("data.alicloud_snapshots.snapshots", "snapshots.0.usage"),
				),
			},
		},
	})
}

const testAccDataSourceSnapshotsConfig = `
data "alicloud_instance_types" "instance_type" {
}

resource "alicloud_vpc" "vpc" {
  name = "tf-testAcc-vpc"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_zones" "zone" {
}

resource "alicloud_vswitch" "vswitch" {
  name = "tf-testAcc-vswitch"
  cidr_block = "192.168.0.0/24"
  availability_zone = "${data.alicloud_zones.zone.zones.0.id}"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "group" {
  name        = "tf-testACC-group"
  description = "New security group"
  vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_disk" "disk" {
  availability_zone = "${alicloud_instance.instance.availability_zone}"
  category          = "cloud_efficiency"
  size              = "20"
}

data "alicloud_images" "sys_images" {
  owners = "system"
}

resource "alicloud_instance" "instance" {
  instance_name   = "tf-testAcc-instance"
  host_name       = "tf-testAcc"
  image_id        = "${data.alicloud_images.sys_images.images.0.id}"
  instance_type   = "${data.alicloud_instance_types.instance_type.instance_types.0.id}"
  security_groups = ["${alicloud_security_group.group.id}"]
  vswitch_id      = "${alicloud_vswitch.vswitch.id}"
}

resource "alicloud_disk_attachment" "instance-attachment" {
  disk_id     = "${alicloud_disk.disk.id}"
  instance_id = "${alicloud_instance.instance.id}"
}

resource "alicloud_snapshot" "snapshot" {
  disk_id = "${alicloud_disk_attachment.instance-attachment.disk_id}"
  name = "tf-testAcc-snapshot"
  description = "TF Test"
  tags = {
    version = "1.0"
  }
}

data "alicloud_snapshots" "snapshots" {
  ids = ["${alicloud_snapshot.snapshot.id}"]
  name_regex = "tf-testAcc-snapshot"
}
`
