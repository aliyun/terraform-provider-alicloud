package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudExpressConnectGrantRuleToCensDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 2999)
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectGrantRuleToCensDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_express_connect_grant_rule_to_cen.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectGrantRuleToCensDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_express_connect_grant_rule_to_cen.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectGrantRuleToCensDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_express_connect_grant_rule_to_cen.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectGrantRuleToCensDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_express_connect_grant_rule_to_cen.default.id}_fake"]`,
		}),
	}
	var existAlicloudExpressConnectGrantRuleToCensDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"cens.#":              "1",
			"cens.0.id":           CHECKSET,
			"cens.0.cen_id":       CHECKSET,
			"cens.0.cen_owner_id": CHECKSET,
			"cens.0.create_time":  CHECKSET,
		}
	}
	var fakeAlicloudExpressConnectGrantRuleToCensDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  "0",
			"cens.#": "0",
		}
	}
	var alicloudExpressConnectGrantRuleToCensCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_express_connect_grant_rule_to_cens.default",
		existMapFunc: existAlicloudExpressConnectGrantRuleToCensDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudExpressConnectGrantRuleToCensDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudExpressConnectGrantRuleToCensCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, allConf)
}

func testAccCheckAlicloudExpressConnectGrantRuleToCensDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccExpressConnectGrantRuleToCen-%d"
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	data "alicloud_express_connect_physical_connections" "default" {
  		name_regex = "^preserved-NODELETING"
	}

	resource "alicloud_express_connect_virtual_border_router" "default" {
  		local_gateway_ip       = "10.0.0.1"
  		peer_gateway_ip        = "10.0.0.2"
  		peering_subnet_mask    = "255.255.255.252"
  		physical_connection_id = data.alicloud_express_connect_physical_connections.default.connections.0.id
  		vlan_id                = %d
  		min_rx_interval        = 1000
  		min_tx_interval        = 1000
  		detect_multiplier      = 10
	}

	resource "alicloud_express_connect_grant_rule_to_cen" "default" {
  		cen_id       = alicloud_cen_instance.default.id
  		cen_owner_id = data.alicloud_account.default.id
  		instance_id  = alicloud_express_connect_virtual_border_router.default.id
	}

	data "alicloud_express_connect_grant_rule_to_cens" "default" {
  		instance_id = alicloud_express_connect_grant_rule_to_cen.default.instance_id
		%s
	}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}
