---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_route_tables"
description: |-
  Provides a list of ENS Network Route Tables to the user.
---

# alicloud_ens_network_route_tables

This data source provides the ENS Network Route Tables of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.291.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

variable "ens_region_id" {
  default = "cn-hangzhou-31"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  description   = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_network_route_table" "default" {
  network_id       = alicloud_ens_network.default.id
  associate_type   = "Gateway"
  route_table_name = var.name
  description      = var.name
}

data "alicloud_ens_network_route_tables" "default" {
  route_table_id = alicloud_ens_network_route_table.default.route_table_id
}
```

## Argument Reference

The following arguments are supported:
* `route_table_id` - (Optional) The ID of the routing table.
* `route_table_name` - (Optional) The name of the routing table.
* `network_id` - (Optional) The network ID.
* `associate_type` - (Optional) The binding type of the routing table. Value: `VSwitch`, `Gateway`.
* `route_table_type` - (Optional) The type of the routing table. Value: `Custom`, `System`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `names` - A list of routing table names.
* `ids` - A list of routing table IDs.
* `tables` - A list of routing tables. Each element contains the following attributes:
  * `id` - The ID of the routing table. Same as `route_table_id`.
  * `route_table_id` - The ID of the routing table.
  * `route_table_name` - The name of the routing table.
  * `description` - The description of the routing table.
  * `associate_type` - The binding type of the routing table.
  * `network_id` - The network ID.
  * `route_table_type` - The type of the routing table.
  * `status` - The status of the routing table.
  * `is_default_gateway_route_table` - Whether it is the default gateway route table.
  * `create_time` - The creation time of the routing table.
