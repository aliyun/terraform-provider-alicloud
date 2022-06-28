---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_network"
sidebar_current: "docs-alicloud-resource-vpc-bgp-network"
description: |-
  Provides a Alicloud VPC Bgp Network resource.
---

# alicloud\_vpc\_bgp\_network

Provides a VPC Bgp Network resource.

For information about VPC Bgp Network and how to use it, see [What is Bgp Network](https://www.alibabacloud.com/help/en/doc-detail/91267.html).

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_express_connect_physical_connections" "default" {}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.default.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = 120
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_vpc_bgp_network" "example" {
  dst_cidr_block = "example_value"
  router_id      = alicloud_express_connect_virtual_border_router.default.id
}
```

## Argument Reference

The following arguments are supported:

* `dst_cidr_block` - (Required, ForceNew) The CIDR block of the virtual private cloud (VPC) or vSwitch that you want to connect to a data center.
* `router_id` - (Required, ForceNew) The ID of the vRouter associated with the router interface.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Bgp Network. The value formats as `<router_id>:<dst_cidr_block>`.
* `status` - The state of the advertised BGP network.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Bgp Network.
* `delete` - (Defaults to 1 mins) Used when delete the Bgp Network.

## Import

VPC Bgp Network can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_bgp_network.example <router_id>:<dst_cidr_block>
```