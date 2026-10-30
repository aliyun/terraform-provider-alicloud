package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudThreatDetectionAssetsDataSource(t *testing.T) {
	rand := tFacctest.RandInt()
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionAssetsDataSourceName(rand, map[string]string{
			"criteria":       `"[{\"name\":\"riskStatus\",\"value\":\"YES\"}]"`,
			"machine_types":  `"ecs"`,
			"no_group_trace": `"false"`,
			"importance":     `"1"`,
			"logical_exp":    `"OR"`,
		}),
		FakeConfig: "",
	}
	var existAlicloudThreatDetectionAssetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    acctest.CHECKSET,
			"assets.#": acctest.CHECKSET,
		}
	}
	var fakeAlicloudThreatDetectionAssetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"assets.#": "0",
		}
	}
	var alicloudThreatDetectionAssetsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_threat_detection_assets.default",
		ExistMapFunc: existAlicloudThreatDetectionAssetsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudThreatDetectionAssetsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudThreatDetectionAssetsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckAlicloudThreatDetectionAssetsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccThreatDetectionAsset-%d"
	}

	data "alicloud_threat_detection_assets" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
