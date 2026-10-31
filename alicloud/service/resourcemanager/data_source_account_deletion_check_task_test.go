package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudResourceManagerAccountDeletionCheckTaskDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	accountIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudResourceManagerAccountDeletionCheckTaskDataSourceName(rand, map[string]string{
			"account_id": `"${alicloud_resource_manager_account.default.id}"`,
		}),
	}
	var existAlicloudResourceManagerAccountDeletionCheckTaskDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"not_allow_reason.#":    "0",
			"abandon_able_checks.#": "0",
			"allow_delete":          acctest.CHECKSET,
			"status":                acctest.CHECKSET,
		}
	}
	var fakeAlicloudResourceManagerAccountDeletionCheckTaskDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"not_allow_reason.#":    "0",
			"abandon_able_checks.#": "0",
			"allow_delete":          acctest.NOSET,
			"status":                acctest.NOSET,
		}
	}
	var alicloudResourceManagerAccountDeletionCheckTaskCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_account_deletion_check_task.default",
		ExistMapFunc: existAlicloudResourceManagerAccountDeletionCheckTaskDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudResourceManagerAccountDeletionCheckTaskDataSourceNameMapFunc,
	}
	alicloudResourceManagerAccountDeletionCheckTaskCheckInfo.DataSourceTestCheck(t, rand, accountIdConf)
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
