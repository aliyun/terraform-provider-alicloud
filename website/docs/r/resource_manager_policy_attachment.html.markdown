---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_policy_attachment"
sidebar_current: "docs-alicloud-resource-resource-manager-policy-attachment"
description: |-
  Provides a Alicloud Resource Manager Policy Attachment resource.
---

# alicloud\_resource\_manager\_policy\_attachment

Provides a Resource Manager Policy Attachment resource to attaches a policy to an object. After you attach a policy to an object, the object has the operation permissions on the current resource group or the resources under the current account. 
For information about Resource Manager Policy Attachment and how to use it, see [How to authorize and manage resource groups](https://www.alibabacloud.com/help/en/doc-detail/94490.htm).

-> **NOTE:** Available in v1.93.0+.

## Example Usage

Basic Usage

```terraform
# Create a RAM user
resource "alicloud_ram_user" "example" {
  name = "tf-testaccramuser"
}

# Create a Resource Manager Policy
resource "alicloud_resource_manager_policy" "example" {
  policy_name     = "tf-testaccrdpolicy"
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

# Create a Resource Group
resource "alicloud_resource_manager_resource_group" "example" {
  display_name = "tf_test"
  name         = "tf_test"
}

# Get Alicloud Account Id
data "alicloud_account" "example" {}

# Attach the custom policy to resource group
resource "alicloud_resource_manager_policy_attachment" "example" {
  policy_name       = alicloud_resource_manager_policy.example.policy_name
  policy_type       = "Custom"
  principal_name    = format("%s@%s.onaliyun.com", alicloud_ram_user.example.name, data.alicloud_account.example.id)
  principal_type    = "IMSUser"
  resource_group_id = alicloud_resource_manager_resource_group.example.id
}
```
## Argument Reference

The following arguments are supported:

* `policy_name` - (Required, ForceNew) The name of the policy. name must be 1 to 128 characters in length and can contain letters, digits, and hyphens (-).
* `policy_type` - - (Required, ForceNew) The type of the policy. Valid values: `Custom`, `System`.
* `principal_name` - (Required, ForceNew) The name of the object to which you want to attach the policy.
* `principal_type` - (Required, ForceNew) The type of the object to which you want to attach the policy. Valid values: `IMSUser`: RAM user, `IMSGroup`: RAM user group, `ServiceRole`: RAM role. 
* `resource_group_id` - (Required, ForceNew) The ID of the resource group or the ID of the Alibaba Cloud account to which the resource group belongs.
    
## Attributes Reference

* `id` - This ID of this resource. It is formatted to `<policy_name>`:`<policy_type>`:`<principal_name>`:`<principal_type>`:`<resource_group_id>`. Before version 1.100.0, the value is `<policy_name>`:`<policy_type>`:`<principal_name>`:`<principal_type>`.

## Import

Resource Manager Policy Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_resource_policy_attachment.example tf-testaccrdpolicy:Custom:tf-testaccrdpolicy@11827252********.onaliyun.com:IMSUser:rg******
```
