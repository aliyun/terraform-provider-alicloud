package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPNGatewayZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVPNGatewayZonesSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existVPNGatewayZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             CHECKSET,
			"zones.#":           CHECKSET,
			"zones.0.zone_id":   CHECKSET,
			"zones.0.zone_name": CHECKSET,
		}
	}

	var fakeVPNGatewayZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var VPNGatewayZonesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpn_gateway_zones.default",
		existMapFunc: existVPNGatewayZonesMapFunc,
		fakeMapFunc:  fakeVPNGatewayZonesMapFunc,
	}

	VPNGatewayZonesRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudVPNGatewayZonesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_vpn_gateway_zones" "default"{
	spec = "5M"
%s
}

`, strings.Join(pairs, "\n   "))
	return config
}
