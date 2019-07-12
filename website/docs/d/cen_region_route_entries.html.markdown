---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_region_route_entries"
sidebar_current: "docs-alicloud-datasource-cen-region-route-entries"
description: |-
    Provides a list of CEN Route Entries from specific region owned by an Alibaba Cloud account.
---

# alicloud\_cen\_region\_route\_entries

This data source provides CEN Regional Route Entries available to the user.

## Example Usage

```
data "alicloud_cen_region_route_entries" "entry" {
  instance_id = "cen-id1"
  region_id   = "cn-beijing"
}

output "first_region_route_entries_route_entry_cidr_block" {
  value = "${data.alicloud_cen_region_route_entries.entry.entries.0.cidr_block}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the CEN instance.
* `region_id` - (Required) ID of the region.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `entries` - A list of CEN Route Entries. Each element contains the following attributes:
  * `cidr_block` - The destination CIDR block of the route entry.
  * `type` - Type of the route entry.
  * `next_hop_id` - ID of the next hop.
  * `next_hop_type` - Type of the next hop.
  * `next_hop_region_id` - ID of the region where the next hop is located.
