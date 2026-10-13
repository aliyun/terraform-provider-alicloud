---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_vendor"
sidebar_current: "docs-alicloud-resource-threat-detection-vendor"
description: |-
  Provides a Alicloud Threat Detection Vendor resource.
---

# alicloud_threat_detection_vendor

Provides a Threat Detection Vendor resource.

For information about Threat Detection Vendor and how to use it, see [What is Vendor](https://www.alibabacloud.com/help/en/security-center/developer-reference/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_threat_detection_vendor" "default" {
  vendor_name = "my-vendor"
  lang        = "en"
}
```

## Argument Reference

The following arguments are supported:

* `vendor_name` - (Required) The name of the vendor.
* `lang` - (Optional) The language type of the API request. Valid values: `en`, `zh`. Default value: `en`.
* `role_for` - (Optional) The user ID when the administrator switches other views.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Vendor. It is the same as `vendor_id`.
* `vendor_id` - The ID of the vendor.
* `vendor_type` - The type of the vendor.
* `create_time` - The creation time of the resource.
* `update_time` - The last modification time of the resource.

## Import

Threat Detection Vendor can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_vendor.example <id>
```
