---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_control_policy_attachment"
description: |-
  Provides a Alicloud Resource Manager Control Policy Attachment resource.
---

# alicloud_resource_manager_control_policy_attachment

Provides a Resource Manager Control Policy Attachment resource.

Control Policy Attachment.

For information about Resource Manager Control Policy Attachment and how to use it, see [What is Control Policy Attachment](https://www.alibabacloud.com/help/en/resource-management/resource-directory/developer-reference/api-resourcemanager-2020-03-31-attachcontrolpolicy).

-> **NOTE:** Available since v1.120.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_resource_manager_control_policy" "default" {
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

resource "alicloud_resource_manager_folder" "default" {
  folder_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_resource_manager_control_policy_attachment" "default" {
  policy_id = alicloud_resource_manager_control_policy.default.id
  target_id = alicloud_resource_manager_folder.default.id
}
```

## Argument Reference

The following arguments are supported:
* `policy_id` - (Required, ForceNew) The ID of the access control policy.
* `target_id` - (Required, ForceNew) The ID of the object from which you want to detach the access control policy. Access control policies can be attached to the following objects:

  - Root folder
  - Subfolders of the Root folder
  - Members

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<policy_id>:<target_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Control Policy Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Control Policy Attachment.

## Import

Resource Manager Control Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_control_policy_attachment.example <policy_id>:<target_id>
```