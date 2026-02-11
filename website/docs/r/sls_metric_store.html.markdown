---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_metric_store"
description: |-
  Provides a Alicloud Log Service (SLS) Metric Store resource.
---

# alicloud_sls_metric_store

Provides a Log Service (SLS) Metric Store resource.



For information about Log Service (SLS) Metric Store and how to use it, see [What is Metric Store](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateMetricStore).

-> **NOTE:** Available since v1.271.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  description  = "example"
  project_name = "sls-sdk-examplep-metricstore"
}

resource "alicloud_sls_metric_store" "default" {
  project_name      = alicloud_log_project.default.project_name
  metering_mode     = "ChargeByFunction"
  mode              = "standard"
  metric_type       = "prometheus"
  metric_store_name = "example-metric"
  ttl               = "7"
  shard_count       = "2"
}
```

## Argument Reference

The following arguments are supported:
* `auto_split` - (Optional, Computed) Determines whether to automatically split a shard. Default to `false`.
* `hot_ttl` - (Optional, Int) The ttl of hot storage. Default to 30, at least 30, hot storage ttl must be less than ttl.
* `infrequent_access_ttl` - (Optional, Int) Low frequency storage time
* `max_split_shard` - (Optional, Int) The maximum number of shards for automatic split, which is in the range of 1 to 256. You must specify this parameter when autoSplit is true.
* `metering_mode` - (Optional, Computed) Metering mode of the metricStore, ChargeByFunction or ChargeByDataIngest
* `metric_store_name` - (Required, ForceNew) The metric store, unique in the same project.
* `metric_type` - (Optional, ForceNew, Computed) Type of metric store, defaults to prometheus.
* `mode` - (Optional, Computed) The mode of storage. Default to `standard`, must be `standard.
* `project_name` - (Required, ForceNew) The project name to the metric store belongs.
* `shard_count` - (Required, ForceNew, Int) The number of shards in the metric store.
* `ttl` - (Required, Int) Ttl for data storage, in days.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<project_name>:<metric_store_name>`.
* `create_time` - Creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Metric Store.
* `delete` - (Defaults to 5 mins) Used when delete the Metric Store.
* `update` - (Defaults to 5 mins) Used when update the Metric Store.

## Import

Log Service (SLS) Metric Store can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_metric_store.example <project_name>:<metric_store_name>
```