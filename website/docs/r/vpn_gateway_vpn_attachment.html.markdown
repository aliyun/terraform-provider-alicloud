---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_gateway_vpn_attachment"
sidebar_current: "docs-alicloud-resource-vpn-gateway-vpn-attachment"
description: |-
  Provides a Alicloud VPN Gateway Vpn Attachment resource.
---

# alicloud\_vpn\_gateway\_vpn\_attachment

Provides a VPN Gateway Vpn Attachment resource.

For information about VPN Gateway Vpn Attachment and how to use it, see [What is Vpn Attachment](https://www.alibabacloud.com/help/zh/virtual-private-cloud/latest/createvpnattachment).

-> **NOTE:** Available in v1.181.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpn_customer_gateway" "default" {
  name        = "${var.name}"
  ip_address  = "42.104.22.210"
  asn         = "45014"
  description = "testAccVpnConnectionDesc"
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

* `bgp_config` - (Optional, Computed) Bgp configuration information. See the following `Block bgp_config`.
* `customer_gateway_id` - (Required, ForceNew) The ID of the customer gateway.
* `effect_immediately` - (Optional) Indicates whether IPsec-VPN negotiations are initiated immediately. Valid values.
* `enable_dpd` - (Optional, Computed) Whether to enable the DPD (peer survival detection) function.
* `enable_nat_traversal` - (Optional, Computed) Allow NAT penetration.
* `health_check_config` - (Optional, Computed) Health check configuration information. See the following `Block health_check_config`.
* `ike_config` - (Optional, Computed) Configuration negotiated in the second stage. See the following `Block ike_config`.
* `ipsec_config` - (Optional, Computed) Configuration negotiated in the second stage. See the following `Block ipsec_config`.
* `local_subnet` - (Required) The CIDR block of the virtual private cloud (VPC).
* `network_type` - (Optional, Computed, ForceNew) The network type of the IPsec connection. Valid values: `public`, `private`.
* `remote_subnet` - (Required) The CIDR block of the on-premises data center.
* `vpn_attachment_name` - (Optional) The name of the vpn attachment.

#### Block ipsec_config

The ipsec_config supports the following: 

* `ipsec_auth_alg` - (Optional) The authentication algorithm of phase-two negotiation. Valid value: md5 | sha1 | sha256 | sha384 | sha512 |. Default value: sha1
* `ipsec_enc_alg` - (Optional) The encryption algorithm of phase-two negotiation. Valid value: aes | aes192 | aes256 | des | 3des. Default value: aes
* `ipsec_lifetime` - (Optional) The SA lifecycle as the result of phase-two negotiation. The valid value is [0, 86400], the unit is second and the default value is 86400.
* `ipsec_pfs` - (Optional) The Diffie-Hellman key exchange algorithm used by phase-two negotiation. Valid value: group1 | group2 | group5 | group14 | group24| disabled. Default value: group2

#### Block ike_config

The ike_config supports the following: 

* `ike_auth_alg` - (Optional) IKE authentication algorithm supports sha1 and MD5.
* `ike_enc_alg` - (Optional) The encryption algorithm of phase-one negotiation. Valid value: aes | aes192 | aes256 | des | 3des. Default Valid value: aes.
* `ike_lifetime` - (Optional) The SA lifecycle as the result of phase-one negotiation. The valid value of n is [0, 86400], the unit is second and the default value is 86400.
* `ike_mode` - (Optional) The negotiation mode of IKE V1. Valid value: main (main mode) | aggressive (aggressive mode). Default value: `main`.
* `ike_pfs` - (Optional) The Diffie-Hellman key exchange algorithm used by phase-one negotiation. Valid value: group1 | group2 | group5 | group14 | group24. Default value: group2
* `ike_version` - (Optional) The version of the IKE protocol. Valid value: `ikev1`, `ikev2`. Default value: `ikev1`.
* `local_id` - (Optional) The local ID, which supports the FQDN and IP formats. The current VPN gateway IP address is selected by default.
* `psk` - (Optional) Used for authentication between the IPsec VPN gateway and the customer gateway.
* `remote_id` - (Optional) The peer ID, which supports FQDN and IP formats. By default, the IP address of the currently selected user gateway.

#### Block health_check_config

The health_check_config supports the following: 

* `dip` - (Optional) The destination IP address that is used for health checks.
* `enable` - (Optional) Specifies whether to enable health checks.
* `interval` - (Optional) The interval between two consecutive health checks. Unit: seconds.
* `retry` - (Optional) The maximum number of health check retries.
* `sip` - (Optional) The source IP address that is used for health checks.
* `policy` - (Optional) Whether to revoke the published route when the health check fails. Valid values: `revoke_route` or `reserve_route`.

#### Block bgp_config

The bgp_config supports the following: 

* `enable` - (Optional, Computed) Whether to enable BGP.
* `local_asn` - (Optional, Computed) The ASN on the Alibaba Cloud side.
* `tunnel_cidr` - (Optional, Computed) The CIDR block of the IPsec tunnel. The CIDR block belongs to 169.254.0.0/16. The mask of the CIDR block is 30 bits in length.
* `local_bgp_ip` - (Optional, Computed)  The BGP IP address on the Alibaba Cloud side.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vpn Attachment.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Vpn Attachment.
* `update` - (Defaults to 1 mins) Used when updating the Vpn Attachment.
* `delete` - (Defaults to 1 mins) Used when deleting the Vpn Attachment.


## Import

VPN Gateway Vpn Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_vpn_gateway_vpn_attachment.example <id>
```