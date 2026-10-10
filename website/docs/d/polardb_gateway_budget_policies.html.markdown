---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_budget_policies"
sidebar_current: "docs-alicloud-datasource-polardb-gateway-budget-policies"
description: |-
  Provides a list of PolarDB AI Gateway Budget Policies.
---

# alicloud_polardb_gateway_budget_policies

Provides a list of PolarDB AI Gateway Budget Policies.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_polardb_gateway_budget_policies" "example" {
  gateway_id            = alicloud_polardb_gateway.example.id
  budget_dimension_type = "Consumer"
}
```

## Argument Reference

* `gateway_id` - (Required, Available since v1.294.0) The ID of the PolarDB AI gateway.
* `ids` - (Optional, Available since v1.294.0) A list of budget policy IDs.
* `status` - (Optional, Available since v1.294.0) The policy status.
* `budget_dimension_type` - (Optional, Available since v1.294.0) The policy dimension. Valid values: `ConsumerGroup`, `Consumer`.
* `budget_dimension_ref_id` - (Optional, Available since v1.294.0) The consumer group or consumer ID.
* `scope_ref_name` - (Optional, Available since v1.294.0) The dimension object name.

## Attributes Reference

* `policies` - A list of budget policies. Each element contains the following attributes:
  * `id` - The budget policy ID.
  * `budget_type` - The budget type.
  * `budget_dimension_type` - The policy dimension.
  * `budget_dimension_ref_id` - The dimension object ID.
  * `reset_day_of_month` - The monthly reset day.
  * `budget_points` - The budget amount in points.
  * `alert_threshold_pct` - The alert threshold percentage.
  * `status` - The policy status.
  * `used_points` - The used points.
  * `alert_triggered` - Whether the alert was triggered.
  * `exceeded` - Whether the budget was exceeded.
  * `create_time` - The creation time.
  * `modify_time` - The last modification time.
