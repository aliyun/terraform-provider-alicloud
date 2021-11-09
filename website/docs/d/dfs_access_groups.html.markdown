---
subcategory: "Apsara File Storage for HDFS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_access_groups"
sidebar_current: "docs-alicloud-datasource-dfs-access-groups"
description: |-
  Provides a list of Apsara File Storage for HDFS Access Groups to the user.
---

# alicloud\_dfs\_access\_groups

This data source provides the Apsara File Storage for HDFS Access Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dfs_access_groups" "ids" {
  ids = ["example_id"]
}
output "dfs_access_group_id_1" {
  value = data.alicloud_dfs_access_groups.ids.groups.0.id
}

data "alicloud_dfs_access_groups" "nameRegex" {
  name_regex = "^my-AccessGroup"
}
output "dfs_access_group_id_2" {
  value = data.alicloud_dfs_access_groups.nameRegex.groups.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Access Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Access Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Access Group names.
* `groups` - A list of Dfs Access Groups. Each element contains the following attributes:
	* `access_group_id` - The length of `description` does not exceed 100 bytes.
	* `access_group_name` - The Name of Access Group. The length Of `access_group_name` does not exceed 100 bytes.
	* `create_time` - The CreateTime of Access Group.
	* `description` - The Description of Access Group. The length Of `description` does not exceed 100 bytes.
	* `id` - The ID of the Access Group.
	* `mount_point_count` - The Number of attached mountpoint.
	* `network_type` - The NetworkType of Access Group. Valid values: `VPC`.
	* `rule_count` - The Number of access rule.
