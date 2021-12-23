---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_synchronization_jobs"
sidebar_current: "docs-alicloud-datasource-dts-synchronization-jobs"
description: |-
  Provides a list of Dts Synchronization Jobs to the user.
---

# alicloud\_dts\_synchronization\_jobs

This data source provides the Dts Synchronization Jobs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dts_synchronization_jobs" "ids" {}

output "dts_synchronization_job_id_1" {
  value = data.alicloud_dts_synchronization_jobs.ids.jobs.0.id
}
```

## Argument Reference

The following arguments supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Synchronization Job IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by synchronization job name.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Synchronizing`, `Suspending`. `Downgrade`, `Failed`, `Finished`, `InitializeFailed`, `Locked`, `Modifying`, `NotConfigured`, `NotStarted`, `PreCheckPass`, `PrecheckFailed`, `Prechecking`, `Retrying`, `Upgrade`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `jobs` - A list of Dts Synchronization Jobs. Each element contains the following attributes:
	* `id` - The ID of synchronizing instance. It's the ID of resource `alicloud_dts_synchronization_instance`.
	* `synchronization_direction` - Synchronization direction. Valid values: `Forward`, `Reverse`. Only when the property `sync_architecture` of the `alicloud_dts_synchronization_instance` was `bidirectional` this parameter should be passed, otherwise this parameter should not be specified.
	* `dts_job_name` - The name of synchronization job.
	* `checkpoint` - Start time in Unix timestamp format.
	* `data_initialization` - Whether to execute DTS supports schema migration, full data migration, or full-data initialization.
	* `data_synchronization` - Whether to perform incremental data migration for migration types or synchronization.
	* `structure_initialization` - Whether to perform a database table structure to migrate or initialization values include:
	* `db_list` - Migration object, in the format of JSON strings. For detailed definition instructions, please refer to [the description of migration, synchronization or subscription objects](https://help.aliyun.com/document_detail/209545.html).
	* `source_endpoint_instance_type` - The type of source instance. Valid values: `CEN`, `DG`, `DISTRIBUTED_DMSLOGICDB`, `ECS`, `EXPRESS`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`.
	* `source_endpoint_engine_name` - The type of source database. Valid values: `AS400`, `DB2`, `DMSPOLARDB`, `HBASE`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `POSTGRESQL`, `TERADATA`.
	* `source_endpoint_instance_id` - The ID of source instance.
	* `source_endpoint_region` - The region of source instance.
	* `source_endpoint_ip` - The ip of source endpoint.
	* `source_endpoint_port` - The port of source endpoint.
	* `source_endpoint_oracle_sid` - The SID of Oracle database.
	* `source_endpoint_database_name` - The name of migrate the database.
	* `source_endpoint_user_name` - The username of database account.
	* `source_endpoint_owner_id` - The Alibaba Cloud account ID to which the source instance belongs.
	* `source_endpoint_role` - The name of the role configured for the cloud account to which the source instance belongs.
	* `destination_endpoint_instance_type` - The type of destination instance. Valid values: `ads`, `CEN`, `DATAHUB`, `DG`, `ECS`, `EXPRESS`, `GREENPLUM`, `MONGODB`, `OTHER`, `PolarDB`, `POLARDBX20`, `RDS`.
	* `destination_endpoint_engine_name` - The type of destination database. Valid values: `ADB20`, `ADB30`, `AS400`, `DATAHUB`, `DB2`, `GREENPLUM`, `KAFKA`, `MONGODB`, `MSSQL`, `MySQL`, `ORACLE`, `PolarDB`, `POLARDBX20`, `POLARDB_O`, `PostgreSQL`.
	* `destination_endpoint_instance_id` - The ID of destination instance.
	* `destination_endpoint_region` - The region of destination instance.
	* `destination_endpoint_ip` - The ip of source endpoint.
	* `destination_endpoint_port` - The port of source endpoint.
	* `destination_endpoint_data_base_name` - The name of migrate the database.
	* `destination_endpoint_user_name` - The username of database account.
	* `destination_endpoint_oracle_sid` - The SID of Oracle database.
	* `status` - The status of the resource. Valid values: `Synchronizing`, `Suspending`. You can stop the task by specifying `Suspending` and start the task by specifying `Synchronizing`.
