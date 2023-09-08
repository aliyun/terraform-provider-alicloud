---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_groups"
sidebar_current: "docs-alicloud-datasource-resource-manager-resource-groups"
description: |-
  Provides a list of resource groups to the user.
---

# alicloud_resource_manager_resource_groups

This data source provides resource groups of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.84.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "example" {
  name_regex = "tf"
}

output "first_resource_group_id" {
  value = data.alicloud_resource_manager_resource_groups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of resource group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by resource group identifier.
* `status` - (Optional, ForceNew) The status of the resource group. Valid values: `Creating`, `Deleted`, `Deleting`, `OK` and `PendingDelete`. **NOTE:** From version 1.114.0, `status` can be set to `Deleting`.
* `enable_details` -(Optional, Available since v1.114.0) Set it to true can output more details. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of resource group IDs.
* `names` - A list of resource group identifiers.
* `groups` - A list of resource groups. Each element contains the following attributes:
  * `id` - The ID of the resource group.
  * `name` - The unique identifier of the resource group.
  * `resource_group_name` - (Available since v1.114.0) The unique identifier of the resource group.
  * `display_name` - The display name of the resource group.
  * `account_id` - The ID of the Alibaba Cloud account to which the resource group belongs.
  * `status` - The status of the resource group.
  * `region_statuses`- (Available since v1.114.0) The status of the resource group in all regions.
    * `region_id` - The region ID.
    * `status` - The status of the regional resource group.
  