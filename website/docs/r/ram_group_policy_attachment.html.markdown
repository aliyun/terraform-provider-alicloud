---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_group_policy_attachment"
sidebar_current: "docs-alicloud-resource-ram-group-policy-attachment"
description: |-
  Provides a RAM Group Policy attachment resource.
---

# alicloud\_ram\_group\_policy\_attachment

Provides a RAM Group Policy attachment resource. 

## Example Usage

```
# Create a RAM Group Policy attachment.
resource "alicloud_ram_group" "group" {
  name = "test_group"
  comments = "this is a group comments."
  force = true
}

resource "alicloud_ram_policy" "policy" {
  name = "test_policy"
  statement = [
      {
        effect = "Allow"
        action = [
          "oss:ListObjects",
          "oss:GetObject"]
        resource = [
          "acs:oss:*:*:mybucket",
          "acs:oss:*:*:mybucket/*"]
      }]
  description = "this is a policy test"
  force = true
}

resource "alicloud_ram_group_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
  group_name = "${alicloud_ram_group.group.name}"
}
```
## Argument Reference

The following arguments are supported:

* `group_name` - (Required, ForceNew) Name of the RAM group. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `policy_name` - (Required, ForceNew) Name of the RAM policy. This name can have a string of 1 to 128 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `policy_type` - (Required, ForceNew) Type of the RAM policy. It must be `Custom` or `System`.

## Attributes Reference

The following attributes are exported:

* `id` - The attachment ID.
* `group_name` - The group name.
* `policy_name` - The policy name.
* `policy_type` - The policy type.