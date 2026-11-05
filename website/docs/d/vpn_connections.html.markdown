---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_connections"
sidebar_current: "docs-alicloud-datasource-vpn-connections"
description: |-
    Provides a list of VPN connections which owned by an Alicloud account.
---

# alicloud_vpn_connections

The VPN connections data source lists lots of VPN connections resource information owned by an Alicloud account.

-> **NOTE:** Available since v1.18.0.

## Example Usage

```terraform
data "alicloud_vpn_connections" "foo" {
  ids                 = ["fake-conn-id"]
  vpn_gateway_id      = "fake-vpn-id"
  customer_gateway_id = "fake-cgw-id"
  output_file         = "/tmp/vpnconn"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) IDs of the VPN connections.
* `vpn_gateway_id` - (Optional) Use the VPN gateway ID as the search key.
* `customer_gateway_id` - (Optional)Use the VPN customer gateway ID as the search key.
* `name_regex` - (Optional) A regex string of VPN connection name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `ids` - IDs of the VPN connections.
* `names` - Names of the VPN connections.
* `connections` - A list of VPN connections. Each element contains the following attributes:
  * `id` - ID of the VPN connection.
  * `customer_gateway_id` - ID of the VPN customer gateway.
  * `vpn_gateway_id` - ID of the VPN gateway.
  * `name` - The name of the VPN connection.
  * `local_subnet` - The local subnet of the VPN connection.
  * `remote_subnet` - The remote subnet of the VPN connection.
  * `create_time` - The creation time of the VPN connection.
  * `effect_immediately` - Whether the VPN connection takes effect immediately.
  * `status` - The status of the VPN connection, valid value:ike_sa_not_established, ike_sa_established, ipsec_sa_not_established, ipsec_sa_established.
  * `enable_dpd` - Specifies whether to enable the dead peer detection (DPD) feature.
  * `enable_nat_traversal` - Specifies whether to enable NAT traversal.
  * `enable_tunnels_bgp` - Enable tunnel bgp.
  * `ike_config` - (Optional) The configurations of phase-one negotiation.
    * `psk` - (Optional) Used for authentication between the IPsec VPN gateway and the customer gateway.
    * `ike_version` - (Optional) The version of the IKE protocol.
    * `ike_mode` - (Optional) The negotiation mode of IKE phase-one.
    * `ike_enc_alg` - (Optional) The encryption algorithm of phase-one negotiation.
    * `ike_auth_alg` - (Optional) The authentication algorithm of phase-one negotiation.
    * `ike_pfs` - (Optional) The Diffie-Hellman key exchange algorithm used by phase-one negotiation.
    * `ike_lifetime` - (Optional) The SA lifecycle as the result of phase-one negotiation.
    * `ike_local_id` - (Optional) The identification of the VPN gateway.
    * `ike_remote_id` - (Optional) The identification of the customer gateway.
  * `ipsec_config` - (Optional) The configurations of phase-two negotiation.
    * `ipsec_enc_alg` - (Optional) The encryption algorithm of phase-two negotiation.
    * `ipsec_auth_alg` - (Optional) The authentication algorithm of phase-two negotiation.
    * `ipsec_pfs` - (Optional) The Diffie-Hellman key exchange algorithm used by phase-two negotiation.
    * `ipsec_lifetime` - (Optional) The SA lifecycle as the result of phase-two negotiation.
  * `vco_health_check` - (Optional) The health check configuration information.
    * `status` - (Optional) The health check status. Valid values: `success`, `failed`.
    * `dip` - (Optional) The destination ip address.
    * `sip` - (Optional) The source ip address.
    * `interval` - (Optional) The time interval between health checks.
    * `retry` - (Optional) The number of retries for health checks issued.
    * `enable` - (Optional) The health check on status. Valid values: `true`, `false`.
  * `vpn_bgp_config` - (Optional) The configuration information for BGP.
    * `status` - (Optional) The negotiation status of the BGP routing protocol. Valid values: `success`, `false`.
    * `peer_bgp_ip` - (Optional) The BGP address on the other side.
    * `peer_asn` - (Optional) The counterpart autonomous system number.
    * `local_asn` - (Optional) The ali cloud side autonomous system.
    * `auth_key` - (Optional) The authentication keys for BGP routing protocols.
    * `tunnel_cidr` - (Optional) The ipsec tunnel segments.
    * `local_bgp_ip` - (Optional) The ali cloud side BGP address.
  * `tunnel_options_specification` - The tunnel_options_specification supports the following:
    * `customer_gateway_id` - The ID of the customer gateway in Tunnel.
    * `enable_dpd` - Wether enable Dpd detection.
    * `enable_nat_traversal` - enable nat traversal.
    * `internet_ip` - The internet IP of the tunnel.
    * `role` - The role of Tunnel.
    * `state` - The state of the tunnel.
    * `status` - The status of the tunnel.
    * `tunnel_id` - The ID of the tunnel.
    * `zone_no` - The zone of the tunnel.
    * `tunnel_bgp_config` - The bgp config of Tunnel.
      * `local_asn` - Local asn.
      * `local_bgp_ip` - Local bgp IP.
      * `tunnel_cidr` - BGP Tunnel CIDR.
      * `bgp_status` - The BGP status.
      * `peer_bgp_ip` - The peer BGP IP.
      * `peer_asn` - The peer ASN.
    * `tunnel_ike_config` - The configuration of Phase 1 negotiations in Tunnel.
      * `ike_auth_alg` - IKE auth Algorithm.
      * `ike_enc_alg` - IKE encript algorithm.
      * `ike_lifetime` - IKE lifetime.
      * `ike_mode` - IKE Mode.
      * `ike_pfs` - DH Group.
      * `ike_version` - IKE Version.
      * `local_id` - The local Id.
      * `psk` - Preshared secret key.
      * `remote_id` - Remote ID.
    * `tunnel_ipsec_config` - IPsec configuration in Tunnel.
      * `ipsec_auth_alg` - IPsec Auth algorithm.
      * `ipsec_enc_alg` - IPsec Encript algorithm.
      * `ipsec_lifetime` - IPsec lifetime.
      * `ipsec_pfs` - DH Group.
