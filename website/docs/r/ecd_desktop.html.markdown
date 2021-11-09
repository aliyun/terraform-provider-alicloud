---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_desktop"
sidebar_current: "docs-alicloud-resource-ecd-desktop"
description: |-
  Provides a Alicloud ECD Desktop resource.
---

# alicloud\_ecd\_desktop

Provides a ECD Desktop resource.

For information about ECD Desktop and how to use it, see [What is Desktop](https://help.aliyun.com/document_detail/188382.html).

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block             = "172.16.0.0/12"
  desktop_access_type    = "Internet"
  office_site_name       = "your_office_site_name"
  enable_internet_access = false
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
  bundle_id       = data.alicloud_ecd_bundles.default.bundles.0.id
  desktop_name    = "your_desktop_name"
  end_user_ids    = [alicloud_ecd_user.default.id]
}

resource "alicloud_ecd_user" "default" {
  end_user_id = "your_end_user_id"
  email       = "your_email"
  phone       = "your_phone"
  password    = "your_password"
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mines) Used when create the Desktop.
* `delete` - (Defaults to 10 mines) Used when delete the Desktop.
* `update` - (Defaults to 20 mines) Used when update the Desktop.

## Import

ECD Desktop can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_desktop.example <id>
```