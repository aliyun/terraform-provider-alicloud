---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_network_interface_permissions"
sidebar_current: "docs-alicloud-datasource-ecs-network-interface-permissions"
description: |-
  Provides a list of Ecs Network Interface Permissions to the user.
---

# alicloud\_ecs\_network\_interface\_permissions

This data source provides the Ecs Network Interface Permissions of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.166.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_network_interface_permissions" "ids" {
  ids                  = ["example_value"]
  network_interface_id = "example_value"
}
output "ecs_network_interface_permission_id_1" {
  value = data.alicloud_ecs_network_interface_permissions.ids.permissions.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Network Interface Permission IDs.
* `network_interface_id` - (Required, ForceNew) The ID of the network interface.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The Status of the Network Interface Permissions. Valid values: `Granted`, `Pending`, `Revoked`, `Revoking`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `permissions` - A list of Ecs Network Interface Permissions. Each element contains the following attributes:
	* `account_id` - Alibaba Cloud Partner (Certified ISV) account ID or individual user ID.
	* `id` - The ID of the Network Interface Permission.
	* `network_interface_id` - The ID of the network interface.
	* `network_interface_permission_id` - The ID of the Network Interface Permissions.
	* `permission` - The permissions of the Network Interface.
	* `service_name` - Alibaba Cloud service name.
	* `status` - The Status of the Network Interface Permissions.