package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudInstancesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudInstancesDataSourceVpcId(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_instances.inst"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_instances.inst", "instances.0.availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_instances.inst", "instances.0.private_ip", "172.16.0.10"),
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
				Config: testAccCheckAlicloudImagesDataSourceTags(EcsInstanceCommonTestCase),
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

func testAccCheckAlicloudInstancesDataSourceVpcId(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckAlicloudInstancesDataSourceVpcId"
	}

	resource "alicloud_instance" "foo" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		private_ip = "172.16.0.10"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
	}

	data "alicloud_instances" "inst" {
		vpc_id = "${alicloud_vpc.default.id}"
		status = "Running"
		vswitch_id = "${alicloud_instance.foo.vswitch_id}"
	}
	`, common)
}

func testAccCheckAlicloudImagesDataSourceTags(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccCheckAlicloudImagesDataSourceTags"
	}

	resource "alicloud_instance" "foo" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		instance_name = "${var.name}"
		system_disk_category = "cloud_efficiency"
		security_groups = ["${alicloud_security_group.default.id}"]
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
	`, common)
}

const testAccCheckAlicloudImagesDataSourceEmpty = `
data "alicloud_instances" "inst" {
	name_regex = "^tf-testacc-fake-name"
}
`
