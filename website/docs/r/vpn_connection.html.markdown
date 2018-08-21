---
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_connection"
sidebar_current: "docs-alicloud-resource-vpn-connection"
description: |-
  Provides a Alicloud VPN connection resource.
---

# alicloud\_vpn

Provides a VPN connection resource.

~> **NOTE:** Terraform will auto build vpn connection while it uses `alicloud_vpn_connection` to build a vpn connection resource.
             The vpn connection depends on VPN and VPN customer gateway.

## Example Usage

Basic Usage

```
resource "alicloud_vpn" "foo" {
        name = "testAccVpnConfig_create"
        vpc_id = "vpc-2ze9wy916mfwpwbf6hx4u"
		bandwidth = "10"
        enable_ssl = true
        instance_charge_type = "postpaid"
        auto_pay = true
		description = "test_create_description"
}

resource "alicloud_vpn_customer_gateway" "foo" {
	name = "testAccVpnCgwName"
	ip_address = "41.104.22.229"
	description = "testAccVpnCgwDesc"
}

resource "alicloud_vpn_connection" "foo" {
	name = "tf-vco_test1"
	vpn_gateway_id = "${alicloud_vpn.foo.id}"
	customer_gateway_id = "${alicloud_vpn_customer_gateway.foo.id}"
	local_subnet = "172.16.0.0/16"
	remote_subnet = "10.0.0.0/8"
	effect_immediately = true
	ike_config = "{\"LocalId\":\"${alicloud_vpn.foo.internet_ip}\",\"IkeAuthAlg\":\"sha1\",\"IkePfs\":\"group2\",\"IkeMode\":\"main\",\"IkeEncAlg\":\"aes\",\"Psk\":\"tf-testvpn1\",\"RemoteId\":\"{$alicloud_vpn_customer_gateway.ip_address}\",\"IkeVersion\":\"ikev2\",\"IkeLifetime\":3600}"
	ipsec_config = "{\"IpsecAuthAlg\":\"sha1\",\"IpsecPfs\":\"group2\",\"IpsecEncAlg\":\"aes\",\"IpsecLifetime\":7200}"
}
```
## Argument Reference

The following arguments are supported:
* `name` - (Optional) The name of the VPN connection. Defaults to null.
* `vpn_gateway_id` - (Required) The VPN instance ID.
* `customer_gateway_id` - (Required) The VPN customer gateway instance id.
* `local_subnet` - (Required) The local subnet of the VPN connection.
* `remote_subnet` - (Required) The remote subnet of the VPN connection.
* `effect_immediately` - (Optional) The Config should take effect immediately.
* `ike_config` - (Required) The ike config of the VPN connection, all fields in the example above except LocalId and RemoteId are needed.
* `ipsec_config` - (Required) The ipsec config of the VPN connection.
* `description` - (Optional) The description of the VPN connection.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPN connection id.
* `status` - The status of VPN connection.




