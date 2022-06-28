---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_change_sets"
sidebar_current: "docs-alicloud-datasource-ros-change-sets"
description: |-
  Provides a list of Ros Change Sets to the user.
---

# alicloud\_ros\_change\_sets

This data source provides the Ros Change Sets of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.105.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ros_change_sets" "example" {
  stack_id   = "example_value"
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_ros_change_set_id" {
  value = data.alicloud_ros_change_sets.example.sets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `change_set_name` - (Optional, ForceNew) The name of the change set.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Change Set IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Change Set name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `stack_id` - (Required, ForceNew) The ID of the stack for which you want to create the change set. ROS generates the change set by comparing the stack information with the information that you submit, such as a modified template or different inputs.
* `status` - (Optional, ForceNew) The status of the change set. Valid Value: `CREATE_COMPLETE`, `CREATE_FAILED`, `CREATE_IN_PROGRESS`, `CREATE_PENDING`, `DELETE_COMPLETE` and `DELETE_FAILED`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Change Set names.
* `sets` - A list of Ros Change Sets. Each element contains the following attributes:
	* `change_set_id` - The ID of the change set.
	* `change_set_name` - The name of the change set.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
	* `change_set_type` - The type of the change set. Valid values:  CREATE: creates a change set for a new stack. UPDATE: creates a change set for an existing stack. IMPORT: creates a change set for a new stack or an existing stack to import non-ROS-managed resources. If you create a change set for a new stack, ROS creates a stack that has a unique stack ID. The stack is in the REVIEW_IN_PROGRESS state until you execute the change set.  You cannot use the UPDATE type to create a change set for a new stack or the CREATE type to create a change set for an existing stack.
	* `description` - The description of the change set. The description can be up to 1,024 bytes in length.
	* `disable_rollback` - Specifies whether to disable rollback on stack creation failure. Default value: false.  Valid values:  true: disables rollback on stack creation failure. false: enables rollback on stack creation failure. Note This parameter takes effect only when ChangeSetType is set to CREATE or IMPORT.
	* `execution_status` - The execution status of change set N. Maximum value of N: 5. Valid values:  UNAVAILABLE AVAILABLE EXECUTE_IN_PROGRESS EXECUTE_COMPLETE EXECUTE_FAILED OBSOLETE.
	* `id` - The ID of the Change Set.
	* `parameters` - Parameters.
		* `parameter_key` - The parameters.
		* `parameter_value` - The parameters.
	* `stack_id` - The ID of the stack for which you want to create the change set. ROS generates the change set by comparing the stack information with the information that you submit, such as a modified template or different inputs.
	* `stack_name` - The name of the stack for which you want to create the change set.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.  Note This parameter takes effect only when ChangeSetType is set to CREATE or IMPORT.
	* `status` - The status of the change set.
	* `template_body` - The structure that contains the template body. The template body must be 1 to 524,288 bytes in length.  If the length of the template body is longer than required, we recommend that you add parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs.  You can specify one of TemplateBody or TemplateURL parameters, but you cannot specify both of them.
	* `timeout_in_minutes` - Timeout In Minutes.
