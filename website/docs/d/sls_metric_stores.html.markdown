---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_metric_stores"
description: |-
  Provides a list of SLS MetricStores to the user.
---

# alicloud_sls_metric_stores

This data source provides the SLS MetricStores of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.291.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_sls_metric_stores" "default" {
  project_name = "your-project-name"
  name_regex   = ".*"
}

output "first_metric_store_name" {
  value = data.alicloud_sls_metric_stores.default.names.0
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required) The name of the SLS project.
* `ids` - (Optional) A list of MetricStore names to filter results.
* `name_regex` - (Optional) A regex string to apply to the MetricStore name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of MetricStore names.
* `metric_stores` - A list of MetricStores. Each element contains the following attributes:
  * `id` - The ID of the MetricStore. The value is formatted as `<project_name>:<metric_store_name>`.
  * `metric_store_name` - The name of the MetricStore.
  * `ttl` - The data retention period in days.
  * `shard_count` - The number of shards.
  * `auto_split` - Whether automatic shard splitting is enabled.
  * `max_split_shard_count` - The maximum number of shards to split.
  * `append_meta` - Whether to record the IP address of the requester.
  * `hot_ttl` - The data retention period in the hot storage tier, in days.
  * `mode` - The type of the MetricStore.
  * `create_time` - The time when the MetricStore was created.
  * `last_modify_time` - The time when the MetricStore was last modified.
