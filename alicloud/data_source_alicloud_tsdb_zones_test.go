package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudTsdbZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudTsdbZonesSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existTsdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              CHECKSET,
			"zones.#":            CHECKSET,
			"zones.0.zone_id":    CHECKSET,
			"zones.0.local_name": CHECKSET,
		}
	}

	var fakeTsdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var tsdbZonesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_tsdb_zones.default",
		existMapFunc: existTsdbZonesMapFunc,
		fakeMapFunc:  fakeTsdbZonesMapFunc,
	}

	tsdbZonesRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudTsdbZonesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_tsdb_zones" "default"{
%s
}

`, strings.Join(pairs, "\n   "))
	return config
}
