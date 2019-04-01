---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_role_policy_attachment"
sidebar_current: "docs-alicloud-resource-ram-role-policy-attachment"
description: |-
  Provides a RAM Role Policy attachment resource.
---

# alicloud\_ram\_role\_policy\_attachment

Provides a RAM Role attachment resource.

## Example Usage

```
# Create a RAM Role Policy attachment.
resource "alicloud_ram_role" "role" {
  name = "test_role"
  ram_users = ["acs:ram::${your_account_id}:root", "acs:ram::${other_account_id}:user/username"]
  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
  description = "this is a role test."
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

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
  role_name = "${alicloud_ram_role.role.name}"
}
```
## Argument Reference

The following arguments are supported:

* `role_name` - (Required, ForceNew) Name of the RAM Role. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-", "_", and must not begin with a hyphen.
* `policy_name` - (Required, ForceNew) Name of the RAM policy. This name can have a string of 1 to 128 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `policy_type` - (Required, ForceNew) Type of the RAM policy. It must be `Custom` or `System`.

## Attributes Reference

The following attributes are exported:

* `id` - The attachment ID.
* `role_name` - The role name.
* `policy_name` - The policy name.
* `policy_type` - The policy type.