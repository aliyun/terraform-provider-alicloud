---
subcategory: "Ecs"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_saving_plan"
description: |-
  Provides a Alicloud Ecs Saving Plan resource.
---

# alicloud_ecs_saving_plan

Provides a Ecs Saving Plan resource. Ecs saving plan.

For information about Ecs Saving Plan and how to use it, see [What is Saving Plan](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_ecs_saving_plan" "default" {
  committed_amount = "0.01"
  purchase_method  = "group"
  offering_type    = "total"
  saving_plan_name = var.name
  plan_type        = "ecs"
  period           = 1
}
```

### Deleting `alicloud_ecs_saving_plan` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ecs_saving_plan`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `committed_amount` - (Required, ForceNew) Commitments.
* `instance_family` - (Optional, ForceNew) Single specification cluster purchase.
* `offering_type` - (Required, ForceNew) Prepayment Type.
* `payment_type` - (Optional, ForceNew, Computed) The payment type of the resource.
* `period` - (Optional, ForceNew) Billing Duration.
* `plan_type` - (Required, ForceNew) Savings Plan Type.
* `purchase_method` - (Optional) PurchaseMethod.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Saving Plan.

## Import

Ecs Saving Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_saving_plan.example <id>
```