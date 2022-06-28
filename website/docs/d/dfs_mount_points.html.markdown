---
subcategory: "Apsara File Storage for HDFS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_mount_points"
sidebar_current: "docs-alicloud-datasource-dfs-mount-points"
description: |-
  Provides a list of Dfs Mount Points to the user.
---

# alicloud\_dfs\_mount\_points

This data source provides the Dfs Mount Points of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dfs_mount_points" "ids" {
  file_system_id = "example_value"
  ids            = ["example_value-1", "example_value-2"]
}
output "dfs_mount_point_id_1" {
  value = data.alicloud_dfs_mount_points.ids.points.0.id
}

```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, ForceNew) The ID of the File System.
* `ids` - (Optional, ForceNew, Computed)  A list of Mount Point IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the Mount Point. Valid values: `Active`, `Inactive`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `points` - A list of Dfs Mount Points. Each element contains the following attributes:
	* `access_group_id` - The ID of the Access Group.
	* `create_time` - The created time of the Mount Point.
	* `description` - The description of the Mount Point.
	* `file_system_id` - The ID of the File System.
	* `id` - The ID of the Mount Point.
	* `mount_point_domain` - The domain name of the Mount Point.
	* `mount_point_id` - The ID of the Mount Point.
	* `network_type` - The network type of the Mount Point. Valid values: `VPC`.
	* `status` - The status of the Mount Point. Valid values: `Active`, `Inactive`.
	* `vpc_id` - The ID of the VPC network.
	* `vswitch_id` - The vswitch id.
