---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_account_alias"
sidebar_current: "docs-alicloud-resource-ram-account-alias"
description: |-
  Provides a RAM cloud account alias.
---

# alicloud_ram_account_alias

Provides a RAM cloud account alias.

-> **NOTE:** Available since v1.0.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ram_account_alias&exampleId=199509ee-dddc-e02e-194a-e933c48fe92a0ec7aefd&activeTab=example&spm=docs.r.ram_account_alias.0.199509eedd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}
resource "alicloud_ram_account_alias" "alias" {
  account_alias = var.name
}
```
## Argument Reference

The following arguments are supported:

* `account_alias` - (Required, ForceNew) Alias of cloud account. This name can have a string of 3 to 32 characters, must contain only alphanumeric characters or hyphens, such as "-", and must not begin with a hyphen.

## Attributes Reference

The following attributes are exported:

* `id` - The account alias ID, it's set to `account_alias`.

## Import
RAM account alias can be imported using the id, e.g.
```shell
$ terraform import alicloud_ram_account_alias.example my-alias
```
