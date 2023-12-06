---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_peer"
sidebar_current: "docs-alicloud-resource-vpc-bgp-peer"
description: |-
  Provides a Alicloud VPC Bgp Peer resource.
---

# alicloud_vpc_bgp_peer

Provides a VPC Bgp Peer resource.

For information about VPC Bgp Peer and how to use it, see [What is Bgp Peer](https://www.alibabacloud.com/help/en/express-connect/developer-reference/api-vpc-2016-04-28-createbgppeer-efficiency-channels).

-> **NOTE:** Available since v1.153.0.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

provider "alicloud" {
  region = var.region
}

data "alicloud_express_connect_physical_connections" "example" {
  name_regex = "^preserved-NODELETING"
}

resource "random_integer" "vlan_id" {
  max = 2999
  min = 1
}

resource "alicloud_express_connect_virtual_border_router" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.example.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = random_integer.vlan_id.id
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_vpc_bgp_group" "example" {
  router_id      = alicloud_express_connect_virtual_border_router.example.id
  peer_asn       = 1111
  is_fake_asn    = true
  auth_key       = "YourPassword+12345678"
  bgp_group_name = var.name
  description    = var.name
}

resource "alicloud_vpc_bgp_peer" "example" {
  bgp_group_id    = alicloud_vpc_bgp_group.example.id
  peer_ip_address = "1.1.1.1"
  ip_version      = "IPV4"
  enable_bfd      = true
  bfd_multi_hop   = "10"
}
```

## Argument Reference

The following arguments are supported:

* `bgp_group_id` - (Required, ForceNew) The ID of the BGP group.
* `peer_ip_address` - (Optional) The IP address of the Bgp Peer.
* `ip_version` - (Optional, ForceNew) The IP version.
* `enable_bfd` - (Optional, Bool) Specifies whether to enable the Bidirectional Forwarding Detection (BFD) feature.
* `bfd_multi_hop` - (Optional, Int) The BFD hop count. Valid values: `1` to `255`. **NOTE:** The attribute is valid when the attribute `enable_bfd` is `true`. The parameter specifies the maximum number of network devices that a packet can traverse from the source to the destination. You can set a proper value based on the factors that affect the physical connection.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bgp Peer.
* `status` - The status of the Bgp Peer.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Bgp Peer.
* `update` - (Defaults to 5 mins) Used when update the Bgp Peer.
* `delete` - (Defaults to 5 mins) Used when delete the Bgp Peer.

## Import

VPC Bgp Peer can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_bgp_peer.example <id>
```
