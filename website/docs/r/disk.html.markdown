---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_disk"
sidebar_current: "docs-alicloud-resource-disk"
description: |-
  Provides a ECS Disk resource.
---

# alicloud\_disk

Provides a ECS disk resource.

-> **DEPRECATED:** This resource has been renamed to [alicloud_ecs_disk](https://www.terraform.io/docs/providers/alicloud/r/ecs_disk) from version 1.122.0.

-> **NOTE:** One of `size` or `snapshot_id` is required when specifying an ECS disk. If all of them be specified, `size` must more than the size of snapshot which `snapshot_id` represents. Currently, `alicloud_disk` doesn't resize disk.

## Example Usage

```terraform
# Create a new ECS disk.
resource "alicloud_disk" "ecs_disk" {
  # cn-beijing
  availability_zone = "cn-beijing-b"
  name              = "New-disk"
  description       = "Hello ecs disk."
  category          = "cloud_efficiency"
  size              = "30"
  encrypted         = true
  kms_key_id        = "2a6767f0-a16c-4679-a60f-13bf*****"
  tags = {
    Name = "TerraformTest"
  }
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_disk&spm=docs.r.disk.example&intl_lang=EN_US)

### Destroying a disk with `force_delete`

The following example creates a pay-as-you-go disk with `force_delete = true`; destroying the resource releases the disk:

```terraform
data "alicloud_zones" "force_delete" {
  available_resource_creation = "VSwitch"
  available_disk_category     = "cloud_essd"
}

resource "alicloud_disk" "force_delete" {
  zone_id      = data.alicloud_zones.force_delete.zones.0.id
  name         = "terraform-example"
  category     = "cloud_essd"
  size         = "40"
  force_delete = true
}
```

### Deleting `alicloud_disk` or removing it from your configuration

By default, destroying an `alicloud_disk` resource or removing it from your configuration only removes the disk from your state file and management, without releasing it, when the disk is PrePaid (`payment_type = "Subscription"`) **and** ECS reports `DeleteWithInstance = true` for it, which is the case for a disk released together with its Subscription instance. This is the only state-only path. Every other case, including a PrePaid disk whose `delete_with_instance` is `false`, calls `DeleteDisk` during destroy; if ECS rejects the release, the destroy fails with the returned error and the disk is kept, rather than being silently removed from the state.

Set `force_delete = true` to release such a disk: during destroy, the provider converts the disk to pay-as-you-go, detaches it with `DeleteWithInstance=false` and releases it once it is available. **`force_delete` deletes the disk and all data on it, and the release cannot be undone.** `force_delete` takes priority over the legacy state-only behavior; when it is `false` or not set, the behavior above is unchanged.

Multi-attach disks (`multi_attach = "Enabled"`) are not supported by `force_delete`: the destroy fails with an actionable error and the disk is kept. The PrePaid to pay-as-you-go conversion requires the disk to be attached to a Subscription instance and is limited to one successful call per disk every five minutes; the refund of the remaining subscription fee is handled by ECS according to your account settings.

If you want to delete a disk without `force_delete`, convert it to `PayAsYouGo` and detach it from its instance before deleting it.

## Argument Reference

The following arguments are supported:
* `advanced_features` - (Optional) The advanced features configured for the disk.
* `bursting_enabled` - (Optional, Bool, Available since v1.237.0) Specifies whether to enable the performance burst feature. Valid values: `true`, `false`. **NOTE:** `bursting_enabled` is only valid when `category` is `cloud_auto`; specifying it for other categories is rejected by the API. When `category` is changed to `cloud_auto` in the same apply (for example from `cloud_essd`), the provider defers the `BurstingEnabled` update until the disk category has been confirmed as `cloud_auto` by `ModifyDiskSpec` and `WaitForState`, because the API rejects `BurstingEnabled` on a non-`cloud_auto` disk.
* `category` - (Optional) The category of the data disk. Default value: `cloud_efficiency`. Valid Values: `cloud`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud_auto`, `cloud_essd_entry`, `elastic_ephemeral_disk_standard`, `elastic_ephemeral_disk_premium`. **NOTE:** When `category` is `cloud_auto`, the `bursting_enabled` and `provisioned_iops` parameters become applicable; they are rejected by the API for other categories. **NOTE:** Category conversion (for example from `cloud_essd` to `cloud_auto`) is applied asynchronously by the API: the disk may report the transitional status `Modifying` while the conversion runs, and even after the status returns to `Available` the category may still report the previous value for some time. The provider therefore waits within a single `DescribeDisks` response until the status is `Available`/`In_use` and the category already equals the requested value before completing the update, and this wait is bounded by the `update` timeout (default 30 mins). Conversion duration depends on the disk (throughput, capacity and the previous category) and has no completion guarantee; set `timeouts.update` explicitly for large disks whose conversion may take longer.
* `delete_auto_snapshot` - (Optional, Bool) Specifies whether to delete the automatic snapshots of the disk when the disk is released. Default value: `false`.
* `delete_with_instance` - (Optional, Bool) Specifies whether to release the disk along with its associated instance. Default value: `false`.
* `description` - (Optional) The description of the disk. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
* `disk_name` - (Optional) The name of the data disk. The name must be 2 to 128 characters in length and can contain letters, digits, colons (:), underscores (_), periods (.), and hyphens (-). The name must start with a letter.
* `dry_run` - (Optional, Bool) Specifies whether to check the validity of the request without actually making the request.request Default value: `false`. Valid values:
  - `true`: The validity of the request is checked, but the request is not made. Check items include the required parameters, request format, service limits, and available ECS resources. If the check fails, the corresponding error message is returned. If the check succeeds, the DryRunOperation error code is returned.
  - `false`: The validity of the request is checked. If the check succeeds, a 2xx HTTP status code is returned and the request is made.
* `enable_auto_snapshot` - (Deprecated, Optional, Bool) Specifies whether the automatic snapshot policy feature is enabled for the cloud disk. Valid values: `true` and `false`. The default value is empty, which indicates that the current value is not changed. **NOTE:** This parameter is deprecated. The automatic snapshot policy feature is enabled by default for a cloud disk after it is created. To use the automatic snapshot policy, apply one to the cloud disk.
* `encrypt_algorithm` - (Optional) The encryption algorithm used to encrypt the disk. **NOTE:** `encrypt_algorithm` is only valid when `encrypted` is `true`.
* `encrypted` - (Optional, ForceNew, Bool) Specifies whether to encrypt the disk. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `force_delete` - (Optional, Bool, Available since v2.1.1) Specifies whether to release the disk during destroy even when the default behavior cannot release it, for example a `payment_type = "Subscription"` disk. Defaults to `false`. When set to `true`, the provider converts the disk to `PayAsYouGo` (`ModifyDiskChargeType`), detaches it with `DeleteWithInstance=false`, waits until the disk is `Available` and then releases it (`DeleteDisk`). **NOTE:** `force_delete` deletes the disk and all data on it, and the release cannot be undone. It takes priority over the legacy behavior that only removes a PrePaid disk with `delete_with_instance = true` from the state. `force_delete` does not support multi-attach disks (`multi_attach = "Enabled"`); the destroy fails with an actionable error and the disk is kept. The billing method conversion requires the disk to be attached to a Subscription instance and is limited to one successful call per disk every five minutes.
* `instance_id` - (Optional, Computed) The ID of the instance to which the created subscription disk is automatically attached.
  * After you specify the instance ID, the specified `resource_group_id`, `tags`, and `kms_key_id` parameters are ignored.
  * One of the `zone_id` and `instance_id` must be set but can not be set at the same time.
* `kms_key_id` - (Optional, ForceNew) The ID of the Key Management Service (KMS) key that is used for the disk. **NOTE:** `kms_key_id` is only valid when `encrypted` is `true`.
* `multi_attach` - (Optional, ForceNew, Available since v1.237.0) Specifies whether to enable the multi-attach feature for the disk. Default value: `Disabled`. Valid values: `Enabled`, `Disabled`. **NOTE:** Currently, `multi_attach` can only be set to `Enabled` when `category` is set to `cloud_essd`.
* `payment_type` - (Optional) The payment type of the disk. Default to `PayAsYouGo`. Valid values: `PayAsYouGo`, `Subscription`. If you want to change the disk payment type, the `instance_id` is required.
* `performance_level` - (Optional) Specifies the performance level of an ESSD when you create the ESSD. Valid values:
  - `PL0`: A single ESSD delivers up to 10,000 random read/write IOPS.
  - `PL1`: A single ESSD delivers up to 50,000 random read/write IOPS.
  - `PL2`: A single ESSD delivers up to 100,000 random read/write IOPS.
  - `PL3`: A single ESSD delivers up to 1,000,000 random read/write IOPS.
* `provisioned_iops` - (Optional, Int, Available since v1.237.0) The provisioned read/write IOPS of the ESSD AutoPL disk. Valid values: 0 to min{50,000, 1,000 × Capacity - Baseline IOPS}. **NOTE:** `provisioned_iops` is only valid when `category` is `cloud_auto`.
* `resource_group_id` - (Optional) The ID of the resource group to which to add the disk.
* `size` - (Optional, Int) The size of the disk. Unit: GiB. This parameter is required. Valid values:
  - If `category` is set to `cloud`. Valid values: `5` to `2000`.
  - If `category` is set to `cloud_efficiency`. Valid values: `20` to `32768`.
  - If `category` is set to `cloud_ssd`. Valid values: `20` to `32768`.
  - If `category` is set to `cloud_auto`. Valid values: `1` to `65536`.
  - If `category` is set to `cloud_essd_entry`. Valid values: `10` to `32768`.
  - If `category` is set to `elastic_ephemeral_disk_standard`. Valid values: `64` to `8192`.
  - If `category` is set to `elastic_ephemeral_disk_premium`. Valid values: `64` to `8192`.
  - If `category` is set to `cloud_essd`, the valid values are related to `performance_level`. Valid values:
    - If `performance_level` is set to `PL0`. Valid values: `1` to `65536`.
    - If `performance_level` is set to `PL1`. Valid values: `20` to `65536`.
    - If `performance_level` is set to `PL2`. Valid values: `461` to `65536`.
    - If `performance_level` is set to `PL3`. Valid values: `1261` to `65536`.
* `snapshot_id` - (Optional, ForceNew) The ID of the snapshot to use to create the disk. **NOTE:** If the size of the snapshot specified by `snapshot_id` is larger than the value of `size`, the size of the created disk is equal to the specified snapshot size. If the size of the snapshot specified by `snapshot_id` is smaller than the value of `size`, the size of the created disk is equal to the value of `size`.
* `storage_set_id` - (Optional, ForceNew) The ID of the storage set.
* `storage_set_partition_number` - (Optional, ForceNew) The number of partitions in the storage set.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `type` - (Optional, Available since v1.122.0) The type to expand cloud disks. Valid Values: `online`, `offline`. Default to `offline`.
  - `offline`: After you resize a disk offline, you must restart the instance by using the console or by calling the RebootInstance operation for the resizing operation to take effect. For more information, see Restart the instance and RebootInstance.
  - `online`: After you resize a disk online, the resizing operation takes effect immediately and you do not need to restart the instance. You can resize ultra disks, standard SSDs, and ESSDs online.
* `zone_id` - (Optional, ForceNew) ID of the free zone to which the disk belongs. One of the `zone_id` and `instance_id` must be set but can not be set at the same time.
* `availability_zone` - (Deprecated since v1.122.0) Field `availability_zone` has been deprecated from provider version 1.122.0. New field `zone_id` instead.
* `name` - (Optional, Deprecated since v1.122.0) Field `name` has been deprecated from provider version 1.122.0. New field `disk_name` instead.

  -> **NOTE:** Disk category `cloud` has been outdated, and it only can be used none I/O Optimized ECS instances. Recommend `cloud_efficiency` and `cloud_ssd` disk.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Disk.
* `create_time` - (Available since v1.237.0) The time when the disk was created.
* `region_id` - (Available since v1.237.0) The ID of the region to which the disk belongs.
* `status` - The status of the disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Disk.
* `delete` - (Defaults to 5 mins) Used when delete the Disk. When `force_delete` is `true`, this timeout applies to each step of the forced sequence: converting the billing method, detaching the disk and waiting until it is available, and releasing it. Adjust it when a step needs more time, for example a slow billing method conversion.
* `update` - (Defaults to 30 mins) Used when update the Disk.

## Import

Cloud disk can be imported using the id, e.g.

```shell
$ terraform import alicloud_disk.example d-abc12345678
```
