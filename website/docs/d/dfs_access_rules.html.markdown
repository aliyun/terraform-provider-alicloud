---
subcategory: "Apsara File Storage for HDFS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_access_rules"
sidebar_current: "docs-alicloud-datasource-dfs-access-rules"
description: |-
  Provides a list of Dfs Access Rules to the user.
---

# alicloud\_dfs\_access\_rules

This data source provides the Dfs Access Rules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dfs_access_rules" "ids" {
  access_group_id = "example_value"
  ids             = ["example_value-1", "example_value-2"]
}
output "dfs_access_rule_id_1" {
  value = data.alicloud_dfs_access_rules.ids.rules.0.id
}

```

## Argument Reference

The following arguments are supported:

* `access_group_id` - (Required, ForceNew) The resource ID of the Access Group.
* `ids` - (Optional, ForceNew, Computed)  A list of Access Rule IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `rules` - A list of Dfs Access Rules. Each element contains the following attributes:
    * `access_group_id` - The resource ID of the Access Group.
    * `access_rule_id` - The ID of the Access Rule.
    * `create_time` - The created time of the Access Rule.
    * `description` - The description of the Access Rule.
    * `id` - The resource ID of Access Rule.
    * `network_segment` - The NetworkSegment of the Access Rule.
    * `priority` - The priority of the Access Rule.
    * `rw_access_type` - RWAccessType of the Access Rule. Valid values: `RDONLY`, `RDWR`.
