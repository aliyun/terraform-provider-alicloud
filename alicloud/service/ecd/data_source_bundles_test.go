package ecd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudECDBundlesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	bundleTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdBundlesDataSourceName(rand, map[string]string{
			"bundle_type": `"SYSTEM"`,
		}),
		FakeConfig: "",
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEcdBundlesDataSourceName(rand, map[string]string{
			"name_regex": `"General"`,
		}),
		FakeConfig: testAccCheckAlicloudEcdBundlesDataSourceName(rand, map[string]string{
			"name_regex": `"General_fake"`,
		}),
	}

	var existAlicloudEcdBundlesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              acctest.CHECKSET,
			"names.#":                            acctest.CHECKSET,
			"bundles.#":                          acctest.CHECKSET,
			"bundles.0.id":                       acctest.CHECKSET,
			"bundles.0.bundle_id":                acctest.CHECKSET,
			"bundles.0.bundle_name":              acctest.CHECKSET,
			"bundles.0.bundle_type":              "SYSTEM",
			"bundles.0.description":              "",
			"bundles.0.desktop_type":             acctest.CHECKSET,
			"bundles.0.desktop_type_attribute.#": acctest.CHECKSET,
			"bundles.0.desktop_type_attribute.0.cpu_count": acctest.CHECKSET,
			"bundles.0.desktop_type_attribute.0.gpu_count": acctest.CHECKSET,
			//todo : The field does not necessarily have a value, so the note
			//"bundles.0.desktop_type_attribute.0.gpu_spec":    acctest.CHECKSET,
			"bundles.0.desktop_type_attribute.0.memory_size": acctest.CHECKSET,
			"bundles.0.disks.#":           acctest.CHECKSET,
			"bundles.0.disks.0.disk_size": acctest.CHECKSET,
			"bundles.0.disks.0.disk_type": acctest.CHECKSET,
			"bundles.0.image_id":          acctest.CHECKSET,
			"bundles.0.os_type":           acctest.CHECKSET,
		}
	}
	var fakeAlicloudEcdBundlesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"bundles.#": "0",
		}
	}
	var alicloudEcdBundlesBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ecd_bundles.default",
		ExistMapFunc: existAlicloudEcdBundlesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEcdBundlesDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
	}
	alicloudEcdBundlesBusesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, bundleTypeConf, nameRegexConf)
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
