---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_synchronization_job"
sidebar_current: "docs-alicloud-resource-dts-synchronization-job"
description: |-
  Provides a Alicloud DTS Synchronization Job resource.
---

# alicloud\_dts\_synchronization\_job

Provides a DTS Synchronization Job resource.

For information about DTS Synchronization Job and how to use it, see [What is Synchronization Job](https://www.alibabacloud.com/product/data-transmission-service).

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dts_synchronization_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "PolarDB"
  source_endpoint_region           = "cn-hangzhou"
  destination_endpoint_engine_name = "ADB30"
  destination_endpoint_region      = "cn-hangzhou"
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

resource "alicloud_dts_synchronization_job" "default" {
  dts_instance_id                    = alicloud_dts_synchronization_instance.default.id
  dts_job_name                       = "tf-testAccCase1"
  source_endpoint_instance_type      = "PolarDB"
  source_endpoint_instance_id        = "pc-xxxxxxxx"
  source_endpoint_engine_name        = "PolarDB"
  source_endpoint_region             = "cn-hangzhou"
  source_endpoint_database_name      = "tf-testacc"
  source_endpoint_user_name          = "root"
  source_endpoint_password           = "password"
  destination_endpoint_instance_type = "ads"
  destination_endpoint_instance_id   = "am-xxxxxxxx"
  destination_endpoint_engine_name   = "ADB30"
  destination_endpoint_region        = "cn-hangzhou"
  destination_endpoint_database_name = "tf-testacc"
  destination_endpoint_user_name     = "root"
  destination_endpoint_password      = "password"
  db_list                            = "{\"tf-testacc\":{\"name\":\"tf-test\",\"all\":true,\"state\":\"normal\"}}"
  structure_initialization           = "true"
  data_initialization                = "true"
  data_synchronization               = "true"
  status                             = "Synchronizing"
}
```

## Argument Reference

The following arguments supported:

* `dts_instance_id` - (Required, ForceNew) The ID of synchronization instance, it must be an ID of `alicloud_dts_synchronization_instance`.
* `dts_job_name` - (Required) The name of synchronization job.
* `data_initialization` - (Required, ForceNew) Whether to perform full data migration or full data initialization. Valid values: `true`, `false`.
* `data_synchronization` - (Required, ForceNew) Whether to perform incremental data migration or synchronization. Valid values: `true`, `false`.
* `structure_initialization` - (Required, ForceNew) Whether to perform library table structure migration or initialization. Valid values: `true`, `false`.
* `db_list` - (Required) Migration object, in the format of JSON strings. For detailed definition instructions, please refer to [the description of migration, synchronization or subscription objects](https://help.aliyun.com/document_detail/209545.html). **NOTE:** From version 1.173.0, `db_list` can be modified.
* `synchronization_direction` - (Optional, ForceNew) Synchronization direction. Valid values: `Forward`, `Reverse`. Only when the property `sync_architecture` of the `alicloud_dts_synchronization_instance` was `bidirectional` this parameter should be passed, otherwise this parameter should not be specified.
* `instance_class` - (Optional, Computed) The instance class. Valid values: `large`, `medium`, `micro`, `small`, `xlarge`, `xxlarge`. You can only upgrade the configuration, not downgrade the configuration. If you downgrade the instance, you need to [submit a ticket](https://selfservice.console.aliyun.com/ticket/category/dts/today).
* `checkpoint` - (Optional, Computed, ForceNew) The start point or synchronization point of incremental data migration, the format is Unix timestamp, and the unit is seconds.
* `reserve` - (Optional, Computed) DTS reserves parameters, the format is a JSON string, you can pass in this parameter to complete the source and target database information (such as the data storage format of the target Kafka database, the instance ID of the cloud enterprise network CEN). For more information, please refer to the parameter [description of the Reserve parameter](https://help.aliyun.com/document_detail/273111.html).
* `source_endpoint_instance_type` - (Required, ForceNew) The type of source instance. If the source instance is a `PolarDB O` engine cluster, the source instance type needs to be `OTHER` or `EXPRESS` as a self-built database, and access via public IP or dedicated line. For the correspondence between supported source and target instances, see [Supported Databases](https://help.aliyun.com/document_detail/131497.htm). When the source instance is a self-built database, you also need to perform corresponding preparations, for details, see [Preparations Overview](https://help.aliyun.com/document_detail/146958.htm). Valid values: `CEN`, `DG`, `DISTRIBUTED_DMSLOGICDB`, `ECS`, `EXPRESS`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`.
* `source_endpoint_engine_name` - (Required, ForceNew) The type of source database. The default value is `MySQL`. For the correspondence between supported source libraries and target libraries, see [Supported Databases](https://help.aliyun.com/document_detail/131497.htm). When the database type of the source instance is `MONGODB`, you also need to pass in some information in the reserved parameter `Reserve`, for the configuration method, see the description of Reserve parameters. Valid values: `AS400`, `DB2`, `DMSPOLARDB`, `HBASE`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `POSTGRESQL`, `TERADATA`.
* `source_endpoint_instance_id` - (Optional, ForceNew) The ID of source instance. If the source instance is a cloud database (such as RDS MySQL), you need to pass in the instance ID of the cloud database (such as the instance ID of RDS MySQL). If the source instance is a self-built database, the value of this parameter changes according to the value of `source_endpoint_instance_type`. For example, the value of `source_endpoint_instance_type` is:
  ** `ECS`, then this parameter needs to be passed into the instance ID of ECS.
  ** `DG`, then this parameter needs to be passed into the ID of database gateway.
  ** `EXPRESS`, `CEN`, then this parameter needs to be passed in the ID of VPC that has been interconnected with the source database. **Note**: when the value is `CEN`, you also need to pass in the ID of CEN instance in the cloud enterprise network with the reserved parameter `reserve`.
* `source_endpoint_region` - (Optional, ForceNew) Source instance area, please refer to the [list of supported areas](https://help.aliyun.com/document_detail/141033.htm) for details. Note if the source is an Alibaba Cloud database, this parameter must be passed in.
* `source_endpoint_ip` - (Optional, ForceNew) The IP of source endpoint. When `source_endpoint_instance_type` is `OTHER`, `EXPRESS`, `DG`, `CEN`, this parameter is available and must be passed in.
* `source_endpoint_port` - (Optional, ForceNew) The port of source endpoint. When the source instance is a self-built database, this parameter is available and must be passed in.
* `source_endpoint_oracle_sid` - (Optional, ForceNew) The SID of Oracle database. When the value of SourceEndpointEngineName is Oracle and the Oracle database is a non-RAC instance, this parameter is available and must be passed in.
* `source_endpoint_database_name` - (Optional, ForceNew) The name of the database to which the migration object belongs in the source instance. Note: this parameter is only available and must be passed in when the source instance, or the database type of the source instance is PolarDB O engine, PostgreSQL, or MongoDB database.
* `source_endpoint_user_name` - (Optional, ForceNew) The username of database account. Note: in most cases, you need to pass in the database account of the source library. The permissions required for migrating or synchronizing different databases are different. For specific permission requirements, see [Preparing database accounts for data migration](https://help.aliyun.com/document_detail/175878.htm) and [Preparing database accounts for data synchronization](https://help.aliyun.com/document_detail/213152.htm).
* `source_endpoint_password` - (Optional) The password of database account.
* `source_endpoint_owner_id` - (Optional, ForceNew) The ID of Alibaba Cloud account to which the source instance belongs. Note: passing in this parameter means performing data migration or synchronization across Alibaba Cloud accounts, and you also need to pass in the `source_endpoint_role` parameter.
* `source_endpoint_role` - (Optional, ForceNew) The name of the role configured for the cloud account to which the source instance belongs. Note: this parameter must be passed in when performing cross Alibaba Cloud account data migration or synchronization. For the permissions and authorization methods required by this role, please refer to [How to configure RAM authorization when cross-Alibaba Cloud account data migration or synchronization](https://help.aliyun.com/document_detail/48468.htm).
* `destination_endpoint_instance_type` - (Required, ForceNew) The type of destination instance. If the target instance is a PolarDB O engine cluster, the target instance type needs to be `OTHER` or `EXPRESS` as a self-built database, and access via public IP or dedicated line. If the target instance is the Kafka version of Message Queuing, the target instance type needs to be `ECS` or `EXPRESS` as a self-built database, and access via ECS or dedicated line. For the correspondence between supported targets and source instances, see [Supported Databases](https://help.aliyun.com/document_detail/131497.htm). When the target instance is a self-built database, you also need to perform corresponding preparations, please refer to the [overview of preparations](https://help.aliyun.com/document_detail/146958.htm). Valid values: `ADS`, `CEN`, `DATAHUB`, `DG`, `ECS`, `EXPRESS`, `GREENPLUM`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`.
* `destination_endpoint_engine_name` - (Required, ForceNew) The type of destination database. The default value is MYSQL. For the correspondence between supported target libraries and source libraries, see [Supported Databases](https://help.aliyun.com/document_detail/131497.htm). When the database type of the target instance is KAFKA or MONGODB, you also need to pass in some information in the reserved parameter `reserve`. For the configuration method, see the description of `reserve` parameters. Valid values: `ADS`, `ADB30`, `AS400`, `DATAHUB`, `DB2`, `GREENPLUM`, `KAFKA`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `PostgreSQL`.
* `destination_endpoint_instance_id` - (Optional, ForceNew) The ID of destination instance. If the target instance is a cloud database (such as RDS MySQL), you need to pass in the instance ID of the cloud database (such as the instance ID of RDS MySQL). If the target instance is a self-built database, the value of this parameter changes according to the value of `destination_endpoint_instance_type`. For example, the value of `destination_endpoint_instance_type` is:
  ** `ECS`, then this parameter needs to be passed into the instance ID of ECS.
  ** `DG`, then this parameter needs to be passed into the ID of database gateway.
  ** `EXPRESS`, `CEN`, then this parameter needs to be passed in the ID of VPC that has been interconnected with the source database. **Note**: when the value is `CEN`, you also need to pass in the ID of CEN instance in the cloud enterprise network with the reserved parameter `reserve`.
* `destination_endpoint_region` - (Optional, ForceNew) The region of destination instance. For the target instance region, please refer to the [list of supported regions](https://help.aliyun.com/document_detail/141033.htm). Note: if the target is an Alibaba Cloud database, this parameter must be passed in.
* `destination_endpoint_ip` - (Optional, ForceNew) The IP of source endpoint. When `destination_endpoint_instance_type` is `OTHER`, `EXPRESS`, `DG`, `CEN`, this parameter is available and must be passed in.
* `destination_endpoint_port` - (Optional, ForceNew) The port of source endpoint. When the target instance is a self-built database, this parameter is available and must be passed in.
* `destination_endpoint_database_name` - (Optional, ForceNew) The name of the database to which the migration object belongs in the target instance. Note: when the target instance or target database type is PolarDB O engine, AnalyticDB PostgreSQL, PostgreSQL, MongoDB database, this parameter is available and must be passed in.
* `destination_endpoint_user_name` - (Optional, ForceNew) The username of database account. Note: in most cases, you need to pass in the database account of the source library. The permissions required for migrating or synchronizing different databases are different. For specific permission requirements, see [Preparing database accounts for data migration](https://help.aliyun.com/document_detail/175878.htm) and [Preparing database accounts for data synchronization](https://help.aliyun.com/document_detail/213152.htm).
* `destination_endpoint_password` - (Optional) The password of database account.
* `destination_endpoint_oracle_sid` - (Optional, ForceNew) The SID of Oracle database. Note: when the value of DestinationEndpointEngineName is Oracle and the Oracle database is a non-RAC instance, this parameter is available and must be passed in.
* `delay_notice` - (Optional, ForceNew) The delay notice. Valid values: `true`, `false`.
* `delay_phone` - (Optional, ForceNew) The delay phone. The mobile phone number of the contact who delayed the alarm. Multiple mobile phone numbers separated by English commas `,`. This parameter currently only supports China stations, and only supports mainland mobile phone numbers, and up to 10 mobile phone numbers can be passed in.
* `delay_rule_time` - (Optional, ForceNew) The delay rule time. When `delay_notice` is set to `true`, this parameter must be passed in. The threshold for triggering the delay alarm. The unit is second and needs to be an integer. The threshold can be set according to business needs. It is recommended to set it above 10 seconds to avoid delay fluctuations caused by network and database load.
* `error_notice` - (Optional, ForceNew) The error notice. Valid values: `true`, `false`.
* `error_phone` - (Optional, ForceNew) The error phone. The mobile phone number of the contact who error the alarm. Multiple mobile phone numbers separated by English commas `,`. This parameter currently only supports China stations, and only supports mainland mobile phone numbers, and up to 10 mobile phone numbers can be passed in.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Synchronizing`, `Suspending`. You can stop the task by specifying `Suspending` and start the task by specifying `Synchronizing`.

-> **NOTE:** From the status of `NotStarted` to `Synchronizing`, the resource goes through the `Prechecking` and `Initializing` phases. Because of the `Initializing` phase takes too long, and once the resource passes to the status of `Prechecking`, it can be considered that the task can be executed normally. Therefore, we treat the status of `Initializing` as an equivalent to `Synchronizing`.

-> **NOTE:** If you want to upgrade the synchronization job specifications by the property `instance_class`, you must also modify the property `instance_class` of it's instance to keep them consistent.

## Notice

1. The expiration time cannot be changed after the work of the annual and monthly subscription suspended;
2. After the pay-as-you-go type job suspended, your job configuration fee will still be charged;
3. If the task suspended for more than 6 hours, the task will not start successfully.
4. Suspending the task will only stop writing to the target library, but will still continue to obtain the incremental log of the source, so that the task can be quickly resumed after the suspension is canceled. Therefore, some resources of the source library, such as bandwidth resources, will continue to be occupied during the period.
5. Charges will continue during the task suspension period. If you need to stop charging, please release the instance
6. When a DTS instance suspended for more than 7 days, the instance cannot be resumed, and the status will change from suspended to failed.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Synchronization Job.

## Import

DTS Synchronization Job can be imported using the id, e.g.

```
$ terraform import alicloud_dts_synchronization_job.example <id>
```
