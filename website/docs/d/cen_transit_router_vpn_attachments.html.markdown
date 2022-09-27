---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpn_attachments"
sidebar_current: "docs-alicloud-datasource-cen-transit-router-vpn-attachments"
description: |-
  Provides a list of Cen Transit Router Vpn Attachments to the user.
---

# alicloud\_cen\_transit\_router\_vpn\_attachments

This data source provides the Cen Transit Router Vpn Attachments of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.183.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cen_transit_router_vpn_attachments" "ids" {
  cen_id = "example_value"
}
output "cen_transit_router_vpn_attachment_id_1" {
  value = data.alicloud_cen_transit_router_vpn_attachments.ids.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The id of the cen.
* `ids` - (Optional, ForceNew, Computed) A list of Transit Router Vpn Attachment IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The Status of Transit Router Vpn Attachment. Valid Value: `Attached`, `Attaching`, `Detaching`.
* `transit_router_id` - (Optional, ForceNew) The ID of the forwarding router instance.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `attachments` - A list of Cen Transit Router Vpn Attachments. Each element contains the following attributes:
  * `auto_publish_route_enabled` - Whether to allow the forwarding router instance to automatically publish routing entries to IPsec connections.
  * `create_time` - The creation time of the resource.
  * `resource_type` - Type of the resource.
  * `status` - The status of the transit router attachment.
  * `transit_router_attachment_description` - The description of the VPN connection.
  * `transit_router_attachment_name` - The name of the VPN connection.
  * `vpn_id` - The id of the vpn.
  * `vpn_owner_id` - The owner id of vpn.
  * `zone` - The list of zone mapping.
    * `zone_id` - The id of the zone.