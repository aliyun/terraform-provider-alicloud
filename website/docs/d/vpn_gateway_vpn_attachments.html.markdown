---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_vpn_attachments"
sidebar_current: "docs-alicloud-datasource-vpn-gateway-vpn-attachments"
description: |-
  Provides a list of Vpn Gateway Vpn Attachment owned by an Alibaba Cloud account.
---

# alicloud_vpn_gateway_vpn_attachments

This data source provides Vpn Gateway Vpn Attachment available to the user.[What is Vpn Attachment](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateVpnAttachment)

-> **NOTE:** Available since v1.245.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-huhehaote"
}

variable "region_id" {
  default = "cn-huhehaote"
}

variable "az2" {
  default = "cn-huhehaote-b"
}

variable "name" {
  default = "example_amp"
}

variable "az1" {
  default = "cn-huhehaote-a"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpn_customer_gateway" "cgw1" {
  ip_address = "2.2.2.2"
  asn        = "1219001"
}

resource "alicloud_vpn_customer_gateway" "cgw2" {
  ip_address            = "43.43.3.22"
  asn                   = "44331"
  customer_gateway_name = "example_amp"
}


resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  local_subnet        = "0.0.0.0/0"
  enable_tunnels_bgp  = true
  vpn_attachment_name = var.name
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.cgw1.id
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_index         = "1"
    tunnel_bgp_config {
      local_asn    = "1219001"
      local_bgp_ip = "169.254.10.1"
      tunnel_cidr  = "169.254.10.0/30"
    }
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes"
      ike_lifetime = "86100"
      ike_mode     = "main"
      ike_pfs      = "group2"
      ike_version  = "ikev1"
      local_id     = "1.1.1.1"
      psk          = "12345678"
      remote_id    = "2.2.2.2"
    }
    tunnel_ipsec_config {
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes"
      ipsec_lifetime = "86200"
      ipsec_pfs      = "group5"
    }
  }
  tunnel_options_specification {
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_index         = "2"
    tunnel_bgp_config {
      local_asn    = "1219001"
      local_bgp_ip = "169.254.20.1"
      tunnel_cidr  = "169.254.20.0/30"
    }
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes"
      ike_lifetime = "86400"
      ike_mode     = "main"
      ike_pfs      = "group5"
      ike_version  = "ikev2"
      local_id     = "4.4.4.4"
      psk          = "32333442"
      remote_id    = "5.5.5.5"
    }
    tunnel_ipsec_config {
      ipsec_auth_alg = "sha256"
      ipsec_enc_alg  = "aes"
      ipsec_lifetime = "86400"
      ipsec_pfs      = "group5"
    }
    customer_gateway_id = alicloud_vpn_customer_gateway.cgw1.id
  }
  remote_subnet = "0.0.0.0/0"
  network_type  = "public"
}

data "alicloud_vpn_gateway_vpn_attachments" "default" {
  ids        = ["${alicloud_vpn_gateway_vpn_attachment.default.id}"]
  name_regex = alicloud_vpn_gateway_vpn_attachment.default.vpn_attachment_name
}

