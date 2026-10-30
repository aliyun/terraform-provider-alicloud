package vpngateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudVPNGatewayZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVPNGatewayZonesSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existVPNGatewayZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             acctest.CHECKSET,
			"zones.#":           acctest.CHECKSET,
			"zones.0.zone_id":   acctest.CHECKSET,
			"zones.0.zone_name": acctest.CHECKSET,
		}
	}

	var fakeVPNGatewayZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var VPNGatewayZonesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_vpn_gateway_zones.default",
		ExistMapFunc: existVPNGatewayZonesMapFunc,
		FakeMapFunc:  fakeVPNGatewayZonesMapFunc,
	}

	VPNGatewayZonesRecordsCheckInfo.DataSourceTestCheck(t, rand, allConf)

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
