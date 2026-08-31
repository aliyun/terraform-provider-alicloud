---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_route_table"
description: |-
  Provides a Alicloud ENS Network Route Table resource.
---

# alicloud_ens_network_route_table

Provides a ENS Network Route Table resource.

Routing table of ENS network.

For information about ENS Network Route Table and how to use it, see [What is Network Route Table](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateEnsRouteTable).

-> **NOTE:** Available since v1.291.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
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
```

## Argument Reference

The following arguments are supported:
* `network_id` - (Required, ForceNew) The network ID.
* `associate_type` - (Required, ForceNew) The binding type of the created routing table. Value: `Gateway`.
* `route_table_name` - (Optional) The name of the routing table. The length is 1 to 128 characters, but it cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the routing table. The length is 1 to 256 characters, but it cannot start with `http://` or `https://`.
* `is_default_gateway_route_table` - (Computed) Whether it is the default gateway route table.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. Same as `route_table_id`.
* `route_table_id` - The ID of the routing table.
* `route_table_type` - The type of the routing table. Value: `Custom`, `System`.
* `status` - The status of the routing table.
* `create_time` - The creation time of the routing table.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Route Table.
* `delete` - (Defaults to 5 mins) Used when delete the Network Route Table.
* `update` - (Defaults to 5 mins) Used when update the Network Route Table.

## Import

ENS Network Route Table can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_network_route_table.example <route_table_id>
```
