---
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_accessgroups"
sidebar_current: "docs-alicloud-datasource-nas-accessgroups"
description: |-
    Provides a list of AccessGroups owned by an Alibaba Cloud account.
---

# alicloud\_nas_accessgroups

This data source provides AccessGroup available to the user.

## Example Usage

```
data "alicloud_nas_accessgroups" "ag" {
  accessgroup_name = "CreateAccessGroup"
  accessgroup_type = "Classic"
  description = "test_AccessGroup"
}

output "first_nas_accessgroups_id" {
  value = "${data.alicloud_nas_accessgroups.nas_accessgroups_ds.accessgroups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `accessgroup_name` - (Required) Filter results by a specific AccessGroupName block. 
* `accessgroup_type` - (Optional) Filter results by a specific AccessGroupType block
* `description` - (Optional) Filter results by a specific Description block
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:


* `accessgroups` - A list of VPCs. Each element contains the following attributes:
  * `accessgroup_name`       - AccessGroupName of the AccessGroup.
  * `rule_count`             - RuleCount of the AccessGroup.
  * `accessgroup_type`       - AccessGroupType of the AccessGroup.
  * `mounttarget_count`      - MountTargetCount block of the AccessGroup
  * `description`            - Destription of the AccessGroup.
