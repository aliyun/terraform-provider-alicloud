---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_route_table"
description: |-
  Provides a Alicloud ENS Network Route Table resource.
---

# alicloud_ens_network_route_table

Provides a ENS Network Route Table resource.



For information about ENS Network Route Table and how to use it, see [What is Network Route Table](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateEnsRouteTable).

-> **NOTE:** Available since v1.287.0.

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
  default = "cn-hangzhou-58"
}

resource "alicloud_ens_network" "defaultasem9i" {
  network_name  = "tf-example-network-route-table"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}


resource "alicloud_ens_network_route_table" "default" {
  description                    = "example1"
  associate_type                 = "Gateway"
  network_id                     = alicloud_ens_network.defaultasem9i.id
  route_table_name               = "tf-example-route-table"
  is_default_gateway_route_table = true
}
```

## Argument Reference

The following arguments are supported:
* `associate_type` - (Required, ForceNew) The binding type of the created routing table. Value:
● VSwitch (default): can be used to bind switches.
● Gateway: can be used to bind a Gateway.
* `description` - (Optional) Description of the routing table.
The name must be 1 to 256 characters in length and cannot start with http:// or https.
* `is_default_gateway_route_table` - (Optional, ForceNew) Is the default gateway route table.
* `network_id` - (Required, ForceNew) The network ID.
* `route_table_name` - (Optional) Name of the routing table.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `route_table_type` - The type of the routing table.
* `status` - Status.

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