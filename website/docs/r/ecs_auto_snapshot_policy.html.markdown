---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_auto_snapshot_policy"
description: |-
  Provides a Alicloud ECS Auto Snapshot Policy resource.
---

# alicloud_ecs_auto_snapshot_policy

Provides a ECS Auto Snapshot Policy resource.

For information about ECS Auto Snapshot Policy and how to use it, see [What is Auto Snapshot Policy](https://www.alibabacloud.com/help/en/doc-detail/25527.htm).

-> **NOTE:** Available since v1.117.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_auto_snapshot_policy&exampleId=6763b317-e12a-3836-8904-15dd975ed483827ebcb0&activeTab=example&spm=docs.r.ecs_auto_snapshot_policy.0.6763b317e1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ecs_auto_snapshot_policy" "example" {
  name            = "terraform-example"
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_auto_snapshot_policy&spm=docs.r.ecs_auto_snapshot_policy.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `association_type` - (Optional, Computed, ForceNew, Available since v1.293.0) The association type between the automatic snapshot policy and target resources. Valid values:
  - `AssociatedWithDisk`: The automatic snapshot policy is applied to disks.
  - `AssociatedWithInstanceTag`: The automatic snapshot policy is associated with ECS instances that have the specified `target_tags`. Changing this parameter creates a new resource.

-> **NOTE:** Before you can set `association_type` to `AssociatedWithInstanceTag`, you must contact the ECS product team to add your Alibaba Cloud account to the `AssociatedWithInstanceTag` whitelist. If your account is not on the whitelist, `CreateAutoSnapshotPolicy` returns the `InvalidOperation.UserNotInWhiteList` error. You can submit a ticket to apply for the whitelist.

* `auto_snapshot_policy_name` - (Optional, Available since v1.236.0) The name of the automatic snapshot policy. The name must be 2 to 128 characters in length. The name must start with a letter and cannot start with http:// or https://. The name can contain letters, digits, colons (:), underscores (_), and hyphens (-).
* `copied_snapshots_retention_days` - (Optional, Int) The retention period of the snapshot copy in the destination region. Unit: days. Valid values:
  - `-1`: The snapshot copy is retained until it is deleted.
  - `1` to `65535`: The snapshot copy is retained for the specified number of days. After the retention period of the snapshot copy expires, the snapshot copy is automatically deleted.
* `copy_encryption_configuration` - (Optional, Set, Available since v1.236.0) The encryption parameters for cross-region snapshot replication. See [`copy_encryption_configuration`](#copy_encryption_configuration) below.
* `enable_cross_region_copy` - (Optional, Bool) Specifies whether to enable cross-region replication for snapshots. Valid values: `true`, `false`.
* `repeat_weekdays` - (Required, List) The days of the week on which to create automatic snapshots. Valid values: `1` to `7`, which correspond to the days of the week. For example, `1` indicates Monday. One or more days can be specified.
* `resource_group_id` - (Optional, Available since v1.236.0) The ID of the resource group. If this parameter is specified to query resources, up to 1,000 resources that belong to the specified resource group can be displayed in the response.
* `retention_days` - (Required, Int) The retention period of the automatic snapshots. Unit: days. Valid values:
  - `-1`: Automatic snapshots are retained until they are deleted.
  - `1` to `65536`: Auto snapshots are retained for the specified number of days. After the retention period of auto snapshots expires, the auto snapshots are automatically deleted.
* `tags` - (Optional, Map) A mapping of tags to assign to the resource.
* `target_copy_regions` - (Optional, List) The destination region to which to copy the snapshot. You can specify only a single destination region.
* `target_tags` - (Optional, Set, Available since v1.293.0) The tags used to associate the automatic snapshot policy with ECS instances. This parameter takes effect only when `association_type` is set to `AssociatedWithInstanceTag`. See [`target_tags`](#target_tags) below.
* `time_points` - (Required, List) The points in time of the day at which to create automatic snapshots.

  The time is displayed in UTC+8. Unit: hours. Valid values: `0` to `23`, which correspond to the 24 points in time on the hour from 00:00:00 to 23:00:00. For example, 1 indicates 01:00:00. Multiple points in time can be specified.

  The parameter value is a JSON array that contains up to 24 points in time separated by commas (,). Example: ["0", "1", ... "23"].

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.236.0). Field `name` has been deprecated from provider version 1.236.0. New field `auto_snapshot_policy_name` instead.

### `copy_encryption_configuration`

The copy_encryption_configuration supports the following:
* `encrypted` - (Optional, Bool) Whether to enable encryption for cross-region snapshot replication. Default value: `false`. Valid values: `true`, `false`.
* `kms_key_id` - (Optional) The ID of the Key Management Service (KMS) key used to encrypt snapshots in cross-region snapshot replication.

### `target_tags`

The target_tags supports the following:
* `tag_key` - (Optional, Available since v1.293.0) The key of target tag N. The tag key cannot be an empty string. The tag key can be up to 128 characters in length and cannot start with `acs:` or `aliyun`. Valid values of N: 1 to 10.
* `tag_value` - (Optional, Available since v1.293.0) The value of target tag N. The tag value can be an empty string. The tag value can be up to 128 characters in length and cannot start with `acs:`.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Auto Snapshot Policy.
* `create_time` - (Available since v1.236.0) The time when the automatic snapshot policy was created. The time follows the ISO 8601 standard in the yyyy-MM-ddThh:mm:ssZ format. The time is displayed in UTC.
* `region_id` - (Available since v1.236.0) The region ID of the automatic snapshot policy.
* `status` - The status of the automatic snapshot policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Auto Snapshot Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Auto Snapshot Policy.
* `update` - (Defaults to 5 mins) Used when update the Auto Snapshot Policy.

## Import

ECS Auto Snapshot Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_auto_snapshot_policy.example <id>
```
