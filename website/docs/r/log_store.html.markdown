---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_store"
sidebar_current: "docs-alicloud-resource-log-store"
description: |-
  Provides a Alicloud log store resource.
---

# alicloud\_log\_store

The log store is a unit in Log Service to collect, store, and query the log data. Each log store belongs to a project,
and each project can create multiple Logstores. [Refer to details](https://www.alibabacloud.com/help/doc-detail/48874.htm)

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "alicloud_log_store" "example" {
  project               = alicloud_log_project.example.name
  name                  = "tf-log-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
```
Encrypt Usage
```
resource "alicloud_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "alicloud_log_store" "example" {
  project               = alicloud_log_project.example.name
  name                  = "tf-log-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
  encrypt_conf {
    enable              = true
    encrypt_type        = "default"
    user_cmk_info {
      cmk_key_id        = "your_cmk_key_id"
      arn               = "your_role_arn"
      region_id         = "you_cmk_region_id"
    }
  }
}
```

## Module Support

You can use the existing [sls module](https://registry.terraform.io/modules/terraform-alicloud-modules/sls/alicloud) 
to create SLS project, store and store index one-click, like ECS instances.

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs.
* `name` - (Required, ForceNew) The log store, which is unique in the same project.
* `retention_period` - (Optional) The data retention time (in days). Valid values: [1-3650]. Default to `30`. Log store data will be stored permanently when the value is `3650`.
* `shard_count` - (Optional) The number of shards in this log store. Default to 2. You can modify it by "Split" or "Merge" operations. [Refer to details](https://www.alibabacloud.com/help/doc-detail/28976.htm)
* `auto_split` - (Optional) Determines whether to automatically split a shard. Default to `false`.
* `max_split_shard_count` - (Optional) The maximum number of shards for automatic split, which is in the range of 1 to 64. You must specify this parameter when autoSplit is true.
* `append_meta` - (Optional) Determines whether to append log meta automatically. The meta includes log receive time and client IP address. Default to `true`.
* `enable_web_tracking` - (Optional) Determines whether to enable Web Tracking. Default `false`.
* `encrypt_conf` (ForceNew, Optional, Available in 1.124.0+) Encrypted storage of data, providing data static protection capability, only supported at creation time.
  * `enable` (Optional) enable encryption. Default `false`
  * `encrypt_type` (Optional) Supported encryption type, only supports `default(AES)`,` m4`
  * `user_cmk_info` (Optional) User bring your own key (BYOK) encryption [Refer to details](https://www.alibabacloud.com/help/zh/doc-detail/187853.htm), the format is as follows:
    
    ```
    {
      "cmk_key_id": "your_cmk_key_id",
      "arn":        "your_role_arn",
      "region_id":  "you_cmk_region_id"
    }
    ```
#### Block user_cmk_info
The user_cmk_info mapping supports the following:

* `cmk_key_id` (Required) User master key id.
* `arn` (Required) role arn.
* `region_id` (Required) Region id where the  user master key id is located.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log project. It formats of `<project>:<name>`.
* `project` - The project name.
* `name` - Log store name.
* `retention_period` - The data retention time.
* `shard_count` - The number of shards.
* `shards` - The shard attribute.
  * `id` - The ID of the shard.
  * `status` - Shard status, only two status of `readwrite` and `readonly`.
  * `begin_key` - The begin value of the shard range(MD5), included in the shard range.
  * `end_key` - The end value of the shard range(MD5), not included in shard range.
* `auto_split` - Determines whether to automatically split a shard.
* `max_split_shard_count` - The maximum number of shards for automatic split.
* `append_meta` - Determines whether to append log meta automatically.
* `enable_web_tracking` - Determines whether to enable Web Tracking.
* `encrypt_conf` - Encryption configuration of logstore.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create`  - (Defaults to 3 mins) Used when Creating LogStore. 
* `delete`  - (Defaults to 3 mins) Used when Deleting LogStore.
* `read`    - (Defaults to 2 mins) Used when Reading LogStore.

## Import

Log store can be imported using the id, e.g.

```
$ terraform import alicloud_log_store.example tf-log:tf-log-store
```
