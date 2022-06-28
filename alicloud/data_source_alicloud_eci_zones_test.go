package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEciZonesDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.EciContainerGroupRegions)
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEciZonesSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existEciZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#":                 CHECKSET,
			"zones.0.region_endpoint": CHECKSET,
			"zones.0.zone_ids.#":      CHECKSET,
		}
	}

	var fakeEciZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}

	var EciZonesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_eci_zones.default",
		existMapFunc: existEciZonesMapFunc,
		fakeMapFunc:  fakeEciZonesMapFunc,
	}

	EciZonesRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudEciZonesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_eci_zones" "default"{
%s
}

`, strings.Join(pairs, "\n   "))
	return config
}
