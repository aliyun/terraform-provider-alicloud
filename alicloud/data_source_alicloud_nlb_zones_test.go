package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNLBZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNlbZonesSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existNlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              CHECKSET,
			"zones.#":            CHECKSET,
			"zones.0.zone_id":    CHECKSET,
			"zones.0.local_name": CHECKSET,
		}
	}

	var fakeNlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var NlbZonesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nlb_zones.default",
		existMapFunc: existNlbZonesMapFunc,
		fakeMapFunc:  fakeNlbZonesMapFunc,
	}

	NlbZonesRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

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
