package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudResourceManagerFoldersDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_resource_manager_folders.default"
	name := fmt.Sprintf("tf-testAcc-ResourceManagerFolder%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceResourceManagerFoldersConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_resource_manager_folder.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_resource_manager_folder.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_resource_manager_folder.default.folder_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_resource_manager_folder.default.folder_name}_fake",
		}),
	}

	parentFolderIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"parent_folder_id": "${alicloud_resource_manager_folder.default.id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"parent_folder_id": "${alicloud_resource_manager_folder.default.id}",
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}_fake"},
		}),
	}

	queryKeywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"query_keyword": "${alicloud_resource_manager_folder.default.folder_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"query_keyword": "${alicloud_resource_manager_folder.default.folder_name}_fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}"},
			"name_regex":       "${alicloud_resource_manager_folder.default.folder_name}",
			"parent_folder_id": "${alicloud_resource_manager_folder.default.parent_folder_id}",
			"query_keyword":    "${alicloud_resource_manager_folder.default.folder_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}_fake"},
			"name_regex":       "${alicloud_resource_manager_folder.default.folder_name}_fake",
			"parent_folder_id": "${alicloud_resource_manager_folder.default.parent_folder_id}",
			"query_keyword":    "${alicloud_resource_manager_folder.default.folder_name}_fake",
		}),
	}

	var existAliCloudResourceManagerFoldersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"names.#":                    "1",
			"folders.#":                  "1",
			"folders.0.id":               acctest.CHECKSET,
			"folders.0.folder_id":        acctest.CHECKSET,
			"folders.0.folder_name":      acctest.CHECKSET,
			"folders.0.parent_folder_id": "",
		}
	}

	var fakeAliCloudResourceManagerFoldersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"folders.#": "0",
		}
	}

	var aliCloudResourceManagerFoldersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_folders.default",
		ExistMapFunc: existAliCloudResourceManagerFoldersDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudResourceManagerFoldersDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudResourceManagerFoldersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, parentFolderIdConf, queryKeywordConf, allConf)
}

func TestAccAliCloudResourceManagerFoldersDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_resource_manager_folders.default"
	name := fmt.Sprintf("tf-testAcc-ResourceManagerFolder%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceResourceManagerFoldersConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_resource_manager_folder.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_resource_manager_folder.default.id}_fake"},
			"enable_details": "true",
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_resource_manager_folder.default.folder_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_resource_manager_folder.default.folder_name}_fake",
			"enable_details": "true",
		}),
	}

	parentFolderIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"parent_folder_id": "${alicloud_resource_manager_folder.default.id}",
			"enable_details":   "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"parent_folder_id": "${alicloud_resource_manager_folder.default.id}",
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}_fake"},
			"enable_details":   "true",
		}),
	}

	queryKeywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"query_keyword":  "${alicloud_resource_manager_folder.default.folder_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"query_keyword":  "${alicloud_resource_manager_folder.default.folder_name}_fake",
			"enable_details": "true",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}"},
			"name_regex":       "${alicloud_resource_manager_folder.default.folder_name}",
			"parent_folder_id": "${alicloud_resource_manager_folder.default.parent_folder_id}",
			"query_keyword":    "${alicloud_resource_manager_folder.default.folder_name}",
			"enable_details":   "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}_fake"},
			"name_regex":       "${alicloud_resource_manager_folder.default.folder_name}_fake",
			"parent_folder_id": "${alicloud_resource_manager_folder.default.parent_folder_id}",
			"query_keyword":    "${alicloud_resource_manager_folder.default.folder_name}_fake",
			"enable_details":   "true",
		}),
	}

	var existAliCloudResourceManagerFoldersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"names.#":                    "1",
			"folders.#":                  "1",
			"folders.0.id":               acctest.CHECKSET,
			"folders.0.folder_id":        acctest.CHECKSET,
			"folders.0.folder_name":      acctest.CHECKSET,
			"folders.0.parent_folder_id": acctest.CHECKSET,
		}
	}

	var fakeAliCloudResourceManagerFoldersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"folders.#": "0",
		}
	}

	var aliCloudResourceManagerFoldersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_resource_manager_folders.default",
		ExistMapFunc: existAliCloudResourceManagerFoldersDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudResourceManagerFoldersDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudResourceManagerFoldersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, parentFolderIdConf, queryKeywordConf, allConf)
}

func dataSourceResourceManagerFoldersConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_resource_manager_folder" "default" {
  		folder_name = "${var.name}-default"
	}

	resource "alicloud_resource_manager_folder" "sub" {
  		folder_name      = "${var.name}-sub"
  		parent_folder_id = alicloud_resource_manager_folder.default.id
	}
`, name)
}
