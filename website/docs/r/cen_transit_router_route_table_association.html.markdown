---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_table_association"
sidebar_current: "docs-alicloud-resource-cen-transit_router_route_table_association"
description: |-
  Provides a Alicloud CEN transit router route table association resource.
---

# alicloud\_cen_transit_router_route_table_association

Provides a CEN transit router route table association resource.[What is Cen Transit Router Route Table Association](https://help.aliyun.com/document_detail/261242.html)

-> **NOTE:** Available in 1.126.0+

## Example Usage

Basic Usage

```terraform
variable "transit_router_attachment_name" {
  default = "sdk_rebot_cen_tr_yaochi"
}

variable "transit_router_attachment_description" {
  default = "sdk_rebot_cen_tr_yaochi"
}

data "alicloud_cen_transit_router_available_resource" "default" {
}

resource "alicloud_vpc" "default" {
  name       = "sdk_rebot_cen_tr_yaochi"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  name              = "sdk_rebot_cen_tr_yaochi"
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "192.168.1.0/24"
  availability_zone = data.alicloud_cen_transit_router_available_resource.default.zones.0.master_zones.0
}

resource "alicloud_vswitch" "default_slave" {
  name              = "sdk_rebot_cen_tr_yaochi"
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "192.168.2.0/24"
  availability_zone = data.alicloud_cen_transit_router_available_resource.default.zones.0.slave_zones.0
}

resource "alicloud_cen_instance" "default" {
  name             = "sdk_rebot_cen_tr_yaochi"
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id            = alicloud_vpc.default.id
  zone_mapping {
    zone_id    = data.alicloud_cen_transit_router_available_resource.default.zones.0.master_zones.0
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mapping {
    zone_id    = data.alicloud_cen_transit_router_available_resource.default.zones.0.slave_zones.0
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_name        = var.transit_router_attachment_name
  transit_router_attachment_description = var.transit_router_attachment_description
}

resource "alicloud_cen_transit_router_route_table_association" "default" {
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
}
```
## Argument Reference

The following arguments are supported:

* `transit_router_route_table_id` - (Required, ForceNew) The ID of the transit router route table.
* `transit_router_attachment_id` - (Required, ForceNew) The ID the transit router attachment.
* `dry_run` - (Optional) The dry run.

-> **NOTE:** The Zone of CEN has MasterZone and SlaveZone, first zone_id of zone_mapping need be MasterZone. We have a API to describeZones[API](https://help.aliyun.com/document_detail/261356.html)

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_attachment_id>`.
* `status` - The associating status of the network.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the cen transit router route table association (until it reaches the initial `Attached` status).
* `delete` - (Defaults to 5 mins) Used when delete the cen transit router route table association.

## Import

CEN transit router route table association can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_route_table_association.default tr-********:tr-attach-********
```
