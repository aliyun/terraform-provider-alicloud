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
  payment_type                        = "PostPaid"
  source_endpoint_engine_name         = "PolarDB"
  source_endpoint_region              = "cn-hangzhou"
  destination_endpoint_engine_name    = "ADB30"
  destination_endpoint_region         = "cn-hangzhou"
  instance_class                      = "small"
  sync_architecture                   = "oneway"
}

resource "alicloud_dts_synchronization_job" "default" {
  dts_instance_id                     = alicloud_dts_synchronization_instance.default.id
  dts_job_name                        = "tf-testAccCase1"
  source_endpoint_instance_type       = "PolarDB"
  source_endpoint_instance_id         = "pc-xxxxxxxx"
  source_endpoint_engine_name         = "PolarDB"
  source_endpoint_region              = "cn-hangzhou"
  source_endpoint_database_name       = "tf-testacc"
  source_endpoint_user_name           = "root"
  source_endpoint_password            = "password"
  destination_endpoint_instance_type  = "ads"
  destination_endpoint_instance_id    = "am-xxxxxxxx"
  destination_endpoint_engine_name    = "ADB30"
  destination_endpoint_region         = "cn-hangzhou"
  destination_endpoint_database_name  = "tf-testacc"
  destination_endpoint_user_name      = "root"
  destination_endpoint_password       = "password"
  db_list                             = "{\"tf-testacc\":{\"name\":\"tf-test\",\"all\":true,\"state\":\"normal\"}}"
  structure_initialization            = "true"
  data_initialization                 = "true"
  data_synchronization                = "true"
  status                              = "Synchronizing"
}
```

## Argument Reference

The following arguments supported:

* `dts_instance_id` - (Required, ForceNew) Synchronizing instance ID. The ID of `alicloud_dts_synchronization_instance`.
* `synchronization_direction` - (Optional, ForceNew) Synchronization direction. Valid values: `Forward`, `Reverse`. Only when the property `sync_architecture` of the `alicloud_dts_synchronization_instance` was `bidirectional` this parameter should be passed, otherwise this parameter should not be specified.
* `dts_job_name` - (Optional, Computed) The name of synchronization job.
* `instance_class` - (Optional, Computed) The instance class. Valid values: `large`, `medium`, `micro`, `small`, `xlarge`, `xxlarge`. You can only upgrade the configuration, not downgrade the configuration. If you downgrade the instance, you need to [submit a ticket](https://selfservice.console.aliyun.com/ticket/category/dts/today).
* `checkpoint` - (Optional, Computed, ForceNew) Start time in Unix timestamp format.
* `data_initialization` - (Required, ForceNew) Whether or not to execute DTS supports schema migration, full data migration, or full-data initialization values include:
* `data_synchronization` - (Required, ForceNew) Whether to perform incremental data migration for migration types or synchronization values include:
* `structure_initialization` - (Required, ForceNew) Whether to perform a database table structure to migrate or initialization values include:
* `db_list` - (Required, ForceNew) Migration object, in the format of JSON strings. For detailed definition instructions, please refer to [the description of migration, synchronization or subscription objects](https://help.aliyun.com/document_detail/209545.html).
* `reserve` - (Optional, ForceNew) DTS reserves parameters, the format is a JSON string, you can pass in this parameter to complete the source and target database information (such as the data storage format of the target Kafka database, the instance ID of the cloud enterprise network CEN). For more information, please refer to the parameter [description of the Reserve parameter](https://help.aliyun.com/document_detail/273111.html).
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
* `destination_endpoint_instance_type` - (Required, ForceNew) The type of destination instance. Valid values: `ads`, `CEN`, `DATAHUB`, `DG`, `ECS`, `EXPRESS`, `GREENPLUM`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`.
* `destination_endpoint_engine_name` - (Required, ForceNew) The type of destination database. Valid values: `ADB20`, `ADB30`, `AS400`, `DATAHUB`, `DB2`, `GREENPLUM`, `KAFKA`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `PostgreSQL`.
* `destination_endpoint_instance_id` - (Optional, ForceNew) The ID of destination instance.
* `destination_endpoint_region` - (Optional, ForceNew) The region of destination instance.
* `destination_endpoint_ip` - (Optional, ForceNew) The ip of source endpoint.
* `destination_endpoint_port` - (Optional, ForceNew) The port of source endpoint.
* `destination_endpoint_data_base_name` - (Optional, ForceNew) The name of migrate the database.
* `destination_endpoint_user_name` - (Optional, ForceNew) The username of database account.
* `destination_endpoint_password` - (Optional, ForceNew) The password of database account.
* `destination_endpoint_oracle_sid` - (Optional, ForceNew) The SID of Oracle database.
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
4. Suspending the task will only stop writing to the target library, but will still continue to obtain the incremental log of the source, so that the task can be quickly resumed after the suspension is cancelled. Therefore, some resources of the source library, such as bandwidth resources, will continue to be occupied during the period.
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