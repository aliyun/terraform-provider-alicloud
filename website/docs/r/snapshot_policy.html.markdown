---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_snapshot_policy"
sidebar_current: "docs-alicloud-resource-snapshot-policy"
description: |-
  Provides an ECS snapshot policy resource.
---

# alicloud\_snapshot\_policy

-> **DEPRECATED:** This resource has been renamed to [alicloud_ecs_auto_snapshot_policy](https://www.terraform.io/docs/providers/alicloud/r/ecs_auto_snapshot_policy) from version 1.117.0.

Provides an ECS snapshot policy resource.

For information about snapshot policy and how to use it, see [Snapshot](https://www.alibabacloud.com/help/doc-detail/25460.html).

-> **NOTE:** Available in 1.42.0+.

## Example Usage

```
resource "alicloud_snapshot_policy" "sp" {
  name            = "tf-testAcc-sp"
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The snapshot policy name.
* `repeat_weekdays` - (Required) The automatic snapshot repetition dates. The unit of measurement is day and the repeating cycle is a week. Value range: [1, 7], which represents days starting from Monday to Sunday, for example 1  indicates Monday. When you want to schedule multiple automatic snapshot tasks for a disk in a week, you can set the RepeatWeekdays to an array.
    - A maximum of seven time points can be selected.
    - The format is  an JSON array of ["1", "2", … "7"]  and the time points are separated by commas (,).
* `retention_days` - (Required) The snapshot retention time, and the unit of measurement is day. Optional values:
    - -1: The automatic snapshots are retained permanently.
    - [1, 65536]: The number of days retained.
    
    Default value: -1.
* `time_points` - (Required) The automatic snapshot creation schedule, and the unit of measurement is hour. Value range: [0, 23], which represents from 00:00 to 24:00,  for example 1 indicates 01:00. When you want to schedule multiple automatic snapshot tasks for a disk in a day, you can set the TimePoints to an array.
    - A maximum of 24 time points can be selected.
    - The format is  an JSON array of ["0", "1", … "23"] and the time points are separated by commas (,).
    
## Attributes Reference

The following attributes are exported:

* `id` - The snapshot policy ID.

## Import

Snapshot can be imported using the id, e.g.

```
$ terraform import alicloud_snapshot_policy.snapshot sp-abc1234567890000
```
