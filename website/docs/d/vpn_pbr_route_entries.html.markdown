---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_pbr_route_entries"
sidebar_current: "docs-alicloud-datasource-vpn-pbr-route-entries"
description: |-
    Provides a list of VPN Pbr Route Entries which owned by an Alicloud account.
---

-> **NOTE:** Available in v1.162.0+.

# alicloud\_vpn\_pbr\_route\_entries

The data source lists a number of VPN Pbr Route Entries resource information owned by an Alicloud account.

## Example Usage

```terraform
data "alicloud_vpn_pbr_route_entries" "ids" {
  vpn_gateway_id = "example_vpn_gateway_id"
  ids            = ["example_id"]
}
output "vpn_ipsec_server_id_1" {
  value = data.alicloud_vpn_pbr_route_entries.ids.entries.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of VPN Pbr Route Entries IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vpn_gateway_id` - (Required, ForceNew) The ID of the VPN gateway.

## Attributes Reference

The following attributes are exported:

* `entries` - A list of VPN Pbr Route Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the VPN Pbr Route Entry. 
  * `id` - The id of the vpn pbr route entry. The value formats as `<vpn_gateway_id>:<next_hop>:<route_source>:<route_dest>`.
  * `vpn_gateway_id` - The ID of the vpn gateway.
  * `next_hop` - The next hop of the policy-based route.
  * `route_source` - The source CIDR block of the policy-based route.
  * `route_dest` - The destination CIDR block of the policy-based route.
  * `weight` - The weight of the policy-based route. Valid values: 0 and 100.
  * `status` - The status of the VPN Pbr Route Entry.
