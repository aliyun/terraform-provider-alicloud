---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_referer"
description: |-
  Provides a Alicloud OSS Bucket Referer resource.
---

# alicloud_oss_bucket_referer

Provides a OSS Bucket Referer resource. Bucket Referer configuration (Hotlink protection).

For information about OSS Bucket Referer and how to use it, see [What is Bucket Referer](https://www.alibabacloud.com/help/en/oss/user-guide/hotlink-protection).

-> **NOTE:** Available since v1.220.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_referer&exampleId=10504655-19f8-f423-2b79-1ffa857bf9c00d496cfb&activeTab=example&spm=docs.r.oss_bucket_referer.0.1050465519&intl_lang=EN_US" target="_blank">
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
* `allow_truncate_query_string` - (Optional, Computed) Specifies whether to truncate the query string in the URL when the Referer is matched. Valid values: true, false.
* `bucket` - (Required, ForceNew) Name of the Bucket.
* `referer_blacklist` - (Optional) The container that holds the Referer blacklist.
* `referer_list` - (Optional) The container that holds the Referer whitelist.
* `truncate_path` - (Optional) Specifies whether to truncate the path and parts that follow the path in the URL when the Referer is matched. Valid values: true, false. If TruncatePath is set to true, the value of AllowTruncateQueryString must be also true because the query string follows the path component. When the path is truncated, the query string is also truncated.

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