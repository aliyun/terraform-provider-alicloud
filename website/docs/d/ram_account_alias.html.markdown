---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_account_alias"
sidebar_current: "docs-alicloud-datasource-ram-account-alias"
description: |-
    Provides an alias of the Alibaba Cloud account.
---

# alicloud_ram_account_alias

This data source provides an alias for the Alibaba Cloud account.

-> **NOTE:** Available since v1.0.0+.

## Example Usage

```terraform
data "alicloud_ram_account_alias" "alias_ds" {
  output_file = "alias.txt"
}

output "account_alias" {
  value = "${data.alicloud_ram_account_alias.alias_ds.account_alias}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `account_alias` - Alias of the account.
