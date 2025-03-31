---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_group_policy_attachment"
description: |-
  Provides a Alicloud RAM Group Policy Attachment resource.
---

# alicloud_ram_group_policy_attachment

Provides a RAM Group Policy Attachment resource.



For information about RAM Group Policy Attachment and how to use it, see [What is Group Policy Attachment](https://next.api.alibabacloud.com/document/Ram/2015-05-01/AttachPolicyToGroup).

-> **NOTE:** Available since v1.0.0+.

## Example Usage

Basic Usage

```terraform
# Create a RAM Group Policy attachment.
resource "alicloud_ram_group" "group" {
  group_name = "groupName"
  comments   = "this is a group comments."
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

resource "alicloud_ram_group_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.policy.policy_name
  policy_type = alicloud_ram_policy.policy.type
  group_name  = alicloud_ram_group.group.name
}
```

## Argument Reference

The following arguments are supported:
* `group_name` - (Required, ForceNew) The name of the group.
* `policy_name` - (Required, ForceNew) The name of the policy.
* `policy_type` - (Required, ForceNew) Policy type.
  - Custom: Custom policy.
  - System: System policy.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `group:<policy_name>:<policy_type>:<group_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Group Policy Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Group Policy Attachment.

## Import

RAM Group Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_group_policy_attachment.example group:<policy_name>:<policy_type>:<group_name>
```