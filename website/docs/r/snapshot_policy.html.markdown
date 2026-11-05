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

-> **NOTE:** Available since v1.42.0.

## Example Usage

```terraform
resource "alicloud_snapshot_policy" "sp" {
  name            = "tf-demoAcc-sp"
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_snapshot_policy&spm=docs.r.snapshot_policy.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `auto_snapshot_policy_name` - (Optional, Computed) The name of the snapshot policy.
* `name` - (Optional, Deprecated) The snapshot policy name. It has been deprecated from provider version 1.236.0, and use field `auto_snapshot_policy_name` instead.
* `association_type` - (Optional, ForceNew) The association type of the automatic snapshot policy.
* `copied_snapshots_retention_days` - (Optional, Computed) The retention period of cross-region copied snapshots. Unit: days.
* `copy_encryption_configuration` - (Optional) The encryption configuration for cross-region copy. See [`copy_encryption_configuration`](#copy_encryption_configuration) below.
* `enable_cross_region_copy` - (Optional) Specifies whether to enable cross-region copy.
* `repeat_weekdays` - (Required) The automatic snapshot repetition dates. The unit of measurement is day and the repeating cycle is a week. Value range: [1, 7], which represents days starting from Monday to Sunday, for example 1  indicates Monday. When you want to schedule multiple automatic snapshot tasks for a disk in a week, you can set the RepeatWeekdays to an array.
  - A maximum of seven time points can be selected.
  - The format is  an JSON array of ["1", "2", … "7"]  and the time points are separated by commas (,).
* `resource_group_id` - (Optional) The ID of the resource group to which the snapshot policy belongs.
* `retention_days` - (Required) The snapshot retention time, and the unit of measurement is day. Optional values:
  - -1: The automatic snapshots are retained permanently.
  - [1, 65536]: The number of days retained.

    Default value: -1.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `target_copy_regions` - (Optional) The target regions for cross-region copy.
* `target_tags` - (Optional) The tags to be added to the snapshot. See [`target_tags`](#target_tags) below.
* `time_points` - (Required) The automatic snapshot creation schedule, and the unit of measurement is hour. Value range: [0, 23], which represents from 00:00 to 24:00,  for example 1 indicates 01:00. When you want to schedule multiple automatic snapshot tasks for a disk in a day, you can set the TimePoints to an array.
  - A maximum of 24 time points can be selected.
  - The format is  an JSON array of ["0", "1", … "23"] and the time points are separated by commas (,).

### `copy_encryption_configuration`

The copy_encryption_configuration supports the following:

* `encrypted` - (Optional) Whether to encrypt the snapshot.
* `kms_key_id` - (Optional) The KMS key ID used for encryption.

### `target_tags`

The target_tags supports the following:

* `tag_key` - (Optional) The key of the tag.
* `tag_value` - (Optional) The value of the tag.

## Attributes Reference

The following attributes are exported:

* `id` - The snapshot policy ID.
* `create_time` - The creation time of the snapshot policy.
* `region_id` - The region ID of the snapshot policy.
* `status` - The status of the snapshot policy.

## Import

Snapshot can be imported using the id, e.g.

```shell
$ terraform import alicloud_snapshot_policy.snapshot sp-abc1234567890000
```
