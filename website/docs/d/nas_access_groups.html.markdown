---
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

```
data "alicloud_nas_access_groups" "ag" {
  name_regex  = "^foo"
  type        = "Classic"
  description = "tf-testAccAccessGroupsdatasource"
}

output "alicloud_nas_access_groups_id" {
  value = "${data.alicloud_nas_access_groups.ag.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Required) A regex string to filter AccessGroups by name. 
* `type` - (Optional) Filter results by a specific AccessGroupType.
* `description` - (Optional) Filter results by a specific Description.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of AccessGroup IDs, the value is set to `names` .
* `names` - A list of AccessGroup names.
* `groups` - A list of AccessGroups. Each element contains the following attributes:
 * `id` - AccessGroupName of the AccessGroup.
 * `rule_count` - RuleCount of the AccessGroup.
 * `type` - AccessGroupType of the AccessGroup.
 * `mount_target_count` - MountTargetCount block of the AccessGroup
 * `description` - Destription of the AccessGroup.
