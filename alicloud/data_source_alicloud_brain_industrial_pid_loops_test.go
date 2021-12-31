package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudBrainIndustrialPidLoopsDataSource(t *testing.T) {
	resourceId := "data.alicloud_brain_industrial_pid_loops.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceBrainIndustrialPidLoopsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"pid_project_id": "${alicloud_brain_industrial_pid_project.default.id}",
			"name_regex":     "${alicloud_brain_industrial_pid_loop.default.pid_loop_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"pid_project_id": "${alicloud_brain_industrial_pid_project.default.id}",
			"name_regex":     "${alicloud_brain_industrial_pid_loop.default.pid_loop_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"pid_project_id": "${alicloud_brain_industrial_pid_project.default.id}",
			"ids":            []string{"${alicloud_brain_industrial_pid_loop.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"pid_project_id": "${alicloud_brain_industrial_pid_project.default.id}",
			"ids":            []string{"${alicloud_brain_industrial_pid_loop.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"pid_project_id": "${alicloud_brain_industrial_pid_project.default.id}",
			"status":         "0",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"pid_project_id": "${alicloud_brain_industrial_pid_project.default.id}",
			"status":         "1",
		}),
	}
	var existBrainIndustrialPidLoopsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"names.#":                     "1",
			"names.0":                     name,
			"loops.#":                     "1",
			"loops.0.pid_loop_dcs_type":   "standard",
			"loops.0.pid_loop_is_crucial": "false",
			"loops.0.pid_loop_id":         CHECKSET,
			"loops.0.id":                  CHECKSET,
			"loops.0.pid_loop_name":       name,
			"loops.0.pid_loop_type":       "0",
			"loops.0.status":              CHECKSET,
		}
	}

	var fakeBrainIndustrialPidLoopsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"loops.#": "0",
		}
	}

	var BrainIndustrialPidLoopsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existBrainIndustrialPidLoopsMapFunc,
		fakeMapFunc:  fakeBrainIndustrialPidLoopsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
	}

	BrainIndustrialPidLoopsInfo.dataSourceTestCheckWithPreCheck(t, 0, preCheck, nameRegexConf, idsConf, statusConf)
}

func dataSourceBrainIndustrialPidLoopsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_brain_industrial_pid_organization" "default" {
		pid_organization_name = "%[1]s"
	}
	resource "alicloud_brain_industrial_pid_project" "default" {
		pid_organization_id = alicloud_brain_industrial_pid_organization.default.id
		pid_project_name = "%[1]s"
	}
	resource "alicloud_brain_industrial_pid_loop" "default" {
	  pid_loop_dcs_type = "standard"
	  pid_loop_desc = "Test For Terraform"
	  pid_loop_is_crucial = false
	  pid_loop_name = "%[1]s"
	  pid_loop_type = "0"
	  pid_project_id = alicloud_brain_industrial_pid_project.default.id
	  pid_loop_configuration = "{\"baseParam\":{\"forwardController\":false,\"integral\":false,\"kd\":{\"tagValue\":\"20\"},\"kp\":{},\"op\":\"PIDBenchmark.FOPDT_OP\",\"opParam\":{\"increment\":{\"max\":10},\"operate\":{\"max\":115,\"min\":-15},\"range\":{\"max\":115,\"min\":-15},\"trend\":0},\"openLoopTime\":150,\"pv\":\"PIDBenchmark.FOPDT_PV\",\"pvRange\":{\"max\":100,\"min\":0},\"sampleTime\":5,\"sp\":\"PIDBenchmark.FOPDT_SP\",\"spOperate\":{\"max\":100,\"min\":0},\"splitRangeControl\":false,\"suitCtrlTime\":100,\"td\":{},\"ti\":{}},\"identParam\":{\"delay\":10,\"modelType\":3},\"resetParam\":{\"ctrlMode\":0,\"ctrlStuc\":1}}"
	}
`, name)
}
