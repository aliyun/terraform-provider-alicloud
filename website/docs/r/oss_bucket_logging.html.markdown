---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_logging"
description: |-
  Provides a Alicloud OSS Bucket Logging resource.
---

# alicloud_oss_bucket_logging

Provides a OSS Bucket Logging resource. After you enable and configure logging for a bucket, Object Storage Service (OSS) generates log objects based on a predefined naming convention. This way, access logs are generated and stored in the specified bucket on an hourly basis.

For information about OSS Bucket Logging and how to use it, see [What is Bucket Logging](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketlogging).

-> **NOTE:** Available since v1.222.0.

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
      logging,
    ]
  }
}

resource "alicloud_oss_bucket_logging" "default" {
  bucket        = alicloud_oss_bucket.CreateBucket.bucket
  target_bucket = alicloud_oss_bucket.CreateBucket.bucket
  target_prefix = "log/"
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `target_bucket` - (Required) The bucket that stores access logs.
* `target_prefix` - (Optional) The prefix of the saved log objects. This element can be left empty.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Logging.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Logging.
* `update` - (Defaults to 5 mins) Used when update the Bucket Logging.

## Import

OSS Bucket Logging can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_logging.example <id>
```