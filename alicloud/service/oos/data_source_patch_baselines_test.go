package oos_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudOOSPatchBaselinesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
		}),
		FakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_patch_baseline.default.patch_baseline_name}_fake"]`,
		}),
	}
	operationSystemConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"operation_system": `"Windows"`,
		}),
		FakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"operation_system": `"AliyunLinux"`,
		}),
	}
	shareTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"share_type": `"Private"`,
		}),
		FakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"share_type": `"Public"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_patch_baseline.default.patch_baseline_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_patch_baseline.default.patch_baseline_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"name_regex":       `"${alicloud_oos_patch_baseline.default.patch_baseline_name}"`,
			"operation_system": `"Windows"`,
			"share_type":       `"Private"`,
		}),
		FakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_oos_patch_baseline.default.patch_baseline_name}_fake"]`,
			"name_regex":       `"${alicloud_oos_patch_baseline.default.patch_baseline_name}_fake"`,
			"operation_system": `"AliyunLinux"`,
			"share_type":       `"Public"`,
		}),
	}
	var existAlicloudOosPatchBaselinesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"baselines.#":                     "1",
			"baselines.0.operation_system":    "Windows",
			"baselines.0.patch_baseline_name": fmt.Sprintf("tf-testAccPatchBaseline-%d", rand),
			"baselines.0.approval_rules":      acctest.CHECKSET,
			"baselines.0.create_time":         acctest.CHECKSET,
			"baselines.0.created_by":          acctest.CHECKSET,
			"baselines.0.description":         fmt.Sprintf("tf-testAccPatchBaseline-%d", rand),
			"baselines.0.is_default":          acctest.CHECKSET,
			"baselines.0.patch_baseline_id":   acctest.CHECKSET,
			"baselines.0.updated_by":          acctest.CHECKSET,
			"baselines.0.id":                  fmt.Sprintf("tf-testAccPatchBaseline-%d", rand),
			"baselines.0.updated_date":        acctest.CHECKSET,
			"baselines.0.share_type":          "Private",
		}
	}
	var fakeAlicloudOosPatchBaselinesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudOosPatchBaselinesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_oos_patch_baselines.default",
		ExistMapFunc: existAlicloudOosPatchBaselinesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudOosPatchBaselinesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudOosPatchBaselinesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, operationSystemConf, shareTypeConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testAccPatchBaseline-%d"
}

resource "alicloud_oos_patch_baseline" "default" {
  operation_system    = "Windows"
  patch_baseline_name = var.name
  description         = var.name
  approval_rules      = "{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}"
}

data "alicloud_oos_patch_baselines" "default" {
  enable_details = true
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
