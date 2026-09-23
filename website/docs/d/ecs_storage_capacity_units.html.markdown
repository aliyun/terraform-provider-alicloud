---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_storage_capacity_units"
sidebar_current: "docs-alicloud-datasource-ecs-storage-capacity-units"
description: |-
  Provides a list of Ecs Storage Capacity Units to the user.
---

# alicloud\_ecs\_storage\_capacity\_units

This data source provides the Ecs Storage Capacity Units of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.155.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_storage_capacity_units" "ids" {}
output "ecs_storage_capacity_unit_id_1" {
  value = data.alicloud_ecs_storage_capacity_units.ids.units.0.id
}

data "alicloud_ecs_storage_capacity_units" "nameRegex" {
  name_regex = "^my-StorageCapacityUnit"
}
output "ecs_storage_capacity_unit_id_2" {
  value = data.alicloud_ecs_storage_capacity_units.nameRegex.units.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Storage Capacity Unit IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Storage Capacity Unit name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of Storage Capacity Unit. Valid values: `Active`, `Creating`, `Expired`, `Pending`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Storage Capacity Unit names.
* `units` - A list of Ecs Storage Capacity Units. Each element contains the following attributes:
	* `allocation_status` - When the AllocationType value is Shared, this parameter indicates the allocation status of Storage Capacity Unit. Valid values: `allocated`, `Ignored`.
	* `capacity` - The capacity of the Storage Capacity Unit.
	* `create_time` - The time when the Storage Capacity Unit was created.
	* `description` - The description of the Storage Capacity Unit.
	* `expired_time` - The time when the Storage Capacity Unit expires.
	* `id` - The ID of the Storage Capacity Unit.
	* `start_time` - The effective time of the Storage Capacity Unit.
	* `status` - The status of Storage Capacity Unit.
	* `storage_capacity_unit_id` - The ID of Storage Capacity Unit.
	* `storage_capacity_unit_name` - The name of the Storage Capacity Unit.