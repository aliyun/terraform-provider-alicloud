---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_tables"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-route-tables"
description: |-
  Provides a list of CEN Transit Router Route Tables owned by an Alibaba Cloud account.
---

# alicloud_cen_transit_router_route_tables

This data source provides CEN Transit Router Route Tables available to the user.[What is Cen Transit Router Route Tables](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-listtransitrouterroutetables)

-> **NOTE:** Available since v1.126.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id              = alicloud_cen_instance.default.id
  transit_router_name = var.name
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id                      = alicloud_cen_transit_router.default.transit_router_id
  transit_router_route_table_description = "desp"
  transit_router_route_table_name        = var.name
}

data "alicloud_cen_transit_router_route_tables" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

output "first_transit_router_route_table_type" {
  value = data.alicloud_cen_transit_router_route_tables.default.tables.0.transit_router_route_table_type
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_id` - (Required, ForceNew) ID of the CEN Transit Router Route Table.
* `transit_router_route_table_names` - (Optional, ForceNew) A list of name of the CEN Transit Router Route Table.  
* `transit_router_route_table_ids` - (Optional, ForceNew) A list of ID of the CEN Transit Router Route Table.
* `transit_router_route_table_type` - (Optional, ForceNew, Available since v1.229.1.) The type of the transit router route table to query. Valid values `System` and `Custom`.
* `transit_router_route_table_status` - (Optional, ForceNew) The status of the transit router route table to query. Valid values `Creating`, `Active` and `Deleting`..
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `name_regex` - (Optional, ForceNew) A regex string to filter CEN Transit Router Route Table by name.
* `ids` - (Optional, ForceNew) A list of CEN Transit Router Route Table IDs.
* `status` - (Optional, ForceNew) The status of the transit router route table to query. Valid values `Creating`, `Active` and `Deleting`..

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:


* `names` - A list of name of CEN Transit Router Route Tables.
* `tables` - A list of CEN Route Entries. Each element contains the following attributes:
    * `status` - The status of the route table.
    * `transit_router_route_table_description` - The description of the transit router route table.
    * `id` - ID of resource.
    * `transit_router_route_table_id` - ID of the trabsit router route table.
    * `transit_router_route_table_name` - Name of the transit router route table.  
    * `transit_router_route_table_type` - Type of the transit router route table.
