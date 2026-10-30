---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_er_route_map"
description: |-
  Provides a Alicloud Eflo Er Route Map resource.
---

# alicloud_eflo_er_route_map

Provides a Eflo Er Route Map resource.

Lingjun HUB routing strategy.

For information about Eflo Er Route Map and how to use it, see [What is Er Route Map](https://next.api.alibabacloud.com/document/eflo/2022-05-30/CreateErRouteMap).

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** The transmission instance and the reception instance must be attached to the Lingjun HUB (`alicloud_eflo_er`) through `alicloud_eflo_er_attachment` before the routing strategy is created.

-> **NOTE:** When `auto_receive_all_route` of an `alicloud_eflo_er_attachment` is set to `true`, the Lingjun HUB automatically creates a routing strategy with the `3000` policy number for `0.0.0.0/0` in both directions. Creating a policy with the same publish instance, receive instance, destination CIDR block and action fails with the `1002` error, so keep `auto_receive_all_route` as `false` when you manage routing strategies manually.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

variable "zone_id" {
  default = "cn-wulanchabu-b"
}

resource "alicloud_eflo_vpd" "transmission_vpd" {
  cidr     = "10.0.0.0/8"
  vpd_name = "tf-transmission-${var.name}"
}

resource "alicloud_eflo_vpd" "reception_vpd" {
  cidr     = "192.168.0.0/16"
  vpd_name = "tf-reception-${var.name}"
}

resource "alicloud_eflo_er" "ER" {
  er_name        = "tf-er-routemap-${var.name}"
  master_zone_id = var.zone_id
}

resource "alicloud_eflo_er_attachment" "transmission_attachment" {
  er_id                  = alicloud_eflo_er.ER.id
  er_attachment_name     = "tf-transmission-attachment-${var.name}"
  instance_id            = alicloud_eflo_vpd.transmission_vpd.id
  instance_type          = "VPD"
  auto_receive_all_route = false
}

resource "alicloud_eflo_er_attachment" "reception_attachment" {
  er_id                  = alicloud_eflo_er.ER.id
  er_attachment_name     = "tf-reception-attachment-${var.name}"
  instance_id            = alicloud_eflo_vpd.reception_vpd.id
  instance_type          = "VPD"
  auto_receive_all_route = false
}

resource "alicloud_eflo_er_route_map" "default" {
  transmission_instance_type  = "VPD"
  action                      = "permit"
  reception_instance_type     = "VPD"
  description                 = "route-map-vpd-to-vpd"
  reception_instance_id       = alicloud_eflo_vpd.reception_vpd.id
  er_id                       = alicloud_eflo_er.ER.id
  reception_instance_owner    = alicloud_eflo_er_attachment.reception_attachment.tenant_id
  transmission_instance_owner = alicloud_eflo_er_attachment.transmission_attachment.tenant_id
  transmission_instance_id    = alicloud_eflo_vpd.transmission_vpd.id
  er_route_map_num            = 1001
  destination_cidr_block      = "0.0.0.0/0"
}
```

## Argument Reference

The following arguments are supported:
* `action` - (Required, ForceNew) Strategic behavior. Valid values: `permit`, `deny`.
* `description` - (Optional) Lingjun HUB routing policy description information.
* `destination_cidr_block` - (Required, ForceNew) Destination network segment.
* `er_id` - (Required, ForceNew) Lingjun HUB ID.
* `er_route_map_num` - (Required, ForceNew, Int) Policy number. Valid values: 1001-2000.
* `reception_instance_id` - (Required, ForceNew) Receive instance ID.
* `reception_instance_owner` - (Optional, Computed, ForceNew) The tenant ID of the receiving instance.
* `reception_instance_type` - (Required, ForceNew) Receive instance type. Valid values: `VPD`, `VCC`.
* `transmission_instance_id` - (Required, ForceNew) Publish instance ID.
* `transmission_instance_owner` - (Optional, Computed, ForceNew) The tenant ID to which the publish instance belongs.
* `transmission_instance_type` - (Required, ForceNew) Publish instance type. Valid values: `VPD`, `VCC`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<er_id>:<er_route_map_id>`.
* `create_time` - Creation time.
* `er_route_map_id` - Routing Policy ID.
* `region_id` - The region ID of the resource.
* `status` - Status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 15 mins) Used when create the Er Route Map.
* `delete` - (Defaults to 20 mins) Used when delete the Er Route Map.
* `update` - (Defaults to 15 mins) Used when update the Er Route Map.

## Import

Eflo Er Route Map can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_er_route_map.example <er_id>:<er_route_map_id>
```
