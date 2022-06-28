---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_table_propagations"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-route-table-propagations"
description: |-
  Provides a list of CEN Transit Router Route Table Propagation owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_router\_route\_table\_propagations

This data source provides CEN Transit Router Route Table Propagations available to the user.[What is Cen Transit Router Route Table Propagations](https://help.aliyun.com/document_detail/261245.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

```
data "alicloud_cen_transit_router_route_table_propagations" "default" {
  transit_router_route_table_id    = "rtb-id1"
}

output "first_transit_router_peer_attachments_transit_router_attachment_resource_type" {
  value = data.alicloud_cen_transit_router_route_table_propagations.default.propagations.0.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_route_table_id` - (Optional) ID of the route table of the VPC or VBR.
* `transit_router_attachment_id` - (Optional) ID of the cen transit router attachment.  
* `status` - (Optional) The status of the route table, including `Active`, `Enabling`, `Disabling`, `Deleted`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN Transit Router Route Table Association IDs.
* `names` - A list of CEN Transit Router Route Table Association Names.
* `propagations` - A list of CEN Transit Router Route Table Propagations. Each element contains the following attributes:
    * `resource_id` - ID of the transit router route table association.
    * `resource_type` - Type of the resource.
    * `status` - The status of the route table.
    * `transit_router_attachment_id` - ID of the transit router attachment.
    * `transit_router_route_table_id` - ID of the transit router route table.

