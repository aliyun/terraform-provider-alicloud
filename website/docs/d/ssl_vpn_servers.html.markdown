---
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_vpn_servers"
sidebar_current: "docs-alicloud-datasource-ssl-vpn-servers"
description: |-
    Provides a list of SSL VPN Servers which owned by an Alicloud account.
---

# alicloud\_ssl_vpn_servers

The SSL VPN servers data source lists a number of SSL VPN servers resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_ssl_vpn_servers" "foo" {
	output_file = "/tmp/sslVpnServers"
}

```

## Argument Reference

The following arguments are supported:

* `vpn_gateway_id` - (Optional) Use the VPN gateway ID as the search key.
* `ssl_vpn_server_id` - (Optional) Use the SSL VPN Server ID as the search key.
* `name_regex` - (Optional) A regex string of VPN name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `ssl_vpn_server_id` - ID of the SSL VPN server.
* `vpn_gateway_id` - ID of the VPN gateway.
* `name` - The name of the SSL VPN server.
* `client_ip_pool` - The client ip pool.
* `local_subnet` - The local subnet.
* `proto` - The protocol of the SSL VPN Server.
* `cipher` - The cipher of the SSL VPN Server.
* `compress` - Enable or Disable compress.
* `connections` - The current connection number of the SSL VPN server.
* `max_connections` - The max connections of the SSL VPN server.
* `internet_ip` - The internet_ip of the SSL VPN server.