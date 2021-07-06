---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_entry"
sidebar_current: "docs-alicloud-resource-cen-transit_router_route_entry"
description: |-
Provides a Alicloud CEN transit router route entry resource.
---

# alicloud\_cen_transit_router_route_entry

Provides a CEN transit router route entry resource.

-> **NOTE:** Available in 1.125.0+

## Example Usage

Basic Usage

```
# Create a new tr-attachment and use it to attach one transit router to a new CEN
variable "name" {
  default = "tf-testAccCenTransitRouter"
}

variable "region" {
  default = "cn-hangzhou"
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
  name       = var.name
  cen_id     = alicloud_cen_instance.cen.id
  region_id  = var.region
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.id
}

resource "alicloud_cen_transit_router_route_entry" "default" {
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.id
  transit_router_route_entry_destination_cidr_block = var.transit_router_route_entry_destination_cidr_block_attachment
  transit_router_route_entry_next_hop_type = "Attachment"
  transit_router_route_entry_name = var.transit_router_route_entry_name
  transit_router_route_entry_description = var.transit_router_route_entry_description
  transit_router_route_entry_next_hop_id = alicloud_cen_transit_router_vpc_attachment.default.id
  delete_parms = "route_entry_id"
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `region_id` - (Required, ForceNew) The Region ID of the Transit Router.
* `type` - (Optional) The Type of the Transit Router. Valid values: `Enterprise`, `Basic`.
* `transit_router_route_table_id` - (Required, ForceNew) The ID of the transit router route table.



## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_route_table_id>:<transit_router_route_entry_id>`.
* `status` - The associating status of the Transit Router.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_route_entry.default vtb-*********:rte-*******
```
