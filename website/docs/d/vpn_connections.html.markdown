---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_connections"
sidebar_current: "docs-alicloud-datasource-vpn-connections"
description: |-
    Provides a list of VPN connections which owned by an Alicloud account.
---

# alicloud\_vpn_connections

The VPN connections data source lists lots of VPN connections resource information owned by an Alicloud account.

## Example Usage

```
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

* `ids` - (Optional) IDs of the VPN connections.
* `names` - (Optional) names of the VPN connections.
* `connections` - A list of VPN connections. Each element contains the following attributes:
  * `id` - ID of the VPN connection.
  * `customer_gateway_id` - ID of the VPN customer gateway.
  * `vpn_gateway_id` - ID of the VPN gateway.
  * `name` - The name of the VPN connection.
  * `local_subnet` - The local subnet of the VPN connection.
  * `remote_subnet` - The remote subnet of the VPN connection.
  * `status` - The status of the VPN connection, valid value:ike_sa_not_established, ike_sa_established, ipsec_sa_not_established, ipsec_sa_established.
  * `ike_config` - The configurations of phase-one negotiation.
  * `ipsec_config` - The configurations of phase-two negotiation.

  ### Block ike_config

  The ike_config mapping supports the following:

  * `psk` - Used for authentication between the IPsec VPN gateway and the customer gateway.
  * `ike_version` - The version of the IKE protocol. 
  * `ike_mode` - The negotiation mode of IKE phase-one. 
  * `ike_enc_alg` - The encryption algorithm of phase-one negotiation. 
  * `ike_auth_alg` - The authentication algorithm of phase-one negotiation. 
  * `ike_pfs` - The Diffie-Hellman key exchange algorithm used by phase-one negotiation. 
  * `ike_lifetime` - The SA lifecycle as the result of phase-one negotiation. 
  * `ike_local_id` - The identification of the VPN gateway.
  * `ike_remote_id` - The identification of the customer gateway.

  ### Block ipsec_config

  The ipsec_config mapping supports the following:

  * `ipsec_enc_alg` - The encryption algorithm of phase-two negotiation. 
  * `ipsec_auth_alg` - The authentication algorithm of phase-two negotiation. 
  * `ipsec_pfs` - The Diffie-Hellman key exchange algorithm used by phase-two negotiation. 
  * `ipsec_lifetime` - The SA lifecycle as the result of phase-two negotiation. 
