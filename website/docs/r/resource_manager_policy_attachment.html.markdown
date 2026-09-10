---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_policy_attachment"
sidebar_current: "docs-alicloud-resource-resource-manager-policy-attachment"
description: |-
  Provides a Alicloud Resource Manager Policy Attachment resource.
---

# alicloud_resource_manager_policy_attachment

Provides a Resource Manager Policy Attachment resource to attaches a policy to an object. After you attach a policy to an object, the object has the operation permissions on the current resource group or the resources under the current account. 
For information about Resource Manager Policy Attachment and how to use it, see [How to authorize and manage resource groups](https://www.alibabacloud.com/help/en/doc-detail/94490.htm).

-> **NOTE:** Available since v1.93.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_policy_attachment&exampleId=52bf5397-9f48-6276-f9ab-e05a1ea4a535923c15cb&activeTab=example&spm=docs.r.resource_manager_policy_attachment.0.52bf53979f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexamplename"
}

resource "alicloud_ram_user" "example" {
  name = var.name
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

data "alicloud_resource_manager_resource_groups" "example" {
  status = "OK"
}

# Get Alicloud Account Id
data "alicloud_account" "example" {}

# Attach the custom policy to resource group
resource "alicloud_resource_manager_policy_attachment" "example" {
  policy_name       = alicloud_resource_manager_policy.example.policy_name
  policy_type       = "Custom"
  principal_name    = format("%s@%s.onaliyun.com", alicloud_ram_user.example.name, data.alicloud_account.example.id)
  principal_type    = "IMSUser"
  resource_group_id = data.alicloud_resource_manager_resource_groups.example.ids.0
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_resource_manager_policy_attachment&spm=docs.r.resource_manager_policy_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `policy_name` - (Required, ForceNew) The name of the policy. name must be 1 to 128 characters in length and can contain letters, digits, and hyphens (-).
* `policy_type` - (Required, ForceNew) The type of the policy. Valid values: `Custom`, `System`.
* `principal_name` - (Required, ForceNew) The name of the object to which you want to attach the policy.
* `principal_type` - (Required, ForceNew) The type of the object to which you want to attach the policy. Valid values: `IMSUser`: RAM user, `IMSGroup`: RAM user group, `ServiceRole`: RAM role. 
* `resource_group_id` - (Required, ForceNew) The ID of the resource group or the ID of the Alibaba Cloud account to which the resource group belongs.
    
## Attributes Reference

* `id` - This ID of this resource. It is formatted to `<policy_name>`:`<policy_type>`:`<principal_name>`:`<principal_type>`:`<resource_group_id>`. Before version 1.100.0, the value is `<policy_name>`:`<policy_type>`:`<principal_name>`:`<principal_type>`.

## Import

Resource Manager Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_policy_attachment.example tf-testaccrdpolicy:Custom:tf-testaccrdpolicy@11827252********.onaliyun.com:IMSUser:rg******
```
