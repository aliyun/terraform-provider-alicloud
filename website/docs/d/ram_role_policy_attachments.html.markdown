---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_role_policy_attachments"
sidebar_current: "docs-alicloud-datasource-ram-role-policy-attachments"
description: |-
  Provides a list of Ram Role Policy Attachment owned by an Alibaba Cloud account.
---

# alicloud_ram_role_policy_attachments

This data source provides Ram Role Policy Attachment available to the user.[What is Role Policy Attachment](https://next.api.alibabacloud.com/document/Ram/2015-05-01/AttachPolicyToRole)

-> **NOTE:** Available since v1.248.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ram_role" "role" {
  name        = "roleName"
  document    = <<EOF
    {
      "Statement": [
        {
          "Action": "sts:AssumeRole",
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "apigateway.aliyuncs.com", 
              "ecs.aliyuncs.com"
            ]
          }
        }
      ],
      "Version": "1"
    }
    EOF
  description = "this is a role test."
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ram_policy" "policy" {
  policy_name     = "tf-example-${random_integer.default.result}"
  policy_document = <<EOF
  {
    "Statement": [
      {
        "Action": [
          "oss:ListObjects",
          "oss:GetObject"
        ],
        "Effect": "Allow",
        "Resource": [
          "acs:oss:*:*:mybucket",
          "acs:oss:*:*:mybucket/*"
        ]
      }
    ],
      "Version": "1"
  }
  EOF
  description     = "this is a policy test"
}

resource "alicloud_ram_role_policy_attachment" "default" {
  policy_name = alicloud_ram_policy.policy.policy_name
  policy_type = alicloud_ram_policy.policy.type
  role_name   = alicloud_ram_role.role.name
}

data "alicloud_ram_role_policy_attachments" "default" {
  ids       = ["${alicloud_ram_role_policy_attachment.default.id}"]
  role_name = alicloud_ram_role.role.id
}

output "alicloud_ram_role_policy_attachment_example_id" {
  value = data.alicloud_ram_role_policy_attachments.default.attachments.0.id
}
```

## Argument Reference

The following arguments are supported:
* `role_name` - (Required, ForceNew) The RAM role name.
* `ids` - (Optional, ForceNew, Computed) A list of Role Policy Attachment IDs. The value is formulated as `role:<policy_name>:<policy_type>:<role_name>`.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Role Policy Attachment IDs.
* `attachments` - A list of Role Policy Attachment Entries. Each element contains the following attributes:
  * `attach_date` - The time when the role was attached to the policy.
  * `description` - The policy description.
  * `policy_name` - The name of the policy.
  * `policy_type` - Policy type.- Custom: Custom policy.- System: System policy.
  * `id` - The ID of the resource supplied above. The value is formulated as `role:<policy_name>:<policy_type>:<role_name>`.
