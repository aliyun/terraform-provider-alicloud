---
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
data "alicloud_vpn_connections" "vpn_conns" {
	output_file = "/tmp/vpnconns"
}

```

## Argument Reference

The following arguments are supported:

* `vpn_gateway_id` - (Optional) Use the VPN gateway ID as the search key.
* `customer_gateway_id` - (Optional)Use the VPN customer gateway ID as the search key.
* `name_regex` - (Optional) A regex string of VPN connection name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `vpn_connection_id` - ID of the VPN connection.
* `customer_gateway_id` - ID of the VPN customer gateway.
* `vpn_gateway_id` - ID of the VPN gateway.
* `name` - The name of the VPN connection.
* `local_subnet` - The local subnet of the VPN connection.
* `remote_subnet` - The remote subnet of the VPN connection.
* `ike_config` - The JSON string of the VPN connection IKE config
* `ipsec_config` - The JSON string of the VPN connection IPSec config
* `status` - The status of the VPN connection, valid value:ike_sa_not_established,ike_sa_established,ipsec_sa_not_established,ipsec_sa_established
