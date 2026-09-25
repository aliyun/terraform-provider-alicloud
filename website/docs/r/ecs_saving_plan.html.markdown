---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_saving_plan"
sidebar_current: "docs-alicloud-resource-ecs-saving-plan"
description: |-
  Provides a Alicloud ECS Saving Plan resource.
---

# alicloud_ecs_saving_plan

Provides a ECS Saving Plan resource.

For information about ECS Saving Plan and how to use it, see [What is Saving Plan](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/purchasesavingplanoffering).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_saving_plan" "default" {
  committed_amount = "0.08"
  plan_type        = "EcsCompute"
  offering_type    = "AllUpfront"
  purchase_method  = "ByInstanceFamily"
  period           = 1
  period_unit      = "Year"
  instance_family  = "ecs.g5"
  start_time       = "2026-01-01 00:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `committed_amount` - (Required, ForceNew) The hourly commitment amount of the saving plan, in USD per hour.
* `plan_type` - (Optional, ForceNew) The type of the saving plan. Valid values: `EcsCompute`, `General`. Default to `EcsCompute`.
* `offering_type` - (Optional, ForceNew) The payment method of the saving plan. Valid values: `AllUpfront`, `HalfUpfront`, `NoUpfront`. Default to `AllUpfront`.
* `purchase_method` - (Optional, ForceNew) The purchase method of the saving plan. Valid values: `ByInstanceFamily`, `ByInstanceFamilySet`. Default to `ByInstanceFamily`.
* `charge_type` - (Optional, ForceNew) The charge type. Default to `PrePaid`.
* `period` - (Optional, ForceNew) The period of the saving plan.
* `period_unit` - (Optional, ForceNew) The unit of the period. Valid values: `Year`, `Month`.
* `instance_family` - (Optional, ForceNew) The instance family. Required when `purchase_method` is `ByInstanceFamily`.
* `instance_family_set` - (Optional, ForceNew) The instance family set. Required when `purchase_method` is `ByInstanceFamilySet`.
* `saving_plan_name` - (Optional, ForceNew) The name of the saving plan.
* `description` - (Optional, ForceNew) The description of the saving plan.
* `start_time` - (Optional, ForceNew) The effective start time of the saving plan.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of the Saving Plan, which equals to the Saving Plan ID.
* `saving_plan_id` - The ID of the Saving Plan.
* `create_time` - The creation time of the saving plan.
* `payment_type` - The payment type of the saving plan.
* `region_id` - The region ID of the saving plan.

## Import

ECS Saving Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_saving_plan.example spn-xxxxxxxxxx
```

-> **NOTE:** Saving Plans are financial commitments and cannot be cancelled or deleted via API. When you run `terraform destroy`, the resource will be removed from Terraform state, but the Saving Plan will remain on the cloud. To cancel a Saving Plan, use the console.
