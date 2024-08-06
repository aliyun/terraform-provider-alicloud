package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGovernanceBaselineDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGovernanceBaselineSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_governance_baseline.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGovernanceBaselineSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_governance_baseline.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGovernanceBaselineSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_governance_baseline.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGovernanceBaselineSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_governance_baseline.default.id}_fake"]`,
		}),
	}

	GovernanceBaselineCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existGovernanceBaselineMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"baselines.#":             "1",
		"baselines.0.baseline_id": CHECKSET,
	}
}

var fakeGovernanceBaselineMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"baselines.#": "0",
	}
}

var GovernanceBaselineCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_governance_baselines.default",
	existMapFunc: existGovernanceBaselineMapFunc,
	fakeMapFunc:  fakeGovernanceBaselineMapFunc,
}

func testAccCheckAlicloudGovernanceBaselineSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccGovernanceBaseline%d"
}
variable "item_password_policy" {
  default = "ACS-BP_ACCOUNT_FACTORY_RAM_USER_PASSWORD_POLICY"
}

variable "baseline_name_update" {
  default = "tf-auto-test-baseline-update"
}

variable "item_services" {
  default = "ACS-BP_ACCOUNT_FACTORY_SUBSCRIBE_SERVICES"
}

variable "baseline_name" {
  default = "tf-auto-test-baseline"
}

variable "item_ram_security" {
  default = "ACS-BP_ACCOUNT_FACTORY_RAM_SECURITY_PREFERENCE"
}



resource "alicloud_governance_baseline" "default" {
  baseline_items {
    version = "1.0"
    name    = var.item_password_policy
    config  = "{\"MinimumPasswordLength\":8,\"RequireLowercaseCharacters\":true,\"RequireUppercaseCharacters\":true,\"RequireNumbers\":true,\"RequireSymbols\":true,\"MaxPasswordAge\":0,\"HardExpiry\":false,\"PasswordReusePrevention\":0,\"MaxLoginAttempts\":0}"
  }
  description   = "tf auto test baseline"
  baseline_name = var.name
}

data "alicloud_governance_baselines" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
