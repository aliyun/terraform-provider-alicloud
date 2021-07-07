---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_table"
sidebar_current: "docs-alicloud-resource-cen-transit_router_route_table"
description: |-
Provides a Alicloud CEN transit router route table resource.
---

# alicloud\_cen_transit_router_route_table

Provides a CEN transit router route table resource.

-> **NOTE:** Available in 1.125.0+

## Example Usage

Basic Usage

```
variable "name" {
  default = "tf-testAccCenTransitRouter"
}

resource "alicloud_cen_instance" "cen" {
  name        = var.name
  description = "terraform01"
}

resource "alicloud_cen_transit_router" "default" {
  name       = var.name
  cen_id     = alicloud_cen_instance.cen.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `type` - (Optional) The Type of the Transit Router. Valid values: `Enterprise`, `Basic`.
* `transit_router_id` - (Required, ForceNew) The ID of the transit router.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_route_table_id>`.
* `status` - The associating status of the Transit Router.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_route_table.default tr-*********:vtb-********
```
