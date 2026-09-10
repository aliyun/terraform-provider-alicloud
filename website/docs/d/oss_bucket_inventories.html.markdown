---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_inventories"
sidebar_current: "docs-alicloud-datasource-oss-bucket-inventories"
description: |-
  Provides a list of Oss Bucket Inventory owned by an Alibaba Cloud account.
---

# alicloud_oss_bucket_inventories

This data source provides Oss Bucket Inventory available to the user.[What is Bucket Inventory](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketInventory)

-> **NOTE:** Available since v1.284.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
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

data "alicloud_oss_bucket_inventories" "default" {
  ids    = ["${alicloud_oss_bucket_inventory.default.id}"]
  bucket = alicloud_oss_bucket.CreateBucket.id
}

output "alicloud_oss_bucket_inventory_example_id" {
  value = data.alicloud_oss_bucket_inventories.default.inventories.0.id
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `ids` - (Optional, ForceNew, Computed) A list of Bucket Inventory IDs. The value is formulated as `<bucket>:<inventory_id>`.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Bucket Inventory IDs.
* `inventories` - A list of Bucket Inventory Entries. Each element contains the following attributes:
    * `destination` - Holds the container that holds the location of the inventory results.
        * `oss_bucket_destination` - The Bucket information stored after the list result is exported.
            * `account_id` - The account ID granted by the Bucket owner.
            * `bucket` - The Bucket where the exported manifest file is stored.
            * `encryption` - The encryption method of the manifest file.
                * `ssekms` - The container that holds the SSE-KMS encryption key.
                    * `key_id` - KMS key ID.
                * `sseoss` - The container that holds the SSE-OSS encryption method.
            * `format` - The file format of the manifest file.
            * `prefix` - The storage path prefix for the manifest file.
            * `role_arn` - The name of the role that has the permission to read all files in the source Bucket and write files to the target Bucket.
    * `filter` - Container for inventory filtering rules.
        * `last_modify_begin_time_stamp` - The start timestamp of the last modification time of the filter file, in seconds.
        * `last_modify_end_time_stamp` - The end timestamp of the last modification time of the filter file, in seconds.
        * `lower_size_bound` - The minimum size of the filter file, in B.
        * `prefix` - The match prefix of the filter rule.
        * `storage_class` - The storage type of the filter file.
        * `upper_size_bound` - The maximum size of the filter file, in B.
    * `included_object_versions` - Whether the Object version information is included in the list.
    * `incremental_inventory` - Configuration container for incremental inventory.
        * `is_enabled` - Incremental inventory enabled.
        * `optional_fields` - Configuration container for incremental manifest file properties.
            * `field` - Incremental inventory export field list.
        * `schedule` - Incremental inventory export cycle container.
            * `frequency` - Describes the period of incremental inventory file export, in seconds, currently fixed at 10 minutes.
    * `inventory_id` - The ID of the inventory rule.
    * `is_enabled` - Identification of whether the manifest feature is enabled.
    * `optional_fields` - Sets the configuration items included in the manifest results.
        * `field` - The configuration items contained in the manifest results.
    * `schedule` - Container for storing inventory export cycle information.
        * `frequency` - Period for manifest file export.
    * `id` - The ID of the resource supplied above.
