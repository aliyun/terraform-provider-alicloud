package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudResourceManagerAccountsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_resource_manager_accounts.default"
	name := fmt.Sprintf("tf-testAcc-rma%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceResourceManagerAccountsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}_fake"},
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status": "CreateSuccess",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status": "CreateFailed",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account_Fake",
			},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status": "CreateSuccess",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"accounts.0.id":                    acctest.CHECKSET,
			"accounts.0.account_id":            acctest.CHECKSET,
			"accounts.0.display_name":          acctest.CHECKSET,
			"accounts.0.type":                  acctest.CHECKSET,
			"accounts.0.folder_id":             acctest.CHECKSET,
			"accounts.0.resource_directory_id": acctest.CHECKSET,
			"accounts.0.status":                acctest.CHECKSET,
			"accounts.0.tags.%":                "2",
			"accounts.0.tags.Created":          "TF",
			"accounts.0.tags.For":              "Account",
			"accounts.0.account_name":          "",
			"accounts.0.payer_account_id":      "",
			"accounts.0.join_method":           acctest.CHECKSET,
			"accounts.0.join_time":             acctest.CHECKSET,
			"accounts.0.modify_time":           acctest.CHECKSET,
			"accounts.0.deletion_status":       acctest.CHECKSET,
		}
	}

	var fakeAliCloudResourceManagerAccountsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"accounts.#": "0",
		}
	}

	var aliCloudResourceManagerAccountsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_accounts.default",
		ExistMapFunc: existAliCloudResourceManagerAccountsDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudResourceManagerAccountsDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
	}

	aliCloudResourceManagerAccountsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, tagsConf, allConf)
}

func TestAccAliCloudResourceManagerAccountsDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_resource_manager_accounts.default"
	name := fmt.Sprintf("tf-testAcc-rma%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceResourceManagerAccountsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"enable_details": "false",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status":         "CreateSuccess",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status":         "CreateSuccess",
			"enable_details": "false",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
			"enable_details": "false",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${data.alicloud_resource_manager_accounts.test.accounts.0.id}"},
			"status": "CreateSuccess",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Account",
			},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"accounts.0.id":                    acctest.CHECKSET,
			"accounts.0.account_id":            acctest.CHECKSET,
			"accounts.0.display_name":          acctest.CHECKSET,
			"accounts.0.type":                  acctest.CHECKSET,
			"accounts.0.folder_id":             acctest.CHECKSET,
			"accounts.0.resource_directory_id": acctest.CHECKSET,
			"accounts.0.status":                acctest.CHECKSET,
			"accounts.0.tags.%":                "2",
			"accounts.0.tags.Created":          "TF",
			"accounts.0.tags.For":              "Account",
			"accounts.0.account_name":          acctest.CHECKSET,
			"accounts.0.payer_account_id":      "",
			"accounts.0.join_method":           acctest.CHECKSET,
			"accounts.0.join_time":             acctest.CHECKSET,
			"accounts.0.modify_time":           acctest.CHECKSET,
		}
	}

	var fakeAliCloudResourceManagerAccountsDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"accounts.#":                       "1",
			"accounts.0.id":                    acctest.CHECKSET,
			"accounts.0.account_id":            acctest.CHECKSET,
			"accounts.0.display_name":          acctest.CHECKSET,
			"accounts.0.type":                  acctest.CHECKSET,
			"accounts.0.folder_id":             acctest.CHECKSET,
			"accounts.0.resource_directory_id": acctest.CHECKSET,
			"accounts.0.status":                acctest.CHECKSET,
			"accounts.0.tags.%":                "2",
			"accounts.0.tags.Created":          "TF",
			"accounts.0.tags.For":              "Account",
			"accounts.0.account_name":          "",
			"accounts.0.payer_account_id":      "",
			"accounts.0.join_method":           acctest.CHECKSET,
			"accounts.0.join_time":             acctest.CHECKSET,
			"accounts.0.modify_time":           acctest.CHECKSET,
			"accounts.0.deletion_status":       acctest.CHECKSET,
		}
	}

	var aliCloudResourceManagerAccountsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_accounts.default",
		ExistMapFunc: existAliCloudResourceManagerAccountsDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudResourceManagerAccountsDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
	}

	aliCloudResourceManagerAccountsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, tagsConf, allConf)
}

func dataSourceResourceManagerAccountsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_accounts" "test"{
      ids = ["1893454309173762"]
	}
`, name)
}
