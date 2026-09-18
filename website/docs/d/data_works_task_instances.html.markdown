---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_task_instances"
sidebar_current: "docs-alicloud-datasource-data-works-task-instances"
description: |-
  Provides a list of Data Works Task Instances to the user.
---

# alicloud_data_works_task_instances

This data source provides the Data Works Task Instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.235.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_task_instances" "default" {
  project_env = "PROD"
  project_id  = "xxxx"
  bizdate     = "1743350400000"
  status      = "Success"
}

output "first_task_instance_id" {
  value = data.alicloud_data_works_task_instances.default.instances.0.task_instance_id
}
```

Filter by ids

```terraform
data "alicloud_data_works_task_instances" "filtered" {
  bizdate     = "1743350400000"
  project_id  = "xxxx"
  project_env = "PROD"
  ids         = ["123456", "789012"]
}
```

## Argument Reference

The following arguments are supported:

* `bizdate` - (Required) The data timestamp, in UNIX timestamp format with milliseconds precision. The value is `00:00:00` of the day before the scheduling time of the instance, e.g. `1743350400000`.
* `project_id` - (Required) The ID of the DataWorks workspace.
* `project_env` - (Optional) The environment of the DataWorks workspace. Valid values: `PROD` or `Dev`.
* `ids` - (Optional, Computed) A list of Task Instance IDs. The data source filters the returned list to the instances whose `task_instance_id` matches one of the values in this list.
* `name_regex` - (Optional) A regex string to filter results by `task_name`.
* `owner` - (Optional) The account ID of the task owner.
* `page_number` - (Optional) The page number of the response. Pages start from page `1`. Defaults to `1`.
* `page_size` - (Optional) The number of entries per page. Defaults to `50`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `sort_by` - (Optional) The field used to sort the returned instances, in the format of `<Sort field> <Desc|Asc>`. Valid sort fields: `TriggerTime`, `StartedTime`, `FinishedTime`, `CreateTime`, `Id`. Defaults to `Id Desc`.
* `status` - (Optional) The status of the task instance. Valid values: `NotRun`, `Running`, `Failure`, `Success`, `WaitTime`, `WaitResource`.
* `task_id` - (Optional) The ID of the task for which the instance is generated.
* `task_name` - (Optional) The name of the task. Fuzzy match is supported.
* `task_type` - (Optional) The type of the task for which the instance is generated.
* `trigger_type` - (Optional) The trigger type. Valid values: `Scheduler` (scheduling cycle-based) or `Manual` (manual trigger).
* `workflow_id` - (Optional) The ID of the workflow to which the instance belongs.
* `workflow_instance_id` - (Optional) The workflow instance ID.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Data Works Task Instances. Each element contains the following attributes:
  * `id` - The ID of the Task Instance, same as `task_instance_id`.
  * `task_instance_id` - The unique ID of the Task Instance.
  * `task_type` - The type of the task, e.g. `NODE`, `SQL`, `MANUAL`.
  * `trigger_type` - The trigger type of the task instance, e.g. `SCHEDULED`, `MANUAL`.
  * `bizdate` - The data timestamp of the task instance, in UNIX timestamp format with milliseconds precision.
  * `task_id` - The ID of the node/task.
  * `project_env` - The environment of the DataWorks project, `PROD` or `DEV`.
  * `owner` - The owner of the task.
  * `workflow_instance_id` - The ID of the workflow instance.
  * `project_id` - The ID of the DataWorks project.
  * `workflow_id` - The ID of the workflow.
  * `workflow_instance_type` - The type of the workflow instance, e.g. `MANUAL`, `AUTO`.
  * `runtime_resource_resource_group_id` - The ID of the resource group used at runtime.
  * `task_name` - The name of the task/node.
  * `region_id` - The region ID of the task instance.
  * `description` - The description of the task.
  * `workflow_name` - The name of the workflow.
  * `timeout` - The timeout (in minutes) of the task.
  * `rerun_mode` - The rerun mode of the task, e.g. `Allowed`, `NotAllowed`.
  * `run_number` - The run number of the task instance.
  * `baseline_id` - The ID of the baseline, if any.
  * `priority` - The priority of the task instance.
  * `runtime_resource_image` - The image used at runtime.
  * `runtime_resource_cu` - The compute unit used at runtime.
  * `runtime_process_id` - The process ID at runtime.
  * `runtime_gateway` - The gateway used at runtime.
  * `trigger_time` - The trigger time of the task instance, in RFC 3339 format.
  * `started_time` - The started time of the task instance, in RFC 3339 format.
  * `finished_time` - The finished time of the task instance, in RFC 3339 format.
  * `tenant_id` - The tenant ID.
  * `create_time` - The creation time of the task instance, in RFC 3339 format.
  * `modify_time` - The modification time of the task instance, in RFC 3339 format.
  * `create_user` - The user who created the task instance.
  * `modify_user` - The user who last modified the task instance.
  * `data_source_name` - The name of the data source used by the task.
  * `status` - The status of the task instance, e.g. `NotRun`, `WaitTime`, `WaitResource`, `Running`, `Checking`, `CheckingCondition`, `Failure`, `Success`.
  * `period_number` - The period number of the task instance.
