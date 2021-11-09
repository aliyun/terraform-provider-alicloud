---
subcategory: "Serverless Workflow"
layout: "alicloud"
page_title: "Alicloud: alicloud_fnf_schedules"
sidebar_current: "docs-alicloud-datasource-fnf-schedules"
description: |-
  Provides a list of Fnf Schedules to the user.
---

# alicloud\_fnf\_schedules

This data source provides the Fnf Schedules of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.105.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_fnf_schedules" "example" {
  flow_name  = "example_value"
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_fnf_schedule_id" {
  value = data.alicloud_fnf_schedules.example.schedules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `flow_name` - (Required, ForceNew) The name of the flow bound to the time-based schedule you want to create.
* `ids` - (Optional, ForceNew, Computed) A list of Schedule IDs.
* `limit` - (Optional, ForceNew, Available in v1.110.0+) The number of resource queries.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Schedule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Schedule names.
* `schedules` - A list of Fnf Schedules. Each element contains the following attributes:
	* `cron_expression` - The CRON expression of the time-based schedule to be created.
	* `description` - The description of the time-based schedule to be created.
	* `enable` - Specifies whether to enable the time-based schedule you want to create.
	* `flow_name` - The name of the flow bound to the time-based schedule you want to create.
	* `id` - The ID of the Schedule.
	* `last_modified_time` - The time when the time-based schedule was last updated.
	* `payload` - The trigger message of the time-based schedule to be created. It must be in JSON object format.
	* `schedule_id` - The ID of the time-based schedule.
	* `schedule_name` - The name of the time-based schedule to be created.
