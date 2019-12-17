---
layout: "alicloud"
page_title: "Alicloud: alicloud_account"
sidebar_current: "docs-alicloud-datasource-caller-identity"
description: |-
    Provides the identity of the current Alibaba Cloud user.
---

# alicloud\_caller\_identity

This data source provides the identity of the current user.

-> **NOTE:** Available in 1.65.0+.

## Example Usage

```
data "alicloud_caller_identity" "current" {
}

output "current_user_arn" {
  value = "${data.alicloud_caller_identity.current.id}"
}
```

## Attributes Reference

The following attributes are exported:

* `id` - Principal ID.
* `arn` - The Alibaba Cloud Resource Name (ARN) of the user making the call.
* `account_id` - Account ID.
* `identity_type` - The type of the princiapal. RAMUser for users.
