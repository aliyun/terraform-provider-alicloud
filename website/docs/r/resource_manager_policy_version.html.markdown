---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_policy_version"
sidebar_current: "docs-alicloud-resource-resource-manager-policy-version"
description: |-
  Provides a Alicloud Resource Manager Policy Version resource.
---

# alicloud\_resource\_manager\_policy\_version

Provides a Resource Manager Policy Version resource. 
For information about Resource Manager Policy Version and how to use it, see [What is Resource Manager Policy Version](https://www.alibabacloud.com/help/en/doc-detail/116817.htm).

-> **NOTE:** Available in v1.84.0+.

-> **NOTE:** It is not recommended to use this resource management policy version, it is recommended to directly use the policy resource to manage your policy. Please refer to the link for usage [resource_manager_policy](https://www.terraform.io/docs/providers/alicloud/r/resource_manager_policy).

## Example Usage

Basic Usage

```terraform
resource "alicloud_resource_manager_policy" "example" {
  policy_name     = "tftest"
  policy_document = <<EOF
		{
			"Statement": [{
				"Action": ["oss:*"],
				"Effect": "Allow",
				"Resource": ["acs:oss:*:*:*"]
			}],
			"Version": "1"
		}
    EOF
}

resource "alicloud_resource_manager_policy_version" "example" {
  policy_name     = alicloud_resource_manager_policy.example.policy_name
  policy_document = <<EOF
		{
			"Statement": [{
				"Action": ["oss:*"],
				"Effect": "Allow",
				"Resource": ["acs:oss:*:*:myphotos"]
			}],
			"Version": "1"
		}
    EOF
}
```
## Argument Reference

The following arguments are supported:

* `policy_name` - (Required, ForceNew) The name of the policy. Name must be 1 to 128 characters in length and can contain letters, digits, and hyphens (-).
* `policy_document` - (Required, ForceNew) The content of the policy. The content must be 1 to 2,048 characters in length.
* `is_default_version` - (Optional, Deprecated from version 1.90.0) Specifies whether to set the policy version as the default version. Default to `false`. 

-> **NOTE:** If set to default version, the resource cannot be deleted. You need to set the other version as the default version in policy before you delete this resource.

## Attributes Reference

* `id` - The resource ID of policy version. The value is "`<policy_name>`:`<version_id>`".
* `version_id` - (Removed from version 1.100.0) The ID of the policy version.
* `create_date` - (Removed from version 1.100.0) The time when the policy version was created.

## Import

Resource Manager Policy Version can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_policy_version.example tftest:v2
```
