---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_route_table_aggregations"
sidebar_current: "docs-alicloud-datasource-cen-transit-route-table-aggregations"
description: |-
  Provides a list of Cen Transit Route Table Aggregations to the user.
---

# alicloud\_cen\_transit\_route\_table\_aggregations

This data source provides the Cen Transit Route Table Aggregations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.202.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_transit_route_table_aggregations" "ids" {
  ids                    = ["example_id"]
  transit_route_table_id = "your_transit_route_table_id"
}

output "cen_transit_router_multicast_domain_id_0" {
  value = data.alicloud_cen_transit_route_table_aggregations.ids.transit_route_table_aggregations.0.id
}

data "alicloud_cen_transit_route_table_aggregations" "nameRegex" {
  name_regex             = "^my-name"
  transit_route_table_id = "your_transit_route_table_id"
}

output "cen_transit_router_multicast_domain_id_1" {
  value = data.alicloud_cen_transit_route_table_aggregations.nameRegex.transit_route_table_aggregations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Transit Route Table Aggregation IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Transit Route Table Aggregation name.
* `transit_route_table_id` - (Required, ForceNew) The ID of the route table of the Enterprise Edition transit router.
* `transit_route_table_aggregation_cidr` - (Optional, ForceNew) The destination CIDR block of the aggregate route.
* `status` - (Optional, ForceNew) The status of Transit Route Table Aggregation. Valid Values: `AllConfigured`, `Configuring`, `ConfigFailed`, `PartialConfigured`, `Deleting`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Transit Route Table Aggregation names.
* `transit_route_table_aggregations` - A list of Cen Transit Route Table Aggregations. Each element contains the following attributes:
  * `id` - The ID of the Transit Route Table Aggregation. It formats as `<transit_route_table_id>:<transit_route_table_aggregation_cidr>`.
  * `transit_route_table_id` - The ID of the route table of the Enterprise Edition transit router.
  * `transit_route_table_aggregation_cidr` - The destination CIDR block of the aggregate route.
  * `transit_route_table_aggregation_scope` - The scope of networks that you want to advertise the aggregate route.
  * `route_type` - The route type of the aggregate route.
  * `transit_route_table_aggregation_name` - The name of the aggregate route.
  * `transit_route_table_aggregation_description` - The description of the aggregate route.
  * `status` - The status of the Transit Route Table Aggregation.
  