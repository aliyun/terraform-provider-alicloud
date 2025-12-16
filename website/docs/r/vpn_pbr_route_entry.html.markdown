---
subcategory: "VPN Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpn_pbr_route_entry"
sidebar_current: "docs-alicloud-resource-vpn-pbr-route-entry"
description: |-
  Provides a Alicloud VPN Pbr Route Entry resource.
---

# alicloud_vpn_pbr_route_entry

Provides a VPN Pbr Route Entry resource.

-> **NOTE:** Available since v1.162.0.

For information about VPN Pbr Route Entry and how to use it, see [What is VPN Pbr Route Entry](https://www.alibabacloud.com/help/en/doc-detail/127248.html).


## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpn_pbr_route_entry&exampleId=d70d911d-bdd6-c5c7-fd9f-b5a86c81e5ca0d638fa1&activeTab=example&spm=docs.r.vpn_pbr_route_entry.0.d70d911dbd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_vpn_gateways" "default" {
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
  vpn_gateway_id      = data.alicloud_vpn_gateways.default.ids.0
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


resource "alicloud_vpn_pbr_route_entry" "default" {
  vpn_gateway_id = data.alicloud_vpn_gateways.default.ids.0
  route_source   = "192.168.1.0/24"
  route_dest     = "10.0.0.0/24"
  next_hop       = alicloud_vpn_connection.default.id
  weight         = 0
  publish_vpc    = false
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpn_pbr_route_entry&spm=docs.r.vpn_pbr_route_entry.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `vpn_gateway_id` - (Required, ForceNew) The ID of the vpn gateway.
* `next_hop` - (Required, ForceNew) The next hop of the policy-based route.
* `publish_vpc` - (Required) Whether to issue the destination route to the VPC.
* `route_source` - (Required, ForceNew) The source CIDR block of the policy-based route.
* `route_dest` - (Required, ForceNew) The destination CIDR block of the policy-based route.
* `weight` - (Required) The weight of the policy-based route. Valid values: 0 and 100.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the vpn pbr route entry. The value formats as `<vpn_gateway_id>:<next_hop>:<route_source>:<route_dest>`.
* `status` - The status of the vpn pbr route entry.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the vpn pbr route entry.
* `update` - (Defaults to 5 mins) Used when update the vpn pbr route entry.
* `delete` - (Defaults to 5 mins) Used when delete the vpn pbr route entry.

## Import

VPN Pbr route entry can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpn_pbr_route_entry.example <vpn_gateway_id>:<next_hop>:<route_source>:<route_dest>
```
