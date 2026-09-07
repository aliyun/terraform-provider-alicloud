---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_control_policy_attachments"
sidebar_current: "docs-alicloud-datasource-resource-manager-control-policy-attachments"
description: |-
  Provides a list of Resource Manager Control Policy Attachments to the user.
---

# alicloud\_resource\_manager\_control\_policy\_attachments

This data source provides the Resource Manager Control Policy Attachments of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_control_policy_attachments" "example" {
  target_id = "example_value"
}

output "first_resource_manager_control_policy_attachment_id" {
  value = data.alicloud_resource_manager_control_policy_attachments.example.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:

* `language` - (Optional, ForceNew) The language. Valid value `zh-CN`, `en`, and `ja`. Default value `zh-CN`
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `target_id` - (Required, ForceNew) The Id of target.
* `policy_type` - (Optional, ForceNew) The policy type of control policy. Valid values: `Custom` and `System`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Control Policy Attachment IDs.
* `attachments` - A list of Resource Manager Control Policy Attachments. Each element contains the following attributes:
	* `attach_date` - The attach date.
	* `description` - The description of policy.
	* `id` - The ID of the Control Policy Attachment.
	* `policy_id` - The ID of policy.
	* `policy_name` - The name of policy.
	* `policy_type` - The type of policy.
