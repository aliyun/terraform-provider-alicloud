---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_referer"
description: |-
  Provides a Alicloud OSS Bucket Referer resource.
---

# alicloud_oss_bucket_referer

Provides a OSS Bucket Referer resource. Bucket Referer configuration.

For information about OSS Bucket Referer and how to use it, see [What is Bucket Referer](https://www.alibabacloud.com/help/en/oss/user-guide/hotlink-protection).

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
      referer_config,
    ]
  }
}


resource "alicloud_oss_bucket_referer" "default" {
  allow_empty_referer = "true"
  referer_blacklist = [
    "*.forbidden.com"
  ]
  bucket                      = alicloud_oss_bucket.CreateBucket.bucket
  truncate_path               = "false"
  allow_truncate_query_string = "true"
  referer_list = [
    "*.aliyun.com",
    "*.example.com"
  ]
}
```

## Argument Reference

The following arguments are supported:
* `allow_empty_referer` - (Required) Whether to allow empty Referer request headers.
* `allow_truncate_query_string` - (Optional, Computed) Allow phase request parameters.
* `bucket` - (Required, ForceNew) Name of the Bucket.
* `referer_blacklist` - (Optional) The container that holds the Referer blacklist.
* `referer_list` - (Optional) The container that holds the Referer whitelist.
* `truncate_path` - (Optional) Name of the bucket.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Referer.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Referer.
* `update` - (Defaults to 5 mins) Used when update the Bucket Referer.

## Import

OSS Bucket Referer can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_referer.example <id>
```