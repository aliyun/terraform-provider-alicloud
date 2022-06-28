---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_image"
sidebar_current: "docs-alicloud-resource-ecd-image"
description: |-
  Provides a Alicloud ECD Image resource.
---

# alicloud\_ecd\_image

Provides a ECD Image resource.

For information about ECD Image and how to use it, see [What is Image](https://help.aliyun.com/document_detail/188382.html).

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = "your_simple_office_site_name"
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = "your_policy_group_name"
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = "example_value"
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "example_value"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}

resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = data.alicloud_ecd_bundles.default.bundles.1.id
  desktop_name    = "your_desktop_name"
}

resource "alicloud_ecd_image" "default" {
  image_name  = "your_image_name"
  desktop_id  = alicloud_ecd_desktop.default.id
  description = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the image.
* `desktop_id` - (Required) The desktop id of the desktop.
* `image_name` - (Optional) The name of the image.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Image.
* `status` - The status of the image. Valid values: `Creating`, `Available`, `CreateFailed`.
### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 15 mines) Used when create the Image.
* `delete` - (Defaults to 5 mines) Used when delete the Image.

## Import

ECD Image can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_image.example <id>
```