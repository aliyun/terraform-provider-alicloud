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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_vpn_attachment&exampleId=b6418f7d-6e4e-df31-e943-8889fe7bb9cf170a6aca&activeTab=example&spm=docs.r.cen_transit_router_vpn_attachment.0.b6418f7d6e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_cen_transit_router_available_resources" "default" {
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
  customer_gateway_id = alicloud_vpn_customer_gateway.example.id
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
    psk          = "tf-examplevpn2"
    ike_pfs      = "group1"
    remote_id    = "examplebob2"
    local_id     = "examplealice2"
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
  transit_router_attachment_name        = var.name
  cen_id                                = alicloud_cen_transit_router.example.cen_id
  transit_router_id                     = alicloud_cen_transit_router_cidr.example.transit_router_id
  vpn_id                                = alicloud_vpn_gateway_vpn_attachment.example.id
  zone {
    zone_id = data.alicloud_cen_transit_router_available_resources.default.resources.0.master_zones.0
  }
}
```

Dual Tunnel Mode Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_vpn_attachment&exampleId=99a96233-f1f8-eba6-1499-0b7d42f900808721ccad&activeTab=example&spm=docs.r.cen_transit_router_vpn_attachment.1.99a96233f1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_account" "default" {
}

resource "alicloud_cen_instance" "defaultbpR5Uk" {
  cen_instance_name = "example-vpn-attachment"
}

resource "alicloud_cen_transit_router" "defaultM8Zo6H" {
  cen_id = alicloud_cen_instance.defaultbpR5Uk.id
}

resource "alicloud_cen_transit_router_cidr" "defaultuUtyCv" {
  cidr              = "192.168.10.0/24"
  transit_router_id = alicloud_cen_transit_router.defaultM8Zo6H.transit_router_id
}

resource "alicloud_vpn_customer_gateway" "defaultMeoCIz" {
  ip_address            = "0.0.0.0"
  customer_gateway_name = "example-vpn-attachment"
  depends_on            = ["alicloud_cen_transit_router_cidr.defaultuUtyCv"]
}

data "alicloud_cen_transit_router_service" "default" {
  enable = "On"
}

resource "alicloud_vpn_gateway_vpn_attachment" "defaultvrPzdh" {
  network_type        = "public"
  local_subnet        = "0.0.0.0/0"
  enable_tunnels_bgp  = "false"
  vpn_attachment_name = var.name
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.defaultMeoCIz.id
    enable_dpd           = "true"
    enable_nat_traversal = "true"
    tunnel_index         = "1"

    tunnel_ike_config {
      remote_id    = "2.2.2.2"
      ike_enc_alg  = "aes"
      ike_mode     = "main"
      ike_version  = "ikev1"
      local_id     = "1.1.1.1"
      ike_auth_alg = "md5"
      ike_lifetime = "86100"
      ike_pfs      = "group2"
      psk          = "12345678"
    }

    tunnel_ipsec_config {
      ipsec_auth_alg = "md5"
      ipsec_enc_alg  = "aes"
      ipsec_lifetime = "86200"
      ipsec_pfs      = "group5"
    }

  }
  tunnel_options_specification {
    enable_nat_traversal = "true"
    tunnel_index         = "2"
    tunnel_ike_config {
      local_id     = "4.4.4.4"
      remote_id    = "5.5.5.5"
      ike_lifetime = "86400"
      ike_pfs      = "group5"
      ike_mode     = "main"
      ike_version  = "ikev2"
      psk          = "32333442"
      ike_auth_alg = "md5"
      ike_enc_alg  = "aes"
    }

    tunnel_ipsec_config {
      ipsec_enc_alg  = "aes"
      ipsec_lifetime = "86400"
      ipsec_pfs      = "group5"
      ipsec_auth_alg = "sha256"
    }

    customer_gateway_id = alicloud_vpn_customer_gateway.defaultMeoCIz.id
    enable_dpd          = "true"
  }

  remote_subnet = "0.0.0.0/0"
}

resource "alicloud_cen_transit_router_vpn_attachment" "default" {
  transit_router_id                     = alicloud_cen_transit_router.defaultM8Zo6H.transit_router_id
  vpn_id                                = alicloud_vpn_gateway_vpn_attachment.defaultvrPzdh.id
  auto_publish_route_enabled            = "false"
  charge_type                           = "POSTPAY"
  transit_router_attachment_name        = "example-vpn-attachment"
  vpn_owner_id                          = data.alicloud_account.default.id
  cen_id                                = alicloud_cen_transit_router.defaultM8Zo6H.cen_id
  transit_router_attachment_description = "example-vpn-attachment"
}
```

## Argument Reference

The following arguments are supported:
* `auto_publish_route_enabled` - (Optional) Specifies whether to allow the transit router to automatically advertise routes to the IPsec-VPN attachment. Valid values:

  - `true` (default): yes
  - `false`: no
* `cen_id` - (Optional, ForceNew) The ID of the Cloud Enterprise Network (CEN) instance.
* `charge_type` - (Optional, ForceNew, Available since v1.246.0) The billing method.
Set the value to `POSTPAY`, which is the default value and specifies the pay-as-you-go billing method.
* `tags` - (Optional, Map) The tag of the resource
* `transit_router_attachment_description` - (Optional) The new description of the VPN attachment.
The description must be 2 to 256 characters in length. The description must start with a letter but cannot start with `http://` or `https://`.
* `transit_router_attachment_name` - (Optional) The name of the VPN attachment.
The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (\_), and hyphens (-). It must start with a letter.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `vpn_id` - (Required, ForceNew) The ID of the IPsec-VPN attachment.
* `vpn_owner_id` - (Optional, ForceNew, Int) The ID of the Alibaba Cloud account to which the IPsec-VPN connection belongs.

  - If you do not set this parameter, the ID of the current Alibaba Cloud account is used.
  - You must set VpnOwnerId if you want to connect the transit router to an IPsec-VPN connection that belongs to another Alibaba Cloud account.
* `zone` - (Optional, ForceNew, Set) The Zone ID in the current region.
System will create resources under the Zone that you specify.
Left blank if associated IPSec connection is in dual-tunnel mode. See [`zone`](#zone) below.

### `zone`

The zone supports the following:
* `zone_id` - (Required, ForceNew) The zone ID of the read-only instance.
You can call the [ListTransitRouterAvailableResource](https://www.alibabacloud.com/help/en/doc-detail/261356.html) operation to query the most recent zone list.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `region_id` - The ID of the region where the transit router is deployed.
* `status` - Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 8 mins) Used when create the Transit Router Vpn Attachment.
* `delete` - (Defaults to 8 mins) Used when delete the Transit Router Vpn Attachment.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Vpn Attachment.

## Import

Cloud Enterprise Network (CEN) Transit Router Vpn Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_vpn_attachment.example <id>
```