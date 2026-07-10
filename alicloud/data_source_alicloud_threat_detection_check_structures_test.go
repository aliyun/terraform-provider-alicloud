package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudThreatDetectionCheckStructuresDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionCheckStructureSourceConfig(rand, map[string]string{
			"lang":         `"zh"`,
			"current_page": `1`,
			"page_size":    `10`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionCheckStructureSourceConfig(rand, map[string]string{
			"lang":         `"zh"`,
			"current_page": `1`,
			"page_size":    `10`,
			"ids":          `["fake-id-not-exist"]`,
		}),
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}
	ThreatDetectionCheckStructureCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, pagingConf)
}

var existThreatDetectionCheckStructureMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"structures.#":               CHECKSET,
		"structures.0.standards.#":   CHECKSET,
		"structures.0.standard_type": CHECKSET,
	}
}

var fakeThreatDetectionCheckStructureMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"structures.#": "0",
	}
}

var ThreatDetectionCheckStructureCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_check_structures.default",
	existMapFunc: existThreatDetectionCheckStructureMapFunc,
	fakeMapFunc:  fakeThreatDetectionCheckStructureMapFunc,
}

func testAccCheckAlicloudThreatDetectionCheckStructureSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccThreatDetectionCheckStructure%d"
}

data "alicloud_threat_detection_check_structures" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
