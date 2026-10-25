---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_change_set"
sidebar_current: "docs-alicloud-resource-ros-change-set"
description: |-
  Provides a Alicloud ROS Change Set resource.
---

# alicloud_ros_change_set

Provides a ROS Change Set resource.

For information about ROS Change Set and how to use it, see [What is Change Set](https://www.alibabacloud.com/help/doc-detail/131051.htm).

-> **NOTE:** Available since v1.105.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ros_change_set&exampleId=1c4bfc18-d876-86f4-2360-4f0b62444fa17265b398&activeTab=example&spm=docs.r.ros_change_set.0.1c4bfc18d8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ros_change_set" "example" {
  change_set_name = "terraform-example-${random_integer.default.result}"
  stack_name      = "terraform-example-${random_integer.default.result}"
  change_set_type = "CREATE"
  description     = "Test From Terraform"
  template_body   = "{\"ROSTemplateFormatVersion\":\"2015-09-01\"}"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ros_change_set&spm=docs.r.ros_change_set.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `change_set_name` - (Required, ForceNew) The name of the change set.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
* `change_set_type` - (Optional, ForceNew) The type of the change set. Valid values:  CREATE: creates a change set for a new stack. UPDATE: creates a change set for an existing stack. IMPORT: creates a change set for a new stack or an existing stack to import non-ROS-managed resources. If you create a change set for a new stack, ROS creates a stack that has a unique stack ID. The stack is in the REVIEW_IN_PROGRESS state until you execute the change set.  You cannot use the UPDATE type to create a change set for a new stack or the CREATE type to create a change set for an existing stack.
* `description` - (Optional, ForceNew) The description of the change set. The description can be up to 1,024 bytes in length.
* `disable_rollback` - (Optional, ForceNew) Specifies whether to disable rollback on stack creation failure. Default value: false.  Valid values:  true: disables rollback on stack creation failure. false: enables rollback on stack creation failure. Note This parameter takes effect only when ChangeSetType is set to CREATE or IMPORT.
* `notification_urls` - (Optional, ForceNew) The notification urls.
* `parallelism` - (Optional, ForceNew, Available since v2.1.0) The maximum number of concurrent resources of the Terraform stack. You can enter an integer greater than or equal to zero. This parameter takes effect only when ChangeSetType is set to CREATE or UPDATE.
* `parameters` - (Optional, ForceNew) Parameters. See [`parameters`](#parameters) below.
* `ram_role_name` - (Optional, ForceNew) The ram role name.
* `replacement_option` - (Optional, ForceNew) The replacement option.
* `resource_group_id` - (Optional, ForceNew, Available since v2.1.0) The ID of the resource group. If you do not specify this parameter, the stack corresponding to the change set joins the default resource group.
* `resources_to_import` - (Optional, ForceNew, Available since v2.1.0) The list of resources to be imported. This parameter takes effect only when ChangeSetType is set to IMPORT. See [`resources_to_import`](#resources_to_import) below.
* `stack_id` - (Optional, ForceNew) The ID of the stack for which you want to create the change set. ROS generates the change set by comparing the stack information with the information that you submit, such as a modified template or different inputs.
* `stack_name` - (Optional, ForceNew) The name of the stack for which you want to create the change set.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.  Note This parameter takes effect only when ChangeSetType is set to CREATE or IMPORT.
* `stack_policy_body` - (Optional, ForceNew) The stack policy body. Only one of the four stack policy parameters can be set per request; available for CREATE and UPDATE change sets.
* `stack_policy_during_update_body` - (Optional, ForceNew) The stack policy during update body. Only one of the four stack policy parameters can be set per request; available for UPDATE change sets only.
* `stack_policy_during_update_url` - (Optional, ForceNew) The stack policy during update url. Only one of the four stack policy parameters can be set per request; available for UPDATE change sets only. The URL must be fetchable anonymously.
* `stack_policy_url` - (Optional, ForceNew) The stack policy url. Only one of the four stack policy parameters can be set per request; available for CREATE and UPDATE change sets. The URL must be fetchable anonymously.
* `tags` - (Optional, ForceNew, Available since v2.1.0) The tags of the change set. The tags are passed to each resource that supports tags in the stack. Maximum number of tags: 20.
* `taint_resources` - (Optional, ForceNew, Available since v2.1.0) The list of tainted resources. For ROS stacks, configure this parameter as the resource name, such as `my_vpc`. For Terraform stacks, configure this parameter as the resource type and resource name, such as `alicloud_vpc.my_vpc`.
* `template_body` - (Optional, ForceNew) The structure that contains the template body. The template body must be 1 to 524,288 bytes in length.  If the length of the template body is longer than required, we recommend that you add parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs.  You can specify one of TemplateBody or TemplateURL parameters, but you cannot specify both of them.
* `template_id` - (Optional, ForceNew, Available since v2.1.0) The ID of the template. Shared templates and private templates are supported. You can specify only one of TemplateBody, TemplateURL, or TemplateId.
* `template_scratch_id` - (Optional, ForceNew, Available since v2.1.0) The ID of the scene template. Only supported when `change_set_type` is `IMPORT` (new stack import).
* `template_url` - (Optional, ForceNew) The template url.
* `template_version` - (Optional, ForceNew, Available since v2.1.0) The template version. This parameter takes effect only when TemplateId is specified.
* `timeout_in_minutes` - (Optional, ForceNew) Timeout In Minutes.
* `use_previous_parameters` - (Optional, ForceNew) The use previous parameters.

### `parameters`

The parameters supports the following:

* `parameter_key` - (Required) The parameter key.
* `parameter_value` - (Required) The parameter value.

### `resources_to_import`

The resources_to_import supports the following:

* `resource_identifier` - (Optional, ForceNew, Available since v2.1.0) The key-value mapping that identifies the resource to be imported. The key is the identifier attribute of the resource (for example, `VpcId` for `ALIYUN::ECS::VPC`), and the value is the resource ID. This parameter takes effect only when ChangeSetType is set to IMPORT.
* `logical_resource_id` - (Optional, ForceNew, Available since v2.1.0) The logical ID of the resource in the template. This parameter takes effect only when ChangeSetType is set to IMPORT.
* `resource_type` - (Optional, ForceNew, Available since v2.1.0) The type of the resource, which must be consistent with the type defined in the template. This parameter takes effect only when ChangeSetType is set to IMPORT.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Change Set. Value as `change_set_id`.
* `status` - The status of the change set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when creating the ROS ChangeSet (until it reaches the initial `CREATE_COMPLETE` status). 

## Import

ROS Change Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_ros_change_set.example <change_set_id>
```
