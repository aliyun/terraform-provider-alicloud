package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsEvaluationTasksDataSource_nameRegex(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEvaluationTasksDataSourceName(rand, map[string]string{
			"workspace":       `"${alicloud_cms_workspace.default.workspace_name}"`,
			"task_name_regex": `"${alicloud_cms_evaluation_task.default.task_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEvaluationTasksDataSourceName(rand, map[string]string{
			"workspace":       `"${alicloud_cms_workspace.default.workspace_name}"`,
			"task_name_regex": `"${alicloud_cms_evaluation_task.default.task_name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEvaluationTasksDataSourceName(rand, map[string]string{
			"workspace": `"${alicloud_cms_workspace.default.workspace_name}"`,
			"ids":       `["${alicloud_cms_evaluation_task.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEvaluationTasksDataSourceName(rand, map[string]string{
			"workspace": `"${alicloud_cms_workspace.default.workspace_name}"`,
			"ids":       `["${alicloud_cms_evaluation_task.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsEvaluationTasksDataSourceName(rand, map[string]string{
			"workspace":       `"${alicloud_cms_workspace.default.workspace_name}"`,
			"task_name_regex": `"${alicloud_cms_evaluation_task.default.task_name}"`,
			"ids":             `["${alicloud_cms_evaluation_task.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsEvaluationTasksDataSourceName(rand, map[string]string{
			"workspace":       `"${alicloud_cms_workspace.default.workspace_name}"`,
			"task_name_regex": `"${alicloud_cms_evaluation_task.default.task_name}_fake"`,
			"ids":             `["${alicloud_cms_evaluation_task.default.id}_fake"]`,
		}),
	}
	var existAlicloudCmsEvaluationTasksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"evaluation_tasks.#":             "1",
			"evaluation_tasks.0.id":          CHECKSET,
			"evaluation_tasks.0.task_name":   CHECKSET,
			"evaluation_tasks.0.workspace":   CHECKSET,
			"evaluation_tasks.0.task_id":     CHECKSET,
			"evaluation_tasks.0.create_time": CHECKSET,
			"evaluation_tasks.0.status":      CHECKSET,
		}
	}
	var fakeAlicloudCmsEvaluationTasksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "0",
			"evaluation_tasks.#": "0",
		}
	}
	var alicloudCmsEvaluationTasksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_evaluation_tasks.default",
		existMapFunc: existAlicloudCmsEvaluationTasksDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsEvaluationTasksDataSourceNameMapFunc,
	}
	alicloudCmsEvaluationTasksCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCmsEvaluationTasksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tftestacccmsevaltask%d"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = var.name
}

resource "alicloud_cms_evaluation_task" "default" {
  workspace = alicloud_cms_workspace.default.workspace_name
  task_name = var.name
  task_mode = "Manual"
  data_type = "Metric"
  description = var.name
  status = "Running"
  evaluators {
    name = var.name
    result_name = var.name
    variable_mapping = "{}"
  }
}

data "alicloud_cms_evaluation_tasks" "default" {
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
