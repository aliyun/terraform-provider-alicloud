---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_evaluation_task"
description: |-
  Provides a Alicloud Cms Evaluation Task resource.
---

# alicloud_cms_evaluation_task

Provides a Cms Evaluation Task resource.

For information about Cms Evaluation Task and how to use it, see [What is Evaluation Task](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateEvaluationTask).

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
  data_filter = var.name
  channel     = var.name
  description = var.name
  status      = "Running"
  config      = jsonencode({ key = "value" })
  tags        = jsonencode({ env = "test" })
  evaluators {
    name             = var.name
    result_name      = var.name
    result_type      = "Metric"
    config           = jsonencode({ scope = "all" })
    filters          = jsonencode({ region = "cn-hangzhou" })
    variable_mapping = jsonencode({ metric = "cpu" })
  }
}
```

## Argument Reference

The following arguments are supported:
* `channel` - (Optional) The channel source of the evaluation task.
* `config` - (Optional) The task configuration, in JSON format.
* `data_filter` - (Optional) The data filter of the evaluation task.
* `data_type` - (Optional, ForceNew) The data type of the evaluation task.
* `description` - (Optional) The description of the evaluation task.
* `evaluators` - (Optional) The evaluator configuration list. See [`evaluators`](#evaluators) below.
* `run_strategies` - (Optional) The run strategies of the evaluation task.
* `status` - (Optional) The status of the evaluation task. Valid values: `Running`, `Pendding`, `Completed`, `Failed`.
* `tags` - (Optional) The attribute tags of the evaluation task, in JSON format.
* `task_mode` - (Optional, ForceNew) The task mode of the evaluation task.
* `task_name` - (Required, ForceNew) The name of the evaluation task.
* `workspace` - (Required, ForceNew) The workspace to which the evaluation task belongs.

### `evaluators`

The `evaluators` block supports:
* `config` - (Optional) The evaluator configuration, in JSON format.
* `filters` - (Optional) The evaluator filters, in JSON format.
* `name` - (Required) The evaluator name.
* `result_name` - (Required) The metric result name of the evaluator.
* `result_type` - (Optional) The result type of the evaluator.
* `variable_mapping` - (Required) The variable mapping of the evaluator, in JSON format.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It is formatted as `<workspace>:<task_id>`.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.
* `task_id` - The task ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Evaluation Task.
* `delete` - (Defaults to 5 mins) Used when delete the Evaluation Task.
* `update` - (Defaults to 5 mins) Used when update the Evaluation Task.

## Import

Cms Evaluation Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_evaluation_task.example <workspace>:<task_id>
```
