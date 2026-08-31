package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDataWorksFoldersDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_data_works_folders.default"
	name := fmt.Sprintf("tf-testacc-dataworksfolder%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataworksFoldersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_id":         "${alicloud_data_works_folder.default.project_id}",
			"ids":                []string{"${alicloud_data_works_folder.default.folder_id}"},
			"parent_folder_path": "Business Flow/tfTestAcc/folderDi",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_id":         "${alicloud_data_works_folder.default.project_id}",
			"ids":                []string{"${alicloud_data_works_folder.default.folder_id}_fake"},
			"parent_folder_path": "Business Flow/tfTestAcc/folderDi",
		}),
	}

	var existDataworksFoldersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"folders.#":            "1",
			"folders.0.project_id": CHECKSET,
			"folders.0.folder_id":  CHECKSET,
		}
	}

	var fakeDataworksFoldersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"folders.#": "0",
			"ids.#":     "0",
		}
	}

	var DataworksFoldersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDataworksFoldersMapFunc,
		fakeMapFunc:  fakeDataworksFoldersMapFunc,
	}

	DataworksFoldersCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceDataworksFoldersConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		
		resource "alicloud_data_works_folder" "default" {
			project_id = "638"
			folder_path = "Business Flow/tfTestAcc/folderDi/tftest1"
		}
		`, name)
}
