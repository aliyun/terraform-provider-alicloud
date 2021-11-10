---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_group_policy_attachment"
sidebar_current: "docs-alicloud-resource-ram-group-policy-attachment"
description: |-
  Provides a RAM Group Policy attachment resource.
---

# alicloud\_ram\_group\_policy\_attachment

Provides a RAM Group Policy attachment resource. 

## Example Usage

```terraform
# Create a RAM Group Policy attachment.
resource "alicloud_ram_group" "group" {
  name     = "groupName"
  comments = "this is a group comments."
  force    = true
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

resource "alicloud_ram_group_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.name
  policy_type = alicloud_ram_policy.policy.type
  group_name  = alicloud_ram_group.group.name
}
```
## Argument Reference

The following arguments are supported:

* `group_name` - (Required, ForceNew) Name of the RAM group. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `policy_name` - (Required, ForceNew) Name of the RAM policy. This name can have a string of 1 to 128 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `policy_type` - (Required, ForceNew) Type of the RAM policy. It must be `Custom` or `System`.

## Attributes Reference

The following attributes are exported:

* `id` - The attachment ID. Composed of policy name, policy type and group name with format `group:<policy_name>:<policy_type>:<group_name>`.

## Import

RAM Group Policy attachment can be imported using the id, e.g.

```
$ terraform import alicloud_ram_group_policy_attachment.example group:my-policy:Custom:my-group
```
