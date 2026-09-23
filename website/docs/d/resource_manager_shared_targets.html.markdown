---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_shared_targets"
sidebar_current: "docs-alicloud-datasource-resource-manager-shared-targets"
description: |-
  Provides a list of Resource Manager Shared Targets to the user.
---

# alicloud_resource_manager_shared_targets

This data source provides the Resource Manager Shared Targets of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_resource_manager_accounts" "default" {
}

resource "alicloud_resource_manager_resource_share" "default" {
  resource_share_name = var.name
}

resource "alicloud_resource_manager_shared_target" "default" {
  resource_share_id = alicloud_resource_manager_resource_share.default.id
  target_id         = data.alicloud_resource_manager_accounts.default.ids.0
}

data "alicloud_resource_manager_shared_targets" "ids" {
  ids = ["${alicloud_resource_manager_shared_target.default.target_id}"]
}

output "first_resource_manager_shared_target_id" {
  value = data.alicloud_resource_manager_shared_targets.ids.targets.0.id
}

data "alicloud_resource_manager_shared_targets" "resourceShareId" {
  resource_share_id = alicloud_resource_manager_shared_target.default.resource_share_id
}

output "second_resource_manager_shared_target_id" {
  value = data.alicloud_resource_manager_shared_targets.resourceShareId.targets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Shared Target IDs.
* `resource_share_id` - (Optional, ForceNew) The resource share ID of resource manager.
* `status` - (Optional, ForceNew) The status of share resource. Valid values: `Associated`, `Associating`, `Disassociated`, `Disassociating` and `Failed`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `targets` - A list of Resource Manager Shared Targets. Each element contains the following attributes:
  * `id` - The ID of the Shared Target.
  * `target_id` - The ID of the Shared Target.
  * `resource_share_id` - The resource shared ID of resource manager.
  * `status` - The status of shared target.
  