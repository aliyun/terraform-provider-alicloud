// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudNasLogAnalysisDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_log_analysis.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_nas_log_analysis.default.id}_fake"]`,
		}),
	}

	fileSystemTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_nas_log_analysis.default.id}"]`,
			"file_system_type": `"standard"`,
		}),
		FakeConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_nas_log_analysis.default.id}"]`,
			"file_system_type": `"extreme"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_nas_log_analysis.default.id}"]`,
			"file_system_type": `"standard"`,
		}),
		FakeConfig: testAccCheckAliCloudNasLogAnalysisSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_nas_log_analysis.default.id}_fake"]`,
			"file_system_type": `"standard"`,
		}),
	}

	AliCloudNasLogAnalysisCheckInfo.DataSourceTestCheck(t, rand, idsConf, fileSystemTypeConf, allConf)
}

var existAliCloudNasLogAnalysisMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"analyses.#":                "1",
		"analyses.0.id":             acctest.CHECKSET,
		"analyses.0.file_system_id": acctest.CHECKSET,
		"analyses.0.logstore":       acctest.CHECKSET,
		"analyses.0.project":        acctest.CHECKSET,
		"analyses.0.region":         acctest.CHECKSET,
		"analyses.0.role_arn":       acctest.CHECKSET,
	}
}

var fakeAliCloudNasLogAnalysisMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"analyses.#": "0",
	}
}

var AliCloudNasLogAnalysisCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_nas_log_analyses.default",
	ExistMapFunc: existAliCloudNasLogAnalysisMapFunc,
	FakeMapFunc:  fakeAliCloudNasLogAnalysisMapFunc,
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
`, acctest.DefaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}
