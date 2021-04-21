package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerFoldersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerFoldersSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_folder.example.folder_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerFoldersSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_folder.example.folder_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerFoldersSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_folder.example.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerFoldersSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_folder.example.id}_fake"]`,
		}),
	}

	parentFolderIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerFoldersSourceConfig(rand, map[string]string{
			"name_regex":       `"${alicloud_resource_manager_folder.example.folder_name}"`,
			"parent_folder_id": `"${alicloud_resource_manager_folder.example.parent_folder_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerFoldersSourceConfig(rand, map[string]string{
			"name_regex":       `"${alicloud_resource_manager_folder.example.folder_name}_fake"`,
			"parent_folder_id": `"${alicloud_resource_manager_folder.example.parent_folder_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerFoldersSourceConfig(rand, map[string]string{
			"name_regex":       `"${alicloud_resource_manager_folder.example.folder_name}"`,
			"ids":              `["${alicloud_resource_manager_folder.example.id}"]`,
			"parent_folder_id": `"${alicloud_resource_manager_folder.example.parent_folder_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerFoldersSourceConfig(rand, map[string]string{
			"name_regex":       `"${alicloud_resource_manager_folder.example.folder_name}_fake"`,
			"ids":              `["${alicloud_resource_manager_folder.example.id}"]`,
			"parent_folder_id": `"${alicloud_resource_manager_folder.example.parent_folder_id}"`,
		}),
	}

	var existResourceManagerFoldersRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"folders.#":                  "1",
			"names.#":                    "1",
			"ids.#":                      "1",
			"folders.0.id":               CHECKSET,
			"folders.0.folder_id":        CHECKSET,
			"folders.0.folder_name":      fmt.Sprintf("tf-testAcc-%d", rand),
			"folders.0.parent_folder_id": CHECKSET,
		}
	}

	var fakeResourceManagerFoldersRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"folders.#": "0",
			"ids.#":     "0",
			"names.#":   "0",
		}
	}

	var foldersRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_folders.example",
		existMapFunc: existResourceManagerFoldersRecordsMapFunc,
		fakeMapFunc:  fakeResourceManagerFoldersRecordsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckEnterpriseAccountEnabled(t)
	}

	foldersRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, parentFolderIdConf, allConf)

}

func testAccCheckAlicloudResourceManagerFoldersSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
resource "alicloud_resource_manager_folder" "example"{
	folder_name = "tf-testAcc-%d"
}

data "alicloud_resource_manager_folders" "example"{
	enable_details = true
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}
