---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_delegated_administrators"
sidebar_current: "docs-alicloud-datasource-resource-manager-delegated-administrators"
description: |-
  Provides a list of Resource Manager Delegated Administrators to the user.
---

# alicloud\_resource\_manager\_delegated\_administrators

This data source provides the Resource Manager Delegated Administrators of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.181.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_delegated_administrators" "ids" {
  ids = ["example_value"]
}
output "resource_manager_delegated_administrator_id_1" {
  value = data.alicloud_resource_manager_delegated_administrators.ids.administrators.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Delegated Administrator IDs.
* `service_principal` - (Optional, ForceNew) The identification of the trusted service.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `administrators` - A list of Resource Manager Delegated Administrators. Each element contains the following attributes:
	* `account_id` - The ID of the member account.
	* `delegation_enabled_time` - The time when the member was specified as a delegated administrator account.
	* `id` - The ID of the Delegated Administrator.
	* `service_principal` - The identity of the trusted service.