package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudResourceManagerDelegatedAdministratorsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerDelegatedAdministratorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_delegated_administrator.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerDelegatedAdministratorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_delegated_administrator.default.id}_fake"]`,
		}),
	}
	servicePrincipalConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerDelegatedAdministratorsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_resource_manager_delegated_administrator.default.id}"]`,
			"service_principal": `"${alicloud_resource_manager_delegated_administrator.default.service_principal}"`,
		}),
		FakeConfig: testAccCheckAlicloudResourceManagerDelegatedAdministratorsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_resource_manager_delegated_administrator.default.id}_fake"]`,
			"service_principal": `"${alicloud_resource_manager_delegated_administrator.default.service_principal}_fake"`,
		}),
	}
	var existAlicloudResourceManagerDelegatedAdministratorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"administrators.#":                         "1",
			"administrators.0.account_id":              acctest.CHECKSET,
			"administrators.0.service_principal":       acctest.CHECKSET,
			"administrators.0.delegation_enabled_time": acctest.CHECKSET,
			"administrators.0.id":                      acctest.CHECKSET,
		}
	}
	var fakeAlicloudResourceManagerDelegatedAdministratorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudResourceManagerDelegatedAdministratorsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_delegated_administrators.default",
		ExistMapFunc: existAlicloudResourceManagerDelegatedAdministratorsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudResourceManagerDelegatedAdministratorsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudResourceManagerDelegatedAdministratorsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, servicePrincipalConf)
}
func testAccCheckAlicloudResourceManagerDelegatedAdministratorsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDelegatedAdministrator-%d"
}

data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}
resource "alicloud_resource_manager_delegated_administrator" "default" {
	account_id = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
	service_principal = "cloudfw.aliyuncs.com"
}

data "alicloud_resource_manager_delegated_administrators" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
