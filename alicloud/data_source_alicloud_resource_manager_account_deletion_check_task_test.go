package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerAccountDeletionCheckTaskDataSource(t *testing.T) {
	rand := acctest.RandInt()
	accountIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerAccountDeletionCheckTaskDataSourceName(rand, map[string]string{
			"account_id": `"${alicloud_resource_manager_account.default.id}"`,
		}),
	}
	var existAlicloudResourceManagerAccountDeletionCheckTaskDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"not_allow_reason.#":    "0",
			"abandon_able_checks.#": "0",
			"allow_delete":          CHECKSET,
			"status":                CHECKSET,
		}
	}
	var fakeAlicloudResourceManagerAccountDeletionCheckTaskDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"not_allow_reason.#":    "0",
			"abandon_able_checks.#": "0",
			"allow_delete":          NOSET,
			"status":                NOSET,
		}
	}
	var alicloudResourceManagerAccountDeletionCheckTaskCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_account_deletion_check_task.default",
		existMapFunc: existAlicloudResourceManagerAccountDeletionCheckTaskDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudResourceManagerAccountDeletionCheckTaskDataSourceNameMapFunc,
	}
	alicloudResourceManagerAccountDeletionCheckTaskCheckInfo.dataSourceTestCheck(t, rand, accountIdConf)
}

func testAccCheckAlicloudResourceManagerAccountDeletionCheckTaskDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "tf-test-%d"
	}
	
	resource "alicloud_resource_manager_account" "default" {
	  display_name = var.name
	}
	
	data "alicloud_resource_manager_account_deletion_check_task" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
