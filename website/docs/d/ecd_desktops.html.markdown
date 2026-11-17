---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_desktops"
description: |-
  Provides a list of Ecd Desktops to the user.
---

# alicloud_ecd_desktops

This data source provides the Ecd Desktops of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.144.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

data "alicloud_ecd_images" "default" {
  image_type            = "SYSTEM"
  os_type               = "Windows"
  desktop_instance_type = "eds.hf.4c8g"
}

data "alicloud_ecd_desktop_types" "default" {
  instance_type_family = "eds.hf"
  cpu_count            = 4
  memory_size          = 8192
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
}

resource "alicloud_ecd_bundle" "default" {
  desktop_type                = data.alicloud_ecd_desktop_types.default.ids.0
  image_id                    = data.alicloud_ecd_images.default.ids.0
  user_disk_size_gib          = [110]
  root_disk_size_gib          = 120
  root_disk_performance_level = "PL1"
  user_disk_performance_level = "PL1"
  bundle_name                 = var.name
  description                 = var.name
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
  end_user_id = var.name
  email       = "${var.name}@aaa.com"
}

resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = alicloud_ecd_bundle.default.id
  desktop_name    = var.name
  end_user_ids    = [alicloud_ecd_user.default.id]
}

data "alicloud_ecd_desktops" "ids" {
  ids = [alicloud_ecd_desktop.default.id]
}

output "ecd_desktops_id_0" {
  value = data.alicloud_ecd_desktops.ids.desktops.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional)  A list of Desktop IDs.
* `name_regex` - (Optional) A regex string to filter results by Desktop name.
* `office_site_id` - (Optional) The ID of the Simple Office Site.
* `policy_group_id` - (Optional) The cloud desktop policy ID.
* `desktop_name` - (Optional) The cloud desktop name.
* `status` - (Optional) The cloud desktop status. Valid values: `Starting`, `Pending`, `Running`, `Rebuilding`, `Stopping`, `Stopped`, `Deleted`, `Expired`.
* `end_user_ids` - (Optional, List) The authorized users of the cloud desktop.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Desktop names.
* `desktops` - A list of Ecd Desktops. Each element contains the following attributes:
  * `id` - The ID of the Desktop.
  * `desktop_id` - The cloud desktop ID.
  * `desktop_type` - The specifications of the cloud desktop.
  * `desktop_name` - The cloud desktop name.
  * `policy_group_id` - The cloud desktop policy ID.
  * `directory_id` - The directory ID.
  * `network_interface_id` - The ID of the secondary elastic network interface (ENI) created by the cloud desktop service for a RAM or AD user.
  * `cpu` - The number of CPUs.
  * `image_id` - The image ID.
  * `memory` - The memory size.
  * `system_disk_size` - The system disk size of the Desktop.
  * `payment_type` - The billing method of the cloud desktop.
  * `status` - The cloud desktop status.
  * `create_time` - The creation time of the cloud desktop.
  * `expired_time` - The expiration time of the subscription cloud desktop.
  * `end_user_ids` - The IDs of the authorized users of the cloud desktop.
  * `disks` - (Available since v2.3.0) The disk information.
    * `disk_id` - The disk ID.
    * `disk_type` - The disk type.
    * `disk_category` - The disk category.
    * `disk_size` - The disk size.
    * `performance_level` - The performance level of the enhanced standard SSD cloud disk.
