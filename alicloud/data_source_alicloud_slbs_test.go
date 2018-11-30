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
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.name", "tf-testAccCheckAlicloudSlbsDataSourceBasic"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.network_type", "vpc"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.vswitch_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.address"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.internet", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_slbs.balancers", "slbs.0.creation_time"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.tags.tag_a", "1"),
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
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.name", "tf-testAccCheckAlicloudSlbsDataSourceFilterByIds"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.tags.tag_a", "1"),
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
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.name", "tf-testAccCheckAlicloudSlbsDataSourceFilterByAllFields"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.0.tags.tag_a", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbsDataSource_filterByTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbsDataSourceFilterByTags,
				Check: resource.ComposeTestCheckFunc(
					// common tag: tag_a="tag_a_1"
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers-1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-1", "slbs.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-1", "slbs.0.name", "tf-testAccCheckAlicloudSlbsDataSourceFilterByTags"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-1", "slbs.0.tags.tag_a", "tag_a_1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-1", "slbs.1.name", "tf-testAccCheckAlicloudSlbsDataSourceFilterByTags"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-1", "slbs.1.tags.tag_a", "tag_a_1"),

					// single instance tag: tag_f="tag_f_6"
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers-2"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-2", "slbs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-2", "slbs.0.name", "tf-testAccCheckAlicloudSlbsDataSourceFilterByTags"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-2", "slbs.0.tags.tag_a", "tag_a_1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-2", "slbs.0.tags.tag_f", "tag_f_6"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-2", "slbs.0.tags.tag_h", "tag_h_8"),

					// single instance tag: tag_f="tag_f_66"
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers-3"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-3", "slbs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-3", "slbs.0.name", "tf-testAccCheckAlicloudSlbsDataSourceFilterByTags"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-3", "slbs.0.tags.tag_a", "tag_a_1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-3", "slbs.0.tags.tag_f", "tag_f_66"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-3", "slbs.0.tags.tag_i", "tag_i_11"),

					// single instance tag: tag_i="tag_i_11"
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers-4"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-4", "slbs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-4", "slbs.0.name", "tf-testAccCheckAlicloudSlbsDataSourceFilterByTags"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-4", "slbs.0.tags.tag_a", "tag_a_1"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-4", "slbs.0.tags.tag_f", "tag_f_66"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-4", "slbs.0.tags.tag_i", "tag_i_11"),

					// single instance tag: tag_a="tag_a_1" or tag_tag_f="tag_f_6"
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers-5"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-5", "slbs.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-5", "slbs.0.name", "tf-testAccCheckAlicloudSlbsDataSourceFilterByTags"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-5", "slbs.0.tags.tag_a", "tag_a_1"),

					// single instance tag: tag_a="tag_a_1" or tag_f="tag_f_66"
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers-6"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-6", "slbs.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-6", "slbs.0.name", "tf-testAccCheckAlicloudSlbsDataSourceFilterByTags"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers-6", "slbs.0.tags.tag_a", "tag_a_1"),
				),
			},
		},
	})
}

func TestAccAlicloudSlbsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSlbsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_slbs.balancers"),
					resource.TestCheckResourceAttr("data.alicloud_slbs.balancers", "slbs.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.region_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.master_availability_zone"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.slave_availability_zone"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.network_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.vpc_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.vswitch_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.address"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.internet"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.creation_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_slbs.balancers", "slbs.0.tags"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSlbsDataSourceBasic = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbsDataSourceBasic"
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
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
  }
}

data "alicloud_slbs" "balancers" {
  name_regex = "${alicloud_slb.sample_slb.name}"
}
`

const testAccCheckAlicloudSlbsDataSourceFilterByIds = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbsDataSourceFilterByIds"
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
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
  }
}

data "alicloud_slbs" "balancers" {
  ids = ["${alicloud_slb.sample_slb.id}"]
}
`

const testAccCheckAlicloudSlbsDataSourceFilterByAllFields = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbsDataSourceFilterByAllFields"
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
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
  }
}

data "alicloud_slbs" "balancers" {
  ids = ["${alicloud_slb.sample_slb.id}"]
  name_regex = "${alicloud_slb.sample_slb.name}"
  network_type = "vpc"
  vpc_id = "${alicloud_vpc.sample_vpc.id}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
  address = "${alicloud_slb.sample_slb.address}"
  tags = {
    tag_a = 1
  }
}
`

const testAccCheckAlicloudSlbsDataSourceFilterByTags = `
variable "name" {
	default = "tf-testAccCheckAlicloudSlbsDataSourceFilterByTags"
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

resource "alicloud_slb" "sample_slb-1" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
  tags = {
    tag_a = "tag_a_1"
    tag_b = "tag_b_2"
    tag_c = "tag_c_3"
    tag_f = "tag_f_6"
    tag_g = "tag_g_7"
    tag_h = "tag_h_8"
  }
}

resource "alicloud_slb" "sample_slb-2" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.sample_vswitch.id}"
  tags = {
    tag_a = "tag_a_1"
    tag_b = "tag_b_2"
    tag_c = "tag_c_3"
    tag_f = "tag_f_66"
    tag_i = "tag_i_11"
  }
}

data "alicloud_slbs" "balancers-1" {
  name_regex = "${alicloud_slb.sample_slb-1.name}"
  tags = {
    tag_a = "tag_a_1"
  }
}

data "alicloud_slbs" "balancers-2" {
  name_regex = "${alicloud_slb.sample_slb-2.name}"
  tags = {
    tag_f = "tag_f_6" 
  }
}

data "alicloud_slbs" "balancers-3" {
  name_regex = "${alicloud_slb.sample_slb-1.name}"
  tags = {
    tag_f = "tag_f_66"
  }
}

data "alicloud_slbs" "balancers-4" {
  name_regex = "${alicloud_slb.sample_slb-1.name}"
  tags = {
    tag_i = "tag_i_11"
  }
}

data "alicloud_slbs" "balancers-5" {
  name_regex = "${alicloud_slb.sample_slb-1.name}"
  tags = {
    tag_a = "tag_a_1"
    tag_f = "tag_f_6"
  }
}

data "alicloud_slbs" "balancers-6" {
  name_regex = "${alicloud_slb.sample_slb-1.name}"
  tags = {
    tag_a = "tag_a_1"
    tag_f = "tag_f_66"
  }
}
`
const testAccCheckAlicloudSlbsDataSourceEmpty = `
data "alicloud_slbs" "balancers" {
  name_regex = "^tf-testAcc-fake-name"
}
`
