---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_stack_group"
sidebar_current: "docs-alicloud-resource-ros-stack-group"
description: |-
  Provides a Alicloud ROS Stack Group resource.
---

# alicloud_ros_stack_group

Provides a ROS Stack Group resource.

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

* `account_ids` - (Optional) The list of target account IDs, in JSON format. A maximum of 20 accounts can be specified.
* `administration_role_name` - (Optional) The name of the RAM administrator role assumed by ROS. ROS assumes this role to perform operations on the stack corresponding to the stack instance in the stack group.
* `description` - (Optional) The description of the stack group.
* `execution_role_name` - (Optional) The name of the RAM execution role assumed by the administrator role. ROS assumes this role to perform operations on the stack corresponding to the stack instance in the stack group.
* `operation_description` - (Optional) The description of the operation.
* `operation_preferences` - (Optional) The operation settings, in JSON format.
* `parameters` - (Optional) The parameters. If the parameter name and value are not specified, ROS will use the default value specified in the template.
* `region_ids` - (Optional) The list of target regions, in JSON format. A maximum of 20 accounts can be specified.
* `stack_group_name` - (Required, ForceNew) The name of the stack group. The name must be unique in a region.
* `template_body` - (Optional) The structure that contains the template body. The template body must be 1 to 524,288 bytes in length. If the length of the template body is longer than required, we recommend that you add parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs.
* `template_url` - (Optional) The URL of the file that contains the template body. The URL must point to a template located in an HTTP or HTTPS web server or an Alibaba Cloud OSS bucket. Examples: oss://ros/template/demo and oss://ros/template/demo?RegionId=cn-hangzhou. The template must be 1 to 524,288 bytes in length. If the region of the OSS bucket is not specified, the RegionId value is used by default.
* `template_version` - (Optional) The version of the template.

#### Block parameters

The parameters supports the following: 

* `parameter_key` - (Required) The parameter key.
* `parameter_value` - (Required) The parameter value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Stack Group. Value as `stack_group_name`.
* `stack_group_id` - The id of Stack Group.
* `status` - The status of Stack Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Stack Group.
* `update` - (Defaults to 6 mins) Used when update the Stack Group.

## Import

ROS Stack Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ros_stack_group.example <stack_group_name>
```
