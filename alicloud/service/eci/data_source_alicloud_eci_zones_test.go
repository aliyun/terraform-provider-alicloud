package eci_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEciZonesDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.EciContainerGroupRegions)
	rand := tfacctest.RandIntRange(1000000, 9999999)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEciZonesSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existEciZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#":                 acctest.CHECKSET,
			"zones.0.region_endpoint": acctest.CHECKSET,
			"zones.0.zone_ids.#":      acctest.CHECKSET,
		}
	}

	var fakeEciZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}

	var EciZonesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_eci_zones.default",
		ExistMapFunc: existEciZonesMapFunc,
		FakeMapFunc:  fakeEciZonesMapFunc,
	}

	EciZonesRecordsCheckInfo.DataSourceTestCheck(t, rand, allConf)

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
