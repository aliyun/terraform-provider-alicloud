---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpn_attachments"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-vpn-attachments"
description: |-
  Provides a list of Cen Transit Router Vpn Attachment owned by an Alibaba Cloud account.
---

# alicloud_cen_transit_router_vpn_attachments

This data source provides Cen Transit Router Vpn Attachment available to the user.[What is Transit Router Vpn Attachment](https://next.api.alibabacloud.com/document/Cbn/2017-09-12/CreateTransitRouterVpnAttachment)

-> **NOTE:** Available since v1.183.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "default" {
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_cidr" "default" {
  cidr               = "192.168.10.0/24"
  transit_router_id  = alicloud_cen_transit_router.default.transit_router_id
  publish_cidr_route = true
}

resource "alicloud_vpn_customer_gateway" "default" {
  ip_address            = "1.1.1.8"
  customer_gateway_name = var.name
  depends_on            = [alicloud_cen_transit_router_cidr.default]
}

resource "alicloud_vpn_gateway_vpn_attachment" "default" {
  network_type        = "public"
  local_subnet        = "0.0.0.0/0"
  enable_tunnels_bgp  = "false"
  vpn_attachment_name = var.name
  tunnel_options_specification {
    customer_gateway_id  = alicloud_vpn_customer_gateway.default.id
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

    customer_gateway_id = alicloud_vpn_customer_gateway.default.id
    enable_dpd          = "true"
  }

  remote_subnet = "0.0.0.0/0"
}
resource "alicloud_cen_transit_router_vpn_attachment" "default" {
  auto_publish_route_enabled            = false
  transit_router_attachment_description = var.name
  transit_router_attachment_name        = var.name
  cen_id                                = alicloud_cen_transit_router.default.cen_id
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  vpn_id                                = alicloud_vpn_gateway_vpn_attachment.default.id
  vpn_owner_id                          = data.alicloud_account.default.id
  charge_type                           = "POSTPAY"
  tags = {
    Created = "TF"
    For     = "VpnAttachment"
  }
}

data "alicloud_cen_transit_router_vpn_attachments" "ids" {
  ids               = [alicloud_cen_transit_router_vpn_attachment.default.id]
  cen_id            = alicloud_cen_transit_router_vpn_attachment.default.cen_id
  transit_router_id = alicloud_cen_transit_router_vpn_attachment.default.transit_router_id
}

output "cen_transit_router_vpn_attachment_id_0" {
  value = data.alicloud_cen_transit_router_vpn_attachments.ids.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (Optional, ForceNew) The ID of the Cloud Enterprise Network (CEN) instance.
* `status` - (Optional, ForceNew) The Status of Transit Router Vpn Attachment. Valid values: `Attached`, `Attaching`, `Detaching`.
* `tags` - (Optional, ForceNew, Available since v1.245.0) The tag of the resource.
* `transit_router_attachment_id` - (Optional, ForceNew, Available since v1.245.0) The ID of the VPN attachment.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `ids` - (Optional, ForceNew, List) A list of Transit Router Vpn Attachment IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Transit Router Vpn Attachment name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of name of Transit Router Vpn Attachments.
* `attachments` - A list of Transit Router Vpn Attachment Entries. Each element contains the following attributes:
  * `auto_publish_route_enabled` - Specifies whether to allow the transit router to automatically advertise routes to the IPsec-VPN attachment.
  * `cen_id` - (Available since v1.245.0) The ID of the Cloud Enterprise Network (CEN) instance.
  * `charge_type` - (Available since v1.245.0) The billing method of the VPN attachment.
  * `create_time` - The time when the VPN connection was created.
  * `resource_type` - The type of resource attached to the transit router.
  * `status` - The status of the VPN connection.
  * `tags` - (Available since v1.245.0) The tag of the resource.
  * `transit_router_attachment_description` - The description of the IPsec-VPN connection.
  * `transit_router_attachment_id` - (Available since v1.245.0) The ID of the VPN attachment.
  * `transit_router_attachment_name` - The name of the VPN attachment.
  * `transit_router_id` - The ID of the transit router.
  * `vpn_id` - The ID of the IPsec-VPN attachment.
  * `vpn_owner_id` - The ID of the Alibaba Cloud account to which the IPsec-VPN connection belongs.
  * `zone` - The Zone ID in the current region.System will create resources under the Zone that you specify.Left blank if associated IPSec connection is in dual-tunnel mode.
    * `zone_id` - The zone ID of the read-only instance.
  * `id` - The ID of the resource supplied above.
