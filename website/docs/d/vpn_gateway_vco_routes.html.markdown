---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_vco_routes"
sidebar_current: "docs-alicloud-datasource-vpn-gateway-vco-routes"
description: |-
  Provides a list of Vpn Gateway Vco Routes to the user.
---

# alicloud_vpn_gateway_vco_routes

This data source provides the Vpn Gateway Vco Routes of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.183.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}
resource "alicloud_cen_transit_router" "default" {
  cen_id                     = alicloud_cen_instance.default.id
  transit_router_description = "desd"
  transit_router_name        = var.name
}
data "alicloud_cen_transit_router_available_resources" "default" {}
resource "alicloud_vpn_customer_gateway" "default" {
  name        = "${var.name}"
  ip_address  = "42.104.22.210"
  asn         = "45014"
  description = "testAccVpnConnectionDesc"
}
resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  network_type        = "public"
  local_subnet        = "0.0.0.0/0"
  remote_subnet       = "0.0.0.0/0"
  effect_immediately  = false
  ike_config {
    ike_auth_alg = "md5"
    ike_enc_alg  = "des"
    ike_version  = "ikev2"
    ike_mode     = "main"
    ike_lifetime = 86400
    psk          = "tf-testvpn2"
    ike_pfs      = "group1"
    remote_id    = "testbob2"
    local_id     = "testalice2"
  }
  ipsec_config {
    ipsec_pfs      = "group5"
    ipsec_enc_alg  = "des"
    ipsec_auth_alg = "md5"
    ipsec_lifetime = 86400
  }
  bgp_config {
    enable       = true
    local_asn    = 45014
    tunnel_cidr  = "169.254.11.0/30"
    local_bgp_ip = "169.254.11.1"
  }
  health_check_config {
    enable   = true
    sip      = "192.168.1.1"
    dip      = "10.0.0.1"
    interval = 10
    retry    = 10
    policy   = "revoke_route"
  }
  enable_dpd           = true
  enable_nat_traversal = true
  vpn_attachment_name  = var.name
}
resource "alicloud_cen_transit_router_vpn_attachment" "default" {
  auto_publish_route_enabled            = false
  transit_router_attachment_description = var.name
  transit_router_attachment_name        = var.name
  cen_id                                = alicloud_cen_transit_router.default.cen_id
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  vpn_id                                = alicloud_vpn_gateway_vpn_attachment.default.id
  zone {
    zone_id = data.alicloud_cen_transit_router_available_resources.default.resources.0.master_zones.0
  }
}

resource "alicloud_vpn_gateway_vco_route" "default" {
  route_dest        = "192.168.12.0/24"
  next_hop          = alicloud_cen_transit_router_vpn_attachment.default.vpn_id
  vpn_connection_id = alicloud_cen_transit_router_vpn_attachment.default.vpn_id
  weight            = 100
}

data "alicloud_vpn_gateway_vco_routes" "default" {
  vpn_connection_id = alicloud_cen_transit_router_vpn_attachment.default.vpn_id
}
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

## Attributes Reference

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