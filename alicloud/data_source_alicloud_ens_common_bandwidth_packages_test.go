package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsCommonBandwidthPackagesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idFilterConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"bandwidth_package_id": `"${alicloud_ens_common_bandwidth_package.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"bandwidth_package_id": `"${alicloud_ens_common_bandwidth_package.default.id}_fake"`,
		}),
	}

	nameFilterConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"name": `"${alicloud_ens_common_bandwidth_package.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"name": `"${alicloud_ens_common_bandwidth_package.default.name}_fake"`,
		}),
	}

	EnsCommonBandwidthPackageCheckInfo.dataSourceTestCheck(t, rand, idFilterConf, nameFilterConf)
}

func testAccCheckAlicloudEnsCommonBandwidthPackagesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
    default = "tf-testAccEnsCommonBandwidthPackagesDataSource%d"
}

resource "alicloud_ens_common_bandwidth_package" "default" {
  bandwidth     = 10
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  description   = var.name
  name          = var.name
}

data "alicloud_ens_common_bandwidth_packages" "default" {
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
%s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existEnsCommonBandwidthPackageMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"packages.#":               "1",
		"packages.0.bandwidth":     "10",
		"packages.0.ens_region_id": "cn-chenzhou-telecom_unicom_cmcc",
		"packages.0.name":          fmt.Sprintf("tf-testAccEnsCommonBandwidthPackagesDataSource%d", rand),
	}
}

var fakeEnsCommonBandwidthPackageMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"packages.#": "0",
	}
}

var EnsCommonBandwidthPackageCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_common_bandwidth_packages.default",
	existMapFunc: existEnsCommonBandwidthPackageMapFunc,
	fakeMapFunc:  fakeEnsCommonBandwidthPackageMapFunc,
}
