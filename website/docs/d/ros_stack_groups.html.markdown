---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_stack_groups"
sidebar_current: "docs-alicloud-datasource-ros-stack-groups"
description: |-
  Provides a list of Ros Stack Groups to the user.
---

# alicloud\_ros\_stack\_groups

This data source provides the Ros Stack Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.107.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ros_stack_groups" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_ros_stack_group_id" {
  value = data.alicloud_ros_stack_groups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Stack Group IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Stack Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of Stack Group. Valid Values: `ACTIVE`, `DELETED`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Stack Group names.
* `groups` - A list of Ros Stack Groups. Each element contains the following attributes:
	* `administration_role_name` - The name of the RAM administrator role assumed by ROS.
	* `description` - The description of the stack group.
	* `execution_role_name` - The name of the RAM execution role assumed by the administrator role.
	* `id` - The ID of the Stack Group.
	* `parameters` - The parameters.
		* `parameter_key` - The parameter key.
		* `parameter_value` - The parameter value.
	* `stack_group_id` - The id of Stack Group.
	* `stack_group_name` - The name of the stack group..
	* `status` - The status of Stack Group.
	* `template_body` - The structure that contains the template body.
