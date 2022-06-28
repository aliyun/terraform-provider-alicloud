---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_peer"
sidebar_current: "docs-alicloud-resource-vpc-bgp-peer"
description: |-
  Provides a Alicloud VPC Bgp Peer resource.
---

# alicloud\_vpc\_bgp\_peer

Provides a VPC Bgp Peer resource.

For information about VPC Bgp Peer and how to use it, see [What is Bgp Peer](https://www.alibabacloud.com/help/en/doc-detail/91267.html).

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
  virtual_border_router_name = "example_value"
  vlan_id                    = 120
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_vpc_bgp_group" "default" {
  auth_key       = "YourPassword+12345678"
  bgp_group_name = "example_value"
  description    = "example_value"
  local_asn      = 64512
  peer_asn       = 1111
  router_id      = alicloud_express_connect_virtual_border_router.default.id
}

resource "alicloud_vpc_bgp_peer" "default" {
  bfd_multi_hop   = "10"
  bgp_group_id    = alicloud_vpc_bgp_group.default.id
  enable_bfd      = true
  ip_version      = "IPV4"
  peer_ip_address = "1.1.1.1"
}
```

## Argument Reference

The following arguments are supported:

* `bfd_multi_hop` - (Optional) The BFD hop count. Valid values: `1` to `255`. **NOTE:** The attribute is valid when the attribute `enable_bfd` is `true`. The parameter specifies the maximum number of network devices that a packet can traverse from the source to the destination. You can set a proper value based on the factors that affect the physical connection.
* `bgp_group_id` - (Required, ForceNew) The ID of the BGP group.
* `enable_bfd` - (Optional) Specifies whether to enable the Bidirectional Forwarding Detection (BFD) feature.
* `ip_version` - (Optional, ForceNew, Computed) The IP version.
* `peer_ip_address` - (Optional) The IP address of the BGP peer.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bgp Peer.
* `status` - The status of the BGP peer.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Bgp Peer.
* `delete` - (Defaults to 5 mins) Used when delete the Bgp Peer.
* `update` - (Defaults to 5 mins) Used when update the Bgp Peer.

## Import

VPC Bgp Peer can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_bgp_peer.example <id>
```