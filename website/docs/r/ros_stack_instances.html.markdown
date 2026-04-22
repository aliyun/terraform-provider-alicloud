---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_stack_instances"
sidebar_current: "docs-alicloud-resource-ros-stack-instances"
description: |-
  Provides a Resource Orchestration Service (ROS) Stack Instances resource to manage stack instances across multiple regions and accounts within a Stack Group.
---

# alicloud_ros_stack_instances

Provides a ROS Stack Instances resource that allows you to deploy, update, or remove stack instances in bulk across multiple target regions and Alibaba Cloud accounts within a Stack Group. It supports both self-managed and service-managed permission models.

For information about ROS Stack Instances and how to use it, see [What is Stack Instance](https://www.alibabacloud.com/help/en/resource-orchestration-service/latest/api-ros-2019-09-10-createstackinstances).

-> **NOTE:** Available since v1.257.0.

-> **Note:** Stack Groups and Stack Instances rely on asynchronous batch operations. Partial failures do not trigger automatic rollback. See the [Important Notes](#important-notes) section for operational guidance.

## Example Usage

Basic Usage - Self-Managed Permissions

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ros_stack_instances&exampleId=basic-self-managed&activeTab=example&spm=docs.r.ros_stack_instances.0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_account" "this" {
}

data "alicloud_ros_regions" "default" {}

resource "alicloud_ros_stack_group" "default" {
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

resource "alicloud_ros_stack_instances" "self_managed" {
  stack_group_name = alicloud_ros_stack_group.default.stack_group_name
  region_ids       = [data.alicloud_ros_regions.default.regions.0.region_id]
  account_ids      = [data.alicloud_account.this.id]

  parameter_overrides {
    parameter_value = "VpcName"
    parameter_key   = "VpcName"
  }
  timeout_in_minutes    = 45
  operation_description = "Batch deployment for production environment"
  disable_rollback      = false
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://github.com/aliyun/terraform-provider-alicloud/tree/master/examples/ros-stack-instances)

## Argument Reference

The following arguments are supported:

* `stack_group_name` - (Required, ForceNew) The name of the stack group to which the stack instances belong. This parameter cannot be modified after creation.

* `region_ids` - (Required, ForceNew) List of target region IDs where stack instances will be deployed. You can specify 1 to 20 regions. This parameter cannot be modified after creation. Example: `["cn-beijing", "cn-shanghai"]`.

* `account_ids` - (Optional, ForceNew) List of target Alibaba Cloud account IDs for self-managed permissions model. You can specify 1 to 50 accounts. This parameter conflicts with `deployment_targets`. This parameter cannot be modified after creation. Example: `["123456789012****", "098765432109****"]`.

* `deployment_targets` - (Optional, ForceNew) Configuration block defining deployment targets for service-managed permissions model. This parameter conflicts with `account_ids`. This parameter cannot be modified after creation. See [`deployment_targets`](#deployment_targets) below.

-> **NOTE:** You must specify either `account_ids` (for self-managed permissions) or `deployment_targets` (for service-managed permissions), but not both.

* `parameter_overrides` - (Optional) A set of parameters to override in the stack instances. See [`parameter_overrides`](#parameter_overrides) below.

* `operation_preferences` - (Optional) Configuration block defining preferences for how the operation is performed across multiple accounts and regions. See [`operation_preferences`](#operation_preferences) below.

* `timeout_in_minutes` - (Optional) The amount of time in minutes that can elapse before the stack operation status is set to `TIMED_OUT`. Valid values: 1 to 1440. Default value: 60.

* `operation_description` - (Optional) Description of the stack instances operation. The description must be 1 to 256 characters in length.

* `disable_rollback` - (Optional, ForceNew) Specifies whether to disable the rollback policy when creating stack instances fails. Valid values: `true`, `false`. Default value: `false`. This parameter cannot be modified after creation.

* `deployment_options` - (Optional, ForceNew) List of deployment options for service-managed permissions. Currently only supports `IgnoreExisting`, which skips existing stack instances during deployment. This parameter cannot be modified after creation. Example: `["IgnoreExisting"]`.

### `deployment_targets`

The `deployment_targets` block supports the following:

* `account_ids` - (Optional) List of Alibaba Cloud account IDs for service-managed permissions. Maximum 50 accounts.
* `rd_folder_ids` - (Optional) List of Resource Directory folder IDs. Maximum 20 folders.

### `parameter_overrides`

The `parameter_overrides` block supports the following:

* `parameter_key` - (Required) The key of the parameter to override.
* `parameter_value` - (Optional) The value of the parameter to override. This field is sensitive and will be masked in logs.

### `operation_preferences`

The `operation_preferences` block supports the following:

* `max_concurrent_count` - (Optional) Maximum number of concurrent operations per region. Valid values: 1 to 20. Conflicts with `max_concurrent_percentage`.
* `max_concurrent_percentage` - (Optional) Maximum percentage of concurrent targets per region. Valid values: 1 to 100. Conflicts with `max_concurrent_count`.
* `failure_tolerance_count` - (Optional) Number of failures tolerated per region. Valid values: 0 to 20. Conflicts with `failure_tolerance_percentage`.
* `failure_tolerance_percentage` - (Optional) Percentage of failures tolerated per region. Valid values: 0 to 100. Conflicts with `failure_tolerance_count`.
* `region_concurrency_type` - (Optional) Concurrency type for regions. Valid values: `SEQUENTIAL`, `PARALLEL`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which is the same as `stack_group_name`.

* `stack_instances` - A list of stack instances with their latest operation tracking information. See [`stack_instances`](#stack_instances) below.

-> **NOTE:** The `stack_instances` attribute may be empty if the deployment targets contain no valid accounts or if all deployments failed.

### `stack_instances`

The `stack_instances` block exports the following:

* `last_operation_id` - The ID of the last operation performed on this stack instance.
* `status` - The status of the stack instance. Valid values: `CURRENT`, `OUTDATED`, `INOPERABLE`, `RUNNING`, `FAILED`, `SUCCEEDED`, etc.
* `stack_group_id` - The ID of the stack group to which this instance belongs.
* `stack_id` - The ID of the underlying stack.
* `drift_detection_time` - The timestamp when drift detection was last performed.
* `stack_drift_status` - The drift status of the stack. Valid values: `NOT_CHECKED`, `IN_SYNC`, `DRIFTED`, `CHECK_FAILED`, etc.
* `status_reason` - The reason for the current status of the stack instance.
* `stack_group_name` - The name of the stack group.
* `account_id` - The Alibaba Cloud account ID where the stack instance is deployed.
* `region_id` - The region ID where the stack instance is deployed.
* `rd_folder_id` - The Resource Directory folder ID (if applicable).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the Stack Instances.
* `update` - (Defaults to 60 mins) Used when updating the Stack Instances.
* `delete` - (Defaults to 60 mins) Used when deleting the Stack Instances.

## Import

ROS Stack Instances cannot be imported.

## Important Notes

1. **Asynchronous Operations**: Stack instance operations (create, update, delete) are asynchronous batch operations. The provider waits for the operation to complete, but partial failures may occur where some instances succeed while others fail.

2. **No Automatic Rollback**: If some stack instances fail to create, successful instances are NOT automatically rolled back. You need to manually handle failed instances or retry the operation.

3. **Operation Conflicts**: Only one stack group operation can run at a time. If you attempt to perform another operation while one is in progress, you will receive a `StackGroupOperationInProgress` error. The provider will automatically retry in this case.

4. **State Management**: If the create operation encounters partial failures, the resource ID will NOT be set in Terraform state. This allows Terraform to retry the entire operation on the next apply. This is intentional to prevent orphaned resources.

5. **Parameter Sensitivity**: The `parameter_value` field in `parameter_overrides` is marked as sensitive. Values will be masked in Terraform output and logs for security.

6. **ForceNew Parameters**: Most parameters (`stack_group_name`, `region_ids`, `account_ids`, `deployment_targets`, `disable_rollback`, `deployment_options`) require resource recreation if modified. Only `parameter_overrides`, `operation_preferences`, `timeout_in_minutes`, and `operation_description` support in-place updates.

7. **Empty Results**: If your deployment targets result in no stack instances being created (e.g., targeting an empty folder), the `stack_instances` attribute will be an empty list. This is expected behavior and does not indicate an error.
