---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_stacks"
sidebar_current: "docs-alicloud-datasource-ros-stacks"
description: |-
  Provides a list of Ros Stacks to the user.
---

# alicloud\_ros\_stacks

This data source provides the Ros Stacks of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.106.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ros_stacks" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_ros_stack_id" {
  value = data.alicloud_ros_stacks.example.stacks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Stack IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Stack name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `parent_stack_id` - (Optional, ForceNew) ParentStackId.
* `show_nested_stack` - (Optional, ForceNew) The show nested stack.
* `stack_name` - (Optional, ForceNew) StackName.
* `status` - (Optional, ForceNew) The status of Stack. Valid Values: `CREATE_COMPLETE`, `CREATE_FAILED`, `CREATE_IN_PROGRESS`, `DELETE_COMPLETE`, `DELETE_FAILED`, `DELETE_IN_PROGRESS`, `ROLLBACK_COMPLETE`, `ROLLBACK_FAILED`, `ROLLBACK_IN_PROGRESS`.
* `tags` - (Optional) Query the instance bound to the tag. The format of the incoming value is `json` string, including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty. Format example `{"key1":"value1"}`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Stack names.
* `stacks` - A list of Ros Stacks. Each element contains the following attributes:
	* `deletion_protection` - Specifies whether to enable deletion protection on the stack.
	* `description` - The Description of the Stack.
	* `disable_rollback` - Specifies whether to disable rollback on stack creation failure..
	* `drift_detection_time` - Drift DetectionTime.
	* `id` - The ID of the Stack.
	* `parent_stack_id` - Parent Stack Id.
	* `ram_role_name` - The RamRoleName.
	* `root_stack_id` - Root Stack Id.
	* `stack_drift_status` - Stack DriftStatus.
	* `stack_id` - Stack Id.
	* `stack_name` - Stack Name.
	* `stack_policy_body` - The structure that contains the stack policy body.
	* `status_reason` - Status Reason.
	* `template_description` - Template Description.
	* `timeout_in_minutes` - Specifies whether to use the values that were passed last time for the parameters that you do not specify in the current request.
	* `parameters` - The parameters.
		* `parameter_key` - The key of parameters.
		* `parameter_value` - The value of parameters.
