---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpc_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_vpc_attachment"
description: |-
Provides a Alicloud CEN transit router VPC attachment resource.
---

# alicloud\_cen_transit_router_vpc_attachment

Provides a CEN transit router VPC attachment resource that associate the VPC with the CEN instance.

-> **NOTE:** Available in 1.125.0+

## Example Usage

Basic Usage

```
variable "transit_router_attachment_name" {
  default = "sdk_rebot_cen_tr_yaochi"
}

variable "transit_router_attachment_description" {
  default = "sdk_rebot_cen_tr_yaochi"
}

data "alicloud_cen_transit_router_available_resource" "default" {
}

resource "alicloud_vpc" "default" {
  name = "sdk_rebot_cen_tr_yaochi"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  name = "sdk_rebot_cen_tr_yaochi"
  vpc_id = alicloud_vpc.default.id
  cidr_block = "192.168.1.0/24"
  availability_zone = data.alicloud_cen_transit_router_available_resource.default.zones.0.master_zones.0
}

resource "alicloud_vswitch" "default_slave" {
  name = "sdk_rebot_cen_tr_yaochi"
  vpc_id = alicloud_vpc.default.id
  cidr_block = "192.168.2.0/24"
  availability_zone = data.alicloud_cen_transit_router_available_resource.default.zones.0.slave_zones.0
}

resource "alicloud_cen_instance" "default" {
  name = "sdk_rebot_cen_tr_yaochi"
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.id
  vpc_id = alicloud_vpc.default.id
  zone_mapping {
    zone_id = data.alicloud_cen_transit_router_available_resource.default.zones.0.master_zones.0
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mapping {
    zone_id = data.alicloud_cen_transit_router_available_resource.default.zones.0.slave_zones.0
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_name        = var.transit_router_attachment_name
  transit_router_attachment_description = var.transit_router_attachment_description
}
```
## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `cen_id` - (Optional, ForceNew) The ID of the CEN.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `zone_mapping` - (Required, ForceNew) The list of zone mapping of the VPC.
* `transit_router_attachment_name` - (Optional) The name of the transit router vbr attachment.
* `transit_router_attachment_description` - (Optional) The description of the transit router vbr attachment.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_attachment_id>`.
* `status` - The associating status of the network.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_vpc_attachment.example tr-********:tr-attach-********
```
