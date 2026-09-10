---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_ipsec_server"
sidebar_current: "docs-alicloud-resource-vpn-ipsec-server"
description: |-
  Provides a Alicloud VPN Ipsec Server resource.
---

# alicloud_vpn_ipsec_server

Provides a VPN Ipsec Server resource.

For information about VPN Ipsec Server and how to use it, see [What is Ipsec Server](https://www.alibabacloud.com/help/en/vpn/sub-product-ssl-vpn/developer-reference/api-vpc-2016-04-28-createipsecserver-ssl-vpn).

-> **NOTE:** Available since v1.161.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpn_ipsec_server&exampleId=37ef967a-53f8-0fe9-8ca2-d517d282dff57818de5d&activeTab=example&spm=docs.r.vpn_ipsec_server.0.37ef967a53&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

resource "alicloud_vpn_ipsec_server" "foo" {
  client_ip_pool    = "10.0.0.0/24"
  ipsec_server_name = var.name
  local_subnet      = "192.168.0.0/24"
  vpn_gateway_id    = alicloud_vpn_gateway.default.id
  psk_enabled       = true
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpn_ipsec_server&spm=docs.r.vpn_ipsec_server.example&intl_lang=EN_US)

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
* `ike_config` - (Optional) The configuration of Phase 1 negotiations. See [`ike_config`](#ike_config) below.
* `ipsec_config` - (Optional) The configuration of Phase 2 negotiations. See [`ipsec_config`](#ipsec_config) below.

### `ike_config`

The ike_config supports the following:

* `ike_auth_alg` - (Optional) The authentication algorithm that is used in Phase 1 negotiations. Default value: `sha1`.
* `ike_enc_alg` - (Optional) The encryption algorithm that is used in Phase 1 negotiations. Default value: `aes`.
* `ike_lifetime` - (Optional) IkeLifetime: the SA lifetime determined by Phase 1 negotiations. Valid values: `0` to `86400`. Default value: `86400`. Unit: `seconds`.
* `ike_mode` - (Optional) The IKE negotiation mode. Default value: `main`.
* `ike_pfs` - (Optional) The Diffie-Hellman key exchange algorithm that is used in Phase 1 negotiations. Default value: `group2`.
* `ike_version` - (Optional) The IKE version. Valid values: `ikev1` and `ikev2`. Default value: `ikev2`.
* `local_id` - (Optional) The identifier of the IPsec server. The value can be a fully qualified domain name (FQDN) or an IP address. The default value is the public IP address of the VPN gateway.
* `remote_id` - (Optional) The identifier of the customer gateway. The value can be an FQDN or an IP address. By default, this parameter is not specified.

### `ipsec_config`

The ipsec_config supports the following:

* `ipsec_enc_alg` - (Optional) The encryption algorithm that is used in Phase 2 negotiations. Default value: `aes`.
* `ipsec_lifetime` - (Optional) The SA lifetime determined by Phase 2 negotiations. Valid values: `0` to `86400`. Default value: `86400`. Unit: `seconds`.
* `ipsec_pfs` - (Optional) Forwards packets of all protocols. The Diffie-Hellman key exchange algorithm used in Phase 2 negotiations. Default value: `group2`.
* `ipsec_auth_alg` - (Optional) The authentication algorithm that is used in Phase 2 negotiations. Default value: `sha1`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ipsec Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Ipsec Server.
* `update` - (Defaults to 1 mins) Used when update the Ipsec Server.
* `delete` - (Defaults to 1 mins) Used when delete the Ipsec Server.

## Import

VPN Ipsec Server can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_ipsec_server.example <id>
```