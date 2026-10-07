---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_oncall_schedule"
description: |-
  Provides a Alicloud Cms OncallSchedule resource.
---

# alicloud_cms_oncall_schedule

Provides a Cms OncallSchedule resource.

For information about Cms OncallSchedule and how to use it, see [What is OncallSchedule](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateOncallSchedule).

-> **NOTE:** Available since v1.277.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cms_oncall_schedule" "default" {
  oncall_schedule_name = var.name
  shift_robot_id       = var.name
  source               = var.name
  substitudes          = [var.name]
  rotations {
    rotation_name              = var.name
    contacts                   = [var.name]
    rotation_start_time        = "2026-01-01 00:00:00"
    rotation_end_time          = "2026-01-01 23:59:59"
    shift_recurrence_frequency = "Daily"
    start_date                 = "2026-01-01"
    time_zone                  = "Asia/Shanghai"
    shift_length               = 1
    active_days                = [1, 2, 3, 4, 5]
  }
}
```

## Argument Reference

The following arguments are supported:

* `oncall_schedule_name` - (Required) The name of the oncall schedule.
* `rotations` - (Optional) The rotation configuration list of the oncall schedule. See [Block rotations](#block-rotations) below.
* `shift_robot_id` - (Optional) The ID of the shift robot.
* `source` - (Optional) The source of the oncall schedule.
* `substitudes` - (Optional) The substitute contact list of the oncall schedule.

#### Block rotations

The rotations block supports:

* `active_days` - (Optional) The active days of the rotation.
* `contacts` - (Optional) The contact list of the rotation.
* `rotation_end_time` - (Optional) The end time of the rotation.
* `rotation_name` - (Optional) The name of the rotation.
* `rotation_start_time` - (Optional) The start time of the rotation.
* `shift_length` - (Optional) The shift length of the rotation.
* `shift_recurrence_frequency` - (Optional) The recurrence frequency of the shift.
* `start_date` - (Optional) The start date of the rotation.
* `time_zone` - (Optional) The time zone of the rotation.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above. It is the value of `oncall_schedule_id`.
* `oncall_schedule_id` - The ID of the oncall schedule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the OncallSchedule.
* `delete` - (Defaults to 5 mins) Used when delete the OncallSchedule.
* `update` - (Defaults to 5 mins) Used when update the OncallSchedule.

## Import

Cms OncallSchedule can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_oncall_schedule.example <oncall_schedule_id>
```
