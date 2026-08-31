package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionHoneypotImageDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids": `["sha256:02882320c9a55303410127c5dc4ae2dc470150f9d7f2483102d994f5e5f4d9df"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids": `["metabase_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids":        `["sha256:02882320c9a55303410127c5dc4ae2dc470150f9d7f2483102d994f5e5f4d9df"]`,
			"name_regex": `"^meta"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids":        `["sha256:02882320c9a55303410127c5dc4ae2dc470150f9d7f2483102d994f5e5f4d9df"]`,
			"name_regex": `"^xxmeta"`,
		}),
	}

	ThreatDetectionHoneypotImageCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existThreatDetectionHoneypotImageMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"images.#":                             "1",
		"images.0.id":                          CHECKSET,
		"images.0.honeypot_image_display_name": CHECKSET,
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
