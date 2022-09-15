---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_vco_routes"
sidebar_current: "docs-alicloud-datasource-vpn-gateway-vco-routes"
description: |-
  Provides a list of Vpn Gateway Vco Routes to the user.
---

# alicloud\_vpn\_gateway\_vco\_routes

This data source provides the Vpn Gateway Vco Routes of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.183.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpn_gateway_vco_routes" "ids" {}
output "vpn_gateway_vco_route_id_1" {
  value = data.alicloud_vpn_gateway_vco_routes.ids.routes.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Vco Route IDs.
* `vpn_connection_id` - (Required, ForceNew) The id of the vpn connection.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `route_entry_type` - (Optional, ForceNew) The Routing input type. Valid values: `custom`, `bgp`.
* `status` - (Optional, ForceNew) The status of the vpn route entry. Valid values: `normal`, `published`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `routes` - A list of Vpn Gateway Vco Routes. Each element contains the following attributes:
  * `as_path` - List of autonomous system numbers through which BGP routing entries pass.
  * `create_time` - The creation time of the VPN destination route.
  * `source` - The source CIDR block of the destination route.
  * `status` - The status of the vpn route entry.
  * `weight` - The weight value of the destination route.
  * `next_hop` - The next hop of the destination route.
  * `vpn_connection_id` - The id of the vpn connection.
  * `route_dest` - The destination network segment of the destination route.
  * `id` - The ID of the Vpn Gateway Vco Routes.