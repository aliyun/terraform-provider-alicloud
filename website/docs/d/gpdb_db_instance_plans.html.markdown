---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_db_instance_plans"
sidebar_current: "docs-alicloud-datasource-gpdb-db-instance-plans"
description: |-
  Provides a list of Gpdb Db Instance Plans to the user.
---

# alicloud\_gpdb\_db\_instance\_plans

This data source provides the Gpdb Db Instance Plans of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.189.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_gpdb_db_instance_plans" "ids" {
  db_instance_id = "example_value"
  ids            = ["example_value"]
}
output "gpdb_db_instance_plan_id_1" {
  value = data.alicloud_gpdb_db_instance_plans.ids.plans.0.id
}

data "alicloud_gpdb_db_instance_plans" "nameRegex" {
  db_instance_id = "example_value"
  name_regex     = "^my-DBInstancePlan"
}
output "gpdb_db_instance_plan_id_2" {
  value = data.alicloud_gpdb_db_instance_plans.nameRegex.plans.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of DB Instance Plan IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by DB Instance Plan name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `db_instance_id` - (Required, ForceNew) The ID of the Database instance.
* `plan_schedule_type` - (Optional, ForceNew) Plan scheduling type. Valid values: `Postpone`, `Regular`.
* `plan_type` - (Optional, ForceNew) The type of the Plan. Valid values: `PauseResume`, `Resize`.
* `status` - (Optional, ForceNew) Planning Status. Valid values: `active`, `cancel`, `deleted`, `finished`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of DB Instance Plan names.
* `plans` - A list of Gpdb Db Instance Plans. Each element contains the following attributes:
	* `id` - The ID of the resource. The value formats as `<db_instance_id>:<plan_id>`.
	* `db_instance_plan_name` - The name of the Plan.
	* `plan_id` - The ID of DB Instance Plan.
	* `plan_config` - Plan configuration information.
		* `scale_out` - Scale out instance plan config.
			* `execute_time` - The executed time of the Plan.
			* `plan_cron_time` - The Cron Time of the plan.
			* `plan_task_status` - The Status of the plan Task.
			* `segment_node_num` - The segment Node Num of the Plan.
		* `pause` - Pause instance plan config.
			* `execute_time` - The executed time of the Plan.
			* `plan_cron_time` - The Cron Time of the plan.
			* `plan_task_status` - The Status of the plan Task.
		* `resume` - Resume instance plan config.
			* `plan_cron_time` - The Cron Time of the plan.
			* `plan_task_status` - The Status of the plan Task.
			* `execute_time` - The executed time of the Plan.
		* `scale_in` - Scale In instance plan config.
			* `plan_task_status` - The Status of the plan Task.
			* `segment_node_num` - The segment Node Num of the Plan.
			* `execute_time` - The executed time of the Plan.
			* `plan_cron_time` - The Cron Time of the plan.
	* `plan_end_date` - The end time of the Plan.
	* `plan_start_date` - The start time of the Plan.
	* `status` - The Status of the Plan.