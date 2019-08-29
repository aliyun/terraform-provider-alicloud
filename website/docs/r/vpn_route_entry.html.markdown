---
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_route_entry"
sidebar_current: "docs-alicloud-resource-vpn-route-entry"
description: |-
  Provides a Alicloud VPN route resource.
---

# alicloud\_vpn_route_entry

Provides a VPN route resource.

-> **NOTE:** Terraform will build vpn route instance  while it uses `alicloud_vpn_route_entry` to build a vpn route resource.

## Example Usage

Basic Usage

```
resource "alicloud_vpn_route_entry" "foo" {
  next_hop             = "vco-bp15oes1py4i66rmd****"
  publish_vpc          = true
  route_dest           = "10.0.0.0/24"
  vpn_gateway_id       = "vpn-gateway-fakeid"
  weight               = 10
}
```
## Argument Reference

The following arguments are supported:

* `vpn_gateway_id` - (Required, ForceNew) The id of the vpn gateway.
* `next_hop` - (Required, ForceNew) The next hop of the destination route.
* `publish_vpc` - (Required) Whether to issue the destination route to the VPC.
* `route_dest` - (Required, ForceNew) The destination network segment of the destination route.
* `weight` - (Required) The value should be 0 or 100.
