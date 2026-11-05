package ros_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudROSRegionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRosRegionsDataSourceName(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existAlicloudRosRegionDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#":                 acctest.CHECKSET,
			"regions.0.region_id":       acctest.CHECKSET,
			"regions.0.region_endpoint": acctest.CHECKSET,
			"regions.0.local_name":      acctest.CHECKSET,
		}
	}
	var fakeRosRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": "0",
		}
	}
	var alicloudRosAccountBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ros_regions.default",
		ExistMapFunc: existAlicloudRosRegionDataSourceNameMapFunc,
		FakeMapFunc:  fakeRosRegionsMapFunc,
	}

	alicloudRosAccountBusesCheckInfo.DataSourceTestCheck(t, rand, allConf)
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
