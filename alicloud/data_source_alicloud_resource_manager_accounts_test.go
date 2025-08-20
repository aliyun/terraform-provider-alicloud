package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudResourceManagerAccountsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_resource_manager_accounts.default"
	name := fmt.Sprintf("tf-testAcc-rma%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceResourceManagerAccountsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status": "CreateSuccess",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status": "CreateFailed",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account_Fake",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status": "CreateSuccess",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}_fake"},
			"status": "CreateFailed",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account_Fake",
			},
		}),
	}

	var existAliCloudResourceManagerAccountsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"accounts.#":                       "1",
			"accounts.0.id":                    CHECKSET,
			"accounts.0.account_id":            CHECKSET,
			"accounts.0.display_name":          CHECKSET,
			"accounts.0.type":                  CHECKSET,
			"accounts.0.folder_id":             CHECKSET,
			"accounts.0.resource_directory_id": CHECKSET,
			"accounts.0.status":                CHECKSET,
			"accounts.0.tags.%":                "2",
			"accounts.0.tags.Created":          "TF",
			"accounts.0.tags.For":              "Account",
			"accounts.0.account_name":          "",
			"accounts.0.payer_account_id":      "",
			"accounts.0.join_method":           CHECKSET,
			"accounts.0.join_time":             CHECKSET,
			"accounts.0.modify_time":           CHECKSET,
		}
	}

	var fakeAliCloudResourceManagerAccountsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"accounts.#": "0",
		}
	}

	var aliCloudResourceManagerAccountsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_accounts.default",
		existMapFunc: existAliCloudResourceManagerAccountsDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudResourceManagerAccountsDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	aliCloudResourceManagerAccountsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, tagsConf, allConf)
}

func TestAccAliCloudResourceManagerAccountsDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_resource_manager_accounts.default"
	name := fmt.Sprintf("tf-testAcc-rma%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceResourceManagerAccountsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"enable_details": "false",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status":         "CreateSuccess",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status":         "CreateSuccess",
			"enable_details": "false",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
			"enable_details": "false",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status": "CreateSuccess",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status": "CreateSuccess",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
			"enable_details": "false",
		}),
	}

	var existAliCloudResourceManagerAccountsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"accounts.#":                       "1",
			"accounts.0.id":                    CHECKSET,
			"accounts.0.account_id":            CHECKSET,
			"accounts.0.display_name":          CHECKSET,
			"accounts.0.type":                  CHECKSET,
			"accounts.0.folder_id":             CHECKSET,
			"accounts.0.resource_directory_id": CHECKSET,
			"accounts.0.status":                CHECKSET,
			"accounts.0.tags.%":                "2",
			"accounts.0.tags.Created":          "TF",
			"accounts.0.tags.For":              "Account",
			"accounts.0.account_name":          CHECKSET,
			"accounts.0.payer_account_id":      "",
			"accounts.0.join_method":           CHECKSET,
			"accounts.0.join_time":             CHECKSET,
			"accounts.0.modify_time":           CHECKSET,
		}
	}

	var fakeAliCloudResourceManagerAccountsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"accounts.#":                       "1",
			"accounts.0.id":                    CHECKSET,
			"accounts.0.account_id":            CHECKSET,
			"accounts.0.display_name":          CHECKSET,
			"accounts.0.type":                  CHECKSET,
			"accounts.0.folder_id":             CHECKSET,
			"accounts.0.resource_directory_id": CHECKSET,
			"accounts.0.status":                CHECKSET,
			"accounts.0.tags.%":                "2",
			"accounts.0.tags.Created":          "TF",
			"accounts.0.tags.For":              "Account",
			"accounts.0.account_name":          "",
			"accounts.0.payer_account_id":      "",
			"accounts.0.join_method":           CHECKSET,
			"accounts.0.join_time":             CHECKSET,
			"accounts.0.modify_time":           CHECKSET,
		}
	}

	var aliCloudResourceManagerAccountsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_accounts.default",
		existMapFunc: existAliCloudResourceManagerAccountsDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudResourceManagerAccountsDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	aliCloudResourceManagerAccountsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, tagsConf, allConf)
}

func dataSourceResourceManagerAccountsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_accounts" "test"{
	}
`, name)
}
