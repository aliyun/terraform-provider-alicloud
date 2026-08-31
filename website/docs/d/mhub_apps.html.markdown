---
subcategory: "Enterprise Mobile Application Studio (MHUB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mhub_apps"
sidebar_current: "docs-alicloud-datasource-mhub-apps"
description: |-
  Provides a list of Mhub Apps to the user.
---

# alicloud\_mhub\_apps

This data source provides the Mhub Apps of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.138.0+.

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
  type         = "2"
}
data "alicloud_mhub_apps" "ids" {}
output "mhub_app_id_1" {
  value = data.alicloud_mhub_apps.ids.apps.0.id
}

data "alicloud_mhub_apps" "nameRegex" {
  name_regex = "^my-App"
}
output "mhub_app_id_2" {
  value = data.alicloud_mhub_apps.nameRegex.apps.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of App IDs. The value formats as `<product_id>:<app_key>`
* `product_id` - (Required)  The ID of the Product.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by App name.
* `os_type` - (Optional, ForceNew) The os type. Valid values: `Android` and `iOS`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of App names.
* `apps` - A list of Mhub Apps. Each element contains the following attributes:
	* `app_key` - Application AppKey, which uniquely identifies an application when requested by the interface
	* `app_name` - The Name of the App.
	* `bundle_id` - iOS application ID. Required when creating an iOS app. **NOTE:** Either `bundle_id` or `package_name` must be set.
	* `create_time` - The CreateTime of the App.
	* `encoded_icon` - Base64 string of picture.
	* `id` - The ID of the App.
	* `industry_id` - The Industry ID of the app. For information about Industry and how to use it, MHUB[Industry](https://help.aliyun.com/document_detail/201638.html).
	* `package_name` - Android App package name.  **NOTE:** Either `bundle_id` or `package_name` must be set.
	* `product_id` - The ID of the Product.
	* `type` - The type of the App. Valid values: `Android` and `iOS`. 
