---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_peer_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_peer_attachment"
description: |-
Provides a Alicloud CEN transit router peer attachment resource.
---

# alicloud\_cen_transit_router_peer_attachment

Provides a CEN transit router peer attachment resource that associate the transit router with the CEN instance.

-> **NOTE:** Available in 1.125.0+

## Example Usage

Basic Usage

```
# Create a new tansit-router-peer-attachment and use it to attach one peer transit router to a new CEN
provider "alicloud" {
  alias = "other_region_id"
  region = var.peer_transit_router_region_id
}

variable "peer_transit_router_region_id" {
  default = "us-east-1"
}

variable "geographic_region_a_id" {
  default = "China"
}

variable "geographic_region_b_id" {
  default = "North-America"
}

variable "bandwidth" {
  default = 2
}

variable "auto_publish_route_enabled" {
  default = true
}

variable "transit_router_attachment_name" {
  default = "sdk_rebot_cen_tr_yaochi"
}

variable "transit_router_attachment_description" {
  default = "sdk_rebot_cen_tr_yaochi"
}

resource "alicloud_cen_bandwidth_package_prepay" "default" {
  name = "lyf_test_bwp_prepay"
  bandwidth = 2
  geographic_region_a_id = var.geographic_region_a_id
  geographic_region_b_id = var.geographic_region_b_id
}

resource "alicloud_cen_instance" "default" {
  name = "sdk_rebot_cen_tr_yaochi"
  protection_level = "REDUCED"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
  instance_id        = alicloud_cen_instance.default.id
  bandwidth_package_id = alicloud_cen_bandwidth_package_prepay.default.id
  depends_on = [
    alicloud_cen_bandwidth_package_prepay.default,
    alicloud_cen_instance.default]
}

resource "alicloud_cen_transit_router" "default_0" {
  cen_id = alicloud_cen_instance.default.id
  depends_on = [
    alicloud_cen_bandwidth_package_attachment.default]
}

resource "alicloud_cen_transit_router" "default_1" {
  provider = alicloud.other_region_id
  cen_id = alicloud_cen_instance.default.id
  depends_on = [
    alicloud_cen_transit_router.default_0]
}

resource "alicloud_cen_transit_router_peer_attachment" "default" {
  cen_id = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default_0.id
  peer_transit_router_region_id = var.peer_transit_router_region_id
  peer_transit_router_id = alicloud_cen_transit_router.default_1.id
  cen_bandwidth_package_id = alicloud_cen_bandwidth_package_prepay.default.id
  bandwidth = var.bandwidth
  auto_publish_route_enabled = var.auto_publish_route_enabled
  transit_router_attachment_name = var.transit_router_attachment_name
  transit_router_attachment_description = var.transit_router_attachment_description
  depends_on = [
    alicloud_cen_transit_router.default_0,
    alicloud_cen_transit_router.default_1]
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router to attach.
* `transit_router_peer_region_id` - (Optional, ForceNew) The region ID of the peer transit router.
* `peer_transit_router_id` - (Required, ForceNew) The ID of the peer transit router.
* `cen_bandwidth_package_id` - (Optional, ForceNew) The ID of the bandwidth package.
* `bandwidth` - (Optional) The bandwidth of the bandwidth package.
* `auto_publish_route_enabled` - (Optional) Auto publish route enabled.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_attachment_id>`. 
* `status` - The associating status of the network.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_peer_attachment.example tr-********:tr-attach-*******
```
