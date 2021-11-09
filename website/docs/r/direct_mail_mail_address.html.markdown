---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_mail_address"
sidebar_current: "docs-alicloud-resource-direct-mail-mail-address"
description: |-
  Provides a Alicloud Direct Mail Mail Address resource.
---

# alicloud\_direct\_mail\_mail\_address

Provides a Direct Mail Mail Address resource.

For information about Direct Mail Mail Address and how to use it, see [What is Mail Address](https://www.aliyun.com/product/directmail).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_direct_mail_mail_address" "example" {
  account_name = "example_value@email.com"
  sendtype     = "batch"
}
```

-> **Note:**
A maximum of 10 mailing addresses can be added.
Individual users: Up to 10 mailing addresses can be deleted within a month.
Enterprise users: Up to 10 mailing addresses can be deleted within a month.

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, ForceNew) The sender address. The email address must be filled in the format of account@domain, and only lowercase letters or numbers can be used.
* `password` - (Optional) Account password. The password must be length 10-20 string, contains numbers, uppercase letters, lowercase letters at the same time.
* `reply_address` - (Optional) Return address.
* `sendtype` - (Required, ForceNew) Account type. Valid values: `batch`, `trigger`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Mail Address.
* `status` - Account Status freeze: 1, normal: 0.

## Import

Direct Mail Mail Address can be imported using the id, e.g.

```
$ terraform import alicloud_direct_mail_mail_address.example <id>
```
