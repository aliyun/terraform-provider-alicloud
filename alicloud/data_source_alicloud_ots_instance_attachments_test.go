package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudOtsInstanceAttachmentsDataSource_basic(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsInstanceAttachmentsDataSource_basic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instance_attachments.attachments"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "attachments.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "attachments.0.id", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instance_attachments.attachments", "attachments.0.domain"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instance_attachments.attachments", "attachments.0.endpoint"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instance_attachments.attachments", "attachments.0.region"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "attachments.0.instance_name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "attachments.0.vpc_name", "testvpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instance_attachments.attachments", "attachments.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "names.0", "testvpc"),
				),
			},
		},
	})
}

func TestAccAlicloudOtsInstanceAttachmentsDataSource_name_regex(t *testing.T) {
	randInt := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsInstanceAttachmentsDataSource_name_regex_exist(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instance_attachments.attachments"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "attachments.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "attachments.0.id", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instance_attachments.attachments", "attachments.0.domain"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instance_attachments.attachments", "attachments.0.endpoint"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instance_attachments.attachments", "attachments.0.region"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "attachments.0.instance_name", fmt.Sprintf("tf-testAcc%d", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "attachments.0.vpc_name", "testvpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_instance_attachments.attachments", "attachments.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "names.0", "testvpc"),
				),
			},
			{
				Config: testAccCheckAlicloudOtsInstanceAttachmentsDataSource_name_regex_fake(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_instance_attachments.attachments"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "attachments.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ots_instance_attachments.attachments", "names.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudOtsInstanceAttachmentsDataSource_basic(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}

	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "Capacity"
	}

	data "alicloud_zones" "foo" {
	  available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "foo" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}

	resource "alicloud_vswitch" "foo" {
	  vpc_id = "${alicloud_vpc.foo.id}"
	  name = "${var.name}"
	  cidr_block = "172.16.1.0/24"
	  availability_zone = "${data.alicloud_zones.foo.zones.0.id}"
	}
	resource "alicloud_ots_instance_attachment" "foo" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  vpc_name = "testvpc"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	
	data "alicloud_ots_instance_attachments" "attachments" {
      instance_name = "${alicloud_ots_instance_attachment.foo.instance_name}"
	}
	`, randInt)
}

func testAccCheckAlicloudOtsInstanceAttachmentsDataSource_name_regex_exist(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}

	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "Capacity"
	}

	data "alicloud_zones" "foo" {
	  available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "foo" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}

	resource "alicloud_vswitch" "foo" {
	  vpc_id = "${alicloud_vpc.foo.id}"
	  name = "${var.name}"
	  cidr_block = "172.16.1.0/24"
	  availability_zone = "${data.alicloud_zones.foo.zones.0.id}"
	}
	resource "alicloud_ots_instance_attachment" "foo" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  vpc_name = "testvpc"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	
	data "alicloud_ots_instance_attachments" "attachments" {
      instance_name = "${alicloud_ots_instance_attachment.foo.instance_name}"
      name_regex = "${alicloud_ots_instance_attachment.foo.vpc_name}"
	}
	`, randInt)
}

func testAccCheckAlicloudOtsInstanceAttachmentsDataSource_name_regex_fake(randInt int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%d"
	}

	resource "alicloud_ots_instance" "foo" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "Capacity"
	}

	data "alicloud_zones" "foo" {
	  available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "foo" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}

	resource "alicloud_vswitch" "foo" {
	  vpc_id = "${alicloud_vpc.foo.id}"
	  name = "${var.name}"
	  cidr_block = "172.16.1.0/24"
	  availability_zone = "${data.alicloud_zones.foo.zones.0.id}"
	}
	resource "alicloud_ots_instance_attachment" "foo" {
	  instance_name = "${alicloud_ots_instance.foo.name}"
	  vpc_name = "testvpc"
	  vswitch_id = "${alicloud_vswitch.foo.id}"
	}
	
	data "alicloud_ots_instance_attachments" "attachments" {
      instance_name = "${alicloud_ots_instance_attachment.foo.instance_name}"
	  name_regex = "${alicloud_ots_instance_attachment.foo.vpc_name}-fake"
	}
	`, randInt)
}
