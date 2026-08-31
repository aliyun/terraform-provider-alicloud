package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionAssetsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionAssetsDataSourceName(rand, map[string]string{
			"criteria":       `"[{\"name\":\"riskStatus\",\"value\":\"YES\"}]"`,
			"machine_types":  `"ecs"`,
			"no_group_trace": `"false"`,
			"importance":     `"1"`,
			"logical_exp":    `"OR"`,
		}),
		fakeConfig: "",
	}
	var existAlicloudThreatDetectionAssetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    CHECKSET,
			"assets.#": CHECKSET,
		}
	}
	var fakeAlicloudThreatDetectionAssetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"assets.#": "0",
		}
	}
	var alicloudThreatDetectionAssetsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_threat_detection_assets.default",
		existMapFunc: existAlicloudThreatDetectionAssetsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudThreatDetectionAssetsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudThreatDetectionAssetsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
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
