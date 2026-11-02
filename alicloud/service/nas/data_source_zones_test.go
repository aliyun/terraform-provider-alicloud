package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudNASZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	regionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasZonesDataSourceName(rand, map[string]string{}),
		FakeConfig:  "",
	}

	fileSystemTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasZonesDataSourceName(rand, map[string]string{
			"file_system_type": `"extreme"`,
		}),
		FakeConfig: "",
	}

	var existAlicloudNasZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": acctest.CHECKSET,
		}
	}
	var fakeNasZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudNasZonesAccountBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_nas_zones.default",
		ExistMapFunc: existAlicloudNasZoneDataSourceNameMapFunc,
		FakeMapFunc:  fakeNasZonesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	alicloudNasZonesAccountBusesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, regionIdConf, fileSystemTypeConf)
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
