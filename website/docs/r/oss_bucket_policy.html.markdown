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
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_oss_bucket_policy&exampleId=8b9446e6-9cc7-8b46-9c87-b16c607c0bcd0827edd3&activeTab=example&spm=docs.r.oss_bucket_policy.0.8b9446e69c" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

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
  depends_on = [alicloud_oss_bucket.CreateBucket]

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