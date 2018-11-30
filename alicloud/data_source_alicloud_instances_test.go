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
				Config: testAccCheckAlicloudInstancesDataSourceVpcId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instances.inst"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.private_ip", "172.16.10.10"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.status", string(Running)),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.name", "tf-testAccCheckAlicloudInstancesDataSourceVpcId"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.instance_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.image_id"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.public_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.eip", ""),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.security_groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.key_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.internet_max_bandwidth_out", "0"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.spot_strategy", string(NoSpot)),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.disk_device_mappings.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.tags.%", "0"),
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
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.availability_zone"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.private_ip"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.status", string(Running)),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.name", "tf-testAccCheckAlicloudImagesDataSourceTags"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.instance_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.image_id"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.public_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.eip", ""),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.security_groups.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.key_name", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.internet_max_bandwidth_out", "0"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.spot_strategy", string(NoSpot)),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.disk_device_mappings.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.tags.%", "7"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.tags.from", "datasource"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.tags.usage1", "test"),
				),
			},
		},
	})
}

func TestAccAlicloudInstancesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudImagesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instances.inst"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.region_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.availability_zone"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.private_ip"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.instance_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.vpc_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.vswitch_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.image_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.public_ip"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.eip"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.security_groups.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.key_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.creation_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.instance_charge_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.internet_max_bandwidth_out"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.spot_strategy"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.disk_device_mappings.#"),
					resource.TestCheckNoResourceAttr("data.alicloud_instances.inst", "instances.0.tags.%"),
				),
			},
		},
	})
}

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
	default = "tf-testAccCheckAlicloudInstancesDataSourceVpcId"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
	name = "${var.name}"
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
	default = "tf-testAccCheckAlicloudImagesDataSourceTags"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  	vpc_id = "${alicloud_vpc.foo.id}"
  	cidr_block = "172.16.0.0/21"
  	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	name = "${var.name}"
}

resource "alicloud_security_group" "tf_test_foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	name = "${var.name}"
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
		usage1 = "test"
		usage2 = "test"
		usage3 = "test"
		usage4 = "test"
		usage5 = "test"
		usage6 = "test"

	}
}

data "alicloud_instances" "inst" {
	tags {
		from = "datasource"
		usage1 = "test"
		usage2 = "test"
		usage3 = "test"
		usage4 = "test"
		usage5 = "test"
	}
	ids = ["${alicloud_instance.foo.id}"]
}
`

const testAccCheckAlicloudImagesDataSourceEmpty = `
data "alicloud_instances" "inst" {
	name_regex = "^tf-testacc-fake-name"
}
`
