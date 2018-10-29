---
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_backup_policy"
sidebar_current: "docs-alicloud-resource-kvstore-backup-policy"
description: |-
  Provides a backup policy for ApsaraDB Redis / Memcache instance resource.
---

# alicloud\_kvstore\_backup\_policy

Provides a backup policy for ApsaraDB Redis / Memcache instance resource. 

## Example Usage

```
resource "alicloud_kvstore_backup_policy" "redisbackup" {
  instance_id             = "${alicloud_kvstore_instance.myredis.id}"
  preferred_backup_time   = "00:00Z-04:00Z"
  preferred_backup_period = "Friday"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The id of ApsaraDB for Redis or Memcache intance.
* `preferred_backup_time`- (Required) Backup time, in the format of HH:mmZ- HH:mm Z
* `preferred_backup_period` - (Required) Backup Cycle. Allowed values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday

## Attributes Reference

The following attributes are exported:

* `id` - The id of the backup policy.
* `instance_id` - The id of ApsaraDB for Redis or Memcache intance.
* `preferred_backup_time`- Backup time, in the format of HH:mmZ- HH:mm Z
* `preferred_backup_period` - Backup Cycle. Allowed values: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday

## Import

KVStore backup policy can be imported using the id, e.g.

```
$ terraform import alicloud_kvstore_backup_policy.example r-abc12345678
```
