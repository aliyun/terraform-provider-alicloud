---
subcategory: "Aligreen"
layout: "alicloud"
page_title: "Alicloud: alicloud_aligreen_image_lib"
description: |-
  Provides a Alicloud Aligreen Image Lib resource.
---

# alicloud_aligreen_image_lib

Provides a Aligreen Image Lib resource.

Image library for image detection.

For information about Aligreen Image Lib and how to use it, see [What is Image Lib](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_aligreen_image_lib&exampleId=0fb744f2-779a-1078-0bfc-b16647566ab5eac5cc1b&activeTab=example&spm=docs.r.aligreen_image_lib.0.0fb744f277&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_aligreen_biz_type" "defaultUalunB" {
  biz_type_name = var.name
}


resource "alicloud_aligreen_image_lib" "default" {
  category       = "BLACK"
  enable         = true
  scene          = "PORN"
  image_lib_name = var.name
  biz_types      = [alicloud_aligreen_biz_type.defaultUalunB.biz_type_name]
}
```

## Argument Reference

The following arguments are supported:
* `biz_types` - (Optional) List of business scenarios. For example: ["bizTypeA", "bizTypeB", "bizTypeC"]
* `category` - (Required, ForceNew) The category of the image library. Valid values: BLACK: a blacklist, WHITE: a whitelist, REVIEW: a review list
* `enable` - (Optional, Computed) Specifies whether to enable the image library. Valid values: true: Enable the image library. This is the default value. false: Disable the image library.
* `image_lib_name` - (Required) The name of the image library defined by the customer. It can contain no more than 20 characters in Chinese, English, and underscore (_).
* `scene` - (Required, ForceNew) The moderation scenario to which the custom image library applies. Valid values: PORN: pornography detection, AD: ad detection, ILLEGAL: terrorist content detection

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Image Lib.
* `delete` - (Defaults to 5 mins) Used when delete the Image Lib.
* `update` - (Defaults to 5 mins) Used when update the Image Lib.

## Import

Aligreen Image Lib can be imported using the id, e.g.

```shell
$ terraform import alicloud_aligreen_image_lib.example <id>
```