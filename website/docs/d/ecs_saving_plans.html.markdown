---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_saving_plans"
sidebar_current: "docs-alicloud-datasource-ecs-saving-plans"
description: |-
  Provides a list of ECS Saving Plans owned by the current user.
---

# alicloud_ecs_saving_plans

This data source provides a list of ECS Saving Plans available to the current user.

-> **NOTE:** Available since v1.244.0.

## Example Usage

```terraform
data "alicloud_ecs_saving_plans" "default" {
  ids         = ["spn-xxxxxxxxxx"]
  output_file = "saving_plans.txt"
}

output "first_plan_id" {
  value = data.alicloud_ecs_saving_plans.default.plans.0.saving_plan_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional) The ID of the saving plan instance. If specified, only the matching plan will be returned.
* `ids` - (Optional) A list of saving plan IDs used to filter the results.
* `name_regex` - (Optional) A regex string to filter saving plans by name.
* `output_file` - (Optional) File name where to save the result after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of saving plan IDs.
* `names` - A list of saving plan names.
* `plans` - A list of saving plans. Each element contains the following attributes:
  * `id` - The ID of the saving plan, same as `saving_plan_id`.
  * `saving_plan_id` - The ID of the saving plan.
  * `committed_amount` - The hourly commitment amount of the saving plan.
  * `create_time` - The creation time of the saving plan.
  * `instance_family` - The instance family of the saving plan.
  * `offering_type` - The payment method of the saving plan.
  * `payment_type` - The payment type of the saving plan.
  * `period` - The period of the saving plan.
  * `plan_type` - The type of the saving plan.
  * `region_id` - The region ID of the saving plan.
  * `start_time` - The effective start time of the saving plan.
