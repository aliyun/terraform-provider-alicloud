---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_groups"
sidebar_current: "docs-alicloud-datasource-resource-manager-resource-groups"
description: |-
    Provides a list of resource groups to the user.
---

# alicloud\_resource\_manager\_resource\_groups

This data source provides resource groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.84.0+.

## Example Usage

```
data "alicloud_resource_manager_resource_groups" "example" {
  name_regex = "tftest"
}

output "first_resource_group_id" {
  value = "${data.alicloud_resource_manager_resource_groups.example.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of resource group IDs.
* `name_regex` - (Optional) A regex string to filter results by resource group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional) The status of the resource group. Possible values:`Creating`,`Deleted`,`OK` and `PendingDelete`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of resource group IDs.
* `names` - A list of resource group names.
* `groups` - A list of resource groups. Each element contains the following attributes:
  * `id` - The ID of the resource group.
  * `name` - The unique identifier of the resource group.
  * `display_name` - The display name of the resource group.
  * `account_id` - The ID of the Alibaba Cloud account to which the resource group belongs.
  * `create_date` - The time when the resource group was created.
  * `status` - The status of the resource group. Possible values:`Creating`,`Deleted`,`OK` and `PendingDelete`.
