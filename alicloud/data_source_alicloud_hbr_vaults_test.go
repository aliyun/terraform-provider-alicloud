package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBRVaultsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	vaultIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrVaultSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_vault.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrVaultSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_vault.default.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrVaultSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_hbr_vault.default.id}"]`,
			"status": `"CREATED"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrVaultSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_hbr_vault.default.id}"]`,
			"status": `"ERROR"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrVaultSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_vault.default.vault_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrVaultSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_vault.default.vault_name}_fake"`,
		}),
	}

	vaultTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrVaultSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_hbr_vault.default.id}"]`,
			"vault_type": `"STANDARD"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrVaultSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_hbr_vault.default.id}"]`,
			"name_regex": `"${alicloud_hbr_vault.default.vault_name}"`,
			"status":     `"CREATED"`,
			"vault_type": `"STANDARD"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrVaultSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_hbr_vault.default.id}"]`,
			"name_regex": `"${alicloud_hbr_vault.default.vault_name}_fake"`,
			"status":     `"ERROR"`,
			"vault_type": `"STANDARD"`,
		}),
	}

	HbrVaultCheckInfo.dataSourceTestCheck(t, rand, vaultIdConf, statusConf, nameRegexConf, vaultTypeConf, allConf)
}

func testAccCheckAlicloudHbrVaultSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccHbrVaultsDataSource%d"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

data "alicloud_hbr_vaults" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existHbrVaultMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vaults.#":             "1",
		"vaults.0.vault_id":    CHECKSET,
		"vaults.0.vault_name":  fmt.Sprintf("tf-testAccHbrVaultsDataSource%d", rand),
		"vaults.0.vault_type":  "STANDARD",
		"vaults.0.description": "",
		"vaults.0.status":      "CREATED",
	}
}

var fakeHbrVaultMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vaults.#": "0",
	}
}

var HbrVaultCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_vaults.default",
	existMapFunc: existHbrVaultMapFunc,
	fakeMapFunc:  fakeHbrVaultMapFunc,
}
