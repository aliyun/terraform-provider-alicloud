---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_stack_instances"
sidebar_current: "docs-alicloud-datasource-ros-stack-instances"
description: |-
  Provides a list of Ros Stack Instances to the user.
---

# alicloud\_ros\_stack\_instances

This data source provides the Ros Stack Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.145.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ros_stack_instances" "ids" {
  stack_group_name = "example_value"
  ids              = ["example_value-1", "example_value-2"]
  enable_details   = true
}
output "ros_stack_instance_id_1" {
  value = data.alicloud_ros_stack_instances.ids.instances.0.id
}

data "alicloud_ros_stack_instances" "status" {
  stack_group_name = "example_value"
  status           = "CURRENT"
  enable_details   = true
}
output "ros_stack_instance_id_2" {
  value = data.alicloud_ros_stack_instances.status.instances.0.id
}

data "alicloud_ros_stack_instances" "regionId" {
  stack_group_name         = "example_value"
  stack_instance_region_id = "example_value"
  enable_details           = true
}
output "ros_stack_instance_id_3" {
  value = data.alicloud_ros_stack_instances.regionId.instances.0.id
}

data "alicloud_ros_stack_instances" "accountId" {
  stack_group_name          = "example_value"
  stack_instance_account_id = "example_value"
  enable_details            = true
}
output "ros_stack_instance_id_4" {
  value = data.alicloud_ros_stack_instances.accountId.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Stack Instance IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `stack_group_name` - (Required, ForceNew) The name of the stack group.
* `stack_instance_account_id` - (Optional, ForceNew) The account to which the stack instance belongs.
* `stack_instance_region_id` - (Optional, ForceNew) The region of the stack instance.
* `status` - (Optional, ForceNew) The status of the stack instance. Valid values: `CURRENT` or `OUTDATED`. 
  * `CURRENT`: The stack corresponding to the stack instance is up to date with the stack group. 
  * `OUTDATED`: The stack corresponding to the stack instance is not up to date with the stack group. The `OUTDATED` state has the following possible causes: 
    * When the CreateStackInstances operation is called to create stack instances, the corresponding stacks fail to be created. 
    * When the UpdateStackInstances or UpdateStackGroup operation is called to update stack instances, the corresponding stacks fail to be updated, or only some of the stack instances are updated. 
    * The create or update operation is not complete.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Ros Stack Instances. Each element contains the following attributes:
  * `id` - The ID of the Stack Instance. The value formats as `<stack_group_name>:<stack_instance_account_id>:<stack_instance_region_id>`.
  * `parameter_overrides` - ParameterOverrides.
    * `parameter_key` - The key of override parameter.
    * `parameter_value` - The value of override parameter.
  * `stack_group_id` - The ID of the stack group.
  * `stack_group_name` - The name of the stack group.
  * `stack_id` - The ID of the stack corresponding to the stack instance.
  * `stack_instance_account_id` - The account to which the stack instance belongs.
  * `stack_instance_region_id` - The region of the stack instance.
  * `status` - The status of the stack instance. Valid values: `CURRENT` or `OUTDATED`. 
    * `CURRENT`: The stack corresponding to the stack instance is up to date with the stack group. 
    * `OUTDATED`: The stack corresponding to the stack instance is not up to date with the stack group. The `OUTDATED` state has the following possible causes: 
      * When the CreateStackInstances operation is called to create stack instances, the corresponding stacks fail to be created. 
      * When the UpdateStackInstances or UpdateStackGroup operation is called to update stack instances, the corresponding stacks fail to be updated, or only some of the stack instances are updated. 
      * The create or update operation is not complete.
  * `status_reason` - The reason why the stack is in its current state.