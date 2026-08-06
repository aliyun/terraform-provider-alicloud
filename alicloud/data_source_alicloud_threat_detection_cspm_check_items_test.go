package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// NOTE: Test depends on data source or hardcoded are not stable and may fail at any time
// Uninitialized resource; assumes the resource already exists.

func TestAccAliCloudThreatDetectionCspmCheckItemsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionCspmCheckItemSourceConfig(rand, map[string]string{
			"output_file": `"./tf-testacc-threat-detection-cspm-check-items.txt"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionCspmCheckItemSourceConfig(rand, map[string]string{
			"ids":          `["fake-id-not-exist"]`,
			"task_sources": `["fake-task-source-not-exist"]`,
		}),
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}
	ThreatDetectionCspmCheckItemCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, pagingConf)
}

var existThreatDetectionCspmCheckItemMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"items.#":                   CHECKSET,
		"items.0.id":                CHECKSET,
		"items.0.check_id":          CHECKSET,
		"items.0.check_show_name":   CHECKSET,
		"items.0.check_type":        CHECKSET,
		"items.0.custom_configs.#":  CHECKSET,
		"items.0.description.#":     CHECKSET,
		"items.0.estimated_count":   CHECKSET,
		"items.0.instance_sub_type": CHECKSET,
		"items.0.instance_type":     CHECKSET,
		"items.0.risk_level":        CHECKSET,
		"items.0.section_ids.#":     CHECKSET,
		"items.0.vendor":            CHECKSET,
	}
}

var fakeThreatDetectionCspmCheckItemMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"items.#": "0",
	}
}

var ThreatDetectionCspmCheckItemCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_cspm_check_items.default",
	existMapFunc: existThreatDetectionCspmCheckItemMapFunc,
	fakeMapFunc:  fakeThreatDetectionCspmCheckItemMapFunc,
}

func testAccCheckAlicloudThreatDetectionCspmCheckItemSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccThreatDetectionCspmCheckItem%d"
}

data "alicloud_threat_detection_cspm_check_items" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
