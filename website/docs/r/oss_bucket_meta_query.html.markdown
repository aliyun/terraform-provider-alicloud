---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_meta_query"
description: |-
  Provides a Alicloud OSS Bucket Meta Query resource.
---

# alicloud_oss_bucket_meta_query

Provides a OSS Bucket Meta Query resource. Enables the metadata management feature for a bucket.

For information about OSS Bucket Meta Query and how to use it, see [What is Bucket Meta Query](https://www.alibabacloud.com/help/en/oss/developer-reference/openmetaquery).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_meta_query&exampleId=316b92c6-d53e-76ee-57a2-1004246bb90eca2990fa&activeTab=example&spm=docs.r.oss_bucket_meta_query.0.316b92c6d5&intl_lang=EN_US" target="_blank">
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


resource "alicloud_oss_bucket_meta_query" "default" {
  bucket = alicloud_oss_bucket.CreateBucket.bucket
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) The name of the bucket.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the metadata index database. The format is mm:ss + TIMEZONE in the YYYY-MM-DDTHH format of RFC 3339. Where YYYY-MM-DD indicates the year, month and day, T indicates the beginning of the time element, HH:mm:ss indicates the hour, minute and second, and TIMEZONE indicates the time zone.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Meta Query.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Meta Query.

## Import

OSS Bucket Meta Query can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_meta_query.example <id>
```