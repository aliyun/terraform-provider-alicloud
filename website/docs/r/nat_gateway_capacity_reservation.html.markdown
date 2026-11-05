---
subcategory: "NAT Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_nat_gateway_capacity_reservation"
description: |-
  Provides a Alicloud NAT Gateway Capacity Reservation resource.
---

# alicloud_nat_gateway_capacity_reservation

Provides a NAT Gateway Capacity Reservation resource. The resource reserves a minimum amount of capacity units (CU) for one NAT gateway, so sustained peak traffic does not depend on the minute-level delay of automatic elastic scaling.

For information about NAT Gateway Capacity Reservation and how to use it, see [What is NAT Resource Reservation](https://help.aliyun.com/zh/nat-gateway/user-guide/nat-resource-reservation).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** Capacity reservation is a whitelisted feature. If the API is not enabled for your account, ask your account manager to enable it before using this resource.

-> **NOTE:** The minimum reservation is 9000 CU and the default quota is 11250 CU per NAT gateway. Scaling up is unlimited, while scaling down is limited to 2 times per day and a cancellation counts as one scale-down. Create and update therefore wait for the reservation to reach the `provisioned` status, and destroy cancels only the reservation; the parent NAT gateway is not deleted.

-> **NOTE:** Reserved CU is billed hourly and a reservation shorter than one hour is billed as one hour. In a billing cycle the higher of the initial and the effective reserved CU is charged, and the billing change of a scale-down or cancellation takes effect in the next billing cycle even though the reservation status changes earlier in the current cycle. See [billing of NAT Gateway](https://help.aliyun.com/zh/nat-gateway/product-billing/) before reserving capacity.

-> **NOTE:** A NAT gateway holds at most one capacity reservation. The NAT gateway ID is the resource ID, so changing `nat_gateway_id` replaces the resource. Importing a NAT gateway without a reservation returns no resource.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_nat_gateway" "default" {
  vpc_id           = alicloud_vpc.default.id
  nat_gateway_name = var.name
  nat_type         = "Enhanced"
  vswitch_id       = alicloud_vswitch.default.id
  eip_bind_mode    = "MULTI_BINDED"
}

resource "alicloud_nat_gateway_capacity_reservation" "default" {
  nat_gateway_id = alicloud_nat_gateway.default.id
  capacity_units = 9000
}
```

## Argument Reference

The following arguments are supported:
* `nat_gateway_id` - (Required, ForceNew, String, Available since v2.3.0) The ID of the NAT gateway whose capacity is reserved.
* `capacity_units` - (Required, Int, Available since v2.3.0) The number of capacity units (CU) to reserve per hour. The value must be positive, and the product currently requires at least 9000 CU. The upper bound is the `nat_instance_quota_reserved_capacity_units` quota of the account (11250 by default, adjustable to 45000); the API rejects values outside the supported range.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above, which is the ID of the NAT gateway.
* `status` - (Computed, String, Available since v2.3.0) The reservation status. Valid values: `pending`, `provisioned`, `failed`, `deleting` and `deleted`.
* `status_reason` - (Computed, String, Available since v2.3.0) The reason why the reservation failed when `status` is `failed`.
* `decrease_requests_remaining` - (Computed, String, Available since v2.3.0) The remaining number of scale-down and cancellation requests for the current day.
* `last_modified_time` - (Computed, String, Available since v2.3.0) The time when the reservation was last modified.
* `bandwidth` - (Computed, Int, Available since v2.3.0) The bandwidth (Gbps) guaranteed by the reserved CU.
* `cps` - (Computed, Int, Available since v2.3.0) The number of new connections per second guaranteed by the reserved CU.
* `pps` - (Computed, Int, Available since v2.3.0) The number of packets per second guaranteed by the reserved CU.
* `session_count` - (Computed, Int, Available since v2.3.0) The number of concurrent sessions guaranteed by the reserved CU.

## Import

NAT Gateway Capacity Reservation can be imported using the id, e.g.

```shell
$ terraform import alicloud_nat_gateway_capacity_reservation.example <nat_gateway_id>
```
