---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_nas_file_systems"
sidebar_current: "docs-alicloud-datasource-ecd-nas-file-systems"
description: |-
  Provides a list of Ecd Nas File Systems to the user.
---

# alicloud\_ecd\_nas\_file\_systems

This data source provides the Ecd Nas File Systems of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.141.0+.

## Example Usage

Basic Usage

```terraform

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block             = "172.16.0.0/12"
  desktop_access_type    = "Internet"
  office_site_name       = "your_office_site_name"
  enable_internet_access = false
}

resource "alicloud_ecd_nas_file_system" "default" {
  description          = "your_description"
  office_site_id       = alicloud_ecd_simple_office_site.default.id
  nas_file_system_name = "your_nas_file_system_name"
}

data "alicloud_ecd_nas_file_systems" "ids" {}
output "ecd_nas_file_system_id_1" {
  value = data.alicloud_ecd_nas_file_systems.ids.systems.0.id
}

data "alicloud_ecd_nas_file_systems" "nameRegex" {
  name_regex = alicloud_ecd_nas_file_system.default.nas_file_system_name
}
output "ecd_nas_file_system_id_2" {
  value = data.alicloud_ecd_nas_file_systems.nameRegex.systems.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Nas File System IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Nas File System name.
* `office_site_id` - (Optional, ForceNew) The ID of office site.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of nas file system. Valid values: `Pending`, `Running`, `Stopped`,`Deleting`, `Deleted`, `Invalid`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Nas File System names.
* `systems` - A list of Ecd Nas File Systems. Each element contains the following attributes:
	* `capacity` - The capacity of nas file system.
	* `create_time` - The create time of nas file system.
	* `description` - The description of nas file system.
	* `file_system_id` - The filesystem id of nas file system.
	* `file_system_type` - The type of nas file system.
	* `id` - The ID of the Nas File System.
	* `metered_size` - The size of metered.
	* `mount_target_domain` - The domain of mount target.
	* `mount_target_status` - The status of mount target. Valid values: `Pending`, `Active`, `Inactive`,`Deleting`,`Invalid`.
	* `nas_file_system_name` - The name of nas file system.
	* `office_site_id` - The ID of office site.
	* `office_site_name` - The name of office site.
	* `status` - The status of nas file system. Valid values: `Pending`, `Running`, `Stopped`,`Deleting`, `Deleted`, `Invalid`.
	* `storage_type` - The storage type of nas file system.
	* `support_acl` - Whether to support Acl.
	* `zone_id` - The zone id of nas file system.
