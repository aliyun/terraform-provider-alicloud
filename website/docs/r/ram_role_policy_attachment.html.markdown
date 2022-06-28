---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_role_policy_attachment"
sidebar_current: "docs-alicloud-resource-ram-role-policy-attachment"
description: |-
  Provides a RAM Role Policy attachment resource.
---

# alicloud\_ram\_role\_policy\_attachment

Provides a RAM Role attachment resource.

## Example Usage

```terraform
# Create a RAM Role Policy attachment.
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
  force       = true
}

resource "alicloud_ram_policy" "policy" {
  name        = "policyName"
  document    = <<EOF
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
  description = "this is a policy test"
  force       = true
}

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.name
  policy_type = alicloud_ram_policy.policy.type
  role_name   = alicloud_ram_role.role.name
}
```

## Argument Reference

The following arguments are supported:

* `role_name` - (Required, ForceNew) Name of the RAM Role. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-", "_", and must not begin with a hyphen.
* `policy_name` - (Required, ForceNew) Name of the RAM policy. This name can have a string of 1 to 128 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `policy_type` - (Required, ForceNew) Type of the RAM policy. It must be `Custom` or `System`.

## Attributes Reference

The following attributes are exported:

* `id` - The attachment ID. Composed of policy name, policy type and role name with format `role:<policy_name>:<policy_type>:<role_name>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins, Available in 1.173.0+) Used when creating the RAM Role Policy attachment.
* `delete` - (Defaults to 1 mins, Available in 1.173.0+) Used when deleting the RAM Role Policy attachment.

## Import

RAM Role Policy attachment can be imported using the id, e.g.

```
$ terraform import alicloud_ram_role_policy_attachment.example role:my-policy:Custom:my-role
```
