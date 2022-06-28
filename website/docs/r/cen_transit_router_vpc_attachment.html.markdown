---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpc_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_vpc_attachment"
description: |-
  Provides a Alicloud CEN transit router VPC attachment resource.
---

# alicloud\_cen_transit_router_vpc_attachment

Provides a CEN transit router VPC attachment resource that associate the VPC with the CEN instance. [What is Cen Transit Router VPC Attachment](https://help.aliyun.com/document_detail/261358.html)

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

data "alicloud_cen_transit_router_available_resources" "default" {

}

resource "alicloud_vpc" "default" {
  vpc_name   = "sdk_rebot_cen_tr_yaochi"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = "sdk_rebot_cen_tr_yaochi"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[0]
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = "sdk_rebot_cen_tr_yaochi"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].slave_zones[0]
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "sdk_rebot_cen_tr_yaochi"
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.id
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
```
## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) The dry run.
* `cen_id` - (Required, ForceNew) The ID of the CEN.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `transit_router_id` - (Optional, ForceNew) The ID of the transit router.
* `transit_router_attachment_name` - (Optional) The name of the transit router vbr attachment.
* `transit_router_attachment_description` - (Optional) The description of the transit router vbr attachment.
* `resource_type` - (Optional) The resource type of transit router vpc attachment. Valid value `VPC`. Default value is `VPC`.
* `route_table_association_enabled` - (Optional) Whether to enabled route table association. The system default value is `true`.
* `route_table_propagation_enabled` - (Optional) Whether to enabled route table propagation. The system default value is `true`.
* `vpc_owner_id` - (Optional,ForceNew) The owner id of vpc.
* `payment_type` - (Optional, ForceNew, Available in 1.168.0+) The payment type of the resource. Valid values: `PayAsYouGo`.
* `zone_mappings` - (Required, ForceNew) The list of zone mapping of the VPC.

-> **NOTE:** The Zone of CEN has MasterZone and SlaveZone, first zone_id of zone_mapping need be MasterZone. We have a API to describeZones[API](https://help.aliyun.com/document_detail/261356.html)

#### ZoneMapping Block

The `zone_mapping` supports the following:

* `vswitch_id` - (Optional, ForceNew) The VSwitch id of attachment.
* `zone_id` - (Optional, ForceNew) The zone Id of VSwitch.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_attachment_id>`.
* `status` - The associating status of the network.
* `transit_router_attachment_id` - The ID of transit router attachment. 

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when creating the cen transit router vpc attachment (until it reaches the initial `Attached` status).
* `update` - (Defaults to 3 mins) Used when update the cen transit router vpc attachment.
* `delete` - (Defaults to 3 mins) Used when delete the cen transit router vpc attachment.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_transit_router_vpc_attachment.example tr-********:tr-attach-********
```
