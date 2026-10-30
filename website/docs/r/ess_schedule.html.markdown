---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_schedule"
sidebar_current: "docs-alicloud-resource-ess-schedule"
description: |-
  Provides a ESS schedule resource.
---

# alicloud_ess_schedule

-> **DEPRECATED:** This resource has been deprecated from v1.45.0.

`alicloud_ess_schedule` is a deprecated alias of [`alicloud_ess_scheduled_task`](./ess_scheduled_task.html.markdown) and shares the same implementation factory. The arguments and attributes below mirror `alicloud_ess_scheduled_task`.

## Argument Reference

The following arguments are supported:

* `scheduled_action` - (Optional) The operation to be performed when a scheduled task is triggered. Enter the unique identifier of a scaling rule. **NOTE:** You cannot specify `scheduled_action` and `scaling_group_id` at the same time.
* `scheduled_task_name` - (Optional) Display name of the scheduled task, which must be 2-40 characters (English or Chinese) long.
* `scaling_group_id` - (Optional, Available in 1.90.0+) The ID of the scaling group where the number of instances is modified when the scheduled task is triggered. After the `ScalingGroupId` parameter is specified, the scaling method of the scheduled task is to specify the number of instances in a scaling group. You must specify at least one of the following parameters: `MinValue`, `MaxValue`, and `DesiredCapacity`. **NOTE:** You cannot specify `scheduled_action` and `scaling_group_id` at the same time.
* `min_value` - (Optional, Available in 1.90.0+) The minimum number of instances in a scaling group when the scaling method of the scheduled task is to specify the number of instances in a scaling group.
* `max_value` - (Optional, Available in 1.90.0+) The maximum number of instances in a scaling group when the scaling method of the scheduled task is to specify the number of instances in a scaling group.
* `desired_capacity` - (Optional, Available in 1.90.0+) The expected number of instances in a scaling group when the scaling method of the scheduled task is to specify the number of instances in a scaling group. **NOTE:** You must specify the `DesiredCapacity` parameter when you create the scaling group.
* `description` - (Optional) Description of the scheduled task, which is 2-200 characters (English or Chinese) long.
* `launch_time` - (Optional) The time at which the scheduled task is triggered. Specify the time in the ISO 8601 standard in the YYYY-MM-DDThh:mmZ format.
The time must be in UTC. You cannot enter a time point later than 90 days from the date of scheduled task creation.
If the `recurrence_type` parameter is specified, the task is executed repeatedly at the time specified by LaunchTime.
Otherwise, the task is only executed once at the date and time specified by LaunchTime.
* `launch_expiration_time` - (Optional) The time period during which a failed scheduled task is retried. Unit: seconds. Valid values: 0 to 1800. Default value: 600
* `recurrence_type` - (Optional) Specifies the recurrence type of the scheduled task. **NOTE:** You must specify `RecurrenceType`, `RecurrenceValue`, and `RecurrenceEndTime` at the same time. Valid values:
  - Daily: The scheduled task is executed once every specified number of days.
  - Weekly: The scheduled task is executed on each specified day of a week.
  - Monthly: The scheduled task is executed on each specified day of a month.
  - Cron: (Available in 1.60.0+) The scheduled task is executed based on the specified cron expression.
* `recurrence_value` - (Optional) Specifies how often a scheduled task recurs. **NOTE:** You must specify `RecurrenceType`, `RecurrenceValue`, and `RecurrenceEndTime` at the same time. The valid value depends on `recurrence_type`
  - Daily: You can enter one value. Valid values: 1 to 31.
  - Weekly: You can enter multiple values and separate them with commas (,). For example, the values 0 to 6 correspond to the days of the week in sequence from Sunday to Saturday.
  - Monthly: You can enter two values in A-B format. Valid values of A and B: 1 to 31. The value of B must be greater than or equal to the value of A.
  - Cron: You can enter a cron expression which is written in UTC and consists of five fields: minute, hour, day of month (date), month, and day of week. The expression can contain wildcard characters including commas (,), question marks (?), hyphens (-), asterisks (*), number signs (#), forward slashes (/), and the L and W letters.
* `recurrence_end_time` - (Optional) Specifies the end time after which the scheduled task is no longer repeated. Specify the time in the ISO 8601 standard in the YYYY-MM-DDThh:mm:ssZ format.
The time must be in UTC. You cannot enter a time point later than 365 days from the date of scheduled task creation. **NOTE:** You must specify `RecurrenceType`, `RecurrenceValue`, and `RecurrenceEndTime` at the same time.
* `task_enabled` - (Optional) Specifies whether to start the scheduled task. Default to true.

## Attributes Reference

The following attributes are exported:

* `id` - The schedule task ID.

## Import

ESS schedule can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_schedule.example abc123456
```
