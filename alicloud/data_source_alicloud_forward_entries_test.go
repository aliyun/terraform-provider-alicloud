package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudForwardEntriesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudForwardEntriesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_forward_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_forward_entries.foo", "entries.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_forward_entries.foo", "entries.0.external_ip"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.external_port", "80"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.internal_ip", "172.16.0.3"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.internal_port", "8080"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "ids.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudForwardEntriesDataSourceConfig_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_forward_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudForwardEntriesDataSourceIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudForwardEntriesDataSourceIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_forward_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_forward_entries.foo", "entries.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_forward_entries.foo", "entries.0.external_ip"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.external_port", "80"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.internal_ip", "172.16.0.3"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.internal_port", "8080"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "ids.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudForwardEntriesDataSourceIds_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_forward_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudForwardEntriesDataSourceExternalIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudForwardEntriesDataSourceExternalIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_forward_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_forward_entries.foo", "entries.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_forward_entries.foo", "entries.0.external_ip"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.external_port", "80"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.internal_ip", "172.16.0.3"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.internal_port", "8080"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "ids.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudForwardEntriesDataSourceExternalIp_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_forward_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudForwardEntriesDataSourceInternalIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudForwardEntriesDataSourceInternalIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_forward_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_forward_entries.foo", "entries.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_forward_entries.foo", "entries.0.external_ip"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.external_port", "80"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.internal_ip", "172.16.0.3"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.internal_port", "8080"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.0.status", "Available"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "ids.#", "1"),
				),
			},
			{
				Config: testAccCheckAlicloudForwardEntriesDataSourceInternalIp_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_forward_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "entries.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_forward_entries.foo", "ids.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudForwardEntriesDataSourceConfig = `
variable "name" {
	default = "tf-testAcc-for-forward-entries-datasource"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
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

resource "alicloud_forward_entry" "foo" {
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "foo" {
	internal_ip = "${alicloud_forward_entry.foo.internal_ip}"
    external_ip = "${alicloud_forward_entry.foo.external_ip}"
    ids = ["${alicloud_forward_entry.foo.id}"]
    forward_table_id = "${alicloud_forward_entry.foo.forward_table_id}"
}
`

const testAccCheckAlicloudForwardEntriesDataSourceConfig_mismatch = `
variable "name" {
	default = "tf-testAcc-for-forward-entries-datasource"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
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

resource "alicloud_forward_entry" "foo" {
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "foo" {
	internal_ip = "${alicloud_forward_entry.foo.internal_ip}-fake"
    external_ip = "${alicloud_forward_entry.foo.external_ip}-fake"
    ids = ["${alicloud_forward_entry.foo.id}-fake"]
    forward_table_id = "${alicloud_forward_entry.foo.forward_table_id}"
}
`

const testAccCheckAlicloudForwardEntriesDataSourceIds = `
variable "name" {
	default = "tf-testAcc-for-forward-entries-datasource"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
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

resource "alicloud_forward_entry" "foo" {
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "foo" {
    ids = ["${alicloud_forward_entry.foo.id}"]
    forward_table_id = "${alicloud_forward_entry.foo.forward_table_id}"
}
`

const testAccCheckAlicloudForwardEntriesDataSourceIds_mismatch = `
variable "name" {
	default = "tf-testAcc-for-forward-entries-datasource"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
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

resource "alicloud_forward_entry" "foo" {
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "foo" {
    ids = ["${alicloud_forward_entry.foo.id}-fake"]
    forward_table_id = "${alicloud_forward_entry.foo.forward_table_id}"
}
`

const testAccCheckAlicloudForwardEntriesDataSourceExternalIp = `
variable "name" {
	default = "tf-testAcc-for-forward-entries-datasource"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
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

resource "alicloud_forward_entry" "foo" {
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "foo" {
    external_ip = "${alicloud_forward_entry.foo.external_ip}"
    forward_table_id = "${alicloud_forward_entry.foo.forward_table_id}"
}
`

const testAccCheckAlicloudForwardEntriesDataSourceExternalIp_mismatch = `
variable "name" {
	default = "tf-testAcc-for-forward-entries-datasource"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
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

resource "alicloud_forward_entry" "foo" {
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "foo" {
    external_ip = "${alicloud_forward_entry.foo.external_ip}-fake"
    forward_table_id = "${alicloud_forward_entry.foo.forward_table_id}"
}
`

const testAccCheckAlicloudForwardEntriesDataSourceInternalIp = `
variable "name" {
	default = "tf-testAcc-for-forward-entries-datasource"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
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

resource "alicloud_forward_entry" "foo" {
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "foo" {
	internal_ip = "${alicloud_forward_entry.foo.internal_ip}"
    forward_table_id = "${alicloud_forward_entry.foo.forward_table_id}"
}
`

const testAccCheckAlicloudForwardEntriesDataSourceInternalIp_mismatch = `
variable "name" {
	default = "tf-testAcc-for-forward-entries-datasource"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
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

resource "alicloud_forward_entry" "foo" {
	forward_table_id = "${alicloud_nat_gateway.foo.forward_table_ids}"
	external_ip = "${alicloud_eip.foo.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "alicloud_forward_entries" "foo" {
	internal_ip = "${alicloud_forward_entry.foo.internal_ip}-fake"
    forward_table_id = "${alicloud_forward_entry.foo.forward_table_id}"
}
`
