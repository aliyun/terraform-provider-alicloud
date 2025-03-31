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

-> **NOTE:** Available since v1.0.0+.

## Example Usage

Basic Usage

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account Alias.
* `update` - (Defaults to 5 mins) Used when update the Account Alias.

## Import

RAM Account Alias can be imported using the id, e.g.

```shell
$ terraform import alicloud_ram_account_alias.example <id>
```