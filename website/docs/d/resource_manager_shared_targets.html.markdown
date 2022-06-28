---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_shared_targets"
sidebar_current: "docs-alicloud-datasource-resource-manager-shared-targets"
description: |-
  Provides a list of Resource Manager Shared Targets to the user.
---

# alicloud\_resource\_manager\_shared\_targets

This data source provides the Resource Manager Shared Targets of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_shared_targets" "example" {
  ids = ["15681091********"]
}

output "first_resource_manager_shared_target_id" {
  value = data.alicloud_resource_manager_shared_targets.example.targets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Shared Target IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_share_id` - (Optional, ForceNew) The resource share ID of resource manager.
* `status` - (Optional, ForceNew) The status of share resource, valid values `Associated`,`Associating`,`Disassociated`,`Disassociating`, and `Failed`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `targets` - A list of Resource Manager Shared Targets. Each element contains the following attributes:
	* `id` - The ID of the Shared Target.
	* `resource_share_id` - The resource shared ID of resource manager.
	* `status` - The status of shared target.
	* `target_id` - The member account ID in resource directory.
