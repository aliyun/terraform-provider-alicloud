---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_vpn_attachments"
sidebar_current: "docs-alicloud-datasource-vpn-gateway-vpn-attachments"
description: |-
  Provides a list of Vpn Gateway Vpn Attachments to the user.
---

# alicloud_vpn_gateway_vpn_attachments

This data source provides the Vpn Gateway Vpn Attachments of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.181.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpn_gateway_vpn_attachments" "ids" {}
output "vpn_gateway_vpn_attachment_id_1" {
  value = data.alicloud_vpn_gateway_vpn_attachments.ids.attachments.0.id
}

data "alicloud_vpn_gateway_vpn_attachments" "nameRegex" {
  name_regex = "^my-VpnAttachment"
}
output "vpn_gateway_vpn_attachment_id_2" {
  value = data.alicloud_vpn_gateway_vpn_attachments.nameRegex.attachments.0.id
}
output "local_id" {
  value = data.alicloud_vpn_gateway_vpn_attachments.vpn_attachments.attachments.0.ike_config.0.local_id
}
output "internet_ip" {
  value = data.alicloud_vpn_gateway_vpn_attachments.vpn_attachments.attachments.0.internet_ip
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Vpn Attachment IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Vpn Attachment name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `init`, `active`, `attaching`, `attached`, `detaching`, `financialLocked`, `provisioning`, `updating`, `upgrading`, `deleted`.
* `vpn_gateway_id` - (Deprecated) The parameter 'vpn_gateway_id' has been deprecated from 1.194.0.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Vpn Attachment names.
* `attachments` - A list of Vpn Gateway Vpn Attachments. Each element contains the following attributes:
  * `bgp_config` - The configurations of the BGP routing protocol.
    * `local_asn` - The ASN on the Alibaba Cloud side.
    * `local_bgp_ip` - The BGP IP address on the Alibaba Cloud side.
    * `status` - The negotiation status of the BGP routing protocol.
    * `tunnel_cidr` - The CIDR block of the IPsec tunnel. The CIDR block belongs to 169.254.0.0/16. The mask of the CIDR block is 30 bits in length.
  * `create_time` - The creation time of the resource.
  * `customer_gateway_id` - The ID of the customer gateway.
  * `effect_immediately` - Indicates whether IPsec-VPN negotiations are initiated immediately. Valid values.
  * `health_check_config` - The health check configurations.
    * `dip` - The destination IP address.
    * `status` - The status of the health check.
    * `interval` - The interval between two consecutive health checks. Unit: seconds.
    * `retry` - The maximum number of health check retries.
    * `sip` - The source IP address.
    * `enable` - Specifies whether to enable health checks.
    * `policy` - (Optional) Whether to revoke the published route when the health check fails.
  * `id` - The ID of the Vpn Attachment.
  * `ike_config` - Configuration negotiated in the second stage.
    * `remote_id` - The identifier of the peer. The default value is the IP address of the VPN gateway. The value can be a fully qualified domain name (FQDN) or an IP address.
    * `ike_lifetime` - The IKE lifetime. Unit: seconds.
    * `ike_pfs` - The DH group.
    * `local_id` - The local ID, which supports the FQDN and IP formats. The current VPN gateway IP address is selected by default. The alicloud_cen_transit_router_vpn_attachment resource will not have a value until after it is created.
    * `psk` - The pre-shared key.
    * `ike_auth_alg` - The IKE authentication algorithm.
    * `ike_enc_alg` - The IKE encryption algorithm.
    * `ike_mode` - The IKE negotiation mode.
    * `ike_version` - The version of the IKE protocol.
  * `ipsec_config` - The configuration of Phase 2 negotiations.
    * `ipsec_enc_alg` - The IPsec encryption algorithm.
    * `ipsec_lifetime` - The IPsec lifetime. Unit: seconds.
    * `ipsec_pfs` - The DH group.
    * `ipsec_auth_alg` - The IPsec authentication algorithm.
  * `local_subnet` - The CIDR block of the virtual private cloud (VPC).
  * `network_type` - The network type.
  * `remote_subnet` - The CIDR block of the on-premises data center.
  * `status` - The status of the resource.
  * `vpn_attachment_name` - The name of the IPsec-VPN connection.
  * `vpn_connection_id` - The first ID of the resource.
  * `connection_status` - The status of the IPsec-VPN connection. 
  * `internet_ip` - The internet ip of the resource. The alicloud_cen_transit_router_vpn_attachment resource will not have a value until after it is created.