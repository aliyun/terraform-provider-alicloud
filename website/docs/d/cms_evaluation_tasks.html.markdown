---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_evaluation_tasks"
sidebar_current: "docs-alicloud-datasource-cms-evaluation-tasks"
description: |-
  Provides a list of Cms Evaluation Tasks to the user.
---

# alicloud\_cms\_evaluation\_tasks

This data source provides the Cms Evaluation Tasks of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = var.name
}

resource "alicloud_cms_evaluation_task" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  task_name   = var.name
  task_mode   = "Manual"
  data_type   = "Metric"
  description = var.name
  status      = "Running"
  evaluators {
    name             = var.name
    result_name      = var.name
    variable_mapping = "{}"
  }
}

data "alicloud_cms_evaluation_tasks" "default" {
  workspace = alicloud_cms_workspace.default.workspace_name
  ids       = [alicloud_cms_evaluation_task.default.id]
}
output "cms_evaluation_task_id_1" {
  value = data.alicloud_cms_evaluation_tasks.default.evaluation_tasks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `channel` - (Optional, ForceNew) The channel source used to filter results.
* `ids` - (Optional, ForceNew, Computed) A list of Evaluation Task IDs. Its element value is formatted as `<workspace>:<task_id>`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status used to filter results. Valid values: `Running`, `Pendding`, `Completed`, `Failed`.
* `task_mode` - (Optional, ForceNew) The task mode used to filter results.
* `task_name_regex` - (Optional, ForceNew) A regex string to filter results by Evaluation Task name.
* `workspace` - (Required, ForceNew) The workspace to which the evaluation tasks belong.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `evaluation_tasks` - A list of Cms Evaluation Tasks. Each element contains the following attributes:
  * `channel` - The channel source of the evaluation task.
  * `create_time` - The creation time of the resource.
  * `data_filter` - The data filter of the evaluation task.
  * `data_type` - The data type of the evaluation task.
  * `description` - The description of the evaluation task.
  * `evaluators` - The evaluator configuration of the evaluation task.
  * `id` - The ID of the resource. It is formatted as `<workspace>:<task_id>`.
  * `run_strategies` - The run strategies of the evaluation task.
  * `status` - The status of the evaluation task.
  * `task_id` - The task ID of the resource.
  * `task_mode` - The task mode of the evaluation task.
  * `task_name` - The name of the evaluation task.
  * `workspace` - The workspace to which the evaluation task belongs.
