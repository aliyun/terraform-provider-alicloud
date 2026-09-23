---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_vpn_servers"
sidebar_current: "docs-alicloud-datasource-ssl-vpn-servers"
description: |-
    Provides a list of SSL-VPN servers which owned by an Alicloud account.
---

# alicloud\_ssl_vpn_servers

The SSL-VPN servers data source lists lots of SSL-VPN servers resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_ssl_vpn_servers" "foo" {
  ids            = ["fake-server-id"]
  vpn_gateway_id = "fake-vpn-id"
  output_file    = "/tmp/sslserver"
  name_regex     = "^foo"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) IDs of the SSL-VPN servers.
* `vpn_gateway_id` - (Optional) Use the VPN gateway ID as the search key.
* `name_regex` - (Optional) A regex string of SSL-VPN server name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of SSL-VPN server IDs.
* `names` - A list of SSL-VPN server names.
* `servers` - A list of SSL-VPN servers. Each element contains the following attributes:
  * `vpn_gateway_id` - The ID of the VPN gateway instance.
  * `id` - The ID of the SSL-VPN server.
  * `name` - The name of the SSL-VPN server.
  * `create_time` - The time of creation.
  * `compress` - Whether to compress.
  * `cipher` - The encryption algorithm used.
  * `proto` - The protocol used by the SSL-VPN server.
  * `port` - The port used by the SSL-VPN server.
  * `client_ip_pool` - The IP address pool of the client.
  * `local_subnet` - The local subnet of the VPN connection.
  * `internet_ip` - The public IP.
  * `connections` - The number of current connections.
  * `max_connections` - The maximum number of connections.
