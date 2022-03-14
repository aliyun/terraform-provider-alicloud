---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_ipsec_server"
sidebar_current: "docs-alicloud-resource-vpn-ipsec-server"
description: |-
  Provides a Alicloud VPN Ipsec Server resource.
---

# alicloud\_vpn\_ipsec\_server

Provides a VPN Ipsec Server resource.

For information about VPN Ipsec Server and how to use it, see [What is Ipsec Server](https://www.alibabacloud.com/help/en/doc-detail/205454.html).

-> **NOTE:** Available in v1.161.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
locals {
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}

resource "alicloud_vpn_gateway" "default" {
  name                 = var.name
  vpc_id               = data.alicloud_vpcs.default.ids.0
  bandwidth            = 10
  enable_ssl           = true
  enable_ipsec         = true
  ssl_connections      = 5
  instance_charge_type = "PrePaid"
  vswitch_id           = local.vswitch_id
}

resource "alicloud_vpn_ipsec_server" "example" {
  client_ip_pool    = "example_value"
  ipsec_server_name = "example_value"
  local_subnet      = "example_value"
  vpn_gateway_id    = alicloud_vpn_gateway.default.id
}
```

## Argument Reference

The following arguments are supported:

* `client_ip_pool` - (Required) The client CIDR block. It refers to the CIDR block that is allocated to the virtual interface of the client.
* `dry_run` - (Optional) The dry run.
* `effect_immediately` - (Optional) Specifies whether you want the configuration to immediately take effect.
* `ipsec_server_name` - (Optional) The name of the IPsec server. The name must be `2` to `128` characters in length, and can contain digits, hyphens (-), and underscores (_). It must start with a letter.
* `local_subnet` - (Required) The local CIDR block. It refers to the CIDR block of the virtual private cloud (VPC) that is used to connect with the client. Separate multiple CIDR blocks with commas (,). Example: `192.168.1.0/24,192.168.2.0/24`.
* `psk` - (Optional) The pre-shared key. The pre-shared key is used to authenticate the VPN gateway and the client. By default, the system generates a random string that is 16 bits in length. You can also specify the pre-shared key. It can contain at most 100 characters.
* `psk_enabled` - (Optional) Whether to enable the pre-shared key authentication method. The value is only `true`, which indicates that the pre-shared key authentication method is enabled.
* `vpn_gateway_id` - (Required, ForceNew) The ID of the VPN gateway.
* `ike_config` - (Optional) The configuration of Phase 1 negotiations. See the following `Block ike_config`.
* `ipsec_config` - (Optional) The configuration of Phase 2 negotiations. See the following `Block ipsec_config`.

#### Block ike_config

The ike_config supports the following:

* `ike_auth_alg` - (Optional) The authentication algorithm that is used in Phase 1 negotiations. Default value: `sha1`.
* `ike_enc_alg` - (Optional) The encryption algorithm that is used in Phase 1 negotiations. Default value: `aes`.
* `ike_lifetime` - (Optional) IkeLifetime: the SA lifetime determined by Phase 1 negotiations. Valid values: `0` to `86400`. Default value: `86400`. Unit: `seconds`.
* `ike_mode` - (Optional) The IKE negotiation mode. Default value: `main`.
* `ike_pfs` - (Optional) The Diffie-Hellman key exchange algorithm that is used in Phase 1 negotiations. Default value: `group2`.
* `ike_version` - (Optional) The IKE version. Valid values: `ikev1` and `ikev2`. Default value: `ikev2`.
* `local_id` - (Optional) The identifier of the IPsec server. The value can be a fully qualified domain name (FQDN) or an IP address. The default value is the public IP address of the VPN gateway.
* `remote_id` - (Optional) The identifier of the customer gateway. The value can be an FQDN or an IP address. By default, this parameter is not specified.

#### Block ipsec_config

The ipsec_config supports the following:

* `ipsec_enc_alg` - (Optional) The encryption algorithm that is used in Phase 2 negotiations. Default value: `aes`.
* `ipsec_lifetime` - (Optional) The SA lifetime determined by Phase 2 negotiations. Valid values: `0` to `86400`. Default value: `86400`. Unit: `seconds`.
* `ipsec_pfs` - (Optional) Forwards packets of all protocols. The Diffie-Hellman key exchange algorithm used in Phase 2 negotiations. Default value: `group2`.
* `ipsec_auth_alg` - (Optional) The authentication algorithm that is used in Phase 2 negotiations. Default value: `sha1`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ipsec Server.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Ipsec Server.
* `update` - (Defaults to 1 mins) Used when update the Ipsec Server.
* `delete` - (Defaults to 1 mins) Used when delete the Ipsec Server.

## Import

VPN Ipsec Server can be imported using the id, e.g.

```
$ terraform import alicloud_vpn_ipsec_server.example <id>
```