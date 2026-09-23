---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_policy_attachments"
sidebar_current: "docs-alicloud-datasource-resource-manager-policy-attachments"
description: |-
    Provides a list of Resource Manager Policy Attachments to the user.
---

# alicloud\_resource\_manager\_policy\_attachments

This data source provides the Resource Manager Policy Attachments of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.93.0+.

## Example Usage

```terraform
data "alicloud_resource_manager_policy_attachments" "example" {}

output "first_attachment_id" {
  value = "${data.alicloud_resource_manager_policy_attachments.example.attachments.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `policy_name` - (Optional, ForceNew) The name of the policy. The name must be 1 to 128 characters in length and can contain letters, digits, and hyphens (-).
* `policy_type` - (Optional, ForceNew) The type of the policy. Valid values: `Custom` and `System`.
* `principal_name` - (Optional, ForceNew) The name of the object to which the policy is attached.
* `principal_type` - (Optional, ForceNew) The type of the object to which the policy is attached. If you do not specify this parameter, the system lists all types of objects. Valid values: `IMSUser`: RAM user, `IMSGroup`: RAM user group, `ServiceRole`: RAM role. 
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group or the ID of the Alibaba Cloud account to which the resource group belongs. If you do not specify this parameter, the system lists all policy attachment records under the current account.
* `language` - (Optional, ForceNew) The language that is used to return the description of the system policy. Valid values:`en`: English, `zh-CN`: Chinese, `ja`: Japanese.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Resource Manager Policy Attachment IDs.
* `attachments` - A list of Resource Manager Policy Attachment. Each element contains the following attributes:
    * `id` - The ID of the Resource Manager Policy Attachment.
    * `policy_name`- The name of the policy.
    * `policy_type`- The type of the policy.
    * `principal_name`- The name of the object to which the policy is attached.
    * `principal_type`- The type of the object to which the policy is attached.
    * `description` - The description of the policy.
    * `attach_date` - The time when the policy was attached.
    * `resource_group_id` - The ID of the resource group or the ID of the Alibaba Cloud account to which the resource group belongs.
