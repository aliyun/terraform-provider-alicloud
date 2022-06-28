---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_auto_snapshot_policies"
sidebar_current: "docs-alicloud-datasource-ecs-auto-snapshot-policies"
description: |-
  Provides a list of Ecs Auto Snapshot Policies to the user.
---

# alicloud\_ecs\_auto\_snapshot\_policies

This data source provides the Ecs Auto Snapshot Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.117.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_auto_snapshot_policies" "example" {
  ids        = ["sp-bp14e66xxxxxxxx"]
  name_regex = "tf-testAcc"
}

output "first_ecs_auto_snapshot_policy_id" {
  value = data.alicloud_ecs_auto_snapshot_policies.example.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Auto Snapshot Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Auto Snapshot Policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of Auto Snapshot Policy. Valid Values: `Expire`, `Normal`.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Auto Snapshot Policy names.
* `policies` - A list of Ecs Auto Snapshot Policies. Each element contains the following attributes:
	* `auto_snapshot_policy_id` - The ID of the Auto Snapshot Policy.
	* `copied_snapshots_retention_days` - The retention period of the snapshot copied across regions.
	* `disk_nums` - The number of disks to which the automatic snapshot policy is applied.
	* `enable_cross_region_copy` - Specifies whether to enable the system to automatically copy snapshots across regions.
	* `id` - The ID of the Auto Snapshot Policy.
	* `name` - The snapshot policy name..
	* `repeat_weekdays` - The automatic snapshot repetition dates.
	* `retention_days` - The snapshot retention time, and the unit of measurement is day.
	* `status` - The status of Auto Snapshot Policy.
	* `tags` - A mapping of tags to assign to the resource.
	* `target_copy_regions` - The destination region to which the snapshot is copied.
	* `time_points` - The automatic snapshot creation schedule, and the unit of measurement is hour.
	* `volume_nums` - The number of extended volumes on which this policy is enabled.
