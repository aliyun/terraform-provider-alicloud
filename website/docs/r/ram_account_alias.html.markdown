---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_account_alias"
sidebar_current: "docs-alicloud-resource-ram-account-alias"
description: |-
  Provides a RAM cloud account alias.
---

# alicloud\_ram\_account\_alias

Provides a RAM cloud account alias.


## Example Usage

```terraform
# Create a alias for cloud account.
resource "alicloud_ram_account_alias" "alias" {
  account_alias = "hallo"
}
```
## Argument Reference

The following arguments are supported:

* `account_alias` - (Required, ForceNew) Alias of cloud account. This name can have a string of 3 to 32 characters, must contain only alphanumeric characters or hyphens, such as "-", and must not begin with a hyphen.

## Attributes Reference

The following attributes are exported:

* `id` - The account alias ID, it's set to `account_alias`.
* `account_alias` - The account alias.

## Import
RAM account alias can be imported using the id, e.g.
```shell
$ terraform import alicloud_ram_account_alias.example my-alias
```
