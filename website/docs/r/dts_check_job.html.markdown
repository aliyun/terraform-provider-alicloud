---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_check_job"
sidebar_current: "docs-alicloud-resource-dts-check-job"
description: |-
  Provides a Alicloud DTS Check Job resource.
---

# alicloud_dts_check_job

Provides a DTS Check Job resource.

A DTS check job (JobType = `CHECK`) performs data verification between the source and destination databases. It is created through `ConfigureDtsJob` with `JobType` set to `CHECK`. When the job type is `CHECK`, `DataInitialization`, `DataSynchronization` and `StructureInitialization` are forced to `false`, and the verification behaviour is controlled by `data_check_configure`.

For information about DTS Check Job and how to use it, see [What is DTS Job](https://www.alibabacloud.com/help/en/dts/developer-reference/api-configuredtsjob).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}
data "alicloud_regions" "example" {
  current = true
}
data "alicloud_db_zones" "example" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "Basic"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "example" {
  zone_id                  = data.alicloud_db_zones.example.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "Basic"
  db_instance_storage_type = "cloud_essd"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.example.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_security_group" "example" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.example.id
}

resource "alicloud_db_instance" "example" {
  count                    = 2
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.example.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.example.instance_classes.0.storage_range.min
  instance_charge_type     = "Postpaid"
  instance_name            = format("%s_%d", var.name, count.index + 1)
  vswitch_id               = alicloud_vswitch.example.id
  monitoring_period        = "60"
  db_instance_storage_type = "cloud_essd"
  security_group_ids       = [alicloud_security_group.example.id]
}

resource "alicloud_rds_account" "example" {
  count            = 2
  db_instance_id   = alicloud_db_instance.example[count.index].id
  account_name     = format("example_name_%d", count.index + 1)
  account_password = format("example_password_%d", count.index + 1)
}

resource "alicloud_db_database" "example" {
  count       = 2
  instance_id = alicloud_db_instance.example[count.index].id
  name        = format("%s_%d", var.name, count.index + 1)
}

resource "alicloud_db_account_privilege" "example" {
  count        = 2
  instance_id  = alicloud_db_instance.example[count.index].id
  account_name = alicloud_rds_account.example[count.index].account_name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.example[count.index].name]
}

resource "alicloud_dts_migration_instance" "example" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.example.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.example.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

