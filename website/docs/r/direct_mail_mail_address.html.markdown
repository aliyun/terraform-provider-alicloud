---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_mail_address"
sidebar_current: "docs-alicloud-resource-direct-mail-mail-address"
description: |-
  Provides a Alicloud Direct Mail Mail Address resource.
---

# alicloud_direct_mail_mail_address

Provides a Direct Mail Mail Address resource.

For information about Direct Mail Mail Address and how to use it, see [What is Mail Address](https://www.alibabacloud.com/help/en/directmail/latest/set-up-sender-addresses).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_direct_mail_mail_address&exampleId=f0eba849-80c0-1031-2686-81e7401f99d1bc584a05&activeTab=example&spm=docs.r.direct_mail_mail_address.0.f0eba84980&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "account_name" {
  default = "tfexample"
}
variable "domain_name" {
  default = "alicloud-provider.online"
}
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_direct_mail_domain" "example" {
  domain_name = var.domain_name
}
resource "alicloud_direct_mail_mail_address" "example" {
  account_name = format("%s@%s", var.account_name, alicloud_direct_mail_domain.example.domain_name)
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

```shell
$ terraform import alicloud_direct_mail_mail_address.example <id>
```
