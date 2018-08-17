---
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_vpn_client_certs"
sidebar_current: "docs-alicloud-datasource-ssl-vpn-client-certs"
description: |-
    Provides a list of SSL VPN client certs which owned by an Alicloud account.
---

# alicloud\_ssl_vpn_client_certs

The SSL VPN client certs data source lists a number of SSL VPN client certs resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_ssl_vpn_client_certs" "foo" {
	output_file = "/tmp/vpnClientCerts"
}

```

## Argument Reference

The following arguments are supported:

* `ssl_vpn_server_id` - (Optional) Use the SSL VPN server ID as the search key.
* `name_regex` - (Optional) A regex string of VPN name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `ssl_vpn_client_cert_id` - ID of the SSL VPN client certs instance.
* `name` - The name of the SSL VPN client certs instance.
* `status` - The status of the VPN client certs instance.
