// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudNasLogAnalysisDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_log_analysis.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_log_analysis.default.id}_fake"]`,
		}),
	}

	fileSystemTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_nas_log_analysis.default.id}"]`,
			"file_system_type": `"standard"`,
		}),
		fakeConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_nas_log_analysis.default.id}"]`,
			"file_system_type": `"extreme"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_nas_log_analysis.default.id}"]`,
			"file_system_type": `"standard"`,
		}),
		fakeConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_nas_log_analysis.default.id}_fake"]`,
			"file_system_type": `"standard"`,
		}),
	}

	AliCloudNasLogAnalysisCheckInfo.dataSourceTestCheck(t, rand, idsConf, fileSystemTypeConf, allConf)
}

var existAliCloudNasLogAnalysisMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"analyses.#":                "1",
		"analyses.0.id":             CHECKSET,
		"analyses.0.file_system_id": CHECKSET,
		"analyses.0.logstore":       CHECKSET,
		"analyses.0.project":        CHECKSET,
		"analyses.0.region":         CHECKSET,
		"analyses.0.role_arn":       CHECKSET,
	}
}

var fakeAliCloudNasLogAnalysisMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"analyses.#": "0",
	}
}

var AliCloudNasLogAnalysisCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nas_log_analyses.default",
	existMapFunc: existAliCloudNasLogAnalysisMapFunc,
	fakeMapFunc:  fakeAliCloudNasLogAnalysisMapFunc,
}

func testAccCheckAliCloudNasLogAnalysisSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testacc%snasloganalysis%d"
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}

resource "alicloud_nas_log_analysis" "default" {
  file_system_id = alicloud_nas_file_system.default.id
}

data "alicloud_nas_log_analyses" "default" {
%s
}
`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}
