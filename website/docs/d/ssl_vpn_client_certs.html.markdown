---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_vpn_client_certs"
sidebar_current: "docs-alicloud-datasource-ssl-vpn-client-certs"
description: |-
    Provides a list of SSL-VPN client certificates which owned by an Alicloud account.
---

# alicloud\_ssl_vpn_client_certs

The SSL-VPN client certificates data source lists lots of SSL-VPN client certificates resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_ssl_vpn_client_certs" "foo" {
  ids               = ["fake-cert-id"]
  ssl_vpn_server_id = "fake-server-id"
  output_file       = "/tmp/clientcert"
  name_regex        = "^foo"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) IDs of the SSL-VPN client certificates.
* `ssl_vpn_server_id` - (Optional) Use the SSL-VPN server ID as the search key.
* `name_regex` - (Optional) A regex string of SSL-VPN client certificate name.
* `output_file` - (Optional) Save the result to the file.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of SSL-VPN client cert IDs.
* `names` - A list of SSL-VPN client cert names.
* `certs` - A list of SSL-VPN client certificates. Each element contains the following attributes:
  * `id` - ID of the SSL-VPN client certificate.
  * `ssl_vpn_server_id` - ID of the SSL-VPN Server.
  * `name` - The name of the SSL-VPN client certificate.
  * `create_time` - The time of creation.
  * `end_time` - The expiration time of the client certificate.
  * `status` - The status of the client certificate. valid value:expiring-soon, normal, expired.
