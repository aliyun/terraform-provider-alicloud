---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_stack_group"
description: |-
  Provides a Alicloud ROS Stack Group resource.
---

# alicloud_ros_stack_group

Provides a ROS Stack Group resource.

Resource stack Group.

For information about ROS Stack Group and how to use it, see [What is Stack Group](https://www.alibabacloud.com/help/en/doc-detail/151333.htm).

-> **NOTE:** Available since v1.107.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ros_stack_group&exampleId=63f69414-5509-07dd-bf0d-9e271f7c5cdb5bbd0242&activeTab=example&spm=docs.r.ros_stack_group.0.63f6941455&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ros_stack_group" "example" {
  stack_group_name = "example_value"
  template_body    = <<EOF
    {
    	"ROSTemplateFormatVersion": "2015-09-01"
    }
    EOF
}

```

## Argument Reference

The following arguments are supported:
* `administration_role_name` - (Optional, Computed) The name of the RAM role that you specify for the administrator account in ROS when you create the self-managed stack group. If you do not specify this parameter, the default value AliyunROSStackGroupAdministrationRole is used. You can use the administrator role in ROS to assume the execution role AliyunROSStackGroupExecutionRole to perform operations on the stacks that correspond to stack instances in the stack group.
* `auto_deployment` - (Optional, List, Available since v1.257.0) Automatic deployment setting information. Description
This parameter is required only if the PermissionModel is SERVICE_MANAGED. See [`auto_deployment`](#auto_deployment) below.
* `capabilities` - (Optional, List, Available since v1.257.0) The list of resource stack group options. The maximum length is 1.
* `description` - (Optional) The description of the stack group.
* `execution_role_name` - (Optional, Computed) The name of the RAM role that you specify for the execution account when you create the self-managed stack group. You can use the administrator role AliyunROSStackGroupAdministrationRole to assume the execution role. If you do not specify this parameter, the default value AliyunROSStackGroupExecutionRole is used. You can use this role in ROS to perform operations on the stacks that correspond to stack instances in the stack group.
* `parameters` - (Optional, List) Parameters See [`parameters`](#parameters) below.
* `permission_model` - (Optional, Available since v1.257.0) The permission model.
* `resource_group_id` - (Optional, Computed, Available since v1.257.0) The ID of the resource group.
* `stack_group_name` - (Required, ForceNew) StackGroupName
* `tags` - (Optional, Map, Available since v1.257.0) The label of the resource stack group.
* `template_body` - (Optional) The template body.
* `template_id` - (Optional, Available since v1.257.0) The ID of the template.
* `template_url` - (Optional) The location of the file that contains the template body. The URL must point to the template (1 to 524,288 bytes) located in the HTTP Web server (HTTP or HTTPS) or Alibaba Cloud OSS bucket. The URL of the OSS bucket, such as oss:// ros/template/demo or oss:// ros/template/demo? RegionId = cn-hangzhou. If the OSS region is not specified, the RegionId of the interface is the same by default.

-> **NOTE:** You must and can specify only one of the parameters of TemplateBody, TemplateURL, or TemplateId.

* `template_version` - (Optional) The version of the template.

### `auto_deployment`

The auto_deployment supports the following:
* `enabled` - (Optional, Available since v1.257.0) Enable or disable automatic deployment. Valid Values:
  - `true`: Enable automatic deployment. After automatic deployment is enabled, if a member account is added to the target folder, the stack group will automatically deploy the stack instance to the account. If a member account is deleted from the target folder, the stack group automatically deletes the stack instance in the account.
  - `false`: disable automatic deployment. After automatic deployment is disabled, the stack instance will not change when the member accounts in the target folder change.
* `retain_stacks_on_account_removal` - (Optional, Available since v1.257.0) Whether to retain the stack in the member account when the member account is deleted from the target folder. Valid values:
  - `true`: Keep the resource stack.
  - `false`: delete the resource stack.

-> **NOTE:** When Enabled is true, `retain_stacks_on_account_removal` is required.

* `account_ids` - (Removed since v1.257.0). Field 'account_ids' has been deprecated from provider version 1.257.0.
* `region_ids` - (Removed since v1.257.0). Field 'region_ids' has been deprecated from provider version 1.257.0.
* `operation_description` - (Removed since v1.257.0). Field 'operation_description' has been deprecated from provider version 1.257.0. You should use resource alicloud_ros_stack_instance's field 'operation_description'.
* `operation_preferences` - (Removed since v1.257.0). Field 'operation_preferences' has been deprecated from provider version 1.257.0. You should use resource alicloud_ros_stack_instance's field 'operation_preferences'.

### `parameters`

The parameters supports the following:
* `parameter_key` - (Required) The key of parameter N. If you do not specify the key and value of the parameter, ROS uses the default key and value in the template.
* `parameter_value` - (Required) The value of parameter N.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `stack_group_id` - The ID of stack group.
* `status` - The status of the stack group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Stack Group.
* `delete` - (Defaults to 5 mins) Used when delete the Stack Group.
* `update` - (Defaults to 5 mins) Used when update the Stack Group.

## Import

ROS Stack Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ros_stack_group.example <id>
```