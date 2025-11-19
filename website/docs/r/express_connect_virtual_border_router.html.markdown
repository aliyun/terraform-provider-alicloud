---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_virtual_border_router"
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
  default = "terraform-example"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}

resource "random_integer" "default" {
  min = 1
  max = 2999
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.default.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = random_integer.default.id
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Optional, Int) The bandwidth of the VBR instance. Unit: Mbps. Valid values:
  - When creating a VBR instance for an exclusive leased line, the values are `50`, `100`, `200`, `300`, `400`, `500`, `1000`, `2048`, `5120`, `8192`, `10240`, `20480`, `40960`, `50120`, `61440`, and `102400`.
  - When creating a VBR instance for a shared line, you do not need to configure it. The bandwidth of the VBR is the bandwidth set when creating a shared physical line.
* `circuit_code` - (Optional) The circuit code provided by the operator for the physical connection.
* `description` - (Optional) The description information of the VBR.
* `detect_multiplier` - (Optional, Int) Multiple of detection time.
  That is the maximum number of connection packet losses allowed by the receiver to send messages, which is used to detect whether the link is normal.
  Valid values: `3` to `10`.
* `enable_ipv6` - (Optional, Bool) Whether IPv6 is enabled.
  - `true`: on.
  - `false`: closed.
* `local_gateway_ip` - (Required) The IPv4 address on the Alibaba Cloud side of the VBR instance.
* `local_ipv6_gateway_ip` - (Optional) The IPv6 address on the Alibaba Cloud side of the VBR instance.
* `min_rx_interval` - (Optional, Int) Configure the receiving interval of BFD packets. Valid values: `200` to `1000`.
* `min_tx_interval` - (Optional, Int) Configure the sending interval of BFD packets. Valid values: `200` to `1000`.
* `mtu` - (Optional, Int, Available since v1.263.0) Maximum transmission unit.
* `peer_gateway_ip` - (Required) The IPv4 address of the client side of the VBR instance.
* `peer_ipv6_gateway_ip` - (Optional) The IPv6 address of the client side of the VBR instance.
* `peering_ipv6_subnet_mask` - (Optional) The subnet masks of the Alibaba Cloud-side IPv6 and the customer-side IPv6 of The VBR instance.
* `peering_subnet_mask` - (Required) The subnet masks of the Alibaba Cloud-side IPv4 and the customer-side IPv4 of The VBR instance.
* `physical_connection_id` - (Required, ForceNew) The ID of the physical connection to which the VBR belongs.
* `resource_group_id` - (Optional, Available since v1.263.0) The ID of the resource group.
* `sitelink_enable` - (Optional, Bool, Available since v1.263.0) Whether to allow inter-IDC communication. Valid values: `true`, `false`.
* `status` - (Optional) The status of the VBR. Valid values: `active`, `terminated`.
* `tags` - (Optional, Map, Available since v1.263.0) The tag of the resource.
* `vbr_owner_id` - (Optional) The account ID of the VBR instance owner. The default value is the logon Alibaba Cloud account ID.
* `virtual_border_router_name` - (Optional) The name of the VBR instance.
* `vlan_id` - (Required, Int) The VLAN ID of the VBR instance. Valid values: `0` to `2999`.
* `associated_physical_connections` - (Deprecated since v1.263.0) Field `associated_physical_connections` has been deprecated from provider version 1.263.0. Please use the resource `alicloud_express_connect_vbr_pconn_association` instead.
* `include_cross_account_vbr` - (Removed since v1.263.0) Field `include_cross_account_vbr` has been removed from provider version 1.263.0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - (Available since v1.263.0) The creation time of the VBR.
* `route_table_id` - (Available since v1.166.0) The Route Table ID Of the Virtual Border Router.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Virtual Border Router.
* `delete` - (Defaults to 5 mins) Used when delete the Virtual Border Router.
* `update` - (Defaults to 5 mins) Used when update the Virtual Border Router.

## Import

Express Connect Virtual Border Router can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_virtual_border_router.example <id>
```
