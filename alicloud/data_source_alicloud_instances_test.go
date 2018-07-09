package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudInstancesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstancesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instances.inst"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.name", "testAccCheckAlicloudInstancesDataSourceBasic"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.status", "Running"),
				),
			},
		},
	})
}

func TestAccAlicloudInstancesDataSource_vpcId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstancesDataSourceVpcId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instances.inst"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.private_ip", "172.16.10.10"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.status", "Running"),
				),
			},
		},
	})
}

func TestAccAlicloudInstancesDataSource_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudImagesDataSourceTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instances.inst"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.tags.from", "datasource"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.tags.usage", "test"),
				),
			},
		},
	})
}

const testAccCheckAlicloudInstancesDataSourceBasic = `
data "alicloud_images" "images" {
	name_regex = "ubuntu*"
}
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccCheckAlicloudInstancesDataSourceBasic"
}
resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
}

resource "alicloud_instance" "foo" {
	image_id = "${data.alicloud_images.images.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	security_groups = ["${alicloud_security_group.tf_test_foo.*.id}"]
	instance_name = "${var.name}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_instances" "inst" {
	image_id = "${alicloud_instance.foo.image_id}"
	name_regex = "${alicloud_instance.foo.instance_name}"
	availability_zone = "${alicloud_instance.foo.availability_zone}"
}
`

const testAccCheckAlicloudInstancesDataSourceVpcId = `
data "alicloud_images" "images" {
	name_regex = "ubuntu*"
}
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccCheckAlicloudInstancesDataSourceVpcId"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/16"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	vswitch_id = "${alicloud_vswitch.foo.id}"
	private_ip = "172.16.10.10"
	image_id = "${data.alicloud_images.images.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
}

data "alicloud_instances" "inst" {
        vpc_id = "${alicloud_vpc.foo.id}"
        status = "Running"
        vswitch_id = "${alicloud_instance.foo.vswitch_id}"
}
`

const testAccCheckAlicloudImagesDataSourceTags = `
data "alicloud_images" "images" {
	name_regex = "ubuntu*"
}
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
variable "name" {
	default = "testAccCheckAlicloudImagesDataSourceTags"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/21"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	vswitch_id = "${alicloud_vswitch.foo.id}"
	image_id = "${data.alicloud_images.images.images.0.id}"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
	tags {
		from = "datasource"
		usage = "test"
	}
}

data "alicloud_instances" "inst" {
	tags {
		from = "datasource"
	}
	ids = ["${alicloud_instance.foo.id}"]
}
`
