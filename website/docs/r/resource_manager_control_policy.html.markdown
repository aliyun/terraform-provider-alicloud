---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_control_policy"
description: |-
  Provides a Alicloud Resource Manager Control Policy resource.
---

# alicloud_resource_manager_control_policy

Provides a Resource Manager Control Policy resource.



For information about Resource Manager Control Policy and how to use it, see [What is Control Policy](https://www.alibabacloud.com/help/en/resource-management/latest/api-resourcedirectorymaster-2022-04-19-createcontrolpolicy).

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_control_policy&exampleId=4853b74e-7473-f7d9-ae4d-554d6448cab81e064016&activeTab=example&spm=docs.r.resource_manager_control_policy.0.4853b74e74&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_resource_manager_control_policy" "example" {
  control_policy_name = var.name
  description         = var.name
  effect_scope        = "RAM"
  policy_document     = <<EOF
  {
    "Version": "1",
    "Statement": [
      {
        "Effect": "Deny",
        "Action": [
          "ram:UpdateRole",
          "ram:DeleteRole",
          "ram:AttachPolicyToRole",
          "ram:DetachPolicyFromRole"
        ],
        "Resource": "acs:ram:*:*:role/ResourceDirectoryAccountAccessRole"
      }
    ]
  }
  EOF
}

```

## Argument Reference

The following arguments are supported:
* `control_policy_name` - (Required) The new name of the access control policy.
The name must be 1 to 128 characters in length. The name can contain letters, digits, and hyphens (-) and must start with a letter.
* `description` - (Optional) The new description of the access control policy.
The description must be 1 to 1,024 characters in length. The description can contain letters, digits, underscores (\_), and hyphens (-) and must start with a letter.
* `effect_scope` - (Required, ForceNew) The effective scope of the access control policy. Valid values:

  - All: The access control policy is in effect for Alibaba Cloud accounts, RAM users, and RAM roles.
  - RAM: The access control policy is in effect only for RAM users and RAM roles.
* `policy_document` - (Required, JsonString) The new document of the access control policy.
The document can be a maximum of 4,096 characters in length.
For more information about the languages of access control policies, see [Languages of access control policies](https://www.alibabacloud.com/help/en/doc-detail/179096.html).
For more information about the examples of access control policies, see [Examples of custom access control policies](https://www.alibabacloud.com/help/en/doc-detail/181474.html).
* `tags` - (Optional, Map, Available since v1.260.1) The tags.
You can specify a maximum of 20 tags.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the access control policy was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Control Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Control Policy.
* `update` - (Defaults to 5 mins) Used when update the Control Policy.

## Import

Resource Manager Control Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_control_policy.example <id>
```