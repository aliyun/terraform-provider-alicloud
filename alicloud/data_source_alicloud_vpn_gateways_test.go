package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVpnsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpnsDataCfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpn_gateways.vpn_gateways"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "names.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.internet_ip"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.create_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.end_time"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.name", "tf-testAccVpnDatasource"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.specification", "10M"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.description", "test_create_description"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.enable_ssl", "enable"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.enable_ipsec", "enable"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.status", "Active"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.business_status", "Normal"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.ssl_connections", "5"),
				),
			},
		},
	})
}

func TestAccAlicloudVpnsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpnsDataEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_vpn_gateways.vpn_gateways"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "names.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.vpc_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.internet_ip"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.create_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.end_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.specification"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.enable_ssl"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.enable_ipsec"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.business_status"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.instance_charge_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_vpn_gateways.vpn_gateways", "gateways.0.ssl_connections"),
				),
			},
		},
	})
}

const testAccCheckAlicloudVpnsDataCfg = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "tf-testAccVpcConfig"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "tf-testAccVswitchConfig"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "tf-testAccVpnDatasource"
	vpc_id = "${alicloud_vswitch.foo.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

data "alicloud_vpn_gateways" "vpn_gateways" {
	vpc_id = "${alicloud_vpc.foo.id}"
	ids = ["${alicloud_vpn_gateway.foo.id}"]
	status = "Active"
	business_status = "Normal"
	name_regex = "tf-testAcc*"
}
`

const testAccCheckAlicloudVpnsDataEmpty = `
data "alicloud_vpn_gateways" "vpn_gateways" {
	name_regex = "tf-testAcc-fake"
}
`

func TestAccAlicloudVpnGatewaysDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vpn_gateway.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vpn_gateway.default.id}_fake" ]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}_fake"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alicloud_vpn_gateway.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alicloud_vpn_gateway.default.vpc_id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}"`,
			"status":     `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}"`,
			"status":     `"Init"`,
		}),
	}

	businessStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"business_status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"business_status": `"FinancialLocked"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids":             `[ "${alicloud_vpn_gateway.default.id}" ]`,
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"vpc_id":          `"${alicloud_vpn_gateway.default.vpc_id}"`,
			"status":          `"Active"`,
			"business_status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids":             `[ "${alicloud_vpn_gateway.default.id}" ]`,
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"vpc_id":          `"${alicloud_vpn_gateway.default.vpc_id}"`,
			"status":          `"Active"`,
			"business_status": `"FinancialLocked"`,
		}),
	}

	vpnGatewaysCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, vpcIdConf, statusConf, businessStatusConf, allConf)
}

func testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpcGatewayConfig%d"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "${var.name}"
}

data "alicloud_vpn_gateways" "default" {
	%s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVpnGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":                      "1",
		"ids.#":                           "1",
		"names.#":                         "1",
		"gateways.0.id":                   CHECKSET,
		"gateways.0.vpc_id":               CHECKSET,
		"gateways.0.internet_ip":          CHECKSET,
		"gateways.0.create_time":          CHECKSET,
		"gateways.0.end_time":             CHECKSET,
		"gateways.0.name":                 fmt.Sprintf("tf-testAccVpcGatewayConfig%d", rand),
		"gateways.0.specification":        "10M",
		"gateways.0.description":          fmt.Sprintf("tf-testAccVpcGatewayConfig%d", rand),
		"gateways.0.enable_ssl":           "enable",
		"gateways.0.enable_ipsec":         "enable",
		"gateways.0.status":               "Active",
		"gateways.0.business_status":      "Normal",
		"gateways.0.instance_charge_type": string(PostPaid),
		"gateways.0.ssl_connections":      "5",
	}
}

var fakeVpnGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"gateways.#": "0",
	}
}

var vpnGatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpn_gateways.default",
	existMapFunc: existVpnGatewaysMapFunc,
	fakeMapFunc:  fakeVpnGatewaysMapFunc,
}
