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

-> **NOTE:** Available since v2.1.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_gateway_budget_policy&exampleId=63ec8dc2-48f3-0e9e-02dd-d13d1d41f43ddc1b4aee&activeTab=example&spm=docs.r.polardb_gateway_budget_policy.0.63ec8dc248&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_gateway_budget_policy&spm=docs.r.polardb_gateway_budget_policy.example&intl_lang=EN_US)


## Argument Reference

* `gateway_id` - (Required, ForceNew, Available since v2.1.0) The ID of the PolarDB AI gateway.
* `budget_type` - (Required, ForceNew, Available since v2.1.0) The budget type. Valid values: `GlobalTotal`, `ConsumerTotal`, `ConsumerGroupTotal`.
* `budget_dimension_ref_id` - (Optional, ForceNew, Available since v2.1.0) The consumer or consumer group ID. It is required for `ConsumerTotal` and `ConsumerGroupTotal` and must be omitted for `GlobalTotal`.
* `reset_day_of_month` - (Required, Available since v2.1.0) The monthly reset day. Valid values: `1` to `28`.
* `budget_points` - (Required, Available since v2.1.0) The budget amount in points.
* `alert_threshold_pct` - (Optional, Available since v2.1.0) The alert threshold percentage. Valid values: `0` to `100`. The service does not support clearing this value after it is set; removing it from the configuration retains the remote value.

## Attributes Reference

* `id` - The resource ID in the format `<gateway_id>:<budget_policy_id>`.
* `budget_policy_id` - (Available since v2.1.0) The budget policy ID.
* `status` - (Available since v2.1.0) The policy status.
* `budget_dimension_type` - (Available since v2.1.0) The policy dimension.
* `used_points` - (Available since v2.1.0) The used points.
* `alert_triggered` - (Available since v2.1.0) Whether the alert threshold was triggered.
* `exceeded` - (Available since v2.1.0) Whether the budget was exceeded.
* `create_time` - (Available since v2.1.0) The creation time.
* `modify_time` - (Available since v2.1.0) The last modification time.

## Import

```shell
$ terraform import alicloud_polardb_gateway_budget_policy.example pg-abc:budget-abc
```
