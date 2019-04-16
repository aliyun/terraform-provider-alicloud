package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSnatEntriesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSnatEntriesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_snat_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_snat_entries.foo", "entries.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_snat_entries.foo", "entries.0.snat_ip"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.0.source_cidr", "172.16.0.0/21"),
				),
			},
			{
				Config: testAccCheckAlicloudSnatEntriesDataSourceConfig_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_snat_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudSnatEntriesDataSourceCidr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSnatEntriesDataSourceCidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_snat_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_snat_entries.foo", "entries.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_snat_entries.foo", "entries.0.snat_ip"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.0.source_cidr", "172.16.0.0/21"),
				),
			},
			{
				Config: testAccCheckAlicloudSnatEntriesDataSourceCidr_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_snat_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudSnatEntriesDataSourceIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSnatEntriesDataSourceIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_snat_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_snat_entries.foo", "entries.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_snat_entries.foo", "entries.0.snat_ip"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.0.source_cidr", "172.16.0.0/21"),
				),
			},
			{
				Config: testAccCheckAlicloudSnatEntriesDataSourceIp_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_snat_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_snat_entries.foo", "entries.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSnatEntriesDataSourceConfig = `
variable "name" {
	default = "tf-testAcc-for-snat-entries-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo" {
	snat_table_id = "${alicloud_nat_gateway.foo.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.foo.id}"
	snat_ip = "${alicloud_eip.foo.ip_address}"
}

data "alicloud_snat_entries" "foo" {
    source_cidr = "${alicloud_vswitch.foo.cidr_block}"
    snat_ip = "${alicloud_snat_entry.foo.snat_ip}"
    snat_table_id = "${alicloud_snat_entry.foo.snat_table_id}"
}
`
const testAccCheckAlicloudSnatEntriesDataSourceConfig_mismatch = `
variable "name" {
	default = "tf-testAcc-for-snat-entries-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo" {
	snat_table_id = "${alicloud_nat_gateway.foo.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.foo.id}"
	snat_ip = "${alicloud_eip.foo.ip_address}"
}

data "alicloud_snat_entries" "foo" {
    source_cidr = "${alicloud_vswitch.foo.cidr_block}-fake"
    snat_ip = "${alicloud_snat_entry.foo.snat_ip}-fake"
    snat_table_id = "${alicloud_snat_entry.foo.snat_table_id}-fake"
}
`

const testAccCheckAlicloudSnatEntriesDataSourceCidr = `
variable "name" {
	default = "tf-testAcc-for-snat-entries-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo" {
	snat_table_id = "${alicloud_nat_gateway.foo.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.foo.id}"
	snat_ip = "${alicloud_eip.foo.ip_address}"
}

data "alicloud_snat_entries" "foo" {
    source_cidr = "${alicloud_vswitch.foo.cidr_block}"
    snat_table_id = "${alicloud_snat_entry.foo.snat_table_id}"
}
`
const testAccCheckAlicloudSnatEntriesDataSourceCidr_mismatch = `
variable "name" {
	default = "tf-testAcc-for-snat-entries-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo" {
	snat_table_id = "${alicloud_nat_gateway.foo.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.foo.id}"
	snat_ip = "${alicloud_eip.foo.ip_address}"
}

data "alicloud_snat_entries" "foo" {
    source_cidr = "${alicloud_vswitch.foo.cidr_block}-fake"
    snat_table_id = "${alicloud_snat_entry.foo.snat_table_id}"
}
`
const testAccCheckAlicloudSnatEntriesDataSourceIp = `
variable "name" {
	default = "tf-testAcc-for-snat-entries-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo" {
	snat_table_id = "${alicloud_nat_gateway.foo.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.foo.id}"
	snat_ip = "${alicloud_eip.foo.ip_address}"
}

data "alicloud_snat_entries" "foo" {
    snat_ip = "${alicloud_snat_entry.foo.snat_ip}"
    snat_table_id = "${alicloud_snat_entry.foo.snat_table_id}"
}
`
const testAccCheckAlicloudSnatEntriesDataSourceIp_mismatch = `
variable "name" {
	default = "tf-testAcc-for-snat-entries-datasource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
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

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "foo" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "foo" {
	allocation_id = "${alicloud_eip.foo.id}"
	instance_id = "${alicloud_nat_gateway.foo.id}"
}

resource "alicloud_snat_entry" "foo" {
	snat_table_id = "${alicloud_nat_gateway.foo.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.foo.id}"
	snat_ip = "${alicloud_eip.foo.ip_address}"
}

data "alicloud_snat_entries" "foo" {
    snat_ip = "${alicloud_snat_entry.foo.snat_ip}-fake"
    snat_table_id = "${alicloud_snat_entry.foo.snat_table_id}"
}
`
