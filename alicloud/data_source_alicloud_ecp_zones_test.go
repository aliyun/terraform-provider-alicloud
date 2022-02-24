package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcpZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	checkoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	regionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcpZonesDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existAlicloudEcpZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": CHECKSET,
		}
	}
	var fakeEcpZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudEcpZonesAccountBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecp_zones.default",
		existMapFunc: existAlicloudEcpZoneDataSourceNameMapFunc,
		fakeMapFunc:  fakeEcpZonesMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	alicloudEcpZonesAccountBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, regionIdConf)
}

func testAccCheckAlicloudEcpZonesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_ecp_zones" "default" {  
   %s
}
`, strings.Join(pairs, " \n "))
	return config
}
