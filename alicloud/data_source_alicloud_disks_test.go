package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDisksDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", "tf-testAccCheckAlicloudDisksDataSourceBasic"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.description", "tf-testAccCheckAlicloudDisksDataSourceBasic_description"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.type", "data"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.category", "cloud_efficiency"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted", "off"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.size", "20"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.image_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.instance_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.attached_time", ""),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.tags.Name", "TerraformTest"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_filterByAllFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSourceFilterByAllFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", "tf-testAccCheckAlicloudDisksDataSourceFilterByAllFields"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_filterByInstanceId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSourceFilterByInstanceId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.name", "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.0.status", "In_use"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_disks.disks", "disks.0.attached_time"),
				),
			},
		},
	})
}

func TestAccAlicloudDisksDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDisksDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_disks.disks"),
					resource.TestCheckResourceAttr("data.alicloud_disks.disks", "disks.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.region_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.availability_zone"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.type"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.category"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.encrypted"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.size"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.image_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.snapshot_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.instance_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.creation_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.attached_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.detached_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.expiration_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_disks.disks", "disks.0.tags.%"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDisksDataSourceBasic = `
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSourceBasic"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
    name_regex = "${alicloud_disk.sample_disk.name}"
}
`

const testAccCheckAlicloudDisksDataSourceFilterByAllFields = `
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSourceFilterByAllFields"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
	tags {
	    Name = "TerraformTest"
	    Name1 = "TerraformTest"
	}
}

data "alicloud_disks" "disks" {
    ids = ["${alicloud_disk.sample_disk.id}"]
    name_regex = "${alicloud_disk.sample_disk.name}"
    type = "data"
    category = "cloud_efficiency"
    encrypted = "off"
    tags = {
        Name = "TerraformTest"
    }
}
`

const testAccCheckAlicloudDisksDataSourceFilterByInstanceId = `
variable "name" {
	default = "tf-testAccCheckAlicloudDisksDataSourceFilterByInstanceId"
}

data "alicloud_zones" "az" {
	"available_resource_creation"= "VSwitch"
}
data "alicloud_images" "images" {
	name_regex = "ubuntu*"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

resource "alicloud_disk" "sample_disk" {
	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
    description = "${var.name}_description"
	size = "20"
}

resource "alicloud_vpc" "sample_vpc" {
	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "sample_vswitch" {
	name = "${var.name}"
  	vpc_id = "${alicloud_vpc.sample_vpc.id}"
  	cidr_block = "172.16.0.0/16"
  	availability_zone = "${data.alicloud_zones.az.zones.0.id}"
}

resource "alicloud_security_group" "sample_security_group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.sample_vpc.id}"
}

resource "alicloud_instance" "sample_instance" {
	vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
	private_ip = "172.16.10.10"
	image_id = "${data.alicloud_images.images.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
	security_groups = ["${alicloud_security_group.sample_security_group.id}"]
}

resource "alicloud_disk_attachment" "sample_disk_attachment" {
  disk_id = "${alicloud_disk.sample_disk.id}"
  instance_id = "${alicloud_instance.sample_instance.id}"
}

data "alicloud_disks" "disks" {
    instance_id = "${alicloud_disk_attachment.sample_disk_attachment.instance_id}"
    type = "data"
}
`

const testAccCheckAlicloudDisksDataSourceEmpty = `
data "alicloud_disks" "disks" {
    name_regex = "^tf-testacc-fake-name"
}
`
