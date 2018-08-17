---
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_vpn_server"
sidebar_current: "docs-alicloud-resource-ssl-vpn-server"
description: |-
  Provides a Alicloud SSL VPN server resource.
---

# alicloud\_ssl_vpn_server

Provides a SSL VPN server resource.

~> **NOTE:** Terraform will auto build ssl vpn server while it uses `alicloud_ssl_vpn_server` to build a ssl vpn server resource.

## Example Usage

Basic Usage

```
resource "alicloud_vpn" "foo" {
        name = "testAccVpnConfig_create"
        vpc_id = "vpc-2ze9wy916mfwpwbf6hx4u"
		bandwidth = "10"
        enable_ssl = true
        instance_charge_type = "postpaid"
        auto_pay = true
		description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
    name = "testAccSslVpnServerConfig_create"
    vpn_gateway_id = "${alicloud_vpn.foo.id}"
    client_ip_pool = "192.168.10.0/24"
    local_subnet = "172.16.0.0/24"
    proto = "UDP"
    cipher = "AES-192-CBC"
    port = "1194"
    compress = "false"
}
```
## Argument Reference

The following arguments are supported:
* `name` - (Optional) The name of the SSL VPN server. Defaults to null.
* `vpn_gateway_id` - (Required, ForceNew) The VPN gateway ID.
* `client_ip_pool` - (Required) The remote client virtual interface will get ip from the pool.The pool can't conflict with VPC CIDR on cloud.
* `local_subnet` - (Required) The VPC CIDR on the cloud.
* `proto` - (Optional) The protocol type of the SSL VPN server.
* `cipher` - (Optional) The valid value is AES-128-CBC, AES-192-CBC, AES-256-CBC.
* `port` - (Optional) The default value is 1194.
* `compress`  - (Optional) The default value is false.
* `description` - (Optional) The description of the VPN instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SSL VPN instance id.
* `internet_ip` - The internet ip of the SSL VPN.




