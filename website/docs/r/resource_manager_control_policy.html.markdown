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
* `control_policy_name` - (Required) The name of the access control policy.
* `description` - (Optional) The description of the access control policy.
* `effect_scope` - (Required, ForceNew) The effective scope of the access control policy. Valid values:
  - All: The access control policy is in effect for Alibaba Cloud accounts, RAM users, and RAM roles.
  - RAM: The access control policy is in effect only for RAM users and RAM roles.
* `policy_document` - (Required) The document of the access control policy.
* `tags` - (Optional, Map, Available since v1.259.0) The tag of the resource

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - CreateTime

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