package nlb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudNLBZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNlbZonesSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existNlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              acctest.CHECKSET,
			"zones.#":            acctest.CHECKSET,
			"zones.0.zone_id":    acctest.CHECKSET,
			"zones.0.local_name": acctest.CHECKSET,
		}
	}

	var fakeNlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var NlbZonesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_nlb_zones.default",
		ExistMapFunc: existNlbZonesMapFunc,
		FakeMapFunc:  fakeNlbZonesMapFunc,
	}

	NlbZonesRecordsCheckInfo.DataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudNlbZonesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_nlb_zones" "default"{
%s
}

`, strings.Join(pairs, "\n   "))
	return config
}
