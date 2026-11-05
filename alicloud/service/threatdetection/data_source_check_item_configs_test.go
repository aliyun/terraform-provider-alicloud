package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// NOTE: Test depends on data source or hardcoded are not stable and may fail at any time
// Uninitialized resource; assumes the resource already exists.

func TestAccAliCloudThreatDetectionCheckItemConfigsDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	pagingConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionCheckItemConfigSourceConfig(rand, map[string]string{
			"lang":        `"zh"`,
			"page_number": `1`,
			"page_size":   `10`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionCheckItemConfigSourceConfig(rand, map[string]string{
			"lang":        `"zh"`,
			"page_number": `1`,
			"page_size":   `10`,
			"ids":         `["fake-id-not-exist"]`,
		}),
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}
	ThreatDetectionCheckItemConfigCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, pagingConf)
}

var existThreatDetectionCheckItemConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"configs.#":                   acctest.CHECKSET,
		"configs.0.section_ids.#":     acctest.CHECKSET,
		"configs.0.description.#":     acctest.CHECKSET,
		"configs.0.check_show_name":   acctest.CHECKSET,
		"configs.0.vendor":            acctest.CHECKSET,
		"configs.0.instance_sub_type": acctest.CHECKSET,
		"configs.0.check_id":          acctest.CHECKSET,
		"configs.0.custom_configs.#":  acctest.CHECKSET,
		"configs.0.check_type":        acctest.CHECKSET,
		"configs.0.estimated_count":   acctest.CHECKSET,
		"configs.0.instance_type":     acctest.CHECKSET,
		"configs.0.risk_level":        acctest.CHECKSET,
	}
}

var fakeThreatDetectionCheckItemConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"configs.#": "0",
	}
}

var ThreatDetectionCheckItemConfigCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_check_item_configs.default",
	ExistMapFunc: existThreatDetectionCheckItemConfigMapFunc,
	FakeMapFunc:  fakeThreatDetectionCheckItemConfigMapFunc,
}

func testAccCheckAlicloudThreatDetectionCheckItemConfigSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccThreatDetectionCheckItemConfig%d"
}

data "alicloud_threat_detection_check_item_configs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
