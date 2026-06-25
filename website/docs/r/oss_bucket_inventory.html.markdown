---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_inventory"
description: |-
  Provides a Alicloud OSS Bucket Inventory resource.
---

# alicloud_oss_bucket_inventory

Provides a OSS Bucket Inventory resource. Bucket inventory periodically exports a list of objects and their metadata in a bucket to a destination bucket in CSV, ORC, or Parquet format, for object auditing, storage cost analysis, and compliance.

For information about OSS Bucket Inventory and how to use it, see [What is Bucket Inventory](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketInventory).

-> **NOTE:** Available since v1.283.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_account" "this" {}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-src-${random_integer.default.result}"
}

resource "alicloud_oss_bucket" "DestBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-dst-${random_integer.default.result}"
}

resource "alicloud_ram_role" "role" {
  name     = "${var.name}-${random_integer.default.result}"
  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {"Service": ["oss.aliyuncs.com"]}
      }
    ],
    "Version": "1"
  }
  EOF
}

resource "alicloud_oss_bucket_inventory" "default" {
  bucket                   = alicloud_oss_bucket.CreateBucket.id
  inventory_id             = "report1"
  is_enabled               = true
  included_object_versions = "All"

  schedule {
    frequency = "Daily"
  }

  optional_fields {
    field = ["Size", "LastModifiedDate"]
  }

  filter {
    prefix = "frontends/"
  }

  destination {
    oss_bucket_destination {
      format     = "CSV"
      account_id = data.alicloud_account.this.id
      role_arn   = alicloud_ram_role.role.arn
      bucket     = "acs:oss:::${alicloud_oss_bucket.DestBucket.id}"
      prefix     = "inventory/"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of the source bucket.
* `inventory_id` - (Required, ForceNew) The name of the inventory rule. It must be globally unique within the bucket.
* `is_enabled` - (Required) Whether to enable the bucket inventory. Valid values: `true`, `false`.
* `included_object_versions` - (Required) Whether to include all object versions in the inventory. Valid values: `All`, `Current`.
* `schedule` - (Required) The export schedule of the inventory. See [`schedule`](#schedule) below.
* `destination` - (Required) The destination where the inventory report is exported. See [`destination`](#destination) below.
* `optional_fields` - (Optional) The configuration fields to include in the exported inventory report. See [`optional_fields`](#optional_fields) below.
* `filter` - (Optional) The container that stores filter conditions for inventoried objects. See [`filter`](#filter) below.

### `schedule`

The schedule supports the following:

* `frequency` - (Required) The export frequency of the inventory. Valid values: `Daily`, `Weekly`.

### `optional_fields`

The optional_fields supports the following:

* `field` - (Optional) The configuration fields included in the inventory report. Valid values: `Size`, `LastModifiedDate`, `ETag`, `StorageClass`, `IsMultipartUploaded`, `EncryptionStatus`.

### `filter`

The filter supports the following:

* `prefix` - (Optional) The prefix that is used to filter inventoried objects.

### `destination`

The destination supports the following:

* `oss_bucket_destination` - (Required) The information about the bucket to which the inventory report is exported. See [`oss_bucket_destination`](#destination-oss_bucket_destination) below.

### `destination-oss_bucket_destination`

The oss_bucket_destination supports the following:

* `format` - (Required, ForceNew) The format of the exported inventory report. Valid values: `CSV`, `ORC`, `Parquet`.
* `account_id` - (Required, ForceNew) The ID of the account to which permissions are granted by the destination bucket.
* `role_arn` - (Required, ForceNew) The Alibaba Cloud Resource Name (ARN) of the role that has permission to access the destination bucket.
* `bucket` - (Required, ForceNew) The destination bucket that stores the exported inventory report, in `acs:oss:::bucket_name` format.
* `prefix` - (Optional) The prefix of the path where the inventory report is exported.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bucket Inventory. It formats as `<bucket>:<inventory_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Bucket Inventory.
* `update` - (Defaults to 5 mins) Used when update the Bucket Inventory.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Inventory.

## Import

OSS Bucket Inventory can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_inventory.example <bucket>:<inventory_id>
```
