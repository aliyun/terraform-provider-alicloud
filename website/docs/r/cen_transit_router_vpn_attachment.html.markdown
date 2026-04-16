---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpn_attachment"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Vpn Attachment resource.
---

# alicloud_cen_transit_router_vpn_attachment

Provides a Cloud Enterprise Network (CEN) Transit Router Vpn Attachment resource.



For information about Cloud Enterprise Network (CEN) Transit Router Vpn Attachment and how to use it, see [What is Transit Router Vpn Attachment](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createtransitroutervpnattachment).

-> **NOTE:** Available since v1.183.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "example" {
  cen_id                     = alicloud_cen_instance.example.id
  transit_router_description = var.name
  transit_router_name        = var.name
}

resource "alicloud_vpn_customer_gateway" "example" {
  customer_gateway_name = var.name
  ip_address            = "42.104.22.210"
  asn                   = "45014"
  description           = var.name
}

resource "alicloud_vpn_gateway_vpn_attachment" "example" {
  network_type       = "public"
  local_subnet       = "0.0.0.0/0"
  remote_subnet      = "0.0.0.0/0"
  effect_immediately = false
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.example.id
    role                 = "master"
    tunnel_index         = 1
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "des"
      ike_version  = "ikev2"
      ike_mode     = "main"
      ike_lifetime = 86400
      psk          = "tf-examplevpn1"
      ike_pfs      = "group1"
      remote_id    = "examplebob1"
      local_id     = "examplealice1"
    }
    tunnel_ipsec_config {
      ipsec_pfs      = "group5"
      ipsec_enc_alg  = "des"
      ipsec_auth_alg = "md5"
      ipsec_lifetime = 86400
    }
  }
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.example.id
    role                 = "slave"
    tunnel_index         = 2
    enable_dpd           = true
    enable_nat_traversal = true
    tunnel_ike_config {
      ike_auth_alg = "md5"
      ike_enc_alg  = "des"
      ike_version  = "ikev2"
      ike_mode     = "main"
      ike_lifetime = 86400
      psk          = "tf-examplevpn2"
      ike_pfs      = "group1"
      remote_id    = "examplebob2"
      local_id     = "examplealice2"
    }
    tunnel_ipsec_config {
      ipsec_pfs      = "group5"
      ipsec_enc_alg  = "des"
      ipsec_auth_alg = "md5"
      ipsec_lifetime = 86400
    }
  }
  vpn_attachment_name = var.name
}

resource "alicloud_cen_transit_router_cidr" "example" {
  transit_router_id        = alicloud_cen_transit_router.example.transit_router_id
  cidr                     = "192.168.0.0/16"
  transit_router_cidr_name = var.name
  description              = var.name
  publish_cidr_route       = true
}

resource "alicloud_cen_transit_router_vpn_attachment" "example" {
  auto_publish_route_enabled            = false
  transit_router_attachment_description = var.name
  transit_router_vpn_attachment_name    = var.name
  cen_id                                = alicloud_cen_transit_router.example.cen_id
  transit_router_id                     = alicloud_cen_transit_router_cidr.example.transit_router_id
  vpn_id                                = alicloud_vpn_gateway_vpn_attachment.example.id
}
```

## Argument Reference

The following arguments are supported:
* `auto_publish_route_enabled` - (Optional) Specifies whether to allow the transit router to automatically advertise routes to the IPsec-VPN attachment. Valid values:

  - `true` (default): yes
  - `false`: no
* `cen_id` - (Optional, ForceNew, Computed) The ID of the Cloud Enterprise Network (CEN) instance.
* `charge_type` - (Optional, ForceNew, Computed, Available since v1.246.0) The billing method.
Set the value to `POSTPAY`, which is the default value and specifies the pay-as-you-go billing method.
* `order_type` - (Optional, Computed, Available since v1.276.0) The entity that pays the fees of the network instance. Valid values:

  - `PayByCenOwner`: the Alibaba Cloud account that owns the CEN instance.
  - `PayByResourceOwner`: the Alibaba Cloud account that owns the network instance.
* `tags` - (Optional, Map) The tag of the resource
* `transit_router_attachment_description` - (Optional) The new description of the VPN attachment.
The description must be 2 to 256 characters in length. The description must start with a letter but cannot start with `http://` or `https://`.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `transit_router_vpn_attachment_name` - (Optional, Available since v1.274.0) The name of the VPN attachment.
The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (\_), and hyphens (-). It must start with a letter.
* `vpn_id` - (Required, ForceNew) The ID of the IPsec-VPN attachment.
* `vpn_owner_id` - (Optional, ForceNew, Computed, Int) The ID of the Alibaba Cloud account to which the IPsec-VPN connection belongs.

  - If you do not set this parameter, the ID of the current Alibaba Cloud account is used.
  - You must set VpnOwnerId if you want to connect the transit router to an IPsec-VPN connection that belongs to another Alibaba Cloud account.
* `zone` - (Optional, ForceNew, List) The Zone ID in the current region.
System will create resources under the Zone that you specify.
Left blank if associated IPSec connection is in dual-tunnel mode. See [`zone`](#zone) below.

The following arguments will be discarded. Please use new fields as soon as possible:
* `transit_router_attachment_name` - (Deprecated since v1.274.0). Field 'transit_router_attachment_name' has been deprecated from provider version 1.274.0. New field 'transit_router_vpn_attachment_name' instead.

### `zone`

The zone supports the following:
* `zone_id` - (Required, ForceNew) The zone ID of the read-only instance.
You can call the [ListTransitRouterAvailableResource](https://www.alibabacloud.com/help/en/doc-detail/261356.html) operation to query the most recent zone list.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - The creation time of the resource.
* `region_id` - The ID of the region where the transit router is deployed.
* `status` - Status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 17 mins) Used when create the Transit Router Vpn Attachment.
* `delete` - (Defaults to 8 mins) Used when delete the Transit Router Vpn Attachment.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Vpn Attachment.

## Import

Cloud Enterprise Network (CEN) Transit Router Vpn Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_vpn_attachment.example <transit_router_attachment_id>
```