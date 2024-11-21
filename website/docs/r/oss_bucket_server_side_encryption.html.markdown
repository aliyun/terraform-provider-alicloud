---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_server_side_encryption"
description: |-
  Provides a Alicloud OSS Bucket Server Side Encryption resource.
---

# alicloud_oss_bucket_server_side_encryption

Provides a OSS Bucket Server Side Encryption resource. Server-side encryption rules of the bucket.

For information about OSS Bucket Server Side Encryption and how to use it, see [What is Bucket Server Side Encryption](https://www.alibabacloud.com/help/en/oss/developer-reference/putbucketencryption).

-> **NOTE:** Available since v1.222.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_server_side_encryption&exampleId=9b44ce5d-da59-85c9-407b-eaae723eaf9f3bb89158&activeTab=example&spm=docs.r.oss_bucket_server_side_encryption.0.9b44ce5dda&intl_lang=EN_US" target="_blank">
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
  lifecycle {
    ignore_changes = [
      server_side_encryption_rule,
    ]
  }
}

resource "alicloud_kms_key" "GetKMS" {
  origin                 = "Aliyun_KMS"
  protection_level       = "SOFTWARE"
  description            = var.name
  key_spec               = "Aliyun_AES_256"
  key_usage              = "ENCRYPT/DECRYPT"
  automatic_rotation     = "Disabled"
  pending_window_in_days = 7
}


resource "alicloud_oss_bucket_server_side_encryption" "default" {
  kms_data_encryption = "SM4"
  kms_master_key_id   = alicloud_kms_key.GetKMS.id
  bucket              = alicloud_oss_bucket.CreateBucket.bucket
  sse_algorithm       = "KMS"
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.
* `kms_data_encryption` - (Optional) The algorithm used to encrypt objects. If this element is not specified, objects are encrypted by using AES256. This element is valid only when the value of SSEAlgorithm is set to KMS.
* `kms_master_key_id` - (Optional) The CMK ID that must be specified when SSEAlgorithm is set to KMS and a specified CMK is used for encryption. In other cases, this element must be set to null.
* `sse_algorithm` - (Required) The server-side encryption method. Valid Values: KMS, AES256.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Server Side Encryption.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Server Side Encryption.
* `update` - (Defaults to 5 mins) Used when update the Bucket Server Side Encryption.

## Import

OSS Bucket Server Side Encryption can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_server_side_encryption.example <id>
```