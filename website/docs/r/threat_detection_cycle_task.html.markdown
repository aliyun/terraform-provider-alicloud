---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_cycle_task"
description: |-
  Provides a Alicloud Threat Detection Cycle Task resource.
---

# alicloud_threat_detection_cycle_task

Provides a Threat Detection Cycle Task resource.

Configure periodic tasks in Security Center.

For information about Threat Detection Cycle Task and how to use it, see [What is Cycle Task](https://next.api.alibabacloud.com/document/Sas/2018-12-03/CreateCycleTask).

-> **NOTE:** Available since v1.253.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_threat_detection_cycle_task" "default" {
  target_end_time   = "6"
  task_type         = "VIRUS_VUL_SCHEDULE_SCAN"
  target_start_time = "0"
  source            = "console_batch"
  task_name         = "VIRUS_VUL_SCHEDULE_SCAN"
  first_date_str    = "1650556800000"
  period_unit       = "day"
  interval_period   = "7"
  param             = jsonencode({ "targetInfo" : [{ "type" : "groupId", "name" : "TI HOST", "target" : 10597 }, { "type" : "groupId", "name" : "expense HOST", "target" : 10597 }] })
  enable            = "1"
}
```

## Argument Reference

The following arguments are supported:
* `enable` - (Required, Int) Whether to enable. Value:
  - `1`: On
  - `0`: Closed
* `first_date_str` - (Required, Int) First execution time.
* `interval_period` - (Required, Int) Interval period.
* `param` - (Optional) Extended information field.
* `period_unit` - (Required) Unit of scan cycle, value:
  - `day`: day.
  - `hour`: hours.
* `source` - (Optional) Added the source of the task.
* `target_end_time` - (Required, Int) Task end time (hours).
* `target_start_time` - (Required, Int) Task start time (hours).
* `task_name` - (Required, ForceNew) The task name.
  - **VIRUS_VUL_SCHEDULE_SCAN**: scans for viruses.
  - **IMAGE_SCAN**: Image scan.
  - **EMG_VUL_SCHEDULE_SCAN**: Emergency vulnerability scanning.
* `task_type` - (Required, ForceNew) The task type.
  - **VIRUS_VUL_SCHEDULE_SCAN**: scans for viruses.
  - **IMAGE_SCAN**: Image scan.
  - **EMG_VUL_SCHEDULE_SCAN**: Emergency vulnerability scanning.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cycle Task.
* `delete` - (Defaults to 5 mins) Used when delete the Cycle Task.
* `update` - (Defaults to 5 mins) Used when update the Cycle Task.

## Import

Threat Detection Cycle Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_cycle_task.example <id>
```