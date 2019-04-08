package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSslVpnServersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSslVpnServersDataCfg_id,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_servers.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.name", "tf-testAccSslVpnServersDataResource"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.vpn_gateway_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.compress", "false"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.port", "1194"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.proto", "UDP"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.local_subnet", "172.16.1.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.client_ip_pool", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.internet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.max_connections"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.connections"),
				),
			},
			{
				Config: testAccCheckAlicloudSslVpnServersDataCfg_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_servers.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.name", "tf-testAccSslVpnServersDataResource"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.vpn_gateway_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.compress", "false"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.port", "1194"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.proto", "UDP"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.local_subnet", "172.16.1.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.client_ip_pool", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.internet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.max_connections"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.connections"),
				),
			},
			{
				Config: testAccCheckAlicloudSslVpnServersDataCfg_both,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_servers.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.name", "tf-testAccSslVpnServersDataResource"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.vpn_gateway_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.compress", "false"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.port", "1194"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.proto", "UDP"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.local_subnet", "172.16.1.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.client_ip_pool", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.internet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.max_connections"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.connections"),
				),
			},
			{
				Config: testAccCheckAlicloudSslVpnServersDataCfg_vpnid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_servers.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.name", "tf-testAccSslVpnServersDataResource"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.vpn_gateway_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.create_time"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.compress", "false"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.port", "1194"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.proto", "UDP"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.local_subnet", "172.16.1.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.client_ip_pool", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.internet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.max_connections"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "servers.0.connections"),
				),
			},
			{
				Config: testAccCheckAlicloudSslVpnServersDataCfg_miss_id,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_servers.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "names.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.vpn_gateway_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.create_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.compress"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.cipher"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.port"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.proto"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.local_subnet"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.client_ip_pool"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.internet_ip"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.max_connections"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.connections"),
				),
			},
			{
				Config: testAccCheckAlicloudSslVpnServersDataCfg_miss_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_servers.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "names.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.vpn_gateway_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.create_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.compress"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.cipher"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.port"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.proto"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.local_subnet"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.client_ip_pool"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.internet_ip"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.max_connections"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.connections"),
				),
			},
			{
				Config: testAccCheckAlicloudSslVpnServersDataCfg_miss_vpnid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_servers.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "names.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.vpn_gateway_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.create_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.compress"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.cipher"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.port"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.proto"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.local_subnet"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.client_ip_pool"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.internet_ip"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.max_connections"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_servers.foo", "servers.0.connections"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSslVpnServersDataCfg_name = `
variable "name" {
	default = "tf-testAccSslVpnServersDataResource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name="${var.name}"
	vpn_gateway_id="${alicloud_vpn_gateway.foo.id}"
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

data "alicloud_ssl_vpn_servers" "foo" {
	name_regex = "${alicloud_ssl_vpn_server.foo.name}"
}
`

const testAccCheckAlicloudSslVpnServersDataCfg_id = `
variable "name" {
	default = "tf-testAccSslVpnServersDataResource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name="${var.name}"
	vpn_gateway_id="${alicloud_vpn_gateway.foo.id}"
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

data "alicloud_ssl_vpn_servers" "foo" {
	ids = ["${alicloud_ssl_vpn_server.foo.id}"]
}
`

const testAccCheckAlicloudSslVpnServersDataCfg_both = `
variable "name" {
	default = "tf-testAccSslVpnServersDataResource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name="${var.name}"
	vpn_gateway_id="${alicloud_vpn_gateway.foo.id}"
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

data "alicloud_ssl_vpn_servers" "foo" {
	name_regex = "${alicloud_ssl_vpn_server.foo.name}"
	ids = ["${alicloud_ssl_vpn_server.foo.id}"]
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
}
`

const testAccCheckAlicloudSslVpnServersDataCfg_vpnid = `
variable "name" {
	default = "tf-testAccSslVpnServersDataResource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name="${var.name}"
	vpn_gateway_id="${alicloud_vpn_gateway.foo.id}"
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

data "alicloud_ssl_vpn_servers" "foo" {
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
}
`

const testAccCheckAlicloudSslVpnServersDataCfg_miss_id = `
variable "name" {
	default = "tf-testAccSslVpnServersDataResource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name="${var.name}"
	vpn_gateway_id="${alicloud_vpn_gateway.foo.id}"
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

data "alicloud_ssl_vpn_servers" "foo" {
	ids = ["${alicloud_ssl_vpn_server.foo.id}-fake"]
}
`

const testAccCheckAlicloudSslVpnServersDataCfg_miss_name = `
variable "name" {
	default = "tf-testAccSslVpnServersDataResource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name="${var.name}"
	vpn_gateway_id="${alicloud_vpn_gateway.foo.id}"
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

data "alicloud_ssl_vpn_servers" "foo" {
	name_regex = "${alicloud_ssl_vpn_server.foo.name}-fake"
}
`

const testAccCheckAlicloudSslVpnServersDataCfg_miss_vpnid = `
variable "name" {
	default = "tf-testAccSslVpnServersDataResource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name="${var.name}"
	vpn_gateway_id="${alicloud_vpn_gateway.foo.id}"
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

data "alicloud_ssl_vpn_servers" "foo" {
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}-fake"
}
`
