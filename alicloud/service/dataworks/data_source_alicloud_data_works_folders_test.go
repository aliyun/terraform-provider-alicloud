package dataworks_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudDataWorksFoldersDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_data_works_folders.default"
	name := fmt.Sprintf("tf-testacc-dataworksfolder%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceDataworksFoldersConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"project_id":         "${alicloud_data_works_folder.default.project_id}",
			"ids":                []string{"${alicloud_data_works_folder.default.folder_id}"},
			"parent_folder_path": "Business Flow/tfTestAcc/folderDi",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"project_id":         "${alicloud_data_works_folder.default.project_id}",
			"ids":                []string{"${alicloud_data_works_folder.default.folder_id}_fake"},
			"parent_folder_path": "Business Flow/tfTestAcc/folderDi",
		}),
	}

	var existDataworksFoldersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"folders.#":            "1",
			"folders.0.project_id": acctest.CHECKSET,
			"folders.0.folder_id":  acctest.CHECKSET,
		}
	}

	var fakeDataworksFoldersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"folders.#": "0",
			"ids.#":     "0",
		}
	}

	var DataworksFoldersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDataworksFoldersMapFunc,
		FakeMapFunc:  fakeDataworksFoldersMapFunc,
	}

	DataworksFoldersCheckInfo.DataSourceTestCheck(t, rand, idsConf)
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
