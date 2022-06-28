package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOOSPatchBaselinesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_patch_baseline.default.patch_baseline_name}_fake"]`,
		}),
	}
	operationSystemConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"operation_system": `"Windows"`,
		}),
		fakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"operation_system": `"AliyunLinux"`,
		}),
	}
	shareTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"share_type": `"Private"`,
		}),
		fakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"share_type": `"Public"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_patch_baseline.default.patch_baseline_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_patch_baseline.default.patch_baseline_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_oos_patch_baseline.default.patch_baseline_name}"]`,
			"name_regex":       `"${alicloud_oos_patch_baseline.default.patch_baseline_name}"`,
			"operation_system": `"Windows"`,
			"share_type":       `"Private"`,
		}),
		fakeConfig: testAccCheckAlicloudOosPatchBaselinesDataSourceName(rand, map[string]string{
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
			"baselines.0.approval_rules":      CHECKSET,
			"baselines.0.create_time":         CHECKSET,
			"baselines.0.created_by":          CHECKSET,
			"baselines.0.description":         fmt.Sprintf("tf-testAccPatchBaseline-%d", rand),
			"baselines.0.is_default":          CHECKSET,
			"baselines.0.patch_baseline_id":   CHECKSET,
			"baselines.0.updated_by":          CHECKSET,
			"baselines.0.id":                  fmt.Sprintf("tf-testAccPatchBaseline-%d", rand),
			"baselines.0.updated_date":        CHECKSET,
			"baselines.0.share_type":          "Private",
		}
	}
	var fakeAlicloudOosPatchBaselinesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudOosPatchBaselinesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_oos_patch_baselines.default",
		existMapFunc: existAlicloudOosPatchBaselinesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudOosPatchBaselinesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudOosPatchBaselinesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, operationSystemConf, shareTypeConf, nameRegexConf, allConf)
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
  approval_rules      = "{\"PatchRules\":[{\"PatchFilterGroup\":[{\"Key\":\"PatchSet\",\"Values\":[\"OS\"]},{\"Key\":\"ProductFamily\",\"Values\":[\"Windows\"]},{\"Key\":\"Product\",\"Values\":[\"Windows 10\",\"Windows 7\"]},{\"Key\":\"Classification\",\"Values\":[\"Security Updates\",\"Updates\",\"Update Rollups\",\"Critical Updates\"]},{\"Key\":\"Severity\",\"Values\":[\"Critical\",\"Important\",\"Moderate\"]}],\"ApproveAfterDays\":7,\"ApproveUntilDate\":\"\",\"EnableNonSecurity\":true,\"ComplianceLevel\":\"Medium\"}]}"
}

data "alicloud_oos_patch_baselines" "default" {
  enable_details = true
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
