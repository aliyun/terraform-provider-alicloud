---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_control_policy_attachment"
sidebar_current: "docs-alicloud-resource-resource-manager-control-policy-attachment"
description: |-
  Provides a Alicloud Resource Manager Control Policy Attachment resource.
---

# alicloud\_resource\_manager\_control\_policy\_attachment

Provides a Resource Manager Control Policy Attachment resource.

For information about Resource Manager Control Policy Attachment and how to use it, see [What is Control Policy Attachment](https://help.aliyun.com/document_detail/208330.html).

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
// Enable the control policy
resource "alicloud_resource_manager_resource_directory" "example" {
  status = "Enabled"
}

resource "alicloud_resource_manager_control_policy" "example" {
  control_policy_name = "tf-testAccName"
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

resource "alicloud_resource_manager_folder" "example" {
  folder_name = "tf-testAccName"
}

resource "alicloud_resource_manager_control_policy_attachment" "example" {
  policy_id  = alicloud_resource_manager_control_policy.example.id
  target_id  = alicloud_resource_manager_folder.example.id
  depends_on = [alicloud_resource_manager_resource_directory.example]
}

```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew) The ID of control policy.
* `target_id` - (Required, ForceNew) The ID of target.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Control Policy Attachment. The value is formatted `<policy_id>:<target_id>`.

## Import

Resource Manager Control Policy Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_control_policy_attachment.example <policy_id>:<target_id>
```
