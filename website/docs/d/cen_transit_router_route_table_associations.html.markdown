---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_table_associations"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-route-table-associations"
description: |-
  Provides a list of CEN Transit Router Route Table Association owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_router\_route\_table\_associations

This data source provides CEN Transit Router Route Table Associations available to the user.[What is Cen Transit Router Route Table Associations](https://help.aliyun.com/document_detail/261243.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

```
data "alicloud_cen_transit_router_route_table_associations" "default" {
  transit_router_route_table_id    = "rtb-id1"
}

output "first_transit_router_peer_attachments_transit_router_attachment_resource_type" {
  value = data.alicloud_cen_transit_router_route_table_associations.default.associations.0.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_route_table_id` - (Optional) The ID of the route table of the Enterprise Edition transit router.
* `transit_router_attachment_id` - (Optional) The ID of the network instance connection. 
* `transit_router_attachment_resource_id` - (Optional) The ID of the next hop.
* `transit_router_attachment_resource_type` - (Optional) The type of next hop. Valid values:
  * `VPC`: virtual private cloud (VPC)
  * `VBR`: virtual border router (VBR)
  * `TR`: transit router
  * `VPN`: VPN attachment
  
* `status` - (Optional) The status of the route table, including `Active`, `Associating`, `Dissociating`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN Transit Router Route Table Association IDs.
* `names` - A list of CEN Transit Router Route Table Association Names.
* `associations` - A list of CEN Transit Router Route Table Associations. Each element contains the following attributes:
    * `resource_id` - ID of the transit router route table association.
    * `resource_type` - Type of the resource.
    * `status` - The status of the route table.
    * `transit_router_attachment_id` - ID of the transit router attachment.
    * `transit_router_route_table_id` - ID of the transit router route table.

