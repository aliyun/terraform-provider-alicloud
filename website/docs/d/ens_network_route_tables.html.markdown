---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_route_tables"
sidebar_current: "docs-alicloud-datasource-ens-network-route-tables"
description: |-
  Provides a list of Ens Network Route Tables to the user.
---

# alicloud\_ens\_network\_route\_tables

This data source provides the Ens Network Route Tables of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ens_network_route_tables" "default" {
  ids            = ["example-rt-id"]
  network_id     = "example-network-id"
  associate_type = "Gateway"
  output_file    = "./route_tables.txt"
}
output "ens_route_table_id_1" {
  value = data.alicloud_ens_network_route_tables.default.tables.0.id
}
```

## Argument Reference

The following arguments are supported:

* `associate_type` - (Optional) The binding type of the route table. Valid values: `VSwitch`, `Gateway`.
* `ids` - (Optional) A list of Route Table IDs.
* `network_id` - (Optional) The ID of the network to which the route table belongs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `route_table_id` - (Optional) The ID of the route table.
* `route_table_name` - (Optional) The name of the route table.
* `route_table_type` - (Optional) The type of the route table.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Route Table IDs.
* `tables` - A list of Ens Network Route Tables. Each element contains the following attributes:
  * `associate_type` - The binding type of the route table.
  * `create_time` - The creation time of the route table.
  * `description` - The description of the route table.
  * `id` - The ID of the route table.
  * `is_default_gateway_route_table` - Whether the route table is the default gateway route table.
  * `network_id` - The ID of the network to which the route table belongs.
  * `route_table_id` - The ID of the route table.
  * `route_table_name` - The name of the route table.
  * `route_table_type` - The type of the route table.
  * `status` - The status of the route table.
