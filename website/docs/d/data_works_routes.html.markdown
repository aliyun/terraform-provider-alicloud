---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_routes"
description: |-
  Provides a list of Data Works Route resources to the user.
---

# alicloud_data_works_routes

This data source provides the Data Works Route resources available to the user.

Resource group network routing rules.

For information about Data Works Route and how to use it, see [What is Route](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-listroutes).

-> **NOTE:** Available since v1.288.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_data_works_routes" "default" {
  resource_group_id = "your_resource_group_id"
  ids               = ["1234567890"]
}

output "data_works_route_id_1" {
  value = data.alicloud_data_works_routes.default.routes.0.id
}
```

## Argument Reference

The following arguments are supported:
* `resource_group_id` - (Required) The ID of the resource group.
* `network_id` - (Optional) The ID of the network resource to which the routes belong.
* `sort_by` - (Optional) The fields used for sorting. The value is in the format of `SortField Desc/Asc`. Valid sort fields are `Id`, `DestinationCidr`, and `CreateTime`. Default value: `CreateTime Asc`.
* `ids` - (Optional) A list of route IDs used to filter the results.
* `output_file` - (Optional) File name where to write the result. If not set, the result is only shown in the Terraform state.

## Attributes Reference

The following attributes are exported:
* `id` - The data source ID.
* `ids` - A list of route IDs.
* `routes` - A list of Data Works Route resources. Each element contains the following attributes:
  * `id` - The ID of the route.
  * `route_id` - The ID of the route.
  * `destination_cidr` - The CIDR block of the destination route.
  * `network_id` - The ID of the network resource to which the route belongs.
  * `resource_id` - The identifier of the network resource to which the route belongs.
  * `create_time` - The time when the route was created.
  * `dw_resource_group_id` - The ID of the resource group to which the route belongs.
  * `region_id` - The region to which the route belongs.
