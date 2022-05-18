---
subcategory: "Enterprise Mobile Application Studio (MHUB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mhub_app"
sidebar_current: "docs-alicloud-resource-mhub-app"
description: |-
  Provides a Alicloud MHUB App resource.
---

# alicloud\_mhub\_app

Provides a MHUB App resource.

For information about MHUB App and how to use it, see [What is App](https://help.aliyun.com/product/65109.html).

-> **NOTE:** Available in v1.138.0+.

-> **NOTE:** At present, the resource only supports cn-shanghai region.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example_value"
}
resource "alicloud_mhub_app" "default" {
  app_name     = var.name
  product_id   = alicloud_mhub_product.default.id
  package_name = "com.test.android"
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

```
$ terraform import alicloud_mhub_app.example <product_id>:<app_key>
```
