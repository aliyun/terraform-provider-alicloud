package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCddcHostEcsLevelInfosDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcHostEcsLevelInfosDataSourceName(rand, map[string]string{
			"db_type":        `"mysql"`,
			"zone_id":        `"${data.alicloud_cddc_zones.default.ids.0}"`,
			"storage_type":   `"cloud_essd"`,
			"image_category": `"AliLinux"`,
		}),
		fakeConfig: "",
	}

	var existAlicloudCddcHostEcsLevelInfosDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"infos.#":                CHECKSET,
			"infos.0.res_class_code": CHECKSET,
			"infos.0.ecs_class_code": CHECKSET,
			"infos.0.ecs_class":      CHECKSET,
			"infos.0.description":    CHECKSET,
		}
	}
	var fakeCddcHostEcsLevelInfosMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"infos.#": "0",
		}
	}
	var alicloudCddcHostEcsLevelInfoCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cddc_host_ecs_level_infos.default",
		existMapFunc: existAlicloudCddcHostEcsLevelInfosDataSourceNameMapFunc,
		fakeMapFunc:  fakeCddcHostEcsLevelInfosMapFunc,
	}

	alicloudCddcHostEcsLevelInfoCheckInfo.dataSourceTestCheck(t, rand, allConf)
}
func testAccCheckAlicloudCddcHostEcsLevelInfosDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {	
  %s
}`, strings.Join(pairs, " \n "))
	return config
}
