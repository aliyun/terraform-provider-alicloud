package vpngateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudVPNCgwsDataSourceBasic(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnCustomerGatewaysConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vpn_customer_gateway.default.id}" ]`,
		}),
		FakeConfig: testAccCheckAlicloudVpnCustomerGatewaysConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vpn_customer_gateway.default.id}_fake" ]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnCustomerGatewaysConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_customer_gateway.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnCustomerGatewaysConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_customer_gateway.default.name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnCustomerGatewaysConfig(rand, map[string]string{
			"ids":        `[ "${alicloud_vpn_customer_gateway.default.id}" ]`,
			"name_regex": `"${alicloud_vpn_customer_gateway.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnCustomerGatewaysConfig(rand, map[string]string{
			"ids":        `[ "${alicloud_vpn_customer_gateway.default.id}" ]`,
			"name_regex": `"${alicloud_vpn_customer_gateway.default.name}_fake"`,
		}),
	}

	vpnCustomerGatewaysCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudVpnCustomerGatewaysConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
resource "alicloud_vpn_customer_gateway" "default" {
	name = "tf-testAccVpnCgwNameDataResource%d"
	ip_address = "40.104.22.228"
	description = "tf-testAccVpnCgwNameDataResource%d"
}

data "alicloud_vpn_customer_gateways" "default" {
	%s
}
`, rand, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVpnCustomerGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":             "1",
		"ids.#":                  "1",
		"names.#":                "1",
		"gateways.0.id":          acctest.CHECKSET,
		"gateways.0.name":        fmt.Sprintf("tf-testAccVpnCgwNameDataResource%d", rand),
		"gateways.0.ip_address":  "40.104.22.228",
		"gateways.0.description": fmt.Sprintf("tf-testAccVpnCgwNameDataResource%d", rand),
		"gateways.0.create_time": acctest.CHECKSET,
	}
}

var fakeVpnCustomerGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"gateways.#": "0",
	}
}

var vpnCustomerGatewaysCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_vpn_customer_gateways.default",
	ExistMapFunc: existVpnCustomerGatewaysMapFunc,
	FakeMapFunc:  fakeVpnCustomerGatewaysMapFunc,
}
