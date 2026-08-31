package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudALBZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbZonesSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existAlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              CHECKSET,
			"zones.#":            CHECKSET,
			"zones.0.zone_id":    CHECKSET,
			"zones.0.local_name": CHECKSET,
		}
	}

	var fakeAlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var AlbZonesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_zones.default",
		existMapFunc: existAlbZonesMapFunc,
		fakeMapFunc:  fakeAlbZonesMapFunc,
	}

	AlbZonesRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudAlbZonesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_alb_zones" "default"{
%s
}

`, strings.Join(pairs, "\n   "))
	return config
}
