---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_account_alias"
description: |-
  Provides a Alicloud RAM Account Alias resource.
---

# alicloud_ram_account_alias

Provides a RAM Account Alias resource.



For information about RAM Account Alias and how to use it, see [What is Account Alias](https://next.api.alibabacloud.com/document/Ram/2015-05-01/SetAccountAlias).

-> **NOTE:** Available since v1.0.0.

## Example Usage

Basic Usage

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

### Deleting `alicloud_ram_account_alias` or removing it from your configuration

Terraform cannot destroy resource `alicloud_ram_account_alias`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `account_alias` - (Required) The alias of the account.
It can be 3 to 32 characters in length and can contain lowercase letters, digits, and dashes (-).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. This field is set to your Alibaba Cloud Account ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account Alias.
* `update` - (Defaults to 5 mins) Used when update the Account Alias.

## Import

RAM Account Alias can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_account_alias.example <id>
```