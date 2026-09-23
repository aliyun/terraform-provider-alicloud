package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudROSRegionsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRosRegionsDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existAlicloudRosRegionDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#":                 CHECKSET,
			"regions.0.region_id":       CHECKSET,
			"regions.0.region_endpoint": CHECKSET,
			"regions.0.local_name":      CHECKSET,
		}
	}
	var fakeRosRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": "0",
		}
	}
	var alicloudRosAccountBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ros_regions.default",
		existMapFunc: existAlicloudRosRegionDataSourceNameMapFunc,
		fakeMapFunc:  fakeRosRegionsMapFunc,
	}

	alicloudRosAccountBusesCheckInfo.dataSourceTestCheck(t, rand, allConf)
}
func testAccCheckAlicloudRosRegionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_ros_regions" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}
