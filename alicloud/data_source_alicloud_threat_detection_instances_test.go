package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionInstanceDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionInstanceSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	ThreatDetectionInstanceCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existThreatDetectionInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                    "1",
		"instances.#":              "1",
		"instances.0.id":           CHECKSET,
		"instances.0.create_time":  CHECKSET,
		"instances.0.instance_id":  CHECKSET,
		"instances.0.payment_type": CHECKSET,
		"instances.0.status":       CHECKSET,
	}
}

var fakeThreatDetectionInstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"instances.#": "0",
	}
}

var ThreatDetectionInstanceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_instances.default",
	existMapFunc: existThreatDetectionInstanceMapFunc,
	fakeMapFunc:  fakeThreatDetectionInstanceMapFunc,
}

func testAccCheckAlicloudThreatDetectionInstanceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionInstance%d"
}

data "alicloud_threat_detection_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
