---
subcategory: "PAIPlugin"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_plugin_user"
description: |-
  This data source provides the PAIPlugin user of the current account.
---

# alicloud_pai_plugin_user

This data source provides the PAIPlugin user of the current account (singleton).

For information about PAIPlugin user and how to use it, see
[What is GetUser](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiplugin-2022-01-12-getuser).

-> **NOTE:** Available since v1.249.0.

## Example Usage

```terraform
data "alicloud_pai_plugin_user" "default" {}

output "account_status" {
  value = data.alicloud_pai_plugin_user.default.account_status
}
```

## Argument Reference

> **NOTE:** This data source has no arguments; it always reflects the current account.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source.
* `account_status` - The account status.
