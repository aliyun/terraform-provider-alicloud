---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_table_propagation"
sidebar_current: "docs-alicloud-resource-cen-transit_router_route_table_propagation"
description: |-
  Provides a Alicloud CEN transit router route table propagation resource.
---

# alicloud_cen_transit_router_route_table_propagation

Provides a CEN transit router route table propagation resource.[What is Cen Transit Router Route Table Propagation](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-enabletransitrouterroutetablepropagation)

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_cen_transit_router_available_resources" "default" {}
locals {
  master_zone = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[0]
  slave_zone  = data.alicloud_cen_transit_router_available_resources.default.resources[0].slave_zones[1]
}
resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}
resource "alicloud_vswitch" "example_master" {
  vswitch_name = var.name
  cidr_block   = "192.168.1.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.master_zone
}
resource "alicloud_vswitch" "example_slave" {
  vswitch_name = var.name
  cidr_block   = "192.168.2.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = local.slave_zone
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "example" {
  transit_router_name = var.name
  cen_id              = alicloud_cen_instance.example.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "example" {
  cen_id            = alicloud_cen_instance.example.id
  transit_router_id = alicloud_cen_transit_router.example.transit_router_id
  vpc_id            = alicloud_vpc.example.id
  zone_mappings {
    zone_id    = local.master_zone
    vswitch_id = alicloud_vswitch.example_master.id
  }
  zone_mappings {
    zone_id    = local.slave_zone
    vswitch_id = alicloud_vswitch.example_slave.id
  }
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
}
resource "alicloud_cen_transit_router_route_table" "example" {
  transit_router_id = alicloud_cen_transit_router.example.transit_router_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "example" {
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.example.transit_router_route_table_id
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.example.transit_router_attachment_id
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the cen transit router route table propagation (until it reaches the initial `Attached` status).
* `delete` - (Defaults to 5 mins) Used when delete the cen transit router route table propagation.

## Import

CEN transit router route table propagation can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_route_table_propagation.default tr-********:tr-attach-********
```
