---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_policy"
sidebar_current: "docs-alicloud-resource-resource-manager-policy"
description: |-
  Provides a Alicloud Resource Manager Policy resource.
---

# alicloud\_resource\_manager\_policy

Provides a Resource Manager Policy resource.  
For information about Resource Manager Policy and how to use it, see [What is Resource Manager Policy](https://www.alibabacloud.com/help/en/doc-detail/93732.htm).

-> **NOTE:** Available in v1.83.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_resource_manager_policy" "example" {
  policy_name     = "abc12345"
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
```
## Argument Reference

The following arguments are supported:

* `policy_name` - (Required, ForceNew) The name of the policy. name must be 1 to 128 characters in length and can contain letters, digits, and hyphens (-).
* `policy_document` - (Required, ForceNew) The content of the policy. The content must be 1 to 2,048 characters in length.
* `description` - (Optional, ForceNew) The description of the policy. The description must be 1 to 1,024 characters in length.
* `default_version` - (Optional, Computed, Deprecated from version 1.90.0) The version of the policy. Default to v1.
    
## Attributes Reference

* `id` - The resource ID of policy. The value is same as `policy_name`.
* `policy_type` - The type of the policy. Valid values: `Custom`, `System`.

## Import

Resource Manager Policy can be imported using the id, e.g.

```
$ terraform import alicloud_resource_policy.example abc12345
```
