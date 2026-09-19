---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_route"
description: |-
  Provides a Alicloud Data Works Route resource.
---

# alicloud_data_works_route

Provides a Data Works Route resource.

Resource group network routing rules.

For information about Data Works Route and how to use it, see [What is Route](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createroute).

-> **NOTE:** Available since v1.288.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "default5Bia4h" {
  description = var.name
  vpc_name    = var.name
  cidr_block  = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultss7s7F" {
  description  = var.name
  vpc_id       = alicloud_vpc.default5Bia4h.id
  zone_id      = "cn-beijing-g"
  vswitch_name = format("%s1", var.name)
  cidr_block   = "10.0.0.0/24"
}

resource "alicloud_data_works_dw_resource_group" "defaultVJvKvl" {
  payment_duration_unit = "Month"
  payment_type          = "PostPaid"
  specification         = "500"
  default_vswitch_id    = alicloud_vswitch.defaultss7s7F.id
  remark                = var.name
  resource_group_name   = "route_openapi_example01"
  default_vpc_id        = alicloud_vpc.default5Bia4h.id
}

resource "alicloud_vpc" "defaulte4zhaL" {
  description = var.name
  vpc_name    = format("%s3", var.name)
  cidr_block  = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default675v38" {
  description  = var.name
  vpc_id       = alicloud_vpc.defaulte4zhaL.id
  zone_id      = "cn-beijing-g"
  vswitch_name = format("%s4", var.name)
  cidr_block   = "172.16.0.0/24"
}


resource "alicloud_data_works_network" "default" {
  vpc_id               = alicloud_vpc.defaulte4zhaL.id
  vswitch_id           = alicloud_vswitch.default675v38.id
  dw_resource_group_id = alicloud_data_works_dw_resource_group.defaultVJvKvl.id
}

resource "alicloud_data_works_route" "default" {
  network_id       = alicloud_data_works_network.default.id
  destination_cidr = "192.168.10.0/24"
}
```

## Argument Reference

The following arguments are supported:
* `destination_cidr` - (Required) The CIDR block of the destination route.
* `network_id` - (Required, ForceNew) The ID of the network resource to which the route belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `route_id` - The ID of the route.
* `create_time` - The time when the route was created.
* `dw_resource_group_id` - The ID of the resource group to which the route belongs.
* `region_id` - The region to which the route belongs.
* `resource_id` - The identifier of the network resource to which the route belongs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Route.
* `update` - (Defaults to 5 mins) Used when update the Route.
* `delete` - (Defaults to 5 mins) Used when delete the Route.

## Import

Data Works Route can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_route.example <id>
```
