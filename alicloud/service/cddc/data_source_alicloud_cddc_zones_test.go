package cddc_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCddcZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCddcZonesDataSourceName(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existAlicloudCddcZonesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             acctest.CHECKSET,
			"zones.#":           acctest.CHECKSET,
			"zones.0.region_id": os.Getenv("ALICLOUD_REGION"),
			"zones.0.zone_id":   acctest.CHECKSET,
			"zones.0.id":        acctest.CHECKSET,
		}
	}
	var fakeCddcZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudCddcZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cddc_zones.default",
		ExistMapFunc: existAlicloudCddcZonesDataSourceNameMapFunc,
		FakeMapFunc:  fakeCddcZonesMapFunc,
	}

	alicloudCddcZonesCheckInfo.DataSourceTestCheck(t, rand, allConf)
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
