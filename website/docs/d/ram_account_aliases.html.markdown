---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_account_aliases"
sidebar_current: "docs-alicloud-datasource-ram-account-alias"
description: |-
    Provides an alias of the cloud account.
---

# alicloud\_ram\_account\_aliases

The Ram Account Alias data source provides an alias for the Alicloud account.

## Example Usage

```
data "alicloud_ram_account_aliases" "alias" {
  output_file = "alias.txt"
}

```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) The name of file that can save alias data source after running `terraform plan`.

## Attributes Reference

* `account_alias` - Alias of the account.