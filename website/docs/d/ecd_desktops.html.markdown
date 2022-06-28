---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_desktops"
sidebar_current: "docs-alicloud-datasource-ecd-desktops"
description: |-
  Provides a list of Ecd Desktops to the user.
---

# alicloud\_ecd\_desktops

This data source provides the Ecd Desktops of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = "your_office_site_name"

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

data "alicloud_ecd_desktops" "ids" {
  ids = [alicloud_ecd_desktop.default.id]
}
output "ecd_desktop_id_1" {
  value = data.alicloud_ecd_desktops.ids.desktops.0.id
}

data "alicloud_ecd_desktops" "nameRegex" {
  name_regex = alicloud_ecd_desktop.default.desktop_name
}
output "ecd_desktop_id_2" {
  value = data.alicloud_ecd_desktops.nameRegex.desktops.0.id
}
```

## Argument Reference

The following arguments are supported:

* `desktop_name` - (Optional, ForceNew) The desktop name.
* `office_site_id` - (Optional, ForceNew) The ID of the Simple Office Site.
* `ids` - (Optional, ForceNew, Computed)  A list of Desktop IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Desktop name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `policy_group_id` - (Optional, ForceNew) The policy group id of the Desktop.
* `status` - (Optional, ForceNew) The status of the Desktop. Valid values: `Deleted`, `Expired`, `Pending`, `Running`, `Starting`, `Stopped`, `Stopping`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Desktop names.
* `desktops` - A list of Ecd Desktops. Each element contains the following attributes:
	* `cpu` - The number of CPUs.
	* `create_time` - The creation time of the Desktop.
	* `desktop_id` - The desktop id of the Desktop.
	* `desktop_name` - The desktop name of the Desktop.
	* `desktop_type` - The desktop type of the Desktop.
	* `directory_id` - The directory id of the Desktop.
	* `end_user_ids` - The desktop end user id of the Desktop.
	* `expired_time` - The expired time of the Desktop.
	* `id` - The ID of the Desktop.
	* `image_id` - The image id of the Desktop.
	* `memory` - The memory of the Desktop.
	* `network_interface_id` - The network interface id of the Desktop.
	* `payment_type` - The payment type of the Desktop.
	* `policy_group_id` - The policy group id of the Desktop.
	* `status` - The status of the Desktop. Valid values: `Deleted`, `Expired`, `Pending`, `Running`, `Starting`, `Stopped`, `Stopping`.
	* `system_disk_size` - The system disk size of the Desktop.