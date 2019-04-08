---
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_vpn_server"
sidebar_current: "docs-alicloud-resource-ssl-vpn-server"
description: |-
  Provides a Alicloud SSL VPN server resource.
---

# alicloud\_ssl_vpn_server

Provides a SSL VPN server resource. [Refer to details](https://www.alibabacloud.com/help/doc-detail/64960.htm)

-> **NOTE:** Terraform will auto build ssl vpn server while it uses `alicloud_ssl_vpn_server` to build a ssl vpn server resource.

## Example Usage

Basic Usage

```
resource "alicloud_vpn_gateway" "foo" {
	name = "testAccVpnConfig_create"
	vpc_id = "vpc-fake-id"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name = "testAccSslVpnServerConfig_create"
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
	client_ip_pool = "192.168.0.0/16"
	local_subnet = "172.16.0.0/21"
	protocol = "UDP"
	cipher = "AES-128-CBC"
	port = 1194
	compress = "false"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the SSL-VPN server.
* `vpn_gateway_id` - (Required, ForceNew) The ID of the VPN gateway.
* `client_ip_pool` - (Required) The CIDR block from which access addresses are allocated to the virtual network interface card of the client.
* `local_subnet` - (Required) The CIDR block to be accessed by the client through the SSL-VPN connection.
* `protocol` - (Optional) The protocol used by the SSL-VPN server. Valid value: UDP(default) |TCP
* `cipher` - (Optional) The encryption algorithm used by the SSL-VPN server. Valid value: AES-128-CBC (default)| AES-192-CBC | AES-256-CBC | none
* `port` - (Optional) The port used by the SSL-VPN server. The default value is 1194.The following ports cannot be used: [22, 2222, 22222, 9000, 9001, 9002, 7505, 80, 443, 53, 68, 123, 4510, 4560, 500, 4500].
* `compress`  - (Optional) Specify whether to compress the communication. Valid value: true (default) | false

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SSL-VPN server.
* `internet_ip` - The internet IP of the SSL-VPN server.
* `connections` - The number of current connections.
* `max_connections` - The maximum number of connections.

## Import

SSL-VPN server can be imported using the id, e.g.

```
$ terraform import alicloud_ssl_vpn_server.example vss-abc123456
```


