---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_capacity_reservation"
sidebar_current: "docs-alicloud-resource-ecs-capacity-reservation"
description: |-
  Provides a Alicloud Ecs Capacity Reservation resource.
---

# alicloud_ecs_capacity_reservation

Provides a Ecs Capacity Reservation resource.

For information about Ecs Capacity Reservation and how to use it, see [What is Capacity Reservation](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/createcapacityreservation).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_capacity_reservation&exampleId=4a481996-e05c-c906-4e97-aeb527ad94b3ac6990ca&activeTab=example&spm=docs.r.ecs_capacity_reservation.0.4a481996e0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g5"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
  available_instance_type     = data.alicloud_instance_types.default.ids.0
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_ecs_capacity_reservation" "default" {
  description               = "terraform-example"
  platform                  = "linux"
  capacity_reservation_name = "terraform-example"
  end_time_type             = "Unlimited"
  resource_group_id         = data.alicloud_resource_manager_resource_groups.default.ids.0
  instance_amount           = 1
  instance_type             = data.alicloud_instance_types.default.ids.0
  match_criteria            = "Open"
  tags = {
    Created = "terraform-example"
  }
  zone_ids = [data.alicloud_zones.default.zones[0].id]
}
```

## Argument Reference

The following arguments are supported:
* `capacity_reservation_name` - (Optional) Capacity reservation service name.
* `description` - (Optional) description of the capacity reservation instance.
* `end_time` - (Optional) end time of the capacity reservation. the capacity reservation will be  released at the end time automatically if set. otherwise it will last until manually released
* `end_time_type` - (Optional) Release mode of capacity reservation service. Value range:Limited: release at specified time. The EndTime parameter must be specified at the same time.Unlimited: manual release. No time limit.
* `instance_amount` - (Required) The total number of instances that need to be reserved within the capacity reservation.
* `instance_type` - (Required,ForceNew) Instance type. Currently, you can only set the capacity reservation service for one instance type. 
* `match_criteria` - (ForceNew,Optional) The type of private resource pool generated after the capacity reservation service takes effect. Value range:Open: Open mode.Target: dedicated mode.Default value: Open
* `platform` - (Optional) platform of the capacity reservation, value range `windows`, `linux`.
* `resource_group_id` - (ForceNew,Optional) The resource group id.
* `tags` - (Optional) The tag of the resource.
* `dry_run` - (Optional) Specifies whether to pre-check the API request. Valid values: `true` and `false`.
* `zone_ids` - (ForceNew,Required) The ID of the zone in the region to which the capacity reservation service belongs. Currently, it is only supported to create a capacity reservation service in one zone.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `payment_type` - The payment type of the resource
* `start_time_type` - The capacity is scheduled to take effect. Possible values:-Now: Effective immediately.-Later: the specified time takes effect.
* `status` - The status of the capacity reservation.
* `time_slot` - This parameter is under test and is not yet open for use.
* `start_time` - time of the capacity reservation which become active.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Capacity Reservation.
* `delete` - (Defaults to 5 mins) Used when delete the Capacity Reservation.
* `update` - (Defaults to 5 mins) Used when update the Capacity Reservation.

## Import

Ecs Capacity Reservation can be imported using the id, e.g.

```shell
$terraform import alicloud_ecs_capacity_reservation.example <id>
```