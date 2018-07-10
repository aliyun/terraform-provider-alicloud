package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDisk_basic(t *testing.T) {
	var v ecs.Disk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDiskConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists(
						"alicloud_disk.foo", &v),
					resource.TestCheckResourceAttr(
						"alicloud_disk.foo",
						"category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr(
						"alicloud_disk.foo",
						"size",
						"30"),
					resource.TestCheckResourceAttr(
						"alicloud_disk.foo",
						"encrypted",
						"false"),
				),
			},
		},
	})
}

func TestAccAlicloudDisk_withTags(t *testing.T) {
	var v ecs.Disk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		//module name
		IDRefreshName: "alicloud_disk.bar",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDiskConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists("alicloud_disk.bar", &v),
					resource.TestCheckResourceAttr(
						"alicloud_disk.bar",
						"tags.Name",
						"TerraformTest"),
				),
			},
		},
	})
}

func TestAccAlicloudDisk_encrypted(t *testing.T) {
	var v ecs.Disk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk.encrypted",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDiskConfigEncrypted,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDiskExists(
						"alicloud_disk.encrypted", &v),
					resource.TestCheckResourceAttr(
						"alicloud_disk.encrypted",
						"category",
						"cloud_efficiency"),
					resource.TestCheckResourceAttr(
						"alicloud_disk.encrypted",
						"size",
						"30"),
					resource.TestCheckResourceAttr(
						"alicloud_disk.encrypted",
						"encrypted",
						"true"),
				),
			},
		},
	})
}

func testAccCheckDiskExists(n string, disk *ecs.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Disk ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		d, err := client.DescribeDiskById("", rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("While checking disk existing, describing disk got an error: %#v.", err)
		}

		*disk = d
		return nil
	}
}

func testAccCheckDiskDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_disk" {
			continue
		}

		// Try to find the Disk
		client := testAccProvider.Meta().(*AliyunClient)

		d, err := client.DescribeDiskById("", rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("While checking disk destroy, describing disk got an error: %#v.", err)
		}

		if d.DiskId != "" {
			return fmt.Errorf("Error ECS Disk still exist")
		}
	}

	return nil
}

const testAccDiskConfig = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "testAccDiskConfig"
}
resource "alicloud_disk" "foo" {
	# cn-beijing
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
	description = "Hello ecs disk."
	category = "cloud_efficiency"
  	size = "30"
}
`
const testAccDiskConfigWithTags = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "testAccDiskConfigWithTags"
}
resource "alicloud_disk" "bar" {
	# cn-beijing
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	size = "20"
	tags {
	    Name = "TerraformTest"
	}
}
`
const testAccDiskConfigEncrypted = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
}
variable "name" {
	default = "testAccDiskConfigEncrypted"
}
resource "alicloud_disk" "encrypted" {
	# cn-beijing
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
	description = "Hello ecs disk."
	category = "cloud_efficiency"
  	size = "30"
	encrypted = true
}
`
