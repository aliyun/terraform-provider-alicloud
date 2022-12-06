---
subcategory: "Ecs"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_capacity_reservations"
sidebar_current: "docs-alicloud-datasource-ecs-capacity-reservations"
description: |-
  Provides a list of Ecs Capacity Reservation owned by an Alibaba Cloud account.
---

# alicloud_ecs_capacity_reservations

This data source provides Ecs Capacity Reservation available to the user.

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_ecs_capacity_reservations" "default" {
  ids           = ["${alicloud_ecs_capacity_reservation.default.id}"]
  name_regex    = alicloud_ecs_capacity_reservation.default.name
  instance_type = "ecs.c6.large"
  platform      = "linux"
}

output "alicloud_ecs_capacity_reservation_example_id" {
  value = data.alicloud_ecs_capacity_reservations.default.reservations.0.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_type` - (ForceNew,Optional) Instance type. Currently, you can only set the capacity reservation service for one instance type. 
* `payment_type` - (ForceNew,Optional) The payment type of the resource. value range `PostPaid`, `PrePaid`.
* `platform` - (ForceNew,Optional) platform of the capacity reservation , value range `windows`, `linux`, `all`.
* `resource_group_id` - (ForceNew,Optional) The resource group id.
* `status` - (ForceNew,Optional) The status of the capacity reservation. value range `All`, `Pending`, `Preparing`, `Prepared`, `Active`, `Released`.
* `tags` - (ForceNew,Optional) The tag of the resource.
* `ids` - (Optional, ForceNew, Computed) A list of Capacity Reservation IDs.
* `capacity_reservation_names` - (Optional, ForceNew) The name of the Capacity Reservation. You can specify at most 10 names.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Capacity Reservation IDs.
* `names` - A list of name of Capacity Reservations.
* `reservations` - A list of Capacity Reservation Entries. Each element contains the following attributes:
  * `capacity_reservation_id` - Capacity Reservation id
  * `capacity_reservation_name` - Capacity reservation service name.
  * `description` - description of the capacity reservation instance
  * `end_time` - end time of the capacity reservation. the capacity reservation will be  released at the end time automatically if set. otherwise it will last until manually released
  * `end_time_type` - Release mode of capacity reservation service. Value range:Limited: release at specified time. The EndTime parameter must be specified at the same time.Unlimited: manual release. No time limit.
  * `instance_amount` - The total number of instances that need to be reserved within the capacity reservation
  * `instance_type` - Instance type. Currently, you can only set the capacity reservation service for one instance type.
  * `match_criteria` - The type of private resource pool generated after the capacity reservation service takes effect. Value range:Open: Open mode.Target: dedicated mode.Default value: Open
  * `payment_type` - The payment type of the resource
  * `platform` - platform of the capacity reservation.
  * `resource_group_id` - The resource group id
  * `start_time` -  time of the capacity reservation which become active
  * `start_time_type` - The capacity is scheduled to take effect. Possible values:-Now: Effective immediately.-Later: the specified time takes effect.
  * `status` - The status of the capacity reservation.
  * `tags` - A mapping of tags to assign to the Capacity Reservation.
  * `time_slot` - This parameter is under test and is not yet open for use.
  * `zone_ids` - The ID of the zone in the region to which the capacity reservation service belongs. Currently, it is only supported to create a capacity reservation service in one zone.
  * `id` - The ID of the Capacity Reservation.
