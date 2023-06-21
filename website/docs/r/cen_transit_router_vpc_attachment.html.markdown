---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpc_attachment"
sidebar_current: "docs-alicloud-resource-cen-transit_router_vpc_attachment"
description: |-
  Provides a Alicloud CEN transit router VPC attachment resource.
---

# alicloud_cen_transit_router_vpc_attachment

Provides a CEN transit router VPC attachment resource that associate the VPC with the CEN instance. [What is Cen Transit Router VPC Attachment](https://www.alibabacloud.com/help/en/cloud-enterprise-network/latest/api-doc-cbn-2017-09-12-api-doc-createtransitroutervpcattachment)

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
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
* `route_table_association_enabled` - (Optional, Deprecated) Whether to enabled route table association. The system default value is `true`. **NOTE:** "Field `route_table_association_enabled` has been deprecated from provider version 1.192.0. Please use the resource `alicloud_cen_transit_router_route_table_association` instead, [how to use alicloud_cen_transit_router_route_table_association](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_route_table_association)."
* `route_table_propagation_enabled` - (Optional, Deprecated) Whether to enabled route table propagation. The system default value is `true`. **NOTE:** "Field `route_table_propagation_enabled` has been deprecated from provider version 1.192.0. Please use the resource `alicloud_cen_transit_router_route_table_propagation` instead, [how to use alicloud_cen_transit_router_route_table_propagation](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cen_transit_router_route_table_propagation)."
* `vpc_owner_id` - (Optional, ForceNew) The owner id of vpc.
* `payment_type` - (Optional, ForceNew, Available in 1.168.0+) The payment type of the resource. Valid values: `PayAsYouGo`.
* `zone_mappings` - (Required) The list of zone mapping of the VPC. **NOTE:** From version 1.184.0, `zone_mappings` can be modified. See [`zone_mappings`](#zone_mappings) below.
-> **NOTE:** The Zone of CEN has MasterZone and SlaveZone, first zone_id of zone_mapping need be MasterZone. We have a API to describeZones[API](https://help.aliyun.com/document_detail/261356.html)
* `auto_publish_route_enabled` - (Optional, Computed, Available in v1.204.0+) Whether the transit router is automatically published to the VPC instance. Default value: `false`. Valid values:
  - `true`: Enable.
  - `false`: Disable.
* `tags` - (Optional, Available in v1.193.1+) A mapping of tags to assign to the resource.

### `zone_mappings`

The `zone_mappings` supports the following:

* `vswitch_id` - (Optional, ForceNew) The VSwitch id of attachment.
* `zone_id` - (Optional, ForceNew) The zone Id of VSwitch.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, It is formatted to `<transit_router_id>:<transit_router_attachment_id>`.
* `status` - The associating status of the network.
* `transit_router_attachment_id` - The ID of transit router attachment. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when creating the cen transit router vpc attachment (until it reaches the initial `Attached` status).
* `update` - (Defaults to 3 mins) Used when update the cen transit router vpc attachment.
* `delete` - (Defaults to 3 mins) Used when delete the cen transit router vpc attachment.

## Import

CEN instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_vpc_attachment.example tr-********:tr-attach-********
```
