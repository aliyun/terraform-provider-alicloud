---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_hybrid_monitor_datas"
sidebar_current: "docs-alicloud-datasource-cms-hybrid-monitor-datas"
description: |-
  Provides a list of Cms Hybrid Monitor Datas to the user.
---

# alicloud\_cms\_hybrid\_monitor\_datas

This data source provides the Cms Hybrid Monitor Datas of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.177.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_hybrid_monitor_datas" "default" {
  namespace = "example_value"
  prom_sql  = "AliyunEcs_cpu_total"
  start     = "1657505665"
  end       = "1657520065"
}
output "cms_metric_rule_template_id_1" {
  value = data.alicloud_cms_hybrid_monitor_datas.default.datas.0
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `namespace` - (Required, ForceNew) The name of the namespace.
* `prom_sql` - (Required, ForceNew) The name of the metric. Note PromQL statements are supported.
* `start` - (Required, ForceNew) The timestamp that specifies the beginning of the time range to query.
* `end` - (Required, ForceNew) The timestamp that specifies the end of the time range to query.
* `period` - (Optional, ForceNew) The interval at which monitoring data is collected. Unit: seconds.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `datas` - A list of Cms Hybrid Monitor Datas. Each element contains the following attributes:
  * `metric_name` - The name of the monitoring indicator.
  * `values` - The metric values that are collected at different timestamps.
    * `ts` - The timestamp that indicates the time when the metric value is collected. Unit: seconds.
    * `value` - The value of the monitoring indicator.
  * `labels` - The label of the time dimension.
    * `key` - Label key.
    * `value` - Label value.