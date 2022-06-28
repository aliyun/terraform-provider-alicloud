package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECDBundlesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	bundleTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdBundlesDataSourceName(rand, map[string]string{
			"bundle_type": `"SYSTEM"`,
		}),
		fakeConfig: "",
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdBundlesDataSourceName(rand, map[string]string{
			"name_regex": `"General"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdBundlesDataSourceName(rand, map[string]string{
			"name_regex": `"General_fake"`,
		}),
	}

	var existAlicloudEcdBundlesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              CHECKSET,
			"names.#":                            CHECKSET,
			"bundles.#":                          CHECKSET,
			"bundles.0.id":                       CHECKSET,
			"bundles.0.bundle_id":                CHECKSET,
			"bundles.0.bundle_name":              CHECKSET,
			"bundles.0.bundle_type":              "SYSTEM",
			"bundles.0.description":              "",
			"bundles.0.desktop_type":             CHECKSET,
			"bundles.0.desktop_type_attribute.#": CHECKSET,
			"bundles.0.desktop_type_attribute.0.cpu_count": CHECKSET,
			"bundles.0.desktop_type_attribute.0.gpu_count": CHECKSET,
			//todo : The field does not necessarily have a value, so the note
			//"bundles.0.desktop_type_attribute.0.gpu_spec":    CHECKSET,
			"bundles.0.desktop_type_attribute.0.memory_size": CHECKSET,
			"bundles.0.disks.#":           CHECKSET,
			"bundles.0.disks.0.disk_size": CHECKSET,
			"bundles.0.disks.0.disk_type": CHECKSET,
			"bundles.0.image_id":          CHECKSET,
			"bundles.0.os_type":           CHECKSET,
		}
	}
	var fakeAlicloudEcdBundlesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"bundles.#": "0",
		}
	}
	var alicloudEcdBundlesBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_bundles.default",
		existMapFunc: existAlicloudEcdBundlesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdBundlesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
	}
	alicloudEcdBundlesBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, bundleTypeConf, nameRegexConf)
}
func testAccCheckAlicloudEcdBundlesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testaccdesktop%d"
}

data "alicloud_ecd_bundles" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
