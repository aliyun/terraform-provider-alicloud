package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudThreatDetectionVendorsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionVendorsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vendor.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionVendorsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vendor.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionVendorsSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionVendorsSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionVendorsSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_threat_detection_vendor.default.id}"]`,
			"name_regex": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionVendorsSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_threat_detection_vendor.default.id}_fake"]`,
			"name_regex": `"${var.name}_fake"`,
		}),
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Shanghai})
	}
	ThreatDetectionVendorsDataSourceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

var existThreatDetectionVendorsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vendors.#":             "1",
		"vendors.0.id":          CHECKSET,
		"vendors.0.vendor_id":   CHECKSET,
		"vendors.0.vendor_name": CHECKSET,
		"vendors.0.vendor_type": CHECKSET,
		"vendors.0.create_time": CHECKSET,
		"vendors.0.update_time": CHECKSET,
	}
}

var fakeThreatDetectionVendorsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vendors.#": "0",
	}
}

var ThreatDetectionVendorsDataSourceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_vendors.default",
	existMapFunc: existThreatDetectionVendorsMapFunc,
	fakeMapFunc:  fakeThreatDetectionVendorsMapFunc,
}

func testAccCheckAlicloudThreatDetectionVendorsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccThreatDetectionVendors%d"
}

resource "alicloud_threat_detection_vendor" "default" {
  vendor_name = var.name
  lang        = "en"
}

data "alicloud_threat_detection_vendors" "default" {
%s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
