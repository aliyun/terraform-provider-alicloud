---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpn_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit-router-vpn-attachment"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Vpn Attachment resource.
---

# alicloud\_cen\_transit\_router\_vpn\_attachment

Provides a Cloud Enterprise Network (CEN) Transit Router Vpn Attachment resource.

For information about Cloud Enterprise Network (CEN) Transit Router Vpn Attachment and how to use it, see [What is Transit Router Vpn Attachment](https://help.aliyun.com/document_detail/443993.html).

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
```

## Argument Reference

The following arguments are supported:

* `auto_publish_route_enabled` - (Optional, Computed) Whether to allow the forwarding router instance to automatically publish routing entries to IPsec connections.
* `cen_id` - (Optional, ForceNew) The id of the cen.
* `transit_router_attachment_description` - (Optional) The description of the VPN connection. The description can contain `2` to `256` characters. The description must start with English letters, but cannot start with `http://` or `https://`.
* `transit_router_attachment_name` - (Optional) The name of the VPN connection. The name must be `2` to `128` characters in length, and can contain digits, underscores (_), and hyphens (-). It must start with a letter.
* `transit_router_id` - (Required, ForceNew) The ID of the forwarding router instance.
* `vpn_id` - (Required, ForceNew) The id of the vpn.
* `vpn_owner_id` - (Optional, Computed, ForceNew) The owner id of vpn.
* `zone` - (Required, ForceNew) The list of zone mapping. See the following `Block zone`.
* `tags` - (Optional, Available in v1.193.1+) A mapping of tags to assign to the resource.

#### zone Block

The `zone` supports the following:

* `zone_id` - (Required, ForceNew) The id of the zone.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router Vpn Attachment.
* `status` - The associating status of the network.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 40 mins) Used when create the Transit Router Vpn Attachment.
* `update` - (Defaults to 1 mins) Used when update the Transit Router Vpn Attachment.
* `delete` - (Defaults to 30 mins) Used when delete the Transit Router Vpn Attachment.

## Import

Cloud Enterprise Network (CEN) Transit Router Vpn Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_vpn_attachment.example <id>
```
