---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_elasticity_assurances"
sidebar_current: "docs-alicloud-datasource-ecs-elasticity-assurances"
description: |-
  Provides a list of Ecs Elasticity Assurance owned by an Alibaba Cloud account.
---

# alicloud_ecs_elasticity_assurances

This data source provides Ecs Elasticity Assurance available to the user.

-> **NOTE:** Available in 1.196.0+

## Example Usage

```
data "alicloud_ecs_elasticity_assurances" "default" {
  ids = ["${alicloud_ecs_elasticity_assurance.default.id}"]
}

output "alicloud_ecs_elasticity_assurance_example_id" {
  value = data.alicloud_ecs_elasticity_assurances.default.assurances.0.id
}
```

## Argument Reference

The following arguments are supported:
* `private_pool_options_ids` - (ForceNew,Optional) The ID of the elastic protection service.
* `resource_group_id` - (ForceNew,Optional) The ID of the resource group.
* `status` - (ForceNew,Optional) The status of flexible guarantee services. Possible values: `All`, `Preparing`, `Prepared`, `Active`, `Released`.
* `tags` - (ForceNew,Optional) The tag key-value pair information bound by the elastic guarantee service.
* `ids` - (Optional, ForceNew, Computed) A list of Elasticity Assurance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Elasticity Assurance IDs.
* `assurances` - A list of Elasticity Assurance Entries. Each element contains the following attributes:
  * `allocated_resources` - Details of resource allocation.
    * `instance_type` - Instance type.
    * `total_amount` - The total number of instances that need to be reserved within an instance type.
    * `used_amount` - The number of instances that have been used.
    * `zone_id` - The zone ID.
  * `id` - ID of flexible guarantee service.
  * `description` - Description of flexible guarantee service.
  * `elasticity_assurance_id` - The first ID of the resource
  * `end_time` - Flexible guarantee service failure time.
  * `instance_charge_type` - The billing method of the instance. Possible value: PostPaid. Currently, only pay-as-you-go is supported.
  * `private_pool_options_id` - The ID of the elasticity assurance.
  * `private_pool_options_name` - The name of the elasticity assurance.
  * `private_pool_options_match_criteria` - The matching mode of flexible guarantee service. Possible values:-Open: flexible guarantee service for Open mode.-Target: specifies the flexible guarantee service of the mode.
  * `resource_group_id` - The ID of the resource group.
  * `start_time` - Flexible guarantee service effective time.
  * `start_time_type` - Flexible guarantee effective way. Possible values:-Now: Effective immediately.-Later: the specified time takes effect.
  * `status` - The status of flexible guarantee services. Possible values:-Preparing: in preparation.-Prepared: to take effect.-Active: in effect.-Released: Released.
  * `tags` - A mapping of tags to assign to the Capacity Reservation.
  * `total_assurance_times` - The total number of flexible guarantee services.
  * `used_assurance_times` - This parameter is not yet available.
