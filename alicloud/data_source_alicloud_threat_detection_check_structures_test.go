// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionCheckStructureDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	ThreatDetectionCheckStructureCheckInfo.dataSourceTestCheck(t, rand)
}

var existThreatDetectionCheckStructureMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"structures.#":               "1",
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
