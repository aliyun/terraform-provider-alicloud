---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_desktop"
sidebar_current: "docs-alicloud-resource-ecd-desktop"
description: |-
  Provides a Alicloud ECD Desktop resource.
---

# alicloud_ecd_desktop

Provides a ECD Desktop resource.

For information about ECD Desktop and how to use it, see [What is Desktop](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-ecd-2020-09-30-createdesktops)

-> **NOTE:** Available since v1.144.0.

## Example Usage


Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_desktop&exampleId=3ff92ccf-ee32-0331-2fef-f16f44573556f17005be&activeTab=example&spm=docs.r.ecd_desktop.0.3ff92ccfee&intl_lang=EN_US" target="_blank">
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
```

## Argument Reference

The following arguments are supported:

* `amount` - (Optional) The amount of the Desktop.
* `auto_pay` - (Optional) The auto-pay of the Desktop whether to pay automatically. values: `true`, `false`.
* `auto_renew` - (Optional) The auto-renewal of the Desktop whether to renew automatically. It takes effect only when the parameter ChargeType is set to PrePaid. values: `true`, `false`.
* `bundle_id` - (Required) The bundle id of the Desktop.
* `desktop_name` - (Optional) The desktop name of the Desktop.
* `desktop_type` - (Optional) The desktop type of the Desktop.
* `office_site_id` - (Required, ForceNew) The ID of the Simple Office Site.
* `end_user_ids` - (Optional, ForceNew) The desktop end user id of the Desktop.
* `host_name` - (Optional) The hostname of the Desktop.
* `payment_type` - (Optional, Computed) The payment type of the Desktop. Valid values: `PayAsYouGo`, `Subscription`. Default to `PayAsYouGo`.
* `period` - (Optional) The period of the Desktop.
* `period_unit` - (Optional) The period unit of the Desktop.
* `policy_group_id` - (Required) The policy group id of the Desktop.
* `root_disk_size_gib` - (Optional) The root disk size gib of the Desktop.
* `status` - (Optional, Computed) The status of the Desktop. Valid values: `Deleted`, `Expired`, `Pending`, `Running`, `Starting`, `Stopped`, `Stopping`.
* `stopped_mode` - (Optional) The stopped mode of the Desktop.
* `user_assign_mode` - (Optional) The user assign mode of the Desktop. Valid values: `ALL`, `PER_USER`. Default to `ALL`.
* `user_disk_size_gib` - (Optional) The user disk size gib of the Desktop.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Desktop.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mines) Used when create the Desktop.
* `delete` - (Defaults to 10 mines) Used when delete the Desktop.
* `update` - (Defaults to 20 mines) Used when update the Desktop.

## Import

ECD Desktop can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_desktop.example <id>
```