---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_public_access_block"
description: |-
  Provides a Alicloud OSS Bucket Public Access Block resource.
---

# alicloud_oss_bucket_public_access_block

Provides a OSS Bucket Public Access Block resource. Blocking public access at the bucket-level.

For information about OSS Bucket Public Access Block and how to use it, see [What is Bucket Public Access Block](https://www.alibabacloud.com/help/en/).

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

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
}


resource "alicloud_oss_bucket_public_access_block" "default" {
  bucket              = alicloud_oss_bucket.CreateBucket.bucket
  block_public_access = true
}
```

## Argument Reference

The following arguments are supported:
* `block_public_access` - (Required) Whether AlibabaCloud OSS should block public bucket policies and ACL for this bucket.
* `bucket` - (Required, ForceNew) The name of the bucket.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Public Access Block.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Public Access Block.
* `update` - (Defaults to 5 mins) Used when update the Bucket Public Access Block.

## Import

OSS Bucket Public Access Block can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_public_access_block.example <id>
```