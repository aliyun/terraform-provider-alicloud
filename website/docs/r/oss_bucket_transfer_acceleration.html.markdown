---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_transfer_acceleration"
description: |-
  Provides a Alicloud OSS Bucket Transfer Acceleration resource.
---

# alicloud_oss_bucket_transfer_acceleration

Provides a OSS Bucket Transfer Acceleration resource. Transfer acceleration configuration of a bucket.

For information about OSS Bucket Transfer Acceleration and how to use it, see [What is Bucket Transfer Acceleration](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}


resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "${var.name}-${random_integer.default.result}"
  lifecycle {
    ignore_changes = [
      transfer_acceleration,
    ]
  }
}


resource "alicloud_oss_bucket_transfer_acceleration" "default" {
  bucket  = alicloud_oss_bucket.CreateBucket.bucket
  enabled = true
}
```

### Deleting `alicloud_oss_bucket_transfer_acceleration` or removing it from your configuration

Terraform cannot destroy resource `alicloud_oss_bucket_transfer_acceleration`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the Bucket.
* `enabled` - (Optional) Specifies whether to enable transfer acceleration for the bucket. Valid values: true: transfer acceleration for the bucket is enabled. false: transfer acceleration for the bucket is disabled.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Transfer Acceleration.
* `update` - (Defaults to 5 mins) Used when update the Bucket Transfer Acceleration.

## Import

OSS Bucket Transfer Acceleration can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_transfer_acceleration.example <id>
```