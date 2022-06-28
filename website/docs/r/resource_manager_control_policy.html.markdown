---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_control_policy"
sidebar_current: "docs-alicloud-resource-resource-manager-control-policy"
description: |-
  Provides a Alicloud Resource Manager Control Policy resource.
---

# alicloud\_resource\_manager\_control\_policy

Provides a Resource Manager Control Policy resource.

For information about Resource Manager Control Policy and how to use it, see [What is Control Policy](https://help.aliyun.com/document_detail/208287.html).

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_resource_manager_control_policy" "example" {
  control_policy_name = "tf-testAccRDControlPolicy"
  description         = "tf-testAccRDControlPolicy"
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

* `control_policy_name` - (Required) The name of control policy.
* `description` - (Optional) The description of control policy.
* `effect_scope` - (Required, ForceNew) The effect scope. Valid values `RAM`.
* `policy_document` - (Required) The policy document of control policy.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Control Policy.

## Import

Resource Manager Control Policy can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_control_policy.example <id>
```
