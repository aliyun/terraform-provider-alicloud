package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudDataWorksDataServiceGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_data_works_data_service_groups.default"
	name := fmt.Sprintf("tf_testacc_dsg%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataWorksDataServiceGroupsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alicloud_data_works_data_service_group.default.project_id}",
			"ids":        []string{"${alicloud_data_works_data_service_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alicloud_data_works_data_service_group.default.project_id}",
			"ids":        []string{"${alicloud_data_works_data_service_group.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alicloud_data_works_data_service_group.default.project_id}",
			"name_regex": "${alicloud_data_works_data_service_group.default.data_service_group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alicloud_data_works_data_service_group.default.project_id}",
			"name_regex": "${alicloud_data_works_data_service_group.default.data_service_group_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alicloud_data_works_data_service_group.default.project_id}",
			"ids":        []string{"${alicloud_data_works_data_service_group.default.id}"},
			"name_regex": "${alicloud_data_works_data_service_group.default.data_service_group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_id": "${alicloud_data_works_data_service_group.default.project_id}",
			"ids":        []string{"${alicloud_data_works_data_service_group.default.id}_fake"},
			"name_regex": "${alicloud_data_works_data_service_group.default.data_service_group_name}",
		}),
	}

	var existDataWorksDataServiceGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"groups.#":                         "1",
			"groups.0.id":                      CHECKSET,
			"groups.0.data_service_group_id":   CHECKSET,
			"groups.0.data_service_group_name": CHECKSET,
			"groups.0.api_gateway_group_id":    CHECKSET,
			"groups.0.project_id":              CHECKSET,
			"groups.0.description":             CHECKSET,
			"groups.0.create_time":             CHECKSET,
			"groups.0.creator_id":              CHECKSET,
			"groups.0.modified_time":           CHECKSET,
		}
	}

	var fakeDataWorksDataServiceGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
		}
	}

	var DataWorksDataServiceGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDataWorksDataServiceGroupsMapFunc,
		fakeMapFunc:  fakeDataWorksDataServiceGroupsMapFunc,
	}

	DataWorksDataServiceGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceDataWorksDataServiceGroupsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_data_works_project" "default" {
  project_name     = var.name
  display_name     = var.name
  description      = var.name
  pai_task_enabled = false
}

resource "alicloud_api_gateway_group" "default" {
  name        = var.name
  description = var.name
}

resource "alicloud_data_works_data_service_group" "default" {
  project_id              = alicloud_data_works_project.default.id
  api_gateway_group_id    = alicloud_api_gateway_group.default.id
  data_service_group_name = var.name
  description             = var.name
}
`, name)
}
