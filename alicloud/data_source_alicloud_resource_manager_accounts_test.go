package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerAccountsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	id := os.Getenv("ALICLOUD_RESOURCE_MANAGER_ACCOUNT_ID")
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerAccountsSourceConfig(rand, map[string]string{
			"ids":    fmt.Sprintf(`["%s"]`, id),
			"status": `"CreateSuccess"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerAccountsSourceConfig(rand, map[string]string{
			"ids":    fmt.Sprintf(`["%s_fake"]`, id),
			"status": `"CreateFailed"`,
		}),
	}

	var existResourceManagerAccountsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#":                       "1",
			"ids.#":                            "1",
			"accounts.0.id":                    CHECKSET,
			"accounts.0.account_id":            id,
			"accounts.0.display_name":          CHECKSET,
			"accounts.0.folder_id":             CHECKSET,
			"accounts.0.join_method":           CHECKSET,
			"accounts.0.join_time":             CHECKSET,
			"accounts.0.modify_time":           CHECKSET,
			"accounts.0.resource_directory_id": CHECKSET,
			"accounts.0.status":                CHECKSET,
			"accounts.0.type":                  CHECKSET,
			"accounts.0.payer_account_id":      "",
			"accounts.0.account_name":          CHECKSET,
		}
	}

	var fakeResourceManagerAccountsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#": "0",
			"ids.#":      "0",
		}
	}

	var accountsRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_accounts.default",
		existMapFunc: existResourceManagerAccountsRecordsMapFunc,
		fakeMapFunc:  fakeResourceManagerAccountsRecordsMapFunc,
	}

	var preCheck = func() {
		testAccPreCheck(t)
		testAccPreCheckWithResourceManagerAccountsSetting(t)
	}

	accountsRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)

}

func testAccCheckAlicloudResourceManagerAccountsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_resource_manager_accounts" "default"{
	enable_details = true
%s
}
`, strings.Join(pairs, "\n   "))
	return config
}
