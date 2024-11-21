---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_policy"
sidebar_current: "docs-alicloud-resource-resource-manager-policy"
description: |-
  Provides a Alicloud Resource Manager Policy resource.
---

# alicloud_resource_manager_policy

Provides a Resource Manager Policy resource.  
For information about Resource Manager Policy and how to use it, see [What is Resource Manager Policy](https://www.alibabacloud.com/help/en/doc-detail/93732.htm).

-> **NOTE:** Available since v1.83.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_policy&exampleId=5b76ce81-6d7f-86a6-7167-426dfcb852ecffe686f8&activeTab=example&spm=docs.r.resource_manager_policy.0.5b76ce816d&intl_lang=EN_US" target="_blank">
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

```shell
$ terraform import alicloud_resource_policy.example abc12345
```
