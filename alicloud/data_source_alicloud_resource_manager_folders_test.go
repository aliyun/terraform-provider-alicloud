package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudResourceManagerFoldersDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_resource_manager_folders.default"
	name := fmt.Sprintf("tf-testAcc-ResourceManagerFolder%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceResourceManagerFoldersConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_resource_manager_folder.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_resource_manager_folder.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_resource_manager_folder.default.folder_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_resource_manager_folder.default.folder_name}_fake",
		}),
	}

	parentFolderIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"parent_folder_id": "${alicloud_resource_manager_folder.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"parent_folder_id": "${alicloud_resource_manager_folder.default.id}",
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}_fake"},
		}),
	}

	queryKeywordConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"query_keyword": "${alicloud_resource_manager_folder.default.folder_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"query_keyword": "${alicloud_resource_manager_folder.default.folder_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}"},
			"name_regex":       "${alicloud_resource_manager_folder.default.folder_name}",
			"parent_folder_id": "${alicloud_resource_manager_folder.default.parent_folder_id}",
			"query_keyword":    "${alicloud_resource_manager_folder.default.folder_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"folders.0.id":               CHECKSET,
			"folders.0.folder_id":        CHECKSET,
			"folders.0.folder_name":      CHECKSET,
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

	var aliCloudResourceManagerFoldersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_folders.default",
		existMapFunc: existAliCloudResourceManagerFoldersDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudResourceManagerFoldersDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudResourceManagerFoldersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, parentFolderIdConf, queryKeywordConf, allConf)
}

func TestAccAliCloudResourceManagerFoldersDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_resource_manager_folders.default"
	name := fmt.Sprintf("tf-testAcc-ResourceManagerFolder%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceResourceManagerFoldersConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_resource_manager_folder.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_resource_manager_folder.default.id}_fake"},
			"enable_details": "true",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_resource_manager_folder.default.folder_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_resource_manager_folder.default.folder_name}_fake",
			"enable_details": "true",
		}),
	}

	parentFolderIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"parent_folder_id": "${alicloud_resource_manager_folder.default.id}",
			"enable_details":   "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"parent_folder_id": "${alicloud_resource_manager_folder.default.id}",
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}_fake"},
			"enable_details":   "true",
		}),
	}

	queryKeywordConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"query_keyword":  "${alicloud_resource_manager_folder.default.folder_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"query_keyword":  "${alicloud_resource_manager_folder.default.folder_name}_fake",
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_resource_manager_folder.default.id}"},
			"name_regex":       "${alicloud_resource_manager_folder.default.folder_name}",
			"parent_folder_id": "${alicloud_resource_manager_folder.default.parent_folder_id}",
			"query_keyword":    "${alicloud_resource_manager_folder.default.folder_name}",
			"enable_details":   "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"folders.0.id":               CHECKSET,
			"folders.0.folder_id":        CHECKSET,
			"folders.0.folder_name":      CHECKSET,
			"folders.0.parent_folder_id": CHECKSET,
		}
	}

	var fakeAliCloudResourceManagerFoldersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"folders.#": "0",
		}
	}

	var aliCloudResourceManagerFoldersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_folders.default",
		existMapFunc: existAliCloudResourceManagerFoldersDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudResourceManagerFoldersDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudResourceManagerFoldersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, parentFolderIdConf, queryKeywordConf, allConf)
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
