package ecp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEcpZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ECPSupportRegions)
	regionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcpZonesDataSourceName(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existAlicloudEcpZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": acctest.CHECKSET,
		}
	}
	var fakeEcpZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudEcpZonesAccountBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecp_zones.default",
		ExistMapFunc: existAlicloudEcpZoneDataSourceNameMapFunc,
		FakeMapFunc:  fakeEcpZonesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	alicloudEcpZonesAccountBusesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, regionIdConf)
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
