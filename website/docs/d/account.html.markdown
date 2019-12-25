---
layout: "alicloud"
page_title: "Alicloud: alicloud_account"
sidebar_current: "docs-alicloud-datasource-account"
description: |-
    Provides information about the current Alibaba Cloud account.
---

# alicloud\_account

This data source provides information about the current account.

## Example Usage

```
data "alicloud_account" "current" {
}

output "current_account_id" {
  value = data.alicloud_account.current.id
}
```

## Attributes Reference

The following attributes are exported:

* `id` - Account ID (e.g. "1239306421830812"). It can be used to construct an ARN.
