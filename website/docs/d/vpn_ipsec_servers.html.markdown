---
subcategory: "VPN"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_ipsec_servers"
sidebar_current: "docs-alicloud-datasource-vpn-ipsec-servers"
description: |-
  Provides a list of Vpn Ipsec Servers to the user.
---

# alicloud\_vpn\_ipsec\_servers

This data source provides the Vpn Ipsec Servers of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.161.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpn_ipsec_servers" "ids" {
  ids = ["example_id"]
}
output "vpn_ipsec_server_id_1" {
  value = data.alicloud_vpn_ipsec_servers.ids.servers.0.id
}

data "alicloud_vpn_ipsec_servers" "nameRegex" {
  name_regex = "^my-IpsecServer"
}
output "vpn_ipsec_server_id_2" {
  value = data.alicloud_vpn_ipsec_servers.nameRegex.servers.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ipsec Server IDs.
* `ipsec_server_name` - (Optional, ForceNew) The name of the IPsec server.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Ipsec Server name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `vpn_gateway_id` - (Optional, ForceNew) The ID of the VPN gateway.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ipsec Server names.
* `servers` - A list of Vpn Ipsec Servers. Each element contains the following attributes:
	* `client_ip_pool` - The CIDR block of the client, which is assigned an access address to the virtual NIC of the client.
	* `create_time` - The creation time of the IPsec server. T represents the delimiter, and Z represents UTC, which is World Standard Time.
	* `effect_immediately` - Indicates whether the current IPsec tunnel is deleted and negotiations are reinitiated.
	* `id` - The ID of the Ipsec Server.
	* `idaas_instance_id` - The ID of the Identity as a Service (IDaaS) instance.
	* `ike_config` - The configurations of Phase 1 negotiations.
		* `ike_auth_alg` - The IKE authentication algorithm.
		* `ike_enc_alg` - The IKE encryption algorithm.
		* `ike_lifetime` - The IKE lifetime. Unit: seconds.
		* `ike_mode` - The IKE negotiation mode.
		* `ike_pfs` - Diffie-Hellman key exchange algorithm.
		* `ike_version` - The IKE version.
		* `local_id` - IPsec server identifier. Supports the format of FQDN and IP address. The public IP address of the VPN gateway is selected by default.
		* `remote_id` - The peer identifier. Supports the format of FQDN and IP address, which is empty by default.
	* `internet_ip` - The public IP address of the VPN gateway.
	* `ipsec_config` - The configuration of Phase 2 negotiations.
		* `ipsec_enc_alg` - IPsec encryption algorithm.
		* `ipsec_lifetime` - IPsec survival time. Unit: seconds.
		* `ipsec_pfs` - Diffie-Hellman key exchange algorithm.
		* `ipsec_auth_alg` - IPsec authentication algorithm.
	* `ipsec_server_id` - The ID of the IPsec server.
	* `ipsec_server_name` - The name of the IPsec server.
	* `local_subnet` - Local network segment: the network segment on The VPC side that needs to be interconnected with the client network segment.
	* `max_connections` - The number of SSL connections of the VPN gateway. SSL-VPN the number of SSL connections shared with the IPsec server. For example, if the number of SSL connections is 5 and you have three SSL clients connected to the SSL-VPN, you can also use two clients to connect to the IPsec server.
	* `multi_factor_auth_enabled` - Whether the two-factor authentication function has been turned on.
	* `online_client_count` - The number of clients that have connected to the IPsec server.
	* `psk` - The pre-shared key.
	* `psk_enabled` - Whether to enable the pre-shared key authentication method. The value is only `true`, which indicates that the pre-shared key authentication method is enabled.
	* `vpn_gateway_id` - The ID of the VPN gateway.