---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_oncall_schedules"
description: |-
  Provides a list of Cms OncallSchedule owned by an Alibaba Cloud account.
---

# alicloud_cms_oncall_schedules

This data source provides the Cms OncallSchedules of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.277.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_oncall_schedules" "default" {
  oncall_schedule_name_regex = alicloud_cms_oncall_schedule.default.oncall_schedule_name
  ids                        = [alicloud_cms_oncall_schedule.default.id]
}

output "schedule_id" {
  value = data.alicloud_cms_oncall_schedules.default.schedules.0.oncall_schedule_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of OncallSchedule IDs.
* `oncall_schedule_name_regex` - (Optional) A regex string to filter results by the oncall schedule name.
* `source` - (Optional) The source of the oncall schedule.
* `workspace` - (Optional) The workspace of the oncall schedule.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported:

* `ids` - A list of OncallSchedule IDs.
* `schedules` - A list of OncallSchedules. Each element contains the following attributes:
  * `oncall_schedule_id` - The ID of the oncall schedule.
  * `oncall_schedule_name` - The name of the oncall schedule.
  * `shift_robot_id` - The ID of the shift robot.
  * `rotations` - The rotation configuration list of the oncall schedule.
