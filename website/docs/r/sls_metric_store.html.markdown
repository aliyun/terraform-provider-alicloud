---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_metric_store"
description: |-
  Provides a Alicloud SLS MetricStore resource.
---

# alicloud_sls_metric_store

Provides a SLS MetricStore resource. A MetricStore is a dedicated time-series data store in Simple Log Service for storing and querying metrics.

For information about SLS MetricStore and how to use it, see [MetricStore](https://www.alibabacloud.com/help/en/sls/developer-reference/api-sls-2020-12-30-dir-metricstore/).

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

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = "terraform example project for metric store"
}

resource "alicloud_sls_metric_store" "default" {
  project_name          = alicloud_log_project.default.project_name
  metric_store_name     = "${var.name}-metric-store"
  ttl                   = 30
  shard_count           = 2
  auto_split            = true
  max_split_shard_count = 64
  append_meta           = false
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The name of the SLS project.
* `metric_store_name` - (Required, ForceNew) The name of the MetricStore.
* `ttl` - (Required) The data retention period in days. Minimum value is 1.
* `shard_count` - (Required, ForceNew) The number of shards for the MetricStore.
* `auto_split` - (Optional) Specifies whether to automatically split shards. Default value: `true`.
* `max_split_shard_count` - (Optional) The maximum number of shards to split when auto-split is enabled. Valid values: `0` to `256`. Default value: `64`.
* `append_meta` - (Optional) Specifies whether to record the IP address of the requester. Default value: `false`.
* `hot_ttl` - (Optional) The data retention period in the hot storage tier, in days. Minimum value is `7`. Set to `-1` to retain data in the hot storage tier for the entire TTL.
* `infrequent_access_ttl` - (Optional) The data retention period in the Infrequent Access (IA) storage tier, in days. It must be greater than `60`, and `hot_ttl + infrequent_access_ttl` must not be greater than `ttl`.
* `mode` - (Optional, ForceNew) The type of the MetricStore. Valid values: `standard`. Default value: `standard`.
* `encrypt_conf` - (Optional) The encryption configuration. See [`encrypt_conf`](#encrypt_conf) below.

### `encrypt_conf`

The `encrypt_conf` block supports:

* `enable` - (Required) Specifies whether to enable encryption.
* `encrypt_type` - (Optional, ForceNew) The encryption algorithm. Valid values: `default`.
* `user_cmk_info` - (Optional, ForceNew) The BYOK (Bring Your Own Key) configuration. See [`user_cmk_info`](#encrypt_conf-user_cmk_info) below.

### `encrypt_conf-user_cmk_info`

The encrypt_conf-user_cmk_info supports the following:

* `cmk_key_id` - (Optional, ForceNew) The ID of the CMK (Customer Master Key).
* `arn` - (Optional) The ARN of the RAM role that is authorized to use the CMK.
* `region_id` - (Optional, ForceNew) The region ID of the CMK.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MetricStore. The value is formatted as `<project_name>:<metric_store_name>`.
* `create_time` - The time when the MetricStore was created.
* `last_modify_time` - The time when the MetricStore was last modified.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the MetricStore.
* `update` - (Defaults to 5 minutes) Used when updating the MetricStore.
* `delete` - (Defaults to 5 minutes) Used when deleting the MetricStore.

## Import

SLS MetricStore can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_metric_store.example <project_name>:<metric_store_name>
```
