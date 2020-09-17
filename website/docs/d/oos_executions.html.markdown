---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_executions"
sidebar_current: "docs-alicloud-datasource-oos-executions"
description: |-
    Provides a list of OOS Executions.
---

# alicloud\_oos\_executions

This data source provides a list of OOS Executions in an Alibaba Cloud account according to the specified filters.
 
-> **NOTE:** Available in v1.93.0+.

## Example Usage

```
# Declare the data source

data "alicloud_oos_executions" "example" {
  ids = ["execution_id"]
  template_name = "name"
  status = "Success"
}

output "first_execution_id" {
  value = "${data.alicloud_oos_executions.example.executions.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of OOS Execution ids.
* `category` - (Optional) The category of template. Valid: `AlarmTrigger`, `EventTrigger`, `Other` and `TimerTrigger`.
* `end_date` - (Optional) The time when the execution was ended.
* `end_date_after` - (Optional) Execution whose end time is less than or equal to the specified time.
* `executed_by` - (Optional) The user who execute the template.
* `include_child_execution` - (Optional) Whether to include sub-execution.
* `mode` - (Optional) The mode of OOS Execution. Valid: `Automatic`, `Debug`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `parent_execution_id` - (Optional) The id of parent OOS Execution.
* `ram_role` - (Optional) The role that executes the current template.
* `sort_field` - (Optional) The sort field.
* `sort_order` - (Optional) The sort order.
* `start_date_after` - (Optional) The execution whose start time is greater than or equal to the specified time.
* `start_date_before` - (Optional) The execution with start time less than or equal to the specified time.
* `status` - (Optional) The Status of OOS Execution. Valid: `Cancelled`, `Failed`, `Queued`, `Running`, `Started`, `Success`, `Waiting`.
* `template_name` - (Optional) The name of execution template.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of OOS Execution ids.
* `executions` - A list of OOS Executions. Each element contains the following attributes:
  * `id` - ID of the OOS Executions.
  * `counters` - The counters of OOS Execution.
  * `create_date` - The time when the execution was created.
  * `execution_id` - ID of the OOS Executions.
  * `is_parent` - Whether to include subtasks.
  * `outputs` - The outputs of OOS Executions.
  * `parameters` - The parameters required by the template
  * `start_date` - The time when the template was started.
  * `status_message` - The message of status.
  * `status_reason` - The reason of status.
  * `template_id` - The id of execution template.
  * `template_version` - The version of execution template.
  * `update_date` - The time when the template was updated.
