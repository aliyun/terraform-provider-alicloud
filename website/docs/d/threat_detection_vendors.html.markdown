---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_vendors"
sidebar_current: "docs-alicloud-datasource-threat-detection-vendors"
description: |-
  Provides a list of Threat Detection Vendors to the user.
---

# alicloud_threat_detection_vendors

This data source provides a list of Threat Detection Vendors in an Alibaba Cloud account.

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_vendors" "default" {
  ids        = ["example-vendor-id"]
  name_regex = "my-vendor"
}

output "first_vendor_id" {
  value = data.alicloud_threat_detection_vendors.default.vendors.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Vendor IDs.
* `name_regex` - (Optional) A regex string to filter results by the vendor name.
* `vendor_type` - (Optional) The type of the vendor.
* `lang` - (Optional) The language type of the API request. Valid values: `en`, `zh`. Default value: `en`.
* `role_for` - (Optional) The user ID when the administrator switches other views.
* `output_file` - (Optional) File path where results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of Vendor IDs.
* `vendors` - A list of Vendors. Each element contains the following attributes:
  * `id` - The ID of the Vendor. It is the same as `vendor_id`.
  * `vendor_id` - The ID of the vendor.
  * `vendor_name` - The name of the vendor.
  * `vendor_type` - The type of the vendor.
  * `create_time` - The creation time of the resource.
  * `update_time` - The last modification time of the resource.
