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

-> **NOTE:** Available since v1.245.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "eu-central-1"
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
}

resource "alicloud_vpn_gateway_vpn_attachment" "defaultvrPzdh" {
  customer_gateway_id = alicloud_vpn_customer_gateway.defaultMeoCIz.id
  vpn_attachment_name = "example-vpn-attachment"
  local_subnet        = "10.0.1.0/24"
  remote_subnet       = "10.0.2.0/24"
}


resource "alicloud_cen_transit_router_vpn_attachment" "default" {
  vpn_owner_id                          = alicloud_cen_transit_router.defaultM8Zo6H.id
  cen_id                                = alicloud_cen_transit_router.defaultM8Zo6H.id
  transit_router_attachment_description = "example-vpn-attachment"
  transit_router_id                     = alicloud_cen_transit_router.defaultM8Zo6H.transit_router_id
  vpn_id                                = alicloud_vpn_gateway_vpn_attachment.defaultvrPzdh.id
  auto_publish_route_enabled            = false
  charge_type                           = "POSTPAY"
  transit_router_attachment_name        = "example-vpn-attachment"
}

data "alicloud_cen_transit_router_vpn_attachments" "default" {
  ids               = ["${alicloud_cen_transit_router_vpn_attachment.default.id}"]
  cen_id            = alicloud_cen_transit_router.defaultM8Zo6H.id
  transit_router_id = alicloud_cen_transit_router.defaultM8Zo6H.transit_router_id
}

output "alicloud_cen_transit_router_vpn_attachment_example_id" {
  value = data.alicloud_cen_transit_router_vpn_attachments.default.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (ForceNew, Optional) The ID of the Cloud Enterprise Network (CEN) instance.
* `status` - (Optional, ForceNew) The Status of Transit Router Vpn Attachment. Valid Value: `Attached`, `Attaching`, `Detaching`.
* `tags` - (ForceNew, Optional) The tag of the resource
* `transit_router_attachment_id` - (ForceNew, Optional) The ID of the VPN attachment.
* `transit_router_id` - (ForceNew, Optional) The ID of the transit router.
* `ids` - (Optional, ForceNew, Computed) A list of Transit Router Vpn Attachment IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Transit Router Vpn Attachment IDs.
* `names` - A list of name of Transit Router Vpn Attachments.
* `attachments` - A list of Transit Router Vpn Attachment Entries. Each element contains the following attributes:
  * `auto_publish_route_enabled` - Specifies whether to allow the transit router to automatically advertise routes to the IPsec-VPN attachment. Valid values:*   **true** (default): yes*   **false**: no
  * `cen_id` - The ID of the Cloud Enterprise Network (CEN) instance.
  * `charge_type` - The billing method.Set the value to **POSTPAY**, which is the default value and specifies the pay-as-you-go billing method.
  * `create_time` - The creation time of the resource
  * `resource_type` - The type of the resource. Set the value to **cen**, which specifies a CEN instance.
  * `status` - Status
  * `tags` - The tag of the resource
  * `transit_router_attachment_description` - The new description of the VPN attachment.The description must be 2 to 256 characters in length. The description must start with a letter but cannot start with `http://` or `https://`.
  * `transit_router_attachment_id` - The ID of the VPN attachment.
  * `transit_router_attachment_name` - The name of the VPN attachment.The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (\_), and hyphens (-). It must start with a letter.
  * `transit_router_id` - The ID of the transit router.
  * `vpn_id` - The ID of the IPsec-VPN attachment.
  * `vpn_owner_id` - The ID of the Alibaba Cloud account to which the IPsec-VPN connection belongs.*   If you do not set this parameter, the ID of the current Alibaba Cloud account is used.*   You must set VpnOwnerId if you want to connect the transit router to an IPsec-VPN connection that belongs to another Alibaba Cloud account.
  * `zone` - The Zone ID in the current region.System will create resources under the Zone that you specify.Left blank if associated IPSec connection is in dual-tunnel mode.
    * `zone_id` - The zone ID of the read-only instance.You can call the [ListTransitRouterAvailableResource](https://www.alibabacloud.com/help/en/doc-detail/261356.html) operation to query the most recent zone list.
  * `id` - The ID of the resource supplied above.
