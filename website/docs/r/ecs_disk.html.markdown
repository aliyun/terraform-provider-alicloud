---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_disk"
sidebar_current: "docs-alicloud-resource-ecs-disk"
description: |-
  Provides a Alicloud ECS Disk resource.
---

# alicloud_ecs_disk

Provides an ECS Disk resource.

For information about ECS Disk and how to use it, see [What is Disk](https://www.alibabacloud.com/help/en/doc-detail/25513.htm).

-> **NOTE:** Available since v1.122.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_disk&exampleId=d536ac89-4604-be47-91a8-9a65cc0b8245cb6ba4a8&activeTab=example&spm=docs.r.ecs_disk.0.d536ac8946&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_kms_key" "example" {
  description            = "terraform-example"
  pending_window_in_days = "7"
  status                 = "Enabled"
}

resource "alicloud_ecs_disk" "example" {
  zone_id     = data.alicloud_zones.example.zones.0.id
  disk_name   = "terraform-example"
  description = "terraform-example"
  category    = "cloud_efficiency"
  size        = "30"
  encrypted   = true
  kms_key_id  = alicloud_kms_key.example.id
  tags = {
    Name = "terraform-example"
  }
}
```

### Deleting `alicloud_ecs_disk` or removing it from your configuration

The `alicloud_ecs_disk` resource allows you to manage `payment_type = "Subscription"` and `delete_with_instance = true` disk, 
but Terraform cannot destroy it. Deleting the subscription resource or removing it from your configuration will 
remove it from your state file and management, but will not destroy it.
If you want to delete it, you can change it to `PayAsYouGo` and setting `delete_with_instance = true` and detach it from instance.

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional, ForceNew) ID of the free zone to which the disk belongs. One of the `zone_id` and `instance_id` must be set but can not be set at the same time.
* `availability_zone` - (Deprecated since v1.122.0) Field `availability_zone` has been deprecated from provider version 1.122.0. New field `zone_id` instead.
* `category` - (Optional) Category of the disk. Default value: `cloud_efficiency`. Valid Values: `cloud`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud_auto`, `cloud_essd_entry`, `elastic_ephemeral_disk_standard`, `elastic_ephemeral_disk_premium`.
* `delete_auto_snapshot` - (Optional) Indicates whether the automatic snapshot is deleted when the disk is released. Default value: `false`.
* `delete_with_instance` - (Optional) Indicates whether the disk is released together with the instance. Default value: `false`.
* `description` - (Optional) Description of the disk. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `disk_name` - (Optional) Name of the ECS disk. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with `http://` or `https://`. Default value is `null`.
* `dry_run` - (Optional) Specifies whether to check the validity of the request without actually making the request.request Default value: false. Valid values:
  - `true`: The validity of the request is checked, but the request is not made. Check items include the required parameters, request format, service limits, and available ECS resources. If the check fails, the corresponding error message is returned. If the check succeeds, the DryRunOperation error code is returned.
  - `false`: The validity of the request is checked. If the check succeeds, a 2xx HTTP status code is returned and the request is made.
* `enable_auto_snapshot` - (Optional) Indicates whether to enable creating snapshot automatically.
* `encrypted` - (Optional, ForceNew) Specifies whether to encrypt the disk. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `instance_id` - (Optional, ForceNew) The ID of the instance to which the created subscription disk is automatically attached.
  * After you specify the instance ID, the specified `resource_group_id`, `tags`, and `kms_key_id` parameters are ignored.
  * One of the `zone_id` and `instance_id` must be set but can not be set at the same time.
* `kms_key_id` - (Optional, ForceNew) The ID of the KMS key corresponding to the data disk, The specified parameter `Encrypted` must be `true` when KmsKeyId is not empty.
* `name` - (Optional, Deprecated since v1.122.0) Field `name` has been deprecated from provider version 1.122.0. New field `disk_name` instead.
* `payment_type` - (Optional) Payment method for disk. Valid values: `PayAsYouGo`, `Subscription`. Default to `PayAsYouGo`. If you want to change the disk payment type, the `instance_id` is required.
* `performance_level` - (Optional) Specifies the performance level of an ESSD when you create the ESSD. Valid values:                                                       
  - `PL0`: A single ESSD delivers up to 10,000 random read/write IOPS.
  - `PL1`: A single ESSD delivers up to 50,000 random read/write IOPS.
  - `PL2`: A single ESSD delivers up to 100,000 random read/write IOPS.
  - `PL3`: A single ESSD delivers up to 1,000,000 random read/write IOPS.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional) The Id of resource group which the disk belongs. This attribute only supports adding or updating, not destroying.
* `size` - (Optional) The size of the disk. Unit: GiB. This parameter is required. Valid values:
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
* `type` - (Optional, Available since v1.122.0) The type to expand cloud disks. Valid Values: `online`, `offline`. Default to `offline`.
  - `offline`: After you resize a disk offline, you must restart the instance by using the console or by calling the RebootInstance operation for the resizing operation to take effect. For more information, see Restart the instance and RebootInstance.
  - `online`: After you resize a disk online, the resizing operation takes effect immediately and you do not need to restart the instance. You can resize ultra disks, standard SSDs, and ESSDs online.

-> **NOTE:** Disk category `cloud` has been outdated, and it only can be used none I/O Optimized ECS instances. Recommend `cloud_efficiency` and `cloud_ssd` disk.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Disk.
* `status` - The disk status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Disk.
* `update` - (Defaults to 10 mins) Used when update the Disk.
* `delete` - (Defaults to 2 mins) Used when delete the Disk.

## Import

ECS Disk can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_disk.example d-abcd12345
```
