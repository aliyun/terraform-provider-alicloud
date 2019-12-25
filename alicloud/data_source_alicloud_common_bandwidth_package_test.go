package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudCommonBandwidthPackagesDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `alicloud_common_bandwidth_package.default.name`,
		}),
		fakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${alicloud_common_bandwidth_package.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ alicloud_common_bandwidth_package.default.id ]`,
		}),
		fakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ "${alicloud_common_bandwidth_package.default.id}_fake" ]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ alicloud_common_bandwidth_package.default.id ]`,
			"name_regex": `alicloud_common_bandwidth_package.default.name`,
		}),
		fakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ "${alicloud_common_bandwidth_package.default.id}_fake" ]`,
			"name_regex": `"${alicloud_common_bandwidth_package.default.name}_fake"`,
		}),
	}
	commonBandwidthPackagesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCommonBandwidthPackageDataSource%d"
}

resource "alicloud_common_bandwidth_package" "default" {
  bandwidth = "2"
  name = var.name
  description = "${var.name}_description"
  resource_group_id = "%s"
}

data "alicloud_common_bandwidth_packages" "default"  {
  %s
}
`, rand, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"), strings.Join(pairs, "\n  "))
	return config
}

var existsCommonBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                      "1",
		"names.#":                    "1",
		"packages.#":                 "1",
		"packages.0.id":              CHECKSET,
		"packages.0.isp":             CHECKSET,
		"packages.0.creation_time":   CHECKSET,
		"packages.0.status":          CHECKSET,
		"packages.0.business_status": CHECKSET,
		"packages.0.bandwidth":       "2",
	}
}

var fakeCommonBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"packages.#": "0",
	}
}

var commonBandwidthPackagesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_common_bandwidth_packages.default",
	existMapFunc: existsCommonBandwidthPackagesMapFunc,
	fakeMapFunc:  fakeCommonBandwidthPackagesMapFunc,
}
