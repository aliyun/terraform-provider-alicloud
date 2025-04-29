---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_peer"
description: |-
  Provides a Alicloud Express Connect Bgp Peer resource.
---

# alicloud_vpc_bgp_peer

Provides a Express Connect Bgp Peer resource. 

For information about VPC Bgp Peer and how to use it, see [What is Bgp Peer](https://www.alibabacloud.com/help/en/doc-detail/91267.html).

-> **NOTE:** Available since v1.153.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_bgp_peer&exampleId=9347c15d-7dfc-00dd-f5c0-0b26b2c8e9f9cbea2ca6&activeTab=example&spm=docs.r.vpc_bgp_peer.0.9347c15d7d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
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
  auth_key       = "YourPassword+12345678"
  bgp_group_name = var.name
  description    = var.name
  peer_asn       = 1111
  router_id      = alicloud_express_connect_virtual_border_router.example.id
  is_fake_asn    = true
}

resource "alicloud_vpc_bgp_peer" "example" {
  bfd_multi_hop   = "10"
  bgp_group_id    = alicloud_vpc_bgp_group.example.id
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
* `ip_version` - (Optional, ForceNew) The IP version.
* `peer_ip_address` - (Optional) The IP address of the BGP peer.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `bgp_peer_name` - The name of the BGP neighbor.
* `status` - Status of BGP neighbors.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Bgp Peer.
* `delete` - (Defaults to 5 mins) Used when delete the Bgp Peer.
* `update` - (Defaults to 5 mins) Used when update the Bgp Peer.

## Import

Express Connect Bgp Peer can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_bgp_peer.example <id>
```