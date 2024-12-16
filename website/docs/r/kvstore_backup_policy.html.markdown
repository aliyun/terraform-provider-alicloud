---
subcategory: "Tair (Redis OSS-Compatible) And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_backup_policy"
sidebar_current: "docs-alicloud-resource-kvstore-backup-policy"
description: |-
  Provides a Backup Policy for Tair (Redis OSS-Compatible) And Memcache (KVStore) resource.
---

# alicloud_kvstore_backup_policy

Provides a Backup Policy for Tair (Redis OSS-Compatible) And Memcache (KVStore) resource.

For information about Backup Policy for Tair (Redis OSS-Compatible) And Memcache (KVStore) and how to use it, see [What is Backup Policy](https://www.alibabacloud.com/help/en/redis/developer-reference/api-r-kvstore-2015-01-01-modifybackuppolicy-redis).

-> **NOTE:** Available since v1.15.0.

-> **DEPRECATED:**  This resource  has been deprecated from version `1.104.0`. Please use resource [alicloud_kvstore_instance](https://www.terraform.io/docs/providers/alicloud/r/kvstore_instance).

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_kvstore_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_kvstore_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_kvstore_instance" "default" {
  db_instance_name = var.name
  vswitch_id       = alicloud_vswitch.default.id
  zone_id          = data.alicloud_kvstore_zones.default.zones.0.id
  instance_class   = "redis.master.large.default"
  instance_type    = "Redis"
  engine_version   = "5.0"
  security_ips     = ["10.23.12.24"]
  config = {
    appendonly             = "yes"
    lazyfree-lazy-eviction = "yes"
  }
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_kvstore_backup_policy" "default" {
  instance_id   = alicloud_kvstore_instance.default.id
  backup_period = ["Tuesday", "Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The id of Tair (Redis OSS-Compatible) And Memcache (KVStore) Instance.
* `backup_time` - (Optional) Backup time, in the format of HH:mmZ- HH:mm Z.
* `backup_period` - (Optional) Backup Cycle. Allowed values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Backup Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `update` - (Defaults to 40 mins) Used when update the Backup Policy.

## Import

Tair (Redis OSS-Compatible) And Memcache (KVStore) Backup Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_kvstore_backup_policy.example <id>
```
