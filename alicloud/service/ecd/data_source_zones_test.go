package ecd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEcdZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	regionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdZonesDataSourceName(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existAlicloudEcdZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   acctest.CHECKSET,
			"zones.#": acctest.CHECKSET,
		}
	}
	var fakeEcdZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}
	var alicloudEcdZonesAccountBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecd_zones.default",
		ExistMapFunc: existAlicloudEcdZoneDataSourceNameMapFunc,
		FakeMapFunc:  fakeEcdZonesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	alicloudEcdZonesAccountBusesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, regionIdConf)
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
