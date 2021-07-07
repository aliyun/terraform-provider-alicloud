---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_entries"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-route-entries"
description: |-
Provides a list of CEN Transit Router Route Entries owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_router\_route\_entries

This data source provides CEN Transit Router Route Entries available to the user.

-> **NOTE:** Available in 1.125.0+

## Example Usage

```
data "alicloud_cen_transit_router_route_entries" "entry" {
  transit_router_route_table_id = "vtb-*********"
}

output "transit_router_route_entry_destination_cidr_block"" {
  value = "${data.alicloud_cen_transit_router_route_entries.transit_router_route_entries.0.transit_router_route_entry_destination_cidr_block}"
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_route_table_id` - (Required) ID of the CEN Transit Router Route Table.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `transit_router_route_entries` - A list of CEN Route Entries. Each element contains the following attributes:
    * `transit_router_route_entry_destination_cidr_block` - The destination CIDR block of the route entry.
    * `transit_router_route_entry_next_hop_id` - ID of the next hop.
    * `transit_router_route_entry_next_hop_type` - Type of the next hop.
    * `transit_router_route_entry_type` - Type of the route entry.
    * `transit_router_route_entry_status` - The status of the route entry in CEN.
