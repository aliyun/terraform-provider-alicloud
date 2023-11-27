---
subcategory: "SLS"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_store"
description: |-
  Provides a Alicloud SLS Log Store resource.
---

# alicloud_log_store

Provides a SLS Log Store resource. 

For information about SLS Log Store and how to use it, see [What is Log Store](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "logstore_name" {
  default = "logstore"
}

variable "project_name" {
  default = "terraform-logstore-test"
}

resource "alicloud_log_project" "defaultbRFbyS" {
  description = "terraform-logstore-test"
  name        = var.name

}


resource "alicloud_log_store" "default" {
  hot_ttl       = "7"
  logstore_name = var.name

  project_name    = alicloud_log_project.defaultbRFbyS.name
  max_split_shard = "0"
  ttl             = "20"
  shard_count     = "2"
  mode            = "query"
  telemetry_type  = "None"
  enable_tracking = true
  append_meta     = true
}
```

## Argument Reference

The following arguments are supported:
* `append_meta` - (Optional, Available since v1.0.0) Whether to turn on logging, network IP address.
* `auto_split` - (Optional, Available since v1.0.0) Whether to automatically split the shard.
* `enable_tracking` - (Optional) Whether open webtracking. webtracking network tracing, support the collection of HTML log, H5, Ios and android platforms.
* `hot_ttl` - (Optional, Available since v1.0.0) Hot ttl.
* `logstore_name` - (Required, ForceNew) LogstoreName.
* `max_split_shard` - (Optional) Automatically divide the maximum number of shard, the minimum value is 1, the maximum value is 64.
* `mode` - (Optional, Available since v1.0.0) Mode.
* `project_name` - (Required, ForceNew) Project.
* `shard_count` - (Required, ForceNew, Available since v1.0.0) ShardCount.
* `telemetry_type` - (Optional, ForceNew, Available since v1.0.0) Telemetry type.
* `ttl` - (Required) Ttl.

The following arguments will be discarded. Please use new fields as soon as possible:
* `project` - (Deprecated since v1.214.0). Field 'project' has been deprecated from provider version 1.214.0. New field 'project_name' instead.
* `name` - (Deprecated since v1.214.0). Field 'name' has been deprecated from provider version 1.214.0. New field 'logstore_name' instead.
* `retention_period` - (Deprecated since v1.214.0). Field 'retention_period' has been deprecated from provider version 1.214.0. New field 'ttl' instead.
* `max_split_shard_count` - (Deprecated since v1.214.0). Field 'max_split_shard_count' has been deprecated from provider version 1.214.0. New field 'max_split_shard' instead.
* `enable_web_tracking` - (Deprecated since v1.214.0). Field 'enable_web_tracking' has been deprecated from provider version 1.214.0. New field 'enable_tracking' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_name>:<logstore_name>`.
* `create_time` - Log library creation time. Unix timestamp format that represents the number of seconds from 1970-1-1 00:00:00 UTC calculation.
* `encrypt_conf` - Encrypt config.
  * `enable` - Enable.
  * `encrypt_type` - Encrypt type.
  * `user_cmk_info` - User CMK info.
    * `arn` - Arn.
    * `cmk_key_id` - Cmk key id.
    * `region_id` - Region id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Log Store.
* `delete` - (Defaults to 5 mins) Used when delete the Log Store.
* `update` - (Defaults to 5 mins) Used when update the Log Store.

## Import

SLS Log Store can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_store.example <project_name>:<logstore_name>
```