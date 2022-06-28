---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_shares"
sidebar_current: "docs-alicloud-datasource-resource-manager-resource-shares"
description: |-
  Provides a list of Resource Manager Resource Shares to the user.
---

# alicloud\_resource\_manager\_resource\_shares

This data source provides the Resource Manager Resource Shares of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_shares" "example" {
  resource_share_owner = "Self"
  ids                  = ["example_value"]
  name_regex           = "the_resource_name"
}

output "first_resource_manager_resource_share_id" {
  value = data.alicloud_resource_manager_resource_shares.example.shares.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Resource Share IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Resource Share name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_share_name` - (Optional, ForceNew) The name of resource share.
* `resource_share_owner` - (Required, ForceNew) The owner of resource share, Valid values: `Self` and `OtherAccounts`.
* `status` - (Optional, ForceNew) The status of resource share. Valid values: `Active`,`Deleted` and `Deleting`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Resource Share names.
* `shares` - A list of Resource Manager Resource Shares. Each element contains the following attributes:
	* `id` - The ID of the Resource Share.
	* `resource_share_id` - The ID of the resource share.
	* `resource_share_name` - The name of resource share.
	* `resource_share_owner` - The owner of resource share.
	* `status` - The status of resource share.
