---
subcategory: "OSS"
layout: "alicloud"
page_title: "Alicloud: alicloud_oss_bucket_style"
description: |-
  Provides a Alicloud OSS Bucket Style resource.
---

# alicloud_oss_bucket_style

Provides a OSS Bucket Style resource.

Image styles that contain single or multiple image processing parameters.

For information about OSS Bucket Style and how to use it, see [What is Bucket Style](https://next.api.alibabacloud.com/document/Oss/2019-05-17/PutStyle).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oss_bucket_style&exampleId=7cc355b7-630c-7238-2198-3c3be026163cb25757e8&activeTab=example&spm=docs.r.oss_bucket_style.0.7cc355b763&intl_lang=EN_US" target="_blank">
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

resource "alicloud_oss_bucket_style" "default" {
  bucket     = alicloud_oss_bucket.CreateBucket.id
  style_name = "style-933"
  content    = "image/resize,p_75,w_75"
  category   = "document"
}
```

## Argument Reference

The following arguments are supported:
* `bucket` - (Required, ForceNew) Storage space to which the picture style belongs
* `category` - (Optional, Computed) Style category, valid values: image, document, video.
* `content` - (Required) The Image style content can contain single or multiple image processing parameters.
* `style_name` - (Required, ForceNew) Image Style Name

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<bucket>:<style_name>`.
* `create_time` - Image Style Creation Time

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bucket Style.
* `delete` - (Defaults to 5 mins) Used when delete the Bucket Style.
* `update` - (Defaults to 5 mins) Used when update the Bucket Style.

## Import

OSS Bucket Style can be imported using the id, e.g.

```shell
$ terraform import alicloud_oss_bucket_style.example <bucket>:<style_name>
```