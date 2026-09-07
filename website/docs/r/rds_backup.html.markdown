---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_backup"
description: |-
  Provides a Alicloud RDS Backup resource.
---

# alicloud_rds_backup

Provides a RDS Backup resource.

Backup object at the instance level or database level.

For information about RDS Backup and how to use it, see [What is Backup](https://www.alibabacloud.com/help/en/rds/developer-reference/api-rds-2014-08-15-createbackup).

-> **NOTE:** Available since v1.149.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_backup&exampleId=e8067d24-c7ba-faf2-88e6-b2c5d35cf357bf1d9822&activeTab=example&spm=docs.r.rds_backup.0.e8067d24c7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_rds_backup&spm=docs.r.rds_backup.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `backup_method` - (Optional, ForceNew, Computed) The backup type. Valid values:  
* `Logical`: logical backup (supported only for MySQL)  
* `Physical`: physical backup (supported for MySQL, SQL Server, and PostgreSQL)  
* `Snapshot`: snapshot backup (supported for all database engines)  

Default value: `Physical`.  

-> **NOTE:**  * When using logical backup, the database must contain data (the data cannot be empty).  

-> **NOTE:**  * MariaDB instances support only snapshot backup, but you must specify `Physical` for this parameter.  

* `backup_retention_period` - (Optional, Int, Available since v1.273.0) When the database engine is SQL Server, `BackupStrategy` is set to `db`, `BackupMethod` is `Physical`, and `BackupType` is `FullBackup`, you can specify the retention period for the backup set. Valid values: 7 to 730 days, or - 1 (permanent retention).  

-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `backup_strategy` - (Optional) The backup strategy. Valid values:
* `db`: Single-database backup
* `instance`: Instance-level backup

-> **NOTE:** This parameter takes effect only under the following conditions:

-> **NOTE:**  - MySQL: The `BackupMethod` parameter is specified and set to `Logical`.

-> **NOTE:**  - SQL Server: The `BackupType` parameter is specified and set to `FullBackup`.


-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `backup_type` - (Optional, ForceNew, Computed) The backup type. Valid values:  
  - FullBackup: full backup  
  - IncrementalBackup: incremental backup  
* `db_instance_id` - (Required, ForceNew) The instance ID. You can call DescribeDBInstances to obtain it.
* `db_name` - (Optional) A list of databases, separated by commas (,).  

-> **NOTE:**  This parameter takes effect only when the `BackupStrategy` parameter is specified and its value is `db`.  


-> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `remove_from_state` - (Optional) Remove form state when resource cannot be deleted. Valid values: `true` and `false`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `backup_id` - The backup set ID.
* `status` - The status of the resource.
* `store_status` - Indicates whether the backup can be deleted.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Backup.
* `delete` - (Defaults to 5 mins) Used when delete the Backup.

## Import

RDS Backup can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_backup.example <backup_id>
```