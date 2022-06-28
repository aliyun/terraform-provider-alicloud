---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_shared_resources"
sidebar_current: "docs-alicloud-datasource-resource-manager-shared-resources"
description: |-
  Provides a list of Resource Manager Shared Resources to the user.
---

# alicloud\_resource\_manager\_shared\_resources

This data source provides the Resource Manager Shared Resources of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_shared_resources" "this" {
  resource_share_id = "rs-V2NV******"
  ids               = ["vsw-bp1mzouzpmvie********:VSwitch"]
}

output "first_resource_manager_shared_resource_id" {
  value = data.alicloud_resource_manager_shared_resources.example.resources.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of shared resource ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_share_id` - (Optional, ForceNew) The resource share ID of resource manager.
* `status` - (Optional, ForceNew) The status of share resource, valid values `Associated`,`Associating`,`Disassociated`,`Disassociating`, and `Failed`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `resources` - A list of Resource Manager Shared Resources. Each element contains the following attributes:
	* `id` - The ID of the Shared Resource.
	* `resource_id` - The ID of the shared resource.
	* `resource_share_id` - The resource share ID of resource manager.
	* `resource_type` - The type of shared resource.
	* `status` - The status of shared resource.
