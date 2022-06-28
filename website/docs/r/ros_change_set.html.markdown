---
subcategory: "ROS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ros_change_set"
sidebar_current: "docs-alicloud-resource-ros-change-set"
description: |-
  Provides a Alicloud ROS Change Set resource.
---

# alicloud\_ros\_change\_set

Provides a ROS Change Set resource.

For information about ROS Change Set and how to use it, see [What is Change Set](https://www.alibabacloud.com/help/doc-detail/131051.htm).

-> **NOTE:** Available in v1.105.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ros_change_set" "example" {
  change_set_name = "example_value"
  stack_name      = "tf-testacc"
  change_set_type = "CREATE"
  description     = "Test From Terraform"
  template_body   = "{\"ROSTemplateFormatVersion\":\"2015-09-01\"}"
}

```

## Argument Reference

The following arguments are supported:

* `change_set_name` - (Required, ForceNew) The name of the change set.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
* `change_set_type` - (Optional, ForceNew) The type of the change set. Valid values:  CREATE: creates a change set for a new stack. UPDATE: creates a change set for an existing stack. IMPORT: creates a change set for a new stack or an existing stack to import non-ROS-managed resources. If you create a change set for a new stack, ROS creates a stack that has a unique stack ID. The stack is in the REVIEW_IN_PROGRESS state until you execute the change set.  You cannot use the UPDATE type to create a change set for a new stack or the CREATE type to create a change set for an existing stack.
* `description` - (Optional, ForceNew) The description of the change set. The description can be up to 1,024 bytes in length.
* `disable_rollback` - (Optional, ForceNew) Specifies whether to disable rollback on stack creation failure. Default value: false.  Valid values:  true: disables rollback on stack creation failure. false: enables rollback on stack creation failure. Note This parameter takes effect only when ChangeSetType is set to CREATE or IMPORT.
* `notification_urls` - (Optional, ForceNew) The notification urls.
* `parameters` - (Optional, ForceNew) Parameters.
* `ram_role_name` - (Optional, ForceNew) The ram role name.
* `replacement_option` - (Optional, ForceNew) The replacement option.
* `stack_id` - (Optional, ForceNew) The ID of the stack for which you want to create the change set. ROS generates the change set by comparing the stack information with the information that you submit, such as a modified template or different inputs.
* `stack_name` - (Optional, ForceNew) The name of the stack for which you want to create the change set.  The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.  Note This parameter takes effect only when ChangeSetType is set to CREATE or IMPORT.
* `stack_policy_body` - (Optional, ForceNew) The stack policy body.
* `stack_policy_during_update_body` - (Optional, ForceNew) The stack policy during update body.
* `stack_policy_during_update_url` - (Optional, ForceNew) The stack policy during update url.
* `stack_policy_url` - (Optional, ForceNew) The stack policy url.
* `template_body` - (Optional, ForceNew) The structure that contains the template body. The template body must be 1 to 524,288 bytes in length.  If the length of the template body is longer than required, we recommend that you add parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs.  You can specify one of TemplateBody or TemplateURL parameters, but you cannot specify both of them.
* `template_url` - (Optional, ForceNew) The template url.
* `timeout_in_minutes` - (Optional, ForceNew) Timeout In Minutes.
* `use_previous_parameters` - (Optional, ForceNew) The use previous parameters.

#### Block parameters

The parameters supports the following: 

* `parameter_key` - (Required) The parameter key.
* `parameter_value` - (Required) The parameter value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Change Set. Value as `change_set_id`.
* `status` - The status of the change set.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when creating the ROS ChangeSet (until it reaches the initial `CREATE_COMPLETE` status). 

## Import

ROS Change Set can be imported using the id, e.g.

```
$ terraform import alicloud_ros_change_set.example <change_set_id>
```
