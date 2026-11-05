package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// NOTE: Test depends on data source or hardcoded are not stable and may fail at any time

func TestAccAliCloudThreatDetectionHoneypotImagesDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids": `["sha256:b1a9c97eb7e4b4b9030e038a0d6dfb1838401e16cf53d1ee5b3a880d3186c458"]`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids": `["metabase_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids":        `["sha256:b1a9c97eb7e4b4b9030e038a0d6dfb1838401e16cf53d1ee5b3a880d3186c458"]`,
			"name_regex": `"^counter_honeypot"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneypotImageSourceConfig(rand, map[string]string{
			"ids":        `["sha256:b1a9c97eb7e4b4b9030e038a0d6dfb1838401e16cf53d1ee5b3a880d3186c458"]`,
			"name_regex": `"^xxmeta"`,
		}),
	}

	ThreatDetectionHoneypotImageCheckInfo.DataSourceTestCheck(t, rand, idsConf, allConf)
}

var existThreatDetectionHoneypotImageMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"images.#":                             "3",
		"images.0.id":                          acctest.CHECKSET,
		"images.0.honeypot_image_display_name": "",
		"images.0.honeypot_image_id":           acctest.CHECKSET,
		"images.0.honeypot_image_name":         acctest.CHECKSET,
		"images.0.honeypot_image_type":         acctest.CHECKSET,
		"images.0.honeypot_image_version":      acctest.CHECKSET,
		"images.0.multiports":                  acctest.CHECKSET,
		"images.0.proto":                       acctest.CHECKSET,
		"images.0.template":                    acctest.CHECKSET,
	}
}

var fakeThreatDetectionHoneypotImageMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"images.#": "0",
	}
}

var ThreatDetectionHoneypotImageCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_honeypot_images.default",
	ExistMapFunc: existThreatDetectionHoneypotImageMapFunc,
	FakeMapFunc:  fakeThreatDetectionHoneypotImageMapFunc,
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
