---
subcategory: "PAIPlugin"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_plugin_schedule"
sidebar_current: "docs-alicloud-resource-pai-plugin-schedule"
description: |-
  Provides a PAI Plugin Schedule resource.
---

# alicloud_pai_plugin_schedule

Provides a PAI Plugin Schedule resource.

For information about PAI Plugin Schedule and how to use it, see [What is PAI Plugin Schedule](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiplugin-2022-01-12-createschedule).

-> **NOTE:** Available since v1.232.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_pai_plugin_schedule" "example" {
  name              = "terraform-example"
  group_id          = "0a54e195-03e2-40bd-869d-b71cb302783e"
  signature_id      = "0a54e195-03e2-40bd-869d-b71cb302783e"
  template_id       = "0a54e195-03e2-40bd-869d-b71cb302783e"
  payment_type      = "PayAsYouGo"
  execute_time      = "2099-01-01 12:00:00"
  repeat_cycle      = 1
  repeat_cycle_unit = 2
  repeat_times      = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the schedule.
* `group_id` - (Required, ForceNew) The crowd (audience) ID.
* `signature_id` - (Optional, ForceNew) The signature ID, or specify a signature.
* `template_id` - (Optional, ForceNew) The template ID, or specify a template code.
* `template_code` - (Optional, ForceNew) The template code.
* `sign_name` - (Optional, ForceNew) The signature name.
* `payment_type` - (Optional, ForceNew) The billing type. Valid values: `PayAsYouGo` (postpaid) and `Subscription` (prepaid). Defaults to postpaid when empty.
* `execute_time` - (Optional, ForceNew) The execution time (UTC+8), in the format `2020-01-01 12:00:00`. Empty means execute immediately.
* `end_time` - (Optional, ForceNew) The end time (UTC+8), in the format `2020-01-01`.
* `ai_send_start_date` - (Optional, ForceNew) The AI send start date, e.g. `2022-01-01 12:00:00`. Write-only: accepted on create but not returned by the API.
* `ai_send_end_date` - (Optional, ForceNew) The AI send end date, e.g. `2022-02-01 12:00:00`. Write-only: accepted on create but not returned by the API.
* `repeat_cycle` - (Optional, ForceNew) The repeat cycle, executed together with the repeat cycle unit.
* `repeat_cycle_unit` - (Optional, ForceNew) The repeat cycle unit. Valid values: `0` (never, default), `1` (hour), `2` (day), `3` (week), `4` (month).
* `repeat_times` - (Optional, ForceNew) The repeat times. `-1` means repeat with no end time; a positive `N` means terminate after `N` repeats.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PAI Plugin Schedule.
* `status` - The status of the schedule. `0`: checking, `1`: check success, `2`: check failed, `3`: sending, `4`: send success, `5`: send failed.
* `create_time` - The creation time (UTC+8).
* `updated_time` - The update time (UTC+8).

## Import

PAI Plugin Schedule can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_plugin_schedule.example <schedule_id>
```
