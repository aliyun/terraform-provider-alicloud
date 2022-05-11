---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_virtual_border_router"
sidebar_current: "docs-alicloud-resource-express-connect-virtual-border-router"
description: |-
  Provides a Alicloud Express Connect Virtual Border Router resource.
---

# alicloud\_express\_connect\_virtual\_border\_router

Provides a Express Connect Virtual Border Router resource.

For information about Express Connect Virtual Border Router and how to use it, see [What is Virtual Border Router](https://www.alibabacloud.com/help/en/doc-detail/44854.htm).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^my-PhysicalConnection"
}

resource "alicloud_express_connect_virtual_border_router" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = "example_value"
  vlan_id                    = 1
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
* `detect_multiplier` - (Optional, Computed) Detection time multiplier that recipient allows the sender to send a message of the maximum allowable connections for the number of packets, used to detect whether the link normal. Value: 3~10.
* `enable_ipv6` - (Optional, Computed) Whether to Enable IPv6. Valid values: `false`, `true`.
* `local_gateway_ip` - (Required) Alibaba Cloud-Connected IPv4 address.
* `local_ipv6_gateway_ip` - (Optional) Alibaba Cloud-Connected IPv6 Address.
* `min_rx_interval` - (Optional, Computed) Configure BFD packet reception interval of values include: 200~1000, unit: ms.
* `min_tx_interval` - (Optional, Computed) Configure BFD packet transmission interval maximum value: 200~1000, unit: ms.
* `peer_gateway_ip` - (Required) The Client-Side Interconnection IPv4 Address.
* `peer_ipv6_gateway_ip` - (Optional) The Client-Side Interconnection IPv6 Address.
* `peering_ipv6_subnet_mask` - (Optional) Alibaba Cloud-Connected IPv6 with Client-Side Interconnection IPv6 of Subnet Mask.
* `peering_subnet_mask` - (Required) Alibaba Cloud-Connected IPv4 and Client-Side Interconnection IPv4 of Subnet Mask.
* `physical_connection_id` - (Required, ForceNew) The ID of the Physical Connection to Which the ID.
* `status` - (Optional, Computed) The instance state. Valid values: `active`, `deleting`, `recovering`, `terminated`, `terminating`, `unconfirmed`.
* `vbr_owner_id` - (Optional) The vbr owner id.
* `virtual_border_router_name` - (Optional) The name of VBR. Length is from 2 to 128 characters, must start with a letter or the Chinese at the beginning can contain numbers, the underscore character (_) and dash (-). But do not start with http:// or https:// at the beginning.
* `vlan_id` - (Required) The VLAN ID of the VBR. Value range: 0~2999.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Virtual Border Router.
* `route_table_id` - (Available in v1.166.0+) The Route Table ID Of the Virtual Border Router.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `update` - (Defaults to 2 mins) Used when update the Virtual Border Router.

## Import

Express Connect Virtual Border Router can be imported using the id, e.g.

```
$ terraform import alicloud_express_connect_virtual_border_router.example <id>
```
