---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_stack"
sidebar_current: "docs-alicloud-resource-ros-stack"
description: |-
  Provides a Alicloud ROS Stack resource.
---

# alicloud_ros_stack

Provides a ROS Stack resource.

For information about ROS Stack and how to use it, see [What is Stack](https://www.alibabacloud.com/help/en/doc-detail/132086.htm).

-> **NOTE:** Available since v1.106.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ros_stack&exampleId=aa17e5cf-fef0-e575-c30a-d8ae978bcb35be89d3e1&activeTab=example&spm=docs.r.ros_stack.0.aa17e5cffe&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ros_stack" "example" {
  stack_name        = "tf-testaccstack"
  template_body     = <<EOF
    {
    	"ROSTemplateFormatVersion": "2015-09-01"
    }
    EOF
  stack_policy_body = <<EOF
    {
    	"Statement": [{
    		"Action": "Update:Delete",
    		"Resource": "*",
    		"Effect": "Allow",
    		"Principal": "*"
    	}]
    }
    EOF
}


```

## Argument Reference

The following arguments are supported:

* `create_option` - (Optional, ForceNew) Specifies whether to delete the stack after it is created.
* `deletion_protection` - (Optional, ForceNew) Specifies whether to enable deletion protection on the stack. Valid values: `Disabled`, `Enabled`. Default to: `Disabled`
* `disable_rollback` - (Optional) Specifies whether to disable rollback on stack creation failure. Default to: `false`.
* `notification_urls` - (Optional, ForceNew) The callback URL for receiving stack event N. Only HTTP POST is supported. Maximum value of N: 5.
* `ram_role_name` - (Optional) The name of the RAM role. ROS assumes the specified RAM role to create the stack and call API operations by using the credentials of the role.
* `replacement_option` - (Optional) Specifies whether to enable replacement update after a resource attribute that does not support modification update is changed. Modification update keeps the physical ID of the resource unchanged. However, the resource is deleted and then recreated, and its physical ID is changed if replacement update is enabled.
* `retain_all_resources` - (Optional) The retain all resources.
* `parameters` - (Optional) The parameters. If the parameter name and value are not specified, ROS will use the default value specified in the template.
* `retain_resources` - (Optional) Specifies whether to retain the resources in the stack.
* `stack_name` - (Required, ForceNew) The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
* `stack_policy_body` - (Optional) The structure that contains the stack policy body. The stack policy body must be 1 to 16,384 bytes in length.
* `stack_policy_during_update_body` - (Optional) The structure that contains the body of the temporary overriding stack policy. The stack policy body must be 1 to 16,384 bytes in length.
* `stack_policy_during_update_url` - (Optional) The URL of the file that contains the temporary overriding stack policy. The URL must point to a policy located in an HTTP or HTTPS web server or an Alibaba Cloud OSS bucket. Examples: oss://ros/stack-policy/demo and oss://ros/stack-policy/demo?RegionId=cn-hangzhou. The policy can be up to 16,384 bytes in length and the URL can be up to 1,350 bytes in length. If the region of the OSS bucket is not specified, the RegionId value is used by default.
* `stack_policy_url` - (Optional) The URL of the file that contains the stack policy. The URL must point to a policy located in an HTTP or HTTPS web server or an Alibaba Cloud OSS bucket. Examples: oss://ros/stack-policy/demo and oss://ros/stack-policy/demo?RegionId=cn-hangzhou. The policy can be up to 16,384 bytes in length and the URL can be up to 1,350 bytes in length. If the region of the OSS bucket is not specified, the RegionId value is used by default.
* `template_body` - (Optional) The structure that contains the template body. The template body must be 1 to 524,288 bytes in length. If the length of the template body is longer than required, we recommend that you add parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs.
* `template_url` - (Optional) The URL of the file that contains the template body. The URL must point to a template located in an HTTP or HTTPS web server or an Alibaba Cloud OSS bucket. Examples: oss://ros/template/demo and oss://ros/template/demo?RegionId=cn-hangzhou. The template must be 1 to 524,288 bytes in length. If the region of the OSS bucket is not specified, the RegionId value is used by default.
* `template_version` - (Optional) The version of the template.
* `timeout_in_minutes` - (Optional) The timeout period that is specified for the stack creation request. Default to: `60`.
* `use_previous_parameters` - (Optional) Specifies whether to use the values that were passed last time for the parameters that you do not specify in the current request.
* `tags` - (Optional) A mapping of tags to assign to the resource.

#### Block parameters

The parameters supports the following: 

* `parameter_key` - (Required) The parameter key.
* `parameter_value` - (Required) The parameter value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Stack. Value as `stack_id`.
* `status` - The status of Stack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Stack.
* `delete` - (Defaults to 6 mins) Used when delete the Stack.
* `update` - (Defaults to 11 mins) Used when update the Stack.

## Import

ROS Stack can be imported using the id, e.g.

```shell
$ terraform import alicloud_ros_stack.example <stack_id>
```
