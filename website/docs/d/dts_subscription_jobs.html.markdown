---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_subscription_jobs"
sidebar_current: "docs-alicloud-datasource-dts-subscription-jobs"
description: |-
  Provides a list of Dts Subscription Jobs to the user.
---

# alicloud\_dts\_subscription\_jobs

This data source provides the Dts Subscription Jobs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dts_subscription_jobs" "ids" {}
output "dts_subscription_job_id_1" {
  value = data.alicloud_dts_subscription_jobs.ids.jobs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Subscription Job IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by subscription job name.
* `status` - (Optional, ForceNew) The status of the task. Valid values: `Abnormal`, `Downgrade`, `Locked`, `Normal`, `NotStarted`, `NotStarted`, `PreCheckPass`, `PrecheckFailed`, `Prechecking`, `Retrying`, `Starting`, `Upgrade`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `jobs` - A list of Dts Subscription Jobs. Each element contains the following attributes:
	* `id` - The ID of the Subscription Job.
	* `dts_instance_id` - The ID of subscription job instance.
	* `dts_job_id` - The ID of subscription job instance.
	* `dts_job_name` - The name of subscription job instance.
	* `data_initialization` - Whether to execute DTS supports schema migration, full data migration, or full-data initialization values include:
	* `data_synchronization` - Whether to perform incremental data migration for migration types or synchronization values include:
	* `structure_initialization` - Whether to perform a database table structure to migrate or initialization values include:
	* `synchronization_direction` - Synchronization direction. Valid values: `Forward`, `Reverse`. Only when the property `sync_architecture` of the `alicloud_dts_synchronization_instance` was `bidirectional` this parameter should be passed, otherwise this parameter should not be specified.
	* `create_time` - The creation time of subscription job instance.
	* `checkpoint` - Subscription start time in Unix timestamp format.
	* `db_list` - Subscription object, in the format of JSON strings.
	* `expire_time` -  The Expiration Time. Formatting with yyyy-MM-ddTHH:mm:ssZ(UTC time).
	* `payment_type` - The payment type of the resource. Valid values: `Subscription`, `PayAsYouGo`.
	* `source_endpoint_database_name` - To subscribe to the name of the database.
	* `source_endpoint_engine_name` - The source database type value is MySQL or Oracle.
	* `source_endpoint_user_name` - The username of source database instance account.
	* `source_endpoint_instance_id` - The ID of source instance. Only when the type of source database instance was RDS MySQL, PolarDB-X 1.0, PolarDB MySQL, this parameter can be available and must be set.
	* `source_endpoint_instance_type` - The type of source instance. Valid values: `RDS`, `PolarDB`, `DRDS`, `LocalInstance`, `ECS`, `Express`, `CEN`, `dg`.
	* `source_endpoint_ip` - The IP of source endpoint.
	* `source_endpoint_oracle_sid` - The SID of Oracle Database. When the source database is self-built Oracle and the Oracle database is a non-RAC instance, this parameter is available and must be passed in.
	* `source_endpoint_owner_id` - The Alibaba Cloud account ID to which the source instance belongs. This parameter is only available when configuring data subscriptions across Alibaba Cloud accounts and must be passed in.
	* `source_endpoint_port` - The  port of source database.
	* `source_endpoint_region` - The region of source database.
	* `source_endpoint_role` - Both the authorization roles. When the source instance and configure subscriptions task of the Alibaba Cloud account is not the same as the need to pass the parameter, to specify the source of the authorization roles, to allow configuration subscription task of the Alibaba Cloud account to access the source of the source instance information.
	* `status` - The status of the task. Valid values: `NotStarted`, `Normal`, `Abnormal`. When a task created, it is in this state of `NotStarted`. You can specify this state of `Normal` to start the job, and specify this state of `Abnormal` to stop the job.
	* `subscription_data_type_ddl` - Whether to subscribe the DDL type of data. Valid values: `true`, `false`.
	* `subscription_data_type_dml` - Whether to subscribe the DML type of data. Valid values: `true`, `false`.
	* `subscription_host` - Network information.
		* `private_host` - Classic network address.
		* `public_host` - Public network address.
		* `vpc_host` - VPC network address.
	* `subscription_instance_network_type` - The type of subscription instance network. Valid value: `classic`, `vpc`.
	* `subscription_instance_vpc_id` - The ID of subscription instance vpc.
	* `subscription_instance_vswitch_id` - The ID of subscription instance vswitch.
	* `tags` - The tag of the resource.
		* `tag_key` - The key of the tags.
		* `tag_value` - The value of the tags.
