---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_groups"
sidebar_current: "docs-alicloud-datasource-nas-access-groups"
description: |-
  Provides a list of Access Groups owned by an Alibaba Cloud account.
---

# alicloud\_nas_access_groups

This data source provides user-available access groups. Use when you can create mount points

-> NOTE: Available in 1.35.0+

## Example Usage

```terraform
data "alicloud_nas_access_groups" "example" {
  name_regex        = "^foo"
  access_group_type = "Classic"
  description       = "tf-testAccAccessGroupsdatasource"
}

output "alicloud_nas_access_groups_id" {
  value = data.alicloud_nas_access_groups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter AccessGroups by name. 
* `type` - (Optional, ForceNew, Deprecated in v1.95.0+) Field `type` has been deprecated from version 1.95.0. Use `access_group_type` instead.
* `access_group_type` - (Optional, ForceNew ,Available in 1.95.0+) Filter results by a specific AccessGroupType.
* `description` - (Optional, ForceNew) Filter results by a specific Description.
* `access_group_name` - (Optional, ForceNew, Available in 1.95.0+) The name of access group.
* `file_system_type` - (Optional, ForceNew, Available in 1.95.0+) The type of file system. Valid values: `standard` and `extreme`. Default to `standard`.
* `useutc_date_time` - (Optional, ForceNew, Available in 1.95.0+) Specifies whether the time to return is in UTC. Valid values: true and false.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of AccessGroup IDs, the value is set to `names`. After version 1.95.0 the item value as `<access_group_id>:<file_system_type>`. 
* `names` - A list of AccessGroup names.
* `groups` - A list of AccessGroups. Each element contains the following attributes:
 * `id` - This ID of this AccessGroup. It is formatted to ``<access_group_id>:<file_system_type>``. Before version 1.95.0, the value is `access_group_name`.
 * `rule_count` - RuleCount of the AccessGroup.
 * `type` - (Deprecated in v1.95.0+) AccessGroupType of the AccessGroup. The Field replace by `access_group_type` after version 1.95.0.
 * `mount_target_count` - MountTargetCount block of the AccessGroup
 * `description` - Description of the AccessGroup.
 * `access_group_name` - (Available in 1.95.0+) The name of the AccessGroup.
 * `access_group_type` - (Available in 1.95.0+) The type of the AccessGroup.
