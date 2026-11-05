package governance_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudGovernanceBaselineDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGovernanceBaselineSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_governance_baseline.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGovernanceBaselineSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_governance_baseline.default.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGovernanceBaselineSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_governance_baseline.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGovernanceBaselineSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_governance_baseline.default.id}_fake"]`,
		}),
	}

	GovernanceBaselineCheckInfo.DataSourceTestCheck(t, rand, idsConf, allConf)
}

var existGovernanceBaselineMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"baselines.#":             "1",
		"baselines.0.baseline_id": acctest.CHECKSET,
	}
}

var fakeGovernanceBaselineMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"baselines.#": "0",
	}
}

var GovernanceBaselineCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_governance_baselines.default",
	ExistMapFunc: existGovernanceBaselineMapFunc,
	FakeMapFunc:  fakeGovernanceBaselineMapFunc,
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
