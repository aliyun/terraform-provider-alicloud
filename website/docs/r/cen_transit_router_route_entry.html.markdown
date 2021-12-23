---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_entry"
sidebar_current: "docs-alicloud-resource-cen-transit_router_route_entry"
description: |-
  Provides a Alicloud CEN transit router route entry resource.
---

# alicloud\_cen_transit_router_route_entry

Provides a CEN transit router route entry resource.[What is Cen Transit Router Route Entry](https://help.aliyun.com/document_detail/261238.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

Basic Usage

```terraform
# Create a new tr-attachment and use it to attach one transit router to a new CEN
variable "name" {
  default = "tf-testAccCenTransitRouter"
}

variable "transit_router_route_entry_destination_cidr_block_attachment" {
  default = "192.168.0.0/24"
}

variable "transit_router_route_entry_name" {
  default = "sdk_rebot_cen_tr_yaochi"
}

variable "transit_router_route_entry_description" {
  default = "sdk_rebot_cen_tr_yaochi"
}

resource "alicloud_cen_instance" "cen" {
  name        = var.name
  description = "terraform01"
}

resource "alicloud_cen_transit_router" "default" {
  name   = var.name
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

resource "alicloud_cen_transit_router_route_entry" "default" {
  transit_router_route_table_id                     = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
  transit_router_route_entry_destination_cidr_block = var.transit_router_route_entry_destination_cidr_block_attachment
  transit_router_route_entry_next_hop_type          = "Attachment"
  transit_router_route_entry_name                   = var.transit_router_route_entry_name
  transit_router_route_entry_description            = var.transit_router_route_entry_description
  transit_router_route_entry_next_hop_id            = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
}
```
## Argument Reference

The following arguments are supported:

* `transit_router_route_table_id` - (Required, ForceNew) The ID of the transit router route table.
* `transit_router_route_entry_destination_cidr_block` - (Required, ForceNew) The CIDR of the transit router route entry.
* `transit_router_route_entry_next_hop_type` - (Required, ForceNew) The Type of the transit router route entry next hop,Valid values `Attachment` and `BlackHole`.
* `transit_router_route_entry_name` - (Optional) The name of the transit router route entry.
* `transit_router_route_entry_description` - (Optional) The description of the transit router route entry.
* `transit_router_route_entry_next_hop_id` - (Required, ForceNew) The ID of the transit router route entry next hop.
* `dry_run` - (Optional) The dry run.

-> **NOTE:** If TransitRouterRouteEntryNextHopType is `Attachment`, TransitRouterRouteEntryNextHopId is required.
             If TransitRouterRouteEntryNextHopType is `BlackHole`, TransitRouterRouteEntryNextHopId cannot be filled.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_route_table_id>:<transit_router_route_entry_id>`.
* `status` - The associating status of the Transit Router.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when creating the cen transit router route entry (until it reaches the initial `Active` status).
* `update` - (Defaults to 6 mins) Used when update the cen transit router route entry.
* `delete` - (Defaults to 6 mins) Used when delete the cen transit router route entry.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_route_entry.default vtb-*********:rte-*******
```
