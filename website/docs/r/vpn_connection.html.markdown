---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_connection"
sidebar_current: "docs-alicloud-resource-vpn-connection"
description: |-
  Provides a Alicloud VPN connection resource.
---

# alicloud\_vpn\_connection

Provides a VPN connection resource.

-> **NOTE:** Terraform will auto build vpn connection while it uses `alicloud_vpn_connection` to build a vpn connection resource.
             The vpn connection depends on VPN and VPN customer gateway.

## Example Usage

Basic Usage

```
resource "alicloud_vpn_gateway" "foo" {
  name                 = "testAccVpnConfig_create"
  vpc_id               = "vpc-fake-id"
  bandwidth            = "10"
  enable_ssl           = true
  instance_charge_type = "PostPaid"
  description          = "test_create_description"
}

resource "alicloud_vpn_customer_gateway" "foo" {
  name        = "testAccVpnCgwName"
  ip_address  = "42.104.22.228"
  description = "testAccVpnCgwDesc"
}

resource "alicloud_vpn_connection" "foo" {
  name                = "tf-vco_test1"
  vpn_gateway_id      = "${alicloud_vpn_gateway.foo.id}"
  customer_gateway_id = "${alicloud_vpn_customer_gateway.foo.id}"
  local_subnet        = ["172.16.0.0/24", "172.16.1.0/24"]
  remote_subnet       = ["10.0.0.0/24", "10.0.1.0/24"]
  effect_immediately  = true
  ike_config {
    ike_auth_alg  = "md5"
    ike_enc_alg   = "des"
    ike_version   = "ikev1"
    ike_mode      = "main"
    ike_lifetime  = 86400
    psk           = "tf-testvpn2"
    ike_pfs       = "group1"
    ike_remote_id = "testbob2"
    ike_local_id  = "testalice2"
    }
  ipsec_config {
    ipsec_pfs      = "group5"
    ipsec_enc_alg  = "des"
    ipsec_auth_alg = "md5"
    ipsec_lifetime = 8640
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the IPsec connection.
* `vpn_gateway_id` - (Required, ForceNew) The ID of the VPN gateway.
* `customer_gateway_id` - (Required, ForceNew) The ID of the customer gateway.
* `local_subnet` - (Required, Type:Set) The CIDR block of the VPC to be connected with the local data center. This parameter is used for phase-two negotiation.
* `remote_subnet` - (Required, Type:Set) The CIDR block of the local data center. This parameter is used for phase-two negotiation.
* `effect_immediately` - (Optional) Whether to delete a successfully negotiated IPsec tunnel and initiate a negotiation again. Valid value:true,false.
* `ike_config` - (Optional) The configurations of phase-one negotiation.
* `ipsec_config` - (Optional) The configurations of phase-two negotiation.

### Block ike_config

The ike_config mapping supports the following:

* `psk` - (Optional) Used for authentication between the IPsec VPN gateway and the customer gateway.
* `ike_version` - (Optional) The version of the IKE protocol. Valid value: ikev1 | ikev2. Default value: ikev1
* `ike_mode` - (Optional) The negotiation mode of IKE V1. Valid value: main (main mode) | aggressive (aggressive mode). Default value: main
* `ike_enc_alg` - (Optional) The encryption algorithm of phase-one negotiation. Valid value: aes | aes192 | aes256 | des | 3des. Default Valid value: aes
* `ike_auth_alg` - (Optional) The authentication algorithm of phase-one negotiation. Valid value: md5 | sha1 | sha256 | sha384 | sha512 |. Default value: sha1
* `ike_pfs` - (Optional) The Diffie-Hellman key exchange algorithm used by phase-one negotiation. Valid value: group1 | group2 | group5 | group14 | group24. Default value: group2
* `ike_lifetime` - (Optional) The SA lifecycle as the result of phase-one negotiation. The valid value of n is [0, 86400], the unit is second and the default value is 86400.
* `ike_local_id` - (Optional) The identification of the VPN gateway.
* `ike_remote_id` - (Optional) The identification of the customer gateway.

### Block ipsec_config

The ipsec_config mapping supports the following:

* `ipsec_enc_alg` - (Optional) The encryption algorithm of phase-two negotiation. Valid value: aes | aes192 | aes256 | des | 3des. Default value: aes
* `ipsec_auth_alg` - (Optional) The authentication algorithm of phase-two negotiation. Valid value: md5 | sha1 | sha256 | sha384 | sha512 |. Default value: sha1
* `ipsec_pfs` - (Optional) The Diffie-Hellman key exchange algorithm used by phase-two negotiation. Valid value: group1 | group2 | group5 | group14 | group24| disabled. Default value: group2
* `ipsec_lifetime` - (Optional)  The SA lifecycle as the result of phase-two negotiation. The valid value is [0, 86400], the unit is second and the default value is 86400.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN connection id.
* `status` - The status of VPN connection.
* `ike_config` - The configurations of phase-one negotiation.
* `ipsec_config` - The configurations of phase-two negotiation.

## Import

VPN connection can be imported using the id, e.g.

```
$ terraform import alicloud_vpn_connection.example vco-abc123456
```


