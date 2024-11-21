---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_policy_version"
sidebar_current: "docs-alicloud-resource-resource-manager-policy-version"
description: |-
  Provides a Alicloud Resource Manager Policy Version resource.
---

# alicloud_resource_manager_policy_version

Provides a Resource Manager Policy Version resource. 
For information about Resource Manager Policy Version and how to use it, see [What is Resource Manager Policy Version](https://www.alibabacloud.com/help/en/doc-detail/116817.htm).

-> **NOTE:** Available since v1.84.0.

-> **NOTE:** It is not recommended to use this resource management policy version, it is recommended to directly use the policy resource to manage your policy. Please refer to the link for usage [resource_manager_policy](https://www.terraform.io/docs/providers/alicloud/r/resource_manager_policy).

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_policy_version&exampleId=2d0bb1e5-785a-f3aa-3766-97767a31b15b19cc3967&activeTab=example&spm=docs.r.resource_manager_policy_version.0.2d0bb1e578&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}

resource "alicloud_resource_manager_policy" "example" {
  policy_name     = var.name
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
* `version_id` - (Removed from version 1.100.0) The ID of the policy version.
* `create_date` - (Removed from version 1.100.0) The time when the policy version was created.

-> **NOTE:** If set to default version, the resource cannot be deleted. You need to set the other version as the default version in policy before you delete this resource.

## Attributes Reference

* `id` - The resource ID of policy version. The value is "`<policy_name>`:`<version_id>`".

## Import

Resource Manager Policy Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_policy_version.example tftest:v2
```
