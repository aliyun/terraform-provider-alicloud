---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_check_jobs"
sidebar_current: "docs-alicloud-datasource-dts-check-jobs"
description: |-
  Provides a list of Dts Check Jobs to the user.
---

# alicloud\_dts\_check\_jobs

This data source provides the Dts Check Jobs of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.236.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dts_check_jobs" "ids" {}

output "dts_check_job_id_1" {
  value = data.alicloud_dts_check_jobs.ids.jobs.0.id
}
```

## Argument Reference

The following arguments supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Check Job IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by check job name.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Synchronizing`, `Suspending`, `Downgrade`, `Failed`, `Finished`, `InitializeFailed`, `Locked`, `Modifying`, `NotConfigured`, `NotStarted`, `PreCheckPass`, `PrecheckFailed`, `Prechecking`, `Retrying`, `Upgrade`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Default to false. Set it to true to get more details of the check jobs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `jobs` - A list of Dts Check Jobs. Each element contains the following attributes:
	* `id` - The ID of check job.
	* `dts_job_id` - The ID of DTS job.
	* `dts_job_name` - The name of check job.
	* `checkpoint` - Start time in Unix timestamp format.
	* `data_initialization` - Whether to perform full data initialization. For check jobs, this value is always `false`.
	* `data_synchronization` - Whether to perform incremental data synchronization. For check jobs, this value is always `false`.
	* `structure_initialization` - Whether to perform library table structure initialization. For check jobs, this value is always `false`.
	* `db_list` - Migration object, in the format of JSON strings. For detailed definition instructions, please refer to [the description of migration, synchronization or subscription objects](https://help.aliyun.com/document_detail/209545.html).
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
	* `dts_instance_id` - The ID of DTS instance.
	* `status` - The status of the resource.
	* `create_time` - The creation time of the resource.
	* `expire_time` - The expiration time of the resource.
