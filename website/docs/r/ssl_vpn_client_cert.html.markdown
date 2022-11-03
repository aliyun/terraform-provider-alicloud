---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_vpn_client_cert"
sidebar_current: "docs-alicloud-resource-ssl-vpn-client-cert"
description: |-
  Provides a Alicloud SSL VPN Client Cert resource.
---

# alicloud_ssl_vpn_client_cert

Provides a SSL VPN client cert resource.

-> **NOTE:** Terraform will auto build SSL VPN client certs while it uses `alicloud_ssl_vpn_client_cert` to build a ssl vpn client certs resource.
It depends on VPN instance and SSL VPN Server.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ssl_vpn_client_cert" "foo" {
  ssl_vpn_server_id = "ssl_vpn_server_fake_id"
  name              = "sslVpnClientCertExample"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) The name of the client certificate.
- `ssl_vpn_server_id` - (Required, ForceNew) The ID of the SSL-VPN server.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the SSL-VPN client certificate.
- `status` - The status of the client certificate.
- `ca_cert` - The client ca cert.
- `client_cert` - The client cert.
- `client_key` - The client key.
- `client_config` - The vpn client config.

## Import

SSL-VPN client certificates can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_vpn_client_cert.example vsc-abc123456
```
