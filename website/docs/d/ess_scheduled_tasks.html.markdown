---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_scheduled_tasks"
sidebar_current: "docs-alicloud_ess_scheduled_tasks"
description: |-
    Provides a list of scheduled tasks available to the user.
---

# alicloud_ess_scheduled_tasks

This data source provides available scheduled task resources. 

-> **NOTE:** Available in 1.72.0+

## Example Usage

```
data "alicloud_ess_scheduled_tasks" "ds" {
  scheduled_task_id = "scheduled_task_id"
  name_regex       = "scheduled_task_name"
}

output "first_scheduled_task" {
  value = "${data.alicloud_ess_scheduled_tasks.ds.tasks.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `scheduled_task_id` - (Optional) The id of the scheduled task.
* `scheduled_action` - (Optional) The operation to be performed when a scheduled task is triggered.
* `name_regex` - (Optional) A regex string to filter resulting scheduled tasks by name.
* `ids` - (Optional) A list of scheduled task IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scheduled task ids.
* `names` - A list of scheduled task names.
* `tasks` - A list of scheduled tasks. Each element contains the following attributes:
  * `id` - ID of the scheduled task id.
  * `name` - Name of the scheduled task name.
  * `scheduled_action` - The operation to be performed when a scheduled task is triggered.
  * `description` - Description of the scheduled task.
  * `launch_expiration_time` - The time period during which a failed scheduled task is retried.
  * `launch_time` - The time at which the scheduled task is triggered.
  * `recurrence_type` - Specifies the recurrence type of the scheduled task. 
  * `recurrence_value` - Specifies how often a scheduled task recurs. 
  * `recurrence_end_time` - Specifies the end time after which the scheduled task is no longer repeated.
