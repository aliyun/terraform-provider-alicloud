---
subcategory: "Bss Open Api"
layout: "alicloud"
page_title: "Alicloud: alicloud_bss_open_api_budget"
description: |-
  Provides a Alicloud Bss Open Api Budget resource.
---

# alicloud_bss_open_api_budget

Provides a Bss Open Api Budget resource.

For information about Bss Open Api Budget and how to use it, see [What is Budget](https://next.api.alibabacloud.com/document/BssOpenApi/2023-09-30/CreateBudget).

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_bss_open_api_budget" "default" {
  budget_name        = var.name
  budget_type        = "cost"
  metric             = "Cost"
  cycle_type         = "Month"
  cycle_start_period = "2026-09"
  cycle_end_period   = "2026-12"
  quota_type         = "quota"
  quota              = "80"
  comment            = "terraform example budget"
  cycle_quota {
    cycle_period = "2026-09"
    quota        = "80"
  }
  query_filter {
    code        = "productCode"
    select_type = "equal"
    values      = ["ECS"]
  }
  warn_confs {
    name            = "terraform-warn"
    warn_target     = "Msc"
    threshold_type  = "percent"
    threshold_value = "80"
    event_bridge    = false
    comment         = "terraform warn conf"
  }
}
```

## Argument Reference

The following arguments are supported:
* `budget_name` - (Required, ForceNew) The name of the budget.
* `budget_type` - (Required) The type of the budget.
* `cycle_end_period` - (Required) The end period of the budget cycle.
* `cycle_start_period` - (Required) The start period of the budget cycle.
* `cycle_type` - (Required) The cycle type of the budget.
* `metric` - (Required) The metric of the budget.
* `quota_type` - (Required) The quota type of the budget.
* `comment` - (Optional) The comment of the budget.
* `nbid` - (Optional) The first-level marketplace ID. If empty, the current user marketplace ID is used.
* `quota` - (Optional) The fixed quota value. When the quota type is a quota, the unit is a percentage.
* `cycle_quota` - (Optional) The specified quota for each cycle. See [`cycle_quota`](#cycle_quota) below.
* `query_filter` - (Optional) The filter conditions. See [`query_filter`](#query_filter) below.
* `warn_confs` - (Optional) The alert configurations. See [`warn_confs`](#warn_confs) below.

### `cycle_quota`

The cycle_quota supports the following:
* `cycle_period` - (Optional) The cycle period.
* `quota` - (Optional) The quota.

### `query_filter`

The query_filter supports the following:
* `code` - (Optional) The parameter code.
* `select_type` - (Optional) The select mode.
* `values` - (Optional) The list of filter values.

### `warn_confs`

The warn_confs supports the following:
* `comment` - (Optional) The comment.
* `event_bridge` - (Optional) Whether to enable EventBridge.
* `msc_channels` - (Optional) The list of message center channels.
* `msc_contacts` - (Optional) The list of message center contacts.
* `name` - (Optional) The alert name.
* `threshold_type` - (Optional) The threshold type.
* `threshold_value` - (Optional) The threshold value.
* `warn_target` - (Optional) The alert target.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It is the same as `budget_name`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Budget.
* `delete` - (Defaults to 5 mins) Used when delete the Budget.
* `update` - (Defaults to 5 mins) Used when update the Budget.

## Import

Bss Open Api Budget can be imported using the id, e.g.

```shell
$ terraform import alicloud_bss_open_api_budget.example <budget_name>
```
