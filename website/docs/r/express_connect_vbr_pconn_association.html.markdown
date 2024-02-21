---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_vbr_pconn_association"
sidebar_current: "docs-alicloud-resource-express-connect-vbr-pconn-association"
description: |-
  Provides a Alicloud Express Connect Vbr Pconn Association resource.
---

# alicloud_express_connect_vbr_pconn_association

Provides a Express Connect Vbr Pconn Association resource.

For information about Express Connect Vbr Pconn Association and how to use it, see [What is Vbr Pconn Association](https://www.alibabacloud.com/help/en/express-connect/latest/associatephysicalconnectiontovirtualborderrouter#doc-api-Vpc-AssociatePhysicalConnectionToVirtualBorderRouter).

-> **NOTE:** Available since v1.196.0.

## Example Usage

Basic Usage

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
  enable_ipv6                = true
  local_ipv6_gateway_ip      = "2408:4004:cc:400::1"
  peer_ipv6_gateway_ip       = "2408:4004:cc:400::2"
  peering_ipv6_subnet_mask   = "2408:4004:cc:400::/56"
}

resource "alicloud_express_connect_vbr_pconn_association" "example" {
  peer_gateway_ip          = "10.0.0.6"
  local_gateway_ip         = "10.0.0.5"
  physical_connection_id   = data.alicloud_express_connect_physical_connections.example.connections.2.id
  vbr_id                   = alicloud_express_connect_virtual_border_router.example.id
  peering_subnet_mask      = "255.255.255.252"
  vlan_id                  = random_integer.vlan_id.id + 2
  enable_ipv6              = true
  local_ipv6_gateway_ip    = "2408:4004:cc::3"
  peer_ipv6_gateway_ip     = "2408:4004:cc::4"
  peering_ipv6_subnet_mask = "2408:4004:cc::/56"
}
```

## Argument Reference

The following arguments are supported:
* `enable_ipv6` - (Optional, ForceNew) Whether IPv6 is enabled. Value:
  - **true**: on.
  - **false** (default): Off.
* `local_gateway_ip` - (Optional, ForceNew) The Alibaba cloud IP address of the VBR instance.
* `local_ipv6_gateway_ip` - (Optional, ForceNew) The IPv6 address on the Alibaba Cloud side of the VBR instance.
* `peer_gateway_ip` - (Optional, ForceNew) The client IP address of the VBR instance. This attribute only allows the VBR owner to specify or modify. **NOTE:** Required when creating a VBR instance for the physical connection owner.
* `peer_ipv6_gateway_ip` - (Optional, ForceNew) The IPv6 address of the client side of the VBR instance. This attribute only allows the VBR owner to specify or modify. **NOTE:** Required when creating a VBR instance for the physical connection owner.
* `peering_ipv6_subnet_mask` - (Optional, ForceNew) The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.Two IPv6 addresses must be in the same subnet.
* `peering_subnet_mask` - (Optional, ForceNew) The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.The two IP addresses must be in the same subnet.
* `physical_connection_id` - (Required, ForceNew) The ID of the leased line instance.
* `vbr_id` - (Required, ForceNew) The ID of the VBR instance.
* `vlan_id` - (Required, ForceNew) VLAN ID of the VBR. Valid values: **0 to 2999**. **NOTE:** only the owner of the physical connection can specify this parameter. The VLAN ID of two VBRs under the same physical connection cannot be the same.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `vbr_id:physical_connection_id`.
* `status` - The status of the resource.
* `circuit_code` - The circuit code provided by the operator for the physical connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vbr Pconn Association.
* `delete` - (Defaults to 5 mins) Used when delete the Vbr Pconn Association.

## Import

Express Connect Vbr Pconn Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_vbr_pconn_association.example <VbrId>:<PhysicalConnectionId>
```