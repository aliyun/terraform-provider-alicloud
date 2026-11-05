package ecd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudECDNetworkPackagesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdNetworkPackagesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_network_package.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEcdNetworkPackagesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_network_package.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdNetworkPackagesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_network_package.default.id}"]`,
			"status": `"InUse"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdNetworkPackagesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_network_package.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	var existAlicloudEventBridgeEventBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"ids.0":                           acctest.CHECKSET,
			"packages.#":                      "1",
			"packages.0.id":                   acctest.CHECKSET,
			"packages.0.bandwidth":            "10",
			"packages.0.internet_charge_type": "PayByTraffic",
			"packages.0.status":               acctest.CHECKSET,
			"packages.0.office_site_id":       acctest.CHECKSET,
			"packages.0.office_site_name":     fmt.Sprintf("tf-testacc-ecdnetworkpackage%d", rand),
			"packages.0.expired_time":         acctest.CHECKSET,
			"packages.0.create_time":          acctest.CHECKSET,
			"packages.0.network_package_id":   acctest.CHECKSET,
			"packages.0.eip_addresses.#":      "1",
			"packages.0.eip_addresses.0":      acctest.CHECKSET,
		}
	}
	var fakeAlicloudEventBridgeEventBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudEventBridgeEventBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecd_network_packages.default",
		ExistMapFunc: existAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
	}

	alicloudEventBridgeEventBusesCheckInfo.DataSourceTestCheck(t, rand, idsConf, statusConf)
}
func testAccCheckAlicloudEcdNetworkPackagesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testacc-ecdnetworkpackage%d"
}


resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
}

resource "alicloud_ecd_network_package" "default" {
  bandwidth = "10"
  office_site_id = alicloud_ecd_simple_office_site.default.id
}


data "alicloud_ecd_network_packages" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
