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

```terraform
variable "name" {
  default = "terraform_example"
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
  biz_types      = ["example_007"]
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