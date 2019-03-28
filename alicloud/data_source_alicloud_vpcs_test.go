package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVpcsDataSource_cidr_block(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceCidrBlockConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.0", "tf-testAccVpcsdatasourceNameRegex"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vpc_name", "tf-testAccVpcsdatasourceNameRegex"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vswitch_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.vrouter_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.route_table_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.creation_time"),
				),
			},
			{
				Config: testAccCheckAlicloudVpcsDataSourceCidrBlockConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "0"),
				),
			},
		},
	})
}

func TestAccCheckAlicloudVpcsDataSource_Status(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceStatus,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.0", "tf-testAccVpcsdatasourceStatus"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vpc_name", "tf-testAccVpcsdatasourceStatus"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vswitch_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.vrouter_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.route_table_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.creation_time"),
				),
			},
			{
				Config: testAccCheckAlicloudVpcsDataSourceStatusEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "0"),
				),
			},
		},
	})
}
func TestAccCheckAlicloudVpcsDataSource_Is_Default(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceIsDefault,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.0", "tf-testAccVpcsdatasourceIsDefault"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vpc_name", "tf-testAccVpcsdatasourceIsDefault"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vswitch_ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.vrouter_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.route_table_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.creation_time"),
				),
			},
			{
				Config: testAccCheckAlicloudVpcsDataSourceIsDefaultEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "0"),
				),
			},
		},
	})
}
func TestAccCheckAlicloudVpcsDataSource_VSwitch_ID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceVSwitchID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.0", "tf-testAccVpcsdatasourceVSwitchID"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vpc_name", "tf-testAccVpcsdatasourceVSwitchID"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vswitch_ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.cidr_block", "172.16.0.0/12"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.vrouter_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.route_table_id"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.is_default", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpcs.vpc", "vpcs.0.creation_time"),
				),
			},
			{
				Config: testAccCheckAlicloudVpcsDataSourceVSwitchIDEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "0"),
				),
			},
		},
	})
}
func TestAccAlicloudVpcsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpcs.vpc"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "names.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpcs.vpc", "vpcs.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.region_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vpc_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vswitch_ids"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.cidr_block"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.vrouter_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.route_table_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.is_default"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpcs.vpc", "vpcs.0.create_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpcsDataSourceCidrBlockConfig = `
variable "name" {
  default = "tf-testAccVpcsdatasourceNameRegex"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vpcs" "vpc" {
  name_regex = "testAccVpcsdatasource*"
  cidr_block = "${alicloud_vpc.foo.cidr_block}"
}
`

const testAccCheckAlicloudVpcsDataSourceCidrBlockConfigEmpty = `
variable "name" {
  default = "tf-testAccVpcsdatasourceNameRegex"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vpcs" "vpc" {
  name_regex = "testAccVpcsdatasource*"
  cidr_block = "${alicloud_vpc.foo.cidr_block}-fake"
}
`

const testAccCheckAlicloudVpcsDataSourceStatus = `
variable "name" {
  default = "tf-testAccVpcsdatasourceStatus"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vpcs" "vpc" {
  name_regex = "${alicloud_vpc.foo.name}"
  status = "Available"
}
`

const testAccCheckAlicloudVpcsDataSourceStatusEmpty = `
variable "name" {
  default = "tf-testAccVpcsdatasourceStatus"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vpcs" "vpc" {
  name_regex = "${alicloud_vpc.foo.name}"
  status = "UnAvailable"
}
`
const testAccCheckAlicloudVpcsDataSourceIsDefault = `
variable "name" {
  default = "tf-testAccVpcsdatasourceIsDefault"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vpcs" "vpc" {
  name_regex = "${alicloud_vpc.foo.name}"
  is_default="false"
}
`

const testAccCheckAlicloudVpcsDataSourceIsDefaultEmpty = `
variable "name" {
  default = "tf-testAccVpcsdatasourceIsDefault"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_vpcs" "vpc" {
  name_regex = "${alicloud_vpc.foo.name}"
  is_default="true"
}
`

const testAccCheckAlicloudVpcsDataSourceVSwitchID = `
variable "name" {
  default = "tf-testAccVpcsdatasourceVSwitchID"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_zones" "default" {}
resource "alicloud_vswitch" "vswitch" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
	vpc_id = "${alicloud_vpc.foo.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_vpcs" "vpc" {
  vswitch_id="${alicloud_vswitch.vswitch.id}"
}
`

const testAccCheckAlicloudVpcsDataSourceVSwitchIDEmpty = `
variable "name" {
  default = "tf-testAccVpcsdatasourceVSwitchID"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
data "alicloud_zones" "default" {}
resource "alicloud_vswitch" "vswitch" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
	vpc_id = "${alicloud_vpc.foo.id}"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_vpcs" "vpc" {
  vswitch_id="${alicloud_vswitch.vswitch.id}-fake"
}
`
const testAccCheckAlicloudVpcsDataSourceEmpty = `
data "alicloud_vpcs" "vpc" {
  name_regex = "^tf-fake-name"
}
`
