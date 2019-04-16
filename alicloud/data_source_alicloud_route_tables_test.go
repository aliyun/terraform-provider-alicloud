package alicloud

import (
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRouteTablesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouteTablesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_tables.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "tables.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.route_table_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.router_id"),
					resource.TestMatchResourceAttr("data.alicloud_route_tables.foo", "tables.0.name", regexp.MustCompile("^tf-testAcc-for-route-tables-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_route_tables.foo", "tables.0.description", regexp.MustCompile("^tf-testAcc-for-route-tables-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudRouteTablesDataSourceConfig_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_tables.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "tables.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudRouteTablesDataSourceNameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouteTablesDataSourceNameRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_tables.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "tables.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.route_table_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.router_id"),
					resource.TestMatchResourceAttr("data.alicloud_route_tables.foo", "tables.0.name", regexp.MustCompile("^tf-testAcc-for-route-tables-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_route_tables.foo", "tables.0.description", regexp.MustCompile("^tf-testAcc-for-route-tables-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudRouteTablesDataSourceNameRegex_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_tables.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "tables.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudRouteTablesDataSourceIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouteTablesDataSourceIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_tables.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "tables.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.route_table_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.router_id"),
					resource.TestMatchResourceAttr("data.alicloud_route_tables.foo", "tables.0.name", regexp.MustCompile("^tf-testAcc-for-route-tables-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_route_tables.foo", "tables.0.description", regexp.MustCompile("^tf-testAcc-for-route-tables-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudRouteTablesDataSourceIds_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_tables.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "tables.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "names.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudRouteTablesDataSourceVpcId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.RouteTableNoSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouteTablesDataSourceVpcId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_tables.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "tables.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.route_table_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_tables.foo", "tables.0.router_id"),
					resource.TestMatchResourceAttr("data.alicloud_route_tables.foo", "tables.0.name", regexp.MustCompile("^tf-testAcc-for-route-tables-datasourc")),
					resource.TestMatchResourceAttr("data.alicloud_route_tables.foo", "tables.0.description", regexp.MustCompile("^tf-testAcc-for-route-tables-datasourc")),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "names.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudRouteTablesDataSourceVpcId_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_tables.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "tables.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_route_tables.foo", "names.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRouteTablesDataSourceConfig = `
variable "name" {
  default = "tf-testAcc-for-route-tables-datasource"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_route_tables" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name_regex = "${alicloud_route_table.foo.name}"
  ids = ["${alicloud_route_table.foo.id}"]
}
`
const testAccCheckAlicloudRouteTablesDataSourceConfig_mismatch = `
variable "name" {
  default = "tf-testAcc-for-route-tables-datasource"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_route_tables" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}-fake"
  name_regex = "${alicloud_route_table.foo.name}-fake"
  ids = ["${alicloud_route_table.foo.id}-fake"]
}
`

const testAccCheckAlicloudRouteTablesDataSourceNameRegex = `
variable "name" {
  default = "tf-testAcc-for-route-tables-datasource"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_route_tables" "foo" {
  name_regex = "${alicloud_route_table.foo.name}"
}
`

const testAccCheckAlicloudRouteTablesDataSourceNameRegex_mismatch = `
variable "name" {
  default = "tf-testAcc-for-route-tables-datasource"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_route_tables" "foo" {
  name_regex = "${alicloud_route_table.foo.name}-fake"
}
`

const testAccCheckAlicloudRouteTablesDataSourceIds = `
variable "name" {
  default = "tf-testAcc-for-route-tables-datasource"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_route_tables" "foo" {
  ids = ["${alicloud_route_table.foo.id}"]
}
`

const testAccCheckAlicloudRouteTablesDataSourceIds_mismatch = `
variable "name" {
  default = "tf-testAcc-for-route-tables-datasource"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_route_tables" "foo" {
  ids = ["${alicloud_route_table.foo.id}-fake"]
}
`
const testAccCheckAlicloudRouteTablesDataSourceVpcId = `
variable "name" {
  default = "tf-testAcc-for-route-tables-datasource"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_route_tables" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name_regex = "${alicloud_route_table.foo.name}"
}
`
const testAccCheckAlicloudRouteTablesDataSourceVpcId_mismatch = `
variable "name" {
  default = "tf-testAcc-for-route-tables-datasource"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  name = "${var.name}"
  description = "${var.name}"
}

data "alicloud_route_tables" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}-fake"
  name_regex = "${alicloud_route_table.foo.name}-fake"
}
`
