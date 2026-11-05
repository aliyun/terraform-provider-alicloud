package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudThreatDetectionCheckStructuresDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	pagingConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionCheckStructureSourceConfig(rand, map[string]string{
			"lang":         `"zh"`,
			"current_page": `1`,
			"page_size":    `10`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionCheckStructureSourceConfig(rand, map[string]string{
			"lang":         `"zh"`,
			"current_page": `1`,
			"page_size":    `10`,
			"ids":          `["fake-id-not-exist"]`,
		}),
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}
	ThreatDetectionCheckStructureCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, pagingConf)
}

var existThreatDetectionCheckStructureMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"structures.#":               acctest.CHECKSET,
		"structures.0.standards.#":   acctest.CHECKSET,
		"structures.0.standard_type": acctest.CHECKSET,
	}
}

var fakeThreatDetectionCheckStructureMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"structures.#": "0",
	}
}

var ThreatDetectionCheckStructureCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_check_structures.default",
	ExistMapFunc: existThreatDetectionCheckStructureMapFunc,
	FakeMapFunc:  fakeThreatDetectionCheckStructureMapFunc,
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
