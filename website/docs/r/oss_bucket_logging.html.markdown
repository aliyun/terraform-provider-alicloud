---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_logging"
description: |-
  Provides a Alicloud OSS Bucket Logging resource.
---

# alicloud_oss_bucket_logging

Provides a OSS Bucket Logging resource.

After you enable and configure logging for a bucket, Object Storage Service (OSS) generates log objects based on a predefined naming convention. This way, access logs are generated and stored in the specified bucket on an hourly basis.

For information about OSS Bucket Logging and how to use it, see [What is Bucket Logging](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutBucketLogging).

-> **NOTE:** Available since v1.222.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_logging&exampleId=cc51cd4a-2444-5f8e-220e-c632e98582a7eaaee009&activeTab=example&spm=docs.r.oss_bucket_logging.0.cc51cd4a24&intl_lang=EN_US" target="_blank">
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

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = "resource-example-logging-806"
}

resource "alicloud_oss_bucket" "CreateLoggingBucket" {
  storage_class = "Standard"
  bucket        = "resource-example-logging-153"
  lifecycle {
    # When you use `alicloud_oss_bucket_logging`, you must explicitly ignore the `logging` attribute on the `alicloud_oss_bucket`.
    ignore_changes = [logging]
  }
}


resource "alicloud_oss_bucket_logging" "default" {
  bucket        = alicloud_oss_bucket.CreateBucket.id
  target_bucket = alicloud_oss_bucket.CreateBucket.id
  target_prefix = "log/"
  logging_role  = "example-role"
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket
* `logging_role` - (Optional, Available since v1.262.0) Authorization role used for bucket logging
* `target_bucket` - (Required) The bucket that stores access logs.
* `target_prefix` - (Optional) The prefix of the saved log objects. This element can be left empty.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Logging.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Logging.
* `update` - (Defaults to 5 mins) Used when update the Bucket Logging.

## Import

OSS Bucket Logging can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_logging.example <id>
```