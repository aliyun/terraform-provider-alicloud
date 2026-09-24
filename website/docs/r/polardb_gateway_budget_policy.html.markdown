---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_gateway_budget_policy"
sidebar_current: "docs-alicloud-resource-polardb-gateway-budget-policy"
description: |-
  Provides a PolarDB AI Gateway Budget Policy resource.
---

# alicloud_polardb_gateway_budget_policy

Provides a PolarDB AI Gateway Budget Policy resource.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
resource "alicloud_polardb_gateway_budget_policy" "example" {
  gateway_id              = alicloud_polardb_gateway.example.id
  budget_type             = "ConsumerTotal"
  budget_dimension_ref_id = alicloud_polardb_gateway_consumer.example.consumer_id
  reset_day_of_month      = 1
  budget_points           = "10000"
  alert_threshold_pct     = 80
}
```

## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v1.294.0) The ID of the PolarDB AI gateway.
* `budget_type` - (Required, ForceNew, Available since v1.294.0) The budget type. Valid values: `GlobalTotal`, `ConsumerTotal`, `ConsumerGroupTotal`.
* `budget_dimension_ref_id` - (Optional, ForceNew, Available since v1.294.0) The consumer or consumer group ID. It is required for `ConsumerTotal` and `ConsumerGroupTotal` and must be omitted for `GlobalTotal`.
* `reset_day_of_month` - (Required, Available since v1.294.0) The monthly reset day. Valid values: `1` to `28`.
* `budget_points` - (Required, Available since v1.294.0) The budget amount in points.
* `alert_threshold_pct` - (Optional, Available since v1.294.0) The alert threshold percentage. Valid values: `0` to `100`.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<budget_policy_id>`.
* `budget_policy_id` - (Available since v1.294.0) The budget policy ID.
* `status` - (Available since v1.294.0) The policy status.
* `budget_dimension_type` - (Available since v1.294.0) The policy dimension.
* `used_points` - (Available since v1.294.0) The used points.
* `alert_triggered` - (Available since v1.294.0) Whether the alert threshold was triggered.
* `exceeded` - (Available since v1.294.0) Whether the budget was exceeded.
* `create_time` - (Available since v1.294.0) The creation time.
* `modify_time` - (Available since v1.294.0) The last modification time.

## Import

```shell
$ terraform import alicloud_polardb_gateway_budget_policy.example pg-abc:budget-abc
```
