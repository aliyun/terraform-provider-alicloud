---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_vpn_attachment"
description: |-
  Provides a Alicloud V P N Gateway Vpn Attachment resource.
---

# alicloud_vpn_gateway_vpn_attachment

Provides a V P N Gateway Vpn Attachment resource.

VpnAttachment has been upgraded to dual-tunnel mode. When you create a VpnAttachment in dual tunnel mode, you can configure the following request parameters in addition to the required parameters: vpn_attachment_name, network_type, effectImmediately, tags array, resource_group_id, tunnel_options_specification array, and enable_tunnels_bgp.

For information about V P N Gateway Vpn Attachment and how to use it, see [What is Vpn Attachment](https://www.alibabacloud.com/help/zh/virtual-private-cloud/latest/createvpnattachment).

-> **NOTE:** Available since v1.181.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_vpn_customer_gateway" "default" {
  customer_gateway_name = var.name
  ip_address            = "42.104.22.210"
  asn                   = "45014"
  description           = var.name
}
resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  customer_gateway_id = alicloud_vpn_customer_gateway.default.id
  network_type        = "public"
  local_subnet        = "0.0.0.0/0"
  remote_subnet       = "0.0.0.0/0"
  effect_immediately  = false
  ike_config {
    ike_auth_alg = "md5"
    ike_enc_alg  = "des"
    ike_version  = "ikev2"
    ike_mode     = "main"
    ike_lifetime = 86400
    psk          = "tf-testvpn2"
    ike_pfs      = "group1"
    remote_id    = "testbob2"
    local_id     = "testalice2"
  }
  ipsec_config {
    ipsec_pfs      = "group5"
    ipsec_enc_alg  = "des"
    ipsec_auth_alg = "md5"
    ipsec_lifetime = 86400
  }
  bgp_config {
    enable       = true
    local_asn    = 45014
    tunnel_cidr  = "169.254.11.0/30"
    local_bgp_ip = "169.254.11.1"
  }
  health_check_config {
    enable   = true
    sip      = "192.168.1.1"
    dip      = "10.0.0.1"
    interval = 10
    retry    = 10
    policy   = "revoke_route"

  }
  enable_dpd           = true
  enable_nat_traversal = true
  vpn_attachment_name  = var.name
}
```

## Argument Reference

The following arguments are supported:
* `bgp_config` - (Optional, Computed, List) Bgp configuration information.
  - This parameter is supported when you create an vpn attachment in single-tunnel mode. See [`bgp_config`](#bgp_config) below.
* `customer_gateway_id` - (Optional, Available since v1.196.0) Customer gateway ID.
  - This parameter is required when creating a single-tunnel mode vpn attachment.
* `effect_immediately` - (Optional, Computed) Specifies whether to immediately start IPsec negotiations after the configuration takes effect. Valid values:
  - `true`: immediately starts IPsec negotiations after the configuration is complete.
  - `false` (default): starts IPsec negotiations when inbound traffic is received.
* `enable_dpd` - (Optional, Computed) This parameter is supported if you create an vpn attachment in single-tunnel mode.
Whether to enable the DPD (peer survival detection) function.
  - true (default): enables DPD. The initiator of the IPsec-VPN connection sends DPD packets to check the existence and availability of the peer. If no feedback is received from the peer within the specified period of time, the connection fails. In this case, ISAKMP SA and IPsec SA are deleted along with the security tunnel.
  - false: disables DPD. The initiator of the IPsec-VPN connection does not send DPD packets.
* `enable_nat_traversal` - (Optional, Computed) This parameter is supported if you create an vpn attachment in single-tunnel mode.
Specifies whether to enable NAT traversal. Valid values:
  - true (default): enables NAT traversal. After NAT traversal is enabled, the initiator does not check the UDP ports during IKE negotiations and can automatically discover NAT gateway devices along the vpn attachment tunnel.
  - false: disables NAT traversal.
* `enable_tunnels_bgp` - (Optional, Computed, Available since v1.246.0) You can configure this parameter when you create a vpn attachment in dual-tunnel mode.Whether to enable the BGP function for the tunnel. Value: `true` or `false` (default).

-> **NOTE:**  before adding BGP configuration, we recommend that you understand the working mechanism and usage restrictions of the BGP dynamic routing function. 

* `health_check_config` - (Optional, Computed, List) This parameter is supported if you create an vpn attachment in single-tunnel mode.
Health check configuration information. See [`health_check_config`](#health_check_config) below.
* `ike_config` - (Optional, Computed, List) The configurations of Phase 1 negotiations. 
  - This parameter is supported if you create an vpn attachment in single-tunnel mode. See [`ike_config`](#ike_config) below.
* `ipsec_config` - (Optional, Computed, List) Configuration negotiated in the second stage. 
  - This parameter is supported if you create an vpn attachment in single-tunnel mode. See [`ipsec_config`](#ipsec_config) below.
* `local_subnet` - (Required) The CIDR block on the VPC side. The CIDR block is used in Phase 2 negotiations.Separate multiple CIDR blocks with commas (,). Example: 192.168.1.0/24,192.168.2.0/24.The following routing modes are supported:
  - If you set LocalSubnet and RemoteSubnet to 0.0.0.0/0, the routing mode of the IPsec-VPN connection is set to Destination Routing Mode.
  - If you set LocalSubnet and RemoteSubnet to specific CIDR blocks, the routing mode of the IPsec-VPN connection is set to Protected Data Flows.
* `network_type` - (Optional, ForceNew, Computed) network type
  - `public` (default)
  - `private`
* `remote_subnet` - (Required) The CIDR block on the data center side. This CIDR block is used in Phase 2 negotiations.Separate multiple CIDR blocks with commas (,). Example: 192.168.3.0/24,192.168.4.0/24.The following routing modes are supported:
  - If you set LocalSubnet and RemoteSubnet to 0.0.0.0/0, the routing mode of the IPsec-VPN connection is set to Destination Routing Mode.
  - If you set LocalSubnet and RemoteSubnet to specific CIDR blocks, the routing mode of the IPsec-VPN connection is set to Protected Data Flows.
* `resource_group_id` - (Optional, Computed, Available since v1.246.0) The ID of the resource group
* `tags` - (Optional, Map, Available since v1.246.0) Tags
* `tunnel_options_specification` - (Optional, Computed, Set, Available since v1.246.0) Configure the tunnel.
  - You can configure parameters in the `tunnel_options_specification` array when you create a vpn attachment in dual-tunnel mode.
  - When creating a vpn attachment in dual-tunnel mode, you must add both tunnels for the vpn attachment to ensure that the vpn attachment has link redundancy. Only two tunnels can be added to a vpn attachment. See [`tunnel_options_specification`](#tunnel_options_specification) below.
* `vpn_attachment_name` - (Optional) vpn attachment name

### `bgp_config`

The bgp_config supports the following:
* `enable` - (Optional, Computed) Whether to enable the BGP function. Valid values: true or false (default).
* `local_asn` - (Optional, Computed, Int) The autonomous system number on the Alibaba Cloud side. The value range of autonomous system number is 1~4294967295. Default value: 45104
* `local_bgp_ip` - (Optional, Computed) The BGP address on the Alibaba Cloud side. This address is an IP address in the IPsec tunnel network segment.
  - Before adding the BGP configuration, we recommend that you understand the working mechanism and usage restrictions of the BGP dynamic routing function. For more information, see BGP Dynamic Routing Bulletin.
  - We recommend that you use the private number of the autonomous system number to establish a BGP connection with Alibaba Cloud. Please refer to the documentation for the private number range of the autonomous system number.
* `tunnel_cidr` - (Optional, Computed) IPsec tunnel network segment. This network segment must be a network segment with a mask length of 30 within 169.254.0.0/16

### `health_check_config`

The health_check_config supports the following:
* `dip` - (Optional, Computed) Target IP.
* `enable` - (Optional, Computed) Whether health check is enabled:-`false`: not enabled. - `true`: enabled.
* `interval` - (Optional, Computed, Int) The health check retry interval, in seconds.
* `policy` - (Optional, Computed) Whether to revoke the published route when the health check fails
  - `revoke_route`(default): withdraws published routes.
  - `reserve_route`: does not withdraw published routes.
* `retry` - (Optional, Computed, Int) Number of retries for health check.
* `sip` - (Optional, Computed) SOURCE IP.

### `ike_config`

The ike_config supports the following:
* `ike_auth_alg` - (Optional, Computed) The authentication algorithm negotiated in the first stage. Valid values: md5, sha1, sha256, sha384, sha512. Default value: md5.
* `ike_enc_alg` - (Optional, Computed) The encryption algorithm that is used in Phase 1 negotiations. Valid values: aes, aes192, aes256, des, and 3des. Default value: aes.
* `ike_lifetime` - (Optional, Computed, Int) The SA lifetime as a result of Phase 1 negotiations. Unit: seconds. Valid values: 0 to 86400. Default value: 86400.
* `ike_mode` - (Optional, Computed) IKE mode, the negotiation mode. Valid values: main and aggressive. Default value: main.
* `ike_pfs` - (Optional, Computed) The Diffie-Hellman key exchange algorithm used in the first stage negotiation. Valid values: group1, group2, group5, or group14. Default value: group2.
* `ike_version` - (Optional, Computed) The version of the IKE protocol. Value: ikev1 or ikev2. Default value: ikev1.
* `local_id` - (Optional, Computed) The identifier on the Alibaba Cloud side of the IPsec connection. The length is limited to 100 characters. The default value is leftId-not-exist
* `psk` - (Optional, Computed) A pre-shared key for authentication between the VPN gateway and the local data center. The key length is 1~100 characters.
  - If you do not specify a pre-shared key, the system randomly generates a 16-bit string as the pre-shared key. 
  - The pre-shared key of the IPsec-VPN connection must be the same as the authentication key of the on-premises data center. Otherwise, connections between the on-premises data center and the VPN gateway cannot be established.
* `remote_id` - (Optional, Computed) The identifier of the IPsec connection to the local data center. The length is limited to 100 characters. The default value is the IP address of the user gateway.

### `ipsec_config`

The ipsec_config supports the following:
* `ipsec_auth_alg` - (Optional, Computed) The authentication algorithm negotiated in the second stage. Valid values: md5, sha1, sha256, sha384, sha512. Default value: MD5.
* `ipsec_enc_alg` - (Optional, Computed) The encryption algorithm negotiated in the second stage. Valid values: aes, aes192, aes256, des, or 3des. Default value: aes.
* `ipsec_lifetime` - (Optional, Computed, Int) The life cycle of SA negotiated in the second stage. Unit: seconds. Value range: 0~86400. Default value: 86400.
* `ipsec_pfs` - (Optional, Computed) Diffie-Hellman Key Exchange Algorithm Used in Second Stage Negotiation

### `tunnel_options_specification`

The tunnel_options_specification supports the following:
* `customer_gateway_id` - (Required, Available since v1.246.0) The ID of the user gateway associated with the tunnel.

-> **NOTE:**  This parameter is required when creating a dual-tunnel mode IPsec-VPN connection.

* `enable_dpd` - (Optional, Computed, Available since v1.246.0) Whether the DPD (peer alive detection) function is enabled for the tunnel. Value:
  - `true` (default): enable the DPD function. IPsec initiator will send DPD message to check whether the peer device is alive. If the peer device does not receive a correct response within the set time, it is considered that the peer has been disconnected. IPsec will delete ISAKMP SA and the corresponding IPsec SA, and the security tunnel will also be deleted.
  - `false`: If the DPD function is disabled, the IPsec initiator does not send DPD detection packets.
* `enable_nat_traversal` - (Optional, Computed, Available since v1.246.0) Whether the NAT crossing function is enabled for the tunnel. Value:
  - `true` (default): Enables the NAT Traversal function. When enabled, the IKE negotiation process deletes the verification process of the UDP port number and realizes the discovery function of the NAT gateway device in the tunnel.
  - `false`: does not enable the NAT Traversal function.
* `tunnel_bgp_config` - (Optional, Computed, List, Available since v1.246.0) Add the BGP configuration for the tunnel.

-> **NOTE:**  After you enable the BGP function for IPsec connections (that is, specify `EnableTunnelsBgp` as `true`), you must configure this parameter.
 See [`tunnel_bgp_config`](#tunnel_options_specification-tunnel_bgp_config) below.
* `tunnel_ike_config` - (Optional, Computed, Set, Available since v1.246.0) Configuration information for the first phase negotiation. See [`tunnel_ike_config`](#tunnel_options_specification-tunnel_ike_config) below.
* `tunnel_index` - (Required, Int, Available since v1.246.0) The order in which the tunnel was created.
  - `1`: First tunnel.
  - `2`: The second tunnel.
* `tunnel_ipsec_config` - (Optional, Computed, List, Available since v1.246.0) Configuration information for the second-stage negotiation. See [`tunnel_ipsec_config`](#tunnel_options_specification-tunnel_ipsec_config) below.

### `tunnel_options_specification-tunnel_bgp_config`

The tunnel_options_specification-tunnel_bgp_config supports the following:
* `local_asn` - (Optional, Computed, Int, Available since v1.246.0) The number of the local (Alibaba Cloud) autonomous system of the tunnel. The value range of the autonomous system number is `1` to `4294967295`. Default value: `45104`.

-> **NOTE:**  We recommend that you use the private number of the autonomous system number to establish a BGP connection with Alibaba Cloud. The private number range of the autonomous system number please consult the document yourself.

* `local_bgp_ip` - (Optional, Computed, Available since v1.246.0) The local BGP address of the tunnel (on the Alibaba Cloud side). This address is an IP address in the BGP network segment.
* `tunnel_cidr` - (Optional, Computed, Available since v1.246.0) The BGP network segment of the tunnel. The network segment must be a network segment with a mask length of 30 in 169.254.0.0/16, and cannot be 169.254.0.0/30, 169.254.1.0/30, 169.254.2.0/30, 169.254.3.0/30, 169.254.4.0/30, 169.254.5.0/30, 169.254.6.0/30, and 169.254.169.252/30.

-> **NOTE:**  the network segments of two tunnels under an IPsec connection cannot be the same.


### `tunnel_options_specification-tunnel_ike_config`

The tunnel_options_specification-tunnel_ike_config supports the following:
* `ike_auth_alg` - (Optional, Computed, Available since v1.246.0) The authentication algorithm negotiated in the first stage. Values: `md5`, `sha1`, `sha256`, `sha384`, `sha512`. Default value: `sha1`.
* `ike_enc_alg` - (Optional, Computed, Available since v1.246.0) The encryption algorithm negotiated in the first stage. Value: `aes`, `aes192`, `aes256`, `des`, or `3des`. Default value: `aes`.
* `ike_lifetime` - (Optional, Computed, Int, Available since v1.246.0) The life cycle of SA negotiated in the first stage. Unit: seconds.
Value range: `0` to `86400`. Default value: `86400`.
* `ike_mode` - (Optional, Computed, Available since v1.246.0) IKE version of the negotiation mode. Value: `main` or `aggressive`. Default value: `main`.
  - `main`: main mode, high security during negotiation.
  - `aggressive`: Savage mode, fast negotiation and high negotiation success rate.
* `ike_pfs` - (Optional, Computed, Available since v1.246.0) The first stage negotiates the Diffie-Hellman key exchange algorithm used. Default value: `group2`.
Values: `group1`, `group2`, `group5`, `group14`.
* `ike_version` - (Optional, Computed, Available since v1.246.0) Version of the IKE protocol. Value: `ikev1` or `ikev2`. Default value: `ikev2`.
Compared with IKEv1, IKEv2 simplifies the SA negotiation process and provides better support for multiple network segments.
* `local_id` - (Optional, Computed, Available since v1.246.0) The identifier of the local end of the tunnel (Alibaba Cloud side), which is used for the first phase of negotiation. The length is limited to 100 characters and cannot contain spaces. The default value is the IP address of the tunnel.

  - *LocalId** supports the FQDN format. If you use the FQDN format, we recommend that you select `aggressive` (barbaric mode) as the negotiation mode.
* `psk` - (Optional, Computed, Available since v1.246.0) The pre-shared key is used for identity authentication between the tunnel and the tunnel peer.
  - The key can be 1 to 100 characters in length. It supports numbers, upper and lower case English letters, and characters on the right. It cannot contain spaces. '''~! \'@#$%^& *()_-+ ={}[]|;:',./? '''
  - If you do not specify a pre-shared key, the system randomly generates a 16-bit string as the pre-shared key. 

-> **NOTE:**  The pre-shared key of the tunnel and the tunnel peer must be the same, otherwise the system cannot establish the tunnel normally.

* `remote_id` - (Optional, Computed, Available since v1.246.0) Identifier of the tunnel peer, which is used for the first-stage negotiation. The length is limited to 100 characters and cannot contain spaces. The default value is the IP address of the user gateway associated with the tunnel.
  - `RemoteId` supports the FQDN format. If you use the FQDN format, we recommend that you select `aggressive` (barbaric mode) as the negotiation mode.

### `tunnel_options_specification-tunnel_ipsec_config`

The tunnel_options_specification-tunnel_ipsec_config supports the following:
* `ipsec_auth_alg` - (Optional, Computed, Available since v1.246.0) The second stage negotiated authentication algorithm.
Values: `md5`, `sha1`, `sha256`, `sha384`, `sha512`. Default value: `sha1`.
* `ipsec_enc_alg` - (Optional, Computed, Available since v1.246.0) The encryption algorithm negotiated in the second stage. Value: `aes`, `aes192`, `aes256`, `des`, or `3des`. Default value: `aes`.
* `ipsec_lifetime` - (Optional, Computed, Int, Available since v1.246.0) The life cycle of SA negotiated in the second stage. Unit: seconds.
Value range: `0` to `86400`. Default value: `86400`.
* `ipsec_pfs` - (Optional, Computed, Available since v1.246.0) The second stage negotiates the Diffie-Hellman key exchange algorithm used. Default value: `group2`.
Values: `disabled`, `group1`, `group2`, `group5`, `group14`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `bgp_config` - Bgp configuration information.-tunnel mode.
  * `status` - The negotiation status of Bgp.
* `create_time` - The creation time of the resource
* `health_check_config` - This parameter is supported if you create an vpn attachment in single-tunnel mode.
  * `status` - health check status
* `status` - The status of the resource
* `tunnel_options_specification` - Configure the tunnel._options_specification` array when you create a vpn attachment in dual-tunnel mode.-tunnel mode, you must add both tunnels for the vpn attachment to ensure that the vpn attachment has link redundancy. Only two tunnels can be added to a vpn attachment.
  * `internet_ip` - The local internet IP in Tunnel.
  * `role` - The role of Tunnel.
  * `state` - The state of Tunnel.
  * `status` - The negotiation status of Tunnel. 
  * `tunnel_bgp_config` - Add the BGP configuration for the tunnel.
    * `bgp_status` - BGP status.
    * `peer_asn` - Peer asn.
    * `peer_bgp_ip` - Peer bgp ip.
  * `tunnel_id` - The tunnel ID of IPsec-VPN connection.
  * `zone_no` - The zoneNo of tunnel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpn Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Vpn Attachment.
* `update` - (Defaults to 5 mins) Used when update the Vpn Attachment.

## Import

V P N Gateway Vpn Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_gateway_vpn_attachment.example <id>
```