---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_invocations"
sidebar_current: "docs-alicloud-datasource-ecs-invocations"
description: |-
  Provides a list of Ecs Invocations to the user.
---

# alicloud\_ecs\_invocations

This data source provides the Ecs Invocations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.168.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_invocations" "ids" {
  ids = ["example-id"]
}
output "ecs_invocation_id_1" {
  value = data.alicloud_ecs_invocations.ids.invocations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `command_id` - (Optional, ForceNew) The execution ID of the command.
* `content_encoding` - (Optional, ForceNew) The encoding mode of the CommandContent and Output response parameters. Valid values: `PlainText`, `Base64`.
* `ids` - (Optional, ForceNew, Computed)  A list of Invocation IDs.
* `invoke_status` - (Optional, ForceNew) The overall execution state of the command. The value of this parameter depends on the execution states on all the involved instances. Valid values: `Running`, `Finished`, `Failed`, `PartialFailed`, `Stopped`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `invocations` - A list of Ecs Invocations. Each element contains the following attributes:
  * `command_id` - The ID of the command.
  * `create_time` - The creation time of the resource.
  * `frequency` - The schedule on which the recurring execution of the command takes place. For information about the value specifications, see [Cron expression](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/cron-expression).
  * `id` - The ID of the Invocation.
  * `invocation_id` - The ID of the Invocation.
  * `invocation_status` - The overall execution state of the command. The value of this parameter depends on the execution states on all the involved instances.
  * `command_type` - The type of the command.
  * `invoke_status` - The overall execution state of the command. **Note:** We recommend that you ignore this parameter and check the value of the `invocation_status` response parameter for the overall execution state.
  * `command_content` - The Base64-encoded command content.
  * `command_name` - The name of the command.
  * `parameters` - The custom parameters in the command.
  * `repeat_mode` - Indicates the execution mode of the command.
  * `timed` - Indicates whether the commands are to be automatically run.
  * `username` - The username that was used to run the command on the instance.
  * `invoke_instances` - Execute target instance set type.
    * `creation_time` - The start time of the execution.
    * `update_time` - The time when the execution state was updated.
    * `finish_time` - The end time of the execution.
    * `invocation_status` - The execution state on a single instance. Valid values: `Pending`, `Scheduled`, `Running`, `Success`, `Failed`, `Stopping`, `Stopped`, `PartialFailed`.
    * `repeats` - The number of times that the command is run on the instance.
    * `instance_id` - The ID of the instance.
    * `output` - The output of the command.
    * `dropped` - The size of truncated and discarded text when the value of the Output response parameter exceeds 24 KB in size.
    * `stop_time` - The time when the command stopped being run on the instance. If you call the StopInvocation operation to manually stop the execution, the value is the time when you call the operation.
    * `exit_code` - The exit code of the execution.
    * `start_time` - The time when the command started to be run on the instance.
    * `error_info` - Details about the reason why the command failed to be sent or run.
    * `timed` - Indicates whether the commands are to be automatically run.
    * `error_code	` - The code that indicates why the command failed to be sent or run. 
    * `instance_invoke_status	` - **Note:** We recommend that you ignore this parameter and check the value of the `invocation_status` response parameter for the overall execution state.