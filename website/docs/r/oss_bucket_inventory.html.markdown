---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_inventory"
description: |-
  Provides a Alicloud OSS Bucket Inventory resource.
---

# alicloud_oss_bucket_inventory

Provides a OSS Bucket Inventory resource.

The inventory rule of an OSS bucket. Inventory periodically exports object metadata from a bucket.

For information about OSS Bucket Inventory and how to use it, see [What is Bucket Inventory](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketInventory).

-> **NOTE:** Available since v1.284.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
}


resource "alicloud_oss_bucket_inventory" "default" {
  destination {
  }
  optional_fields {
    field = ["Size", "LastModifiedDate", "ETag"]
  }
  bucket = alicloud_oss_bucket.CreateBucket.id
  filter {
    prefix           = "Pics/"
    lower_size_bound = "256"
    upper_size_bound = "999999"
    storage_class    = "Standard"
  }
  included_object_versions = "Current"
  schedule {
    frequency = "Daily"
  }
  inventory_id = "report01"
  is_enabled   = false
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `destination` - (Optional, Set) Holds the container that holds the location of the inventory results. **Note:** The parameter is immutable after resource creation and the exported storage location cannot be modified without recreating the rule. See [`destination`](#destination) below.
* `filter` - (Optional, Set) Container for inventory filtering rules. See [`filter`](#filter) below.
* `included_object_versions` - (Optional) Whether the Object version information is included in the list. Valid values: All: export All version information of the Object. Current: exports the Current version of the Object.
* `incremental_inventory` - (Optional, Set) Configuration container for incremental inventory. **Note:** The parameter is immutable after resource creation; OSS fixes the incremental export cycle server-side. See [`incremental_inventory`](#incremental_inventory) below.
* `inventory_id` - (Required, ForceNew) The ID of the inventory rule. The ID must be unique in the bucket.
* `is_enabled` - (Optional) Identification of whether the manifest feature is enabled. Valid values: true: Enable the inventory feature. false: Do not enable the manifest feature.
* `optional_fields` - (Optional, Set) Sets the configuration items included in the manifest results. See [`optional_fields`](#optional_fields) below.
* `schedule` - (Optional, Set) Container for storing inventory export cycle information. See [`schedule`](#schedule) below.

### `destination`

The destination supports the following:
* `oss_bucket_destination` - (Optional, Set) The Bucket information stored after the list result is exported. See [`oss_bucket_destination`](#destination-oss_bucket_destination) below.

### `destination-oss_bucket_destination`

The destination-oss_bucket_destination supports the following:
* `account_id` - (Optional) The account ID granted by the Bucket owner.
* `bucket` - (Optional) The Bucket where the exported manifest file is stored.
* `encryption` - (Optional, Set) The encryption method of the manifest file. Valid value: SSE-OSS: Use the OSS fully managed key for encryption and decryption. SSE-KMS: Use the default KMS-managed CMK(Customer Master Key) or a specified CMK for encryption and decryption. See [`encryption`](#destination-oss_bucket_destination-encryption) below.
* `format` - (Optional) The file format of the manifest file.
* `prefix` - (Optional) The storage path prefix for the manifest file.
* `role_arn` - (Optional) The name of the role that has the permission to read all files in the source Bucket and write files to the target Bucket. The format is acs:ram::uid:role/rolename.

### `destination-oss_bucket_destination-encryption`

The destination-oss_bucket_destination-encryption supports the following:
* `ssekms` - (Optional, Set) The container that holds the SSE-KMS encryption key. See [`ssekms`](#destination-oss_bucket_destination-encryption-ssekms) below.
* `sseoss` - (Optional) The container that holds the SSE-OSS encryption method. Set it to an empty string when OSS-managed keys are used.

### `destination-oss_bucket_destination-encryption-ssekms`

The destination-oss_bucket_destination-encryption-ssekms supports the following:
* `key_id` - (Optional) KMS key ID.

### `filter`

The filter supports the following:
* `last_modify_begin_time_stamp` - (Optional, Int) The start timestamp of the last modification time of the filter file, in seconds. Value range:[1262275200, 253402271999]
* `last_modify_end_time_stamp` - (Optional, Int) The end timestamp of the last modification time of the filter file, in seconds. Value range:[1262275200, 253402271999]
* `lower_size_bound` - (Optional, Int) The minimum size of the filter file, in B. Value range: greater than or equal to 0 B, less than or equal to 48.8 TB.
* `prefix` - (Optional) The match prefix of the filter rule.
* `storage_class` - (Optional) The storage type of the filter file. Multiple storage types can be specified. Optional values: Standard: Standard storage IA: low-frequency access Archive: Archive storage ColdArchive: cold Archive storage All (default): All storage types
* `upper_size_bound` - (Optional, Int) The maximum size of the filter file, in B. Value range: greater than 0 B, less than or equal to 48.8 TB.

### `incremental_inventory`

The incremental_inventory supports the following:
* `is_enabled` - (Optional) Incremental inventory enabled
* `optional_fields` - (Optional, Set) Configuration container for incremental manifest file properties See [`optional_fields`](#incremental_inventory-optional_fields) below.
* `schedule` - (Optional, Set) Incremental inventory export cycle container See [`schedule`](#incremental_inventory-schedule) below.

### `incremental_inventory-optional_fields`

The incremental_inventory-optional_fields supports the following:
* `field` - (Optional, List) Incremental inventory export field list

### `incremental_inventory-schedule`

The incremental_inventory-schedule supports the following:
* `frequency` - (Optional, Int) Describes the period of incremental inventory file export, in seconds, currently fixed at 10 minutes.

### `optional_fields`

The optional_fields supports the following:
* `field` - (Optional, List) The configuration items contained in the manifest results.

### `schedule`

The schedule supports the following:
* `frequency` - (Optional) Period for manifest file export.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<bucket>:<inventory_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Inventory.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Inventory.
* `update` - (Defaults to 5 mins) Used when update the Bucket Inventory.

## Import

OSS Bucket Inventory can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_inventory.example <bucket>:<inventory_id>
```