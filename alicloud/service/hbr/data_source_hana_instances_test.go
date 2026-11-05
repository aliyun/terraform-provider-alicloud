package hbr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudHbrHanaInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.HBRSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_hbr_hana_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_hbr_hana_instance.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_hana_instance.default.hana_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_hana_instance.default.hana_name}_fake"`,
		}),
	}
	vaultIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_hbr_hana_instance.default.id}"]`,
			"vault_id": `"${alicloud_hbr_hana_instance.default.vault_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_hbr_hana_instance.default.id}"]`,
			"vault_id": `"${alicloud_hbr_hana_instance.default.vault_id}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_hbr_hana_instance.default.id}"]`,
			"status": `"${alicloud_hbr_hana_instance.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_hbr_hana_instance.default.id}"]`,
			"status": `"INVALID_HANA_NODE"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_hbr_hana_instance.default.id}"]`,
			"status":     `"${alicloud_hbr_hana_instance.default.status}"`,
			"vault_id":   `"${alicloud_hbr_hana_instance.default.vault_id}"`,
			"name_regex": `"${alicloud_hbr_hana_instance.default.hana_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_hbr_hana_instance.default.id}_fake"]`,
			"status":     `"INVALID_HANA_NODE"`,
			"vault_id":   `"${alicloud_hbr_hana_instance.default.vault_id}_fake"`,
			"name_regex": `"${alicloud_hbr_hana_instance.default.hana_name}_fake"`,
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
			"instances.0.resource_group_id":    acctest.CHECKSET,
			"instances.0.use_ssl":              "false",
			"instances.0.user_name":            "admin",
			"instances.0.validate_certificate": "false",
			"instances.0.vault_id":             acctest.CHECKSET,
			"instances.0.id":                   acctest.CHECKSET,
			"instances.0.hana_instance_id":     acctest.CHECKSET,
			"instances.0.status":               acctest.CHECKSET,
			"instances.0.status_message":       acctest.CHECKSET,
		}
	}
	var fakeAlicloudHbrHanaInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudHbrHanaInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_hbr_hana_instances.default",
		ExistMapFunc: existAlicloudHbrHanaInstancesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudHbrHanaInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudHbrHanaInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vaultIdConf, statusConf, allConf)
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
