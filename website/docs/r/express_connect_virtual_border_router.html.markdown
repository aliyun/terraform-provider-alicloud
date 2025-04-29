---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_virtual_border_router"
sidebar_current: "docs-alicloud-resource-express-connect-virtual-border-router"
description: |-
  Provides a Alicloud Express Connect Virtual Border Router resource.
---

# alicloud_express_connect_virtual_border_router

Provides a Express Connect Virtual Border Router resource.

For information about Express Connect Virtual Border Router and how to use it, see [What is Virtual Border Router](https://www.alibabacloud.com/help/en/doc-detail/44854.htm).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_virtual_border_router&exampleId=ebfe3429-41d0-6517-e8ea-0416ffe684ce150dad87&activeTab=example&spm=docs.r.express_connect_virtual_border_router.0.ebfe342941&intl_lang=EN_US" target="_blank">
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
```

## Argument Reference

The following arguments are supported:

* `associated_physical_connections` - (Optional) The associated physical connections.
* `bandwidth` - (Optional) The bandwidth.
* `circuit_code` - (Optional) Operators for physical connection circuit provided coding.
* `description` - (Optional) The description of VBR. Length is from 2 to 256 characters, must start with a letter or the Chinese at the beginning, but not at the http:// Or https:// at the beginning.
* `detect_multiplier` - (Optional) Detection time multiplier that recipient allows the sender to send a message of the maximum allowable connections for the number of packets, used to detect whether the link normal. Value: 3~10.
* `enable_ipv6` - (Optional) Whether to Enable IPv6. Valid values: `false`, `true`.
* `local_gateway_ip` - (Required) Alibaba Cloud-Connected IPv4 address.
* `local_ipv6_gateway_ip` - (Optional) Alibaba Cloud-Connected IPv6 Address.
* `min_rx_interval` - (Optional) Configure BFD packet reception interval of values include: 200~1000, unit: ms.
* `min_tx_interval` - (Optional) Configure BFD packet transmission interval maximum value: 200~1000, unit: ms.
* `peer_gateway_ip` - (Required) The Client-Side Interconnection IPv4 Address.
* `peer_ipv6_gateway_ip` - (Optional) The Client-Side Interconnection IPv6 Address.
* `peering_ipv6_subnet_mask` - (Optional) Alibaba Cloud-Connected IPv6 with Client-Side Interconnection IPv6 of Subnet Mask.
* `peering_subnet_mask` - (Required) Alibaba Cloud-Connected IPv4 and Client-Side Interconnection IPv4 of Subnet Mask.
* `physical_connection_id` - (Required, ForceNew) The ID of the Physical Connection to Which the ID.
* `status` - (Optional) The instance state. Valid values: `active`, `deleting`, `recovering`, `terminated`, `terminating`, `unconfirmed`.
* `vbr_owner_id` - (Optional) The vbr owner id.
* `virtual_border_router_name` - (Optional) The name of VBR. Length is from 2 to 128 characters, must start with a letter or the Chinese at the beginning can contain numbers, the underscore character (_) and dash (-). But do not start with http:// or https:// at the beginning.
* `vlan_id` - (Required) The VLAN ID of the VBR. Value range: 0~2999.
* `include_cross_account_vbr` - (Optional, Available since v1.191.0) Whether cross account border routers are included. Valid values: `false`, `true`. Default: `true`. 

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Virtual Border Router.
* `route_table_id` - (Available since v1.166.0) The Route Table ID Of the Virtual Border Router.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `update` - (Defaults to 2 mins) Used when update the Virtual Border Router.

## Import

Express Connect Virtual Border Router can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_virtual_border_router.example <id>
```
