---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_stack_instance"
sidebar_current: "docs-alicloud-resource-ros-stack-instance"
description: |-
  Provides a Alicloud ROS Stack Instance resource.
---

# alicloud\_ros\_stack\_instance

Provides a ROS Stack Instance resource.

For information about ROS Stack Instance and how to use it, see [What is Stack Instance](https://www.alibabacloud.com/help/en/doc-detail/151338.html).

-> **NOTE:** Available in v1.145.0+.

## Example Usage

Basic Usage

```terraform

data "alicloud_ros_regions" "example" {}

resource "alicloud_ros_stack_group" "example" {
  stack_group_name = var.name
  template_body    = "{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}"
  description      = "test for stack groups"
  parameters {
    parameter_key   = "VpcName"
    parameter_value = "VpcName"
  }
  parameters {
    parameter_key   = "InstanceType"
    parameter_value = "InstanceType"
  }
}

resource "alicloud_ros_stack_instance" "example" {
  stack_group_name          = alicloud_ros_stack_group.example.stack_group_name
  stack_instance_account_id = "example_value"
  stack_instance_region_id  = data.alicloud_ros_regions.example.regions.0.region_id
  operation_preferences     = "{\"FailureToleranceCount\": 1, \"MaxConcurrentCount\": 2}"
  parameter_overrides {
    parameter_value = "VpcName"
    parameter_key   = "VpcName"
  }
}

```

## Argument Reference

The following arguments are supported:

* `operation_description` - (Optional) The operation description.
* `operation_preferences` - (Optional) The operation preferences. The operation settings. The following fields are supported:
  * `FailureToleranceCount` The maximum number of stack group operation failures that can occur. In a stack group operation, if the total number of failures does not exceed the FailureToleranceCount value, the operation succeeds. Otherwise, the operation fails. If the FailureToleranceCount parameter is not specified, the default value 0 is used. You cannot specify both FailureToleranceCount and FailureTolerancePercentage. Valid values: `0` to `20`. 
  * `FailureTolerancePercentage`: The percentage of stack group operation failures that can occur. In a stack group operation, if the percentage of failures does not exceed the FailureTolerancePercentage value, the operation succeeds. Otherwise, the operation fails. You cannot specify both FailureToleranceCount and FailureTolerancePercentage. Valid values: `0` to `100`. 
  * `MaxConcurrentCount`: The maximum number of accounts within which to perform this operation at one time. You cannot specify both MaxConcurrentCount and MaxConcurrentPercentage. Valid values: `1` to `20`. 
  * `MaxConcurrentPercentage`: The maximum percentage of accounts within which to perform this operation at one time. You cannot specify both MaxConcurrentCount and MaxConcurrentPercentage. Valid values: `1` to `100`
* `parameter_overrides` - (Optional, Sensitive) ParameterOverrides. See the following `Block parameter_overrides`.
* `stack_instance_account_id` - (Required) The account to which the stack instance belongs.
* `stack_instance_region_id` - (Required) The region of the stack instance.
* `stack_group_name` - (Required, ForceNew) The name of the stack group.
* `retain_stacks` - (Optional) Specifies whether to retain the stack corresponding to the stack instance.Default value `false`. **NOTE:** When `retain_stacks` is `true`, the stack is retained. If the stack is retained, the corresponding stack is not deleted when the stack instance is deleted from the stack group. 
* `timeout_in_minutes` - (Optional) The timeout period that is specified for the stack creation request. Default value: `60`. Unit: `minutes`.

#### Block parameter_overrides

The parameter_overrides supports the following: 

* `parameter_key` - (Required, Sensitive) The key of override parameter. If you do not specify the key and value of the parameter, ROS uses the key and value that you specified when you created the stack group.
* `parameter_value` - (Required, Sensitive) The value of override parameter. If you do not specify the key and value of the parameter, ROS uses the key and value that you specified when you created the stack group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Stack Instance. The value formats as `<stack_group_name>:<stack_instance_account_id>:<stack_instance_region_id>`.
* `status` - The status of the stack instance. Valid values: `CURRENT` or `OUTDATED`. 
  * `CURRENT`: The stack corresponding to the stack instance is up to date with the stack group. 
  * `OUTDATED`: The stack corresponding to the stack instance is not up to date with the stack group. The `OUTDATED` state has the following possible causes: 
    * When the CreateStackInstances operation is called to create stack instances, the corresponding stacks fail to be created. 
    * When the UpdateStackInstances or UpdateStackGroup operation is called to update stack instances, the corresponding stacks fail to be updated, or only some of the stack instances are updated. 
    * The create or update operation is not complete.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Stack Instance.

## Import

ROS Stack Instance can be imported using the id, e.g.

```
$ terraform import alicloud_ros_stack_instance.example <stack_group_name>:<stack_instance_account_id>:<stack_instance_region_id>
```