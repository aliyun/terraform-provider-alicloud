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

[IPsec-VPN connections support the dual-tunnel mode](https://www.alibabacloud.com/help/en/vpn/product-overview/ipsec-vpn-connections-support-the-dual-tunnel-mode)

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpn_connection&exampleId=ada16e1a-2276-5be4-83e2-555f40750a75eb8fca1b&activeTab=example&spm=docs.r.vpn_connection.0.ada16e1a22&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

variable "spec" {
  default = "5"
}

data "alicloud_vpn_gateway_zones" "default" {
  spec = "5M"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default0" {
  cidr_block = "172.16.0.0/24"
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_vpn_gateway_zones.default.ids.0
}

resource "alicloud_vswitch" "default1" {
  cidr_block = "172.16.1.0/24"
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_vpn_gateway_zones.default.ids.1
}

resource "alicloud_vpn_gateway" "HA-VPN" {
  vpn_type                     = "Normal"
  disaster_recovery_vswitch_id = alicloud_vswitch.default1.id
  vpn_gateway_name             = var.name

  vswitch_id   = alicloud_vswitch.default0.id
  auto_pay     = true
  vpc_id       = alicloud_vpc.default.id
  network_type = "public"
  payment_type = "Subscription"
  enable_ipsec = true
  bandwidth    = var.spec
}

resource "alicloud_vpn_customer_gateway" "defaultCustomerGateway" {
  description           = "defaultCustomerGateway"
  ip_address            = "2.2.2.5"
  asn                   = "2224"
  customer_gateway_name = var.name
}

resource "alicloud_vpn_customer_gateway" "changeCustomerGateway" {
  description           = "changeCustomerGateway"
  ip_address            = "2.2.2.6"
  asn                   = "2225"
  customer_gateway_name = var.name
}

resource "alicloud_vpn_connection" "default" {
  vpn_gateway_id      = alicloud_vpn_gateway.HA-VPN.id
  vpn_connection_name = var.name
  local_subnet = [
    "3.0.0.0/24"
  ]
  remote_subnet = [
    "10.0.0.0/24",
    "10.0.1.0/24"
  ]
  tags = {
    Created = "TF"
    For     = "example"
  }
  enable_tunnels_bgp = "true"
  tunnel_options_specification {
    tunnel_ipsec_config {
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes256"
      ipsec_lifetime = "16400"
      ipsec_pfs      = "group5"
    }

    customer_gateway_id = alicloud_vpn_customer_gateway.defaultCustomerGateway.id
    role                = "master"
    tunnel_bgp_config {
      local_asn    = "1219002"
      tunnel_cidr  = "169.254.30.0/30"
      local_bgp_ip = "169.254.30.1"
    }

    tunnel_ike_config {
      ike_mode     = "aggressive"
      ike_version  = "ikev2"
      local_id     = "localid_tunnel2"
      psk          = "12345678"
      remote_id    = "remote2"
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes256"
      ike_lifetime = "3600"
      ike_pfs      = "group14"
    }

  }
  tunnel_options_specification {
    tunnel_ike_config {
      remote_id    = "remote24"
      ike_enc_alg  = "aes256"
      ike_lifetime = "27000"
      ike_mode     = "aggressive"
      ike_pfs      = "group5"
      ike_auth_alg = "md5"
      ike_version  = "ikev2"
      local_id     = "localid_tunnel2"
      psk          = "12345678"
    }

    tunnel_ipsec_config {
      ipsec_lifetime = "2700"
      ipsec_pfs      = "group14"
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes256"
    }

    customer_gateway_id = alicloud_vpn_customer_gateway.defaultCustomerGateway.id
    role                = "slave"
    tunnel_bgp_config {
      local_asn    = "1219002"
      local_bgp_ip = "169.254.40.1"
      tunnel_cidr  = "169.254.40.0/30"
    }
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
* `name` - (Deprecated since v1.216.0). Field 'name' has been deprecated from provider version 1.216.0. New field 'vpn_connection_name' instead.

### `bgp_config`

The bgp_config supports the following:
* `enable` - (Optional, Computed) specifies whether to enable BGP. Valid values: true and false (default).
* `local_asn` - (Optional, Computed) the autonomous system number (ASN) on the Alibaba Cloud side. 
   Valid values: 1 to 4294967295. Default value: 45104. You can enter a value in two segments separated by a period (.). 
   Each segment is 16 bits in length. Enter the number in each segment in decimal format. 
   For example, if you enter 123.456, the ASN is 8061384. The ASN is calculated by using the following formula: 123 × 65536 + 456 = 8061384.
* `local_bgp_ip` - (Optional, Computed) the BGP address on the Alibaba Cloud side. It must be an IP address that falls within the CIDR block of the IPsec tunnel.
* `tunnel_cidr` - (Optional, Computed) The CIDR block of the IPsec tunnel. The CIDR block must belong to 169.254.0.0/16 and the subnet mask is 30 bits in length.

### `health_check_config`

The health_check_config supports the following:
* `dip` - (Optional, Computed) the destination IP address configured for health checks.
* `enable` - (Optional, Computed) specifies whether to enable health checks. Valid values: true and false. Default value: false.
* `interval` - (Optional, Computed) the time interval of health check retries. Unit: seconds. Default value: 3.
* `retry` - (Optional, Computed) the maximum number of health check retries. Default value: 3.
* `sip` - (Optional, Computed) the source IP address that is used for health checks.

### `ike_config`

The ike_config supports the following:
* `ike_auth_alg` - (Optional, Computed) the authentication algorithm that is used in Phase 1 negotiations. Valid values: md5, sha1, sha2
* `ike_enc_alg` - (Optional, Computed) the encryption algorithm that is used in Phase 1 negotiations. Valid values: aes, aes192, aes256, des, and 3des. Default value: aes.
* `ike_lifetime` - (Optional, Computed) the SA lifetime as a result of Phase 1 negotiations. Unit: seconds. Valid values: 0 to 86400. Default value: 86400.
* `ike_local_id` - (Optional, Computed) the identifier of the VPN gateway. It can contain at most 100 characters. The default value is the IP address of the VPN gateway.
* `ike_mode` - (Optional, Computed) the negotiation mode of IKE. Valid values: main and aggressive. Default value: main.
  - main: This mode offers higher security during negotiations. 
  - aggressive: This mode supports faster negotiations and a higher success rate.
* `ike_pfs` - (Optional, Computed) the Diffie-Hellman key exchange algorithm that is used in Phase 1 negotiations. Valid values: group1, group2, group5, and group14. Default value: group2.
* `ike_remote_id` - (Optional, Computed) the identifier of the customer gateway. It can contain at most 100 characters. The default value is the IP address of the customer gateway.
* `ike_version` - (Optional, Computed) the version of the Internet Key Exchange (IKE) protocol. Valid values: ikev1 and ikev2. Default value: ikev1.
   Compared with IKEv1, IKEv2 simplifies the security association (SA) negotiation process and provides better support for scenarios with multiple CIDR blocks.
* `psk` - (Optional, Computed) the pre-shared key that is used for identity authentication between the VPN gateway and the on-premises data center. The key must be 1 to 100 characters in length and can contain digits, letters, and the following special characters: ~!\`@#$%^&*()_-+={}[]|;:',.<>/? If you do not specify a pre-shared key, the system randomly generates a 16-bit string as the pre-shared key. You can call the DescribeVpnConnection operation to query the pre-shared key that is automatically generated by the system.

### `ipsec_config`

The ipsec_config supports the following:
* `ipsec_auth_alg` - (Optional, Computed) the authentication algorithm that is used in Phase 2 negotiations. Valid values: md5, sha1, sha256, sha384, and sha512. Default value: md5.
* `ipsec_enc_alg` - (Optional, Computed) the encryption algorithm that is used in Phase 2 negotiations. Valid values: aes, aes192, aes256, des, and 3des. Default value: aes.
* `ipsec_lifetime` - (Optional, Computed) the SA lifetime that is determined by Phase 2 negotiations. Unit: seconds. Valid values: 0 to 86400. Default value: 86400.
* `ipsec_pfs` - (Optional, Computed) the DH key exchange algorithm that is used in Phase 2 negotiations. Valid values: disabled, group1, group2, group5, and group14. Default value: group2.

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
* `local_asn` - (Optional) The autonomous system number (ASN) of the tunnel on the Alibaba Cloud side. Valid values: 1 to 4294967295. Default value: 45104.
* `local_bgp_ip` - (Optional) The BGP IP address of the tunnel on the Alibaba Cloud side. The address is an IP address that falls within the BGP CIDR block.
* `tunnel_cidr` - (Optional) The BGP CIDR block of the tunnel. The CIDR block must fall within the 169.254.0.0/16 range. The subnet mask of the CIDR block must be 30 bits in length.

### `tunnel_options_specification-tunnel_ike_config`

The tunnel_options_specification-tunnel_ike_config supports the following:
* `ike_auth_alg` - (Optional) The authentication algorithm that is used in Phase 1 negotiations. Valid values: md5, sha1, sha256, sha384, and sha512. Default value: md5.
* `ike_enc_alg` - (Optional) The encryption algorithm that is used in Phase 1 negotiations. Valid values: aes, aes192, aes256, des, and 3des. Default value: aes.
* `ike_lifetime` - (Optional) The SA lifetime as a result of Phase 1 negotiations. Unit: seconds. Valid values: 0 to 86400. Default value: 86400.
* `ike_mode` - (Optional) The negotiation mode of IKE. Valid values: main and aggressive. Default value: main.
  - main: This mode offers higher security during negotiations. 
  - aggressive: This mode supports faster negotiations and a higher success rate.
* `ike_pfs` - (Optional) The Diffie-Hellman key exchange algorithm that is used in Phase 1 negotiations. Default value: group2.
* `ike_version` - (Optional) The version of the IKE protocol. Valid values: ikev1 and ikev2. Default value: ikev1.
   Compared with IKEv1, IKEv2 simplifies the SA negotiation process and provides better support for scenarios with multiple CIDR blocks.
* `local_id` - (Optional) The identifier of the tunnel on the Alibaba Cloud side, which is used in Phase 1 negotiations. It can contain at most 100 characters. The default value is the IP address of the tunnel.
   LocalId supports fully qualified domain names (FQDNs). If you use an FQDN, we recommend that you set the negotiation mode to aggressive.
* `psk` - (Optional) The pre-shared key that is used for identity authentication between the tunnel and the tunnel peer.
   The key must be 1 to 100 characters in length and can contain digits, letters, and the following special characters: ~!\`@#$%^&*()_-+={}[]|;:',.<>/?
   If you do not specify a pre-shared key, the system randomly generates a 16-bit string as the pre-shared key. You can call the DescribeVpnConnection operation to query the pre-shared key that is automatically generated by the system.
* `remote_id` - (Optional) The identifier of the tunnel peer, which is used in Phase 1 negotiations. It can contain at most 100 characters. The default value is the IP address of the customer gateway that is associated with the tunnel.
   RemoteId supports FQDNs. If you use an FQDN, we recommend that you set the negotiation mode to aggressive.

### `tunnel_options_specification-tunnel_ipsec_config`

The tunnel_options_specification-tunnel_ipsec_config supports the following:
* `ipsec_auth_alg` - (Optional) The authentication algorithm that is used in Phase 2 negotiations. Valid values: md5, sha1, sha256, sha384, and sha512. Default value: md5.
* `ipsec_enc_alg` - (Optional) The encryption algorithm that is used in Phase 2 negotiations. Valid values: aes, aes192, aes256, des, and 3des. Default value: aes.
* `ipsec_lifetime` - (Optional) The SA lifetime as a result of Phase 2 negotiations. Unit: seconds. Valid values: 0 to 86400. Default value: 86400.
* `ipsec_pfs` - (Optional) The Diffie-Hellman key exchange algorithm that is used in Phase 2 negotiations. Default value: group2. Valid values: disabled, group1, group2, group5, and group14.

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