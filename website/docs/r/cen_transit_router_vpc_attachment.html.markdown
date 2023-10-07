---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_transit_router_vpc_attachment"
description: |-
  Provides a Alicloud CEN Transit Router Vpc Attachment resource.
---

# alicloud_cen_transit_router_vpc_attachment

Provides a CEN Transit Router Vpc Attachment resource. 

For information about CEN Transit Router Vpc Attachment and how to use it, see [What is Transit Router Vpc Attachment](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_cen_instance" "defaultJ6HrUE" {
  cen_instance_name = var.name

}

resource "alicloud_cen_transit_router" "defaults5WvfD" {
  cen_id = alicloud_cen_instance.defaultJ6HrUE.cen_id
}

resource "alicloud_vpc" "defaultJLRlxW" {
  vpc_name = var.name

  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaulteKv3Dd" {
  vpc_id     = alicloud_vpc.defaultJLRlxW.id
  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "10.10.0.0/24"
}

resource "alicloud_vswitch" "defaulteKv3Da" {
  vpc_id     = alicloud_vpc.defaultJLRlxW.id
  zone_id    = data.alicloud_zones.default.zones.1.id
  cidr_block = "10.20.0.0/24"
}

resource "alicloud_vswitch" "defaulteKv3Dc" {
  vpc_id     = alicloud_vpc.defaultJLRlxW.id
  zone_id    = data.alicloud_zones.default.zones.2.id
  cidr_block = "10.30.0.0/24"
}

resource "alicloud_vswitch" "defaulteKv3Db" {
  vpc_id     = alicloud_vpc.defaultJLRlxW.id
  zone_id    = data.alicloud_zones.default.zones.2.id
  cidr_block = "10.40.0.0/24"
}


resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  vpc_id = alicloud_vpc.defaultJLRlxW.id
  cen_id = alicloud_cen_instance.defaultJ6HrUE.cen_id
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaulteKv3Da.id
    zone_id    = "cn-hangzhou-j"
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaulteKv3Dd.id
    zone_id    = "cn-hangzhou-h"
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.defaulteKv3Dc.id
    zone_id    = "cn-hangzhou-k"
  }
  transit_router_id                  = alicloud_cen_transit_router.defaults5WvfD.id
  transit_router_vpc_attachment_name = var.name

  transit_router_attachment_description = "test"
  charge_type                           = "POSTPAY"
  auto_publish_route_enabled            = true
}
```

## Argument Reference

The following arguments are supported:
* `auto_publish_route_enabled` - (Optional, Available since v1.126.0) Switch to turn automatic route publish on or off.
* `cen_id` - (Optional, Available since v1.126.0) CenId.
* `charge_type` - (Optional, ForceNew) ChargeType.
* `resource_type` - (Optional, Available since v1.126.0) ResourceType.
* `tags` - (Optional, Map, Available since v1.126.0) The tag of the resource.
* `transit_router_attachment_description` - (Optional, Available since v1.126.0) TransitRouterAttachmentDescription.
* `transit_router_id` - (Optional, ForceNew, Available since v1.126.0) TransitRouterId.
* `transit_router_vpc_attachment_name` - (Optional) TransitRouterAttachmentName.
* `vpc_id` - (Required, ForceNew, Available since v1.126.0) VpcId.
* `zone_mappings` - (Required, Available since v1.126.0) ZoneMappingss. See [`zone_mappings`](#zone_mappings) below.

The following arguments will be discarded. Please use new fields as soon as possible:
* `transit_router_attachment_name` - (Deprecated since v1.212.0). Field 'transit_router_attachment_name' has been deprecated from provider version 1.212.0. New field 'transit_router_vpc_attachment_name' instead.

### `zone_mappings`

The zone_mappings supports the following:
* `vswitch_id` - (Required) VSwitchId.
* `zone_id` - (Required) ZoneId.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - Status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Transit Router Vpc Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Transit Router Vpc Attachment.
* `update` - (Defaults to 5 mins) Used when update the Transit Router Vpc Attachment.

## Import

CEN Transit Router Vpc Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_transit_router_vpc_attachment.example <id>
```