package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

// NOTE: Test depends on data source or hardcoded are not stable and may fail at any time

func TestAccAliCloudThreatDetectionHoneypotImagesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids": `["sha256:b1a9c97eb7e4b4b9030e038a0d6dfb1838401e16cf53d1ee5b3a880d3186c458"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids": `["metabase_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids":        `["sha256:b1a9c97eb7e4b4b9030e038a0d6dfb1838401e16cf53d1ee5b3a880d3186c458"]`,
			"name_regex": `"^counter_honeypot"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids":        `["sha256:b1a9c97eb7e4b4b9030e038a0d6dfb1838401e16cf53d1ee5b3a880d3186c458"]`,
			"name_regex": `"^xxmeta"`,
		}),
	}

	ThreatDetectionHoneypotImageCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existThreatDetectionHoneypotImageMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"images.#":                             "3",
		"images.0.id":                          CHECKSET,
		"images.0.honeypot_image_display_name": "",
		"images.0.honeypot_image_id":           CHECKSET,
		"images.0.honeypot_image_name":         CHECKSET,
		"images.0.honeypot_image_type":         CHECKSET,
		"images.0.honeypot_image_version":      CHECKSET,
		"images.0.multiports":                  CHECKSET,
		"images.0.proto":                       CHECKSET,
		"images.0.template":                    CHECKSET,
	}
}

var fakeThreatDetectionHoneypotImageMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"images.#": "0",
	}
}

var ThreatDetectionHoneypotImageCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_honeypot_images.default",
	existMapFunc: existThreatDetectionHoneypotImageMapFunc,
	fakeMapFunc:  fakeThreatDetectionHoneypotImageMapFunc,
}

func testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionHoneypotImage%d"
}

data "alicloud_threat_detection_honeypot_images" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
