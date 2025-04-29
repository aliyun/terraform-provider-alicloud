---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_public_access_block"
description: |-
  Provides a Alicloud OSS Bucket Public Access Block resource.
---

# alicloud_oss_bucket_public_access_block

Provides a OSS Bucket Public Access Block resource. Blocking public access at the bucket-level.

For information about OSS Bucket Public Access Block and how to use it, see [What is Bucket Public Access Block](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketpublicaccessblock).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_public_access_block&exampleId=d27d6567-1f37-30f8-996d-a7f223723ca6cc8c7924&activeTab=example&spm=docs.r.oss_bucket_public_access_block.0.d27d65671f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Public Access Block.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Public Access Block.
* `update` - (Defaults to 5 mins) Used when update the Bucket Public Access Block.

## Import

OSS Bucket Public Access Block can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_public_access_block.example <id>
```