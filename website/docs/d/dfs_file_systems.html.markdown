---
subcategory: "Apsara File Storage for HDFS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_file_systems"
sidebar_current: "docs-alicloud-datasource-dfs-file-systems"
description: |-
  Provides a list of Dfs File Systems to the user.
---

# alicloud\_dfs\_file\_systems

This data source provides the Dfs File Systems of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dfs_file_systems" "ids" {
  ids = ["example_id"]
}
output "dfs_file_system_id_1" {
  value = data.alicloud_dfs_file_systems.ids.systems.0.id
}

data "alicloud_dfs_file_systems" "nameRegex" {
  name_regex = "^my-FileSystem"
}
output "dfs_file_system_id_2" {
  value = data.alicloud_dfs_file_systems.nameRegex.systems.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of File System IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by File System name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of File System names.
* `systems` - A list of Dfs File Systems. Each element contains the following attributes:
	* `create_time` - The creation time of the File system.
	* `description` - The description of the File system.
	* `file_system_id` - The ID of the File System.
	* `file_system_name` - The name of the File system.
	* `id` - The ID of the File System.
	* `mount_point_count` - The number of Mount points.
	* `number_of_directories` - The number of directories.
	* `number_of_files` - The number of files.
	* `protocol_type` - The protocol type. Valid values: `HDFS`.
	* `provisioned_throughput_in_mi_bps` - The preset throughput of the File system. Valid values: `1` to `1024`, Unit: MB/s.
	* `space_capacity` - The capacity budget of the File system.
	* `storage_package_id` - Storage package Id.
	* `storage_type` - The storage specifications of the File system. Valid values: `PERFORMANCE`, `STANDARD`.
	* `throughput_mode` - The throughput mode of the File system. Valid values: `Provisioned`, `Standard`.
	* `used_space_size` - The used space of the File system.
	* `zone_id` - The zone ID of the File system.
