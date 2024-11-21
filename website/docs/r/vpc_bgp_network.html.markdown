---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_bgp_network"
sidebar_current: "docs-alicloud-resource-vpc-bgp-network"
description: |-
  Provides a Alicloud VPC Bgp Network resource.
---

# alicloud_vpc_bgp_network

Provides a VPC Bgp Network resource.

For information about VPC Bgp Network and how to use it, see [What is Bgp Network](https://www.alibabacloud.com/help/en/doc-detail/91267.html).

-> **NOTE:** Available since v1.153.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_bgp_network&exampleId=7e2fe32d-0fc7-43d0-03af-5b8de718bb3db4437e83&activeTab=example&spm=docs.r.vpc_bgp_network.0.7e2fe32d0f&intl_lang=EN_US" target="_blank">
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

resource "alicloud_vpc_bgp_network" "example" {
  dst_cidr_block = "192.168.0.0/24"
  router_id      = alicloud_express_connect_virtual_border_router.example.id
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Bgp Network.
* `delete` - (Defaults to 1 mins) Used when delete the Bgp Network.

## Import

VPC Bgp Network can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_bgp_network.example <router_id>:<dst_cidr_block>
```