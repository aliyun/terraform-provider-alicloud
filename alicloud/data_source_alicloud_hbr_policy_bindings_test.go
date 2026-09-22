package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHbrPolicyBindingsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(rand, map[string]string{
			"ids": `[alicloud_hbr_policy_binding.default.id]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_policy_binding.default.id}_fake"]`,
		}),
	}

	policyIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(rand, map[string]string{
			"policy_id": `alicloud_hbr_policy.default.id`,
		}),
	}

	sourceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(rand, map[string]string{
			"policy_id":   `alicloud_hbr_policy.default.id`,
			"source_type": `"OSS"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(rand, map[string]string{
			"policy_id":   `alicloud_hbr_policy.default.id`,
			"source_type": `"NAS"`,
		}),
	}

	dataSourceIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(rand, map[string]string{
			"policy_id":       `alicloud_hbr_policy.default.id`,
			"data_source_ids": `[alicloud_oss_bucket.default.id]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrPolicyBindingsSourceConfig(rand, map[string]string{
			"ids":             `[alicloud_hbr_policy_binding.default.id]`,
			"policy_id":       `alicloud_hbr_policy.default.id`,
			"source_type":     `"OSS"`,
			"data_source_ids": `[alicloud_oss_bucket.default.id]`,
		}),
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	HbrPolicyBindingsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, policyIdConf, sourceTypeConf, dataSourceIdsConf, allConf)
}

var existHbrPolicyBindingsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                 "1",
		"bindings.#":                            "1",
		"bindings.0.id":                         CHECKSET,
		"bindings.0.policy_id":                  CHECKSET,
		"bindings.0.source_type":                "OSS",
		"bindings.0.data_source_id":             CHECKSET,
		"bindings.0.disabled":                   "false",
		"bindings.0.source":                     "prefix-example-create/",
		"bindings.0.policy_binding_description": "policy binding example",
		"bindings.0.create_time":                CHECKSET,
	}
}

var fakeHbrPolicyBindingsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"bindings.#": "0",
	}
}

var HbrPolicyBindingsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_policy_bindings.default",
	existMapFunc: existHbrPolicyBindingsMapFunc,
	fakeMapFunc:  fakeHbrPolicyBindingsMapFunc,
}

func testAccCheckAlicloudHbrPolicyBindingsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tfacchbr%d"
}

resource "alicloud_hbr_vault" "default" {
  vault_type = "STANDARD"
  vault_name = var.name
}

resource "alicloud_hbr_policy" "default" {
  policy_name = var.name
  rules {
    rule_type    = "BACKUP"
    backup_type  = "COMPLETE"
    schedule     = "I|0|P1D"
    retention    = "7"
    vault_id     = alicloud_hbr_vault.default.id
    archive_days = "0"
  }
}

resource "alicloud_oss_bucket" "default" {
  storage_class = "Standard"
  bucket        = var.name
}

resource "alicloud_hbr_policy_binding" "default" {
  source_type                = "OSS"
  disabled                   = false
  policy_id                  = alicloud_hbr_policy.default.id
  data_source_id             = alicloud_oss_bucket.default.id
  policy_binding_description = "policy binding example"
  source                     = "prefix-example-create/"
  advanced_options {
    oss_detail {
      ignore_archive_object    = false
      inventory_cleanup_policy = "NO_CLEANUP"
    }
  }
}

data "alicloud_hbr_policy_bindings" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
