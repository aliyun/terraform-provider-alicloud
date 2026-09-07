---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_control_policies"
sidebar_current: "docs-alicloud-datasource-resource-manager-control-policies"
description: |-
  Provides a list of Resource Manager Control Policies to the user.
---

# alicloud\_resource\_manager\_control\_policies

This data source provides the Resource Manager Control Policies of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_control_policies" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_resource_manager_control_policy_id" {
  value = data.alicloud_resource_manager_control_policies.example.policies.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Control Policy IDs.
* `language` - (Optional, ForceNew) The language. Valid value `zh-CN`, `en`, and `ja`. Default value `zh-CN`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Control Policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `policy_type` - (Optional, ForceNew) The policy type of control policy. Valid values `System` and `Custom`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Control Policy names.
* `policies` - A list of Resource Manager Control Policies. Each element contains the following attributes:
	* `attachment_count` - The count of policy attachment.
	* `control_policy_name` - The name of policy.
	* `description` - The description of policy.
	* `effect_scope` - The effect scope.
	* `id` - The ID of the Control Policy.
	* `policy_document` - The policy document.
	* `policy_id` - The ID of policy.
	* `policy_type` - The type of policy.
