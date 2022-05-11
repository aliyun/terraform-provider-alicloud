---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_route_entry"
sidebar_current: "docs-alicloud-resource-vpn-route-entry"
description: |-
  Provides a Alicloud VPN Route Entry resource.
---

# alicloud\_vpn_route_entry

Provides a VPN Route Entry resource.

-> **NOTE:** Terraform will build vpn route entry instance while it uses `alicloud_vpn_route_entry` to build a VPN Route Entry resource.

-> **NOTE:** Available in 1.57.0+.

For information about VPN Route Entry and how to use it, see [What is VPN Route Entry](https://www.alibabacloud.com/help/en/doc-detail/127250.html).


## Example Usage

Basic Usage

```
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name       = "tf_test"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  name              = "tf_test"
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "10.1.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_vpn_gateway" "default" {
  name                 = "tf_vpn_gateway_test"
  vpc_id               = alicloud_vpc.default.id
  bandwidth            = 10
  instance_charge_type = "PayByTraffic"
  enable_ssl           = false
  vswitch_id           = alicloud_vswitch.default.id
}

resource "alicloud_vpn_connection" "default" {
  name                = "tf_vpn_connection_test"
  customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  vpn_gateway_id      = alicloud_vpn_gateway.default.id
  local_subnet        = ["192.168.2.0/24"]
  remote_subnet       = ["192.168.3.0/24"]
}

resource "alicloud_vpn_customer_gateway" "default" {
  name       = "tf_customer_gateway_test"
  ip_address = "192.168.1.1"
}

resource "alicloud_vpn_route_entry" "default" {
  vpn_gateway_id = alicloud_vpn_gateway.default.id
  route_dest     = "10.0.0.0/24"
  next_hop       = alicloud_vpn_connection.default.id
  weight         = 0
  publish_vpc    = false
}
```
## Argument Reference

The following arguments are supported:

* `vpn_gateway_id` - (Required, ForceNew) The id of the vpn gateway.
* `next_hop` - (Required, ForceNew) The next hop of the destination route.
* `publish_vpc` - (Required) Whether to issue the destination route to the VPC.
* `route_dest` - (Required, ForceNew) The destination network segment of the destination route.
* `weight` - (Required) The value should be 0 or 100.

## Attributes Reference

The following attributes are exported:

* `id` - The combination id of the vpn route entry.
* `route_entry_type` - (Available in 1.161.0+) The type of the vpn route entry.
* `status` - (Available in 1.161.0+) The status of the vpn route entry.

#### Timeouts

-> **NOTE:** Available in 1.161.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the vpn route entry.
* `update` - (Defaults to 5 mins) Used when update the vpn route entry.
* `delete` - (Defaults to 5 mins) Used when delete the vpn route entry.

## Import

VPN route entry can be imported using the id(VpnGatewayId +":"+ NextHop +":"+ RouteDest), e.g.

```
$ terraform import alicloud_vpn_route_entry.example vpn-abc123456:vco-abc123456:10.0.0.10/24
```
