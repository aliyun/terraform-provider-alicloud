package alb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudALBZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbZonesSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existAlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              acctest.CHECKSET,
			"zones.#":            acctest.CHECKSET,
			"zones.0.zone_id":    acctest.CHECKSET,
			"zones.0.local_name": acctest.CHECKSET,
		}
	}

	var fakeAlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var AlbZonesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_alb_zones.default",
		ExistMapFunc: existAlbZonesMapFunc,
		FakeMapFunc:  fakeAlbZonesMapFunc,
	}

	AlbZonesRecordsCheckInfo.DataSourceTestCheck(t, rand, allConf)

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
