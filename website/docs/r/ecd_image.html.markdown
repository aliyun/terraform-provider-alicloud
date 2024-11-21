---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_image"
sidebar_current: "docs-alicloud-resource-ecd-image"
description: |-
  Provides a Alicloud ECD Image resource.
---

# alicloud_ecd_image

Provides a ECD Image resource.

For information about ECD Image and how to use it, see [What is Image](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-ecd-2020-09-30-createimage).

-> **NOTE:** Available since v1.146.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_image&exampleId=2031a4eb-bf37-5c73-1093-c676017e32257afe0c11&activeTab=example&spm=docs.r.ecd_image.0.2031a4ebbf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  enable_admin_access = true
  desktop_access_type = "Internet"
  office_site_name    = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard         = "read"
  local_drive       = "read"
  usb_redirect      = "off"
  watermark         = "off"

  authorize_access_policy_rules {
    description = var.name
    cidr_ip     = "1.2.3.45/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = var.name
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "1.2.3.4/24"
  }
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = data.alicloud_ecd_bundles.default.bundles.1.id
  desktop_name    = var.name
}

resource "alicloud_ecd_image" "default" {
  image_name  = var.name
  desktop_id  = alicloud_ecd_desktop.default.id
  description = var.name
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
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 15 mines) Used when create the Image.
* `delete` - (Defaults to 5 mines) Used when delete the Image.

## Import

ECD Image can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_image.example <id>
```