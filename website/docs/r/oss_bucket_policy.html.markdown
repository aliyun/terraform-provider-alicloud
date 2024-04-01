---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_policy"
description: |-
  Provides a Alicloud OSS Bucket Policy resource.
---

# alicloud_oss_bucket_policy

Provides a OSS Bucket Policy resource.  Authorization policy of a bucket.

For information about OSS Bucket Policy and how to use it, see [What is Bucket Policy](https://www.alibabacloud.com/help/en/oss/user-guide/use-bucket-policy-to-grant-permission-to-access-oss).

-> **NOTE:** Available since v1.220.0.

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
      policy,
    ]
  }
}

resource "alicloud_oss_bucket_policy" "default" {
  policy = jsonencode({ "Version" : "1", "Statement" : [{ "Action" : ["oss:PutObject", "oss:GetObject"], "Effect" : "Deny", "Principal" : ["1234567890"], "Resource" : ["acs:oss:*:1234567890:*/*"] }] })
  bucket = alicloud_oss_bucket.CreateBucket.bucket
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the Bucket.
* `policy` - (Required) Json-formatted authorization policies for buckets.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Policy.
* `update` - (Defaults to 5 mins) Used when update the Bucket Policy.

## Import

OSS Bucket Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_policy.example <id>
```