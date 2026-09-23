---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_desktop_groups"
sidebar_current: "docs-alicloud-datasource-ecd-desktop-groups"
description: |-
  Provides a list of Elastic Desktop Service (ECD) Desktop Groups owned by an Alibaba Cloud account.
---

# alicloud_ecd_desktop_groups

This data source provides Elastic Desktop Service (ECD) Desktop Groups available to the user. For information about Elastic Desktop Service (ECD) Desktop Group, see [What is Desktop Group](https://next.api.alibabacloud.com/document/ecd/2020-09-30/CreateDesktopGroup).

-> **NOTE:** Available since v1.291.0.

## Example Usage

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
}

data "alicloud_ecd_desktop_groups" "default" {
  ids        = [alicloud_ecd_desktop_group.default.id]
  name_regex = alicloud_ecd_desktop_group.default.desktop_group_name
}

output "alicloud_ecd_desktop_group_example_id" {
  value = data.alicloud_ecd_desktop_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:
* `desktop_group_id` - (Optional) The ID of the desktop group.
* `desktop_group_name` - (Optional) The name of the desktop group. Fuzzy search is supported.
* `office_site_id` - (Optional) The ID of the office network to which the desktop groups belong.
* `period_unit` - (Optional) The subscription duration unit of the desktop groups. Valid values: `Month` and `Year`.
* `ids` - (Optional, Computed) A list of Desktop Group IDs.
* `name_regex` - (Optional) A regex string to filter results by Desktop Group name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Desktop Group IDs.
* `names` - A list of name of Desktop Groups.
* `groups` - A list of Desktop Group Entries. Each element contains the following attributes:
  * `allow_auto_setup` - **NOTE:** This field is only available when `enable_details` is `true`. Whether cloud desktops can be automatically created for the subscription desktop group.
  * `allow_buffer_count` - **NOTE:** This field is only available when `enable_details` is `true`. The number of cloud desktops that are reserved in the desktop group.
  * `bundle_id` - The ID of the desktop template.
  * `comments` - The remarks of the desktop group.
  * `cpu` - The number of vCPUs.
  * `create_time` - The time when the desktop group was created. The time follows the ISO 8601 standard in UTC.
  * `creator` - The ID of the Alibaba Cloud account that created the desktop group.
  * `data_disk_category` - The category of the data disk.
  * `data_disk_size` - The size of the data disk. Unit: GiB.
  * `desktop_group_id` - The ID of the desktop group.
  * `desktop_group_name` - The name of the desktop group.
  * `directory_id` - The ID of the directory.
  * `directory_type` - The type of the directory.
  * `end_user_count` - The number of end users authorized to use the desktop group.
  * `end_user_ids` - **NOTE:** This field is only available when `enable_details` is `true`. The list of IDs of the end users authorized to use the desktop group.
  * `expired_time` - The time when the subscription desktop group expires. The time follows the ISO 8601 standard in UTC.
  * `gpu_count` - The number of GPUs.
  * `gpu_spec` - The GPU specifications.
  * `keep_duration` - The retention duration of a session after it is disconnected. Unit: milliseconds.
  * `max_desktops_count` - The maximum number of cloud desktops that the desktop group can contain.
  * `memory` - The memory size. Unit: MiB.
  * `min_desktops_count` - The minimum number of cloud desktops that the desktop group automatically creates.
  * `office_site_id` - The ID of the office network.
  * `office_site_name` - The name of the office network.
  * `office_site_type` - The type of the account system of the office network.
  * `own_bundle_name` - The name of the desktop template.
  * `pay_type` - The billing method.
  * `policy_group_id` - The ID of the policy group.
  * `policy_group_name` - The name of the policy group.
  * `res_type` - **NOTE:** This field is only available when `enable_details` is `true`. The type of the resource.
  * `system_disk_category` - The category of the system disk.
  * `system_disk_size` - The size of the system disk. Unit: GiB.
  * `id` - The ID of the resource supplied above.
