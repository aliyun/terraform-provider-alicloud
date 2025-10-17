---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_vbr_pconn_association"
description: |-
  Provides a Alicloud Express Connect Vbr Pconn Association resource.
---

# alicloud_express_connect_vbr_pconn_association

Provides a Express Connect Vbr Pconn Association resource.

VBR multi-pconn Association.

For information about Express Connect Vbr Pconn Association and how to use it, see [What is Vbr Pconn Association](https://www.alibabacloud.com/help/en/express-connect/latest/associatephysicalconnectiontovirtualborderrouter#doc-api-Vpc-AssociatePhysicalConnectionToVirtualBorderRouter).

-> **NOTE:** Available since v1.196.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_vbr_pconn_association&exampleId=09686d90-b700-a2fa-f87e-512fc88726f36d23fbc2&activeTab=example&spm=docs.r.express_connect_vbr_pconn_association.0.09686d90b7&intl_lang=EN_US" target="_blank">
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

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.example.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = 110
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
  physical_connection_id   = data.alicloud_express_connect_physical_connections.example.connections.1.id
  vbr_id                   = alicloud_express_connect_virtual_border_router.default.id
  peering_subnet_mask      = "255.255.255.252"
  vlan_id                  = "1122"
  enable_ipv6              = true
  local_ipv6_gateway_ip    = "2408:4004:cc::3"
  peer_ipv6_gateway_ip     = "2408:4004:cc::4"
  peering_ipv6_subnet_mask = "2408:4004:cc::/56"
}
```

## Argument Reference

The following arguments are supported:
* `enable_ipv6` - (Optional, ForceNew, Computed) Whether IPv6 is enabled. Value:
  - `true`: on.
  - `false` (default): Off.
* `local_gateway_ip` - (Optional, ForceNew) The Alibaba cloud IP address of the VBR instance.
* `local_ipv6_gateway_ip` - (Optional, ForceNew) The IPv6 address on the Alibaba Cloud side of the VBR instance.
* `peer_gateway_ip` - (Optional, ForceNew) The client IP address of the VBR instance.
  - This attribute only allows the VBR owner to specify or modify.
  - Required when creating a VBR instance for the physical connection owner.
* `peer_ipv6_gateway_ip` - (Optional, ForceNew) The IPv6 address of the client side of the VBR instance.
  - This attribute only allows the VBR owner to specify or modify.
  - Required when creating a VBR instance for the physical connection owner.
* `peering_ipv6_subnet_mask` - (Optional, ForceNew) The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.
Two IPv6 addresses must be in the same subnet.
* `peering_subnet_mask` - (Optional, ForceNew) The subnet mask of the Alibaba Cloud side and the client side of the VBR instance.
The two IP addresses must be in the same subnet.
* `physical_connection_id` - (Required, ForceNew) The ID of the leased line instance.
* `vbr_id` - (Required, ForceNew) The ID of the VBR instance.
* `vlan_id` - (Required, ForceNew) VLAN ID of the VBR. Valid values: **0 to 2999**.

-> **NOTE:**  only the owner of the physical connection can specify this parameter. The VLAN ID of two VBRs under the same physical connection cannot be the same.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<vbr_id>:<physical_connection_id>`.
* `status` - The status of the resource
* `circuit_code` - (Optional, ForceNew, Computed) The circuit code provided by the operator for the physical connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vbr Pconn Association.
* `delete` - (Defaults to 5 mins) Used when delete the Vbr Pconn Association.

## Import

Express Connect Vbr Pconn Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_vbr_pconn_association.example <vbr_id>:<physical_connection_id>
```