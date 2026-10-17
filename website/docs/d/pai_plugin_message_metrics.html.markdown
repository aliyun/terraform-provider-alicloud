---
subcategory: "PAIPlugin"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_plugin_message_metrics"
description: |-
  This data source provides PAIPlugin message delivery metrics.
---

# alicloud_pai_plugin_message_metrics

This data source provides PAIPlugin message delivery metrics between the given
dates. Each metrics record is surfaced as parallel arrays (dates, fails,
pendings, rates, successes, totals) whose index aligns across all six lists,
preserving the order returned by the gateway.

For information about PAIPlugin message metrics and how to use it, see
[What is ListMessageMetrics](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiplugin-2022-01-12-listmessagemetrics).

-> **NOTE:** Available since v1.249.0.

## Example Usage

```terraform
data "alicloud_pai_plugin_message_metrics" "default" {
  start_date = "20260101"
  end_date   = "20260102"
}

output "first_date_totals" {
  value = data.alicloud_pai_plugin_message_metrics.default.totals.0
}
```

## Argument Reference

The following arguments are supported:

* `start_date` - (Required) The start date (yyyymmdd), inclusive.
* `end_date` - (Required) The end date (yyyymmdd), inclusive.
* `group_id` - (Optional) The group id filter.
* `schedule_id` - (Optional) The schedule id filter.
* `signature` - (Optional) The signature name filter.
* `signature_id` - (Optional) The signature id filter.
* `template_code` - (Optional) The template code filter.
* `template_id` - (Optional) The template id filter.
* `template_type` - (Optional) The template type filter. 0: verification code, 1: notification, 2: marketing.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source.
* `dates` - The dates of the metrics, in the order returned by the gateway.
* `fails` - The failed counts, aligned by index with `dates`.
* `pendings` - The pending counts, aligned by index with `dates`.
* `rates` - The success rates, aligned by index with `dates`.
* `successes` - The success counts, aligned by index with `dates`.
* `totals` - The total counts, aligned by index with `dates`.
* `total_count` - The total count of metrics records returned.
