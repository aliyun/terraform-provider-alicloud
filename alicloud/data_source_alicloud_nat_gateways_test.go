package alicloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudNatGatewaysDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudNatGatewaysDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nat_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.spec", "Small"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.status", "Available"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.forward_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.snat_table_id"),
					resource.TestMatchResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.name", regexp.MustCompile("^tf-testAcc-for-nat-gateways-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.description", regexp.MustCompile("^tf-testAcc-for-nat-gateways-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudNatGatewaysDataSourceConfig_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nat_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNatGatewaysDataSourceNameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudNatGatewaysDataSourceNameRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nat_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.spec", "Small"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.status", "Available"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.forward_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.snat_table_id"),
					resource.TestMatchResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.name", regexp.MustCompile("^tf-testAcc-for-nat-gateways-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.description", regexp.MustCompile("^tf-testAcc-for-nat-gateways-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudNatGatewaysDataSourceNameRegex_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nat_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNatGatewaysDataSourceIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudNatGatewaysDataSourceIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nat_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.spec", "Small"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.status", "Available"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.forward_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.snat_table_id"),
					resource.TestMatchResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.name", regexp.MustCompile("^tf-testAcc-for-nat-gateways-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.description", regexp.MustCompile("^tf-testAcc-for-nat-gateways-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudNatGatewaysDataSourceIds_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nat_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudNatGatewaysDataSourceVpcId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudNatGatewaysDataSourceVpcId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nat_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.spec", "Small"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.status", "Available"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.forward_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nat_gateways.foo", "gateways.0.snat_table_id"),
					resource.TestMatchResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.name", regexp.MustCompile("^tf-testAcc-for-nat-gateways-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_nat_gateways.foo", "gateways.0.description", regexp.MustCompile("^tf-testAcc-for-nat-gateways-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudNatGatewaysDataSourceVpcId_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nat_gateways.foo"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "gateways.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_nat_gateways.foo", "names.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudNatGatewaysDataSourceConfig = `
variable "name" {
  default = "tf-testAcc-for-nat-gateways-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
    description = "${var.name}"
}

data "alicloud_nat_gateways" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
    name_regex = "${alicloud_nat_gateway.foo.name}"
    ids = ["${alicloud_nat_gateway.foo.id}"]
}
`

const testAccCheckAlicloudNatGatewaysDataSourceConfig_mismatch = `
variable "name" {
  default = "tf-testAcc-for-nat-gateways-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
    description = "${var.name}"
}

data "alicloud_nat_gateways" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}-fake"
    name_regex = "${alicloud_nat_gateway.foo.name}-fake"
    ids = ["${alicloud_nat_gateway.foo.id}-fake"]
}
`

const testAccCheckAlicloudNatGatewaysDataSourceNameRegex = `
variable "name" {
  default = "tf-testAcc-for-nat-gateways-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
    description = "${var.name}"
}

data "alicloud_nat_gateways" "foo" {
    name_regex = "${alicloud_nat_gateway.foo.name}"
}
`

const testAccCheckAlicloudNatGatewaysDataSourceNameRegex_mismatch = `
variable "name" {
  default = "tf-testAcc-for-nat-gateways-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
    description = "${var.name}"
}

data "alicloud_nat_gateways" "foo" {
    name_regex = "${alicloud_nat_gateway.foo.name}-fake"
}
`

const testAccCheckAlicloudNatGatewaysDataSourceIds = `
variable "name" {
  default = "tf-testAcc-for-nat-gateways-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
    description = "${var.name}"
}

data "alicloud_nat_gateways" "foo" {
    ids = ["${alicloud_nat_gateway.foo.id}"]
}
`

const testAccCheckAlicloudNatGatewaysDataSourceIds_mismatch = `
variable "name" {
  default = "tf-testAcc-for-nat-gateways-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
    description = "${var.name}"
}

data "alicloud_nat_gateways" "foo" {
    ids = ["${alicloud_nat_gateway.foo.id}-fake"]
}
`

const testAccCheckAlicloudNatGatewaysDataSourceVpcId = `
variable "name" {
  default = "tf-testAcc-for-nat-gateways-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
    description = "${var.name}"
}

data "alicloud_nat_gateways" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
    name_regex = "${alicloud_nat_gateway.foo.name}"
}
`

const testAccCheckAlicloudNatGatewaysDataSourceVpcId_mismatch = `
variable "name" {
  default = "tf-testAcc-for-nat-gateways-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
    description = "${var.name}"
}

data "alicloud_nat_gateways" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}-fake"
    name_regex = "${alicloud_nat_gateway.foo.name}-fake"
}
`
