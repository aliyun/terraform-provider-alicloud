package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHbrHanaInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.HBRSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_hbr_hana_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_hbr_hana_instance.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"name_regex": `["${alicloud_hbr_hana_instance.default.hana_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"name_regex": `["${alicloud_hbr_hana_instance.default.hana_name}_fake"]`,
		}),
	}
	vaultIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_hbr_hana_instance.default.id}"]`,
			"vault_id": `"${alicloud_hbr_hana_instance.default.vault_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_hbr_hana_instance.default.id}"]`,
			"vault_id": `"${alicloud_hbr_hana_instance.default.vault_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_hbr_hana_instance.default.id}"]`,
			"status": `"${alicloud_hbr_hana_instance.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_hbr_hana_instance.default.id}"]`,
			"status": `"INVALID_HANA_NODE"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_hbr_hana_instance.default.id}"]`,
			"status":     `"${alicloud_hbr_hana_instance.default.status}"`,
			"vault_id":   `"${alicloud_hbr_hana_instance.default.vault_id}"`,
			"name_regex": `["${alicloud_hbr_hana_instance.default.hana_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_hbr_hana_instance.default.id}_fake"]`,
			"status":     `"INVALID_HANA_NODE"`,
			"vault_id":   `"${alicloud_hbr_hana_instance.default.vault_id}_fake"`,
			"name_regex": `["${alicloud_hbr_hana_instance.default.hana_name}_fake"]`,
		}),
	}
	var existAlicloudHbrHanaInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"instances.#":                      "1",
			"instances.0.alert_setting":        "INHERITED",
			"instances.0.hana_name":            fmt.Sprintf("tf-testAccHanaInstance-%d", rand),
			"instances.0.host":                 "1.1.1.1",
			"instances.0.instance_number":      "1",
			"instances.0.resource_group_id":    CHECKSET,
			"instances.0.use_ssl":              "false",
			"instances.0.user_name":            "admin",
			"instances.0.validate_certificate": "false",
			"instances.0.vault_id":             CHECKSET,
			"instances.0.id":                   CHECKSET,
			"instances.0.hana_instance_id":     CHECKSET,
			"instances.0.status":               CHECKSET,
			"instances.0.status_message":       CHECKSET,
		}
	}
	var fakeAlicloudHbrHanaInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudHbrHanaInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_hbr_hana_instances.default",
		existMapFunc: existAlicloudHbrHanaInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudHbrHanaInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudHbrHanaInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vaultIdConf, statusConf, allConf)
}
func testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccHanaInstance-%d"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}
data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_hbr_hana_instance" "default" {
	alert_setting = "INHERITED"
	hana_name = var.name
	host = "1.1.1.1"
	instance_number = "1"
	password = "YouPassword123"
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	sid = "HXE"
	use_ssl = "false"
	user_name = "admin"
	validate_certificate = "false"
	vault_id = alicloud_hbr_vault.default.id
}

data "alicloud_hbr_hana_instances" "default" {	
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
