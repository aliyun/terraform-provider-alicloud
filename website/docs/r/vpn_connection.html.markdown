---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_connection"
description: |-
  Provides a Alicloud VPN connection resource.
---

# alicloud_vpn_connection

Provides a VPN connection resource.

-> **NOTE:** Terraform will auto build vpn connection while it uses `alicloud_vpn_connection` to build a vpn connection resource.
             The vpn connection depends on VPN and VPN customer gateway.

For information about VPN connection and how to use it, see [What is vpn connection](https://www.alibabacloud.com/help/en/doc-detail/120390.html).

-> **NOTE:** Available since v1.14.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "me-east-1"
}

variable "spec" {
  default = "20"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "me-east-1a"
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = "me-east-1a"
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "default" {
  vpn_type         = "Normal"
  vpn_gateway_name = var.name

  vswitch_id   = local.vswitch_id
  auto_pay     = true
  vpc_id       = data.alicloud_vpcs.default.ids.0
  network_type = "public"
  payment_type = "Subscription"
  enable_ipsec = true
  bandwidth    = var.spec
}

resource "alicloud_vpn_customer_gateway" "default" {
  description           = var.name
  ip_address            = "4.3.2.10"
  asn                   = "1219002"
  customer_gateway_name = var.name
}

resource "alicloud_vpn_connection" "default" {
  local_subnet = [
    "3.0.0.0/24"
  ]
  enable_nat_traversal = "true"
  bgp_config {
    local_bgp_ip = "169.254.10.1"
    tunnel_cidr  = "169.254.10.0/30"
    enable       = "true"
    local_asn    = "1219002"
  }

  customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  vpn_gateway_id      = alicloud_vpn_gateway.default.id
  vpn_connection_name = var.name
  effect_immediately  = "true"
  health_check_config {
    enable   = "true"
    dip      = "1.1.1.1"
    retry    = "3"
    sip      = "3.3.3.3"
    interval = "3"
  }

  remote_subnet = [
    "10.0.0.0/24",
    "10.0.1.0/24"
  ]
  ipsec_config {
    ipsec_enc_alg  = "aes"
    ipsec_auth_alg = "sha1"
    ipsec_lifetime = "86400"
    ipsec_pfs      = "group2"
  }

  auto_config_route = "true"
  enable_dpd        = "true"
  ike_config {
    ike_lifetime  = "86400"
    ike_local_id  = "localid1"
    ike_version   = "ikev2"
    ike_mode      = "main"
    psk           = "12345678"
    ike_remote_id = "remoteId2"
    ike_pfs       = "group2"
    ike_auth_alg  = "sha1"
    ike_enc_alg   = "aes"
  }
}
```

## Argument Reference

The following arguments are supported:

* `local_subnet` - (Required) The CIDR block of the VPC to be connected with the local data center. This parameter is used for phase-two negotiation.
* `remote_subnet` - (Required) The CIDR block of the local data center. This parameter is used for phase-two negotiation.
* `auto_config_route` - (Optional) Whether to configure routing automatically. Value:
  - **true**: Automatically configure routes.
  - **false**: does not automatically configure routes.
* `bgp_config` - (Optional, Computed) vpnBgp configuration. See [`bgp_config`](#bgp_config) below.
* `customer_gateway_id` - (Optional, ForceNew) The ID of the customer gateway.
* `effect_immediately` - (Optional) Indicates whether IPsec-VPN negotiations are initiated immediately. Valid values.
* `enable_dpd` - (Optional, Computed) Wether enable Dpd detection.
* `enable_nat_traversal` - (Optional, Computed) enable nat traversal.
* `enable_tunnels_bgp` - (Optional, Computed) Enable tunnel bgp.
* `health_check_config` - (Optional, Computed) Health Check information. See [`health_check_config`](#health_check_config) below.
* `ike_config` - (Optional, Computed) The configuration of Phase 1 negotiations. See [`ike_config`](#ike_config) below.
* `ipsec_config` - (Optional, Computed) IPsec configuration. See [`ipsec_config`](#ipsec_config) below.
* `network_type` - (Optional) The network type of the IPsec connection. Value:
  - **public**: public network, indicating that the IPsec connection establishes an encrypted communication channel through the public network.
  - **private**: private network, indicating that the IPsec connection establishes an encrypted communication channel through the private network.
* `tags` - (Optional, Map) Tags.
* `tunnel_options_specification` - (Optional) The tunnel options of IPsec. See [`tunnel_options_specification`](#tunnel_options_specification) below.
* `vpn_connection_name` - (Optional) The name of the IPsec-VPN connection.
* `vpn_gateway_id` - (Required, ForceNew) The ID of the VPN gateway.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.216.0). Field 'name' has been deprecated from provider version 1.216.0. New field 'vpn_connection_name' instead.

### `bgp_config`

The bgp_config supports the following:
* `enable` - (Optional, Computed) Bgp enable.
* `local_asn` - (Optional, Computed) Local asn.
* `local_bgp_ip` - (Optional, Computed) Local bgp IP.
* `tunnel_cidr` - (Optional, Computed) IPSec tunnel Cidr.

### `health_check_config`

The health_check_config supports the following:
* `dip` - (Optional, Computed) Destination IP.
* `enable` - (Optional, Computed) Specifies whether to enable healthcheck.
* `interval` - (Optional, Computed) Retry interval.
* `retry` - (Optional, Computed) retry times.
* `sip` - (Optional, Computed) Source IP.

### `ike_config`

The ike_config supports the following:
* `ike_auth_alg` - (Optional, Computed) IKE auth Algorithm.
* `ike_enc_alg` - (Optional, Computed) IKE encript algorithm.
* `ike_lifetime` - (Optional, Computed) IKE lifetime.
* `ike_local_id` - (Optional, Computed) The local ID, which supports the FQDN and IP formats, and defaults to the IP address of the selected VPN gateway.
* `ike_mode` - (Optional, Computed) IKE mode, supports main and aggressive mode. The main mode is highly secure. If NAT traversal is enabled, we recommend that you use the aggressive mode.
* `ike_pfs` - (Optional, Computed) DH group.
* `ike_remote_id` - (Optional, Computed) The peer ID. The FQDN and IP address formats are supported. The default value is the IP address of the selected customer gateway.
* `ike_version` - (Optional, Computed) IKE version.
* `psk` - (Optional, Computed) Preshared secret key.

### `ipsec_config`

The ipsec_config supports the following:
* `ipsec_auth_alg` - (Optional, Computed) IPsec authentication algorithm. sha1 and md5 are supported.
* `ipsec_enc_alg` - (Optional, Computed) IPsec Encript algorithm.
* `ipsec_lifetime` - (Optional, Computed) IPsec lifetime.
* `ipsec_pfs` - (Optional, Computed) DH Group.

### `tunnel_options_specification`

The tunnel_options_specification supports the following:
* `customer_gateway_id` - (Optional, ForceNew) The ID of the customer gateway in Tunnel.
* `enable_dpd` - (Optional) Wether enable Dpd detection.
* `enable_nat_traversal` - (Optional) enable nat traversal.
* `role` - (Optional) The role of Tunnel.
* `tunnel_bgp_config` - (Optional) The bgp config of Tunnel. See [`tunnel_bgp_config`](#tunnel_options_specification-tunnel_bgp_config) below.
* `tunnel_ike_config` - (Optional) The configuration of Phase 1 negotiations in Tunnel. See [`tunnel_ike_config`](#tunnel_options_specification-tunnel_ike_config) below.
* `tunnel_ipsec_config` - (Optional) IPsec configuration in Tunnel. See [`tunnel_ipsec_config`](#tunnel_options_specification-tunnel_ipsec_config) below.

### `tunnel_options_specification-tunnel_bgp_config`

The tunnel_options_specification-tunnel_bgp_config supports the following:
* `local_asn` - (Optional) Local asn.
* `local_bgp_ip` - (Optional) Local bgp IP.
* `tunnel_cidr` - (Optional) BGP Tunnel CIDR.

### `tunnel_options_specification-tunnel_ike_config`

The tunnel_options_specification-tunnel_ike_config supports the following:
* `ike_auth_alg` - (Optional) IKE auth Algorithm.
* `ike_enc_alg` - (Optional) IKE encript algorithm.
* `ike_lifetime` - (Optional) IKE lifetime.
* `ike_mode` - (Optional) IKE Mode.
* `ike_pfs` - (Optional) DH Group.
* `ike_version` - (Optional) IKE Version.
* `local_id` - (Optional) The local Id.
* `psk` - (Optional) Preshared secret key.
* `remote_id` - (Optional) Remote ID.

### `tunnel_options_specification-tunnel_ipsec_config`

The tunnel_options_specification-tunnel_ipsec_config supports the following:
* `ipsec_auth_alg` - (Optional) IPsec Auth algorithm.
* `ipsec_enc_alg` - (Optional) IPsec Encript algorithm.
* `ipsec_lifetime` - (Optional) IPsec  lifetime.
* `ipsec_pfs` - (Optional) DH Group.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `bgp_config` - vpnBgp configuration.
  * `status` - Whether BGP function is turned on.
* `create_time` - The time when the IPsec-VPN connection was created.
* `resource_group_id` - The ID of the resource group.
* `status` - The status of the resource.
* `tunnel_options_specification` - The tunnel options of IPsec.
  * `internet_ip` - The local internet IP in Tunnel.
  * `state` - The state of Tunnel.
  * `status` - The negotiation status of Tunnel.
  * `tunnel_id` - The tunnel ID of IPsec-VPN connection.
  * `zone_no` - The zoneNo of tunnel.
  * `tunnel_bgp_config` - The config of bgp.
    * `peer_asn` - Peer asn.
    * `peer_bgp_ip` - Peer bgp ip.
    * `bgp_status` - Whether BGP function is turned on.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpn Connection.
* `delete` - (Defaults to 5 mins) Used when delete the Vpn Connection.
* `update` - (Defaults to 5 mins) Used when update the Vpn Connection.

## Import

VPN connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_connection.example <id>
```