---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_account_aliases"
sidebar_current: "docs-alicloud-datasource-ram-account-alias"
description: |-
    Provides an alias of the Alibaba Cloud account.
---

# alicloud\_ram\_account\_aliases

This data source provides an alias for the Alibaba Cloud account.

## Example Usage

```
data "alicloud_ram_account_aliases" "alias_ds" {
  output_file = "alias.txt"
}

output "account_alias" {
  value = "${data.alicloud_ram_account_aliases.alias_ds.account_alias}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `account_alias` - Alias of the account.