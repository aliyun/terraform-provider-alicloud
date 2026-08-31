---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_migration_jobs"
sidebar_current: "docs-alicloud-datasource-dts-migration-jobs"
description: |-
  Provides a list of Dts Migration Jobs to the user.
---

# alicloud\_dts\_migration\_jobs

This data source provides the Dts Migration Jobs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.157.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dts_migration_jobs" "ids" {
  ids = ["dts_job_id"]
}
output "dts_migration_job_id_1" {
  value = data.alicloud_dts_migration_jobs.ids.jobs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Synchronization Job IDs.
* `name_regex` - A regex string to filter results by Migration Job name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Migration Job names.
* `jobs` - A list of Dts Migration Jobs. Each element contains the following attributes:
    * `id` - The ID of the Migration Job. Its value is same as `dts_job_id`.
    * `dts_job_id` - The ID of the Migration Job.
    * `dts_job_name` - The name of synchronization job.
    * `dts_instance_id` - The Migration instance ID. The ID of `alicloud_dts_migration_instance`.
    * `data_initialization` - Whether or not to execute DTS supports schema migration, full data migration, or full-data initialization.
    * `data_synchronization` - Whether to perform incremental data migration for migration types or synchronization values include:
    * `structure_initialization` - Whether to perform a database table structure to migrate or initialization.
    * `db_list` - The Migration object, in the format of JSON strings.
    * `payment_type` - The payment type of the Migration Instance. 
    * `source_endpoint_instance_type` - The type of source instance.
    * `source_endpoint_engine_name` - The type of source database. 
    * `source_endpoint_instance_id` - The ID of source instance.
    * `source_endpoint_region` - The region of source instance.
    * `source_endpoint_ip` - The ip of source endpoint.
    * `source_endpoint_port` - The port of source endpoint.
    * `source_endpoint_oracle_sid` - The SID of Oracle database.
    * `source_endpoint_database_name` - The name of migrate the database.
    * `source_endpoint_user_name` - The username of database account.
    * `source_endpoint_owner_id` - The Alibaba Cloud account ID to which the source instance belongs.
    * `source_endpoint_role` - The name of the role configured for the cloud account to which the source instance belongs.
    * `destination_endpoint_instance_type` - The type of destination instance. 
    * `destination_endpoint_engine_name` - The type of destination database. 
    * `destination_endpoint_instance_id` - The ID of destination instance.
    * `destination_endpoint_region` - The region of destination instance.
    * `destination_endpoint_ip` - The ip of source endpoint.
    * `destination_endpoint_port` - The port of source endpoint.
    * `destination_endpoint_data_base_name` - The name of migrate the database.
    * `destination_endpoint_user_name` - The username of database account.
    * `destination_endpoint_oracle_sid` - The SID of Oracle database.
    * `status` - The status of the resource. 