---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_route_entries"
sidebar_current: "docs-alicloud-datasource-cen-route-entries"
description: |-
    Provides a list of CEN Route Entries owned by an Alibaba Cloud account.
---

# alicloud\_cen\_route\_entries

This data source provides CEN Route Entries available to the user.

## Example Usage

```
provider "alicloud" {
    alias = "bj"
    region = "cn-beijing"
}
data "alicloud_cen_route_entries" "entry"{
    provider = "alicloud.bj"
	instance_id = "cen-id1"
	route_table_id = "vtb-id1"
}

output "first_route_entries_route_entry_cidr_block" {
  value = "${data.alicloud_cen_bandwidth_packages.entry.route_entries.0.cidr_block}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the CEN instance.
* `route_table_id` - (Required) ID of the route table of the VPC or VBR.
* `cidr_block` - (Optional) The destination CIDR block of the route entry to query.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `route_entries` - A list of CEN Route Entries. Each element contains the following attributes:
  * `route_table_id` - ID of the route table.
  * `cidr_block` - The destination CIDR block of the route entry.
  * `next_hop_id` - ID of the next hop.
  * `next_hop_type` - Type of the next hop, including "Instance", "HaVip" and "RouterInterface".
  * `route_type` - Type of the route entry, including "System", "Custom" and "BGP".
  * `operational_mode` - Whether to allow the route entry to be published or removed to or from CEN.
  * `publish_status` - The publish status of the route entry in CEN, including "Published" and "NonPublished".
  * `conflicts` - A list of conflicted Route Entries. Each element contains the following attributes:
    * `cidr_block` - The destination CIDR block of the conflicted route entry.
    * `region_id` - ID of the region where the conflicted route entry is located.
    * `instance_id` - ID of the CEN child instance.
    * `instance_type` - The type of the CEN child instance.
    * `status` - Reasons of exceptions.