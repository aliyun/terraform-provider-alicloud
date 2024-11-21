---
subcategory: "Enterprise Mobile Application Studio (MHUB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mhub_app"
sidebar_current: "docs-alicloud-resource-mhub-app"
description: |-
  Provides a Alicloud MHUB App resource.
---

# alicloud_mhub_app

Provides a MHUB App resource.

For information about MHUB App and how to use it, see [What is App](https://help.aliyun.com/product/65109.html).

-> **NOTE:** Available since v1.138.0+.

-> **NOTE:** At present, the resource only supports cn-shanghai region.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mhub_app&exampleId=6ba13854-db17-db2b-fcb6-4123543594605aedb177&activeTab=example&spm=docs.r.mhub_app.0.6ba13854db&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

variable "name" {
  default = "example_value"
}

resource "alicloud_mhub_product" "default" {
  product_name = var.name
}

resource "alicloud_mhub_app" "default" {
  app_name     = var.name
  product_id   = alicloud_mhub_product.default.id
  package_name = "com.example.android"
  type         = "Android"
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Required) AppName.
* `bundle_id` - (Optional) The app id of iOS. **NOTE:** Either `bundle_id` or `package_name` must be set.
* `encoded_icon` - (Optional) Base64 string of picture.
* `product_id` - (Required)  The ID of the Product.
* `type` - (Required) The type of the Product. Valid values: `Android` and `iOS`.  
* `industry_id` - (Optional) The Industry ID of the app. For information about Industry and how to use it, MHUB[Industry](https://help.aliyun.com/document_detail/201638.html).
* `package_name` - (Optional) Android App package name. **NOTE:** Either `bundle_id` or `package_name` must be set.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of App. The value formats as `<product_id>:<app_key>`

## Import

MHUB App can be imported using the id, e.g.

```shell
$ terraform import alicloud_mhub_app.example <product_id>:<app_key>
```
