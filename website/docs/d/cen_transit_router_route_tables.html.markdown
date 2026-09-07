---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_tables"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-route-tables"
description: |-
  Provides a list of CEN Transit Router Route Tables to the user.
---

# alicloud_cen_transit_router_route_tables

This data source provides the CEN Transit Router Route Tables of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
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
  transit_router_route_table_description = var.name
  transit_router_route_table_name        = var.name
}

data "alicloud_cen_transit_router_route_tables" "ids" {
  transit_router_id = alicloud_cen_transit_router_route_table.default.transit_router_id
  ids               = [alicloud_cen_transit_router_route_table.default.transit_router_route_table_id]
}

output "cen_transit_router_route_table_id_0" {
  value = data.alicloud_cen_transit_router_route_tables.ids.tables.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Transit Router Route Table IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Transit Router Route Table name.
* `transit_router_id` - (Required, ForceNew) The ID of the Enterprise Edition transit router.
* `transit_router_route_table_type` - (Optional, ForceNew, Available since v1.229.1.) The type of the route table. Valid values: `System`, `Custom`.
* `transit_router_route_table_status` - (Optional, ForceNew) The status of the route table. Valid values: `Creating`, `Active`, `Deleting`.
* `status` - (Optional, ForceNew) The status of the route table. Valid values: `Creating`, `Active`, `Deleting`.
* `transit_router_route_table_ids` - (Optional, ForceNew, List) A list of ID of the CEN Transit Router Route Table.
* `transit_router_route_table_names` - (Optional, ForceNew, List) A list of name of the CEN Transit Router Route Table.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Transit Router Route Table names.
* `tables` - A list of Transit Router Route Tables. Each element contains the following attributes:
  * `id` - The ID of the Transit Router Route Table.
  * `transit_router_route_table_id` - The ID of the Transit Router Route Table.
  * `transit_router_route_table_type` - The type of the route table.
  * `transit_router_route_table_name` - The name of the route table.
  * `transit_router_route_table_description` - The description of the route table.
  * `status` - The status of the route table.
