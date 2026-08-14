---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_desktop_group"
description: |-
  Provides a Alicloud Elastic Desktop Service (ECD) Desktop Group resource.
---

# alicloud_ecd_desktop_group

Provides a Elastic Desktop Service (ECD) Desktop Group resource.

A desktop group is a pool of shared cloud desktops. End users authorized to the
desktop group are assigned an available cloud desktop from the pool when they
connect.

For information about Elastic Desktop Service (ECD) Desktop Group and how to use it, see [What is Desktop Group](https://next.api.alibabacloud.com/document/ecd/2020-09-30/CreateDesktopGroup).

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = var.name
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = var.name
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}

resource "alicloud_ecd_user" "default" {
  end_user_id = "terraform-example-user"
  email       = "terraform-example@example.com"
  password    = "Example12345"
}

resource "alicloud_ecd_desktop_group" "default" {
  desktop_group_name = var.name
  office_site_id     = alicloud_ecd_simple_office_site.default.id
  policy_group_id    = alicloud_ecd_policy_group.default.id
  bundle_id          = data.alicloud_ecd_bundles.default.bundles.0.id
  end_user_ids       = [alicloud_ecd_user.default.id]
  allow_buffer_count = 1
  comments           = var.name
}
```

## Argument Reference

The following arguments are supported:
* `bundle_id` - (Required) The ID of the desktop template.
* `end_user_ids` - (Required, List) The list of IDs of the end users authorized to use the desktop group.
* `office_site_id` - (Required, ForceNew) The ID of the office network to which the desktop group belongs.
* `policy_group_id` - (Required) The ID of the policy group associated with the desktop group.
* `allow_auto_setup` - (Optional, Int) Specifies whether to allow cloud desktops to be automatically created for a subscription desktop group. This parameter takes effect only when the desktop group uses the subscription billing method. Valid values: `0` and `1`.
* `allow_buffer_count` - (Optional, Int) The number of cloud desktops that are reserved in the desktop group. Reserved desktops are kept started and idle, waiting for connections. Valid values: `0` to `100`. `0` indicates that no desktop is reserved.
* `comments` - (Optional) The remarks of the desktop group.
* `desktop_group_name` - (Optional) The name of the desktop group. The name must be 1 to 30 characters in length, and can contain letters, digits, colons (:), underscores (_), periods (.), and hyphens (-). It must start with a letter but cannot start with `http://` or `https://`.
* `directory_id` - (Optional, ForceNew) The ID of the directory.
* `keep_duration` - (Optional, Int) The retention duration of a session after it is disconnected. Unit: milliseconds. Valid values: `180000` to `345600000`. `0` indicates that the session is always retained. If the user does not reconnect within the retention duration, the session is logged off and unsaved data is destroyed.
* `max_desktops_count` - (Optional, Int) The maximum number of cloud desktops that the pay-as-you-go desktop group can contain. Valid values: `0` to `500`.
* `min_desktops_count` - (Optional, Int) The minimum number of cloud desktops that the pay-as-you-go desktop group automatically creates. Valid values: `0` to the value of `max_desktops_count`.
* `scale_strategy_id` - (Optional) The ID of the scaling policy.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `cpu` - The number of vCPUs.
* `create_time` - The time when the desktop group was created. The time follows the ISO 8601 standard in UTC.
* `creator` - The ID of the Alibaba Cloud account that created the desktop group.
* `data_disk_category` - The category of the data disk.
* `data_disk_size` - The size of the data disk. Unit: GiB.
* `directory_type` - The type of the directory.
* `expired_time` - The time when the subscription desktop group expires. The time follows the ISO 8601 standard in UTC.
* `gpu_count` - The number of GPUs.
* `gpu_spec` - The GPU specifications.
* `memory` - The memory size. Unit: MiB.
* `office_site_name` - The name of the office network.
* `office_site_type` - The type of the account system of the office network.
* `own_bundle_name` - The name of the desktop template.
* `pay_type` - The billing method.
* `policy_group_name` - The name of the policy group.
* `res_type` - The type of the resource.
* `system_disk_category` - The category of the system disk.
* `system_disk_size` - The size of the system disk. Unit: GiB.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Desktop Group.
* `delete` - (Defaults to 10 mins) Used when delete the Desktop Group.
* `update` - (Defaults to 5 mins) Used when update the Desktop Group.

## Import

Elastic Desktop Service (ECD) Desktop Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_desktop_group.example <desktop_group_id>
```
