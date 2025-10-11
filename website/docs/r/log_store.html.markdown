---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_store"
description: |-
  Provides a Alicloud Log Service (SLS) Log Store resource.
---

# alicloud_log_store

Provides a Log Service (SLS) Log Store resource.



For information about Log Service (SLS) Log Store and how to use it, see [What is Log Store](https://www.alibabacloud.com/help/doc-detail/48874.htm).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  name        = "terraform-example-${random_integer.default.result}"
  description = "terraform-example"
}

resource "alicloud_log_store" "example" {
  project               = alicloud_log_project.example.name
  name                  = "example-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
```

## Argument Reference

The following arguments are supported:
* `append_meta` - (Optional, Computed) Determines whether to append log meta automatically. The meta includes log receive time and client IP address. Default to `true`.
* `auto_split` - (Optional) Determines whether to automatically split a shard. Default to `false`.
* `enable_web_tracking` - (Optional) Whether open webtracking. webtracking network tracing, support the collection of HTML log, H5, Ios and android platforms.
* `encrypt_conf` - (Optional, Computed, List) Encrypted storage of data, providing data static protection capability, encrypt_conf can be updated since 1.188.0 (only enable change is supported when updating logstore). See [`encrypt_conf`](#encrypt_conf) below.
* `hot_ttl` - (Optional, Int) The ttl of hot storage. Default to 30, at least 30, hot storage ttl must be less than ttl.
* `infrequent_access_ttl` - (Optional, Int, Available since v1.229.0) Low frequency storage time
* `logstore_name` - (Optional, ForceNew, Available since v1.215.0) The log store, which is unique in the same project.
* `max_split_shard_count` - (Optional, Int) The maximum number of shards for automatic split, which is in the range of 1 to 256. You must specify this parameter when autoSplit is true.
* `metering_mode` - (Optional, Computed, Available since v1.215.0) Metering mode. The default metering mode of ChargeByFunction, ChargeByDataIngest traffic mode.
* `mode` - (Optional, Computed) The mode of storage. Default to `standard`, must be `standard` or `query`, `lite`.
* `project_name` - (Optional, ForceNew, Available since v1.215.0) The project name to the log store belongs.
* `retention_period` - (Optional, Computed, Int) The data retention time (in days). Valid values: [1-3650]. Default to 30. Log store data will be stored permanently when the value is 3650.
* `shard_count` - (Optional, ForceNew, Computed, Int) The number of shards in this log store. Default to 2. You can modify it by "Split" or "Merge" operations. [Refer to details](https://www.alibabacloud.com/help/zh/sls/product-overview/shard).
* `telemetry_type` - (Optional, ForceNew) Determines whether store type is metric. `Metrics` means metric store, empty means log store.

The following arguments will be discarded. Please use new fields as soon as possible:
* `project` - (Deprecated since v1.257.0). Field 'project' has been deprecated from provider version 1.257.0. New field 'project_name' instead.
* `name` - (Deprecated since v1.257.0). Field 'name' has been deprecated from provider version 1.257.0. New field 'logstore_name' instead.

### `encrypt_conf`

The encrypt_conf supports the following:
* `enable` - (Optional, Computed) Enable encryption. Default false
* `encrypt_type` - (Optional, ForceNew, Computed) Supported encryption type, only supports `default`(AES), `m4`.
* `user_cmk_info` - (Optional, ForceNew, Computed, List) User bring your own key (BYOK) encryption Refer to details, the format is as follows. See user_cmk_info below. `{ "cmk_key_id": "your_cmk_key_id", "arn": "your_role_arn", "region_id": "you_cmk_region_id" }` See [`user_cmk_info`](#encrypt_conf-user_cmk_info) below.

### `encrypt_conf-user_cmk_info`

The encrypt_conf-user_cmk_info supports the following:
* `arn` - (Optional, ForceNew, Computed) Role arn.
* `cmk_key_id` - (Optional, ForceNew, Computed) User master key id.
* `region_id` - (Optional, ForceNew, Computed) Region id where the user master key id is located.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_name>:<logstore_name>`.
* `shards` - The shard attribute.
  * `id` - The ID of the shard.
  * `status` - Shard status, only two status of `readwrite` and `readonly`.
  * `begin_key` - The begin value of the shard range(MD5), included in the shard range.
  * `end_key` - The end value of the shard range(MD5), not included in shard range.
* `create_time` - Log library creation time. Unix timestamp format that represents the number of seconds from 1970-1-1 00:00:00 UTC calculation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Log Store.
* `delete` - (Defaults to 5 mins) Used when delete the Log Store.
* `update` - (Defaults to 5 mins) Used when update the Log Store.

## Import

Log Service (SLS) Log Store can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_store.example <project_name>:<logstore_name>
```