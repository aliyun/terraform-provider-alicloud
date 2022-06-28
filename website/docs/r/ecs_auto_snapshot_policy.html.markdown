---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_auto_snapshot_policy"
sidebar_current: "docs-alicloud-resource-ecs-auto-snapshot-policy"
description: |-
  Provides a Alicloud ECS Auto Snapshot Policy resource.
---

# alicloud\_ecs\_auto\_snapshot\_policy

Provides a ECS Auto Snapshot Policy resource.

For information about ECS Auto Snapshot Policy and how to use it, see [What is Auto Snapshot Policy](https://www.alibabacloud.com/help/en/doc-detail/25527.htm).

-> **NOTE:** Available in v1.117.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_auto_snapshot_policy" "example" {
  name            = "tf-testAcc"
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
* `target_copy_regions` - (Optional) The destination region to which the snapshot is copied. You can set a destination region.
* `copied_snapshots_retention_days` - (Optional) The retention period of the snapshot copied across regions.
    - -1: The snapshot is permanently retained.
    - [1, 65535]: The automatic snapshot is retained for the specified number of days.     
    Default value: -1.
* `enable_cross_region_copy` - (Optional) Specifies whether to enable the system to automatically copy snapshots across regions.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Auto Snapshot Policy.
* `status` - The status of Auto Snapshot Policy.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Auto Snapshot Policy.
* `delete` - (Defaults to 3 mins) Used when delete the Auto Snapshot Policy.

## Import

ECS Auto Snapshot Policy can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_auto_snapshot_policy.example <id>
```
