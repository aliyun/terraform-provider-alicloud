package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudBrainIndustrialPidProjectsDataSource(t *testing.T) {
	resourceId := "data.alicloud_brain_industrial_pid_projects.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceBrainIndustrialPidProjectsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_brain_industrial_pid_project.default.pid_project_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_brain_industrial_pid_project.default.pid_project_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_brain_industrial_pid_project.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_brain_industrial_pid_project.default.id}-fake"},
		}),
	}
	var existBrainIndustrialPidProjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"ids.0":                          CHECKSET,
			"names.#":                        "1",
			"names.0":                        name,
			"projects.#":                     "1",
			"projects.0.pid_organization_id": CHECKSET,
			"projects.0.pid_project_desc":    "",
			"projects.0.pid_project_id":      CHECKSET,
			"projects.0.id":                  CHECKSET,
			"projects.0.pid_project_name":    name,
		}
	}

	var fakeBrainIndustrialPidProjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"projects.#": "0",
		}
	}

	var BrainIndustrialPidProjectsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existBrainIndustrialPidProjectsMapFunc,
		fakeMapFunc:  fakeBrainIndustrialPidProjectsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
	}

	BrainIndustrialPidProjectsInfo.dataSourceTestCheckWithPreCheck(t, 0, preCheck, nameRegexConf, idsConf)
}

func dataSourceBrainIndustrialPidProjectsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_brain_industrial_pid_organization" "default" {
		pid_organization_name = "%[1]s"
	}
	resource "alicloud_brain_industrial_pid_project" "default" {
		pid_organization_id = alicloud_brain_industrial_pid_organization.default.id
		pid_project_name = "%[1]s"
	}`, name)
}
