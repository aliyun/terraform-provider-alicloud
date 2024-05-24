---
subcategory: "VPN Gateway"
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

-> **NOTE:** Available since v1.15.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "172.16.0.0/16"
}

data "alicloud_vswitches" "default0" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}

data "alicloud_vswitches" "default1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.1
}

resource "alicloud_vpn_gateway" "default" {
  vpn_gateway_name             = var.name
  vpc_id                       = data.alicloud_vpcs.default.ids.0
  bandwidth                    = "10"
  enable_ssl                   = true
  description                  = var.name
  payment_type                 = "Subscription"
  vswitch_id                   = data.alicloud_vswitches.default0.ids.0
  disaster_recovery_vswitch_id = data.alicloud_vswitches.default1.ids.0
}


resource "alicloud_ssl_vpn_server" "default" {
  name           = var.name
  vpn_gateway_id = alicloud_vpn_gateway.default.id
  client_ip_pool = "192.168.0.0/16"
  local_subnet   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 8)
  protocol       = "UDP"
  cipher         = "AES-128-CBC"
  port           = "1194"
  compress       = "false"
}

resource "alicloud_ssl_vpn_client_cert" "default" {
  ssl_vpn_server_id = alicloud_ssl_vpn_server.default.id
  name              = var.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the client certificate.
* `ssl_vpn_server_id` - (Required, ForceNew) The ID of the SSL-VPN server.
* `status` - (Optional) The status of the client certificate.
* `ca_cert` - (Optional) The client ca cert.
* `client_cert` - (Optional) The client cert.
* `client_key` - (Optional) The client key.
* `client_config` - (Optional) The vpn client config.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SSL-VPN client certificate.

## Import

SSL-VPN client certificates can be imported using the id, e.g.

```shell
$ terraform import alicloud_ssl_vpn_client_cert.example vsc-abc123456
```
