---
subcategory: "Tair (Redis OSS-Compatible) And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_redis_backup"
description: |-
  Provides a Alicloud Tair (Redis OSS-Compatible) And Memcache (KVStore) Backup resource.
---

# alicloud_redis_backup

Provides a Tair (Redis OSS-Compatible) And Memcache (KVStore) Backup resource.

Instance level or database level backup objects.

For information about Tair (Redis OSS-Compatible) And Memcache (KVStore) Backup and how to use it, see [What is Backup](https://www.alibabacloud.com/help/en/redis/developer-reference/api-r-kvstore-2015-01-01-modifybackuppolicy-redis).

-> **NOTE:** Available since v1.15.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

variable "zone_id" {
  default = "cn-hangzhou-h"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  zone_id = var.zone_id
  vpc_id  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = var.zone_id
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_kvstore_instance" "default" {
  port           = "6379"
  payment_type   = "PrePaid"
  instance_type  = "Redis"
  password       = "123456_tf"
  engine_version = "5.0"
  zone_id        = var.zone_id
  vswitch_id     = local.vswitch_id
  period         = "1"
  instance_class = "redis.shard.small.2.ce"
}

resource "alicloud_redis_backup" "default" {
  instance_id             = alicloud_kvstore_instance.default.id
  backup_retention_period = 7
}
```

## Argument Reference

The following arguments are supported:
* `backup_retention_period` - (Optional, Int, Available since v1.266.0) The expiration period for this manual backup ranges from 7 to 730 days. When you pass in -1, it indicates that this manual backup will not expire (during the instance's lifecycle). If you do not pass any value (default case), it means the expiration policy will be consistent with the current automatic backup strategy.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `instance_id` - (Required, ForceNew) InstanceId

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<backup_id>`.
* `backup_id` - Backup ID.
* `status` - Backup status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 31 mins) Used when create the Backup.
* `delete` - (Defaults to 5 mins) Used when delete the Backup.

## Import

Tair (Redis OSS-Compatible) And Memcache (KVStore) Backup can be imported using the id, e.g.

```shell
$ terraform import alicloud_redis_backup.example <instance_id>:<backup_id>
```