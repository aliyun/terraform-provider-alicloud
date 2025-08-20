---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_delegate_account"
description: |-
  Provides a Alicloud Cloud SSO Delegate Account resource.
---

# alicloud_cloud_sso_delegate_account

Provides a Cloud SSO Delegate Account resource.

Delegated Administrator Account.

For information about Cloud SSO Delegate Account and how to use it, see [What is Delegate Account](https://next.api.alibabacloud.com/document/cloudsso/2021-05-15/EnableDelegateAccount).

-> **NOTE:** Available since v1.259.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}
resource "alicloud_resource_manager_delegated_administrator" "default" {
  account_id        = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
  service_principal = "cloudsso.aliyuncs.com"
}

resource "alicloud_cloud_sso_delegate_account" "default" {
  account_id = alicloud_resource_manager_delegated_administrator.default.account_id
}
```

## Argument Reference

The following arguments are supported:
* `account_id` - (Required, ForceNew) Delegate administrator account Id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Delegate Account.
* `delete` - (Defaults to 5 mins) Used when delete the Delegate Account.

## Import

Cloud SSO Delegate Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_delegate_account.example <id>
```