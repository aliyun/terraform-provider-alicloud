package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSlbsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.master_availability_zone"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.slave_availability_zone"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.status", "active"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.name", "testAccCheckAlicloudSlbsDataSourceBasic"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.network_type", "vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.address"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.internet", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.creation_time"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbsDataSource_filterByIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbsDataSourceFilterByIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.name", "testAccCheckAlicloudSlbsDataSourceFilterByIds"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbsDataSource_filterByAllFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbsDataSourceFilterByAllFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.name", "testAccCheckAlicloudSlbsDataSourceFilterByAllFields"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSlbsDataSourceBasic = `
variable "name" {
	default = "testAccCheckAlicloudSlbsDataSourceBasic"
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
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

data "alicloud_slbs" "balancers" {
  name_regex = "${alicloud_slb.sample_slb.name}"
}
`

const testAccCheckAlicloudSlbsDataSourceFilterByIds = `
variable "name" {
	default = "testAccCheckAlicloudSlbsDataSourceFilterByIds"
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
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

data "alicloud_slbs" "balancers" {
  ids = ["${alicloud_slb.sample_slb.id}"]
}
`

const testAccCheckAlicloudSlbsDataSourceFilterByAllFields = `
variable "name" {
	default = "testAccCheckAlicloudSlbsDataSourceFilterByAllFields"
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
}

resource "alicloud_slb" "sample_slb" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
}

data "alicloud_slbs" "balancers" {
  ids = ["${alicloud_slb.sample_slb.id}"]
  name_regex = "${alicloud_slb.sample_slb.name}"
  network_type = "vpc"
  vpc_id = "${alicloud_vpc.sample_vpc.id}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
  address = "${alicloud_slb.sample_slb.address}"
}
`
