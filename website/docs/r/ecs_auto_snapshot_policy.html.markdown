---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_auto_snapshot_policy"
description: |-
  Provides a Alicloud ECS Auto Snapshot Policy resource.
---

# alicloud_ecs_auto_snapshot_policy

Provides a ECS Auto Snapshot Policy resource. Automatic snapshot policy.

For information about ECS Auto Snapshot Policy and how to use it, see [What is Auto Snapshot Policy](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_resource_manager_resource_group" "ResourceGroup" {
  display_name        = "test"
  resource_group_name = var.name
}


resource "alicloud_ecs_auto_snapshot_policy" "default" {
  time_points               = ["1"]
  resource_group_id         = alicloud_resource_manager_resource_group.ResourceGroup.id
  retention_days            = 1
  repeat_weekdays           = ["1"]
  auto_snapshot_policy_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `auto_snapshot_policy_name` - (Optional) AutoSnapshotPolicyName.
* `copied_snapshots_retention_days` - (Optional, Available since v1.117.0) CopiedSnapshotsRetentionDays.
* `enable_cross_region_copy` - (Optional, Available since v1.117.0) EnableCrossRegionCopy.
* `repeat_weekdays` - (Required, Available since v1.117.0) RepeatWeekdays.
* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the resource group.
* `retention_days` - (Required, Available since v1.117.0) RetentionDays.
* `tags` - (Optional, Map, Available since v1.117.0) Tags.
* `target_copy_regions` - (Optional, Available since v1.117.0) TargetCopyRegions.
* `time_points` - (Required, Available since v1.117.0) TimePoints.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - CreationTime.
* `status` - Status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Auto Snapshot Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Auto Snapshot Policy.
* `update` - (Defaults to 5 mins) Used when update the Auto Snapshot Policy.

## Import

ECS Auto Snapshot Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_auto_snapshot_policy.example <id>
```