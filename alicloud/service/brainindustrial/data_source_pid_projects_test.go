package brainindustrial_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudBrainIndustrialPidProjectsDataSource(t *testing.T) {
	resourceId := "data.alicloud_brain_industrial_pid_projects.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceBrainIndustrialPidProjectsDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_brain_industrial_pid_project.default.pid_project_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_brain_industrial_pid_project.default.pid_project_name}-fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_brain_industrial_pid_project.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_brain_industrial_pid_project.default.id}-fake"},
		}),
	}
	var existBrainIndustrialPidProjectsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"ids.0":                          acctest.CHECKSET,
			"names.#":                        "1",
			"names.0":                        name,
			"projects.#":                     "1",
			"projects.0.pid_organization_id": acctest.CHECKSET,
			"projects.0.pid_project_desc":    "",
			"projects.0.pid_project_id":      acctest.CHECKSET,
			"projects.0.id":                  acctest.CHECKSET,
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

	var BrainIndustrialPidProjectsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existBrainIndustrialPidProjectsMapFunc,
		FakeMapFunc:  fakeBrainIndustrialPidProjectsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
	}

	BrainIndustrialPidProjectsInfo.DataSourceTestCheckWithPreCheck(t, 0, preCheck, nameRegexConf, idsConf)
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
