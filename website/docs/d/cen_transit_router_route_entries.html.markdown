---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_entries"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-route-entries"
description: |-
  Provides a list of CEN Transit Router Route Entries owned by an Alibaba Cloud account.
---

# alicloud_cen_transit_router_route_entries

This data source provides CEN Transit Router Route Entries available to the user.[What is Cen Transit Router Route Entries](https://www.alibabacloud.com/help/doc-detail/260941.html)

-> **NOTE:** Available since v1.126.0+

## Example Usage

```terraform
data "alicloud_cen_transit_router_route_entries" "entry" {
  transit_router_route_table_id = "vtb-*********"
}

output "transit_router_route_entry_destination_cidr_block" {
  value = data.alicloud_cen_transit_router_route_entries.transit_router_route_entries.0.transit_router_route_entry_destination_cidr_block
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_route_table_id` - (Required) ID of the CEN Transit Router Route Table.
* `transit_router_route_entry_ids` - (Optional) A list of ID of the cen transit router route entry.
* `transit_router_route_entry_names` - (Optional) A list of name of the cen transit router route entry.
* `transit_router_route_entry_status` - (Optional) The status of the resource.Valid values `Creating`, `Active` and `Deleting`.
* `status` - (Optional) The status of the transit router route entry. Valid values: `Active`, `Creating`, `Deleting` and `Updating`.
* `name_regex` - (Optional) A regex string to filter results by transit router route entry name.
* `ids` - (Optional) A list of CEN Transit Router Route Entry IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN Transit Router Route Entry IDs.
* `names` - A list of CEN Transit Router Route Entry Names.
* `entries` - A list of CEN Route Entries. Each element contains the following attributes:
  * `id` - The ID of the CEN Transit Router Route Entry.
  * `transit_router_route_entry_id` - The ID of the route entry.
  * `transit_router_route_entry_name` - The name of the transit router route entry.
  * `transit_router_route_entry_description` - The description of the transit router route entry.
  * `transit_router_route_entry_destination_cidr_block` - The destination CIDR block of the route entry.
  * `transit_router_route_entry_next_hop_id` - ID of the next hop.
  * `transit_router_route_entry_next_hop_type` - Type of the next hop.
  * `transit_router_route_entry_type` - Type of the route entry.
  * `status` - The status of the route entry in CEN.
