---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_tables"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-route-tables"
description: |-
  Provides a list of CEN Transit Router Route Tables owned by an Alibaba Cloud account.
---

# alicloud\_cen\_transit\_router\_route\_tables

This data source provides CEN Transit Router Route Tables available to the user.[What is Cen Transit Router Route Tables](https://help.aliyun.com/document_detail/261237.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

```
data "alicloud_cen_transit_router_route_tables" "entry" {
  transit_router_id = "tr-*********"
}

output "first_transit_router_route_table_type"" {
  value = data.alicloud_cen_transit_router_route_tables.tables.0.transit_router_route_table_type
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_id` - (Required) ID of the CEN Transit Router Route Table.
* `transit_router_route_table_names` - (Optional) A list of name of the CEN Transit Router Route Table.  
* `transit_router_route_table_ids` - (Optional) A list of ID of the CEN Transit Router Route Table.
* `transit_router_route_table_type` - (Optional) The type of the transit router route table to query. Valid values `Creating`, `Active` and `Deleting`..
* `transit_router_route_table_status` - (Optional) The status of the transit router route table to query.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN Transit Router Route Table IDs.
* `names` - A list of name of CEN Transit Router Route Tables.
* `tables` - A list of CEN Route Entries. Each element contains the following attributes:
    * `transit_router_route_table_status` - The status of the route table.
    * `transit_router_route_table_description` - The description of the transit router route table.
    * `id` - ID of resource.
    * `transit_router_route_table_id` - ID of the trabsit router route table.
    * `transit_router_route_table_name` - Name of the transit router route table.  
    * `transit_router_route_table_type` - Type of the transit router route table.
