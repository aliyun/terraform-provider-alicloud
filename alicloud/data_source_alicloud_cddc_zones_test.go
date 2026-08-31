package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCddcZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcZonesDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existAlicloudCddcZonesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             CHECKSET,
			"zones.#":           CHECKSET,
			"zones.0.region_id": os.Getenv("ALICLOUD_REGION"),
			"zones.0.zone_id":   CHECKSET,
			"zones.0.id":        CHECKSET,
		}
	}
	var fakeCddcZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudCddcZonesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cddc_zones.default",
		existMapFunc: existAlicloudCddcZonesDataSourceNameMapFunc,
		fakeMapFunc:  fakeCddcZonesMapFunc,
	}

	alicloudCddcZonesCheckInfo.dataSourceTestCheck(t, rand, allConf)
}
func testAccCheckAlicloudCddcZonesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_cddc_zones" "default" {	
  %s
}
`, strings.Join(pairs, " \n "))
	return config
}
