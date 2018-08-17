---
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_vpn_client_cert"
sidebar_current: "docs-alicloud-resource-ssl-vpn-client-cert"
description: |-
  Provides a Alicloud SSL VPN Client Cert resource.
---

# alicloud\_ssl_vpn_client_cert

Provides a SSL VPN client cert resource.

~> **NOTE:** Terraform will auto build SSL VPN client certs  while it uses `alicloud_ssl_vpn_client_cert` to build a ssl vpn client certs resource.
             It depends on VPN instance and SSL VPN Server.
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
    name = "testAccSslVpnServerConfig"
    vpn_gateway_id = "${alicloud_vpn.foo.id}"
    client_ip_pool = "192.168.0.0/16"
    local_subnet = "172.16.0.0/21"
    proto = "UDP"
    cipher = "AES-128-CBC"
    port = "1194"
    compress = "false"
}

resource "alicloud_ssl_vpn_client_cert" "foo" {
    ssl_vpn_server_id = "${alicloud_ssl_vpn_server.foo.id}"
    name = "test_create_client_cert"
}
```
## Argument Reference

The following arguments are supported:
* `name` - (Optional) The name of the SSL VPN client certs. Defaults to null.
* `ssl_vpn_server_id` - (Required, Forces new resource) The SSL VPN server id.
* `description` - (Optional) The description of the VPN instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN instance id.





