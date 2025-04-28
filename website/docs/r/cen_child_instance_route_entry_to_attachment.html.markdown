---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_child_instance_route_entry_to_attachment"
sidebar_current: "docs-alicloud-resource-cen-child-instance-route-entry-to-attachment"
description: |-
  Provides a Alicloud Cen Child Instance Route Entry To Attachment resource.
---

# alicloud_cen_child_instance_route_entry_to_attachment

Provides a Cen Child Instance Route Entry To Attachment resource.

For information about Cen Child Instance Route Entry To Attachment and how to use it, see [What is Child Instance Route Entry To Attachment](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createcenchildinstancerouteentrytoattachment).

-> **NOTE:** Available since v1.195.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_child_instance_route_entry_to_attachment&exampleId=459e1393-428d-5d0a-e3d4-8d204eaa53116709cc3b&activeTab=example&spm=docs.r.cen_child_instance_route_entry_to_attachment.0.459e139342&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

resource "alicloud_route_table" "example" {
  vpc_id           = alicloud_vpc.example.id
  route_table_name = var.name
  description      = var.name
}

resource "alicloud_cen_child_instance_route_entry_to_attachment" "example" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.example.transit_router_attachment_id
  cen_id                        = alicloud_cen_instance.example.id
  destination_cidr_block        = "10.0.0.0/24"
  child_instance_route_table_id = alicloud_route_table.example.id
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `child_instance_route_table_id` - (Required, ForceNew) The first ID of the resource
* `destination_cidr_block` - (Required, ForceNew) DestinationCidrBlock
* `transit_router_attachment_id` - (Required, ForceNew) TransitRouterAttachmentId
* `dry_run` - (Optional) Whether to perform pre-check on this request, including permission and instance status verification.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.The value is formulated as `<cen_id>:<child_instance_route_table_id>:<transit_router_attachment_id>:<destination_cidr_block>`.
* `service_type` - ServiceType
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Child Instance Route Entry To Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Child Instance Route Entry To Attachment.

## Import

Cen Child Instance Route Entry To Attachment can be imported using the id, e.g.

```shell
$terraform import alicloud_cen_child_instance_route_entry_to_attachment.example <cen_id>:<child_instance_route_table_id>:<transit_router_attachment_id>:<destination_cidr_block>
```