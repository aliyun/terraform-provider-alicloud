---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_user_policy_attachment"
sidebar_current: "docs-alicloud-resource-ram-user-policy-attachment"
description: |-
  Provides a RAM User Policy attachment resource.
---

# alicloud\_ram\_user\_policy\_attachment

Provides a RAM User Policy attachment resource. 

## Example Usage

```
# Create a RAM User Policy attachment.
resource "alicloud_ram_user" "user" {
  name = "userName"
  display_name = "user_display_name"
  mobile = "86-18688888888"
  email = "hello.uuu@aaa.com"
  comments = "yoyoyo"
  force = true
}

resource "alicloud_ram_policy" "policy" {
  name = "policyName"
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

resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
  user_name = "${alicloud_ram_user.user.name}"
}
```
## Argument Reference

The following arguments are supported:

* `user_name` - (Required, ForceNew) Name of the RAM user. This name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin with a hyphen.
* `policy_name` - (Required, ForceNew) Name of the RAM policy. This name can have a string of 1 to 128 characters, must contain only alphanumeric characters or hyphen "-", and must not begin with a hyphen.
* `policy_type` - (Required, ForceNew) Type of the RAM policy. It must be `Custom` or `System`.

## Attributes Reference

The following attributes are exported:

* `id` - The attachment ID.
* `user_name` - The user name.
* `policy_name` - The policy name.
* `policy_type` - The policy type.