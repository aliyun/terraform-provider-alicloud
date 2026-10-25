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

A DTS check job (JobType = `CHECK`) performs data verification between the source and destination databases. The resource purchases a standalone `CHECK` instance through `CreateDtsInstance` and then configures the verification job through `ConfigureDtsJob` with `JobType` set to `CHECK`. When the job type is `CHECK`, `DataInitialization`, `DataSynchronization` and `StructureInitialization` are forced to `false`, and the verification behavior is controlled by `data_check_configure`.

For information about DTS Check Job and how to use it, see [What is DTS Job](https://www.alibabacloud.com/help/en/dts/developer-reference/api-configuredtsjob).

-> **NOTE:** Available since v2.1.1.

-> **NOTE:** `source_endpoint_vswitch_id`, `src_primary_vswitch_id`, `src_secondary_vswitch_id`, `dest_primary_vswitch_id` and `dest_secondary_vswitch_id` are only used when the job is created. They are not refreshed from the server and are not populated by `terraform import`.

-> **NOTE:** The endpoint attributes are refreshed from the job response. In particular, `destination_endpoint_owner_id` and `destination_endpoint_role` are read back from `DestinationEndpoint.AliyunUid` and `RoleName`, so a cross-account destination configuration is preserved by `terraform import` instead of planning a replacement.

-> **NOTE:** When `JobType` is `CHECK`, `data_initialization`, `data_synchronization` and `structure_initialization` are forced to `false` by the API. They are exported as computed attributes and cannot be configured.

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
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "example" {
  zone_id                  = data.alicloud_db_zones.example.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
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
  instance_storage         = data.alicloud_db_instance_classes.example.instance_classes.0.storage_range.0.min
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

resource "alicloud_dts_check_job" "example" {
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
  destination_endpoint_user_name     = alicloud_rds_account.example.1.account_name
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
      fullCheckModel       = 2
      dataCheckDbList = {
        "${alicloud_db_database.example.0.name}" = { name = alicloud_db_database.example.0.name, all = true }
      }
    }
  )
}
```

## Argument Reference

The following arguments supported:

* `dts_job_name` - (Required, String, Available since v2.1.1) The name of check job.
* `db_list` - (Required, String, ForceNew, Available since v2.1.1) Migration object, in the format of JSON strings. For detailed definition instructions, please refer to [the description of migration, synchronization or subscription objects](https://help.aliyun.com/document_detail/209545.html).
* `data_check_configure` - (Optional, String, ForceNew, Available since v2.1.1) The data verification configuration of the check job, in the format of a JSON string, such as parameter limits or alarm configurations. The `dataCheckDbList` field is required when this argument is set. For more information, see the DataCheckConfigure parameter description [datacheckconfigure-parameter](https://help.aliyun.com/zh/dts/developer-reference/datacheckconfigure-parameter).
* `instance_class` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The instance class of the check instance. Valid values: `large`, `medium`, `small`, `xlarge`, `2xlarge`, `4xlarge`. Default value: `small`. Changing this forces a new resource to be created; the class of an existing check job cannot be changed in place.
* `checkpoint` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The start point of the check job, the format is Unix timestamp, and the unit is seconds.
* `reserve` - (Optional, String, ForceNew, Computed, Available since v2.1.1) DTS reserves parameters, the format is a JSON string, you can pass in this parameter to complete the source and target database information (such as the data storage format of the target Kafka database, the instance ID of the cloud enterprise network CEN). For more information, please refer to the parameter [description of the Reserve parameter](https://help.aliyun.com/document_detail/273111.html).

  -> **NOTE:** The `srcSSL` and `destSSL` keys are managed by the properties `source_endpoint_ssl` and `destination_endpoint_ssl`. If either property is set, it overrides the corresponding key here.

* `source_endpoint_ssl` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The connection method of the source instance. Valid values: `0` (an unencrypted connection), `1` (an SSL-secured connection). Only supported when the source endpoint is accessed as a cloud instance or as a self-managed database hosted on ECS. Changing this property replaces the job because a check job cannot change SSL in place.
* `destination_endpoint_ssl` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The connection method of the destination instance. Valid values: `0` (an unencrypted connection), `1` (an SSL-secured connection), `3` (a connection using SCRAM-SHA-256, for a Kafka destination only). `1` is only supported when the destination endpoint is accessed as a cloud instance or as a self-managed database hosted on ECS. Changing this property replaces the job because a check job cannot change SSL in place.
* `source_endpoint_instance_type` - (Required, String, ForceNew, Available since v2.1.1) The type of source instance. Valid values: `CEN`, `DG`, `DISTRIBUTED_DMSLOGICDB`, `ECS`, `EXPRESS`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`, `REDIS`, `DISTRIBUTED_POLARDBX10`.
* `source_endpoint_engine_name` - (Required, String, ForceNew, Available since v2.1.1) The type of source database. Valid values: `AS400`, `DB2`, `DMSPOLARDB`, `HBASE`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `POSTGRESQL`, `TERADATA`, `POLARDB_PG`, `MARIADB`, `POLARDBX10`, `TiDB`, `REDIS`.
* `source_endpoint_instance_id` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The ID of source instance.
* `source_endpoint_region` - (Optional, String, ForceNew, Computed, Available since v2.1.1) Source instance area, please refer to the [list of supported areas](https://help.aliyun.com/document_detail/141033.htm) for details. Note if the source is an Alibaba Cloud database, this parameter must be passed in.
* `source_endpoint_ip` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The IP of source endpoint. When `source_endpoint_instance_type` is `OTHER`, `EXPRESS`, `DG`, `CEN`, this parameter is available and must be passed in.
* `source_endpoint_port` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The port of source endpoint. When the source instance is a self-built database, this parameter is available and must be passed in.
* `source_endpoint_oracle_sid` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The SID of Oracle database. When the value of SourceEndpointEngineName is Oracle and the Oracle database is a non-RAC instance, this parameter is available and must be passed in.
* `source_endpoint_database_name` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The name of the database to which the migration object belongs in the source instance.
* `source_endpoint_user_name` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The username of database account.
* `source_endpoint_password` - (Optional, String, Available since v2.1.1) The password of database account.
* `source_endpoint_owner_id` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The ID of Alibaba Cloud account to which the source instance belongs. Note: passing in this parameter means performing data migration or synchronization across Alibaba Cloud accounts, and you also need to pass in the `source_endpoint_role` parameter.
* `source_endpoint_role` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The name of the role configured for the cloud account to which the source instance belongs. Note: this parameter must be passed in when performing cross Alibaba Cloud account data migration or synchronization.
* `source_endpoint_vswitch_id` - (Optional, String, ForceNew, Available since v2.1.1) The ID of the vSwitch on the source side when the source database is connected through a VPC.
* `src_primary_vswitch_id` - (Optional, String, ForceNew, Available since v2.1.1) The ID of the primary vSwitch on the source side of a VPC NAT connection.
* `src_secondary_vswitch_id` - (Optional, String, ForceNew, Available since v2.1.1) The ID of the secondary vSwitch on the source side of a VPC NAT connection.
* `destination_endpoint_instance_type` - (Required, String, ForceNew, Available since v2.1.1) The type of destination instance. Valid values: `ADS`, `CEN`, `DATAHUB`, `DG`, `ECS`, `EXPRESS`, `GREENPLUM`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`, `REDIS`, `ELK`, `Tablestore`, `ODPS`.
* `destination_endpoint_engine_name` - (Required, String, ForceNew, Available since v2.1.1) The type of destination database. Valid values: `ADB20`, `ADS`, `ADB30`, `AS400`, `DATAHUB`, `DB2`, `GREENPLUM`, `KAFKA`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `PostgreSQL`, `POLARDB_PG`, `MARIADB`, `POLARDBX10`, `ODPS`, `Tablestore`, `ELK`, `REDIS`.
* `destination_endpoint_instance_id` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The ID of destination instance.
* `destination_endpoint_region` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The region of destination instance. For the target instance region, please refer to the [list of supported regions](https://help.aliyun.com/document_detail/141033.htm). Note: if the target is an Alibaba Cloud database, this parameter must be passed in.
* `destination_endpoint_ip` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The IP of destination endpoint. When `destination_endpoint_instance_type` is `OTHER`, `EXPRESS`, `DG`, `CEN`, this parameter is available and must be passed in.
* `destination_endpoint_port` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The port of destination endpoint. When the target instance is a self-built database, this parameter is available and must be passed in.
* `destination_endpoint_database_name` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The name of the database to which the migration object belongs in the target instance.
* `destination_endpoint_user_name` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The username of database account.
* `destination_endpoint_password` - (Optional, String, Available since v2.1.1) The password of database account.
* `destination_endpoint_oracle_sid` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The SID of Oracle database. Note: when the value of DestinationEndpointEngineName is Oracle and the Oracle database is a non-RAC instance, this parameter is available and must be passed in.
* `destination_endpoint_owner_id` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The ID of the Alibaba Cloud account to which the target RDS MySQL instance belongs. can be configured only when the target instance is RDS MySQL. This parameter is used to migrate or synchronize data across Alibaba Cloud accounts. You also need to enter the **destination_endpoint_role** parameter.
* `destination_endpoint_role` - (Optional, String, ForceNew, Computed, Available since v2.1.1) The role name of the Alibaba Cloud account to which the target instance belongs. This parameter must be entered when data migration or synchronization across Alibaba Cloud accounts is performed. For the permissions and authorization methods required by this role.
* `dest_primary_vswitch_id` - (Optional, String, ForceNew, Available since v2.1.1) The ID of the primary vSwitch on the destination side of a VPC NAT connection.
* `dest_secondary_vswitch_id` - (Optional, String, ForceNew, Available since v2.1.1) The ID of the secondary vSwitch on the destination side of a VPC NAT connection.
* `delay_notice` - (Optional, Bool, ForceNew, Available since v2.1.1) The delay notice. Valid values: `true`, `false`.
* `delay_phone` - (Optional, String, ForceNew, Available since v2.1.1) The delay phone. The mobile phone number of the contact who delayed the alarm. Multiple mobile phone numbers separated by English commas `,`. This parameter currently only supports China stations, and only supports mainland mobile phone numbers, and up to 10 mobile phone numbers can be passed in.
* `delay_rule_time` - (Optional, String, ForceNew, Available since v2.1.1) The delay rule time. When `delay_notice` is set to `true`, this parameter must be passed in. The threshold for triggering the delay alarm. The unit is second and needs to be an integer. The threshold can be set according to business needs. It is recommended to set it above 10 seconds to avoid delay fluctuations caused by network and database load.
* `error_notice` - (Optional, Bool, ForceNew, Available since v2.1.1) The error notice. Valid values: `true`, `false`.
* `error_phone` - (Optional, String, ForceNew, Available since v2.1.1) The error phone. The mobile phone number of the contact who error the alarm. Multiple mobile phone numbers separated by English commas `,`. This parameter currently only supports China stations, and only supports mainland mobile phone numbers, and up to 10 mobile phone numbers can be passed in.
* `status` - (Optional, String, Computed, Available since v2.1.1) The status of the resource. Valid values: `Migrating`, `Suspending`. You can stop the task by specifying `Suspending` and start the task by specifying `Migrating`. While the verification is in progress the computed value is `Migrating`, and once the one-shot check has completed it is `Finished`.
* `dedicated_cluster_id` - (Optional, String, ForceNew, Available since v2.1.1) When the ID of the dedicated cluster is input, the task is scheduled to the corresponding cluster.
* `dts_bis_label` - (Optional, String, ForceNew, Available since v2.1.1) The environment label of the DTS instance. The value is: **normal**, **online**.

## Attributes Reference

The following attributes are exported:

* `id` - (String, Available since v2.1.1) The resource ID in terraform of Check Job.
* `dts_instance_id` - (String, Available since v2.1.1) The ID of the DTS check instance created for the check job.
* `data_initialization` - (Bool, Available since v2.1.1) Whether to perform full data initialization. For check jobs, this value is always `false`.
* `data_synchronization` - (Bool, Available since v2.1.1) Whether to perform incremental data synchronization. For check jobs, this value is always `false`.
* `structure_initialization` - (Bool, Available since v2.1.1) Whether to perform library table structure initialization. For check jobs, this value is always `false`.

## Import

DTS Check Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_dts_check_job.example <id>
```
