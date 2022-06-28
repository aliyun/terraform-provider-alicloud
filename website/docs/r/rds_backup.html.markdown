---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_backup"
sidebar_current: "docs-alicloud-resource-rds-backup"
description: |-
  Provides a Alicloud RDS Backup resource.
---

# alicloud\_rds\_backup

Provides a RDS Backup resource.

For information about RDS Backup and how to use it, see [What is Backup](https://www.alibabacloud.com/help/doc-detail/26272.htm).

-> **NOTE:** Available in v1.149.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_db_instance" "example" {
  engine                   = "MySQL"
  engine_version           = "5.6"
  instance_type            = "rds.mysql.t1.small"
  instance_storage         = "30"
  instance_charge_type     = "Postpaid"
  db_instance_storage_type = "local_ssd"
}

resource "alicloud_rds_backup" "example" {
  db_instance_id = alicloud_db_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `backup_method` - (Optional, Computed) The type of backup that you want to perform. Default value: `Physical`. Valid values: `Logical`, `Physical` and `Snapshot`.
* `backup_strategy` - (Optional) The policy that you want to use for the backup task. Valid values:
  * **db**: specifies to perform a database-level backup.
  * **instance**: specifies to perform an instance-level backup.
* `backup_type` - (Optional, Computed) The method that you want to use for the backup task. Default value: `Auto`. Valid values:
  * **Auto**: specifies to automatically perform a full or incremental backup.
  * **FullBackup**: specifies to perform a full backup.
* `db_instance_id` - (Required, ForceNew) The db instance id.
* `db_name` - (Optional) The names of the databases whose data you want to back up. Separate the names of the databases with commas (,).
* `remove_from_state` - (Optional) Remove form state when resource cannot be deleted. Valid values: `true` and `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Backup.
* `backup_id` - The backup id.
* `store_status` - Indicates whether the data backup file can be deleted. Valid values: `Enabled` and `Disabled`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the backup.
* `delete` - (Defaults to 20 mins) Used when deleting the backup.

## Import

RDS Backup can be imported using the id, e.g.

```
$ terraform import alicloud_rds_backup.example <db_instance_id>:<backup_id>
```