resource "alicloud_dts_check_job" "example" {
  dts_instance_id                    = alicloud_dts_migration_instance.example.id
  dts_job_name                       = var.name
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_account_privilege.example.0.instance_id
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = data.alicloud_regions.example.regions.0.id
  source_endpoint_user_name          = alicloud_rds_account.example.0.account_name
  source_endpoint_password           = alicloud_rds_account.example.0.account_password
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = alicloud_db_account_privilege.example.1.instance_id
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = data.alicloud_regions.example.regions.0.id
  destination_endpoint_user_name    = alicloud_rds_account.example.1.account_name
  destination_endpoint_password      = alicloud_rds_account.example.1.account_password
  db_list = jsonencode(
    {
      "${alicloud_db_database.example.0.name}" = { name = alicloud_db_database.example.1.name, all = true }
    }
  )
  data_check_configure = jsonencode(
    {
      fullDataCheck        = true
      incrementalDataCheck = false
      fullCheckModel       = "FLY_CHECK"
    }
  )
}
```

## Argument Reference

The following arguments supported:

* `dts_instance_id` - (Required, ForceNew) The ID of DTS instance, it must be an ID of `alicloud_dts_migration_instance` or `alicloud_dts_synchronization_instance`.
* `dts_job_name` - (Required) The name of check job.
* `db_list` - (Required) Migration object, in the format of JSON strings. For detailed definition instructions, please refer to [the description of migration, synchronization or subscription objects](https://help.aliyun.com/document_detail/209545.html).
* `data_check_configure` - (Optional, ForceNew) The data verification configuration of the check job, in the format of a JSON string, such as parameter limits or alarm configurations. For more information, see the DataCheckConfigure parameter description [datacheckconfigure-parameter](https://help.aliyun.com/zh/dts/developer-reference/datacheckconfigure-parameter).
* `instance_class` - (Optional) The instance class. Valid values: `large`, `medium`, `small`, `xlarge`, `2xlarge`, `4xlarge`. You can only upgrade the configuration, not downgrade the configuration. If you downgrade the instance, you need to [submit a ticket](https://selfservice.console.aliyun.com/ticket/category/dts/today).
* `checkpoint` - (Optional, ForceNew) The start point of the check job, the format is Unix timestamp, and the unit is seconds.
* `reserve` - (Optional) DTS reserves parameters, the format is a JSON string, you can pass in this parameter to complete the source and target database information (such as the data storage format of the target Kafka database, the instance ID of the cloud enterprise network CEN). For more information, please refer to the parameter [description of the Reserve parameter](https://help.aliyun.com/document_detail/273111.html).

  -> **NOTE:** The `srcSSL` and `destSSL` keys are managed by the properties `source_endpoint_ssl` and `destination_endpoint_ssl`. If either property is set, it overrides the corresponding key here.

* `source_endpoint_ssl` - (Optional, Computed) The connection method of the source instance. Valid values: `0` (an unencrypted connection), `1` (an SSL-secured connection). Only supported when the source endpoint is accessed as a cloud instance or as a self-managed database hosted on ECS.
* `destination_endpoint_ssl` - (Optional, Computed) The connection method of the destination instance. Valid values: `0` (an unencrypted connection), `1` (an SSL-secured connection), `3` (a connection using SCRAM-SHA-256, for a Kafka destination only). `1` is only supported when the destination endpoint is accessed as a cloud instance or as a self-managed database hosted on ECS.
* `job_parameters` - (Optional) DTS modifiable runtime parameters, you can modify the parameters of a running DTS (Data Transmission Service) task by providing a JSON array. This allows for real-time adjustments to the task's behavior. Please note that you can only modify these parameters while the task is active; they are not available during the initial setup. For more information, please refer to the parameter [description of the Runtime parameter](https://help.aliyun.com/zh/dts/developer-reference/parameter-description).
* `source_endpoint_instance_type` - (Required, ForceNew) The type of source instance. Valid values: `CEN`, `DG`, `DISTRIBUTED_DMSLOGICDB`, `ECS`, `EXPRESS`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`, `REDIS`, `DISTRIBUTED_POLARDBX10`.
* `source_endpoint_engine_name` - (Required, ForceNew) The type of source database. Valid values: `AS400`, `DB2`, `DMSPOLARDB`, `HBASE`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `POSTGRESQL`, `TERADATA`, `POLARDB_PG`, `MARIADB`, `POLARDBX10`, `TiDB`, `REDIS`.
* `source_endpoint_instance_id` - (Optional, ForceNew) The ID of source instance.
* `source_endpoint_region` - (Optional, ForceNew) Source instance area, please refer to the [list of supported areas](https://help.aliyun.com/document_detail/141033.htm) for details. Note if the source is an Alibaba Cloud database, this parameter must be passed in.
* `source_endpoint_ip` - (Optional, ForceNew) The IP of source endpoint. When `source_endpoint_instance_type` is `OTHER`, `EXPRESS`, `DG`, `CEN`, this parameter is available and must be passed in.
* `source_endpoint_port` - (Optional, ForceNew) The port of source endpoint. When the source instance is a self-built database, this parameter is available and must be passed in.
* `source_endpoint_oracle_sid` - (Optional, ForceNew) The SID of Oracle database. When the value of SourceEndpointEngineName is Oracle and the Oracle database is a non-RAC instance, this parameter is available and must be passed in.
* `source_endpoint_database_name` - (Optional, ForceNew) The name of the database to which the migration object belongs in the source instance.
* `source_endpoint_user_name` - (Optional, ForceNew) The username of database account.
* `source_endpoint_password` - (Optional) The password of database account.
* `source_endpoint_owner_id` - (Optional, ForceNew) The ID of Alibaba Cloud account to which the source instance belongs. Note: passing in this parameter means performing data migration or synchronization across Alibaba Cloud accounts, and you also need to pass in the `source_endpoint_role` parameter.
* `source_endpoint_role` - (Optional, ForceNew) The name of the role configured for the cloud account to which the source instance belongs. Note: this parameter must be passed in when performing cross Alibaba Cloud account data migration or synchronization.
* `destination_endpoint_instance_type` - (Required, ForceNew) The type of destination instance. Valid values: `ADS`, `CEN`, `DATAHUB`, `DG`, `ECS`, `EXPRESS`, `GREENPLUM`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`, `REDIS`, `ELK`, `Tablestore`, `ODPS`.
* `destination_endpoint_engine_name` - (Required, ForceNew) The type of destination database. Valid values: `ADB20`, `ADS`, `ADB30`, `AS400`, `DATAHUB`, `DB2`, `GREENPLUM`, `KAFKA`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `PostgreSQL`, `POLARDB_PG`, `MARIADB`, `POLARDBX10`, `ODPS`, `Tablestore`, `ELK`, `REDIS`.
* `destination_endpoint_instance_id` - (Optional, ForceNew) The ID of destination instance.
* `destination_endpoint_region` - (Optional, ForceNew) The region of destination instance. For the target instance region, please refer to the [list of supported regions](https://help.aliyun.com/document_detail/141033.htm). Note: if the target is an Alibaba Cloud database, this parameter must be passed in.
* `destination_endpoint_ip` - (Optional, ForceNew) The IP of source endpoint. When `destination_endpoint_instance_type` is `OTHER`, `EXPRESS`, `DG`, `CEN`, this parameter is available and must be passed in.
* `destination_endpoint_port` - (Optional, ForceNew) The port of source endpoint. When the target instance is a self-built database, this parameter is available and must be passed in.
* `destination_endpoint_database_name` - (Optional, ForceNew) The name of the database to which the migration object belongs in the target instance.
* `destination_endpoint_user_name` - (Optional, ForceNew) The username of database account.
* `destination_endpoint_password` - (Optional) The password of database account.
* `destination_endpoint_oracle_sid` - (Optional, ForceNew) The SID of Oracle database. Note: when the value of DestinationEndpointEngineName is Oracle and the Oracle database is a non-RAC instance, this parameter is available and must be passed in.
* `destination_endpoint_owner_id` - (Optional, ForceNew) The ID of the Alibaba Cloud account to which the target RDS MySQL instance belongs. can be configured only when the target instance is RDS MySQL. This parameter is used to migrate or synchronize data across Alibaba Cloud accounts. You also need to enter the **destinationendpointrle** parameter.
* `destination_endpoint_role` - (Optional, ForceNew) The role name of the Alibaba Cloud account to which the target instance belongs. This parameter must be entered when data migration or synchronization across Alibaba Cloud accounts is performed. For the permissions and authorization methods required by this role.
* `delay_notice` - (Optional, ForceNew) The delay notice. Valid values: `true`, `false`.
* `delay_phone` - (Optional, ForceNew) The delay phone. The mobile phone number of the contact who delayed the alarm. Multiple mobile phone numbers separated by English commas `,`. This parameter currently only supports China stations, and only supports mainland mobile phone numbers, and up to 10 mobile phone numbers can be passed in.
* `delay_rule_time` - (Optional, ForceNew) The delay rule time. When `delay_notice` is set to `true`, this parameter must be passed in. The threshold for triggering the delay alarm. The unit is second and needs to be an integer. The threshold can be set according to business needs. It is recommended to set it above 10 seconds to avoid delay fluctuations caused by network and database load.
* `error_notice` - (Optional, ForceNew) The error notice. Valid values: `true`, `false`.
* `error_phone` - (Optional, ForceNew) The error phone. The mobile phone number of the contact who error the alarm. Multiple mobile phone numbers separated by English commas `,`. This parameter currently only supports China stations, and only supports mainland mobile phone numbers, and up to 10 mobile phone numbers can be passed in.
* `status` - (Optional) The status of the resource. Valid values: `Migrating`, `Suspending`. You can stop the task by specifying `Suspending` and start the task by specifying `Migrating`.
* `dedicated_cluster_id` - (Optional, ForceNew) When the ID of the dedicated cluster is input, the task is scheduled to the corresponding cluster.
* `dts_bis_label` - (Optional, ForceNew) The environment label of the DTS instance. The value is: **normal**, **online**.

-> **NOTE:** When `JobType` is `CHECK`, `data_initialization`, `data_synchronization` and `structure_initialization` are forced to `false` by the API. They are exported as computed attributes and cannot be configured.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Check Job.
* `data_initialization` - Whether to perform full data initialization. For check jobs, this value is always `false`.
* `data_synchronization` - Whether to perform incremental data synchronization. For check jobs, this value is always `false`.
* `structure_initialization` - Whether to perform library table structure initialization. For check jobs, this value is always `false`.

## Import

DTS Check Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_dts_check_job.example <id>
```
