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

-> **NOTE:** Available since v1.273.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

variable "zone_id" {
  default = "cn-wulanchabu-b"
}

data "alicloud_account" "default" {}

resource "alicloud_eflo_vpd" "transmission_vpd" {
  cidr     = "10.0.0.0/8"
  vpd_name = "tf-transmission-${var.name}"
}

resource "alicloud_eflo_vpd" "reception_vpd" {
  cidr     = "10.1.0.0/8"
  vpd_name = "tf-reception-${var.name}"
}

resource "alicloud_eflo_er" "ER" {
  er_name        = "tf-er-routemap-${var.name}"
  master_zone_id = var.zone_id
}

resource "alicloud_eflo_er_route_map" "default" {
  transmission_instance_type  = "VPD"
  action                      = "permit"
  reception_instance_type     = "VPD"
  description                 = "route-map-vpd-to-vpd"
  reception_instance_id       = alicloud_eflo_vpd.reception_vpd.id
  er_id                       = alicloud_eflo_er.ER.id
  reception_instance_owner    = data.alicloud_account.default.id
  transmission_instance_owner = data.alicloud_account.default.id
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
* `reception_instance_owner` - (Optional, ForceNew) The tenant ID of the receiving instance.
* `reception_instance_type` - (Required, ForceNew) Receive instance type. Valid values: `VPD`, `VCC`.
* `transmission_instance_id` - (Required, ForceNew) Publish instance ID.
* `transmission_instance_owner` - (Optional, ForceNew) The tenant ID to which the publish instance belongs.
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
* `create` - (Defaults to 5 mins) Used when create the Er Route Map.
* `delete` - (Defaults to 5 mins) Used when delete the Er Route Map.
* `update` - (Defaults to 5 mins) Used when update the Er Route Map.

## Import

Eflo Er Route Map can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_er_route_map.example <er_id>:<er_route_map_id>
```
