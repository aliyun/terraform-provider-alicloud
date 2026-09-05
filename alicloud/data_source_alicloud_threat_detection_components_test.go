package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionComponentsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionComponentsSourceConfig(rand, map[string]string{
			"name_regex": `".*"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionComponentsSourceConfig(rand, map[string]string{
			"ids": `["tf-acc-fake-component-name-not-exist-%d"]`,
		}),
	}

	componentNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionComponentsSourceConfig(rand, map[string]string{
			"component_name": `"tf-acc-fake-component-name-not-exist-%d"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionComponentsSourceConfig(rand, map[string]string{
			"component_name": `"tf-acc-fake-component-name-not-exist-%d"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionComponentsSourceConfig(rand, map[string]string{
			"name_regex":  `".*"`,
			"output_file": `"%s"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionComponentsSourceConfig(rand, map[string]string{
			"name_regex": `"^tf-acc-fake-component-name-not-exist"`,
		}),
	}

	ThreatDetectionComponentsCheckInfo.dataSourceTestCheck(t, rand, idsConf, componentNameConf, allConf)
}

var existThreatDetectionComponentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"components.#":                       CHECKSET,
		"components.0.id":                    CHECKSET,
		"components.0.component_name":        CHECKSET,
		"components.0.component_alias":       CHECKSET,
		"components.0.component_description": CHECKSET,
	}
}

var fakeThreatDetectionComponentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"components.#": "0",
	}
}

var ThreatDetectionComponentsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_components.default",
	existMapFunc: existThreatDetectionComponentsMapFunc,
	fakeMapFunc:  fakeThreatDetectionComponentsMapFunc,
}

func testAccCheckAlicloudThreatDetectionComponentsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+fmt.Sprintf(v, rand))
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionComponents%d"
}

data "alicloud_threat_detection_components" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
