---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_migration_job"
sidebar_current: "docs-alicloud-resource-dts-migration-job"
description: |-
  Provides a Alicloud DTS Migration Job resource.
---

# alicloud\_dts\_migration\_job

Provides a DTS Migration Job resource.

For information about DTS Migration Job and how to use it, see [What is Migration Job](https://www.alibabacloud.com/help/en/doc-detail/208399.html).

-> **NOTE:** Available in v1.157.0+.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

variable "name" {
  default = "tftest"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_db_zones" "default" {}

data "alicloud_db_instance_classes" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_db_instance" "default" {
  count            = 2
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage = "10"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  instance_name    = join("", [var.name, count.index])
}

resource "alicloud_rds_account" "default" {
  count            = 2
  db_instance_id   = alicloud_db_instance.default[count.index].id
  account_name     = join("", [var.name, count.index])
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  count       = 2
  instance_id = alicloud_db_instance.default[count.index].id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  count        = 2
  instance_id  = alicloud_db_instance.default[count.index].id
  account_name = alicloud_rds_account.default[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default[count.index].name]
}

resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = var.region
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = var.region
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

resource "alicloud_dts_migration_job" "default" {
  dts_instance_id                    = alicloud_dts_migration_instance.default.id
  dts_job_name                       = var.name
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.default.0.id
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = var.region
  source_endpoint_user_name          = alicloud_rds_account.default.0.name
  source_endpoint_password           = var.password
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = alicloud_db_instance.default.1.id
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = var.region
  destination_endpoint_user_name     = alicloud_rds_account.default.1.name
  destination_endpoint_password      = var.password
  db_list                            = "{\"tftestdatabase\":{\"name\":\"tftestdatabase\",\"all\":true}}"
  structure_initialization           = true
  data_initialization                = true
  data_synchronization               = true
  status                             = "Migrating"
  depends_on                         = [alicloud_db_account_privilege.default]
}
```

## Argument Reference

The following arguments supported:

* `dts_instance_id` - (Required, ForceNew) The Migration instance ID. The ID of `alicloud_dts_migration_instance`.
* `dts_job_name` - (Optional, Computed, ForceNew) The name of migration job.
* `instance_class` - (Optional, Computed) The instance class. Valid values: `large`, `medium`, `micro`, `small`, `xlarge`, `xxlarge`. 
* `checkpoint` - (Optional, Computed, ForceNew) Start time in Unix timestamp format.
* `data_initialization` - (Required, ForceNew) Whether to execute DTS supports schema migration.
* `structure_initialization` - (Required, ForceNew) Whether to perform a database table structure to migrate.
* `data_synchronization` - (Required, ForceNew) Whether to perform incremental data migration.
* `db_list` - (Required, ForceNew) Migration object, in the format of JSON strings. For detailed definition instructions, please refer to [the description of migration, migration or subscription objects](https://help.aliyun.com/document_detail/209545.html).
* `source_endpoint_instance_type` - (Required, ForceNew) The type of source instance. Valid values: `CEN`, `DG`, `DISTRIBUTED_DMSLOGICDB`, `ECS`, `EXPRESS`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`.
* `source_endpoint_engine_name` - (Required, ForceNew) The type of source database. Valid values: `AS400`, `DB2`, `DMSPOLARDB`, `HBASE`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `POSTGRESQL`, `TERADATA`.
* `source_endpoint_instance_id` - (Optional, ForceNew) The ID of source instance.
* `source_endpoint_region` - (Optional, ForceNew) The region of source instance.
* `source_endpoint_ip` - (Optional, ForceNew) The ip of source endpoint.
* `source_endpoint_port` - (Optional, ForceNew) The port of source endpoint.
* `source_endpoint_oracle_sid` - (Optional, ForceNew) The SID of Oracle database.
* `source_endpoint_database_name` - (Optional, ForceNew) The name of migrate the database.
* `source_endpoint_user_name` - (Optional, ForceNew) The username of database account.
* `source_endpoint_password` - (Optional) The password of database account.
* `source_endpoint_owner_id` - (Optional, ForceNew) The Alibaba Cloud account ID to which the source instance belongs.
* `source_endpoint_role` - (Optional, ForceNew) The name of the role configured for the cloud account to which the source instance belongs.
* `destination_endpoint_instance_type` - (Required, ForceNew) The type of destination instance. Valid values: `ADS`, `CEN`, `DATAHUB`, `DG`, `ECS`, `EXPRESS`, `GREENPLUM`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`.
* `destination_endpoint_engine_name` - (Required, ForceNew) The type of destination database. Valid values: `ADS`, `ADB30`, `AS400`, `DATAHUB`, `DB2`, `GREENPLUM`, `KAFKA`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `PostgreSQL`.
* `destination_endpoint_instance_id` - (Optional, ForceNew) The ID of destination instance.
* `destination_endpoint_region` - (Optional, ForceNew) The region of destination instance.
* `destination_endpoint_ip` - (Optional, ForceNew) The ip of source endpoint.
* `destination_endpoint_port` - (Optional, ForceNew) The port of source endpoint.
* `destination_endpoint_database_name` - (Optional, ForceNew) The name of migrate the database.
* `destination_endpoint_user_name` - (Optional, ForceNew) The username of database account.
* `destination_endpoint_password` - (Optional) The password of database account.
* `destination_endpoint_oracle_sid` - (Optional, ForceNew) The SID of Oracle database.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Migrating`, `Suspending`. You can suspend the task by specifying `Suspending` and start the task by specifying `Migrating`.

## Notice

1. The expiration time cannot be changed after the work of the annual and monthly subscription suspended;
2. After the pay-as-you-go type job suspended, your job configuration fee will still be charged;
3. If the task suspended for more than 6 hours, the task will not start successfully.
4. Suspending the task will only stop writing to the target library, but will still continue to obtain the incremental log of the source, so that the task can be quickly resumed after the suspension is canceled. Therefore, some resources of the source library, such as bandwidth resources, will continue to be occupied during the period.
5. Charges will continue during the task suspension period. If you need to stop charging, please release the instance
6. When a DTS instance suspended for more than 7 days, the instance cannot be resumed, and the status will change from suspended to failed.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Migration Job.

## Import

DTS Migration Job can be imported using the id, e.g.

```
$ terraform import alicloud_dts_migration_job.example <id>
```