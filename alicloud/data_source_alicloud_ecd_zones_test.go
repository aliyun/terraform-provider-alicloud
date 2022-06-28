package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcdZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	regionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdZonesDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existAlicloudEcdZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   CHECKSET,
			"zones.#": CHECKSET,
		}
	}
	var fakeEcdZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}
	var alicloudEcdZonesAccountBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_zones.default",
		existMapFunc: existAlicloudEcdZoneDataSourceNameMapFunc,
		fakeMapFunc:  fakeEcdZonesMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	alicloudEcdZonesAccountBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, regionIdConf)
}

func testAccCheckAlicloudEcdZonesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_ecd_zones" "default" {  
   %s
}
`, strings.Join(pairs, " \n "))
	return config
}
