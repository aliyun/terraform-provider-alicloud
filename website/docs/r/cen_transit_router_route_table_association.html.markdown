---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_route_table_association"
sidebar_current: "docs-alicloud-resource-cen-transit_router_route_table_association"
description: |-
  Provides a Alicloud Cloud Enterprise Network (CEN) Transit Router Route Table Association resource.
---

# alicloud_cen_transit_router_route_table_association

Provides a Cloud Enterprise Network (CEN) Transit Router Route Table Association resource.

For information about Cloud Enterprise Network (CEN) Transit Router Route Table Association and how to use it, see [What is Transit Router Route Table Association](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-associatetransitrouterattachmentwithroutetable)

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_transit_router_route_table_association&exampleId=1404d232-00a0-8df7-6439-13be0307890b7432a7ba&activeTab=example&spm=docs.r.cen_transit_router_route_table_association.0.1404d23200&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

locals {
  master_zone = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[0]
  slave_zone  = data.alicloud_cen_transit_router_available_resources.default.resources[0].slave_zones[1]
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = var.name
  cidr_block   = "192.168.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = local.master_zone
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = var.name
  cidr_block   = "192.168.2.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = local.slave_zone
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  transit_router_name = var.name
  cen_id              = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id                                = alicloud_cen_instance.default.id
  transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  vpc_id                                = alicloud_vpc.default.id
  transit_router_vpc_attachment_name    = var.name
  transit_router_attachment_description = var.name
  zone_mappings {
    zone_id    = local.master_zone
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mappings {
    zone_id    = local.slave_zone
    vswitch_id = alicloud_vswitch.default_slave.id
  }
}

resource "alicloud_cen_transit_router_route_table_association" "default" {
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
}
```

## Argument Reference

The following arguments are supported:

* `transit_router_route_table_id` - (Required, ForceNew) The ID of the Transit Router Route Table.
* `transit_router_attachment_id` - (Required, ForceNew) The ID the Transit Router Attachment.
* `dry_run` - (Optional, Bool) The dry run.

-> **NOTE:** The Zone of CEN has MasterZone and SlaveZone, first zone_id of zone_mapping need be MasterZone. We have a API to describeZones[API](https://help.aliyun.com/document_detail/261356.html)

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transit Router Route Table Association. It formats as `<transit_router_id>:<transit_router_attachment_id>`.
* `status` - The status of the Transit Router Route Table Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the cen transit router route table association (until it reaches the initial `Attached` status).
* `delete` - (Defaults to 5 mins) Used when delete the cen transit router route table association.

## Import

Cloud Enterprise Network (CEN) Transit Router Route Table Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_route_table_association.example <transit_router_id>:<transit_router_attachment_id>
```
