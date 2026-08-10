---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_metric_stores"
sidebar_current: "docs-alicloud-datasource-sls-metric-stores"
description: |-
  Provides a list of Sls Metric Store owned by an Alibaba Cloud account.
---

# alicloud_sls_metric_stores

This data source provides Sls Metric Store available to the user. [What is Metric Store](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateMetricStore)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  description = "example"
  name        = "sls-sdk-examplep-metricstore"
}

resource "alicloud_sls_metric_store" "default" {
  project_name      = alicloud_log_project.default.name
  metering_mode     = "ChargeByFunction"
  mode              = "standard"
  metric_type       = "prometheus"
  metric_store_name = var.name
  ttl               = "7"
  shard_count       = "2"
}

data "alicloud_sls_metric_stores" "default" {
  project_name      = alicloud_log_project.default.name
  metric_store_name = alicloud_sls_metric_store.default.metric_store_name
  offset            = 0
  size              = 100
}

output "metric_store_example_id" {
  value = "${data.alicloud_sls_metric_stores.default.metric_stores.0.id}"
}
```

## Argument Reference

The following arguments are supported:
* `project_name` - (Required, ForceNew) The project name to the metric store belongs.
* `metric_store_name` - (Optional, ForceNew) The name of the metric store used to filter results.
* `mode` - (Optional, ForceNew) The mode of storage used to filter results. Valid value: `standard`.
* `offset` - (Optional, ForceNew) Query start row. The default value is `0`.
* `size` - (Optional, ForceNew) The number of rows per page set for a pagination query. The default value is `100`, and the maximum value is `500`.
* `ids` - (Optional, ForceNew, Computed) A list of Metric Store IDs. The value is formulated as `<project_name>:<metric_store_name>`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by metric store name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Metric Store IDs.
* `names` - A list of name of Metric Stores.
* `metric_stores` - A list of Metric Store Entries. Each element contains the following attributes:
  * `id` - The ID of the Metric Store. The value is formulated as `<project_name>:<metric_store_name>`.
  * `project_name` - The project name to the metric store belongs.
  * `metric_store_name` - The metric store name, unique in the same project.
  * `create_time` - Creation time of the metric store.
  * `last_modify_time` - Last updated time of the metric store.
  * `ttl` - Ttl for data storage, in days.
  * `shard_count` - The number of shards in the metric store.
  * `max_split_shard` - The maximum number of shards for automatic split.
  * `auto_split` - Determines whether to automatically split a shard.
  * `metric_type` - Type of metric store, defaults to `prometheus`.
  * `mode` - The mode of storage. Default to `standard`.
