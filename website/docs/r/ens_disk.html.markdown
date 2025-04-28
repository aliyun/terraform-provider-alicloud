---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_disk"
description: |-
  Provides a Alicloud ENS Disk resource.
---

# alicloud_ens_disk

Provides a ENS Disk resource.

The disk. When you use it for the first time, please contact the product classmates to add a resource whitelist.

For information about ENS Disk and how to use it, see [What is Disk](https://www.alibabacloud.com/help/en/ens/developer-reference/api-ens-2017-11-10-createdisk).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_disk&exampleId=7fe7f716-126b-df67-1588-d35a2318e7e1db18acca&activeTab=example&spm=docs.r.ens_disk.0.7fe7f71612&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_disk" "default" {
  category      = "cloud_ssd"
  size          = "20"
  payment_type  = "PayAsYouGo"
  ens_region_id = "cn-chongqing-11"
}
```

## Argument Reference

The following arguments are supported:
* `category` - (Required, ForceNew) The category of the disk. Valid values: `cloud_efficiency` (high-efficiency cloud disk), `cloud_ssd` (full Flash cloud disk), `local_hdd` (local HDD), `local_ssd` (local ssd).
* `disk_name` - (Optional) The name of the disk.
* `encrypted` - (Optional, ForceNew) Specifies whether to encrypt the new system disk. Valid values: `true`, `false`(default).
* `ens_region_id` - (Required, ForceNew) The ID of the edge node.
* `kms_key_id` - (Optional, ForceNew) The ID of the KMS key used by the cloud disk. If `encrypted` is set to `true`, the service default key is used when KMSKeyId is empty.
* `payment_type` - (Required, ForceNew) The billing method of the instance. Valid values: `PayAsYouGo`.
* `size` - (Optional, Int) The size of the disk instance. Unit: GiB.
* `snapshot_id` - (Optional, ForceNew) The ID of the snapshot used to create the cloud disk.

The SnapshotId and Size parameters have the following limitations:
  - If the snapshot capacity corresponding to the `snapshot_id` parameter is greater than the specified `size` parameter, the Size of the cloud disk created is the Size of the specified snapshot.
  - If the snapshot capacity corresponding to the `snapshot_id` parameter is less than the set `size` parameter value, the Size of the cloud disk created is the specified `size` parameter value.
* `tags` - (Optional, Map, Available since v1.248.0) The label to which the instance is bound.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the disk was created.
* `status` - The status of the disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 7 mins) Used when create the Disk.
* `delete` - (Defaults to 5 mins) Used when delete the Disk.
* `update` - (Defaults to 5 mins) Used when update the Disk.

## Import

ENS Disk can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_disk.example <id>
```
