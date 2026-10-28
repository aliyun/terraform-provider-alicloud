package cddc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCddcHostEcsLevelInfosDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCddcHostEcsLevelInfosDataSourceName(rand, map[string]string{
			"db_type":        `"mysql"`,
			"zone_id":        `"${data.alicloud_cddc_zones.default.ids.0}"`,
			"storage_type":   `"cloud_essd"`,
			"image_category": `"AliLinux"`,
		}),
		FakeConfig: "",
	}

	var existAlicloudCddcHostEcsLevelInfosDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"infos.#":                acctest.CHECKSET,
			"infos.0.res_class_code": acctest.CHECKSET,
			"infos.0.ecs_class_code": acctest.CHECKSET,
			"infos.0.ecs_class":      acctest.CHECKSET,
			"infos.0.description":    acctest.CHECKSET,
		}
	}
	var fakeCddcHostEcsLevelInfosMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"infos.#": "0",
		}
	}
	var alicloudCddcHostEcsLevelInfoCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cddc_host_ecs_level_infos.default",
		ExistMapFunc: existAlicloudCddcHostEcsLevelInfosDataSourceNameMapFunc,
		FakeMapFunc:  fakeCddcHostEcsLevelInfosMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.CDDCSupportRegions)
	}

	alicloudCddcHostEcsLevelInfoCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
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
