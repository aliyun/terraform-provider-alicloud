package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerDelegatedAdministratorsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerDelegatedAdministratorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_delegated_administrator.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerDelegatedAdministratorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_delegated_administrator.default.id}_fake"]`,
		}),
	}
	servicePrincipalConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerDelegatedAdministratorsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_resource_manager_delegated_administrator.default.id}"]`,
			"service_principal": `"${alicloud_resource_manager_delegated_administrator.default.service_principal}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerDelegatedAdministratorsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_resource_manager_delegated_administrator.default.id}_fake"]`,
			"service_principal": `"${alicloud_resource_manager_delegated_administrator.default.service_principal}_fake"`,
		}),
	}
	var existAlicloudResourceManagerDelegatedAdministratorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"administrators.#":                         "1",
			"administrators.0.account_id":              CHECKSET,
			"administrators.0.service_principal":       CHECKSET,
			"administrators.0.delegation_enabled_time": CHECKSET,
			"administrators.0.id":                      CHECKSET,
		}
	}
	var fakeAlicloudResourceManagerDelegatedAdministratorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudResourceManagerDelegatedAdministratorsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_delegated_administrators.default",
		existMapFunc: existAlicloudResourceManagerDelegatedAdministratorsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudResourceManagerDelegatedAdministratorsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudResourceManagerDelegatedAdministratorsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, servicePrincipalConf)
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
