package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	regionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasZonesDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}

	fileSystemTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasZonesDataSourceName(rand, map[string]string{
			"file_system_type": `"extreme"`,
		}),
		fakeConfig: "",
	}

	var existAlicloudNasZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": CHECKSET,
		}
	}
	var fakeNasZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudNasZonesAccountBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nas_zones.default",
		existMapFunc: existAlicloudNasZoneDataSourceNameMapFunc,
		fakeMapFunc:  fakeNasZonesMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	alicloudNasZonesAccountBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, regionIdConf, fileSystemTypeConf)
}

func testAccCheckAlicloudNasZonesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_nas_zones" "default" {  
   %s
}
`, strings.Join(pairs, " \n "))
	return config
}
