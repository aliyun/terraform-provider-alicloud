---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_vco_route"
sidebar_current: "docs-alicloud-resource-vpn-gateway-vco-route"
description: |-
  Provides a Alicloud VPN Gateway Vco Route resource.
---

# alicloud\_vpn\_gateway\_vco\_route

Provides a VPN Gateway Vco Route resource.

For information about VPN Gateway Vco Route and how to use it, see [What is Vco Route](https://www.alibabacloud.com/help/zh/virtual-private-cloud/latest/createvcorouteentry).

-> **NOTE:** Available in v1.183.0+.

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
```

## Argument Reference

The following arguments are supported:

* `weight` - (Required, ForceNew) The weight value of the destination route. Valid values: `0`, `100`.
* `next_hop` - (Required, ForceNew) The next hop of the destination route.
* `vpn_connection_id` - (Required, ForceNew) The id of the vpn connection.
* `route_dest` - (Required, ForceNew) The destination network segment of the destination route.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Vco Route. The value formats as `<vpn_connection_id>:<route_dest>:<next_hop>:<weight>`.
* `status` - The status of the vpn route entry.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the VPN Gateway Vco Route.
* `delete` - (Defaults to 1 mins) Used when deleting the VPN Gateway Vco Route.


## Import

VPN Gateway Vco Route can be imported using the id, e.g.

```
$ terraform import alicloud_vpn_gateway_vco_route.example <vpn_connection_id>:<route_dest>:<next_hop>:<weight>
```