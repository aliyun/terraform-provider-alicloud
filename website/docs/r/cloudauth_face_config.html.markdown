---
subcategory: "Cloudauth"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloudauth_face_config"
sidebar_current: "docs-alicloud-resource-cloudauth-face-config"
description: |-
  Provides a Alicloud Cloudauth Face Config resource.
---

# alicloud_cloudauth_face_config

Provides a Cloudauth Face Config resource.

For information about Cloudauth Face Config and how to use it, see [What is Face Config](https://help.aliyun.com/zh/id-verification/cloudauth/product-overview/end-of-integration-announcement-on-id-verification).

-> **NOTE:** Available since v1.137.0.

-> **NOTE:** In order to provide you with more perfect product capabilities, the real person certification service has stopped access, it is recommended that you use the upgraded version of the [real person certification financial real person certification service](https://help.aliyun.com/zh/id-verification/product-overview/what-is-id-verification-for-financial-services). Users that have access to real person authentication are not affected.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloudauth_face_config&exampleId=7668b7e4-1823-81e1-b8aa-0e3c777150af2b7fee6c&activeTab=example&spm=docs.r.cloudauth_face_config.0.7668b7e418&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
provider "alicloud" {
  region = "cn-hangzhou"
}
resource "random_integer" "default" {
  max = 99999
  min = 10000
}
resource "alicloud_cloudauth_face_config" "example" {
  biz_name = format("%s-biz", var.name)
  biz_type = format("type-%s", random_integer.default.result)
}
```

## Argument Reference

The following arguments are supported:

* `biz_name` - (Required) Scene name.
* `biz_type` - (Required) Scene type. **NOTE:** The biz_type cannot exceed 32 characters and can only use English letters, numbers and dashes (-).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Face Config. Its value is same as `biz_type`.
* `gmt_modified` - Last Modified Date.

## Import

Cloudauth Face Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloudauth_face_config.example <lang>
```
