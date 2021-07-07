---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vbr_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_vbr_attachment"
description: |-
Provides a Alicloud CEN transit router VBR attachment resource.
---

# alicloud\_cen_transit_router_vbr_attachment

Provides a CEN transit router VBR attachment resource that associate the VBR with the CEN instance.

-> **NOTE:** Available in 1.125.0+

## Example Usage

Basic Usage

```
# Create a new instance-attachment and use it to attach one child instance to a new CEN
variable "name" {
  default = "tf-testAccCenTransitRouterVbrAttachment"
}

variable "vbr_id" {
  default = "vbr-xxxxxxxxxx"
}

variable "transit_router_attachment_name" {
  default = "sdk_rebot_cen_tr_yaochi"
}

variable "transit_router_attachment_description" {
  default = "sdk_rebot_cen_tr_yaochi"
}

resource "alicloud_cen_instance" "cen" {
  name        = var.name
  description = "terraform01"
}

resource "alicloud_transit_router" "tr" {
  name       = var.name
  cen_id     = alicloud_cen_instance.cen.id
}

resource "alicloud_cen_transit_router_vbr_attachment" "foo" {
  vbr_id                                = var.vbr_id
  cen_id                                = alicloud_cen_instance.cen.id
  transit_router_id                     = alicloud_transit_router.tr.id
  region_id                             = "cn-hangzhou"
  auto_publish_route_enabled            = true
  transit_router_attachment_name        = var.transit_router_attachment_name
  transit_router_attachment_description = var.transit_router_attachment_description
}
```
## Argument Reference

The following arguments are supported:

* `vbr_id` - (Required, ForceNew) The ID of the VBR.
* `cen_id` - (Optional, ForceNew) The ID of the CEN.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `region_id` - (Optional) The region ID of the child instance to attach.
* `auto_publish_route_enabled` - (Optional) Auto publish route enabled.
* `transit_router_attachment_name` - (Optional) The name of the transit router vbr attachment.
* `transit_router_attachment_description` - (Optional) The description of the transit router vbr attachment.

->**NOTE:** Ensure that the vbr is not used in Express Connect.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_attachment_id>`. 
* `status` - The associating status of the network.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_vbr_attachment.example tr-********:tr-attach-********
```