output "alicloud_vpn_gateway_vpn_attachment_example_id" {
  value = data.alicloud_vpn_gateway_vpn_attachments.default.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:
* `page_number` - (ForceNew, Optional) Current page number.
* `page_size` - (ForceNew, Optional) Number of records per page.
* `ids` - (Optional, ForceNew, Computed) A list of Vpn Attachment IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `init`, `active`, `attaching`, `attached`, `detaching`, `financialLocked`, `provisioning`, `updating`, `upgrading`, `deleted`.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Vpn Attachment IDs.
* `names` - A list of name of Vpn Attachments.
* `attachments` - A list of Vpn Attachment Entries. Each element contains the following attributes:
  * `attach_type` - attach type- **CEN**: indicates that the IPsec-VPN connection is associated with a transit router of a Cloud Enterprise Network (CEN) instance.- **NO_ASSOCIATED**: indicates that the IPsec-VPN connection is not associated with any resource.
  * `bgp_config` - Bgp configuration information.- This parameter is supported when you create an vpn attachment in single-tunnel mode.
    * `local_asn` - The autonomous system number on the Alibaba Cloud side. The value range of autonomous system number is 1~4294967295. Default value: 45104.
    * `local_bgp_ip` - The BGP address on the Alibaba Cloud side. This address is an IP address in the IPsec tunnel network segment.- Before adding the BGP configuration, we recommend that you understand the working mechanism and usage restrictions of the BGP dynamic routing function. For more information, see BGP Dynamic Routing Bulletin.- We recommend that you use the private number of the autonomous system number to establish a BGP connection with Alibaba Cloud. Please refer to the documentation for the private number range of the autonomous system number.
    * `status` - The negotiation status of Bgp.- success- failed .
    * `tunnel_cidr` - IPsec tunnel network segment. This network segment must be a network segment with a mask length of 30 within 169.254.0.0/16.
  * `connection_status` - IPsec connection status- **ike_sa_not_established**: Phase 1 negotiations failed.- **ike_sa_established**: Phase 1 negotiations succeeded.- **ipsec_sa_not_established**: Phase 2 negotiations failed.- **ipsec_sa_established**: Phase 2 negotiations succeeded.
  * `create_time` - The creation time of the resource
  * `customer_gateway_id` - Customer gateway ID.- This parameter is required when creating a single-tunnel mode vpn attachment.
  * `effect_immediately` - Specifies whether to immediately start IPsec negotiations after the configuration takes effect. Valid values:- **true**: immediately starts IPsec negotiations after the configuration is complete.- **false** (default): starts IPsec negotiations when inbound traffic is received.
  * `enable_dpd` - This parameter is supported if you create an vpn attachment in single-tunnel mode.Whether to enable the DPD (peer survival detection) function.- true (default): enables DPD. The initiator of the IPsec-VPN connection sends DPD packets to check the existence and availability of the peer. If no feedback is received from the peer within the specified period of time, the connection fails. In this case, ISAKMP SA and IPsec SA are deleted along with the security tunnel.- false: disables DPD. The initiator of the IPsec-VPN connection does not send DPD packets.
  * `enable_nat_traversal` - This parameter is supported if you create an vpn attachment in single-tunnel mode.Specifies whether to enable NAT traversal. Valid values:- true (default): enables NAT traversal. After NAT traversal is enabled, the initiator does not check the UDP ports during IKE negotiations and can automatically discover NAT gateway devices along the vpn attachment tunnel.- false: disables NAT traversal.
  * `enable_tunnels_bgp` - You can configure this parameter when you create a vpn attachment in dual-tunnel mode.Whether to enable the BGP function for the tunnel. Value: **true** or **false** (default).> before adding BGP configuration, we recommend that you understand the working mechanism and usage restrictions of the BGP dynamic routing function. 
  * `health_check_config` - This parameter is supported if you create an vpn attachment in single-tunnel mode.Health check configuration information.
    * `dip` - Target IP.
    * `enable` - Whether health check is enabled:-**false**: not enabled. -**true**: enabled.
    * `interval` - The health check retry interval, in seconds.
    * `policy` - Whether to revoke the published route when the health check fails- **revoke_route**(default): withdraws published routes.- **reserve_route**: does not withdraw published routes.
    * `retry` - Number of retries for health check.
    * `sip` - SOURCE IP.
    * `status` - health check status.
  * `ike_config` - The configurations of Phase 1 negotiations. - This parameter is supported if you create an vpn attachment in single-tunnel mode.
    * `ike_auth_alg` - The authentication algorithm negotiated in the first stage. Valid values: md5, sha1, sha256, sha384, sha512. Default value: md5.
    * `ike_enc_alg` - The encryption algorithm that is used in Phase 1 negotiations. Valid values: aes, aes192, aes256, des, and 3des. Default value: aes.
    * `ike_lifetime` - The SA lifetime as a result of Phase 1 negotiations. Unit: seconds. Valid values: 0 to 86400. Default value: 86400.
    * `ike_mode` - IKE mode, the negotiation mode. Valid values: main and aggressive. Default value: main.
    * `ike_pfs` - The Diffie-Hellman key exchange algorithm used in the first stage negotiation. Valid values: group1, group2, group5, or group14. Default value: group2.
    * `ike_version` - The version of the IKE protocol. Value: ikev1 or ikev2. Default value: ikev1.
    * `local_id` - The identifier on the Alibaba Cloud side of the IPsec connection. The length is limited to 100 characters. The default value is leftId-not-exist.
    * `psk` - A pre-shared key for authentication between the VPN gateway and the local data center. The key length is 1~100 characters.- If you do not specify a pre-shared key, the system randomly generates a 16-bit string as the pre-shared key. - The pre-shared key of the IPsec-VPN connection must be the same as the authentication key of the on-premises data center. Otherwise, connections between the on-premises data center and the VPN gateway cannot be established.
    * `remote_id` - The identifier of the IPsec connection to the local data center. The length is limited to 100 characters. The default value is the IP address of the user gateway.
  * `internet_ip` - VPN gateway IP
  * `ipsec_config` - Configuration negotiated in the second stage. - This parameter is supported if you create an vpn attachment in single-tunnel mode.
    * `ipsec_auth_alg` - The authentication algorithm negotiated in the second stage. Valid values: md5, sha1, sha256, sha384, sha512. Default value: MD5.
    * `ipsec_enc_alg` - The encryption algorithm negotiated in the second stage. Valid values: aes, aes192, aes256, des, or 3des. Default value: aes.
    * `ipsec_lifetime` - The life cycle of SA negotiated in the second stage. Unit: seconds. Value range: 0~86400. Default value: 86400.
    * `ipsec_pfs` - Diffie-Hellman Key Exchange Algorithm Used in Second Stage Negotiation.
  * `local_subnet` - The CIDR block on the VPC side. The CIDR block is used in Phase 2 negotiations.Separate multiple CIDR blocks with commas (,). Example: 192.168.1.0/24,192.168.2.0/24.The following routing modes are supported:- If you set LocalSubnet and RemoteSubnet to 0.0.0.0/0, the routing mode of the IPsec-VPN connection is set to Destination Routing Mode.- If you set LocalSubnet and RemoteSubnet to specific CIDR blocks, the routing mode of the IPsec-VPN connection is set to Protected Data Flows.
  * `network_type` - network type- **public** (default)- **private**
  * `remote_subnet` - The CIDR block on the data center side. This CIDR block is used in Phase 2 negotiations.Separate multiple CIDR blocks with commas (,). Example: 192.168.3.0/24,192.168.4.0/24.The following routing modes are supported:- If you set LocalSubnet and RemoteSubnet to 0.0.0.0/0, the routing mode of the IPsec-VPN connection is set to Destination Routing Mode.- If you set LocalSubnet and RemoteSubnet to specific CIDR blocks, the routing mode of the IPsec-VPN connection is set to Protected Data Flows.
  * `resource_group_id` - The ID of the resource group
  * `status` - The status of the resource
  * `tags` - Tags
  * `tunnel_options_specification` - Configure the tunnel.-You can configure parameters in the **tunnel_options_specification** array when you create a vpn attachment in dual-tunnel mode.-When creating a vpn attachment in dual-tunnel mode, you must add both tunnels for the vpn attachment to ensure that the vpn attachment has link redundancy. Only two tunnels can be added to a vpn attachment.
    * `customer_gateway_id` - The ID of the user gateway associated with the tunnel.> This parameter is required when creating a dual-tunnel mode IPsec-VPN connection.
    * `enable_dpd` - Whether the DPD (peer alive detection) function is enabled for the tunnel. Value:-**true** (default): enable the DPD function. IPsec initiator will send DPD message to check whether the peer device is alive. If the peer device does not receive a correct response within the set time, it is considered that the peer has been disconnected. IPsec will delete ISAKMP SA and the corresponding IPsec SA, and the security tunnel will also be deleted.-**false**: If the DPD function is disabled, the IPsec initiator does not send DPD detection packets.
    * `enable_nat_traversal` - Whether the NAT crossing function is enabled for the tunnel. Value:-**true** (default): Enables the NAT Traversal function. When enabled, the IKE negotiation process deletes the verification process of the UDP port number and realizes the discovery function of the NAT gateway device in the tunnel.-**false**: does not enable the NAT Traversal function.
    * `internet_ip` - The local internet IP in Tunnel.
    * `role` - The role of Tunnel.
    * `state` - The state of Tunnel.
    * `status` - The negotiation status of Tunnel. - **ike_sa_not_established**: Phase 1 negotiations failed.- **ike_sa_established**: Phase 1 negotiations succeeded.- **ipsec_sa_not_established**: Phase 2 negotiations failed.- **ipsec_sa_established**: Phase 2 negotiations succeeded.
    * `tunnel_bgp_config` - Add the BGP configuration for the tunnel.> After you enable the BGP function for IPsec connections (that is, specify **EnableTunnelsBgp** as **true**), you must configure this parameter.
      * `bgp_status` - BGP status.
      * `local_asn` - The number of the local (Alibaba Cloud) autonomous system of the tunnel. The value range of the autonomous system number is **1** to **4294967295**. Default value: **45104**.> We recommend that you use the private number of the autonomous system number to establish a BGP connection with Alibaba Cloud. The private number range of the autonomous system number please consult the document yourself.
      * `local_bgp_ip` - The local BGP address of the tunnel (on the Alibaba Cloud side). This address is an IP address in the BGP network segment.
      * `peer_asn` - Peer asn.
      * `peer_bgp_ip` - Peer bgp ip.
      * `tunnel_cidr` - The BGP network segment of the tunnel. The network segment must be a network segment with a mask length of 30 in 169.254.0.0/16, and cannot be 169.254.0.0/30, 169.254.1.0/30, 169.254.2.0/30, 169.254.3.0/30, 169.254.4.0/30, 169.254.5.0/30, 169.254.6.0/30, and 169.254.169.252/30.> the network segments of two tunnels under an IPsec connection cannot be the same.
    * `tunnel_id` - The tunnel ID of IPsec-VPN connection.
    * `tunnel_ike_config` - Configuration information for the first phase negotiation.
      * `ike_auth_alg` - The authentication algorithm negotiated in the first stage. Values: **md5**, **sha1**, **sha256**, **sha384**, **sha512**. Default value: **sha1**.
      * `ike_enc_alg` - The encryption algorithm negotiated in the first stage. Value: **aes**, **aes192**, **aes256**, **des**, or **3des**. Default value: **aes**.
      * `ike_lifetime` - The life cycle of SA negotiated in the first stage. Unit: seconds.Value range: **0** to **86400**. Default value: **86400**.
      * `ike_mode` - IKE version of the negotiation mode. Value: **main** or **aggressive**. Default value: **main**.-**main**: main mode, high security during negotiation.-**aggressive**: Savage mode, fast negotiation and high negotiation success rate.
      * `ike_pfs` - The first stage negotiates the Diffie-Hellman key exchange algorithm used. Default value: **group2**.Values: **group1**, **group2**, **group5**, **group14**.
      * `ike_version` - Version of the IKE protocol. Value: **ikev1** or **ikev2**. Default value: **ikev2**.Compared with IKEv1, IKEv2 simplifies the SA negotiation process and provides better support for multiple network segments.
      * `local_id` - The identifier of the local end of the tunnel (Alibaba Cloud side), which is used for the first phase of negotiation. The length is limited to 100 characters and cannot contain spaces. The default value is the IP address of the tunnel.**LocalId** supports the FQDN format. If you use the FQDN format, we recommend that you select **aggressive** (barbaric mode) as the negotiation mode.
      * `psk` - The pre-shared key is used for identity authentication between the tunnel and the tunnel peer.-The key can be 1 to 100 characters in length. It supports numbers, upper and lower case English letters, and characters on the right. It cannot contain spaces. '''~! \'@#$%^& *()_-+ ={}[]|;:',./? '''-If you do not specify a pre-shared key, the system randomly generates a 16-bit string as the pre-shared key. > The pre-shared key of the tunnel and the tunnel peer must be the same, otherwise the system cannot establish the tunnel normally.
      * `remote_id` - Identifier of the tunnel peer, which is used for the first-stage negotiation. The length is limited to 100 characters and cannot contain spaces. The default value is the IP address of the user gateway associated with the tunnel.- **RemoteId** supports the FQDN format. If you use the FQDN format, we recommend that you select **aggressive** (barbaric mode) as the negotiation mode.
    * `tunnel_index` - The order in which the tunnel was created.-**1**: First tunnel.-**2**: The second tunnel.
    * `tunnel_ipsec_config` - Configuration information for the second-stage negotiation.
      * `ipsec_auth_alg` - The second stage negotiated authentication algorithm.Values: **md5**, **sha1**, **sha256**, **sha384**, **sha512**. Default value: **sha1**.
      * `ipsec_enc_alg` - The encryption algorithm negotiated in the second stage. Value: **aes**, **aes192**, **aes256**, **des**, or **3des**. Default value: **aes**.
      * `ipsec_lifetime` - The life cycle of SA negotiated in the second stage. Unit: seconds.Value range: **0** to **86400**. Default value: **86400**.
      * `ipsec_pfs` - The second stage negotiates the Diffie-Hellman key exchange algorithm used. Default value: **group2**.Values: **disabled**, **group1**, **group2**, **group5**, **group14**.
    * `zone_no` - The zoneNo of tunnel.
  * `vpn_attachment_name` - vpn attachment name
  * `vpn_connection_id` - The first ID of the resource
  * `id` - The ID of the resource supplied above.
