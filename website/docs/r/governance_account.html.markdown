---
subcategory: "Governance"
layout: "alicloud"
page_title: "Alicloud: alicloud_governance_account"
description: |-
  Provides a Alicloud Governance Account resource.
---

# alicloud_governance_account

Provides a Governance Account resource.

Member account created by the Cloud Governance Center account factory.

For information about Governance Account and how to use it, see [What is Account](https://next.api.aliyun.com/document/governance/2021-01-20/EnrollAccount).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_account" "default" {
}

data "alicloud_governance_baselines" "default" {
}

data "alicloud_resource_manager_folders" "default" {
}

resource "alicloud_governance_account" "default" {
  account_name_prefix = "${var.name}-${random_integer.default.result}"
  folder_id           = data.alicloud_resource_manager_folders.default.ids.0
  baseline_id         = data.alicloud_governance_baselines.default.ids.0
  payer_account_id    = data.alicloud_account.default.id
  display_name        = "${var.name}-${random_integer.default.result}"
}
```

### Deleting `alicloud_governance_account` or removing it from your configuration

Terraform cannot destroy resource `alicloud_governance_account`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `account_id` - (Optional, ForceNew, Computed, Int) The ID of the enrolled account.
  - If you are creating a new resource account, this parameter is not required.
  - If you are enrolling a existing account to account factory, this parameter is required.
* `account_name_prefix` - (Optional) Account name prefix.
  - This parameter is required if you are creating a new resource account.
  - If the registration application is applied to an existing account, this parameter does not need to be filled in.
* `account_tags` - (Optional, ForceNew, List, Available since v1.234.0) The tags of the account See [`account_tags`](#account_tags) below.
* `baseline_id` - (Required) The baseline ID.

  If it is left blank, the system default baseline is used by default.
* `display_name` - (Optional) The account display name.
  - This parameter is required if you are creating a new resource account.
  - If the registration application is applied to an existing account, this parameter does not need to be filled in.
* `folder_id` - (Optional) The ID of the parent resource folder.

  If you want to create a new resource account and leave this parameter blank, the account is created in the Root folder by default.

  If the registration application is applied to an existing account, this parameter does not need to be filled in.
* `payer_account_id` - (Optional, Int) The ID of the billing account. If you leave this parameter empty, the current account is used as the billing account.
* `default_domain_name` - (Optional, Available since v1.231.0) The domain name is used to qualify the login name of RAM users and RAM roles.

### `account_tags`

The account_tags supports the following:
* `tag_key` - (Optional, ForceNew, Available since v1.234.0) The key of the tags
* `tag_value` - (Optional, ForceNew, Available since v1.234.0) The value of the tags

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - Account registration status. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Account.
* `update` - (Defaults to 5 mins) Used when update the Account.

## Import

Governance Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_governance_account.example <id>
```