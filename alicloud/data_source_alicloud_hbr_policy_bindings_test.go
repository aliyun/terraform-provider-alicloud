package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBRPolicyBindingsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacchbrpb%d", rand)

	bindingIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"ids": `["${alicloud_hbr_policy_binding.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"ids": `["${alicloud_hbr_policy_binding.default.id}_fake"]`,
		}),
	}

	policyIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"policy_id": `"${alicloud_hbr_policy_binding.default.policy_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"policy_id": `"${alicloud_hbr_policy_binding.default.policy_id}_fake"`,
		}),
	}

	sourceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"policy_id":   `"${alicloud_hbr_policy_binding.default.policy_id}"`,
			"source_type": `"OSS"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"policy_id":   `"${alicloud_hbr_policy_binding.default.policy_id}"`,
			"source_type": `"NAS"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"policy_id":  `"${alicloud_hbr_policy_binding.default.policy_id}"`,
			"name_regex": `"policy binding example"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"policy_id":  `"${alicloud_hbr_policy_binding.default.policy_id}"`,
			"name_regex": `"policy binding example_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"ids":         `["${alicloud_hbr_policy_binding.default.id}"]`,
			"policy_id":   `"${alicloud_hbr_policy_binding.default.policy_id}"`,
			"source_type": `"OSS"`,
			"name_regex":  `"policy binding example"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name, rand, map[string]string{
			"ids":         `["${alicloud_hbr_policy_binding.default.id}_fake"]`,
			"policy_id":   `"${alicloud_hbr_policy_binding.default.policy_id}"`,
			"source_type": `"OSS"`,
			"name_regex":  `"policy binding example_fake"`,
		}),
	}

	HbrPolicyBindingCheckInfo.dataSourceTestCheck(t, rand, bindingIdConf, policyIdConf, sourceTypeConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudHbrPolicyBindingsSourceConfig(name string, rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_hbr_vault" "defaultyk84Hc" {
  vault_type = "STANDARD"
  vault_name = var.name
}

resource "alicloud_hbr_policy" "defaultoqWvHQ" {
  policy_name = var.name
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|0|P1D"
    retention    = "7"
    vault_id     = alicloud_hbr_vault.defaultyk84Hc.id
    archive_days = "0"
  }
}

resource "alicloud_oss_bucket" "defaultKtt2XY" {
  storage_class = "Standard"
}

resource "alicloud_hbr_policy_binding" "default" {
  source_type                = "OSS"
  disabled                   = "false"
  policy_id                  = alicloud_hbr_policy.defaultoqWvHQ.id
  data_source_id             = alicloud_oss_bucket.defaultKtt2XY.id
  policy_binding_description = "policy binding example"
  source                     = "prefix-example-create/"
}

data "alicloud_hbr_policy_bindings" "default" {
%s
}
`, name, strings.Join(pairs, "\n   "))
	return config
}

var existHbrPolicyBindingMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policy_bindings.#":                "1",
		"policy_bindings.0.id":             CHECKSET,
		"policy_bindings.0.policy_id":      CHECKSET,
		"policy_bindings.0.source_type":    "OSS",
		"policy_bindings.0.data_source_id": CHECKSET,
		"policy_bindings.0.disabled":       "false",
		"policy_bindings.0.create_time":    CHECKSET,
	}
}

var fakeHbrPolicyBindingMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policy_bindings.#": "0",
	}
}

var HbrPolicyBindingCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_policy_bindings.default",
	existMapFunc: existHbrPolicyBindingMapFunc,
	fakeMapFunc:  fakeHbrPolicyBindingMapFunc,
}
